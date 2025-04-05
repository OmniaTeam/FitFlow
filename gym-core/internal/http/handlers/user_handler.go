package handlers

import (
	"errors"
	"gym-core/internal/http/utils"
	"gym-core/internal/models"
	"gym-core/internal/repositories"
	"gym-core/internal/services"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserProfile godoc
// @Summary Get user profile
// @Description Get user profile by ID from X-User-Id header
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/profile [get]
func (uh *UserHandler) GetUserProfile(c *gin.Context) {
	userId, err := utils.ExtractUserIdFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uh.userService.GetUserById(c.Request.Context(), userId)
	if err != nil {
		if errors.Is(err, repositories.ErrorUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		slog.Error("Failed to get user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user from header
// @Description Update user by ID from X-User-Id header
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.UserUpdateRequest true "User data for update"
// @Success 200 {object} models.User "User updated successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [put]
func (uh *UserHandler) UpdateUser(c *gin.Context) {
	userId, err := utils.ExtractUserIdFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data: " + err.Error()})
		return
	}

	updatedUser, err := uh.userService.UpdateUser(c.Request.Context(), userId, req)
	if err != nil {
		if errors.Is(err, repositories.ErrorUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		slog.Error("Failed to update user", "error", err, "userId", userId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}
