// Package kafka Хелпер для работы с кафкой
package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
	"log"
)

var _ KafkaConsumerInterface = (*KafkaConsumer)(nil)

type KafkaConsumerInterface interface {
	RunConsume(ctx context.Context, handlerFunc func(ctx context.Context, key string, value string) error) error
	Close() error
}

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	topic    string
}

func NewConsumer(brokerList []string, topic string) (*KafkaConsumer, error) {

	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}

	// Create consumer group
	kafkaConsumerGroup := topic + "-consumer-group"
	consumerGroup, err := sarama.NewConsumerGroup(brokerList, kafkaConsumerGroup, config)
	if err != nil {
		return &KafkaConsumer{}, errors.Wrap(err, "Starting consumer group")
	}

	kafkaConsumer := &KafkaConsumer{
		consumer: consumerGroup,
		topic:    topic,
	}

	return kafkaConsumer, nil
}

func (c *KafkaConsumer) RunConsume(ctx context.Context, handlerFunc func(ctx context.Context, key string, value string) error) error {
	consumerGroupHandler := Consumer{
		handlerFunc: handlerFunc,
	}

	err := c.consumer.Consume(ctx, []string{c.topic}, &consumerGroupHandler)
	if err != nil {
		return fmt.Errorf("failed to consume: %w", err)
	}

	return nil
}

func (c *KafkaConsumer) Close() error {
	err := c.consumer.Close()
	if err != nil {
		return err
	}

	return nil
}

// Consumer represents a Sarama consumer group consumer.
type Consumer struct {
	handlerFunc func(ctx context.Context, key string, value string) error
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("consumer - setup")
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("consumer - cleanup")
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		err := consumer.handlerFunc(session.Context(), string(message.Key), string(message.Value))
		if err == nil {
			session.MarkMessage(message, "")
		}
	}

	return nil
}
