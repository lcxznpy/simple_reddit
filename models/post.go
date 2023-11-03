package models

import "time"

// todo : b站内存对齐视频
type Post struct {
	//多加一个string类型，可以保证前端传json数据来的时候，先转成string，再转成int64，
	Id          int64     `json:"id,string" db:"post_id"`
	AuthorId    int64     `json:"author_id" db:"author_id"`
	CommunityId int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

// 帖子详情接口
type ApiPostDetail struct {
	AuthorName       string             `json:"author_name"`
	VoteNum          int64              `json:"vote_num"`
	*Post                               //嵌入post信息
	*CommunityDetail `json:"community"` //嵌入社区信息
}
