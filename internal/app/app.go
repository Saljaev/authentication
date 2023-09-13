package app

import (
	"authentication/internal/handlers"
	"authentication/internal/usecase"
	"authentication/internal/usecase/repo"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"strconv"
	"time"
)

func Run() {
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("Error to load .env file: %w", err)
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	userOpts := options.Client().ApplyURI(os.Getenv("user_db_url")).SetServerAPIOptions(serverAPI)
	userClient, err := mongo.Connect(context.TODO(), userOpts)
	if err != nil {
		fmt.Errorf("Error to connect to user database: %w", err)
	}
	defer userClient.Disconnect(context.TODO())

	sessionOpts := options.Client().ApplyURI(os.Getenv("session_db_url")).SetServerAPIOptions(serverAPI)
	sessionClient, err := mongo.Connect(context.TODO(), sessionOpts)
	if err != nil {
		fmt.Errorf("Error to connect to session database: %w", err)
	}
	defer sessionClient.Disconnect(context.TODO())

	err = userClient.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("User crush")
	}

	err = sessionClient.Ping(context.Background(), nil)
	if err != nil {
		fmt.Println("Session crush")
	}

	userRepo := repo.NewUsersRepo(*userClient)
	userSessionRepo := repo.NewUsersSessionsRepo(*sessionClient)

	sessionTime := (time.Hour * 24 * 7)

	jwtManager := handlers.NewJWTManager(os.Getenv("secret"), time.Hour*1)

	tokenLength, err := strconv.ParseInt(os.Getenv("token_length"), 10, 64)
	if err != nil {
		fmt.Errorf("Error to convert token length to int")
	}

	handler := handlers.NewUserHandler(usecase.NewUsersUseCase(userRepo), usecase.NewUserSessionsUseCase(userSessionRepo, int(tokenLength), sessionTime), *jwtManager)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only post", http.StatusMethodNotAllowed)
			return
		}

		err := handler.Register(context.Background(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only post", http.StatusMethodNotAllowed)
			return
		}

		err := handler.Login(context.Background(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/validate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only get", http.StatusMethodNotAllowed)
			return
		}

		err := handler.Validate(context.Background(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only get", http.StatusMethodNotAllowed)
			return
		}

		err := handler.RefreshToken(context.Background(), w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":8080", nil)
}
