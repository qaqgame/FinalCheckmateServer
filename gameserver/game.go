package gameserver

import (
	"code.holdonbush.top/ServerFramework/fsplite"
)

type MyGameInstance struct {
	*fsplite.FSPGame
	APQueue     *fsplite.Queue
}

func NewMyGameInstance(_port int,gameid uint32) *MyGameInstance {
	myGameInstance := new(MyGameInstance)
	defaultparam := fsplite.NewDefaultFspParam("120.79.240.163", _port)
	myGameInstance.FSPGame = fsplite.NewFSPGame(gameid, defaultparam)

	myGameInstance.APQueue = fsplite.NewQueue()

	return myGameInstance
}

func (mygame *MyGameInstance) OnRoundBeginCallBack()  {

}

func (mygame *MyGameInstance) OnRoundEndCallBack() {

}