package auth

import (
	"net/http"
	"time"

	"github.com/joojf/travel-planner-api/internal/notification"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo                Repository
	tokenBlacklist      map[string]time.Time
	notificationService *notification.Service
}

func NewHandler(repo Repository, notificationService *notification.Service) *Handler {
	return &Handler{
		repo:                repo,
		tokenBlacklist:      make(map[string]time.Time),
		notificationService: notificationService,
	}
}

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (h *Handler) Register(c echo.Context) error {
	var user User
	if err := c.Bind(&user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := c.Validate(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := h.repo.CreateUser(&user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":    user.ID,
		"email": user.Email,
	})
}

func (h *Handler) Login(c echo.Context) error {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&loginRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.repo.GetUserByEmail(loginRequest.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	match, err := VerifyPassword(loginRequest.Password, user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to verify password")
	}

	if !match {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate token")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (h *Handler) Logout(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "No token provided")
	}

	h.tokenBlacklist[token] = time.Now().Add(24 * time.Hour)

	return c.NoContent(http.StatusOK)
}

func (h *Handler) ResetPassword(c echo.Context) error {
	var resetRequest struct {
		Email string `json:"email"`
	}

	if err := c.Bind(&resetRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := h.repo.GetUserByEmail(resetRequest.Email)
	if err != nil {
		return c.NoContent(http.StatusOK)
	}

	resetToken, err := GenerateResetToken(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate reset token")
	}

	resetLink := "http://localhost:5000/reset-password?token=" + resetToken
	message := "Click the following link to reset your password: " + resetLink
	err = h.notificationService.SendNotification(user.Email, "Password Reset", message)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to send reset email")
	}

	return c.NoContent(http.StatusOK)
}

func (h *Handler) SetNewPassword(c echo.Context) error {
	var setPasswordRequest struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.Bind(&setPasswordRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, err := ValidateResetToken(setPasswordRequest.Token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired reset token")
	}

	hashedPassword, err := HashPassword(setPasswordRequest.NewPassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	user, err := h.repo.GetUserByID(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user")
	}

	user.Password = hashedPassword
	if err := h.repo.UpdateUser(user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update password")
	}

	return c.NoContent(http.StatusOK)
}
