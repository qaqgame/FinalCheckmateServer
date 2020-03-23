package TestServer1

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/ServerManager"
	log "github.com/sirupsen/logrus"
	"time"
)

type TestServer1 struct {
	*ServerManager.ServerModule
	Id       int
	Port     int
	Name     string
	logger   *log.Entry
	created  bool
	status   int
	close    chan int
	trigger  bool
	ipc      *IPCWork.IPCManager
}

func NewTestServer1() *TestServer1 {
	t := new(TestServer1)
	t.Id = 1
	t.Port = 4050
	t.Name = "TestServer1"
	t.logger = log.WithFields(log.Fields{"Server":"TestServer1"})
	t.created = false
	t.close = make(chan int, 2)
	t.trigger = false
	t.ipc = IPCWork.NewIPCManager(t.Id)

	t.ServerModule = new(ServerManager.ServerModule)
	t.ServerModule.MInfo = ServerManager.ServerModuleInfo{
		Id:   t.Id,
		Name: t.Name,
		Port: t.Port,
	}

	return t
}

func (server *TestServer1)GetId() int {
	return server.Id
}
func (server *TestServer1)Create() {
	if server.status == ServerManager.UnCreated || server.status == ServerManager.Released {
		server.created = true
		server.status = ServerManager.Created
		server.logger.Info("Server Created")
	}
}
func (server *TestServer1)Release() {
	if server.status == ServerManager.Released {
		return
	}
	server.status = ServerManager.Released
	server.logger.Info("Server Released")
}
func (server *TestServer1)Start() {
	if server.status == ServerManager.Running {
		return
	}
	server.status = ServerManager.Running
	go func(server *TestServer1) {
		for true {
			select {
			case _ = <-server.close:
				return
			default:
				server.Tick()
				time.Sleep(time.Second)
			}
		}
	}(server)
	server.logger.Info("Server Started")
}
func (server *TestServer1)Stop() {
	if server.status == ServerManager.Stopped {
		return
	}
	server.status = ServerManager.Stopped
	server.close <- 1
	server.logger.Info("Server Stopped")
}
func (server *TestServer1)Tick() {
	server.logger.Info("Tick")
	if !server.trigger {
		server.trigger = true
		args := DataFormat.Args{
			Phase:  1,
			Phase2: "v",
		}
		reply := DataFormat.Reply{V:0}
		server.ipc.CallRpc(&args,&reply,4051,"TestServer2.TestFunc")
		server.logger.Info(reply)
	}

}
func (server *TestServer1)IsCreated() bool {
	return server.created
}

func (server *TestServer1)GetModuleInfo() ServerManager.ServerModuleInfo {
	return server.ServerModule.MInfo
}

func (server *TestServer1)GetStatus() int {
	return server.status
}