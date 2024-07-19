package main

import (
	"log/slog"
	"message-service/internal/config/config"
	"message-service/internal/controller"
	"message-service/internal/kafka/consumer"
	producer "message-service/internal/kafka/producer"
	"message-service/internal/lib/logger"
	"message-service/internal/lib/logger/sl"
	"message-service/internal/service"
	"message-service/internal/storage/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()

	log := logger.New(cfg.Env)

	log.Info("server running...")

	// Init producer
	producer, err := producer.New(cfg.Kafka)
	if err != nil {
		log.Error("failed to connect to broker", sl.Error(err))
		return
	}

	// Init & start consumer
	consumer, err := consumer.New(cfg.Kafka)
	if err != nil {
		log.Error("failed to connect to broker", sl.Error(err))
		return
	}

	err = consumer.Consume()
	if err != nil {
		log.Error("consumer error", sl.Error(err))
	}

	// Data layer
	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Error("failed to connect to storage", sl.Error(err))
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
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGABRT)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", slog.String("port", cfg.Server.Port), sl.Error(err))
		}
	}()

	log.Info("server is running...", slog.String("port", cfg.Server.Port))

	<-stop

	log.Info("stopping server...")

	srv.Close()
	storage.Close()
	producer.Close()
	consumer.Close()

	log.Info("server stopped")
}
