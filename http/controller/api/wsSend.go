package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"wss/helper"
	"wss/http/protocol"
	"wss/models"
	"wss/websocket/connect"
)

type Conn struct {
	Username   string `json:"username"`
	Icon   string `json:"icon"`
	Ip         string `json:"ip"`
}

func GetConnects(c *gin.Context) {
	resp := &protocol.Resp{Code: 200, Msg: "", Data: ""}
	var all []connect.Conn
	all = connect.GetAll()
	resp.Data = all
	c.JSON(http.StatusOK, resp)
}

func SendMsg(c *gin.Context) {
	resp := &protocol.Resp{Code: 200, Msg: "", Data: ""}
	system := c.Request.FormValue("system")
	if helper.IsEmpty(system) {
		resp.Code = 101
		resp.Msg = "系统不能为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	platform := c.Request.FormValue("platform")
	if helper.IsEmpty(platform) {
		resp.Code = 101
		resp.Msg = "平台不能为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	sendData := c.Request.FormValue("sendData")
	if helper.IsEmpty(sendData) {
		resp.Code = 101
		resp.Msg = "发送数据不能为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	userStr := c.Request.FormValue("userList")
	if helper.IsEmpty(userStr) {
		resp.Code = 101
		resp.Msg = "发送人不能为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	merchantCode := c.Request.FormValue("merchantCode")
	if helper.IsEmpty(merchantCode) {
		merchantCode = "1220817001"
	}
	userList := make([]string,10)
	err := json.Unmarshal([]byte(userStr), &userList)
	if err != nil {
		resp.Code = 101
		resp.Msg = "用户列表格式不正确"
		c.JSON(http.StatusOK, resp)
		return
	}
	var msg models.Message
	msg.System = system
	msg.MerchantCode = merchantCode
	msg.Platform = platform
	msg.SendData = sendData
	msg.SendUserList = userStr
	msg.Save()
	var ms models.MessageStatus
	ms.Mid = msg.Id
	ms.MerchantCode = merchantCode
	ms.ErrorMsg = "已发送"
	ms.IsSend = 1
	ms.Save()
	var sendMsg connect.SendMsg
	sendMsg.Platform = platform
	sendMsg.SendData = sendData
	sendMsg.UserList = &userList
	sendMsg.MerchantCode = merchantCode
	sendMsg.MsgId = msg.Id
	sendMsg.SendChannel()
	c.JSON(http.StatusOK, resp)
}

func SendPlatform(c *gin.Context)  {
	resp := &protocol.Resp{Code: 200, Msg: "", Data: ""}
	system := c.Request.FormValue("system")
	if helper.IsEmpty(system) {
		resp.Code = 101
		resp.Msg = "系统不能为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	platform := c.Request.FormValue("platform")
	if helper.IsEmpty(platform) {
		resp.Code = 101
		resp.Msg = "平台不能为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	sendData := c.Request.FormValue("sendData")
	if helper.IsEmpty(sendData) {
		resp.Code = 101
		resp.Msg = "发送数据不能为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	merchantCode := c.Request.FormValue("merchantCode")
	if helper.IsEmpty(merchantCode) {
		merchantCode = "1220817001"
	}
	var msg models.Message
	msg.System = system
	msg.Platform = platform
	msg.SendData = sendData
	msg.MerchantCode = merchantCode
	msg.Save()
	connect.SendPlatform(platform,sendData,merchantCode)
	c.JSON(http.StatusOK, resp)
}

func GetOnline(c *gin.Context)  {
	resp := &protocol.Resp{Code: 200, Msg: "", Data: ""}
	userStr := c.Request.FormValue("users")
	if helper.IsEmpty(userStr) {
		resp.Code = 101
		resp.Msg = "用户列表不能为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	merchantCode := c.Request.FormValue("merchantCode")
	if helper.IsEmpty(merchantCode) {
		merchantCode = "1220817001"
	}
	users := strings.Split(userStr,",")
	connect.GetOnLineUser(users,resp)
	c.JSON(http.StatusOK, resp)
}

func IpSend(c *gin.Context) {
	resp := &protocol.Resp{Code: 200, Msg: "", Data: ""}
	ip := c.Request.FormValue("ip")
	if helper.IsEmpty(ip) {
		resp.Code = 101
		resp.Msg = "发送IP不能为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	sendData := c.Request.FormValue("sendData")
	if helper.IsEmpty(sendData) {
		resp.Code = 101
		resp.Msg = "发送数据不能为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	notice := c.Request.FormValue("notice")
	var im models.IpMessage
	im.Ip = ip
	im.SendData = sendData
	im.IsSend = 1
	im.Notice = notice
	im.IpMsg()
	connect.IpSend(ip,sendData,im.Id)
	c.JSON(http.StatusOK, resp)
}

func SendIpList(c *gin.Context) {
	resp := &protocol.Resp{Code: 200, Msg: "", Data: ""}
	ipList := c.Request.FormValue("ipList")
	if helper.IsEmpty(ipList) {
		resp.Code = 101
		resp.Msg = "发送IP不能为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	sendData := c.Request.FormValue("sendData")
	if helper.IsEmpty(sendData) {
		resp.Code = 101
		resp.Msg = "发送数据不能为空"
		c.JSON(http.StatusOK, resp)
		return
	}
	notice := c.Request.FormValue("notice")
	list := strings.Split(ipList,",")
	for _,v := range list{
		var im models.IpMessage
		im.Ip = v
		im.SendData = sendData
		im.IsSend = 1
		im.Notice = notice
		im.IpSave()
		connect.IpSend(v,sendData,im.Id)
	}
	c.JSON(http.StatusOK, resp)
}