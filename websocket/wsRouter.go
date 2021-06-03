package websocket

import (
	. "wss/websocket/controller"
)

func InitTcpRouter()  {
	RouterRegister("heart",Heart)
	RouterRegister("msgAck",MsgAck)
	RouterRegister("rollAck",RollAck)
}