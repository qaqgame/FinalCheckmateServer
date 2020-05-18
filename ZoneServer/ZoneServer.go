package ZoneServer

import (
	"bufio"
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/Server"
	"code.holdonbush.top/ServerFramework/ServerManager"
	"code.holdonbush.top/ServerFramework/common"
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
	rpcalls       chan *rpc.Call
	timerStop     chan int
	timerIsRun    bool
	rpcMoniter    bool
}

// NewZoneServer :
func NewZoneServer(id, port int, name ...string) *ZoneServer {
	logger := log.WithFields(log.Fields{"Server": "ZoneServer"})
	zoneServer := new(ZoneServer)

	tname := "ZoneServer"
	if len(name) >= 1 {
		tname = name[0]
	}

	Info := new(common.ServerModuleInfo)
	Info.Id = id
	Info.Port = port
	Info.Name = tname
	c := make(chan int, 2)
	zoneServer.context = new(ServerContext)
	zoneServer.context.Net = Server.NewNetManager(port, logger)
	zoneServer.context.Ipc = IPCWork.NewIPCManager(Info)

	zoneServer.ServerModule = ServerManager.NewServerModule(Info, logger, ServerManager.UnCreated, c, zoneServer.context.Ipc)



	tonlineManager := NewOnlineManager(zoneServer.context)

	zoneServer.onlineManager = tonlineManager
	zoneServer.roomManager = NewRoomManager(zoneServer.context)
	zoneServer.rpcalls = make(chan *rpc.Call,1)
	zoneServer.timerStop = make(chan int, 5)
	zoneServer.timerIsRun = false
	zoneServer.rpcMoniter = false

	zoneServer.context.Net.RegisterRPCMethods(zoneServer, reflect.ValueOf(zoneServer), "UpdateRoomList", "CreateRoom", "JoinRoom", "ExitRoom", "RoomReady", "ChangeTeam", "UpdateRoomInfo")

	// go zoneServer.ShowDump()
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

// JoinRoom : Player join a specified room
func (zoneServer *ZoneServer) JoinRoom(session Server.ISession, roomID uint32) {
	zoneServer.Logger.Info("Invoke RPC function: JoinRoom")
	uid := session.GetUid()
	room := zoneServer.roomManager.GetRoom(roomID)
	if room != nil && !room.IsAllReady() {
		udcom := zoneServer.onlineManager.GetUserDataByID(uid)

		err := room.AddPlayer(uid, udcom.Userdata.GetName(), session)

		if err != nil {
			zoneServer.context.Net.ReturnError("room is full", roomID)
		}
		// listSession := room.GetSessionList()
		zoneServer.context.Net.Return(room.Data)
		// zoneServer.context.Net.InvokeBroadCast(listSession, "NotifyRoomUpdate", room.Data)
	} else {
		if room != nil && room.IsAllReady() {
			zoneServer.context.Net.ReturnError("room is in game", roomID)
		} else {
			zoneServer.context.Net.ReturnError("room not exist", roomID)
		}
	}
}

// ExitRoom : Player Exit Room
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


// ChangeTeam : Change Player's team
func (zoneServer *ZoneServer) ChangeTeam(session Server.ISession, roomID uint32, team uint32) {
	zoneServer.Logger.Info("Invoke RPC function: ChangeTeam")
	userID := session.GetUid()
	room := zoneServer.roomManager.GetRoom(roomID)
	if room != nil {
		for _, v := range room.Data.Players {
			if v.Uid == userID {
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

// OnFinishRPCCall : asynchronous rpc call
func (zoneServer *ZoneServer) OnFinishRPCCall(room *Room, playerTeamData *DataFormat.PlayerTeamData,replycall *rpc.Call) {
	zoneServer.rpcMoniter = true
	for true {
		//if zoneServer.rpcalls == nil {
		//	continue
		//}
		select {
		case replyCall := <- replycall.Done:
			fmt.Println("finish rpc")
			reply,_ := replyCall.Reply.(*DataFormat.Reply)
			// - we use synchronous rpc call in this server.
			for _, v := range room.Data.Players {
				session := room.GetSession(v.GetUid())
				fmt.Println("v: ",reply.P2S)
				reply.Fspparam.Sid = reply.P2S[v.GetUid()]
				zoneServer.Logger.Warn("NotifyGameStart: player id in game: ", v.Id, "session: ", session.GetUid())
				zoneServer.context.Net.Invoke(session, "NotifyGameStart", playerTeamData, v.Id, reply.Fspparam)
			}
			zoneServer.rpcMoniter = false
			return
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
				playerTeamData := new(DataFormat.PlayerTeamData)
				playerTeamData.Masks = make([]*DataFormat.MaskData, 0)

				for _,v := range room.Data.Players {
					maskData := new(DataFormat.MaskData)
					maskData.Pid = v.Id
					maskData.EnemyMask = ^(0x01 << v.Id)
					maskData.FriendMask = 0x01 << v.Id

					playerTeamData.Masks = append(playerTeamData.Masks, maskData)
				}

				//// send msg to client
				//maskData := new(DataFormat.MaskData)
				//// Pid: Id in Game
				//maskData.Pid = 1
				//maskData.EnemyMask = 0xff00
				//maskData.FriendMask = 0x00ff
				//
				//maskData2 := new(DataFormat.MaskData)
				//maskData2.Pid = 0
				//maskData2.EnemyMask = 0xff00
				//maskData2.FriendMask = 0x00ff
				//
				//maskData3 := new(DataFormat.MaskData)
				//maskData3.Pid = 2
				//maskData3.EnemyMask = 0xff00
				//maskData3.FriendMask = 0x00ff
				//
				//
				//playerTeamData.Masks = append(playerTeamData.Masks, maskData)
				//playerTeamData.Masks = append(playerTeamData.Masks, maskData2)
				//playerTeamData.Masks = append(playerTeamData.Masks, maskData3)

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
			creategame.PlayerList = make(map[uint32]uint32)
			creategame.RoomID = room.Data.Id
			creategame.AuthID = -2
			for _, v := range room.Data.Players {
				creategame.PlayerList[v.Uid] = v.Id
			}

			// New a RPC Reply
			reply := new(DataFormat.Reply)

			// start a synchronous rpc call
			ok := zoneServer.context.Ipc.CallRpc(creategame, reply, 4051, "GameManager.RPCStartGame")
			//replycall := zoneServer.context.Ipc.CallRpcAsync(creategame, reply, 4051, "GameManager.RPCStartGame")

			//if !zoneServer.rpcMoniter {
			//	go zoneServer.OnFinishRPCCall(room, playerTeamData,replycall)
			//}
			if ok == true {
				for _, v := range room.Data.Players {
					session := room.GetSession(v.GetUid())
					reply.Fspparam.Sid = reply.P2S[v.GetUid()]
					zoneServer.Logger.Warn("NotifyGameStart: player id in game: ", v.Id, "session: ", session.GetUid())
					zoneServer.context.Net.Invoke(session, "NotifyGameStart", playerTeamData, v.Id, reply.Fspparam)
				}
			} else {
				zoneServer.Logger.Error("RPC call RPCStartGame failed")
			}
		} else {
			zoneServer.context.Net.ReturnError("player unready", room.Data.Id)
		}
	}
}
