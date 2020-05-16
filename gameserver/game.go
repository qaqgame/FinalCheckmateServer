package gameserver

import (
	"bytes"
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/fsplite"
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"sort"
)

type MyGameInstance struct {
	*fsplite.FSPGame
	APQueue     *fsplite.Queue
	ExeQueue    *fsplite.Queue
	Received    map[uint32]bool
}

type PlayerAP struct {
	playerId     uint32
	playerAP     int32
}

func NewMyGameInstance(_port int,gameid uint32) *MyGameInstance {
	myGameInstance := new(MyGameInstance)
	defaultparam := fsplite.NewDefaultFspParam("120.79.240.163", _port)
	myGameInstance.FSPGame = fsplite.NewFSPGame(gameid, defaultparam)

	myGameInstance.APQueue = fsplite.NewQueue()
	myGameInstance.ExeQueue = fsplite.NewQueue()
	myGameInstance.Received = make(map[uint32]bool)
	// set UpperController ---------- FSPGameI
	myGameInstance.FSPGame.UpperController = myGameInstance

	return myGameInstance
}

func (mygame *MyGameInstance) OnGameBeginCallBack(player *fsplite.FSPPlayer, message *fsplite.FSPMessage) {
	v := new(PlayerAP)
	v.playerId = player.ID
	v.playerAP = DataFormat.DefaultAP

	if !mygame.Received[v.playerId] {
		mygame.Received[v.playerId] = true
		mygame.ExeQueue.Push(v)
	}
}

func (mygame *MyGameInstance) OnGameBeginMsgAddCallBack() {
	for k, _ := range mygame.Received {
		mygame.Received[k] = false
	}
}

func (mygame *MyGameInstance) CreateRoundMsg() (b []byte)  {
	newBoolRound := len(mygame.GetPlayerMap()) == mygame.ExeQueue.Len()
	binary.BigEndian.PutUint32(b, uint32(0))
	if newBoolRound {
		binary.BigEndian.PutUint32(b, uint32(1))
	}

	return
}

func (mygame *MyGameInstance) CreateControlStartMsg() (b []byte) {
	playerUId := mygame.ExeQueue.Pop().(*PlayerAP).playerId
	binary.BigEndian.PutUint32(b, playerUId)
	return
}

func (mygame *MyGameInstance) OnRoundEndCallBack(player *fsplite.FSPPlayer, message *fsplite.FSPMessage) {
	// 判断CMD
	if message.Cmd != fsplite.RoundEnd {
		return
	}

	playerap := new(PlayerAP)
	playerap.playerId = player.ID
	//todo: 处理content
	buf := bytes.NewBuffer(message.Content)
	err := binary.Read(buf, binary.BigEndian, &playerap.playerAP)
	if err != nil {
		logrus.Fatal("error while read int")
	}
	// playerap.playerAP = binary.BigEndian.Uint32(message.Content)
	if !mygame.Received[playerap.playerId] {
		mygame.Received[playerap.playerId] = true
		mygame.APQueue.Push(playerap)
	}

	// 处理信息
	if mygame.APQueue.Len() == len(mygame.GetPlayerMap()) {
		list := make([]*PlayerAP, len(mygame.GetPlayerMap()))
		for i := 0; i < len(list); i++ {
			list[i] = mygame.APQueue.Pop().(*PlayerAP)
		}
		sort.SliceStable(list, func(i, j int) bool {
			return list[i].playerAP > list[j].playerAP
		})
		for i := 0; i < len(list); i++ {
			mygame.ExeQueue.Push(list[i])
		}
	}

	// change FSPGameI's State.
	// Our Developer's Actions
	flag := mygame.GetFlag("roundEndFlag")
	for _, v := range mygame.GetPlayerMap() {
		mygame.SetFlag(v.IdInGame, flag, "roundEndFlag")
	}
}

func (mygame *MyGameInstance) OnRoundEndMsgAddCallBack() {
	for k, _ := range mygame.Received {
		mygame.Received[k] = false
	}
}