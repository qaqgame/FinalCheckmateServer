package main

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	_ "code.holdonbush.top/FinalCheckmateServer/Roles"
	"code.holdonbush.top/FinalCheckmateServer/ZoneServer"
	"flag"
	"fmt"

	// _ "code.holdonbush.top/FinalCheckmateServer/db"
	"code.holdonbush.top/FinalCheckmateServer/gameserver"
	"code.holdonbush.top/ServerFramework"
	_ "code.holdonbush.top/ServerFramework/Server"
	"code.holdonbush.top/ServerFramework/ServerManager"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	flag.IntVar(&DataFormat.IpModel,"ipModel",0,"选择使用公网IP--输入0(这也是默认选项),还是局域网IP--输入1")
	flag.Parse()
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
}

func main() {
	fmt.Println("IPMODEL", DataFormat.IpModel)
	servermanager := ServerManager.NewServerManager()

	serverConfig :=ServerManager.GetAllServerModuleInfo()
	for _,v := range serverConfig {
		switch v.Name {
		case "ZoneServer":
			servermanager.AddServer(ZoneServer.NewZoneServer(v.Id, v.Port, v.Name))
		case "GameServer":
			servermanager.AddServer(gameserver.NewGameServer(v.Id, v.Port, v.Name))
		}
	}

	servermanager.StartAllServer1()
	ServerFramework.Run(servermanager)
}
