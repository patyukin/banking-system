package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"sync"
)

// KafkaConsumer is the interface for Kafka consumer
type KafkaConsumer interface {
	Consume() error
	Close() error
}

// Consumer implements the KafkaConsumer interface
type Consumer struct {
	consumer sarama.Consumer
	messages chan *sarama.ConsumerMessage
	topic    string
	done     chan bool
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(brokers []string, topic string) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client, err: %w", err)
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer, err: %w", err)
	}

	messages := make(chan *sarama.ConsumerMessage)

	return &Consumer{
		consumer: consumer,
		messages: messages,
		topic:    topic,
		done:     make(chan bool),
	}, nil
}

// Consume starts consuming messages from Kafka topic
func (c *Consumer) Consume() error {
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return fmt.Errorf("failed to get partitions: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(partitions))
	for _, partition := range partitions {
		go func(partition int32) {
			defer wg.Done()
			consumer, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
			if err != nil {
				log.Println("Failed to start consumer for partition", partition, ":", err)
				return
			}

			defer func(consumer sarama.PartitionConsumer) {
				err = consumer.Close()
				if err != nil {
					log.Println("Failed to close consumer:", err)
				}
			}(consumer)

			for message := range consumer.Messages() {
				c.messages <- message
			}
		}(partition)
	}

	go func() {
		for message := range c.messages {
			processMessage(message)
		}
	}()

	wg.Wait()

	return nil
}

func processMessage(message *sarama.ConsumerMessage) {
	log.Printf("Received message: topic=%s, partition=%d, offset=%d, message=%s\n",
		message.Topic, message.Partition, message.Offset, string(message.Value))
}

// Close stops the consumer
func (c *Consumer) Close() error {
	close(c.messages)
	close(c.done)

	return nil
}
