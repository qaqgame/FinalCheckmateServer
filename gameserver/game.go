package gameserver

import "code.holdonbush.top/ServerFramework/fsplite"

type MyGameInstance struct {
	*fsplite.FSPGame
}

func NewMyGameInstance(_port int,gameid uint32) *MyGameInstance {
	myGameInstance := new(MyGameInstance)
	defaultparam := fsplite.NewDefaultFspParam("120.79.240.163", _port)
	myGameInstance.FSPGame = fsplite.NewFSPGame(gameid, defaultparam)

	return myGameInstance
}

func (mygame *MyGameInstance) OnStateRoundEnd() {

}