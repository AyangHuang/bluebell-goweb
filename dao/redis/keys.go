package redis

// redis key注意使用命名空间的方式,方便查询和拆分

const (
	Prefix           = "bluebell:"   // 项目key前缀
	KeyPostTimeZSet  = "post:time"   // zset;贴子及发帖时间
	KeyPostScoreZSet = "post:score"  // zset;贴子及投票的总分数，可负
	KeyPostVotedZSet = "post:voted:" // zset;记录用户及投票类型;参数是 post id; source是投票类型，可用zcout 快速计算赞成票的个数
	// 1 表示赞成票，-1 表示反对票，0 表示无投票
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}
