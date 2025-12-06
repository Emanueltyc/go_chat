package models

import (
	"time"
)

type Chat struct {
	ID                string    `bson:"_id,omitempty" json:"id"`
	Name              string    `bson:"name" json:"name" validate:"required"`
	IsGroupChat       bool      `bson:"isGroupChat" json:"is_group_chat" validate:"required"`
	Users             []string  `bson:"users" json:"-" validate:"required"`
	LatestMessage     string    `bson:"latest_message_id" json:"-"`
	GroupAmin         string    `bson:"group_admin_id" json:"group_admin,omitempty"`
	CreatedAt         time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at" json:"updated_at"`
	
	UsersData         []User    `bson:"users_data,omitempty" json:"users"`
	LatestMessageData Message   `bson:"latest_message_data,omitempty" json:"latest_message"`
}
