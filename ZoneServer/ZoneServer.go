package ZoneServer

import (
	"bufio"
	"os"

	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/Server"
	"code.holdonbush.top/ServerFramework/ServerManager"
	log "github.com/sirupsen/logrus"
)

// ZoneServer : zoneserver
type ZoneServer struct {
	*ServerManager.ServerModule
	context         *ServerContext
	onlineManager   *OnlineManager
}

// NewZoneServer :
func NewZoneServer(id, port int, name ...string) *ZoneServer {
	logger := log.WithFields(log.Fields{"Server":"ZoneServer"})
	zoneServer := new(ZoneServer)

	netmanager := Server.NewNetManager(port, logger)
	ipcmanager := IPCWork.NewIPCManager(id)

	tcontext := new(ServerContext)
	tcontext.Net = netmanager
	tcontext.Ipc = ipcmanager

	tonlineManager := NewOnlineManager(tcontext)
	tname := "ZoneServer"
	if len(name) >= 1 {
		tname = name[0]
	}
	Info := ServerManager.ServerModuleInfo{
		Id: id,
		Name: tname,
		Port: port,
	}
	c := make(chan int, 2)
	zoneServer.ServerModule = ServerManager.NewServerModule(Info,logger,ServerManager.UnCreated,c,ipcmanager)
	zoneServer.context = tcontext
	zoneServer.onlineManager = tonlineManager
	
	go zoneServer.ShowDump()
	// zoneServer.context.Ipc.RegisterRPC(zoneServer)
	return zoneServer
}


// Stop : override base Stop
func (zoneServer *ZoneServer) Stop() {
	if zoneServer.context.Net != nil {
		zoneServer.context.Net.Clean()
		zoneServer.context.Net = nil
	}

	if zoneServer.context.Ipc != nil {
		zoneServer.context.Ipc.Clean()
		zoneServer.context.Ipc = nil
	}
	zoneServer.onlineManager.Clean()
}

// Tick : tick
func (zoneServer *ZoneServer) Tick() {
	zoneServer.context.Net.Tick()
}

// ShowDump :
func (zoneServer *ZoneServer) ShowDump() {
	for true {
		reader := bufio.NewReader(os.Stdin)
		str,_ := reader.ReadString('\n')
		if str[:len(str)-2] == "1" {
			zoneServer.Logger.Info("invoke onlinemanager Dump")
			zoneServer.onlineManager.Dump()
		} else if str[:len(str)-2] == "2" {
			zoneServer.Logger.Info("invoke Gateway Dump")
			zoneServer.context.Net.Gateway.Dump()
		} else {
			zoneServer.Logger.Info("input value error")
		}
	}
}