package controller

import (
	"github.com/gorilla/websocket"
	"wss/helper"
	"wss/websocket/connect"
	"wss/websocket/parser"
)

func Heart(pack interface{},conn *websocket.Conn,ext interface{},cxt *map[string]string)  {
	var aes parser.AesPack
	var data helper.Package
	data.Cmd = "heart"
	aes.Data = &data
	_ = conn.WriteMessage(websocket.BinaryMessage, aes.PackEncode())
	connect.LastHeart((*cxt)["username"],(*cxt)["platform"],(*cxt)["merchantCode"])
}