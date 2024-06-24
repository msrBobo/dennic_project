package app

import (
	"dennic_user_service/internal/delivery/grpc/kafka/handlers"
	"dennic_user_service/internal/infrastructure/kafka"
	"dennic_user_service/internal/infrastructure/repository/postgresql/user"
	"dennic_user_service/internal/pkg/config"
	logPkg "dennic_user_service/internal/pkg/logger"
	"dennic_user_service/internal/pkg/postgres"
	"dennic_user_service/internal/usecase"
	"dennic_user_service/internal/usecase/event"
	"fmt"

	"go.uber.org/zap"
)

type UserCreateConsumerCLI struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	BrokerConsumer event.BrokerConsumer
}

func NewUserCreateConsumerCLI(config *config.Config, logger *zap.Logger, db *postgres.PostgresDB, brokerConsumer event.BrokerConsumer) (*UserCreateConsumerCLI, error) {
	logger, err := logPkg.New(config.LogLevel, config.Environment, config.APP+"_cli"+".log")
	if err != nil {
		return nil, err
	}

	consumer := kafka.NewConsumer(logger)

	db, err = postgres.New(config)
	if err != nil {
		return nil, err
	}

	return &UserCreateConsumerCLI{
		Config:         config,
		DB:             db,
		Logger:         logger,
		BrokerConsumer: consumer,
	}, nil
}

func (c *UserCreateConsumerCLI) Run() error {
	fmt.Print("consume is running ....")
	// repo init
	userRepo := postgresql.NewUserRepo(c.DB)

	// usecase init
	userUsecase := usecase.NewUserService(c.DB.Config().ConnConfig.ConnectTimeout, userRepo)

	eventHandler := handlers.NewUserCreateHandler(c.Config, c.BrokerConsumer, c.Logger, userUsecase)

	return eventHandler.HandlerEvents()
}

func (c *UserCreateConsumerCLI) Close() {
	c.BrokerConsumer.Close()

	c.Logger.Sync()
}
