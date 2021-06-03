package models

import (
	"time"
	db "wss/db"
)

type Connect struct {
	Id           int    `json:"id" form:"id"`
	Username     string `json:"username" form:"username"`
	Platform     string `json:"platform" form:"platform"`
	MerchantCode string `json:"merchant_code"`
	Status       int `json:"status" form:"status"`
	CreatedAt    string `json:"created_at" form:"created_at"`
}

func GetConnect() []Connect {
	connects := make([]Connect,0)
	rows, err := db.GetConn().Raw("SELECT id, username, platform,status,created_at FROM dlwss_connect_log_2021_06_01 limit 0,10").Rows()
	if err != nil {
		return connects
	}
	for rows.Next() {
		var c Connect
		_ = rows.Scan(&c.Id, &c.Username, &c.Platform, &c.Status, &c.CreatedAt)
		connects = append(connects, c)
	}
	if err = rows.Err(); err != nil {
		return connects
	}
	return connects
}

func (c *Connect)Connects()  {
	t := time.Now()
	table := "dlwss_connect_log_" + t.Format("2006_01_02")
	c.CreatedAt =  t.Format("2006-01-02 15:04:05.000")
	conn := db.GetConn()
	conn.Table(table).Create(c)
}

type IpConnect struct {
	Id           int    `json:"id" form:"id"`
	Ip			 string `json:"ip"`
	Status       int `json:"status" form:"status"`
	CreatedAt    string `json:"created_at" form:"created_at"`
}

func (ip *IpConnect) IpConnects()  {
	t := time.Now()
	table := "dlwss_ip_connect_log_" + t.Format("2006_01_02")
	ip.CreatedAt =  t.Format("2006-01-02 15:04:05.000")
	conn := db.GetConn()
	conn.Table(table).Create(ip)
}