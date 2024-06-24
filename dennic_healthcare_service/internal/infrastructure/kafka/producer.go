package kafka

import (
	"Healthcare_Evrone/internal/entity"
	"Healthcare_Evrone/internal/pkg/config"
	"encoding/json"
	// "Healthcare_Evrone/internal/pkg/otlp"
	"context"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type producer struct {
	logger     *zap.Logger
	userCreate *kafka.Writer
}

func NewProducer(config *config.Config, logger *zap.Logger) *producer {
	return &producer{
		logger: logger,
		userCreate: &kafka.Writer{
			Addr:                   kafka.TCP(config.Kafka.Address...),
			Topic:                  config.Kafka.Topic.Healthcare,
			Balancer:               &kafka.Hash{},
			RequiredAcks:           kafka.RequireAll,
			AllowAutoTopicCreation: true,
			Async:                  true, // make the writer asynchronous
			Completion: func(messages []kafka.Message, err error) {
				if err != nil {
					logger.Error("kafka userCreate", zap.Error(err))
				}
				for _, message := range messages {
					logger.Sugar().Info(
						"kafka userCreate message",
						zap.Int("partition", message.Partition),
						zap.Int64("offset", message.Offset),
						zap.String("key", string(message.Key)),
						zap.String("value", string(message.Value)),
					)
				}
			},
		},
	}
}

func (p *producer) buildMessage(key string, value []byte) kafka.Message {
	return kafka.Message{
		Key:   []byte(key),
		Value: value,
	}
}

func (p *producer) ProduceUserToCreate(ctx context.Context, key string, value *entity.Doctor) error {

	byteValue, err := json.Marshal(&value)
	if err != nil {
		return err
	}
	message := p.buildMessage(key, byteValue)
	return p.userCreate.WriteMessages(ctx, message)
}

func (p *producer) Close() {
	if err := p.userCreate.Close(); err != nil {
		p.logger.Error("error during close writer articleCategoryCreated", zap.Error(err))
	}
}
