package services

import (
	"context"
	"errors"
	"go_chat/src/models"
	"go_chat/src/repositories"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type MessageService struct {
	repo     *repositories.MessageRepository
	chatRepo *repositories.ChatRepository
}

func NewMessageService(repo *repositories.MessageRepository, chatRepo *repositories.ChatRepository) *MessageService {
	return &MessageService{
		repo:     repo,
		chatRepo: chatRepo,
	}
}

func (s *MessageService) Create(ctx context.Context, message *models.Message) error {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")

	message.ID = id
	message.CreatedAt = time.Now()

	message, err := s.repo.Create(ctx, message)
	if err != nil {
		return err
	}

	return s.chatRepo.UpdateByID(ctx, message.Chat, bson.D{{Key: "$set", Value: bson.D{{Key: "latest_message_id", Value: message.ID}}}})
}

func (s *MessageService) Fetch(ctx context.Context, userId string, chatId string, limit int64, offset int64) (*[]models.Message, error) {
	limit = min(limit, 100)

	chat, err := s.chatRepo.FindByID(ctx, chatId)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(chat.Users, userId) {
		return nil, errors.New("User does not have access to the chat!")
	}

	return s.repo.GetMessages(ctx, chatId, limit, offset)
}
