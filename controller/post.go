package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

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
		ResponseError(c, CodeServerBusy)
	} else {
		ResponseSuccess(c, nil)
	}
}
