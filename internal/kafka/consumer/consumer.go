package consumer

import (
	"log"
	"message-service/internal/config/config"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer sarama.Consumer
	topic    string
}

func New(cfg *config.Kafka) (*Consumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{consumer: consumer, topic: cfg.Topic}, nil
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
				log.Printf("Consumed message: %s", string(message.Value))
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
