package websocket

type router struct {
	cmd string
	routerController RouterInterface
}

var routerTable []router

func RouterRegister(cmd string,routerInterface RouterInterface)  {
	router := router{cmd:cmd,routerController:routerInterface}
	routerTable = append(routerTable,router)
}