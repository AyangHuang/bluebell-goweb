package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web_app/logger"
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

	e.Use(logger.GinLogger(), logger.GinRecovery(true))

	e.GET("/hello", func(context *gin.Context) {
		context.String(http.StatusOK, "hello, gin")
	})
	return e
}
