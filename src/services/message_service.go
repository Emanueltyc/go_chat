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
