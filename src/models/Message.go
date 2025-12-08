package models

import (
	"time"
)

type Message struct {
	ID         string    `bson:"_id,omitempty" json:"id"`
	Sender     string    `bson:"sender_id" json:"sender_id" validate:"required"`
	Content    string    `bson:"content" json:"content" validate:"required"`
	Chat       string    `bson:"chat_id" json:"chat_id" validate:"required"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}

func NewMessage(sender string, Content string, chat string) *Message {
	return &Message{
		Sender:    sender,
		Content:   Content,
		Chat:      chat,
		CreatedAt: time.Now(),
	}
}
