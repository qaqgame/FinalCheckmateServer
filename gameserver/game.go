package gameserver

import (
	"code.holdonbush.top/ServerFramework/fsplite"
	"fmt"
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

func (mygame *MyGameInstance) OnStateGameCreate()  {

}

func (mygame *MyGameInstance) OnStateGameBegin() {
	fmt.Println("Use user defined func")
	panic("test fault")
	return
}

func (mygame *MyGameInstance) OnRoundBeginCallBack()  {

}

func (mygame *MyGameInstance) OnRoundEndCallBack() {

}