package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"message-service/internal/config/config"
	"message-service/internal/controller"
	"message-service/internal/kafka/consumer"
	producer "message-service/internal/kafka/producer"
	"message-service/internal/lib/logger"
	"message-service/internal/lib/logger/sl"
	"message-service/internal/service"
	"message-service/internal/storage/postgres"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	log.Info("server running...")

	// Data layer
	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Error("failed to connect to storage", sl.Error(err))
		return
	}

	// Init producer
	producer, err := producer.New(cfg.Kafka)
	if err != nil {
		log.Error("failed to connect to broker", sl.Error(err))
		return
	}

	// Init consumer
	consumer, err := consumer.New(cfg.Kafka, storage)
	if err != nil {
		log.Error("failed to connect to broker", sl.Error(err))
		return
	}

	// Service layr
	service := service.New(log, storage, producer)

	// Controller layer
	controller := controller.New(log, service)

	// Init router
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Route("/messages", controller.Register())

	srv := http.Server{
		Addr:        cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:     r,
		IdleTimeout: cfg.Server.Timeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// Starting server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", slog.String("port", cfg.Server.Port), sl.Error(err))
			os.Exit(1)
		}
	}()

	// Starting consumer
	go func() {
		if err := consumer.Consume(); err != nil {
			log.Error("failed to start consumer", sl.Error(err))
			os.Exit(1)
		}
	}()

	log.Info("server is running...", slog.String("port", cfg.Server.Port))

	<-stop

	log.Info("stopping server...")

	// Graceful shutdown
	storage.Close()

	if err := srv.Close(); err != nil {
		log.Error("failed to close server", sl.Error(err))
	}

	if err := producer.Close(); err != nil {
		log.Error("failed to close producer", sl.Error(err))
	}

	if err := consumer.Close(); err != nil {
		log.Error("failed to close consumer", sl.Error(err))
	}

	log.Info("server stopped")
}
