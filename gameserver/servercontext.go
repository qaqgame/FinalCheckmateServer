package gameserver

import (
	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/fsplite"
)
// ServerContext :
type ServerContext struct {
	Fsp       *fsplite.FSPManager
	Ipc       *IPCWork.IPCManager
}
