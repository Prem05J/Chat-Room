package handler

import (
	"github.com/chatroom-go/Service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service.ChatRoom
}

func NewHandler(ser *Service.ChatRoom) *Handler {
	return &Handler{
		service: ser,
	}
}

func (s *Handler) UnProtectedHandler(app *fiber.App) {
	app.Put("/join", s.Join)
	app.Put("/send", s.Send)
	app.Delete("/leave", s.Leave)
	app.Get("/message", s.Message)
}

func (s *Handler) Join(r *fiber.Ctx) error {
	return s.service.Join(r)
}

func (s *Handler) Send(r *fiber.Ctx) error {
	return s.service.SendMessage(r)
}

func (s *Handler) Leave(r *fiber.Ctx) error {
	return s.service.Leave(r)
}

func (s *Handler) Message(r *fiber.Ctx) error {
	return s.service.GetMessages(r)
}
