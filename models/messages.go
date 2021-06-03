package models

import (
	"time"
	"wss/db"
)

type Message struct {
	Id           int    `json:"id" form:"id"`
	System       string `json:"system"`
	MerchantCode string `json:"merchant_code"`
	Platform     string `json:"platform"`
	SendData     string `json:"send_data"`
	SendUserList string `json:"send_user_list"`
	CreatedAt    string `json:"created_at"`
}

type MessageStatus struct {
	Id           int    `json:"id"`
	Mid          int    `json:"mid"`
	MerchantCode string `json:"merchant_code"`
	Username     string `json:"username"`
	IsSend       int    `json:"is_send"`
	Platform     string `json:"platform"`
	ErrorMsg     string `json:"error_msg"`
	Notice       string `json:"notice"`
	CreatedAt    string `json:"created_at"`
}

func (m *Message)Save()  {
	t := time.Now()
	table := "dlwss_message_" + t.Format("2006_01")
	m.CreatedAt = t.Format("2006-01-02 15:04:05")
	conn := db.GetConn()
	conn.Table(table).Create(m)
}

func (ms *MessageStatus)Save()  {
	t := time.Now()
	table := "dlwss_message_status_" + t.Format("2006_01_02")
	ms.CreatedAt = t.Format("2006-01-02 15:04:05")
	conn := db.GetConn()
	conn.Table(table).Create(ms)
}

func (ms *MessageStatus)UpdateByMid()  {
	t := time.Now()
	table := "dlwss_message_status_" + t.Format("2006_01_02")
	conn := db.GetConn()
	conn.Table(table).Where("mid=?",ms.Mid).Update(ms)
}

type IpMessage struct {
	Id           int    `json:"id"`
	Ip			 string `json:"ip"`
	SendData     string `json:"send_data"`
	IsSend       int    `json:"is_send"`
	ErrorMsg     string `json:"error_msg"`
	Notice       string `json:"notice"`
	CreatedAt    string `json:"created_at"`
}

func (im *IpMessage)IpSave() {
	t := time.Now()
	table := "dlwss_roll_broadcast_" + t.Format("2006_01_02")
	im.CreatedAt = t.Format("2006-01-02 15:04:05")
	conn := db.GetConn()
	conn.Table(table).Create(im)
}

func (im *IpMessage)IpMsg()  {
	t := time.Now()
	table := "dlwss_roll_broadcast_" + t.Format("2006_01_02")
	conn := db.GetConn()
	conn.Table(table).Where("id=?",im.Id).Update(IpMessage{IsSend:2,ErrorMsg:"已回复"})
}