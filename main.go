package main

import (
	"code.holdonbush.top/FinalCheckmateServer/DataFormat"
	"code.holdonbush.top/FinalCheckmateServer/ServerDemo"
	_ "code.holdonbush.top/ServerFramework/Server"
	"fmt"
	"reflect"
	"time"
)

func main() {
	//logger.SetFormatter(&logger.TextFormatter{
	//	DisableColors: false,
	//	FullTimestamp: true,
	//})
	//logger.SetReportCaller(true)

	server := ServerDemo.NewServerDemo(8080)
	t := DataFormat.LoginMsg{}
	fmt.Println(reflect.TypeOf(t))
	for true {
		server.Tick()
		time.Sleep(1*time.Millisecond)
	}
}
