package models

import "time"

// currently using as DTO but its better to keep models and DTOs seperate
type User struct {
	ID       string `json:"_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SessionDTO struct {
	ID         string `json:"_id,omitempty"`
	UserID     string `json:"user_id"`
	TokenHash  string `json:"token_hash"`
	ExpiresAt  string `json:"expires_at"`
	CreatedAt  string `json:"created_at"`
	LastUsedAt string `json:"last_used_at"`
}

type DrawingEventDTO struct {
	EventID   string    `bson:"eventId" json:"eventId"`
	BoardID   string    `bson:"boardId" json:"boardId"`
	UserID    string    `bson:"userId" json:"userId"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	EventType string    `bson:"eventType" json:"eventType"`
	Tool      string    `bson:"tool" json:"tool"`

	// Common properties for all tools
	Color     string `bson:"color,omitempty" json:"color,omitempty"`
	Thickness int    `bson:"thickness,omitempty" json:"thickness,omitempty"`

	// Position and size properties
	X      float64 `bson:"x,omitempty" json:"x,omitempty"`
	Y      float64 `bson:"y,omitempty" json:"y,omitempty"`
	Width  float64 `bson:"width,omitempty" json:"width,omitempty"`
	Height float64 `bson:"height,omitempty" json:"height,omitempty"`

	// Points for freehand drawing and erasing
	Points []Point `bson:"points,omitempty" json:"points,omitempty"`

	// Text properties
	Text     string `bson:"text,omitempty" json:"text,omitempty"`
	FontSize int    `bson:"fontSize,omitempty" json:"fontSize,omitempty"`

	// Object reference for deletion
	ObjectID string `bson:"objectId,omitempty" json:"objectId,omitempty"`
}

type Point struct {
	X float64 `bson:"x" json:"x"`
	Y float64 `bson:"y" json:"y"`
}

type RefreshTokenDTO struct {
	AuthToken        string `json:"auth-token,omitempty"`
	AuthExpiresAt    string `json:"auth-token-expiry,omitempty"`
	RefreshToken     string `json:"refresh-token,omitempty"`
	RefreshExpriesAt string `json:"refresh-token-expiry,omitempty"`
}
