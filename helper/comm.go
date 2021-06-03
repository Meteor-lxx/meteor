package helper

type Package struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
	Ext  map[string]interface{} `json:"ext"`
}
