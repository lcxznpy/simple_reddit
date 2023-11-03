package redis

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"go.uber.org/zap"

	"github.com/redis/go-redis/v9"
)

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 //每一票的分数

)

var (
	ErrVoteTimeExpired = errors.New("该帖子投票时间已过")
	ErrVoteRepeated    = errors.New("不允许重复投票")
	ctx                = context.Background()
)

// CreatePost 在redis里面创建帖子的信息，帖子的创建时间 和 帖子的分数 和 社区中的帖子
func CreatePost(postId, communityId int64) error {
	pipeline := client.TxPipeline()

	// 1. 添加帖子创建时间记录
	pipeline.ZAdd(ctx, getKeyString(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})
	// 2. 帖子分数
	pipeline.ZAdd(ctx, getKeyString(KeyPostScoreZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postId,
	})
	// 3. 把帖子加到社区的set中
	comKey := getKeyString(KeyCommunitySetPF + strconv.Itoa(int(communityId)))
	pipeline.SAdd(ctx, comKey, postId)
	_, err := pipeline.Exec(ctx)
	if err != nil {
		zap.L().Error("redis error", zap.Error(err))
	}
	return err
}

func VoteForPost(userId, postId string, value float64) error {
	// 1. 判断投票限制
	//获取帖子时间，判断还能不能投票
	postTime := client.ZScore(ctx, getKeyString(KeyPostTimeZSet), postId).Val()
	if float64(time.Now().Unix())-postTime > oneWeekInSeconds {
		fmt.Println(float64(time.Now().Unix()))
		fmt.Println(postTime)
		fmt.Println(float64(time.Now().Unix()) - postTime)
		return ErrVoteTimeExpired
	}
	// 2 和 3 要开启pipeline事务
	// 2. 更新帖子分数
	//获取当前用户给帖子的投票记录
	oldv := client.ZScore(ctx, getKeyString(KeyPostVotedZSetPrefix+postId), userId).Val()
	var op float64
	//不允许重复投票
	if oldv == value {
		return ErrVoteRepeated
	}
	if value > oldv {
		op = 1
	} else {
		op = -1
	}
	diff := math.Abs(oldv - value)
	pipeline := client.TxPipeline()
	pipeline.ZIncrBy(ctx, getKeyString(KeyPostScoreZSet), op*diff*scorePerVote, postId).Result()

	// 3. 记录当前用户为该帖子投票的数据
	if value == 0 {
		//取消投票，删除该帖子表中用户投票记录
		pipeline.ZRem(ctx, getKeyString(KeyPostVotedZSetPrefix+postId), userId)
	} else {
		//添加最新的用户投票记录，会覆盖之前的记录
		pipeline.ZAdd(ctx, getKeyString(KeyPostVotedZSetPrefix+postId), redis.Z{
			Score:  value,
			Member: userId,
		})
	}
	_, err := pipeline.Exec(ctx)
	return err
}
