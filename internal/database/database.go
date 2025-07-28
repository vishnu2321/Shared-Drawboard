package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/shared-drawboard/internal/models"
	"github.com/shared-drawboard/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name,omitempty"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
}

type DB interface {
	SaveUserDB(ctx context.Context, u models.User) (id string, err error)
	FindBy(ctx context.Context, field string, value interface{}) (*User, error)
}

type MongoDB struct {
	client *mongo.Client
	db     *mongo.Database
}

const (
	USER_COLLECTION = "users"
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

	_, err = col.InsertOne(ctx, u)
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
		fmt.Println(err)
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
