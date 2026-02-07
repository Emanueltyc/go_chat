package dto

import (
	"time"
)

type ChatDTO struct {
	ID            string     `bson:"_id,omitempty" json:"id"`
	Name          string     `bson:"name" json:"name"`
	IsGroupChat   bool       `bson:"isGroupChat" json:"is_group_chat"`
	GroupAmin     string     `bson:"group_admin_id" json:"group_admin,omitempty"`
	CreatedAt     time.Time  `bson:"created_at" json:"created_at"`
	Users         []UserDTO  `bson:"users,omitempty" json:"users"`
	LatestMessage MessageDTO `bson:"latest_message,omitempty" json:"latest_message"`
}

type MessageDTO struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	SenderID  string    `bson:"sender_id,omitempty" json:"-"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	Sender    UserDTO   `bson:"sender,omitempty" json:"sender"`
}

type UserDTO struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email" validate:"required,email"`
	Picture   string    `bson:"picture" json:"picture"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
