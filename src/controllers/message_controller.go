package controllers

import (
	"context"
	"encoding/json"
	"go_chat/src/dto"
	"go_chat/src/services"
	"go_chat/src/types"
	"net/http"
	"strconv"
)

type MessageController struct {
	service *services.MessageService
}

func NewMessageController(service *services.MessageService) *MessageController {
	return &MessageController{service}
}

func (c *MessageController) Fetch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	const userContextKey types.ContextJWTClaimKey = "userID"

	userID, ok := r.Context().Value(userContextKey).(string)
	if !ok {
		return
	}

	var chatId string
	var limit int64 = 0
	var offset int64 = 0
	var errors []string

	values, ok := r.URL.Query()["chatId"]
	if ok {
		chatId = values[0]

		if chatId == "" {
			errors = append(errors, "parameter chatId is required!")
		}
	} else {
		errors = append(errors, "parameter chatId is required!")
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsed, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			errors = append(errors, "parameter limit must be a valid integer number!")
		} else {
			limit = parsed
		}
	} else {
		errors = append(errors, "parameter limit must be a valid integer number!")
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		parsed, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			errors = append(errors, "parameter offset must be a valid integer number!")
		} else {
			offset = parsed
		}
	} else {
		errors = append(errors, "parameter offset must be a valid integer number!")
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"message": errors,
		})

		return
	}

	messages, err := c.service.Fetch(context.Background(), userID, chatId, limit, offset)
	if err != nil {
		http.Error(w, "There was an error while fetching the messages: "+err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(dto.MessageResponseDTO{
		Messages: *messages,
		Total:    len(*messages),
		Offset:   offset,
		Limit:    min(100, limit),
	})
}
