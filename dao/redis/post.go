package redis

import (
	"goweb/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// GetPostIdsInOrder 从redis中获取id
func GetPostIdsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis中获取帖子id
	// 1. 根据order确定要查询的key
	key := getKeyString(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getKeyString(KeyPostScoreZSet)
	}
	//
	return getIdsFromKey(key, p.Page, p.Size)
}

// getIdsFromKey 从redis中按规则读取数据
func getIdsFromKey(key string, page, size int64) ([]string, error) {
	// 2. 确定索引位置
	start := (page - 1) * size
	end := start + size - 1
	//3. 按照key从大到小查询 获取post_id列表
	return client.ZRevRange(ctx, key, start, end).Result()
}

// GetPostVoteData  根据ids获取每个帖子的赞成数
func GetPostVoteData(ids []string) (data []int64, err error) {
	//通过pipeline一次性发送多条命令，减少RTT网络交互
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getKeyString(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val() //断言一下
		data = append(data, v)
	}
	return
}

// GetCommunityPostIdsInOrder 按社区查询ids
func GetCommunityPostIdsInOrder(p *models.ParamPostList) ([]string, error) {
	//order规则的key
	orderKey := getKeyString(KeyPostTimeZSet)
	if p.Order == KeyPostScoreZSet {
		orderKey = getKeyString(KeyPostScoreZSet)
	}
	//社区key
	cKey := getKeyString(KeyCommunitySetPF + strconv.Itoa(int(p.CommunityID)))

	//通过ZInterStore生成的key,利用缓存key减少ZInterStore的执行次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	//key不存在,执行zinterstore
	if client.Exists(ctx, key).Val() < 1 {
		pipeline := client.Pipeline()
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Keys:      []string{cKey, key},
			Aggregate: "MAX",
		})
		pipeline.Expire(ctx, key, 60*time.Second)
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	//如果60s内，key存在，直接找就行，不用执行zinterstore
	return getIdsFromKey(key, p.Page, p.Size)
}
