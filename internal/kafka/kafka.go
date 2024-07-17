package kafka

import (
	"context"
	"message-service/internal/config/config"
	"message-service/internal/models"
	"time"

	"github.com/IBM/sarama"
)

type Kafka struct {
	producer sarama.SyncProducer
	topic    string
}

func New(cfg *config.Kafka) (*Kafka, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &Kafka{producer: producer, topic: cfg.Topic}, nil
}

func (k *Kafka) ProduceMessage(ctx context.Context, msg *models.Message) error {
	m := &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.StringEncoder(msg.Content),
	}

	_, _, err := k.producer.SendMessage(m)
	if err != nil {
		return err
	}

	return nil
}

func (k *Kafka) Close() error {
	err := k.producer.Close()
	if err != nil {
		return err
	}

	return nil
}
