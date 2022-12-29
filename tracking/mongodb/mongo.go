package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	connectTimeout = 10 * time.Second
)

type MongoDBConfig struct {
	ConnectionString string `yaml:"connectionString" json:"connectionString"`
	DatabaseName     string `yaml:"databaseName" json:"databaseName"`
}

var client *mongo.Client
var Database *mongo.Database

func GetCollection(collectionName string) *mongo.Collection {
	return Database.Collection(collectionName)
}

func InitClient(ctx context.Context, cfg *MongoDBConfig) error {
	if client != nil {
		zap.S().Info("mongodb client already initialized")
		return nil
	}
	if client == nil {
		newClient, err := mongo.NewClient(options.Client().ApplyURI(cfg.ConnectionString))
		if err != nil {
			return fmt.Errorf("create mongo client: %w", err)
		}
		deadLine, cancel := context.WithTimeout(ctx, connectTimeout)
		defer cancel()
		if err := newClient.Connect(deadLine); err != nil {
			return fmt.Errorf("connect mongo client: %w", err)
		}
		client = newClient
		Database = client.Database(cfg.DatabaseName)
	}
	zap.S().Info("Mongo client created")
	return nil
}

func Shutdown(ctx context.Context) error {
	if client == nil {
		return nil
	}
	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("disconnect mongo client: %w", err)
	}
	return nil
}
