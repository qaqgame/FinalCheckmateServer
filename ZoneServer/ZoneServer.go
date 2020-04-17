package ZoneServer

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"reflect"

	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/Server"
	"code.holdonbush.top/ServerFramework/ServerManager"
	log "github.com/sirupsen/logrus"
)

// ZoneServer : zoneserver
type ZoneServer struct {
	*ServerManager.ServerModule
	context       *ServerContext
	onlineManager *OnlineManager
	roomManager   *RoomManager
	rpcalls       *rpc.Call
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
	Info := ServerManager.ServerModuleInfo{
		Id:   id,
		Name: tname,
		Port: port,
	}
	c := make(chan int, 2)
	zoneServer.ServerModule = ServerManager.NewServerModule(Info, logger, ServerManager.UnCreated, c, ipcmanager)
	zoneServer.context = tcontext
	zoneServer.onlineManager = tonlineManager
	zoneServer.roomManager = NewRoomManager(tcontext)
	zoneServer.rpcalls = nil

	zoneServer.context.Net.RegisterRPCMethods(zoneServer, reflect.ValueOf(zoneServer), "UpdateRoomList", "CreateRoom", "JoinRoom", "ExitRoom", "RoomReady", "StartGame", "ChangeTeam", "UpdateRoomInfo")

	go zoneServer.ShowDump()
	// zoneServer.context.Ipc.RegisterRPC(zoneServer)
	zoneServer.DefaultCreateRoom()

	go zoneServer.OnFininshRpcall()
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
	} else {
		zoneServer.context.Net.ReturnError("room not exist", roomID)
	}
}

// StartGame :
func (zoneServer *ZoneServer) StartGame(session Server.ISession, roomID uint32) {
	// userID := session.GetUid()
	zoneServer.Logger.Info("Invoke RPC function: StartGame")
	room := zoneServer.roomManager.GetRoom(roomID)
	if room != nil {
		if room.CanStartGame() {
			// Create RPC Args
			creategame := new(DataFormat.CreateGame)
			creategame.PlayerList = make([]uint32, 0)
			creategame.RoomID = roomID
			creategame.AuthID = -2
			for _,v := range room.Data.Players {
				creategame.PlayerList = append(creategame.PlayerList, v.Uid)
			}

			// New a RPC Reply
			reply := new(DataFormat.Reply)
			ok := zoneServer.context.Ipc.CallRpc(&creategame, reply, 4051, "RPCStartGame")
			if ok == true {
				param := new(DataFormat.PVPStartParam)
				// TODO:
				var a interface{} = reply.Fspparam
				param.Fspparam = a.(*DataFormat.FSPParam)
				param.GameParam = room.GetGameParam()
				param.Players = room.Data.GetPlayers()

				// listSession := room.GetSessionList()
				for _,v := range param.Players {
					session := room.GetSeesion(v.GetUid())
					param.Fspparam.Sid = reply.P2S[v.GetUid()]
					zoneServer.context.Net.Invoke(session, "NotifyGameStart", param)
				}
				// zoneServer.context.Net.InvokeBroadCast(listSession, "NotifyGameStart", param)
			}
			// zoneServer.rpcalls = zoneServer.context.Ipc.CallRpcAsync(creategame, reply, 4051, "RPCStartGame")

			

			// invoke gameserver via IPC

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
