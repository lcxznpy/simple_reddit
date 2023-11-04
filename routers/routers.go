package routers

import (
	"goweb/controllers"
	_ "goweb/docs" // 千万不要忘了导入把你上一步生成的docs
	"goweb/logger"
	"goweb/middlewares"
	"net/http"

	"github.com/gin-contrib/pprof"

	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	gs "github.com/swaggo/gin-swagger"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	r.LoadHTMLFiles("./index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	//middlewares.RateLimitMiddleware(time.Second*10, 1)  每10s放一个令牌
	//r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(time.Second*2, 1))
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")

	v1.GET("/posts", controllers.GetPostListHandler) //获取帖子列表

	//登录
	v1.POST("/login", controllers.LoginHandler)
	//注册
	v1.POST("/signup", controllers.SignupHandler)
	v1.Use(middlewares.JWTAuthMiddleware())

	{
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)

		v1.POST("/post", controllers.CreatePostHandler)       //创建帖子
		v1.GET("/post/:id", controllers.GetPostDetailHandler) //获取帖子详细信息

		v1.GET("/posts2", controllers.GetPostListHandler2) //根据时间或分数获取帖子列表

		v1.POST("/vote", controllers.PostVoteHandler)
	}
	//为性能分析pprof注册路由
	pprof.Register(r)
	return r
}
