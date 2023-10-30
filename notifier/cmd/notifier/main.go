package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/patyukin/banking-system/notifier/internal/queue/kafka"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	kc, err := kafka.NewKafkaConsumer([]string{"localhost:9092"}, "my_topic", 0)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := kc.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	doneCh := make(chan struct{})

	go func() {
		for {
			select {
			case <-signals:
				doneCh <- struct{}{}

			default:
				msg, err := kc.ReadMessage()
				if err != nil {
					log.Println("No message available..")
				} else {
					fmt.Printf("Received message: offset-%d, key-%s, value-%s\n", msg.Offset, string(msg.Key), string(msg.Value))
				}
			}
		}
	}()

	log.Println("Consumer started...")
	<-doneCh
	log.Println("Consumer stopped")

	return nil
}
