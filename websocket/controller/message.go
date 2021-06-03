package controller

import (
	"github.com/gorilla/websocket"
	"wss/models"
)

func MsgAck(pack interface{},conn *websocket.Conn, ext interface{},cxt *map[string]string)  {
	m := ext.(map[string]interface{})
	mid := m["msgId"]
	msgId := mid.(float64)
	var ms models.MessageStatus
	ms.Mid = int(msgId)
	ms.Username = (*cxt)["username"]
	ms.MerchantCode = (*cxt)["merchantCode"]
	ms.Platform = (*cxt)["platform"]
	ms.ErrorMsg = "已回复"
	ms.IsSend = 2
	ms.UpdateByMid()
}

func RollAck(pack interface{},conn *websocket.Conn, ext interface{},cxt *map[string]string)  {
	m := ext.(map[string]interface{})
	mid := m["msgId"]
	msgId := mid.(float64)
	var im models.IpMessage
	im.Id = int(msgId)
	im.IsSend = 2
	im.ErrorMsg = "已回复"
	im.IpMsg()
}