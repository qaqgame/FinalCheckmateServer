package ZoneServer

import (
	"bufio"
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/Server"
	"code.holdonbush.top/ServerFramework/ServerManager"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/rpc"
	"os"
	"reflect"
	"sort"
	"time"
)

// ZoneServer : zoneserver
type ZoneServer struct {
	*ServerManager.ServerModule
	context       *ServerContext
	onlineManager *OnlineManager
	roomManager   *RoomManager
	rpcalls       *rpc.Call
	timerStop     chan int
	timerIsRun    bool
}

// NewZoneServer :
func NewZoneServer(id, port int, name ...string) *ZoneServer {
	logger := log.WithFields(log.Fields{"Server": "ZoneServer"})
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
	Info := new(ServerManager.ServerModuleInfo)
	Info.Id = id
	Info.Port = port
	Info.Name = tname
	c := make(chan int, 2)
	zoneServer.ServerModule = ServerManager.NewServerModule(Info, logger, ServerManager.UnCreated, c, ipcmanager)
	zoneServer.context = tcontext
	zoneServer.onlineManager = tonlineManager
	zoneServer.roomManager = NewRoomManager(tcontext)
	zoneServer.rpcalls = nil
	zoneServer.timerStop = make(chan int, 5)
	zoneServer.timerIsRun = false

	zoneServer.context.Net.RegisterRPCMethods(zoneServer, reflect.ValueOf(zoneServer), "UpdateRoomList", "CreateRoom", "JoinRoom", "ExitRoom", "RoomReady", "ChangeTeam", "UpdateRoomInfo")

	go zoneServer.ShowDump()
	// zoneServer.context.Ipc.RegisterRPC(zoneServer)
	zoneServer.DefaultCreateRoom()

	return zoneServer
}

// DefaultCreateRoom :
func (zoneServer *ZoneServer) DefaultCreateRoom() {
	room := new(Room)
	room.mapUserID2Session = make(map[uint32]Server.ISession)
	room.Data = new(DataFormat.RoomData)
	room.Data.Id = NewRoomID()
	room.Data.Name = "default"
	room.Data.Mode = "none"
	room.Data.MapName = "none"
	room.Data.Team = []int32{5, 5}
	room.Data.Maxplayercount = 0
	for _, v := range room.Data.Team {
		room.Data.Maxplayercount += v
	}

	zoneServer.roomManager.listRoom = append(zoneServer.roomManager.listRoom, room)
}

// Stop : override base Stop
func (zoneServer *ZoneServer) Stop() {
	if zoneServer.context.Net != nil {
		zoneServer.context.Net.UnRegisterRPCListener(zoneServer)
		zoneServer.context.Net.Clean()
		zoneServer.context.Net = nil
	}

	if zoneServer.context.Ipc != nil {
		zoneServer.context.Ipc.Clean()
		zoneServer.context.Ipc = nil
	}
	zoneServer.onlineManager.Clean()
	zoneServer.roomManager.Clean()
}

// Tick : tick
func (zoneServer *ZoneServer) Tick() {
	// zoneServer.Logger.Info("new server tick")
	zoneServer.context.Net.Tick()
}

// ShowDump :
func (zoneServer *ZoneServer) ShowDump() {
	for true {
		zoneServer.Logger.Debug("test input")
		reader := bufio.NewReader(os.Stdin)
		str, _ := reader.ReadString('\n')
		if str[:len(str)-2] == "1" {
			zoneServer.Logger.Info("invoke onlinemanager Dump")
			zoneServer.onlineManager.Dump()
		} else if str[:len(str)-2] == "2" {
			zoneServer.Logger.Info("invoke Gateway Dump")
			zoneServer.context.Net.Gateway.Dump()
		} else if str[:len(str)-2] == "3" {
			zoneServer.Logger.Info("show rooms")
			zoneServer.roomManager.Dump()
		} else {
			zoneServer.Logger.Info("input value error")
		}
	}
}

// RPC

// UpdateRoomList :
func (zoneServer *ZoneServer) UpdateRoomList(session Server.ISession) {
	zoneServer.Logger.Info("Invoke RPC function: UpdateRoomList")
	list := make([]*DataFormat.RoomData, 0)
	for _, v := range zoneServer.roomManager.listRoom {
		list = append(list, v.Data)
	}
	data := new(DataFormat.RoomListData)
	data.Rooms = list

	zoneServer.context.Net.Return(data)
}

// CreateRoom :
func (zoneServer *ZoneServer) CreateRoom(session Server.ISession, roomName, mapName, modeName string, teams []int32) {
	zoneServer.Logger.Info("Invoke RPC function: CrateRoom")
	userID := session.GetUid()
	udcom := zoneServer.onlineManager.GetUserDataByID(userID)

	room := CreateRoom(userID, udcom.Userdata.GetName(), session, roomName, modeName, mapName, teams)
	zoneServer.roomManager.listRoom = append(zoneServer.roomManager.listRoom, room)

	zoneServer.context.Net.Return(room.Data)
}

// JoinRoom :
func (zoneServer *ZoneServer) JoinRoom(session Server.ISession, roomID uint32) {
	zoneServer.Logger.Info("Invoke RPC function: JoinRoom")
	uid := session.GetUid()
	room := zoneServer.roomManager.GetRoom(roomID)
	if room != nil {
		udcom := zoneServer.onlineManager.GetUserDataByID(uid)

		err := room.AddPlayer(uid, udcom.Userdata.GetName(), session)

		if err != nil {
			zoneServer.context.Net.ReturnError("room is full", roomID)
		}
		// listSession := room.GetSessionList()
		zoneServer.context.Net.Return(room.Data)
		// zoneServer.context.Net.InvokeBroadCast(listSession, "NotifyRoomUpdate", room.Data)
	} else {
		zoneServer.context.Net.ReturnError("room not exist", roomID)
	}
}

// ExitRoom :
func (zoneServer *ZoneServer) ExitRoom(session Server.ISession, roomID uint32) {
	zoneServer.Logger.Info("Invoke RPC function: ExitRoom")
	userID := session.GetUid()
	room := zoneServer.roomManager.GetRoom(roomID)
	if room != nil {
		room.RemovePlayer(userID)
		if room.GetPlayerCount() > 0 {
			listSession := room.GetSessionList()
			zoneServer.context.Net.InvokeBroadCast(listSession, "NotifyRoomUpdate", room.Data)
		}
	}
}

// RoomReady :
func (zoneServer *ZoneServer) RoomReady(session Server.ISession, roomID uint32, ready bool) {
	zoneServer.Logger.Info("Invoke RPC function: RoomReady")
	userID := session.GetUid()
	room := zoneServer.roomManager.GetRoom(roomID)
	if room != nil {
		fmt.Println("status: ", ready)
		room.SetReady(userID, ready)
		listSession := room.GetSessionList()
		zoneServer.context.Net.InvokeBroadCast(listSession, "NotifyRoomUpdate", room.Data)

		if room.IsAllReady() {
			// 开始倒计时。
			room.Data.Ready = true
			if !zoneServer.timerIsRun {
				go zoneServer.Timer(listSession, room)
			}
		} else {
			room.Data.Ready = false
			if zoneServer.timerIsRun {
				zoneServer.timerStop <- 1
			}
		}

	} else {
		zoneServer.context.Net.ReturnError("room not exist", roomID)
	}
}

// StartGame :
func (zoneServer *ZoneServer) StartGame(session Server.ISession, roomID uint32) {
	zoneServer.Logger.Info("Invoke RPC function: StartGame")
	//room := zoneServer.roomManager.GetRoom(roomID)
}

// StartGame1 :
func (zoneServer *ZoneServer) StartGame1(session Server.ISession, roomID uint32) {
	// userID := session.GetUid()
	zoneServer.Logger.Info("Invoke RPC function: StartGame1")
	room := zoneServer.roomManager.GetRoom(roomID)
	if room != nil {
		if room.CanStartGame() {
			// Create RPC Args
			creategame := new(DataFormat.CreateGame)
			creategame.PlayerList = make([]uint32, 0)
			creategame.RoomID = roomID
			creategame.AuthID = -2
			for _, v := range room.Data.Players {
				creategame.PlayerList = append(creategame.PlayerList, v.Uid)
			}

			// New a RPC Reply
			reply := new(DataFormat.Reply)
			ok := zoneServer.context.Ipc.CallRpc(&creategame, reply, 4051, "RPCStartGame")
			if ok == true {
				param := new(DataFormat.PVPStartParam)
				//
				var a interface{} = reply.Fspparam
				param.Fspparam = a.(*DataFormat.FSPParam)
				param.GameParam = room.GetGameParam()
				param.Players = room.Data.GetPlayers()
				for _, v := range param.Players {
					session := room.GetSeesion(v.GetUid())
					param.Fspparam.Sid = reply.P2S[v.GetUid()]
					zoneServer.context.Net.Invoke(session, "NotifyGameStart", param)
				}
			}

		} else {
			zoneServer.context.Net.ReturnError("player unready", roomID)
		}
	} else {
		zoneServer.context.Net.ReturnError("room not exist", roomID)
	}
}

// ChangeTeam :
func (zoneServer *ZoneServer) ChangeTeam(session Server.ISession, roomID uint32, team uint32) {
	zoneServer.Logger.Info("Invoke RPC function: ChangeTeam")
	userID := session.GetUid()
	room := zoneServer.roomManager.GetRoom(roomID)
	if room != nil {
		for _, v := range room.Data.Players {
			if v.Id == userID {
				v.Teamid = team
				break
			}
		}
		sort.Sort(room)
		for i := 0; i < len(room.Data.GetPlayers()); i++ {
			room.Data.Players[i].Id = uint32(i + 1)
		}
		listSession := room.GetSessionList()
		zoneServer.context.Net.InvokeBroadCast(listSession, "NotifyRoomUpdate", room.Data)
	}
}

// UpdateRoomInfo :
func (zoneServer *ZoneServer) UpdateRoomInfo(session Server.ISession, roomID uint32) {
	room := zoneServer.roomManager.GetRoom(roomID)
	if room != nil {
		listSession := room.GetSessionList()

		zoneServer.context.Net.InvokeBroadCast(listSession, "NotifyRoomUpdate", room.Data)
	} else {
		zoneServer.context.Net.ReturnError("room not exist", roomID)
	}
}

// OnFininshRpcall :
func (zoneServer *ZoneServer) OnFininshRpcall() {
	for true {
		if zoneServer.rpcalls == nil {
			continue
		}
		select {
		case replyCall := <-zoneServer.rpcalls.Done:
			// TODO: sync rpc
			_ = replyCall.Reply.(*DataFormat.Reply)
		}
	}
}

// Timer : Assistant Function - invoke client's update each second.
func (zoneServer *ZoneServer) Timer(listSession []Server.ISession, room *Room) {
	zoneServer.timerIsRun = true
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	var count int32 = 3
	for {
		select {
		case <-tick.C:
			room.Data.Time = count
			zoneServer.Logger.Warn("Session len： ",len(listSession))
			for _,v := range listSession {
				zoneServer.Logger.Warn("session info: ", v.GetId(), v.GetUid(), v.IsActive())
			}
			// invoke client's update
			zoneServer.context.Net.InvokeBroadCast(listSession, "NotifyRoomUpdate", room.Data)
			if count == 0 {
				// send msg to client
				maskData := new(DataFormat.MaskData)
				maskData.Pid = 1
				maskData.EnemyMask = 0x00ff
				maskData.FriendMask = 0xff00

				maskData2 := new(DataFormat.MaskData)
				maskData2.Pid = 0
				maskData2.EnemyMask = 0x00ff
				maskData2.FriendMask = 0xff00

				playerTeamData := new(DataFormat.PlayerTeamData)
				playerTeamData.Masks = make([]*DataFormat.MaskData, 0)
				playerTeamData.Masks = append(playerTeamData.Masks, maskData)
				playerTeamData.Masks = append(playerTeamData.Masks, maskData2)

				// for _, v := range listSession {
				// 	zoneServer.context.Net.Invoke(v, "NotifyGameStart", playerTeamData, v.GetUid())
				// }

				// start fsp server
				zoneServer.startFspServer(room, playerTeamData)

				zoneServer.timerIsRun = false
				return
			}
			count--
		case <-zoneServer.timerStop:
			// invoke client's update
			zoneServer.context.Net.InvokeBroadCast(listSession, "NotifyRoomUpdate", room.Data)
			zoneServer.timerIsRun = false
			return
		}
	}
}

func (zoneServer *ZoneServer) startFspServer(room *Room, playerTeamData *DataFormat.PlayerTeamData) {
	if room != nil {
		if room.CanStartGame() {
			// Create RPC Args
			creategame := new(DataFormat.CreateGame)
			creategame.PlayerList = make([]uint32, 0)
			creategame.RoomID = room.Data.Id
			creategame.AuthID = -2
			for _, v := range room.Data.Players {
				creategame.PlayerList = append(creategame.PlayerList, v.Uid)
			}

			// New a RPC Reply
			reply := new(DataFormat.Reply)
			ok := zoneServer.context.Ipc.CallRpc(creategame, reply, 4051, "GameManager.RPCStartGame")
			if ok == true {
				param := new(DataFormat.PVPStartParam)

				// var a interface{} = reply.Fspparam
				// param.Fspparam = (*DataFormat.FSPParam)(unsafe.Pointer(reply.Fspparam))
				// param.GameParam = room.GetGameParam()
				param.Players = room.Data.GetPlayers()
				// todo: 玩家id分配位置
				var idingame uint32 = 0
				for _, v := range param.Players {
					session := room.GetSeesion(v.GetUid())
					// param.Fspparam.Sid = reply.P2S[v.GetUid()]
					reply.Fspparam.Sid = reply.P2S[v.GetUid()]
					// v.Id = idingame
					// zoneServer.Logger.Warn("NotifyGameStart target: ", v.Name, session.GetRemoteEndPoint())
					zoneServer.Logger.Warn("NotifyGameStart: player id in game: ", v.Id, "session: ", session.GetUid())
					zoneServer.context.Net.Invoke(session, "NotifyGameStart", playerTeamData, v.Id, reply.Fspparam)
					idingame++
				}
			} else {
				zoneServer.Logger.Error("RPC call RPCStartGame failed")
			}
		} else {
			zoneServer.context.Net.ReturnError("player unready", room.Data.Id)
		}
	}
}
