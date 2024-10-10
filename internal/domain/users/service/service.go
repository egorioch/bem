package service

import (
	"bem/internal/config"
	"bem/internal/domain/users/models"
	"bem/internal/domain/users/repository"
	"bem/pkg/jwt_auth"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/exp/slog"
	"time"
)

const location = "users.service"

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type UserService struct {
	userRepo      *repository.UserRepository
	logger        *slog.Logger
	accessSecret  string
	refreshSecret string
	config        *config.Config
}

func NewUserService(userRepo *repository.UserRepository, config *config.Config, logger *slog.Logger, accessSecret, refreshSecret string) *UserService {
	return &UserService{
		userRepo:      userRepo,
		config:        config,
		logger:        logger,
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (us *UserService) CreateUser(ctx context.Context, user *models.User) error {
	existingUser, err := us.userRepo.UserExists(ctx, user.Email)

	if err != nil {
		return fmt.Errorf("error checking for existing user: %v", err)
	}
	if existingUser == 1 {
		return errors.New("user already exists")
	}
	if user.Username == "" {
		return errors.New("username isn't defined")
	}
	if len(user.Username) < 8 {
		return errors.New("must be at least 8 characters")
	}
	if user.Password == "" {
		return errors.New("password isn't defined")
	}
	if len(user.Password) < 8 {
		return errors.New("must be at least 8 characters")
	}
	if user.Email == "" {
		return errors.New("email isn't defined")
	}

	fmt.Printf("new user: %+v", user)
	if user.AdminToken == us.config.AccessPermission.AdminToken {
		user.Role = "admin"
	} else {
		user.Role = "default"
	}

	err = us.userRepo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("create user error: %s, %s", err, location)
	}

	return nil
}

func (us *UserService) Authenticate(ctx context.Context, email, password string) (user *models.User, accessToken string, refreshToken string, err error) {
	user, err = us.userRepo.FindOne(ctx, email)
	if err != nil || user == nil {
		return nil, "", "", fmt.Errorf("user not found")
	}
	if user.Password != password {
		return nil, "", "", fmt.Errorf("invalid password")
	}

	accessToken, err = jwt_auth.GenerateToken(user.Email, us.accessSecret, time.Hour*1)
	refreshToken, err = jwt_auth.GenerateToken(user.Email, us.refreshSecret, time.Hour*24*3)
	if err != nil {
		us.logger.Error(fmt.Sprintf("generate token error: %s", err))
		fmt.Printf("generate token error: %s", err)
		return nil, "", "", fmt.Errorf("could not create token")
	}

	return user, accessToken, refreshToken, nil
}

func (s *UserService) RefreshToken(refreshToken string) (string, error) {
	claims := &Claims{}
	fmt.Println("input token: ", refreshToken)
	tkn, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.refreshSecret), nil
	})
	if err != nil || !tkn.Valid {
		return err.Error(), errors.New("invalid token")
	}

	newToken, err := jwt_auth.GenerateToken(claims.Email, s.accessSecret, time.Hour*24*3)
	if err != nil {
		fmt.Printf("error: %s", err)
		s.logger.Error(err.Error())
		return err.Error(), err
	}

	fmt.Println("new token: " + newToken)

	return newToken, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.userRepo.FindOne(ctx, email)
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.userRepo.FindAll(ctx)
}
