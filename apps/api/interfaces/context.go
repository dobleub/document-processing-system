package interfaces

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type contextKey string

const (
	ConfigKey  contextKey = "config"
	LoggerKey  contextKey = "logger"
	MongodbKey contextKey = "mongodb"
	MailerKey  contextKey = "mailer"
	StateKey   contextKey = "state"
	McpClient  contextKey = "mcpClient"
)

// ConnectMongoDB connects to the MongoDB database
type MongoDBContext struct {
	Client *mongo.Client
	DB     *mongo.Database
}

type LoggerContext struct {
	Log *zap.Logger
}
