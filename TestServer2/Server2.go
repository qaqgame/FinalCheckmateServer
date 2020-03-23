package TestServer2

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/ServerManager"
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

type TestServer2 struct {
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

func NewTestServer2() *TestServer2 {
	t := new(TestServer2)
	t.Id = 2
	t.Port = 40501
	t.Name = "TestServer2"
	t.logger = log.WithFields(log.Fields{"Server":"TestServer2"})
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
	t.ipc.RegisterRPC(t)
	return t
}

func (server *TestServer2)GetId() int {
	return server.Id
}
func (server *TestServer2)Create() {
	if server.status == ServerManager.UnCreated || server.status == ServerManager.Released {
		server.created = true
		server.status = ServerManager.Created
		server.logger.Info("Server Created")
	}
}
func (server *TestServer2)Release() {
	if server.status == ServerManager.Released {
		return
	}
	server.status = ServerManager.Released
	server.logger.Info("Server Released")
}
func (server *TestServer2)Start() {
	if server.status == ServerManager.Running {
		return
	}
	server.status = ServerManager.Running
	go func(server *TestServer2) {
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
func (server *TestServer2)Stop() {
	if server.status == ServerManager.Stopped {
		return
	}
	server.status = ServerManager.Stopped
	server.close <- 1
	server.logger.Info("Server Stopped")
}
func (server *TestServer2)Tick() {
	server.logger.Info("Tick")
}
func (server *TestServer2)IsCreated() bool {
	return server.created
}

func (server *TestServer2)GetModuleInfo() ServerManager.ServerModuleInfo {
	return server.ServerModule.MInfo
}

func (server *TestServer2)GetStatus() int {
	return server.status
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