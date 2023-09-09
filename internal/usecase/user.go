package usecase

import (
	"authentication/internal/entity"
	"context"
	"fmt"
)

type UsersUseCase struct {
	repo UsersRepo
}

func NewUsersUseCase(repo UsersRepo) *UsersUseCase {
	return &UsersUseCase{repo}
}

var _ Users = (*UsersUseCase)(nil)

func (uc UsersUseCase) Create(ctx context.Context, u entity.User) (*entity.User, error) {
	e, err := uc.repo.Create(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("UsersUseCase - Create: %w", err)
	}

	return e, nil
}

func (uc UsersUseCase) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	e, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("UsersUseCase - GetByEmail: %w", err)
	}

	return e, nil
}
