package main

import (
	"time"

	"code.holdonbush.top/FinalCheckmateServer/ZoneServer"
	_ "code.holdonbush.top/ServerFramework/Server"
	"code.holdonbush.top/ServerFramework/ServerManager"
)

func init() {

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
	servermanager.AddServer(zoneServer)

	servermanager.StartAllServer1()

	for true {
		servermanager.Tick()
		time.Sleep(time.Second)
	}
}
