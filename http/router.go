package http

import (
	"github.com/gin-gonic/gin"
	. "wss/http/controller"
	"wss/http/controller/api"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/home", IndexApi)
	v1 := router.Group("v1")
	v1.POST("/send", api.SendMsg)
	v1.POST("/sendIp", api.IpSend)
	v1.POST("/sendIpList", api.SendIpList)
	v1.POST("/sendAll", api.SendPlatform)
	v1.POST("/getOnline", api.GetOnline)
	v1.POST("/conn", api.GetConnects)

	return router
}
