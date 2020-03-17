package main

import (
	"code.holdonbush.top/FinalCheckmateServer/ServerDemo"
	"time"
)

func main() {
	server := ServerDemo.NewServerDemo(8080)

	for true {
		server.Tick()
		time.Sleep(1*time.Second)
	}
}
