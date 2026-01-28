package services

import (
	"context"
	"go_chat/src/models"
	"go_chat/src/repositories"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MessageService struct {
	repo *repositories.MessageRepository
}

func NewMessageService(repo *repositories.MessageRepository) *MessageService {
    return &MessageService{
        repo: repo,
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
	
	return s.repo.GetMessages(ctx, chatId, limit, offset)
}
