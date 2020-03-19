package ServerDemo

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/Server"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

type ServerDemo struct {
	netManager *Server.NetManager
	logger     *log.Entry
}

func NewServerDemo(port int) *ServerDemo {
	loggerInfo := log.WithFields(log.Fields{"ServerName":"ServerDemo"})
	loggerInfo.Info("Server Created")
	serverDemo := new(ServerDemo)
	serverDemo.logger = loggerInfo
	serverDemo.netManager = Server.NewNetManager(port,loggerInfo)
	serverDemo.netManager.SetAuthCmd(DataFormat.LoginReq)
	serverDemo.netManager.AddListener(DataFormat.LoginReq, serverDemo.OnLogin, new(DataFormat.LoginMsg))

	// log.Println("create serverdemo")
	//logger.WithFields(logger.Fields{
	//
	//}).Info("ServerDemo started")
	return serverDemo
}

func (serverDemo *ServerDemo) Tick() {
	// log.Println("tick")
	serverDemo.netManager.Tick()
}

func (serverDemo *ServerDemo) OnLogin(session Server.ISession, index uint32, tmsg proto.Message) {
	//log.Println("ServerDemo OnLogin",reflect.ValueOf(tmsg))
	t := tmsg.(*DataFormat.LoginMsg)
	//if t == nil {
	//	log.Println("data format error: ",t)
	//}
	//log.Println("id ",t.Uid," name ", t.Name)
	//log.Println(time.Now().Unix(),"received")
	serverDemo.logger.Debug("OnLogin of ServerDemo, uid",t.Uid,"name",t.Name)
	res := DataFormat.LoginRsp{}
	res.Ret = session.GetId()
	res.Msg = "success"

	serverDemo.netManager.Send(session,index,DataFormat.LoginRes,&res)
}
