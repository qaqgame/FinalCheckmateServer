package ZoneServer

import (
	"reflect"
	"time"

	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/Server"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

// OnlineManager : online function module
type OnlineManager struct {
	mapUserData map[uint32]*DataFormat.ComData
	netmanager  *Server.NetManager
	logger      *log.Entry
}

// NewOnlineManager : new a online manager
func NewOnlineManager(context *ServerContext) *OnlineManager {
	onlineManager := new(OnlineManager)

	onlineManager.netmanager = context.Net
	onlineManager.netmanager.SetAuthCmd(DataFormat.LoginReq)
	onlineManager.netmanager.AddListener(DataFormat.LoginReq, onlineManager.OnLoginRequest, new(DataFormat.LoginMsg))
	onlineManager.netmanager.AddListener(DataFormat.HeartBeatRequset, onlineManager.OnHeartBeatRequest, new(DataFormat.HeartBeatReq))
	onlineManager.netmanager.RegisterRPCMethod(onlineManager, reflect.ValueOf(onlineManager), "Logout")

	onlineManager.mapUserData = make(map[uint32]*DataFormat.ComData)
	onlineManager.logger = log.WithFields(log.Fields{"Server": "ZoneServer OnlineManager"})

	return onlineManager
}

// Clean : clean
func (onlineManager *OnlineManager) Clean() {
	onlineManager.mapUserData = make(map[uint32]*DataFormat.ComData)
	onlineManager.netmanager = nil
}

// Dump :
func (onlineManager *OnlineManager) Dump() {
	for _, v := range onlineManager.mapUserData {
		onlineManager.logger.Info("user id :", v.Userdata.Id, "user name :", v.Userdata.Name)
	}
	onlineManager.logger.Info("num of player :", len(onlineManager.mapUserData))
}

// OnLoginRequest :
func (onlineManager *OnlineManager) OnLoginRequest(session Server.ISession, index uint32, tmsg proto.Message) {
	onlineManager.logger.Info("invoke into OnLoginRequest")

	success := false
	loginmsg := tmsg.(*DataFormat.LoginMsg)

	onlineManager.logger.Info("Func OnLoginResquest, loginmsg: ", loginmsg)
	onlineManager.logger.Info("login message id: ", loginmsg.Uid)
	// Get a DataFormat.ComData type info
	ud := onlineManager.GetUserDataByName(loginmsg.Name)
	if ud == nil {
		// create a DataFormat.ComData type info
		ud = onlineManager.CreateUserData(session.GetId(), loginmsg.Name)
		activeUser(session, ud)
		onlineManager.logger.Info("login new user. user id: ", ud.Userdata.Id, " session.uid", session.GetUid())
		success = true
	} else {
		if loginmsg.Uid == ud.Userdata.Id {
			// 重连
			activeUser(session, ud)
			onlineManager.logger.Info("reconnect. user id: ", ud.Userdata.Id, " session.uid", session.GetUid())
			success = true
		} else {
			// todo : - 考虑是否需要该功能 - 抢占别人的名字
			if !ud.Serveruserdata.IfOnline {
				activeUser(session, ud)
				onlineManager.logger.Info("Hijack other's name: user id", ud.Userdata.Id, " session.uid", session.GetUid())
				success = true
			}

		}
	}
	if success {
		response := new(DataFormat.LoginRsp)
		response.Ret = &DataFormat.SuccessReturn
		response.Userdata = &ud.Userdata
		onlineManager.logger.Debug("return message: ", response)
		onlineManager.netmanager.Send(session, index, DataFormat.LoginRes, response)
	} else {
		response := new(DataFormat.LoginRsp)
		response.Ret = new(DataFormat.ReturnCode)
		response.Ret.Code = 1
		response.Ret.Info = "登陆失败，名字被占用"
		onlineManager.netmanager.Send(session, index, DataFormat.LoginRes, response)
	}
}

func activeUser(session Server.ISession, userdata *DataFormat.ComData) {
	userdata.Serveruserdata.LastHeartBeatTime = time.Now().UnixNano() / int64(time.Millisecond)
	userdata.Serveruserdata.Sid = session.GetId()
	session.SetAuth(userdata.Userdata.Id)
}

// OnHeartBeatRequest :
func (onlineManager *OnlineManager) OnHeartBeatRequest(session Server.ISession, index uint32, tmsg proto.Message) {
	onlineManager.logger.Warn("invoke into OnHeartBeatRequest", session.GetUid())
	heartbeatreq := tmsg.(*DataFormat.HeartBeatReq)
	onlineManager.logger.Warn("HeartBeatReq info: ", heartbeatreq, " ", time.Now().UnixNano()/int64(time.Millisecond))
	ud := onlineManager.GetUserDataByID(session.GetUid())
	if ud != nil {
		ud.Serveruserdata.LastHeartBeatTime = time.Now().UnixNano() / int64(time.Millisecond)
		session.SetPing(uint32(heartbeatreq.Ping))
		heartres := new(DataFormat.HeartBeatRsp)
		heartres.Ret = &DataFormat.SuccessReturn
		heartres.Timestamp = heartbeatreq.Timestamp
		onlineManager.logger.Debug("heartbeat: ", heartres)
		onlineManager.netmanager.Send(session, index, DataFormat.HeartBeatRsponse, heartres)
	} else {
		onlineManager.logger.Info("找不到session 对应的 UserData, session :", session)
	}
}

// Logout : logout from server
func (onlineManager *OnlineManager) Logout(session Server.ISession) {
	onlineManager.logger.Info("client request to logout")

	onlineManager.ReleaseUserData(session.GetUid())
	onlineManager.netmanager.Return()
}

// CreateUserData : create ComData type info
func (onlineManager *OnlineManager) CreateUserData(id uint32, name string) *DataFormat.ComData {
	comdata := new(DataFormat.ComData)
	comdata.Userdata.Name = name
	comdata.Userdata.Id = id
	onlineManager.mapUserData[id] = comdata
	return comdata
}

// ReleaseUserData : release info in map
func (onlineManager *OnlineManager) ReleaseUserData(id uint32) {
	delete(onlineManager.mapUserData, id)
}

// GetUserDataByID :
func (onlineManager *OnlineManager) GetUserDataByID(id uint32) *DataFormat.ComData {
	return onlineManager.mapUserData[id]
}

// GetUserDataByName :
func (onlineManager *OnlineManager) GetUserDataByName(name string) *DataFormat.ComData {
	for _, v := range onlineManager.mapUserData {
		if v.Userdata.Name == name {
			return v
		}
	}
	return nil
}
