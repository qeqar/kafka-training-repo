package main

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	brokers       = "172.18.0.3:9092" // add 172.18.0.3 kafka to your host file
	consumerGroup = "testConsumer"
	topic         = "test"
	kafka_ver     = "2.3.0"
)

func main() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	version, err := sarama.ParseKafkaVersion(kafka_ver)
	if err != nil {
		panic(err)
	}

	config := sarama.NewConfig()
	config.ClientID = "testConsumer"
	config.Version = version
	config.Consumer.MaxWaitTime = 10 * time.Second
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 20 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer := Consumer{
		ready: make(chan bool),
	}

	consumerGroup, err := sarama.NewConsumerGroup([]string{brokers}, consumerGroup, config)
	if err != nil {
		log.Panicf("Error creating Consumer %v", err)

	}

	go func(){
		consumerGroup.Consume(context.Background(), []string{topic}, &consumer)
		if err != nil {
			log.Panicf("consume %v", err)
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm // Await a sigterm signal before safely closing the consumer

	err = consumerGroup.Close()
	if err != nil {
		panic(err)
	}
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}

	return nil
}