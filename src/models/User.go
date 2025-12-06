package models

import (
	"time"
)

type User struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name" validate:"required"`
	Email     string    `bson:"email" json:"email" validate:"required,email"`
	Password  string    `bson:"password" json:"-" validate:"required"`
	Picture   string    `bson:"picture" json:"picture"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
