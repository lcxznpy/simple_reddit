package controllers

import (
	"goweb/models"
	"goweb/service"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	// 1.获取参数和参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON(p) failed", zap.Any("err", err))
		zap.L().Error("create post with invalid param")
		ResponseError(c, CodeInvalidParam)
		return
	}
	//获取当前用户id
	userId, err := GetCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	p.AuthorId = userId
	// 2.创建帖子
	if err := service.CreatePost(p); err != nil {
		zap.L().Error("service.CreatePost(p) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, CodeSuccess)
}

// GetPostDetailHandler 获取帖子详情
func GetPostDetailHandler(c *gin.Context) {
	// 1. 获取帖子id
	pidStr := c.Param("id")
	id, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 2. 根据id查询帖子
	data, err := service.GetPostById(id)
	if err != nil {
		zap.L().Error("service.GetPostById(pid) failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, data)
}

// GetPostListHandler 通过mysql获取帖子列表
func GetPostListHandler(c *gin.Context) {

	// 1. 获取分页参数
	page, size := getPageInfo(c)
	// 2. service获取数据
	date, err := service.GetPostList(page, size)
	if err != nil {
		zap.L().Error("service.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 2. 返回响应
	ResponseSuccess(c, date)
}

// 获取分页参数
func getPageInfo(c *gin.Context) (int64, int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	var (
		page int64
		size int64
		err  error
	)

	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}

// GetPostListHandler2  根据前端传来的参数动态获取帖子列表   参数:创建时间、帖子评分
func GetPostListHandler2(c *gin.Context) {
	// 1.获取参数
	//get请求从  url里面获取参数 /api/v1/post2?page=2&size=1&order=time
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime, //magic string
	}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	// 2.从redis中获取id列表

	date, err := service.GetPostListNew(p)
	if err != nil {
		zap.L().Error("service.GetPostList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3.根据id列表从数据库中查帖子数据
	// 2. 返回响应
	ResponseSuccess(c, date)
}

//func GetCommunityPostListHandler(c *gin.Context) {
//	// 1.获取参数
//	//get请求从  url里面获取参数 /api/v1/post2?page=2&size=1&order=time
//	p := &models.ParamCommunityPostList{
//		Page:  1,
//		Size:  10,
//		Order: models.OrderTime, //magic string
//	}
//	if err := c.ShouldBindQuery(p); err != nil {
//		zap.L().Error("GetCommunityPostListHandler with invalid param", zap.Error(err))
//		ResponseError(c, CodeInvalidParam)
//		return
//	}
//
//	// 2.从redis中获取id列表
//
//	date, err := service.GetCommunityPostList(p)
//	if err != nil {
//		zap.L().Error("service.GetCommunityPostList(p) failed", zap.Error(err))
//		ResponseError(c, CodeServerBusy)
//		return
//	}
//
//	// 3.根据id列表从数据库中查帖子数据
//	// 2. 返回响应
//	ResponseSuccess(c, date)
//}
