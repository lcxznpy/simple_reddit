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
