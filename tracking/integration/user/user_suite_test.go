package user_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/athosone/golib/pkg/config"
	app "github.com/athosone/projectraven/tracking/internal/application"
	"github.com/athosone/projectraven/tracking/internal/domain"
	"github.com/athosone/projectraven/tracking/internal/infrastructure"
	"github.com/athosone/projectraven/tracking/mongodb"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

type testConfig struct {
	Database struct {
		ConnectionString string `yaml:"connectionString" json:"connectionString"`
		DatabaseName     string `yaml:"databaseName" json:"databaseName"`
	}
}

var testCfg *testConfig
var application *app.Application
var userRepo domain.UserRepository

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Suite")
}

var _ = BeforeSuite(func() {
	_ = viper.BindEnv("database.connectionString", "DATABASE_CONNECTION_STRING")
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "../../config/dev/"
	}
	fmt.Println("found configPath:", configPath)
	cfg, err := config.LoadConfig[testConfig](configPath)
	Expect(err).To(BeNil())
	testCfg = cfg
	fmt.Println(testCfg)

	application = NewApplication()

	Expect(mongodb.InitClient(context.Background(), (*mongodb.MongoDBConfig)(&testCfg.Database))).To(Succeed())
	By("Checking connectivity to the database")
	Expect(mongodb.Database.Client().Ping(context.TODO(), nil)).To(Succeed())

	r, err := infrastructure.NewUserRepository(mongodb.Database)
	Expect(err).To(BeNil())
	userRepo = r
})

func NewApplication() *app.Application {
	Expect(mongodb.InitClient(context.Background(), (*mongodb.MongoDBConfig)(&testCfg.Database))).To(Succeed())
	Expect(mongodb.Database.Client().Ping(context.TODO(), nil)).To(Succeed())

	r, err := infrastructure.NewUserRepository(mongodb.Database)
	Expect(err).To(BeNil())
	registerHandler, err := app.NewRegisterUserCommandHandler(r)
	userCommands := app.UserCommands{
		RegisterUser: registerHandler,
	}

	return &app.Application{
		Commands: app.Commands{
			UserCommands: userCommands,
		},
	}
}
