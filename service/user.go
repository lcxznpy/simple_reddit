package service

import (
	"goweb/dao/mysql"
	"goweb/models"
	"goweb/pkg/snowflake"
)

// 用户注册
func SignUp(p *models.ParamSignUp) error {
	// 1.判断是否存在
	if err := mysql.CheckUserExist(p.UserName); err != nil {
		return err
	}

	// 2.生成uid
	userId := snowflake.GenID()
	//构造user实例
	user := &models.User{
		UserId:   userId,
		UserName: p.UserName,
		Password: p.Password,
	}
	return mysql.InsertUser(user)
}

func Login(p *models.ParamLogin) error {
	user := &models.User{
		UserName: p.UserName,
		Password: p.Password,
	}

	if err := mysql.LoginByUserName(user); err != nil {
		return err
	}
	return nil
}
