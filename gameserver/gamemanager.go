package gameserver

import (
	"errors"

	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
)

// GameManager :
type GameManager struct {
	context           *ServerContext
}

// NewGameManager :
func NewGameManager(_context *ServerContext) *GameManager {
	gamemanager := new(GameManager)
	gamemanager.context = _context

	gamemanager.context.Ipc.RegisterRPC(gamemanager)
	return gamemanager
}

// Clean :
func (gamemanager *GameManager) Clean() {
	if gamemanager.context != nil {
		gamemanager.context.Ipc.Clean()
		gamemanager.context = nil
	}
}


// RPCStartGame :
func (gamemanager *GameManager) RPCStartGame(args *DataFormat.CreateGame, reply *DataFormat.Reply) error {
	gamemanager.context.Fsp.CreateGame(args.RoomID)
	
	reply.P2S = gamemanager.context.Fsp.AddPlayers(args.RoomID,args.PlayerList)
	reply.Fspparam = gamemanager.context.Fsp.GetParam()

	if reply.Fspparam == nil {
		return errors.New("StartGame error: not Fspparam")
	}
	return nil
}