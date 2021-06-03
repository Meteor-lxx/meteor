package middleware

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"wss/config"
	"wss/helper"
	"wss/http/protocol"
	"wss/loger"
)

type CustomClaims struct {
	Iss           string      `json:"iss"`
	Cid           int         `json:"cid"`
	MCid          int         `json:"mcid"`
	Name          string      `json:"name"`
	MerchantsCode string      `json:"merchantsCode"`
	IsConsumer    int         `json:"isConsumer"`
	AppKey        interface{} `json:"appkey"`
	Iat           int         `json:"iat"`
	Exp           int         `json:"exp"`
	Nbf           int         `json:"nbf"`
	Jti           string      `json:"jti"`
}

type AucData struct {
	MCid int    `json:"mcid"`
	Cid  int    `json:"cid"`
	Path string `json:"path"`
}

type AucRes struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func AucAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("M-Token")
		path := c.FullPath()
		// parseToken 解析token包含的信息
		claims := ParseToken(token)
		data := AucData{Cid: claims.Cid, MCid: claims.MCid, Path: path}
		res := helper.PostRequest(config.Configs.AucHost+"/auth/api_auth", data, "application/json")
		var aucRes AucRes
		jsonRes := json.Unmarshal([]byte(res), &aucRes)
		if jsonRes != nil {
			loger.Loggers.Error("json convert fail!")
			return
		}
		if aucRes.Code != 200 {
			resp := &protocol.Resp{Code: aucRes.Code, Msg: aucRes.Msg, Data: ""}
			c.JSON(http.StatusOK, resp)
			c.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		c.Set("claims", claims)
	}
}

func ParseToken(tokenString string) (claims *CustomClaims) {
	arr := strings.Split(tokenString, ".")
	tokenDe, _ := base64.StdEncoding.DecodeString(arr[1])
	jsonRes := json.Unmarshal(tokenDe, &claims)
	if jsonRes != nil {
		loger.Loggers.Error("json convert fail!")
		return
	}
	return claims
}
