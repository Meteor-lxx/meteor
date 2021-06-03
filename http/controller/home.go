package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wss/config"
	"wss/db"
	"wss/http/protocol"
)

func IndexApi(c *gin.Context) {
	resp := &protocol.Resp{Code: 200, Msg: "hello wss", Data: ""}
	RedisPool := db.RedisPool
	conn := RedisPool.Get()
	defer conn.Close()
	_, err := conn.Do("Publish", config.SendChannel, "hello")
	if err != nil {
		fmt.Println("redis Publish failed.")
	}

	_, err = conn.Do("Publish", config.IpChannel, "hello")
	if err != nil {
		fmt.Println("redis Publish failed.")
	}
	c.JSON(http.StatusOK, resp)
}
