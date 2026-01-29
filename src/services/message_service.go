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
)

type MessageService struct {
	repo *repositories.MessageRepository
	chatRepo *repositories.ChatRepository
}

func NewMessageService(repo *repositories.MessageRepository, chatRepo *repositories.ChatRepository) *MessageService {
    return &MessageService{
        repo: repo,
		chatRepo: chatRepo,
    }
}

func (s *MessageService) Create(ctx context.Context, message *models.Message) (*models.Message, error) {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	
	message.ID = id
	message.CreatedAt = time.Now()
	
	return s.repo.Create(ctx, message)
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
