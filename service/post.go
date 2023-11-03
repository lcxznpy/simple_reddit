package service

import (
	"goweb/dao/mysql"
	"goweb/dao/redis"
	"goweb/models"
	"goweb/pkg/snowflake"

	"go.uber.org/zap"
)

// todo 如何保证写mysql和写redis都一起成功
// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	// 1. 生成post_id
	p.Id = int64(snowflake.GenID())
	// 2. 保存数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return
	}
	err = redis.CreatePost(p.Id, p.CommunityId)
	return

}

// GetPostById 根据帖子id查询帖子详情
func GetPostById(pid int64) (data *models.ApiPostDetail, err error) {
	data = new(models.ApiPostDetail)
	post, err := mysql.GetPostDetailById(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailById(pid) failed",
			zap.Int64("pid", pid),
			zap.Error(err))
		return
	}

	user, err := mysql.GetUserById(post.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorId) failed",
			zap.Int64("author_id", post.AuthorId),
			zap.Error(err))
		return
	}
	com, err := mysql.GetCommunityDetailById(post.CommunityId)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailById(post.CommunityId) failed",
			zap.Int64("community_id", post.CommunityId),
			zap.Error(err))
		return
	}
	data.AuthorName = user.UserName
	data.Post = post
	data.CommunityDetail = com
	return
}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorId) failed",
				zap.Int64("author_id", post.AuthorId),
				zap.Error(err))
			continue
		}
		com, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailById(post.CommunityId) failed",
				zap.Int64("community_id", post.CommunityId),
				zap.Error(err))
			continue
		}
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			Post:            post,
			CommunityDetail: com,
		}
		data = append(data, postDetail)
	}
	return
}
