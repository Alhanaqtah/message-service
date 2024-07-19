package consumer

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"message-service/internal/config/config"

	"github.com/IBM/sarama"
)

type Storage interface {
	MarkMessageAsProcessed(ctx context.Context, id string) error
}

type Consumer struct {
	consumer sarama.Consumer
	topic    string
	storage  Storage
}

func New(cfg *config.Kafka, storage Storage) (*Consumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: consumer, topic: cfg.Topic, storage: storage}, nil
}

func (c *Consumer) Consume() error {
	partitionList, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return err
	}

	for _, partition := range partitionList {
		pc, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			for message := range pc.Messages() {
				log.Printf("Consumed message: %s\n", message.Value)
				err = c.storage.MarkMessageAsProcessed(context.Background(), string(message.Key[:]))
			}
		}(pc)
	}

	// Wait for termination signal to exit
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	return c.Close()
}

func (c *Consumer) Close() error {
	err := c.consumer.Close()
	if err != nil {
		return err
	}

	return nil
}
