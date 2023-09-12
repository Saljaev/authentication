package entity

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             uuid.UUID
	Email          string
	HashedPassword string
}

func NewUser(email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Cannot hash password: %w", err)
	}

	user := &User{
		Id:             uuid.New(),
		Email:          email,
		HashedPassword: string(hashedPassword),
	}

	return user, nil
}

func (u *User) IsPasswordCorrect(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}
