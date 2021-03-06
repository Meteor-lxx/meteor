package config

import (
	"gopkg.in/ini.v1"
)

type SysConfig struct {
	Env               string `ini:"env"`
	Debug             bool   `ini:"debug"`
	DBDriver          string `ini:"db_driver"`
	DBHost            string `ini:"db_host"`
	DBPort            string `ini:"db_port"`
	DBUser            string `ini:"db_user"`
	DBPass            string `ini:"db_pass"`
	DBName            string `ini:"db_name"`
	DBDebug           bool   `ini:"db_debug"`
	RedisHost         string `ini:"redis_host"`
	RedisPwd          string `ini:"redis_pwd"`
	RedisDb           int    `ini:"redis_db"`
	RedisCacheVersion string `ini:"redis_cache_version"`

	HttpListenPort   string `ini:"http_listen_port"`
	TcpListenPort    string `ini:"tcp_listen_port"`

	JwtSecret        string `ini:"jwt_secret"`
	JwtExprTime      int64  `ini:"jwt_expr_time"`
	LogDir           string `ini:"log_dir"`
	LogLevel         string `ini:"log_level"`
	GaSecret         string `ini:"ga_secret"`
	EsHost           string `ini:"es_host"`
	MongoHost        string `ini:"mongo_host"`
	MongoConnTimeout int64  `ini:"mongo_conn_timeout"`

	AucHost string `ini:"auc_host"`
}

var Configs *SysConfig = &SysConfig{}

//加载系统配置文件
func Default(configFileName string) {
	config := &SysConfig{}
	conf, err := ini.Load(configFileName) //加载配置文件
	if err != nil {
		panic(err)
	}
	conf.BlockMode = false
	err = conf.MapTo(&config) //解析成结构体
	if err != nil {
		panic(err)
	}
	Configs = config
}
