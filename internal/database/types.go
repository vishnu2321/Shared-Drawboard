package database

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type Session struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserID     string             `bson:"user_id" json:"user_id"`
	TokenHash  string             `bson:"token_hash" json:"token_hash"`
	ExpiresAt  string             `bson:"expires_at" json:"expires_at"`
	CreatedAt  string             `bson:"created_at" json:"created_at"`
	LastUsedAt string             `bson:"last_used_at" json:"last_used_at"`
}
