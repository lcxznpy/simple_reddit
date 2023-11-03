package mysql

import (
	"goweb/models"
	"strings"

	"github.com/jmoiron/sqlx"
)

// CreatePost  数据库内创建帖子
func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post(
                post_id,title,content,author_id,community_id)
				values (?,?,?,?,?)`

	_, err = db.Exec(sqlStr, p.Id, p.Title, p.Content, p.AuthorId, p.CommunityId)
	return
}

// GetPostDetailById 通过id获取帖子详情
func GetPostDetailById(id int64) (data *models.Post, err error) {
	data = new(models.Post)
	sqlStr := `select 
		post_id,author_id,community_id,title,content,create_time
		from post 
		where post_id = ?`

	err = db.Get(data, sqlStr, id)
	return
}

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
		from post
		order by create_time desc 
		limit ?,?`

	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}

func GetPostListByIds(ids []string) (postList []*models.Post, err error) {
	//通过find_in_set保证id的顺序与ids的顺序一致
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
			from post
			where post_id in (?)
			order by FIND_IN_SET(post_id,?)   
			`

	//返回被修改过的查询语句和参数，方便数据库执行
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}
