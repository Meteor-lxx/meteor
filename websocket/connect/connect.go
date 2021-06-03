package connect

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"time"
	"wss/config"
	"wss/db"
	"wss/helper"
	"wss/http/protocol"
	"wss/models"
	"wss/websocket/parser"
)

type WebsocketConnect struct {
	Username      string
	MerchantCode  string
	Platform  	  string
	LastTime	  time.Time
	Conn          *websocket.Conn
}

var ClientConn map[string][]WebsocketConnect

func init() {
	ClientConn = make(map[string][]WebsocketConnect,8)
}

func Online(conn *websocket.Conn, username string,platform string,merchantCode string) {
	user := WebsocketConnect{Username: username, MerchantCode:merchantCode, Platform:platform, LastTime:time.Now(), Conn: conn}
	ClientConn[platform] = append(ClientConn[platform],user)
	var c models.Connect
	c.Username = username
	c.Platform = platform
	c.MerchantCode = merchantCode
	c.Status = 1
	c.Connects()
}

func Offline(conn *websocket.Conn,username string,platform string,merchantCode string) {
	var i = 0
	for i = 0; i < len(ClientConn[platform]); i++ {
		if ClientConn[platform][i].Username == username && ClientConn[platform][i].MerchantCode == merchantCode {
			ClientConn[platform] = append(ClientConn[platform][:i], ClientConn[platform][i+1:]...)
		}
	}
	var c models.Connect
	c.Username = username
	c.Platform = platform
	c.MerchantCode = merchantCode
	c.Status = 0
	c.Connects()
	_ = conn.Close()
}

type Conn struct {
	Username      string `json:"username"`
	MerchantCode  string `json:"merchant_code"`
	LastTime	  string
	Platform string `json:"platform"`
}

func GetAll() []Conn  {
	var all []Conn
	var conn Conn
	for k,_ := range ClientConn {
		for _,v := range ClientConn[k] {
			conn.Platform = k
			conn.Username = v.Username
			conn.LastTime = v.LastTime.Format("2006-01-02 15:04:05")
			conn.MerchantCode = v.MerchantCode
			all = append(all, conn)
		}
	}
	return all
}

type SendMsg struct {
	Platform string `json:"platform"`
	SendData string `json:"send_data"`
	UserList *[]string `json:"user_list"`
	MerchantCode string `json:"merchant_code"`
	MsgId int `json:"msg_id"`
}

func (sm *SendMsg)SendChannel()  {
	str, _ :=json.Marshal(sm)
	RedisPool := db.RedisPool
	conn := RedisPool.Get()
	defer conn.Close()
	_, _ = conn.Do("Publish", config.SendChannel, str)
}
func (sm *SendMsg)Send()  {
	var aes parser.AesPack
	var data helper.Package
	ext := make(map[string]interface{})
	ext["msgId"] = sm.MsgId
	data.Cmd = "server"
	data.Data = sm.SendData
	data.Ext = ext
	aes.Data = &data
	if sm.Platform == "0" {
		for _,username := range *sm.UserList {
			for k,_ := range ClientConn {
				for _,conn := range ClientConn[k] {
					if conn.MerchantCode == sm.MerchantCode && conn.Username == username {
						_ = conn.Conn.WriteMessage(websocket.BinaryMessage, aes.PackEncode())
					}
				}
			}
		}
		return
	}
	for _,username := range *sm.UserList {
		for _,conn := range ClientConn[sm.Platform] {
			if conn.MerchantCode == sm.MerchantCode && conn.Username == username {
				_ = conn.Conn.WriteMessage(websocket.BinaryMessage, aes.PackEncode())
			}
		}
	}
}

func LastHeart(username string,platform string,merchantCode string)  {
	for _,conn := range ClientConn[platform]  {
		if conn.MerchantCode == merchantCode && conn.Username == username {
			conn.LastTime = time.Now()
		}
	}
}

func SendPlatform(platform string, sendData string,merchantCode string)  {
	var aes parser.AesPack
	var data helper.Package
	data.Cmd = "server"
	data.Data = sendData
	aes.Data = &data
	for _,conn := range ClientConn[platform] {
		if conn.MerchantCode == merchantCode {
			_ = conn.Conn.WriteMessage(websocket.BinaryMessage, aes.PackEncode())
		}
	}
}

func GetOnLineUser(users []string,resp *protocol.Resp) *protocol.Resp {
	a := make(map[string]interface{})
	for _,v := range users {
		arr := make(map[string]string)
		for _,conn := range ClientConn["ios"] {
			if conn.Username == v {
				arr["ios"] = "online"
				a[v] = arr
			}
		}
		for _,conn := range ClientConn["android"] {
			if conn.Username == v {
				arr["android"] = "online"
				a[v] = arr
			}
		}
	}
	resp.Data = a
	return resp
}