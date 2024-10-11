package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/exp/slog"
	//_ "icu/cmd/app/docs"
	"bem/internal/config"
	"bem/internal/domain/users/handler"
	"bem/internal/domain/users/repository"
	"bem/internal/domain/users/service"
	"bem/internal/middleware"
	"bem/pkg/jwt_auth"
	"bem/pkg/lib/logger/sl"
	"log"
	"net/http"
)

// @title TODO App API
// @version 0.1
// @description BEM server API
// @host localhost:8080
// @BasePath /
func main() {
	// Подключение к базе данных
	cfg := config.MustLoad()
	logger := createLogger("config.yml")
	r := gin.Default()

	db := createDatabase(cfg, logger)
	defer db.Close()

	// Создание маршрутизации
	err := createUserHandlers(db, cfg, logger, r)
	c := createCors()
	httpHandler := c.Handler(r)

	// Запуск сервера
	log.Println("Starting server on :8080")
	if err = http.ListenAndServe(":8080", httpHandler); err != nil {
		log.Fatal(err)
	}
}

func createDocsHandler(db *sql.DB, cfg *config.Config, logger *slog.Logger, r *gin.Engine) {

}

func createUserHandlers(db *sql.DB, cfg *config.Config, logger *slog.Logger, r *gin.Engine) error {
	accessSecret, err := jwt_auth.GenerateRandomSecret()
	refreshSecret, err := jwt_auth.GenerateRandomSecret()
	userRepo := repository.NewUserRepository(db, logger)
	userService := service.NewUserService(userRepo, cfg, logger, accessSecret, refreshSecret)

	userHandler := handler.NewUserHandler(userService, logger)
	jwtMiddleware := middleware.NewJWTMiddleware(userService, "your_jwt_secret_key")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/api/auth/sign_up", userHandler.CreateUser)
	r.POST("/api/auth/sign_in", userHandler.Login)
	r.POST("/api/auth/refresh", userHandler.RefreshToken)
	r.GET("/profile", jwtMiddleware.Authenticate(), userHandler.GetProfile)
	r.GET("/users/all", jwtMiddleware.Authenticate(), userHandler.GetAllUsers)

	return err
}

func createDatabase(cfg *config.Config, logger *slog.Logger) *sql.DB {
	db, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			cfg.Storage.Username, cfg.Storage.Password, cfg.Storage.Database))
	if err != nil {
		logger.Error(fmt.Sprintf("could not connect to the database: %v", err))
		log.Fatalf("could not connect to the database: %v", err)
	}

	return db
}

func createCors() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Разрешенные источники
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
	})
}

func createLogger(path string) *slog.Logger {
	logger, err := sl.NewLogger(path)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	logger.Info("Logger is created!")

	return logger
}
