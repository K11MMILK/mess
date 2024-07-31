package kafka

import (
	"context"
	"log"
	"mess/internal/storage"
	"time"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	Ready   chan bool
	Storage *storage.Storage
}

func NewConsumer(storage *storage.Storage) *Consumer {
	return &Consumer{
		Ready:   make(chan bool),
		Storage: storage,
	}
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer setup")
	close(c.Ready)
	return nil
}

func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	log.Println("Consumer cleanup")
	return nil
}

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Println("Starting ConsumeClaim")
	for msg := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(msg.Value), msg.Timestamp, msg.Topic)

		// Обновление статуса сообщения на "processed"
		err := c.Storage.UpdateMessageStatus(string(msg.Value), "processed")
		if err != nil {
			log.Printf("Failed to update message status: %v", err)
		} else {
			log.Printf("Message status updated to 'processed': %s", string(msg.Value))
		}

		sess.MarkMessage(msg, "")
	}
	return nil
}

func StartConsumer(ctx context.Context, brokers []string, groupID, topic string, storage *storage.Storage) error {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Consumer.Return.Errors = true
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 20 * time.Second

	consumer := NewConsumer(storage)

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}
	defer client.Close()

	go func() {
		for {
			log.Println("Starting consumer loop")
			if err := client.Consume(ctx, []string{topic}, consumer); err != nil {
				log.Printf("Error from consumer: %v", err)
				if ctx.Err() != nil {
					log.Println("Context error, stopping consumer loop")
					return
				}
				time.Sleep(2 * time.Second)
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready
	log.Println("Consumer up and running!")

	go func() {
		for err := range client.Errors() {
			log.Printf("Error from consumer client: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Context cancelled, shutting down consumer...")

	return nil
}
