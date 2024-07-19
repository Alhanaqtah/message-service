package models

type Message struct {
	ID      string `json:"id,omitempty"`
	Content string `json:"content,omitempty"`
	Status  string `json:"status,omitempty"`
}
