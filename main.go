package main

import (
	"code.holdonbush.top/FinalCheckmateServer/gameserver"
	"os"
	"time"

	"code.holdonbush.top/FinalCheckmateServer/ZoneServer"
	_ "code.holdonbush.top/FinalCheckmateServer/db"

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
	// ServerFramework.Run(servermanager)
	// server2 := TestServer2.NewTestServer2(2,4051)
	// server1 := TestServer1.NewTestServer1(1,4050)

	// servermanager.AddServer(server1)
	// servermanager.AddServer(server2)

	// servermanager.StartAllServer1()
	zoneServer := ZoneServer.NewZoneServer(1, 4050)
	// TODO:
	gameServer := gameserver.NewGameServer(2, 4051)
	servermanager.AddServer(zoneServer)
	servermanager.AddServer(gameServer)

	servermanager.StartAllServer1()

	for true {
		// servermanager.Tick()
		time.Sleep(time.Millisecond)
	}
}
