package ServerDemo

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/Server"
	"github.com/golang/protobuf/proto"
	"log"
)

type ServerDemo struct {
	netManager *Server.NetManager
}

func NewServerDemo(port int) *ServerDemo {
	serverDemo := new(ServerDemo)
	serverDemo.netManager = Server.NewNetManager(port)
	serverDemo.netManager.SetAuthCmd(DataFormat.LoginReq)
	serverDemo.netManager.AddListener(DataFormat.LoginReq, serverDemo.OnLogin, &DataFormat.LoginMsg{})

	log.Println("create serverdemo")
	return serverDemo
}

func (serverDemo *ServerDemo) Tick() {
	log.Println("tick")
	serverDemo.netManager.Tick()
}

func (serverDemo *ServerDemo) OnLogin(session Server.ISession, index uint32, tmsg proto.Message) {
	err := tmsg.(*DataFormat.LoginMsg)
	if err != nil {
		log.Println("data format error")
	}
	log.Println("id ",tmsg.(*DataFormat.LoginMsg).Id," name ", tmsg.(*DataFormat.LoginMsg).Name)

	res := DataFormat.LoginRsp{}
	res.Ret = 0
	res.Msg = "success"

	serverDemo.netManager.Send(session,index,DataFormat.LoginRes,&res)
}
