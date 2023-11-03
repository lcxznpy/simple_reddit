package controllers

import (
	"goweb/models"
	"goweb/service"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

//帖子投票

func PostVoteHandler(c *gin.Context) {
	// 1. 参数校验和绑定
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		errs, ok := err.(validator.ValidationErrors) //类型断言
		if !ok {
			ResponseError(c, CodeInvalidParam) //传的参数类型错了
			return
		}
		//传的参数类型没错，但是不符合规定
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData) //翻译并去除后端的结构体信息
		return
	}
	//获取用户id
	userId, err := GetCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 2. 业务处理
	if err := service.VoteForPost(userId, p); err != nil {
		zap.L().Error("service.VoteForPost() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, 123123)
}
