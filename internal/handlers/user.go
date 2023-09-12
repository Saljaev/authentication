package handlers

import (
	"authentication/internal/entity"
	"authentication/internal/usecase"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type UserHandler struct {
	users   usecase.Users
	session usecase.UserSessions
	jwt     JWTManager
}

func NewUserHandler(u usecase.Users, us usecase.UserSessions, jwt JWTManager) *UserHandler {
	return &UserHandler{u, us, jwt}
}

func (h *UserHandler) Register(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	user, err := entity.NewUser(r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		return err
	}

	storeUser, err := h.users.Create(ctx, *user)
	if err != nil {
		return err
	}

	jsonResp, err := json.Marshal(storeUser)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

	return nil
}

func (h *UserHandler) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	user, err := h.users.GetByEmail(ctx, r.FormValue("email"))
	if err != nil {
		return err
	}

	if !user.IsPasswordCorrect(r.FormValue("password")) {
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

	// TODO: set session time limit 7 day
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

	return nil
}

func (h *UserHandler) RefreshToken(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	user, err := h.users.GetByEmail(ctx, r.FormValue("email"))
	if err != nil {
		return nil
	}

	refreshFromCookie, err := r.Cookie("RefreshToken")
	if err != nil {
		return fmt.Errorf("No refresh token")
	}

	session, err := h.session.Get(ctx, user.Id)
	if err != nil {
		return nil
	}

	if session.RefreshToken != refreshFromCookie.Value {
		return fmt.Errorf("Can't compare refresh token with token in session")
	}

	refreshToken, err := h.jwt.NewRefreshToken()
	if err != nil {
		return err
	}

	// TODO: change session repo
	_, err = h.session.Refresh(ctx, uuid.MustParse(r.FormValue("refresh_token")), user.Id, refreshToken)
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
