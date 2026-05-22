// Package main is the entry point for the video website application.
// It initializes the database connections and starts the HTTP server.
package main

import (
	"net/http"
	"time"
	"video_website/biz/dal/mysql"
	"video_website/biz/dal/redis"
	handler "video_website/biz/handler"
	chatHandler "video_website/biz/handler/chat"
	router "video_website/biz/router"
	chatService "video_website/biz/service/chat"
	"video_website/config"
	"video_website/pkg/jwt"
	"video_website/pkg/utils"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func main() {
	config.InitConfig()

	jwt.Init(config.Conf.JWT.Secret)

	if err := utils.InitSnowflake("2023-01-01", 1); err != nil {
		hlog.Fatal("初始化雪花ID失败", err)
	}

	if err := mysql.Init(config.Conf.MySQL.DSN); err != nil {
		hlog.Fatal("数据库初始化失败", err)
	}

	if err := redis.Init(); err != nil {
		hlog.Fatal("Redis初始化失败:", err)
	}

	chatService.InitHub()

	h := server.Default(server.WithHostPorts(config.Conf.Server.Address))

	h.Static("/static", "./")

	router.GeneratedRegister(h)

	h.GET("/ping", handler.Ping)

	go func() {
		http.HandleFunc("/ws", chatHandler.HTTPWebSocketHandler)
		hlog.Infof("WebSocket 服务已启动: ws://localhost:6666/ws")
		srv := &http.Server{
			Addr:         ":6666",
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			hlog.Fatalf("WebSocket服务启动失败: %v", err)
		}
	}()

	h.Spin()
}
