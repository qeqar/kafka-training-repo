package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

const (
	brokers = "172.18.0.3:9092" // add 172.18.0.3 kafka to your host file
	topic = "test"
	message = "test message"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Timeout = 10
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer([]string{brokers}, config)
	if err != nil {
		log.Panicf("Error creating Producer %v", err)
	}

	msg := &sarama.ProducerMessage {
        Topic: topic,
        Value: sarama.StringEncoder(message),
    }
    partition, offset, err := producer.SendMessage(msg)
    if err != nil {
        fmt.Println("Error publish: ", err.Error())
    }

    fmt.Println("Partition: ", partition)
    fmt.Println("Offset: ", offset)
}