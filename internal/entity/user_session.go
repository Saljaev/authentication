package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserSession struct {
	Id           uuid.UUID
	UserId       uuid.UUID
	RefreshToken string
	ExpiresAt    time.Time
	UpdateAt     time.Time
	CreatedAt    time.Time
}
