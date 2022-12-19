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

	v1 := e.Group("/api/v1")

	v1.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1.GET("/hello", func(context *gin.Context) {
		context.String(http.StatusOK, "hello, gin")
	})
	return e
}
