package ws

import (
	"context"
	"encoding/json"
	"go_chat/src/services"
	"go_chat/src/types"
	"log"
)

type Hub struct {
	Clients        map[string]*Client
	Register       chan *Client
	Unregister     chan *Client
	Broadcast      chan *types.WebsocketMessage
	MessageService *services.MessageService
	ChatService    *services.ChatService
}

func NewHub(messageService *services.MessageService, chatService *services.ChatService) *Hub {
	return &Hub{
		Clients:        make(map[string]*Client),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		Broadcast:      make(chan *types.WebsocketMessage),
		MessageService: messageService,
		ChatService:    chatService,
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			h.Clients[client.UserID] = client

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
			}

		case payload := <-h.Broadcast:
			h.sendToChatMembers(payload)
		}
	}
}

func (h *Hub) sendToChatMembers(payload *types.WebsocketMessage) {
	users, err := h.ChatService.GetUsersID(context.Background(), payload.Message.Chat)
	if err != nil {
		return
	}

	if users == nil {
		log.Println("no users found")
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return
	}

	for _, userID := range users {
		if client, ok := h.Clients[userID]; ok {
			if client.UserID != payload.Message.Sender {
				client.Send <- jsonPayload
			}
		}
	}
}
