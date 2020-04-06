package ZoneServer

import (
	log "github.com/sirupsen/logrus"
)

type RoomManager struct {
	listRoom          []*Room
	context           *ServerContext
	logger            *log.Entry
}

// NewRoomManager : new a room manager
func NewRoomManager(context *ServerContext) *RoomManager {
	roomManager := new(RoomManager)
	roomManager.context = new(ServerContext)

	roomManager.context.Net = context.Net
	roomManager.context.Ipc = context.Ipc

	roomManager.listRoom = make([]*Room, 0)

	roomManager.logger = log.WithFields(log.Fields{"Server":"RoomManager of ZoneServer"})

	// roomManager.context.Net.RegisterRPCListener(roomManager)
	// roomManager.context.Ipc.RegisterRPC(roomManager)

	roomManager.logger.Info("New room manager successfully")
	return roomManager
}

// Clean :
func (roomManager *RoomManager)Clean() {
	roomManager.listRoom = make([]*Room, 0)
}

// Dump : show detail info for each room
func (roomManager *RoomManager)Dump() {
	for _,v := range roomManager.listRoom {
		v.Dump()
	}

	roomManager.logger.Info("rooms num :", len(roomManager.listRoom))
}

// GetRoom :
func (roomManager *RoomManager)GetRoom(roomID uint32) *Room {
	for _,v := range roomManager.listRoom {
		if v.Data.Id == roomID {
			return v
		}
	}
	return nil
}