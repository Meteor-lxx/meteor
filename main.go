package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"wss/config"
	"wss/db"
	"wss/helper"
	ginHtp "wss/http"
	"wss/loger"
	"wss/process"
)

func init() {
	appDir := helper.GetAppDir()
	config.Default(appDir + "/config.ini")
	db.Default()
	db.DefaultRedis()
	process.RedisSubscriberStart()
}

var wsUpgrade = websocket.Upgrader{
	CheckOrigin: func (r *http.Request) bool {
		return true
	},
}

func main() {
	if config.Configs.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	loger.Default()
	wsRouterInit()
	router := ginHtp.InitRouter()
	router.GET("/", wsConnHandle)
	router.GET("/ip", IpConnHandle)
	err := router.Run(":80")
	if err!=nil {
		loger.Loggers.Error(err.Error())
	}
}