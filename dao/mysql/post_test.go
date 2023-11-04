package mysql

import (
	"goweb/models"
	"goweb/settings"
	"testing"
)

func init() {
	dbCfg := settings.MySQLConfig{
		Host:         "47.102.119.88",
		User:         "root",
		Password:     "dhxdl666",
		DB:           "sql_demo",
		Port:         3306,
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		Id:          10,
		AuthorId:    123,
		CommunityId: 1,
		Title:       "test",
		Content:     "just a test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert into db failed ,err:%v\n ", err)
	}
	t.Log("CreatePost insert into db success")
}
