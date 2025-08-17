package models

// currently using as DTO but its better to keep models and DTOs seperate
type User struct {
	ID       string `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type SessionDTO struct {
	ID         string `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID     string `json:"user_id" bson:"user_id"`
	TokenHash  string `json:"token_hash" bson:"token_hash"`
	ExpiresAt  string `json:"expires_at" bson:"expires_at"`
	CreatedAt  string `json:"created_at" bson:"created_at"`
	LastUsedAt string `json:"last_used_at" bson:"last_used_at"`
}

type EventType string

const (
	FreehandDraw EventType = "freehandDraw"
	ShapeCreate  EventType = "shapeCreate"
	TextAdd      EventType = "textAdd"
	ObjectDelete EventType = "objectDelete"
	BoardClear   EventType = "boardClear"
)

type Event struct {
	Type      EventType   `json:"type" bson:"type"`
	Tool      string      `json:"tool" bson:"tool"`
	CreatedAt string      `json:"timestamp,omitempty" bson:"created_at,omitempty"`
	Data      interface{} `json:"data" bson:"data"`
}

// Event Data Structures
type FreehandDrawData struct {
	Color     string  `json:"color" bson:"color"`
	Thickness float64 `json:"thickness" bson:"thickness"`
	Points    []Point `json:"points" bson:"points"`
}

type ShapeCreateData struct {
	Color     string  `json:"color" bson:"color"`
	Thickness float64 `json:"thickness" bson:"thickness"`
	X         float64 `json:"x" bson:"x"`
	Y         float64 `json:"y" bson:"y"`
	Width     float64 `json:"width" bson:"width"`
	Height    float64 `json:"height" bson:"height"`
}

type TextAddData struct {
	Color     string  `json:"color" bson:"color"`
	Thickness float64 `json:"thickness" bson:"thickness"`
	X         float64 `json:"x" bson:"x"`
	Y         float64 `json:"y" bson:"y"`
	Text      string  `json:"text" bson:"text"`
}

type ObjectDeleteData struct {
	Index int `json:"index" bson:"index"`
}

type Point struct {
	X float64 `json:"x" bson:"x"`
	Y float64 `json:"y" bson:"y"`
}

type RefreshTokenDTO struct {
	UserID           string `json:"user_id,omitempty" bson:"user_id,omitempty"`
	AuthToken        string `json:"auth-token,omitempty" bson:"auth-token,omitempty"`
	AuthExpiresAt    string `json:"auth-token-expiry,omitempty" bson:"auth-token-expiry,omitempty"`
	RefreshToken     string `json:"refresh-token,omitempty" bson:"refresh-token,omitempty"`
	RefreshExpriesAt string `json:"refresh-token-expiry,omitempty" bson:"refresh-token-expiry,omitempty"`
}
