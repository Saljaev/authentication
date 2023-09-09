package entity

import (
	"github.com/google/uuid"
	"time"
)

type UserSession struct {
	Id        uuid.UUID
	UserdId   uint32
	ExpiresAt time.Time
	UpdateAt  time.Time
	CreatedAt time.Time
}
