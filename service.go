package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"wss/models"
	. "wss/websocket"
	"wss/websocket/connect"
)

func wsRouterInit()  {
	InitTcpRouter()
}

func wsConnHandle(c *gin.Context) {
	token := c.Request.FormValue("token")
	parts := strings.SplitN(token, ".", 3)
	var username,platform,merchantCode string
	if len(parts) == 3 {
		jwt,err := models.JwtAuth(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized,err.Error())
			return
		}
		username = jwt.Username
		platform = jwt.Platform
		merchantCode = strconv.Itoa(jwt.MerchantsCode)
	}else{
		tmpToken := strings.SplitN(token, " ", 2)
		if len(tmpToken) == 2 {
			token = tmpToken[1]
		}
		user := models.TokenAuth(token)
		username = user.Username
		platform = user.Platform
		merchantCode = user.MerchantCode
	}
	wsUpgrade.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := wsUpgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer connect.Offline(ws,username,platform,merchantCode)
	connect.Online(ws,username,platform,merchantCode)
	cxt := make(map[string]string)
	cxt["username"] = username
	cxt["platform"] = platform
	cxt["merchantCode"] = merchantCode
	for {
		//读取ws中的数据
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		Dispatch(string(message[4:]),ws,&cxt)
	}
}

func IpConnHandle(c *gin.Context)  {
	ip := c.Request.FormValue("ip")
	if ok, _ := regexp.MatchString("^(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)$", ip); !ok {
		c.JSON(http.StatusUnauthorized,"ip 不正确")
		return
	}
	wsUpgrade.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := wsUpgrade.Upgrade(c.Writer, c.Request, nil)
	defer connect.IpOffLine(ws,ip)
	connect.IpOnLine(ws,ip)
	if err != nil {
		return
	}
	cxt := make(map[string]string)
	cxt["ip"] = ip
	for {
		//读取ws中的数据
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		Dispatch(string(message[4:]),ws,&cxt)
	}
}