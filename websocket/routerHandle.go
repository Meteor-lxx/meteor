package websocket

import (
	"github.com/gorilla/websocket"
)

type RouterInterface func(pack interface{}, conn *websocket.Conn, ext interface{},cxt *map[string]string)

func RouterHandle(routerHandle RouterInterface,data interface{}, conn *websocket.Conn, ext interface{},cxt *map[string]string)  {
	routerHandle(data,conn,ext,cxt)
}