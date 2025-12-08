package repositories

import (
	"context"
	"go_chat/src/dto"
	"go_chat/src/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepository struct {
	collection *mongo.Collection
}

func NewChatRepository(db *mongo.Database) *ChatRepository {
	return &ChatRepository{
		collection: db.Collection("chats"),
	}
}

func (r *ChatRepository) Create(ctx context.Context, chat *models.Chat) (*dto.ChatDTO, error) {
	result, err := r.collection.InsertOne(ctx, chat)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": result.InsertedID.(string)}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "users"},
			{Key: "localField", Value: "users"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "users"},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if !cursor.Next(ctx) {
		return nil, nil
	}

	var newChat dto.ChatDTO

	err = cursor.Decode(&newChat)
	if err != nil {
		return nil, err
	}

	return &newChat, nil
}

func (r *ChatRepository) FindByID(ctx context.Context, chatID string) (*models.Chat, error) {
	var chat models.Chat

	err := r.collection.FindOne(ctx, bson.M{"_id": chatID}).Decode(&chat)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &chat, nil
}

func (r *ChatRepository) Find(ctx context.Context, pipeline mongo.Pipeline) (*dto.ChatDTO, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}

	if !cursor.Next(ctx) {
		return nil, nil
	}

	var chat dto.ChatDTO

	err = cursor.Decode(&chat)
	if err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *ChatRepository) FindMany(ctx context.Context, pipeline mongo.Pipeline) (*[]dto.ChatDTO, error) {
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	if !cursor.Next(ctx) {
		return nil, nil
	}

	var chats []dto.ChatDTO

	if err = cursor.All(ctx, &chats); err != nil {
		return nil, err
	}

	return &chats, nil
}

func (r *ChatRepository) GetUsersID(ctx context.Context, chatID string) ([]string, error) {
	var chat models.Chat

	projection := bson.D{{Key: "users", Value: 1}}
	opts := options.FindOne().SetProjection(projection)

	err := r.collection.FindOne(ctx, bson.M{"_id": chatID}, opts).Decode(&chat)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return chat.Users, nil
}
