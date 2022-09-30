package main

import (
	"context"
	"fmt"

	app "github.com/athosone/projectraven/tracking/internal/application"
	"github.com/athosone/projectraven/tracking/internal/infrastructure"
	"github.com/athosone/projectraven/tracking/mongodb"
)

func NewApplication(ctx context.Context, appCfg *AppConfig) (*app.Application, error) {
	mongoCfg := mongodb.MongoDBConfig(appCfg.Database)
	mongodb.InitClient(ctx, &mongoCfg)
	userRepo, err := infrastructure.NewUserRepository(mongodb.Database)
	if err != nil {
		return nil, fmt.Errorf("wiring repo: %w", err)
	}
	newRegisterUserCommandHandler, err := app.NewRegisterUserCommandHandler(userRepo)
	if err != nil {
		return nil, fmt.Errorf("wiring handler: %w", err)
	}

	return &app.Application{
		Commands: app.Commands{
			UserCommands: app.UserCommands{
				RegisterUser: newRegisterUserCommandHandler,
			},
		},
	}, nil
}
