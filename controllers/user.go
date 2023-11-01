package controllers

import (
	"errors"
	"fmt"
	"goweb/dao/mysql"
	"goweb/models"
	"goweb/service"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func SignupHandler(c *gin.Context) {
	//1. 获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("signup with invalid param", zap.Error(err))
		//判断是不是参数类型的错误
		errs, ok := err.(validator.ValidationErrors)
		//不是参数类型错误
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		//是参数类型错误
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return

	}
	fmt.Println(p)
	//2. 业务处理
	if err := service.SignUp(p); err != nil {
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3. 返回响应
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("login with invalid param", zap.String("username", p.UserName), zap.Error(err))
		//判断是不是参数类型的错误
		errs, ok := err.(validator.ValidationErrors)
		//不是参数类型错误
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		//是参数类型错误
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	token, err := service.Login(p)
	if err != nil {
		zap.L().Error("login failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	ResponseSuccess(c, token)

}
