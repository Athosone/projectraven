package app

type Application struct {
	Commands Commands
}

type Commands struct {
	DeviceOwner  DeviceOwnerCommands
	UserCommands UserCommands
}

type DeviceOwnerCommands struct{}
type UserCommands struct {
	RegisterUser *RegisterUserCommandHandler
}

