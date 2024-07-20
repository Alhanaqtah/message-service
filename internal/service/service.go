package service

import (
	"context"
	"log/slog"
	"time"

	"message-service/internal/lib/logger/sl"
	"message-service/internal/models"
)

type Storage interface {
	SaveMessage(ctx context.Context, msg *models.Message) (*models.Message, error)
	MarkMessageAsProcessed(ctx context.Context, id string) error
	MarkMessageAsFailed(ctx context.Context, id string) error
	FetchStats(ctx context.Context) (*models.Stats, error)
}

type Broker interface {
	ProduceMessage(msg *models.Message) error
}

type Service struct {
	log     *slog.Logger
	storage Storage
	broker  Broker
}

func New(log *slog.Logger, storage Storage, broker Broker) *Service {
	return &Service{
		log:     log,
		storage: storage,
		broker:  broker,
	}
}

func (s *Service) SaveMessage(ctx context.Context, msg *models.Message) error {
	const op = "service.SaveMessage"

	log := s.log.With(slog.String("op", op))

	savedMsg, err := s.storage.SaveMessage(ctx, msg)
	if err != nil {
		log.Error("failed to save message to storage", sl.Error(err))
		return err
	}

	err = s.broker.ProduceMessage(savedMsg)
	if err != nil {
		log.Error("failed to produce message to broker", sl.Error(err))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = s.storage.MarkMessageAsFailed(ctx, msg.ID)
		if err != nil {
			log.Error(`failed to mark message as 'failed'`, sl.Error(err))
			return err
		}

		return err
	}

	return nil
}

func (s *Service) FetchStats(ctx context.Context) (*models.Stats, error) {
	const op = "service.FetchStats"

	log := s.log.With(slog.String("op", op))

	stats, err := s.storage.FetchStats(ctx)
	if err != nil {
		log.Error("failed to fetch statisctics", sl.Error(err))
		return nil, err
	}

	return stats, nil
}
