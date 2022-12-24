package routers

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	switch mode {
	case gin.ReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	case gin.DebugMode:
		gin.SetMode(gin.DebugMode)
	case gin.TestMode:
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}
	e := gin.New()

	v1 := e.Group("/api/v1")

	v1.Use(logger.GinLogger(), logger.GinRecovery(true))
	// 注册
	v1.POST("/signup", controller.SignUpHandler)
	// 登录
	v1.POST("/login", controller.LoginHandler)

	// 携带 refreshToken 请求 accessToken，即自动登录
	v1.POST("/autologin", middlewares.JWTMiddleWare(), controller.AutoLoginHandler)
	// 后面都是需要 JWT 认证登录后才能访问的
	v1.Use(middlewares.JWTMiddleWare())
	return e
}
