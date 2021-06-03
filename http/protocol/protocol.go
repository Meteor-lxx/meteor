package protocol

type Resp struct {
	Code      int       `json:"code"`
	Data      interface{} `json:"data"`
	Msg       string      `json:"msg"`
}

type AdminJwtTokenInfo struct {
	AdminId    int32
	Username   string
	ExpireTime int32
}
