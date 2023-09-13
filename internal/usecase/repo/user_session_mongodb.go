package repo

import (
	"authentication/internal/entity"
	"authentication/internal/usecase"
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserSessionsRepo struct {
	db mongo.Client
}

func NewUsersSessionsRepo(db mongo.Client) UserSessionsRepo {
	return UserSessionsRepo{db}
}

var _ usecase.UserSessionRepo = (*UserSessionsRepo)(nil)

func (r UserSessionsRepo) Create(ctx context.Context, us entity.UserSession) (*entity.UserSession, error) {
	collection := r.db.Database("user_session").Collection("user_session")

	_, err := collection.InsertOne(ctx, us)
	if err != nil {
		return nil, fmt.Errorf("MongoDBUserSessionStore - Create - Insert: %w", err)
	}

	return &us, nil
}

func (r UserSessionsRepo) Get(ctx context.Context, userId uuid.UUID) (*entity.UserSession, error) {
	collection := r.db.Database("user_session").Collection("user_session")

	filter := bson.D{{"userid", userId}}

	userSession := entity.UserSession{}
	err := collection.FindOne(ctx, filter).Decode(&userSession)
	if err != nil {
		return nil, fmt.Errorf("MongoDBUserSessionStore - Get - FindOne: %w", err)
	}

	return &userSession, nil
}

func (r UserSessionsRepo) Update(ctx context.Context, sessionId uuid.UUID, newUserSession entity.UserSession) (*entity.UserSession, error) {
	collection := r.db.Database("user_session").Collection("user_session")

	filter := bson.D{{"id", sessionId}, {"userid", newUserSession.UserId}}

	update := bson.D{
		{"$set", newUserSession},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("MongoDBUserSessionStore - Update - UpdateOne: %w", err)
	}

	return &newUserSession, nil
}

func (r UserSessionsRepo) Delete(ctx context.Context, sessionId uuid.UUID) (*entity.UserSession, error) {
	collection := r.db.Database("user_session").Collection("user_session")

	filter := bson.D{{"id", sessionId}}

	userSessios := entity.UserSession{}
	err := collection.FindOne(ctx, filter).Decode(&userSessios)
	if err != nil {
		return nil, fmt.Errorf("MongoDBUserSessionStore - Delete - FindOne: %w", err)
	}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("MongoDBUserSessionStore - Delete - DeleteOne: %w", err)
	}

	return &userSessios, nil
}
