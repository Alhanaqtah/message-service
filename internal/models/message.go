package models

import "time"

type Message struct {
	ID          string    `json:"id,omitempty"`
	Content     string    `json:"content,omitempty"`
	Topic       string    `json:"topic,omitempty"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	ProcessedAt time.Time `json:"processed_at,omitempty"`
}
