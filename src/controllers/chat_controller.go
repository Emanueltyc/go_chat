package controllers

import (
	"context"
	"encoding/json"
	"go_chat/src/models"
	"go_chat/src/services"
	"go_chat/src/types"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatController struct {
	service *services.ChatService
}

type ChatRequest struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func NewChatController(service *services.ChatService) *ChatController {
	return &ChatController{service}
}

func (c *ChatController) AccessChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	receiverID := strings.Trim(r.URL.Query().Get("receiverID"), "")

	if receiverID == "" {
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]any{
			"status":  "error",
			"message": "receiverID param not sent with request",
		})

		return
	}

	const userContextKey types.ContextJWTClaimKey = "userID"

	userID, ok := r.Context().Value(userContextKey).(string)
	if !ok {
		return
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"isGroupChat": false,
			"$and": bson.A{
				bson.M{"users": bson.M{"$elemMatch": bson.M{"$eq": userID}}},
				bson.M{"users": bson.M{"$elemMatch": bson.M{"$eq": receiverID}}},
			},
		}}},

		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "users"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "users_data"},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "messages"},
			{Key: "localField", Value: "latest_message_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "latest_message_data"},
		}}},

		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$latest_message_data"},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "let", Value: bson.D{{Key: "sender_id", Value: "$latest_message_data.sender_id"}}},
			{Key: "pipeline", Value: mongo.Pipeline{
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$_id", "$$sender_id"}},
					}},
				}}},
			}},
			{Key: "as", Value: "latest_message_data.sender_data"},
		}}},

		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$latest_message_data.sender_data"},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},
	}

	chat, err := c.service.Find(context.Background(), pipeline)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	if chat == nil {
		chat = &models.Chat{
			Name:        "sender",
			IsGroupChat: false,
			Users:       []string{userID, receiverID},
		}

		chat, err = c.service.Create(context.Background(), chat)
		if err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	json.NewEncoder(w).Encode(map[string]any{
		"Chat": chat,
	})
}

func (c *ChatController) FetchChats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	const userContextKey types.ContextJWTClaimKey = "userID"

	userID, ok := r.Context().Value(userContextKey).(string)
	if !ok {
		return
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "users", Value: userID},
		}}},

		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "users"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "users_data"},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "messages"},
			{Key: "localField", Value: "latest_message_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "latest_message_data"},
		}}},

		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$latest_message_data"},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},

		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "latest_message_data.sender_id"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "sender_data"},
		}}},

		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$sender_data"},
			{Key: "preserveNullAndEmptyArrays", Value: false},
		}}},

		bson.D{{Key: "$sort", Value: bson.D{{Key: "users.updated_at", Value: -1}}}},
	}

	chats, err := c.service.FindMany(context.Background(), pipeline)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"Chats": chats,
	})
}
