package models

//定义请求体的结构体参数

// ParamSignUp 用户注册信息
type ParamSignUp struct {
	UserName   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录参数
type ParamLogin struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData投票数据
type ParamVoteData struct {
	//UserId   直接从c上下文获取
	PostId    string `json:"post_id" binding:"required"`                       //帖子id
	Direction int8   `json:"direction,string" binding:"required,oneof=1 0 -1"` //赞成(1)or反对(-1)or取消投票(0)
}
