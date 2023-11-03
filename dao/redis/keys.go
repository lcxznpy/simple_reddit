package redis

//存 redis里的key

const (
	KeyPrefix              = "reddit"
	KeyPostTimeZSet        = "post:time"   //记录帖子id(value)和发帖时间(时间戳为score)
	KeyPostScoreZSet       = "post:score"  //记录帖子id(value)和分数(score)
	KeyPostVotedZSetPrefix = "post:voted:" //记录每个帖子的用户投票情况 用户id(value)和投票情况(score 1,-1)  key(post:voted:帖子id)

	KeyCommunitySetPF = "community:" // set;保存每个分区下帖子的id
)

// getKeyString 给rediskey加前缀
func getKeyString(key string) string {
	return KeyPrefix + key
}
