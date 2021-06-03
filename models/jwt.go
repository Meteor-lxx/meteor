package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

type JwtHeader struct {
	Type string `json:"type"`
	Alg string `json:"alg"`
	Jti string `json:"jti"`
}

type CustomClaims struct {
	Iss           string      `json:"iss"`
	Username      string      `json:"username"`
	MerchantsCode int         `json:"merchant_code"`
	Version       string      `json:"version"`
	Platform      string      `json:"platform"`
	Ip       	  string      `json:"ip"`
	UdId       	  string      `json:"udid"`
	jwt.StandardClaims
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

type JwtError struct {
	msg string
}
func (err *JwtError) Error() string {
	strFormat := `jwt解析失败: %v`
	return fmt.Sprintf(strFormat, err.msg)
}

func ParseToken(tokenString string) (claims *CustomClaims,err error) {
	arr := strings.Split(tokenString, ".")
	tokenDe, _ := base64.RawStdEncoding.DecodeString(arr[1])
	jsonRes := json.Unmarshal(tokenDe, &claims)
	var jwtErr JwtError
	if jsonRes != nil {
		jwtErr.msg = "json 格式不正确"
		return nil,&jwtErr
	}
	return claims,nil
}


func ParserWebToken(tokenString string) (*CustomClaims,error) {
	var jwtErr JwtError
	claims,err := ParseToken(tokenString)
	if err != nil {
		jwtErr.msg = err.Error()
		return nil, &jwtErr
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		key, err := TokenKeyByIss(claims.Iss)
		if err != nil {
			jwtErr.msg = err.Error()
			return nil, &jwtErr
		}
		if token.Header["alg"] == "HS256" {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				jwtErr.msg = "Unexpected signing method"
				return nil, &jwtErr
			}
			return key.Secret, nil
		}else{
			pb, err := jwt.ParseRSAPublicKeyFromPEM([]byte(key.RsaPublicKey))
			if err!=nil {
				jwtErr.msg = "私钥解析失败"
				return nil, &jwtErr
			}
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				jwtErr.msg = "token 验签失败"
				return nil, &jwtErr
			}
			return pb, nil
		}
	})
	valid := token.Valid
	if valid {
		return claims,nil
	}
	jwtErr.msg = "验签失败"
	return nil, &jwtErr
}

func JwtAuth(token string) (*CustomClaims,error) {
	parts := strings.SplitN(token, " ", 2)
	tokenNew := token
	if len(parts) == 2 {
		tokenNew = parts[1]
	}
	claims, err :=ParserWebToken(tokenNew)
	if err != nil {
		return nil, err
	}
	return claims,nil
}
