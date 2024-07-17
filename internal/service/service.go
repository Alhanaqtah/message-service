package service

import (
	"context"
	"log/slog"
	"message-service/internal/lib/logger/sl"
	"message-service/internal/models"
)

type Storage interface {
	SaveMessage(ctx context.Context, msg *models.Message) error
}

type Broker interface {
	ProduceMessage(ctx context.Context, msg *models.Message) error
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

	err := s.storage.SaveMessage(ctx, msg)
	if err != nil {
		log.Error("failed to save message to storage", sl.Error(err))
		return err
	}

	err = s.broker.ProduceMessage(ctx, msg)
	if err != nil {
		log.Error("failed to produce message to broker", sl.Error(err))
		return err
	}

	return nil
}
