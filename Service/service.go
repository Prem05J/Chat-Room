package Service

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/chatroom-go/Helper"
	"github.com/chatroom-go/Model"
	"github.com/gofiber/fiber/v2"
)

type ChatRoom struct {
	mu        sync.RWMutex
	clients   map[string]*Model.Client
	broadcast chan string
	join      chan *Model.Client
	leave     chan *Model.Client
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients:   make(map[string]*Model.Client),
		broadcast: make(chan string),
		join:      make(chan *Model.Client),
		leave:     make(chan *Model.Client),
	}
}

func (c *ChatRoom) Run() {

	for {
		select {
		case client := <-c.join:
			c.mu.Lock()
			c.clients[client.ID] = client
			log.Printf("Added clients: %+v\n", c.clients)
			c.mu.Unlock()
			log.Printf("Client %s Joined in chat Room", client.ID)
		case client := <-c.leave:
			c.mu.Lock()
			fmt.Print(c.clients)
			delete(c.clients, client.ID)
			close(client.Message)
			log.Printf("Latest clients: %+v\n", c.clients)
			c.mu.Unlock()
			log.Printf("Client %s Left from chat room", client.ID)

		case message := <-c.broadcast:
			c.mu.Lock()
			for _, client := range c.clients {
				select {
				case client.Message <- message:
				default:
				}
			}
			c.mu.Unlock()
		}

	}

}

func (s *ChatRoom) Join(r *fiber.Ctx) error {
	clientId := r.Query("id")

	if clientId == "" {
		Helper.WriteErrorJson(r, fiber.StatusBadRequest, "Client Id is empty")
	}
	client := &Model.Client{
		ID:      clientId,
		Message: make(chan string, 100),
	}

	s.join <- client

	return Helper.WriteFiberMap(r, fiber.StatusOK, "message", "Client Joined Successfully")

}

func (s *ChatRoom) SendMessage(r *fiber.Ctx) error {
	clientId := r.Query("id")
	message := r.Query("message")

	s.mu.RLock()
	_, isExits := s.clients[clientId]
	s.mu.RUnlock()

	if !isExits {
		return Helper.WriteFiberMap(r, fiber.StatusNotFound, "message", "Client Not found")
	}

	if clientId == "" || message == "" {
		return Helper.WriteErrorJson(r, fiber.ErrBadRequest.Code, "Client information is Empty")
	}

	compMessage := fmt.Sprintf("%s : %s", clientId, message)

	s.broadcast <- compMessage

	return Helper.WriteFiberMap(r, fiber.StatusOK, "message", "Message Sent Successfully")

}

func (s *ChatRoom) Leave(r *fiber.Ctx) error {
	clientId := r.Query("id")

	s.mu.RLock()
	client, isExits := s.clients[clientId]
	s.mu.RUnlock()

	if !isExits {
		return Helper.WriteFiberMap(r, fiber.StatusNotFound, "message", "Client Not found")
	}

	s.leave <- client

	return Helper.WriteFiberMap(r, fiber.StatusOK, "message", "Client Removeed")

}

func (s *ChatRoom) GetMessages(r *fiber.Ctx) error {
	clientId := r.Query("id")

	s.mu.RLock()
	client, isExist := s.clients[clientId]
	s.mu.RUnlock()

	if !isExist {
		return Helper.WriteFiberMap(r, fiber.StatusNotFound, "message", "Client Not found")
	}

	select {
	case mess := <-client.Message:
		return Helper.WriteFiberMap(r, fiber.StatusOK, "message", mess)
	case <-time.After(10 * time.Second):
		return Helper.WriteErrorJson(r, fiber.StatusNoContent, "No Content")
	}
}
