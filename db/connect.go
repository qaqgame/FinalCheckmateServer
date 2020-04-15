package db

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Connection :
type Connection struct {
	clientOptions *options.ClientOptions
	client        *mongo.Client
	logger        *logrus.Entry
}

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Process": "Connecting into Mongodb"}).Error("Connect error: ", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{"Process": "Connecting into Mongodb"}).Error("check the connection: ", err)
	}

	logrus.WithFields(logrus.Fields{"Process": "Connecting into Mongodb"}).Info("Connection success")

	err = client.Disconnect(context.TODO())
	if err != nil {
		logrus.WithFields(logrus.Fields{"Process": "Connecting into Mongodb"}).Error("Close Connect error: ", err)
	}

	logrus.WithFields(logrus.Fields{"Process": "Connecting into Mongodb"}).Info("close connect success")
}

// NewConnection :
func NewConnection() *Connection {
	var err error
	connection := new(Connection)
	connection.logger = logrus.WithFields(logrus.Fields{"Process": "Connection of Mongodb"})
	connection.clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")
	connection.client, err = mongo.Connect(context.TODO(), connection.clientOptions)
	if err != nil {
		connection.logger.Error("Connect error: ", err)
	}

	err = connection.client.Ping(context.TODO(), nil)
	if err != nil {
		connection.logger.Error("check the connection: ", err)
	}

	connection.logger.Info("Connection success")

	return connection
}

// CloseConnection :
func (connection *Connection) CloseConnection() {
	err := connection.client.Disconnect(context.TODO())
	if err != nil {
		connection.logger.Error("Close Connect error: ", err)
	}

	connection.logger.Info("close connect success")
}
