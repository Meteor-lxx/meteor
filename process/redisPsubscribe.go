package process

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"unsafe"
	"wss/config"
	"wss/db"
	"wss/websocket/connect"
)

func RedisSubscriberStart()  {
	var sub Subscriber
	sub.Connect()
	sub.Subscribe(config.SendChannel, SendMsg)
	sub.Subscribe(config.IpChannel, IpMsg)
}

type SubscribeCallback func (channel, message string)

type Subscriber struct {
	client redis.PubSubConn
	cbMap map[string]SubscribeCallback
}
func (c *Subscriber) Connect()  {
	RedisPool := db.RedisPool
	conn := RedisPool.Get()
	subConn := redis.PubSubConn{Conn: conn}
	c.client = subConn
	c.cbMap = make(map[string]SubscribeCallback)
	go func() {
		for {
			fmt.Println("wait...")
			switch res := c.client.Receive().(type) {
			case redis.Message:
				channel := (*string)(unsafe.Pointer(&res.Channel))
				message := (*string)(unsafe.Pointer(&res.Data))
				c.cbMap[*channel](*channel, *message)
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", res.Channel, res.Kind, res.Count)
			case error:
				fmt.Println(res.Error())
				continue
			}
		}
	}()
}

func (c *Subscriber) Close() {
	err := c.client.Close()
	if err != nil{
		fmt.Println("redis close error.")
	}
}

func (c *Subscriber) Subscribe(channel interface{}, cb SubscribeCallback) {
	err := c.client.Subscribe(channel)
	if err != nil{
		fmt.Println("redis Subscribe error.")
	}
	c.cbMap[channel.(string)] = cb
}

func SendMsg(channel, msg string){
	var sm connect.SendMsg
	_ = json.Unmarshal([]byte(msg), &sm)
	sm.Send()
	fmt.Println("SendMsg channel : ", channel, " message : ", msg)
}

func IpMsg(channel, msg string){
	fmt.Println("IpMsg channel : ", channel, " message : ", msg)
}