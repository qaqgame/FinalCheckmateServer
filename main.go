package main

import (
	"code.holdonbush.top/FinalCheckmateServer/ServerDemo"
	"time"
)

func main() {
	server := ServerDemo.NewServerDemo()

	for true {
		server.Tick()
		time.Sleep(1*time.Second)
	}
}
