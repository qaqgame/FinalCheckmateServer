package main

import (
	"code.holdonbush.top/FinalCheckmateServer/ServerDemo"
	_ "code.holdonbush.top/ServerFramework/Server"
	"time"
)

func init() {

}

func main() {
	server := ServerDemo.NewServerDemo(8080)
	for true {
		server.Tick()
		time.Sleep(1*time.Millisecond)
	}
}
