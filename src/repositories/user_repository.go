package repositories

import (
	"context"
	"go_chat/src/models"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var newUser models.User

	filter := bson.M{"_id": result.InsertedID}

	err = r.collection.FindOne(context.TODO(), filter).Decode(&newUser)
	if err != nil {
		log.Fatal(err)
	}

	return &newUser, nil
}

func (r *UserRepository) Find(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Search(ctx context.Context, filter bson.M) (*[]models.User, error) {
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var users []models.User
	
	if err = cursor.All(ctx, &users); err != nil { /* handle error */
	}
	
	if len(users) == 0 {
		return nil, nil
	}

	return &users, nil
}
