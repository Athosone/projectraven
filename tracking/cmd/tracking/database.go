package main

import (
	"context"
	"time"

	"github.com/athosone/projectraven/tracking/internal/config"
	"github.com/athosone/projectraven/tracking/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

func newMongoDB(cfg *config.AppConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mCfg := mongodb.MongoDBConfig(cfg.MongoDB)

	err := mongodb.InitClient(ctx, &mCfg)
	if err != nil {
		return nil, err
	}
	return mongodb.Database, nil
}
