package models

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"wss/db"
	"wss/helper"
)

type UserInfo struct {
	Username string `json:"userName"`
	Platform string `json:"platform"`
	MerchantCode string `json:"merchantCode"`
	VipGrade int `json:"vipGrade"`
	Ip string `json:"ip"`
}

func TokenAuth(token string) *UserInfo {
	RedisPool := db.RedisPool
	conn := RedisPool.Get()
	defer conn.Close()
	str, err := redis.String(conn.Do("GET", token))
	if err != nil {
		fmt.Println(err)
	}
	var user UserInfo
	err = json.Unmarshal([]byte(str), &user)
	if err != nil {
		fmt.Println(err)
	}
	if helper.IsEmpty(user.MerchantCode) {
		user.MerchantCode = "1220817001"
	}
	return &user
}

type IssGetConsumer struct {
	Key string `json:"key"`
	Secret string `json:"secret"`
	RsaPublicKey string `json:"rsa_public_key"`
}
func TokenKeyByIss(iss string) (*IssGetConsumer,error)  {
	RedisPool := db.RedisPool
	conn := RedisPool.Get()
	defer conn.Close()
	res, err := redis.String(conn.Do("HGET", "issGetConsumer" ,iss))
	if err != nil {
		return nil,err
	}
	var issG IssGetConsumer
	err = json.Unmarshal([]byte(res), &issG)
	if err != nil {
		return nil,err
	}
	return &issG,nil
}