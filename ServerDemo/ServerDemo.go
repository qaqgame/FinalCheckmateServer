package ServerDemo

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/Server"
	"fmt"
	"github.com/golang/protobuf/proto"
	logger "github.com/sirupsen/logrus"
	"log"
	"reflect"
	"time"
)

type ServerDemo struct {
	netManager *Server.NetManager
}

func NewServerDemo(port int) *ServerDemo {
	serverDemo := new(ServerDemo)
	serverDemo.netManager = Server.NewNetManager(port)
	serverDemo.netManager.SetAuthCmd(DataFormat.LoginReq)
	serverDemo.netManager.AddListener(DataFormat.LoginReq, serverDemo.OnLogin, new(DataFormat.LoginMsg))

	log.Println("create serverdemo")
	logger.Info("create serverdemo")
	return serverDemo
}

func (serverDemo *ServerDemo) Tick() {
	// log.Println("tick")
	serverDemo.netManager.Tick()
}

func (serverDemo *ServerDemo) OnLogin(session Server.ISession, index uint32, tmsg proto.Message) {
	log.Println("ServerDemo OnLogin",reflect.ValueOf(tmsg))
	t := tmsg.(*DataFormat.LoginMsg)
	fmt.Println("tfttt",t)
	//if t == nil {
	//	log.Println("data format error: ",t)
	//}
	log.Println("id ",t.Uid," name ", t.Name)
	log.Println(time.Now().Unix(),"received")
	res := DataFormat.LoginRsp{}
	res.Ret = session.GetId()
	res.Msg = "success"

	serverDemo.netManager.Send(session,index,DataFormat.LoginRes,&res)
}
