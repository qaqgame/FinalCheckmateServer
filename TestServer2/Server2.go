package TestServer2

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/ServerManager"
	"fmt"
	log "github.com/sirupsen/logrus"
)

type TestServer2 struct {
	*ServerManager.ServerModule
	Id       int
	Port     int
	Name     string
}

func NewTestServer2(id,port int, name ...string) *TestServer2 {
	t := new(TestServer2)
	t.Id = id
	t.Port = port
	if len(name) >= 1 {
		t.Name = name[0]
	} else {
		t.Name = "TestServer2"
	}
	logger := log.WithFields(log.Fields{"Server":"TestServer2"})
	c := make(chan int, 2)
	ipc := IPCWork.NewIPCManager(t.Id)
	status := ServerManager.UnCreated

	info := ServerManager.ServerModuleInfo{
		Id:   t.Id,
		Name: t.Name,
		Port: t.Port,
	}
	t.ServerModule = ServerManager.NewServerModule(info,logger,status,c,ipc)
	t.Ipc.RegisterRPC(t)
	return t
}

func (server *TestServer2)TestFunc(args *DataFormat.Args, reply *DataFormat.Reply) error {
	fmt.Println(args)
	reply.V = 3
	return nil
}

func (server *TestServer2)TestFunc1(args *DataFormat.Args, reply *DataFormat.Reply) error {
	fmt.Println(args)
	reply.V = 3
	return nil
}