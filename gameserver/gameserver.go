package gameserver

import (
	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/ServerManager"
	"code.holdonbush.top/ServerFramework/fsplite"
	"github.com/sirupsen/logrus"
)

// GameServer : 游戏服务器
type GameServer struct {
	*ServerManager.ServerModule
	context       *ServerContext
	logger        *logrus.Entry
}

// NewGameServer :
func NewGameServer(id, port int, name ...string) *GameServer {
	gameserver := new(GameServer)
	gameserver.context = new(ServerContext)

	gameserver.context.Fsp = fsplite.NewFSPManager(port)
	gameserver.context.Ipc = IPCWork.NewIPCManager(id)

	gameserver.logger = logrus.WithFields(logrus.Fields{"Server":"GameServer"})

	tname := "GameServer"
	if len(name) >= 1 {
		tname = name[0]
	}

	Info := ServerManager.ServerModuleInfo{
		Id: id,
		Name: tname,
		Port: port,
	}
	c := make(chan int, 2)
	gameserver.ServerModule = ServerManager.NewServerModule(Info, gameserver.logger, ServerManager.UnCreated, c, gameserver.context.Ipc)


	return gameserver
}

// Stop : stop gameserver
func (gameserver *GameServer) Stop() {
	if gameserver.context.Fsp != nil {
		gameserver.context.Fsp.Clean()
		gameserver.context.Fsp = nil
	}

	if gameserver.context.Ipc != nil {
		gameserver.context.Ipc.Clean()
		gameserver.context.Ipc = nil
	}
	
}

// Tick : tick
func (gameserver *GameServer) Tick() {
	gameserver.context.Fsp.Tick()
}