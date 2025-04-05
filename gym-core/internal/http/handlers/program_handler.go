package handlers

import (
	"errors"
	"gym-core/internal/http/utils"
	"gym-core/internal/repositories"
	"gym-core/internal/services"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProgramHandler struct {
	programService services.ProgramService
}

func NewProgramHandler(programService services.ProgramService) *ProgramHandler {
	return &ProgramHandler{
		programService: programService,
	}
}

// GetUserProgram godoc
// @Summary Get user program
// @Description Get user program with all workouts, exercises and sets by user ID from X-User-Id header
// @Tags programs
// @Accept json
// @Produce json
// @Success 200 {object} models.Program
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Program not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /programs/user [get]
func (ph *ProgramHandler) GetUserProgram(c *gin.Context) {
	userId, err := utils.ExtractUserIdFromHeader(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	program, err := ph.programService.GetProgramByUserId(c.Request.Context(), userId)
	if err != nil {
		if errors.Is(err, repositories.ErrorProgramNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Program not found"})
			return
		}

		slog.Error("Failed to get program", "error", err, "userId", userId)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get program"})
		return
	}

	c.JSON(http.StatusOK, program)
}
