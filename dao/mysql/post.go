package mysql

import "goweb/models"

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

func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id,title,content,author_id,community_id,create_time
		from post
		limit ?,?`

	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return
}
