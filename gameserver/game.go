package gameserver

import "code.holdonbush.top/ServerFramework/fsplite"

type MyGameInstance struct {
	*fsplite.FSPGame
}

func NewMyGameInstance(gameid uint32, param *fsplite.FSPParam) *MyGameInstance {
	myGameInstance := new(MyGameInstance)
	myGameInstance.FSPGame = fsplite.NewFSPGame(gameid, param)

	return myGameInstance
}

func (mygame *MyGameInstance) OnStateRoundEnd() {

}