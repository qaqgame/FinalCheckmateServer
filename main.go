package main

import (
	"code.holdonbush.top/FinalCheckmateServer/ZoneServer"
	// _ "code.holdonbush.top/FinalCheckmateServer/db"
	"code.holdonbush.top/FinalCheckmateServer/gameserver"
	"code.holdonbush.top/ServerFramework"
	"os"

	_ "code.holdonbush.top/ServerFramework/Server"
	"code.holdonbush.top/ServerFramework/ServerManager"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
}

func main() {
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
	//zoneServer := ZoneServer.NewZoneServer(1, 4050)
	//gameServer := gameserver.NewGameServer(2, 4051)
	//servermanager.AddServer(zoneServer)
	//servermanager.AddServer(gameServer)

	servermanager.StartAllServer1()
	ServerFramework.Run(servermanager)
	//for true {
	//	// servermanager.Tick()
	//	time.Sleep(time.Millisecond)
	//}
}
