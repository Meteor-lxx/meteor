package websocket

import (
	"github.com/gorilla/websocket"
	"wss/helper"
	"wss/websocket/parser"
)

type PackParser interface {
	PackDecode() *helper.Package
	PackEncode() []byte
}

func Dispatch(str string, conn *websocket.Conn,cxt *map[string]string) {
	var aes parser.AesPack
	aes.AesData = str
	msg := aes.PackDecode()
	for _,v := range routerTable{
		if v.cmd == msg.Cmd {
			RouterHandle(v.routerController, msg.Data, conn, msg.Ext,cxt)
			break
		}
	}
}
