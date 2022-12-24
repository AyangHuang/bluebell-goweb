package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

/*
SignUpHandler 注册
localhost:8080/api/v1/signup

	{
	    "username": "ayang",
	    "password": "123",
	    "re_password": "123"
	}
*/
func SignUpHandler(c *gin.Context) {
	// 1. 参数校验（controller）
	p := new(models.ParamSignUp)
	// 底层调用 ShouldBind，只能调用一次
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp with invalid params", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, errs)
		return
	}
	// 2. 具体业务处理（service）// 3. 持久化（dao）
	if err := logic.SignUp(p); err != nil {
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
		} else {
			ResponseError(c, CodeServerBusy)
		}
		return
	}
	// 4. 返回响应
	ResponseSuccess(c, nil)

}

// localhost:8080/api/v1/login
/*
{
    "username": "ayang",
    "password": "123"
}
*/
func LoginHandler(c *gin.Context) {
	// 1. 参数校验
	p := new(models.ParamLogin)
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
	// 2. 获取验证账户密码
	// 2.业务逻辑处理
	acToken, reToken, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 3.返回响应
	ResponseSuccess(c, gin.H{
		"access_token":  acToken,
		"refresh_token": reToken,
	})
}
