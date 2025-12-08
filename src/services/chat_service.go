package services

import (
	"context"
	"go_chat/src/dto"
	"go_chat/src/models"
	"go_chat/src/repositories"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatService struct {
	repo *repositories.ChatRepository
}

func NewChatService(repo *repositories.ChatRepository) *ChatService {
    return &ChatService{
        repo: repo,
    }
}

func (s *ChatService) Create(ctx context.Context, chat *models.Chat) (*models.Chat, error) {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	
	chat.ID = id
	chat.CreatedAt = time.Now()
	
	return s.repo.Create(ctx, chat)
}

func (s *ChatService) Find(ctx context.Context, pipeline mongo.Pipeline) (*models.Chat, error) {
	return s.repo.Find(ctx, pipeline)
}

func (s *ChatService) FindMany(ctx context.Context, pipeline mongo.Pipeline) (*[]dto.ChatDTO, error) {
	return s.repo.FindMany(ctx, pipeline)
}

func (s *ChatService) GetUsersID(ctx context.Context, chatID string) ([]string, error) {
	return s.repo.GetUsersID(ctx, chatID)
}
