// global data's format between client and server

package DataFormat

import (
	"time"
)

const (
	// Timeout : mean for the time internal for checking client status
	Timeout int64 = 40
	// ZoneServer : zone server
	ZoneServer = 1
	// GameServer : game server
	GameServer = 2
)

var (
	// SuccessReturn : Success response
	SuccessReturn = ReturnCode{}
	// UnknownError : unknown error
	UnknownError = ReturnCode{Code: 1, Info: "UnknownError"}
)

// ComData : data communications between client and server
type ComData struct {
	Userdata       UserData
	Serveruserdata ServerUserData
	OnlineTimeout  int64
}

// ServerUserData : UserData info in Server
type ServerUserData struct {
	Sid               uint32 // stroe the session id
	LastHeartBeatTime int64
	IfOnline          bool
}

// CheckOnline : check client status(if online)
func (serveruserdata *ServerUserData) CheckOnline() bool {
	if serveruserdata.IfOnline {
		t := (time.Now().UnixNano()/int64(time.Millisecond)) - serveruserdata.LastHeartBeatTime
		if t > Timeout {
			serveruserdata.IfOnline = false
		}
	}
	return serveruserdata.IfOnline
}
