package services

import (
	"context"
	"go_chat/src/models"
	"go_chat/src/repositories"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
    return &UserService{
        repo: repo,
    }
}

func (s *UserService) Register(ctx context.Context, user *models.User) (*models.User, error) {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	user.ID = id
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	return s.repo.Create(ctx, user)
}

func (s *UserService) Find(ctx context.Context, filter bson.M) (*models.User, error) {
	return s.repo.Find(ctx, filter)
}

func (s *UserService) MatchPassword(user *models.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (s *UserService) GenerateToken(user *models.User) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"picture": user.Picture,
		"createdAt": user.CreatedAt,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", nil
	}

	return tokenString, nil
}

func (s *UserService) SearchUsers(ctx context.Context, filter bson.M) (*[]models.User, error) {
	return s.repo.Search(ctx, filter)
}
