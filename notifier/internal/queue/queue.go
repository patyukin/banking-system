package queue

import "github.com/IBM/sarama"

type Queue interface {
	Consume(string, chan sarama.ConsumerMessage)
}
