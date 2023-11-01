package routers

import (
	"goweb/logger"
	"goweb/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, settings.Conf.Version)
	})
	return r
}
