package controller

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Service interface {
}

type Controller struct {
	log     *slog.Logger
	service Service
}

func New(log *slog.Logger, service Service) *Controller {
	return &Controller{
		log:     log,
		service: service,
	}
}

func (c *Controller) Register() func(r chi.Router) {
	return func(r chi.Router) {
		r.Post("/", c.getMessage)
		r.Get("/stats", c.getStats)
	}
}

func (c *Controller) getMessage(w http.ResponseWriter, r *http.Request) {

}

func (c *Controller) getStats(w http.ResponseWriter, r *http.Request) {

}
