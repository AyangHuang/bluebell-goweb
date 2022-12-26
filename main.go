package main

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/other_serve/persist"
	"bluebell/routers"
	"bluebell/settings"
	"bluebell/utils/jwt"
	"bluebell/utils/snowflake"
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// gin web 项目脚手架

func main() {
	// 1. 加载配置
	if err := settings.Init(); err != nil {
		log.Fatalf("setting.Init err: %s", err)
	}
	jwt.Init()
	if err := snowflake.Init(settings.Conf.Time, 1); err != nil {
		log.Fatalf("snowflake.Init err：%s", err)
	}

	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		log.Fatalf("logger.Init err: %s", err)
	}

	// 3. 初始化数据库连接（包括 mysql 和 redis）
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		zap.L().Fatal("mysql.Init err: %s", zap.Error(err))
	}
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		zap.L().Fatal("redis.Init err: %s", zap.Error(err))
	}

	// 4. 注册路由
	e := routers.Setup(gin.DebugMode)

	// 5. 启动服务（优雅关机）
	srv := &http.Server{
		Addr:    ":8080",
		Handler: e,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listing", zap.Error(err))
		}
	}()

	// 6.从 redis 把过了投票期的刷回 mysql
	go func() {
		persist.PersistVotesToDB()
		time.Sleep(time.Hour)
	}()

	quit := make(chan os.Signal, 1)
	// 注册监听 chan，收到信号操作系统会调用 go 库函数往 quit 里面发送一个 os.Signal
	// syscall.SIGINT kill -2 pid
	signal.Notify(quit, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown", zap.Error(err))
	}
	zap.L().Info("Server exiting!")
}
