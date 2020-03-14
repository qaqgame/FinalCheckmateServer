package main

import (
	"code.holdonbush.top/ServerFramework"
	"code.holdonbush.top/ServerFramework/ServerManager"
)

func main() {
	serverManager := ServerManager.ServerManager{}
	serverManager.Init()
	serverManager.StartServer(1)
	serverManager.StartServer(1)
	ServerFramework.Run(&serverManager)   // 主循环

	serverManager.StopAllServer()
}
