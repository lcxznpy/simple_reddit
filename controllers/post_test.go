package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)

	reqBody := `{
		"community_id":1,
		"title":"test",
		"content":"test test test test"
	}`
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(reqBody)))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, 200)
	//方法1  判断响应体中是否包含 "需要登录"
	//assert.Contains(t, w.Body.String(), "需要登录")

	//方法2  反序列化响应体，判断rescode是否与预期一致
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json unmarshal w.body failed,err:%v\n", err)
	}
	assert.Equal(t, res.Code, CodeNeedLogin)

}
