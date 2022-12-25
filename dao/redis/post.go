package redis

import (
	"github.com/go-redis/redis"
	"strconv"
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

// GetPostTime 获取发布时间
func GetPostTime(postID string) (float64, error) {
	postTime := rdb.ZScore(getRedisKey(KeyPostTimeZSet), postID).Val()
	if postTime == 0 {
		return 0, ErrorPostNotExist
	}
	return postTime, nil
}

// GetPostVotes 获取票数
func GetPostVotes(postID string) float64 {
	return rdb.ZScore(getRedisKey(KeyPostScoreZSet), postID).Val()
}

// GetPostUserVoted 获取该帖子指定用户的投票情况
func GetPostUserVoted(postID string, userID int64) float64 {
	userIDStr := strconv.FormatInt(userID, 10)
	return rdb.ZScore(getRedisKey(KeyPostVotedZSet)+postID, userIDStr).Val()
}

// PostVote 为帖子投票，事务
func PostVote(userID int64, postID string, votes float64, direction int8) error {
	userIDStr := strconv.FormatInt(userID, 10)
	// 事务
	pipeline := rdb.TxPipeline()
	// 记录票数
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), votes, postID)
	// 记录帖子投票类型
	pipeline.ZAdd(getRedisKey(KeyPostVotedZSet)+postID, redis.Z{
		Score:  float64(direction),
		Member: userIDStr,
	})
	_, err := pipeline.Exec()
	return err
}
