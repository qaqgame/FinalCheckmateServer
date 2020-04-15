module code.holdonbush.top/FinalCheckmateServer

require (
	code.holdonbush.top/ServerFramework v0.0.0
	github.com/golang/protobuf v1.4.0-rc.4
	github.com/sirupsen/logrus v1.4.2
	go.mongodb.org/mongo-driver v1.3.1
	google.golang.org/protobuf v1.20.1
)

replace code.holdonbush.top/ServerFramework => ../ServerFramework

go 1.13
