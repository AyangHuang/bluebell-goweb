package redis

import (
	"github.com/go-redis/redis"
	"strconv"
	"time"
)

func CreatePost(postID int64) error {
	postIDStr := strconv.FormatInt(postID, 10)
	// 事务
	pipeline := rdb.TxPipeline()
	// 帖子发布时间
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: postIDStr,
	})

	// 帖子投票总数
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  0,
		Member: postIDStr,
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

// GetVotes 获取按时间排序 [start, stop) 的帖子
func GetVotes(start, stop int64, s *[]redis.Z) {
	*s = rdb.ZRangeWithScores(getRedisKey(KeyPostTimeZSet), start, stop-1).Val()
}

// GetPostsVotes 批量获取帖子的票数
func GetPostsVotes(postIDs *[]redis.Z, postVoteNum *[]float64) error {
	ansCmd := make([]*redis.FloatCmd, len(*postIDs))

	pipeline := rdb.Pipeline()
	for i, postID := range *postIDs {
		ansCmd[i] = pipeline.ZScore(getRedisKey(KeyPostScoreZSet), postID.Member.(string))
	}
	_, err := pipeline.Exec()
	if err != nil {
		return err
	}

	for i, length := 0, len(*postIDs); i < length; i++ {
		(*postVoteNum)[i] = ansCmd[i].Val()
	}
	return nil
}

// DelPost 删除帖子相关的3个zset
func DelPost(s *[]redis.Z, postIsExpired *[]bool) (err error) {
	pipeline := rdb.TxPipeline()
	for i, length := 0, len(*s); i < length; i++ {
		if (*postIsExpired)[i] {
			pipeline.ZRem(getRedisKey(KeyPostTimeZSet), (*s)[i].Member.(string))
			pipeline.Del(getRedisKey(KeyPostVotedZSet) + (*s)[i].Member.(string))
			pipeline.ZRem(getRedisKey(KeyPostScoreZSet), (*s)[i].Member.(string))
		}
	}
	_, err = pipeline.Exec()
	return err
}
