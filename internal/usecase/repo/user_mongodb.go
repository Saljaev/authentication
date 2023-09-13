package repo

import (
	"authentication/internal/entity"
	"authentication/internal/usecase"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRepo struct {
	db mongo.Client
}

func NewUsersRepo(db mongo.Client) UsersRepo {
	return UsersRepo{db}
}

var _ usecase.UsersRepo = (*UsersRepo)(nil)

func (r UsersRepo) Create(ctx context.Context, user entity.User) (*entity.User, error) {
	collection := r.db.Database("user").Collection("user")

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("MongoDBUserStore - Connect - Insert: %w", err)
	}

	return &user, nil
}

func (r UsersRepo) Update(ctx context.Context, u entity.User, email string) (*entity.User, error) {
	collection := r.db.Database("user").Collection("user")

	filter := bson.D{{"email", email}}

	update := bson.D{
		{"$set", u},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("MongoDBUserStore - Update - UpdateOne: %w", err)
	}

	return &u, nil
}

func (r UsersRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	collection := r.db.Database("user").Collection("user")

	filter := bson.D{{"email", email}}

	user := entity.User{}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("MongoDBUserStore - GetByEmail - FindOne: %w", err)
	}

	return &user, nil
}
