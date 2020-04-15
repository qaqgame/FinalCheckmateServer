package testServer1

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/ServerManager"
	log "github.com/sirupsen/logrus"
)

type TestServer1 struct {
	*ServerManager.ServerModule
	Id      int
	Port    int
	Name    string
	trigger bool
}

func NewTestServer1(id, port int, name ...string) *TestServer1 {
	t := new(TestServer1)
	t.Id = id
	t.Port = port
	if len(name) >= 1 {
		t.Name = name[0]
	} else {
		t.Name = "TestServer1"
	}
	logger := log.WithFields(log.Fields{"Server": "TestServer1"})
	c := make(chan int, 2)
	t.trigger = false
	ipc := IPCWork.NewIPCManager(t.Id)
	status := ServerManager.UnCreated

	Info := ServerManager.ServerModuleInfo{
		Id:   t.Id,
		Name: t.Name,
		Port: t.Port,
	}
	t.ServerModule = ServerManager.NewServerModule(Info, logger, status, c, ipc)

	return t
}

func (server *TestServer1) Tick() {
	server.Logger.Info("Tick")
	if !server.trigger {
		server.trigger = true
		args := DataFormat.Args{
			Phase:  1,
			Phase2: "v",
		}
		reply := DataFormat.Reply{V: 0}
		server.Ipc.CallRpc(&args, &reply, 4051, "TestServer2.TestFunc")
		server.Logger.Info(reply)
	}

}
