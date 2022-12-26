package persist

import (
	"bluebell/dao/mysql"
	myredis "bluebell/dao/redis"
	"bluebell/logic"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

const perGet = 100

func PersistVotesToDB() {

	var start int64 = 0
	var stop int64 = perGet
	var nonRe = true

	s := make([]redis.Z, 0)
	postVoteNum := make([]float64, perGet, perGet)
	postIsExpired := make([]bool, perGet, perGet)

	for nonRe {
		var expriedNum int64 = 0
		nonRe = false

		// 1. 批量获取 100 个帖子数据
		myredis.GetVotes(start, stop, &s)
		if err := myredis.GetPostsVotes(&s, &postVoteNum); err != nil {
			zap.L().Error("myredis.GetPostsVotes error", zap.Error(err))
		}
		for index, postZ := range s {
			// 2. 判断是否过期
			if logic.JudgePostExpired(postZ.Score, postVoteNum[index]) {
				nonRe = true
				postIsExpired[index] = true
				expriedNum++
			}
		}
		// 3. 把过期的持久化数据库
		if err := mysql.UpdatePostVoteNum(&s, &postVoteNum, &postIsExpired); err != nil {
			zap.L().Error("mysql.UpdatePostVoteNum error ", zap.Error(err))
			return
		}
		// 4. 删除redis
		if err := myredis.DelPost(&s, &postIsExpired); err != nil {
			zap.L().Error("myredis.DelPost error ", zap.Error(err))
			return
		}

		for i, length := 0, len(s); i < length; i++ {
			if postIsExpired[i] {
				zap.L().Info("PersistVotesToDB success", zap.String("postID", s[i].Member.(string)))
			}
		}

		start += perGet - expriedNum
		stop += perGet - expriedNum
	}
}
