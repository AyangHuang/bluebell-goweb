package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/utils/snowflake"
	"errors"
	"go.uber.org/zap"
	"time"
)

// ReleasePost 发布帖子
func ReleasePost(post *models.ParamReleasePost) (err error) {

	post.ID = snowflake.GenID()

	// 1. 存入数据库
	if err = mysql.InsertPost(post); err != nil {
		return
	}
	// 2. 存入redis
	err = redis.CreatePost(post.ID)
	if err == nil {
		zap.L().Info("发布帖子成功",
			zap.Int64("userID", post.AuthorID),
			zap.Int64("postID", post.ID))
	}
	return
}

// VotePost
/* 投票的几种情况：
direction=1时，有两种情况：
	1. 之前没有投过票，现在投赞成票
	2. 之前投反对票，现在改投赞成票
direction=0时，有两种情况：
	1. 之前投过赞成票，现在要取消投票
	2. 之前投过反对票，现在要取消投票
direction=-1时，有两种情况：
	1. 之前没有投过票，现在投反对票
	2. 之前投赞成票，现在改投反对票

投票的限制：
每个贴子自发表之日起一个星期之内允许用户投票，超过一个星期就不允许再投票了。
	1. 到期之后将redis中保存的赞成票数及反对票数存储到mysql表中
	2. 投票的总数能够增加或减少帖子的剩余投票时间
*/

const (
	oneWeekInSeconds = 7 * 24 * 3600
	scorePerVote     = 432 // 每一票延长或减少多少秒
)

var (
	ErrorPostExpired    = errors.New("帖子已过期")
	ErrorPostZero       = errors.New("用户未投过票，无法取消")
	ErrorPostVoteRepeat = errors.New("重复投票")
)

func VotePost(userID int64, vote *models.ParamVoteData) (err error) {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", vote.PostID),
		zap.Int8("direction", vote.Direction))
	//redis.VotePost(vote)
	// 1.取发帖时间
	var postTime float64
	if postTime, err = redis.GetPostTime(vote.PostID); postTime == 0 {
		zap.L().Error("post 不存在 post:time 中", zap.String("postID", vote.PostID))
		// 帖子不存在
		return
	}

	// 2.获取票数
	votes := redis.GetPostVotes(vote.PostID)

	// 3.查看过期时间与否 现在-发布 >= 7 天过期时间+投票增加或减少时间
	if JudgePostExpired(postTime, votes) {
		zap.L().Info("post 已过期", zap.String("postID", vote.PostID))
		return ErrorPostExpired
	}

	// 4.投票逻辑
	// 4.1 获取用户之前对此帖子的投票情况
	voted := redis.GetPostUserVoted(vote.PostID, userID) // 同下，且 0 表示没投过票
	direction := vote.Direction                          // 1 表示赞同，-1 表示反对， 0 表示取消投票

	// 重复投票
	if voted != 0 && voted == float64(direction) {
		return ErrorPostVoteRepeat
	}

	// 1.未投过票
	if voted == 0 {
		if direction == 0 {
			// 1.1 未投过票，无法取消
			return ErrorPostZero
		} else if direction == 1 {
			// 1.2 投赞成票
			err = redis.PostVote(userID, vote.PostID, 1, direction)
		} else if direction == -1 {
			// 1.3 投反对票
			err = redis.PostVote(userID, vote.PostID, -1, direction)
		}
		// 2.以前投过赞成票
	} else if voted == 1 {
		// 2.1 取消投票
		if direction == 0 {
			err = redis.PostVote(userID, vote.PostID, -1, 0)
			// 2.2 改投反对票
		} else if direction == -1 {
			err = redis.PostVote(userID, vote.PostID, -2, direction)
		}
		// 3. 以前投过反对票
	} else if voted == -1 {
		// 2.1 取消投票
		if direction == 0 {
			err = redis.PostVote(userID, vote.PostID, 1, 0)
			// 2.2 改投赞成票
		} else if direction == 1 {
			err = redis.PostVote(userID, vote.PostID, 2, direction)
		}
	}
	if err == nil {
		zap.L().Info("vote success",
			zap.String("postID", vote.PostID),
			zap.Int64("userID", userID),
			zap.Float64("pre direction", voted),
			zap.Int8("vote direction", direction))
	}
	return
}

// JudgePostExpired 判断是否过期
func JudgePostExpired(postTime, votes float64) bool {
	return float64(time.Now().Unix())-postTime >= oneWeekInSeconds+votes*scorePerVote
}

// JudgePostOutTime 判断是否过了预定日期，不包含投票增加减少的
func JudgePostOutTime(postTime float64) bool {
	return float64(time.Now().Unix())-postTime >= oneWeekInSeconds
}
