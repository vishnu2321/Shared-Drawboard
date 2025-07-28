package models

// currently using as DTO but its better to keep models and DTOs seperate
type User struct {
	ID       string `json:"_id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
