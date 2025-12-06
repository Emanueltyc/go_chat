package types

import "go_chat/src/models"

type WebsocketMessage struct {
	Name    string          `json:"name"`
	Message *models.Message `json:"message"`
}
