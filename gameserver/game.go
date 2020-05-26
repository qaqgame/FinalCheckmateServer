package gameserver

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/ServerFramework/fsplite"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"sort"
)

type MyGameInstance struct {
	*fsplite.FSPGame
	APQueue     *fsplite.Queue
	ExeQueue    *fsplite.Queue
	Received    map[uint32]bool
	WinnerQueue *fsplite.Queue
	RPCCallerF   RPCCaller
}

type PlayerAP struct {
	playerId     uint32
	playerAP     int32
}

type RPCCaller func(args *DataFormat.CreateGame, reply *DataFormat.Reply) error

func (mygame *MyGameInstance) SetRPCCaller(caller RPCCaller) {
	mygame.RPCCallerF = caller
}

func NewMyGameInstance(_port int,gameid uint32) *MyGameInstance {
	myGameInstance := new(MyGameInstance)
	defaultparam := fsplite.NewDefaultFspParam("120.79.240.163", _port)
	myGameInstance.FSPGame = fsplite.NewFSPGame(gameid, defaultparam)

	myGameInstance.APQueue = fsplite.NewQueue()
	myGameInstance.ExeQueue = fsplite.NewQueue()
	myGameInstance.WinnerQueue = fsplite.NewQueue()
	myGameInstance.Received = make(map[uint32]bool)
	// set UpperController ---------- FSPGameI
	myGameInstance.FSPGame.UpperController = myGameInstance

	logrus.Info("port: ",_port,"Upper: ",myGameInstance.FSPGame.UpperController)
	return myGameInstance
}

func (mygame *MyGameInstance) Release() {
	creategame := new(DataFormat.CreateGame)
	creategame.RoomID = mygame.GetGameID()

	reply := new(DataFormat.Reply)
	ok := mygame.RPCCallerF(creategame,reply)
	if ok == nil {
		mygame.APQueue.Clear()
		mygame.WinnerQueue.Clear()
		mygame.ExeQueue.Clear()
		mygame.Received = nil
		mygame.FSPGame.Release()
	}
}

func (mygame *MyGameInstance) OnGameBeginCallBack(player *fsplite.FSPPlayer, message *fsplite.FSPMessage) {
	v := new(PlayerAP)
	v.playerId = player.IdInGame
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
	b = make([]byte,4)
	binary.BigEndian.PutUint32(b, uint32(0))
	if newBoolRound {
		binary.BigEndian.PutUint32(b, uint32(1))
	}
	return
}

func (mygame *MyGameInstance) CreateControlStartMsg() (b []byte) {
	b = make([]byte,4)
	playerIdInGame := mygame.ExeQueue.Pop().(*PlayerAP).playerId
	logrus.Info("ControlStartPID: ",playerIdInGame)
	binary.BigEndian.PutUint32(b, playerIdInGame)
	return
}

func (mygame *MyGameInstance) OnRoundEndCallBack(player *fsplite.FSPPlayer, message *fsplite.FSPMessage) {
	// 判断CMD
	if message.Cmd != fsplite.RoundEnd {
		return
	}

	playerap := new(PlayerAP)
	playerap.playerId = player.IdInGame
	//todo: 处理content
	ap := new(DataFormat.APPoint)
	err := proto.Unmarshal(message.Content, ap)
	playerap.playerAP = ap.AP
	//buf := bytes.NewBuffer(message.Content)
	//err := binary.Read(buf, binary.BigEndian, &playerap.playerAP)
	if err != nil {
		logrus.Fatal("error while Unmarshal Content")
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

func (mygame *MyGameInstance) OnGameEndCallBack(player *fsplite.FSPPlayer, message *fsplite.FSPMessage)  {
	v := new(DataFormat.UintCntMsg)
	err := proto.Unmarshal(message.Content, v)
	if err != nil {
		fmt.Println("err: ",err)
	}

	if !mygame.WinnerQueue.Contain(v.Value) {
		mygame.WinnerQueue.Push(v.Value)
	}
}

func (mygame *MyGameInstance) CreateGameEndMsg() (b []byte) {
	//winnerList := make([]*fsplite.FSPPlayer,0)
	//for mygame.WinnerQueue.Len() > 0 {
	//	IdInGame := mygame.WinnerQueue.Pop().(uint32)
	//	winnerList = append(winnerList, mygame.GetPlayerWithIdInGame(IdInGame))
	//}

	//base := winnerList[0].IdInGame
	//for _, v := range winnerList {
	//	if ((0x01 << base) & v.FriendMask) >> base != 1 {
	//		winnerIdList := new(DataFormat.UintListCntMsg)
	//		winnerIdList.Values = make([]uint32,0)
	//		res, err := proto.Marshal(winnerIdList)
	//		if err != nil {
	//			fmt.Println("error marshal: ",err)
	//		}
	//		return res
	//	}
	//}
	base := mygame.WinnerQueue.Pop().(uint32)
	basePlayer := mygame.GetPlayerWithIdInGame(base)
	for mygame.WinnerQueue.Len() > 0 {
		IdInGame := mygame.WinnerQueue.Pop().(uint32)
		item := mygame.GetPlayerWithIdInGame(IdInGame)
		if ((0x01<<base)&item.FriendMask)>>base != 1 {
			winnerIdList := new(DataFormat.UintListCntMsg)
			winnerIdList.Values = make([]uint32,0)
			res, err := proto.Marshal(winnerIdList)
			if err != nil {
				fmt.Println("error marshal: ",err)
			}
			return res
		}
	}

	baseFriendMask := basePlayer.FriendMask
	var tmpId uint32 = 0
	winnerIdList := new(DataFormat.UintListCntMsg)
	winnerIdList.Values = make([]uint32,0)
	for (baseFriendMask & (0x01 << tmpId)) >> tmpId == 1 {
		winnerIdList.Values = append(winnerIdList.Values, tmpId)
		tmpId++
	}
	b, err := proto.Marshal(winnerIdList)
	if err != nil {
		fmt.Println("error marshal: ",err)
	}
	return
}

func (mygame *MyGameInstance) OnGameEndMsgAddCallBack() {
	// mygame.Release()
}