package connect

import (
	"github.com/gorilla/websocket"
	"time"
	"wss/helper"
	"wss/models"
	"wss/websocket/parser"
)

type IpConn struct {
	Ip 			  string
	LastTime	  time.Time
	Conn          *websocket.Conn
}

var IpAllConn []IpConn

func IpOnLine(conn *websocket.Conn,ip string)  {
	ipConn := IpConn{Ip:ip,Conn:conn,LastTime:time.Now()}
	IpAllConn = append(IpAllConn,ipConn)
	var ipc models.IpConnect
	ipc.Ip = ip
	ipc.Status = 1
	ipc.IpConnects()
}

func IpOffLine(conn *websocket.Conn,ip string)  {
	var i int
	for i = 0 ; i < len(IpAllConn); i++ {
		if IpAllConn[i].Ip == ip {
			IpAllConn = append(IpAllConn[:i], IpAllConn[i+1:]...)
		}
	}
	var ipc models.IpConnect
	ipc.Ip = ip
	ipc.Status = 0
	ipc.IpConnects()
	_ = conn.Close()
}

func IpSend(ip string, sendData string,msgId int)  {
	var aes parser.AesPack
	var data helper.Package
	ext := make(map[string]interface{})
	ext["msgId"] = msgId
	data.Cmd = "server"
	data.Data = sendData
	data.Ext = ext
	aes.Data = &data
	for _,v := range IpAllConn {
		if v.Ip == ip {
			_ = v.Conn.WriteMessage(websocket.BinaryMessage, aes.PackEncode())
		}
	}
}