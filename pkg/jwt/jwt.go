package jwt

import (
	"errors"
	"time"

	"github.com/spf13/viper"

	"github.com/golang-jwt/jwt/v4"
)

// CustomSecret 用于加盐的字符串
var mySecret = []byte("dhxdl666")

type MyClaims struct {
	// 可根据需要自行添加字段
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`

	jwt.RegisteredClaims // 内嵌标准的声明
}

// GenToken 生成JWT
func GenToken(userID int64, username string) (string, error) {
	// 创建一个我们自己的声明的数据
	c := MyClaims{
		userID,
		"username", // 自定义字段
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(viper.GetInt("auth.jwt_expire")) * time.Hour)),
			Issuer:    "dhxdl666", // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid { // 校验token,正常token
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
