package handlers

import (
	"context"
	"encoding/json"
	"dennic_user_service/internal/entity"
	"dennic_user_service/internal/infrastructure/kafka"
	"dennic_user_service/internal/pkg/config"
	"dennic_user_service/internal/usecase"
	"dennic_user_service/internal/usecase/event"

	"go.uber.org/zap"
)

type userCreateHandler struct {
	config         *config.Config
	brokerConsumer event.BrokerConsumer
	logger         *zap.Logger
	userUsecase    usecase.UserStorageI
}

func NewUserCreateHandler(config *config.Config,
	brokerConsumer event.BrokerConsumer,
	logger *zap.Logger,
	userUsecase usecase.UserStorageI) *userCreateHandler {
	return &userCreateHandler{
		config:         config,
		brokerConsumer: brokerConsumer,
		logger:         logger,
		userUsecase:    userUsecase,
	}
}

func (h *userCreateHandler) HandlerEvents() error {
	consumerConfig := kafka.NewConsumerConfig(
		h.config.Kafka.Address,
		"api.user.create",
		"1",
		func(ctx context.Context, key, value []byte) error {
			var user *entity.User

			if err := json.Unmarshal(value, &user); err != nil {
				return err
			}

			if _, err := h.userUsecase.Create(ctx, user); err != nil {
				return err
			}

			return nil
		},
	)

	h.brokerConsumer.RegisterConsumer(consumerConfig)
	h.brokerConsumer.Run()

	return nil

}
