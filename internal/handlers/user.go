package handlers

import (
	"authentication/internal/entity"
	"authentication/internal/usecase"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type UserHandler struct {
	users   usecase.Users
	session usecase.UserSessions
	jwt     JWTManager
}

type userJSON struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewUserHandler(u usecase.Users, us usecase.UserSessions, jwt JWTManager) *UserHandler {
	return &UserHandler{u, us, jwt}
}

func (h *UserHandler) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)

	var u userJSON

	err := decoder.Decode(&u)
	if err != nil {
		return err
	}

	user, err := entity.NewUser(u.Email, u.Password)
	if err != nil {
		return err
	}

	storeUser, err := h.users.Create(ctx, *user)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(storeUser.Email))

	return nil
}

func (h *UserHandler) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)

	var u userJSON

	err := decoder.Decode(&u)
	if err != nil {
		fmt.Errorf("Error to decode JSON: %w", err)
	}

	user, err := h.users.GetByEmail(ctx, u.Email)
	if err != nil {
		return err
	}

	if !user.IsPasswordCorrect(u.Password) {
		return err
	}

	accessToken, err := h.jwt.Generate(user)
	if err != nil {
		return err
	}

	refreshToken, err := h.jwt.NewRefreshToken()
	if err != nil {
		return err
	}

	session, err := h.session.Add(ctx, user.Id, refreshToken)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "RefreshToken",
		Value:   refreshToken,
		Expires: session.ExpiresAt,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfull login"))

	return nil
}

func (h *UserHandler) RefreshToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	decoder := json.NewDecoder(r.Body)

	var u userJSON

	err := decoder.Decode(&u)
	if err != nil {
		return fmt.Errorf("Error to decode JSON: %w", err)
	}

	user, err := h.users.GetByEmail(ctx, u.Email)
	if err != nil {
		return fmt.Errorf("Error to get user: %w", err)
	}

	refreshFromCookie, err := r.Cookie("RefreshToken")
	if err != nil {
		return fmt.Errorf("No refresh token")
	}

	session, err := h.session.Get(ctx, user.Id)
	if err != nil {
		fmt.Errorf("Error get user session: %w", err)
	}

	if session.RefreshToken != refreshFromCookie.Value {
		return fmt.Errorf("Can't compare refresh token with token in session: %w", http.StatusUnprocessableEntity)
	}

	refreshToken, err := h.jwt.NewRefreshToken()
	if err != nil {
		return err
	}

	_, err = h.session.Refresh(ctx, session.Id, user.Id, refreshToken)
	if err != nil {
		return err
	}

	accessToken, err := h.jwt.Generate(user)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Token",
		Value:    accessToken,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "RefreshToken",
		Value:   refreshToken,
		Expires: session.ExpiresAt,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Coockies refresh"))
	return nil
}

func (h *UserHandler) Validate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	token, err := r.Cookie("Token")
	if err != nil {
		return fmt.Errorf("No access token")
	}

	claims, err := h.jwt.Parse(token.Value)
	if err != nil {
		return err
	}

	if claims == nil || claims.Email == "" {
		return err
	}

	user, err := h.users.GetByEmail(ctx, claims.Email)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(user.Email))

	return nil
}
