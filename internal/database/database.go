package database

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shared-drawboard/internal/models"
	"github.com/shared-drawboard/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB interface {
	SaveUserDB(ctx context.Context, u models.User) (id string, err error)
	FindBy(ctx context.Context, field string, value interface{}) (*User, error)
	CreateSession(ctx context.Context, session models.SessionDTO) (string, error)
	UpdateSession(ctx context.Context, uid string, newToken string) (string, error)
	BatchSave(ctx context.Context, batch []interface{}) error
}

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

const (
	USER_COLLECTION    = "users"
	SESSION_COLLECTION = "sessions"
	EVENTS_COLLECTION  = "events"
)

func New() (*MongoDB, error) {
	dbConfig := Config()
	if dbConfig["uri"] == "" || dbConfig["db_name"] == "" {
		return nil, fmt.Errorf("DB: Update the environment values")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOps := options.Client().ApplyURI(dbConfig["uri"].(string))
	client, err := mongo.Connect(ctx, clientOps)
	if err != nil {
		return nil, fmt.Errorf("DB: %w", err)
	}

	db := client.Database(dbConfig["db_name"].(string))
	logger.Info("Database connected successfully")

	return &MongoDB{client: client, db: db}, nil
}

func (m *MongoDB) SaveUserDB(ctx context.Context, u models.User) (id string, err error) {
	col := m.db.Collection(USER_COLLECTION)

	user := User{Name: u.Name, Email: u.Email, Password: u.Password}

	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	_, err = col.InsertOne(ctx, user)
	if err != nil {
		logger.Error("Insert failed: %v", err)
		return "", fmt.Errorf("failed to insert user: %w", err)
	}

	return user.ID.Hex(), nil
}

func (m *MongoDB) FindBy(ctx context.Context, field string, value interface{}) (*User, error) {
	col := m.db.Collection(USER_COLLECTION)

	var user User

	if err := col.FindOne(ctx, bson.M{strings.ToLower(field): value}).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func ClearPreviousSessions(ctx context.Context, col *mongo.Collection, userID string) error {
	filter := bson.M{"user_id": userID}

	_, err := col.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoDB) CreateSession(ctx context.Context, s models.SessionDTO) (string, error) {
	col := m.db.Collection(SESSION_COLLECTION)

	session := Session{
		UserID:     s.UserID,
		TokenHash:  s.TokenHash,
		ExpiresAt:  s.ExpiresAt,
		CreatedAt:  s.CreatedAt,
		LastUsedAt: s.LastUsedAt,
	}

	if session.ID.IsZero() {
		session.ID = primitive.NewObjectID()
	}

	err := ClearPreviousSessions(ctx, col, s.UserID)
	if err != nil {
		logger.Error("Insert failed: %v", err)
		return "", fmt.Errorf("failed to clear previous sessions: %w", err)
	}

	_, err = col.InsertOne(ctx, session)
	if err != nil {
		logger.Error("Insert failed: %v", err)
		return "", fmt.Errorf("failed to insert user: %w", err)
	}
	return session.ID.Hex(), nil
}

func (m *MongoDB) UpdateSession(ctx context.Context, uid string, newTokenHash string) (string, error) {
	col := m.db.Collection(SESSION_COLLECTION)

	filter := bson.M{"user_id": uid}
	update := bson.M{
		"$set": bson.M{
			"token_hash":   newTokenHash,
			"created_at":   strconv.FormatInt(time.Now().Unix(), 10),
			"last_used_at": strconv.FormatInt(time.Now().Unix(), 10),
			"expires_at":   strconv.FormatInt(time.Now().Add(24*7*time.Hour).Unix(), 10),
		},
	}

	doc, err := col.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Update failed: %v", err)
		return "", fmt.Errorf("failed to insert user: %w", err)
	}

	return doc.UpsertedID.(string), nil
}

// func (m *MongoDB) DeleteSession(ctx context.Context) {
// 	col := m.db.Collection(SESSION_COLLECTION)

// 	_, err := col.DeleteOne()

// }

func (m *MongoDB) BatchSave(ctx context.Context, batch []interface{}) error {
	col := m.db.Collection(EVENTS_COLLECTION)

	_, err := col.InsertMany(ctx, batch)
	if err != nil {
		logger.Error("Update failed: %v", err)
		return fmt.Errorf("failed to insert batch data: %w", err)
	}

	return nil
}
