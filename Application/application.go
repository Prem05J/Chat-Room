package Application

import (
	"log"

	handler "github.com/chatroom-go/Handler"
	"github.com/chatroom-go/Service"
	"github.com/gofiber/fiber/v2"
)

func Run() error {
	app := fiber.New()
	service := Service.NewChatRoom()
	go service.Run()
	handler := handler.NewHandler(service)
	handler.UnProtectedHandler(app)
	log.Println("Listening-On", ":8080")
	return app.Listen(":8080")
}
