package gameserver

import (
	"errors"
	"fmt"

	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
)

// GameManager :
type GameManager struct {
	context           *ServerContext
	port              int
}

// NewGameManager :
func NewGameManager(_port int,_context *ServerContext) *GameManager {
	gamemanager := new(GameManager)

	gamemanager.context = _context
	gamemanager.port = _port

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
	// gamemanager.context.Fsp.CreateGameI(args.RoomID)
	myGameInstance := NewMyGameInstance(gamemanager.port,args.RoomID)
	gamemanager.context.Fsp.AddUDefinedGame(myGameInstance)

	//key: playerId   value: id in game
	fmt.Println("playerlist: ", args.PlayerList)

	reply.P2S = gamemanager.context.Fsp.AddPlayers(args.RoomID,args.PlayerList, args.MapFriendMask, args.MapEnemyMask)
	reply.Fspparam = gamemanager.context.Fsp.GetParam()

	if reply.Fspparam == nil {
		return errors.New("StartGame error: not Fspparam")
	}
	return nil
}