package gameserver

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/ServerManager"
	"code.holdonbush.top/ServerFramework/common"
	"code.holdonbush.top/ServerFramework/fsplite"
	"github.com/sirupsen/logrus"
)

// GameServer : 游戏服务器
type GameServer struct {
	*ServerManager.ServerModule
	context     *ServerContext
	gamemanager *GameManager
	logger      *logrus.Entry
}

// NewGameServer :
func NewGameServer(id, port int, name ...string) *GameServer {
	gameserver := new(GameServer)
	gameserver.context = new(ServerContext)

	tname := "GameServer"
	if len(name) >= 1 {
		tname = name[0]
	}

	Info := new(common.ServerModuleInfo)
	Info.Id = id
	Info.Port = port
	Info.Name = tname
	c := make(chan int, 2)

	gameserver.context.Fsp = fsplite.NewFSPManager(port, DataFormat.IpModel)
	gameserver.context.Ipc = IPCWork.NewIPCManager(Info)

	gameserver.gamemanager = NewGameManager(port,gameserver.context)

	gameserver.logger = logrus.WithFields(logrus.Fields{"Server": "GameServer"})


	gameserver.ServerModule = ServerManager.NewServerModule(Info, gameserver.logger, ServerManager.UnCreated, c, gameserver.context.Ipc)

	// gameserver.context.Ipc.RegisterRPC(gameserver)
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

	gameserver.gamemanager.Clean()

}

// Tick : tick
func (gameserver *GameServer) Tick() {
	gameserver.context.Fsp.Tick()
}
