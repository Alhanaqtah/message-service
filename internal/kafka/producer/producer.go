package producer

import (
	"message-service/internal/config/config"
	"message-service/internal/models"
	"time"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
	topic    string
}

func New(cfg *config.Kafka) (*Producer, error) {
	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer, topic: cfg.Topic}, nil
}

func (k *Producer) ProduceMessage(msg *models.Message) error {
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

func (k *Producer) Close() error {
	err := k.producer.Close()
	if err != nil {
		return err
	}

	return nil
}
