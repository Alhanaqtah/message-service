package models

import "time"

type Stats struct {
	TotalMessages         int           `json:"total_messages"`
	ProcessedMessages     int           `json:"processed_messages"`
	AverageProcessingTime time.Duration `json:"average_processing_time"`
}
