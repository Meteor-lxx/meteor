package db

import (
	redisGo "github.com/gomodule/redigo/redis"
	"time"
	"wss/config"
)

var RedisPool *redisGo.Pool

func DefaultRedis()  {
	RedisPool = PoolInitRedis(config.Configs.RedisHost,config.Configs.RedisPwd)
}

// redis pool
func PoolInitRedis(server string, password string) *redisGo.Pool {
	return &redisGo.Pool{
		MaxIdle:     4,//空闲数
		IdleTimeout: 240 * time.Second,
		MaxActive:   8,//最大数
		Dial: func() (redisGo.Conn, error) {
			redisGo.DialConnectTimeout(5)
			redisGo.DialReadTimeout(2)
			redisGo.DialWriteTimeout(2)
			c, err := redisGo.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					err = c.Close()
					return nil, err
				}
			}
			redisGo.DialDatabase(config.Configs.RedisDb)
			return c, err
		},
		TestOnBorrow: func(c redisGo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
