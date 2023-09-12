package usecase

import (
	"authentication/internal/entity"
	"context"
	"github.com/google/uuid"
)

type (
	Users interface {
		Create(ctx context.Context, u entity.User) (*entity.User, error)
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
	}

	UsersRepo interface {
		Create(ctx context.Context, u entity.User) (*entity.User, error)
		GetByEmail(ctx context.Context, email string) (*entity.User, error)
		Update(ctx context.Context, u entity.User, email string) (*entity.User, error)
	}

	UserSessions interface {
		Add(ctx context.Context, userId uuid.UUID, refreshToken string) (*entity.UserSession, error)
		Refresh(ctx context.Context, sessionId, userId uuid.UUID, refreshToken string) (*entity.UserSession, error)
		Get(ctx context.Context, userId uuid.UUID) (*entity.UserSession, error)
	}

	UserSessionRepo interface {
		Create(ctx context.Context, us entity.UserSession) (*entity.UserSession, error)
		Get(ctx context.Context, userId uuid.UUID) (*entity.UserSession, error)
		Update(ctx context.Context, sessionId uuid.UUID, newUserSession entity.UserSession) (*entity.UserSession, error)
		Delete(ctx context.Context, sessionId uuid.UUID) (*entity.UserSession, error)
	}
)
