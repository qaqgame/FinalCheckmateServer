package TestServer2

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/IPCWork"
	"code.holdonbush.top/ServerFramework/ServerManager"
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
func (server *TestServer2)Create(info ServerManager.ServerModuleInfo) {
	server.created = true
	server.status = DataFormat.Created
	server.logger.Info("Server Created")
}
func (server *TestServer2)Release() {
	server.status = DataFormat.Released
	server.logger.Info("Server Released")
}
func (server *TestServer2)Start() {

	server.status = DataFormat.Started
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
	server.status = DataFormat.Stopped
	server.close <- 1
	server.logger.Info("Server Stopped")
}
func (server *TestServer2)IsCreated() bool {
	return server.created
}

func (server *TestServer2)GetModuleInfo() ServerManager.ServerModuleInfo {
	return server.ServerModule.MInfo
}

func (server *TestServer2) TestFunc(args *DataFormat.Args, reply *DataFormat.Reply) error {
	server.logger.Info(args.Phase," ",args.Phase2)
	reply.V = 2
	return nil
}