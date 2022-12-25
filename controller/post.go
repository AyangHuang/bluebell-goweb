package controller

import (
	"bluebell/dao/redis"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// localhost:8080/api/v1/post/release
/*
{
    "title": "我是title2",
    "content": "我是content2"
}
*/
func ReleasePostHandler(c *gin.Context) {
	// 1. 验证参数
	p := new(models.ParamReleasePost)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid params", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs)
		return
	}
	// 2. 逻辑处理
	userID, _ := getCurrentUserID(c)
	p.AuthorID = userID
	if err := logic.ReleasePost(p); err != nil {
		zap.L().Error("logic.ReleasePost", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	} else {
		ResponseSuccess(c, nil)
	}
}

// localhost:8080/api/v1/post/vote
/*
{
    "post_id": "2459427894464512",
    "direction": "1"
}
*/
func VotePostHandler(c *gin.Context) {
	// 1. 验证参数
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login with invalid params", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs)
		return
	}
	userID, _ := getCurrentUserID(c)
	if err := logic.VotePost(userID, p); err != nil {
		if errors.Is(err, redis.ErrorPostNotExist) {
			ResponseError(c, CodePostIDNotExist)
		} else if errors.Is(err, logic.ErrorPostExpired) {
			ResponseError(c, CodePostExpired)
		} else if errors.Is(err, logic.ErrorPostZero) {
			ResponseError(c, CodePostZero)
		} else if errors.Is(err, logic.ErrorPostVoteRepeat) {
			ResponseError(c, CodePostVoteRepeat)
		} else {
			ResponseError(c, CodeServerBusy)
		}
		return
	}
	ResponseSuccess(c, nil)
}
