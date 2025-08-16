package models

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

type EventType string

const (
	FreehandDraw EventType = "freehandDraw"
	ShapeCreate  EventType = "shapeCreate"
	TextAdd      EventType = "textAdd"
	ObjectDelete EventType = "objectDelete"
	BoardClear   EventType = "boardClear"
)

type Event struct {
	Type EventType   `json:"type"`
	Tool string      `json:"tool"`
	Data interface{} `json:"data"`
}

// Event Data Structures
type FreehandDrawData struct {
	Color     string  `json:"color"`
	Thickness float64 `json:"thickness"`
	Points    []Point `json:"points"`
}

type ShapeCreateData struct {
	Color     string  `json:"color"`
	Thickness float64 `json:"thickness"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Width     float64 `json:"width"`
	Height    float64 `json:"height"`
}

type TextAddData struct {
	Color     string  `json:"color"`
	Thickness float64 `json:"thickness"`
	X         float64 `json:"x"`
	Y         float64 `json:"y"`
	Text      string  `json:"text"`
}

type ObjectDeleteData struct {
	Index int `json:"index"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type RefreshTokenDTO struct {
	UserID           string `json:"user_id,omitempty"`
	AuthToken        string `json:"auth-token,omitempty"`
	AuthExpiresAt    string `json:"auth-token-expiry,omitempty"`
	RefreshToken     string `json:"refresh-token,omitempty"`
	RefreshExpriesAt string `json:"refresh-token-expiry,omitempty"`
}
