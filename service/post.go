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

// GetPostListNew  将两个查询帖子列表逻辑合二为一的函数,对请求逻辑进行分发
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 根据请求参数的不同，执行不同的逻辑。
	if p.CommunityID == 0 {
		// 查所有
		data, err = GetPostList2(p)
	} else {
		// 根据社区id查询
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}

// GetPostList 根据page,size调用dao层的函数从mysql中获取帖子列表
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

// GetPostList2
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2. 去redis查post_id列表
	ids, err := redis.GetPostIdsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIdsInOrder return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))

	// 提前取得每个帖子的投票赞成数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 3. 根据post_id列表从mysql中获取数据
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))

	//将帖子的作者和所属社区填充
	for idx, post := range posts {
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
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: com,
		}
		data = append(data, postDetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2. 去redis查post_id列表
	ids, err := redis.GetCommunityPostIdsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIdsInOrder return 0 data")
		return
	}
	zap.L().Debug("GetCommunityPostList", zap.Any("ids", ids))

	// 提前取得每个帖子的投票赞成数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 3. 根据post_id列表从mysql中获取数据
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		return
	}

	data = make([]*models.ApiPostDetail, 0, len(posts))

	//将帖子的作者和所属社区填充
	for idx, post := range posts {
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
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: com,
		}
		data = append(data, postDetail)
	}
	return
}
