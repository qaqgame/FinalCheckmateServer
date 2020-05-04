package ZoneServer

import (
	"errors"
	"sort"
	"strconv"

	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/Server"
	log "github.com/sirupsen/logrus"
)

var lastRoomID uint32

func init() {
	lastRoomID = 0
}

// NewRoomID : Allocate a new room id
func NewRoomID() uint32 {
	lastRoomID++
	return lastRoomID
}

// Room : room
type Room struct {
	Data              *DataFormat.RoomData
	mapUserID2Session map[uint32]Server.ISession
	logger            *log.Entry
}

// Dump : show connect infos
func (room *Room) Dump() {
	room.logger.Info("room id: ", room.Data.Id, " room name :", room.Data.Name, " player_counts: ", len(room.Data.GetPlayers()))
	for _, v := range room.Data.Players {
		room.logger.Info("player id: ", v.Id, " player name: ", v.Name)
	}
	for k, v := range room.mapUserID2Session {
		room.logger.Info("user id: ", k, " session id", v.GetId(), " session active: ", v.IsActive())
	}
}

// CreateRoom : create a room
func CreateRoom(userID uint32, userName string, session Server.ISession, roomName, modeName, mapName string, teams []int32) *Room {
	room := new(Room)
	room.mapUserID2Session = make(map[uint32]Server.ISession)
	room.Data = new(DataFormat.RoomData)
	room.Data.Id = NewRoomID()
	room.Data.Name = roomName
	room.Data.Mode = modeName
	room.Data.MapName = mapName
	room.Data.Team = teams
	room.Data.Maxplayercount = 0
	for _, v := range teams {
		room.Data.Maxplayercount += v
	}
	// max player count; team  --- in room.Data

	room.logger = log.WithFields(log.Fields{"Server": "Room" + strconv.Itoa(int(room.Data.Id))})

	room.AddPlayer(userID, userName, session)

	return room
}

// AddPlayer : add a player to room
func (room *Room) AddPlayer(userID uint32, userName string, session Server.ISession) error {
	playerData := room.GetPlayerDataByUserID(userID)
	if playerData == nil {
		playerData = new(DataFormat.PlayerData)
		// allocate id in game. here we use the lenght of palyerdata
		playerData.Id = uint32(len(room.Data.GetPlayers()))
		// data.Sid = session.GetId()
		if len(room.Data.Players) < int(room.Data.Maxplayercount) {
			room.Data.Players = append(room.Data.Players, playerData)
		} else {
			return errors.New("fulled")
		}
	}
	sort.Sort(room)
	for i := 0; i < len(room.Data.GetPlayers()); i++ {
		room.Data.Players[i].Id = uint32(i)
	}

	// playerData = room.GetPlayerDataByUserID(userID)

	playerData.IsReady = false
	playerData.Uid = userID
	playerData.Name = userName

	room.mapUserID2Session[userID] = session
	return nil
}

// RemovePlayer : remove a player
func (room *Room) RemovePlayer(userID uint32) {
	i := room.GetPlayerIndexByUserID(userID)
	if i >= 0 {
		room.Data.Players = append(room.Data.Players[0:i], room.Data.Players[i+1:]...)
	}

	delete(room.mapUserID2Session, userID)
}

// GetPlayerDataByUserID :
func (room *Room) GetPlayerDataByUserID(userID uint32) *DataFormat.PlayerData {
	for _, v := range room.Data.Players {
		if v.Uid == userID {
			return v
		}
	}
	return nil
}

// GetPlayerIndexByUserID :
func (room *Room) GetPlayerIndexByUserID(userID uint32) int {
	for k, v := range room.Data.Players {
		if v.Uid == userID {
			return k
		}
	}
	return -1
}

// GetPlayerCount :
func (room *Room) GetPlayerCount() int {
	return len(room.Data.GetPlayers())
}

// GetSessionList :
func (room *Room) GetSessionList() []Server.ISession {
	list := make([]Server.ISession, 0)
	for _, v := range room.Data.Players {
		list = append(list, room.mapUserID2Session[v.Uid])
	}
	return list
}

// GetSeesion :
func (room *Room) GetSeesion(userID uint32) Server.ISession {
	return room.mapUserID2Session[userID]
}

// CanStartGame :
func (room *Room) CanStartGame() bool {
	if len(room.Data.GetPlayers()) > 0 && room.IsAllReady() {
		return true
	}
	return false
}

// IsAllReady :
func (room *Room) IsAllReady() bool {
	isAllready := true
	for _, v := range room.Data.Players {
		if !v.IsReady {
			isAllready = false
			break
		}
	}

	return isAllready
}

// SetReady :
func (room *Room) SetReady(userID uint32, value bool) {
	info := room.GetPlayerDataByUserID(userID)
	if info != nil {
		info.IsReady = value
	}
}

// GetGameParam :
func (room *Room) GetGameParam() *DataFormat.GameParam {
	param := new(DataFormat.GameParam)
	param.Id = room.Data.GetId()
	param.LimitedTime = DataFormat.LimiedTime
	param.Randseed = 0
	param.Gamemode = DataFormat.GameModes_TimeLimitPvP

	param.Map.Id = 1
	param.Map.Name = "default map"

	return param
}

// Sort

// Len :
func (room *Room) Len() int { return len(room.Data.Players) }

// Swap :
func (room *Room) Swap(i, j int) {
	room.Data.Players[i], room.Data.Players[j] = room.Data.Players[j], room.Data.Players[i]
}

// Less :
func (room *Room) Less(i, j int) bool {
	if room.Data.Players[i].Teamid != room.Data.Players[j].Teamid {
		return room.Data.Players[i].Teamid < room.Data.Players[j].Teamid
	}
	return room.Data.Players[i].Id < room.Data.Players[j].Id
}
