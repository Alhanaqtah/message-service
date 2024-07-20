package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Env string
	*Server
	*Storage
	*Kafka
}

type Server struct {
	Host    string
	Port    string
	Timeout time.Duration
}

type Storage struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

type Kafka struct {
	Brokers []string
	Topic   string
}

func MustLoad() *Config {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Panic("error loading .env file")
	// }

	timeout, err := strconv.Atoi(os.Getenv("SERVER_TIMEOUT"))
	if err != nil {
		log.Panicf(`error loading 'timeout' from .env file: %s`, err)
	}

	return &Config{
		os.Getenv("ENV"),
		&Server{
			Host:    os.Getenv("SERVER_HOST"),
			Port:    os.Getenv("SERVER_PORT"),
			Timeout: time.Duration(timeout),
		},
		&Storage{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
		},
		&Kafka{
			Brokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ", "),
			Topic:   os.Getenv("KAFKA_TOPIC"),
		},
	}
}
