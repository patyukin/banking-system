package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type Consumer struct {
	Consumer          sarama.Consumer
	PartitionConsumer sarama.PartitionConsumer
}

func NewKafkaConsumer(brokers []string, topic string, partition int32) (*Consumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetOldest)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		Consumer:          consumer,
		PartitionConsumer: partitionConsumer,
	}, nil
}

func (kc *Consumer) Close() error {
	if err := kc.PartitionConsumer.Close(); err != nil {
		return err
	}

	if err := kc.Consumer.Close(); err != nil {
		return err
	}

	return nil
}

func (kc *Consumer) ReadMessage() (*sarama.ConsumerMessage, error) {
	select {
	case msg := <-kc.PartitionConsumer.Messages():
		return msg, nil
	default:
		return nil, fmt.Errorf("No message available")
	}
}
