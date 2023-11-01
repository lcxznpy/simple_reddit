package controllers

import (
	"fmt"
	"goweb/models"
	"goweb/service"
	"net/http"

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
			c.JSON(http.StatusOK, gin.H{
				"msg": "请求参数有误",
			})
			return
		}
		//是参数类型错误
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)),
		})
		return

	}
	fmt.Println(p)
	//2. 业务处理
	if err := service.SignUp(p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册失败",
		})
		return
	}
	//3. 返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func LoginHandler(c *gin.Context) {

}
