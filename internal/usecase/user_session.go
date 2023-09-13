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
	refreshTokenLength int
	sessionDuration    time.Duration
}

func NewUserSessionsUseCase(repo UserSessionRepo, refreshTokenLength int, sessionDuration time.Duration) *UserSessionUseCase {
	return &UserSessionUseCase{repo, refreshTokenLength, sessionDuration}
}

var _ UserSessions = (*UserSessionUseCase)(nil)

func (uc UserSessionUseCase) Add(ctx context.Context, userId uuid.UUID, refreshToken string) (*entity.UserSession, error) {
	now := time.Now()

	e, err := uc.repo.Create(ctx, entity.UserSession{
		Id:           uuid.New(),
		UserId:       userId,
		ExpiresAt:    now.Add(uc.sessionDuration),
		UpdateAt:     now,
		CreatedAt:    now,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("UserSessionCase - Add: %w", err)
	}

	return e, nil
}

func (uc UserSessionUseCase) Refresh(ctx context.Context, sessionId, userId uuid.UUID, refreshToken string) (*entity.UserSession, error) {
	now := time.Now()

	e, err := uc.repo.Update(ctx, sessionId, entity.UserSession{
		UserId:       userId,
		ExpiresAt:    now.Add(uc.sessionDuration),
		UpdateAt:     now,
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, fmt.Errorf("UserSessionCase - Refresh: %w", err)
	}

	return e, nil
}
func (uc UserSessionUseCase) Get(ctx context.Context, userId uuid.UUID) (*entity.UserSession, error) {
	session, err := uc.repo.Get(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("UserSessionCase - Get: %w", err)
	}

	return session, nil
}
