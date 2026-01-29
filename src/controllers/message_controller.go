package controllers

import (
	"context"
	"encoding/json"
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

	var limit int64 = 0
	var offset int64 = 0

	values, ok := r.URL.Query()["chatId"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"message": "parameter chatId is required!",
		})

		return
	}

	chatId := values[0]

	if chatId == "" {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]any{
			"message": "parameter chatId is required!",
		})

		return
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		parsed, err := strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]any{
				"message": "parameter limit must be a valid integer number!",
			})

			return
		}

		limit = parsed
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		parsed, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]any{
				"message": "parameter offset must be a valid integer number!",
			})

			return
		}

		offset = parsed
	}

	messages, err := c.service.Fetch(context.Background(), userID, chatId, limit, offset)
	if err != nil {
		http.Error(w, "There was an error while fetching the messages: " + err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"messages": messages,
	})
}
