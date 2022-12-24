package redis

import (
	"github.com/go-redis/redis"
	"time"
)

func CreatePost(postID int64) error {
	// 事务
	pipeline := rdb.TxPipeline()
	// 帖子发布时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postID,
	})

	// 帖子投票总数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  0,
		Member: postID,
	})

	_, err := pipeline.Exec()
	return err
}
