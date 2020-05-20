package DataFormat

import "code.holdonbush.top/ServerFramework/fsplite"

// LoginReq : info type
const (
	LoginReq         uint32 = 1
	LoginRes         uint32 = 2
	HeartBeatRequset uint32 = 3
	HeartBeatRsponse uint32 = 4
)

// Args : IPCWord using
type Args struct {
	Phase  int
	Phase2 string
}

// Reply : IPCWork using
type Reply struct {
	V        int
	Fspparam *fsplite.FSPParam
	P2S      map[uint32]uint32
}

// CreateGame :
type CreateGame struct {
	RoomID         uint32
	AuthID         int32
	PlayerList     map[uint32]uint32      //key: playerId   value: id in game
	MapFriendMask  map[uint32]uint32      //key: playerId   value: player's friend mask
	MapEnemyMask   map[uint32]uint32      //key: playerId   value: player's enemy mask
}
