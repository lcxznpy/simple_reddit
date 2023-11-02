package controllers

import (
	"goweb/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区相关

// CommunityHandler 获取社区列表
func CommunityHandler(c *gin.Context) {
	//查询所有社区的(com_id,com_name)以列表形式返回
	data, err := service.GetCommunityList()
	if err != nil {
		zap.L().Error("service.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) //不要把服务端的错误暴露给外部
		return
	}
	ResponseSuccess(c, data)

}

// CommunityDetailHandler 获取社区详情
func CommunityDetailHandler(c *gin.Context) {
	// 1. 获取社区id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//	2.根据id获取社区详情
	data, err := service.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("service.GetCommunityDetail() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)

}
