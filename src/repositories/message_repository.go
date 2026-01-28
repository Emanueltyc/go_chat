package repositories

import (
	"context"
	"go_chat/src/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) *MessageRepository {
	return &MessageRepository{
		collection: db.Collection("messages"),
	}
}

func (r *MessageRepository) Create(ctx context.Context, message *models.Message) (*models.Message, error) {
	result, err := r.collection.InsertOne(ctx, message)
	if err != nil {
		log.Print("There was an error trying to create the message (repository): ", err)
		return nil, err
	}

	var newMessage models.Message

	filter := bson.M{"_id": result.InsertedID}

	err = r.collection.FindOne(context.TODO(), filter).Decode(&newMessage)
	if err != nil {
		log.Fatal("There was an error trying to find the message (repository): ", err)
		return nil, err
	}

	return &newMessage, nil
}

func (r *MessageRepository) GetMessages(ctx context.Context, chatID string, limit int64, offset int64) (*[]models.Message, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"chat_id": chatID}, options.Find().SetLimit(limit).SetSkip(offset * limit))

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	defer cursor.Close(ctx)

	var messages *[]models.Message
	if err = cursor.All(ctx, &messages); err != nil { /* handle error */
	}

	return messages, nil
}