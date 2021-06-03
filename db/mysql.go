package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"wss/config"
)

var dbPool *gorm.DB

func Default() {
	s := config.Configs.DBUser + ":" + config.Configs.DBPass + "@tcp(" + config.Configs.DBHost + ":" + config.Configs.DBPort + ")/" + config.Configs.DBName + "?charset=utf8&loc=Local&parseTime=true"
	conn, err := gorm.Open(config.Configs.DBDriver, s)
	if err != nil {
		panic(err.Error() + config.Configs.DBDriver)
	}
	conn.DB().SetMaxOpenConns(20)
	conn.DB().SetMaxIdleConns(5)
	if config.Configs.DBDebug {
		conn = conn.Debug()
	}
	dbPool = conn
}

func GetConn() *gorm.DB {
	return dbPool
}
