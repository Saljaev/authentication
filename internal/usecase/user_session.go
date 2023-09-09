package usecase

import (
	"authentication/internal/entity"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type UserSessionUseCase struct {
	repo               UserSessionRepo
	refreshTokenLength uint32
	sessionDuration    time.Duration
}

func NewUserSessionsUseCase(repo UserSessionRepo, refreshTokenLength uint32, sessionDuration time.Duration) *UserSessionUseCase {
	return &UserSessionUseCase{repo, refreshTokenLength, sessionDuration}
}

var _ UserSessions = (*UserSessionUseCase)(nil)

func (uc UserSessionUseCase) Add(ctx context.Context, userId uint32) (*entity.UserSession, error) {
	now := time.Now()

	e, err := uc.repo.Create(ctx, entity.UserSession{
		UserdId:   userId,
		ExpiresAt: now.Add(uc.sessionDuration),
		UpdateAt:  now,
		CreatedAt: now,
	})
	if err != nil {
		return nil, fmt.Errorf("UserSessionCase - Add: %w", err)
	}

	return e, nil
}

func (uc UserSessionUseCase) Refresh(ctx context.Context, sessionId uuid.UUID, userId uint32) (*entity.UserSession, error) {
	now := time.Now()

	e, err := uc.repo.Update(ctx, sessionId, entity.UserSession{
		UserdId:   userId,
		ExpiresAt: now.Add(uc.sessionDuration),
		UpdateAt:  now,
	})
	if err != nil {
		return nil, fmt.Errorf("UserSessionCase - Refresh: %w", err)
	}

	return e, nil
}
