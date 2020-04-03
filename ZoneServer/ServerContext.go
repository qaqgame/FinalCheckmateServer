package ZoneServer

import (
	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/Server"
)

// ServerContext : common context of ZoneServer Module
type ServerContext struct {
	Net             *Server.NetManager
	Ipc             *IPCWork.IPCManager
}