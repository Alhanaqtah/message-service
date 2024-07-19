package controller

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"message-service/internal/lib/logger/sl"
	"message-service/internal/lib/response"
	"message-service/internal/models"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Service interface {
	SaveMessage(ctx context.Context, msg *models.Message) error
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
		r.Post("/", c.saveMessage)
		r.Get("/stats", c.getStats)
	}
}

func (c *Controller) saveMessage(w http.ResponseWriter, r *http.Request) {
	const op = "controller.getMessage"

	log := c.log.With(
		slog.String("op", op),
		slog.String("req_id", middleware.GetReqID(r.Context())),
	)

	log.Debug("saving new message")

	var message models.Message
	err := render.Decode(r, &message)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Err("Invalid request body"))
		return
	}

	// Validation
	if message.Content == "" {
		log.Error("validation failed", sl.Error(errors.New("invalid message's content")))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Err("Invalid message's content"))
		return
	}

	err = c.service.SaveMessage(r.Context(), &message)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, response.Err("Internal error"))
		return
	}

	log.Debug("message saved succesfully")

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response.Ok("Message saved succesfully"))
}

func (c *Controller) getStats(w http.ResponseWriter, r *http.Request) {

}
