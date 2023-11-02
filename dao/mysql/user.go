package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"goweb/models"
)

const secret = "dhxdl666"

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

func InsertUser(user *models.User) error {
	//对密码加密
	user.Password = encryptPassword(user.Password)
	//保存数据库
	sqlstr := `insert into user(user_id,username,password) 
	values(?,?,?)`
	_, err := db.Exec(sqlstr, user.UserId, user.UserName, user.Password)
	return err
}

// encryptPassword 密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func LoginByUserName(login *models.User) error {
	password := encryptPassword(login.Password)
	sqlStr := `select user_id,username,password from user where username = ?`
	err := db.Get(login, sqlStr, login.UserName)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	// 数据库错误
	if err != nil {
		return err
	}
	if password != login.Password {
		return ErrorInvalidPassword
	}
	return nil
}

func GetUserById(id int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id,username from user where user_id = ?`
	err = db.Get(user, sqlStr, id)
	return
}
