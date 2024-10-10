package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
	_ "github.com/swaggo/swag/example/celler/httputil"
	"golang.org/x/exp/slog"
	"icu/internal/domain/users/models"
	"icu/internal/domain/users/service"
	"net/http"
)

type AuthorizedUser struct {
	UserData     *models.User `json:"user_data"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
}

type credentials struct {
	Email string `json:"email"`
	//Username string `json:"username"`
	Password string `json:"password"`
}

type UserHandler struct {
	userService *service.UserService
	logger      *slog.Logger
}

func NewUserHandler(userService *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the input payload
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user      body    models.User     true        "User Data"
// @Failure      400  {object}  error
// @Router /api/auth/sign_up [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userService.CreateUser(c.Request.Context(), &user)
	if err != nil {
		if err.Error() == "user already exists" {
			c.JSON(http.StatusConflict, gin.H{"User already exists": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Could not create user": err.Error()})
		}
	}

	c.Status(http.StatusCreated)
}

// Login godoc
// @Summary Login a user
// @Description Login with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   credentials      body    credentials     true        "User Credentials"
// @Success 200 {object} AuthorizedUser
// @Failure 401 {object} error
// @Router /api/auth/sign_in [post]
func (h *UserHandler) Login(c *gin.Context) {
	creds := credentials{}

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid input": err.Error()})
		return
	}

	user, accessToken, refreshToken, err :=
		h.userService.Authenticate(c.Request.Context(), creds.Email, creds.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Invalid email or password": err.Error()})
		return
	}

	c.Header("Authorization", "Bearer "+accessToken)
	authUser := AuthorizedUser{UserData: user, AccessToken: accessToken, RefreshToken: refreshToken}
	h.logger.Info(fmt.Sprintf("authuser by access token: %s", authUser))

	c.JSON(http.StatusOK, &authUser)
}

func (h *UserHandler) RefreshToken(c *gin.Context) {
	var authUser AuthorizedUser
	//var refreshToken string
	if err := c.ShouldBindJSON(&authUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Invalid token format": err.Error()})
		return
	}

	//refreshToken := authUser.RefreshToken
	newAccessToken, err := h.userService.RefreshToken(authUser.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Can't create new access token": err.Error()})
		return
	}

	c.Header("Authorization", "Bearer "+newAccessToken)
	authUser.AccessToken = newAccessToken
	fmt.Printf("отправляемые refresh-data: %s\n", authUser)
	c.JSON(http.StatusOK, &authUser)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	// Получаем пользователя из контекста
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Преобразуем пользователя к нужному типу
	userModel, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting user data"})
		return
	}

	c.JSON(http.StatusOK, userModel)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Could not retrieve users": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, users)
}
