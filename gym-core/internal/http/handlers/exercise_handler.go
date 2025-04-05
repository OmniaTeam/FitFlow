package handlers

import (
	"errors"
	"gym-core/internal/repositories"
	"gym-core/internal/services"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExerciseHandler struct {
	exerciseService services.ExerciseService
}

func NewExerciseHandler(exerciseService services.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseService: exerciseService,
	}
}

// GetExerciseById godoc
// @Summary Get exercise by ID
// @Description Get exercise by ID with all related muscles
// @Tags exercises
// @Accept json
// @Produce json
// @Param id path int true "Exercise ID"
// @Success 200 {object} models.Exercise
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Exercise not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /exercises/{id} [get]
func (eh *ExerciseHandler) GetExerciseById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		slog.Error("Failed to parse exercise ID", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID format"})
		return
	}

	exercise, err := eh.exerciseService.GetExerciseById(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrorExerciseNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
			return
		}

		slog.Error("Failed to get exercise", "error", err, "exerciseId", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get exercise"})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

// GetAllExercises godoc
// @Summary Get all exercises
// @Description Get all exercises with their related muscles
// @Tags exercises
// @Accept json
// @Produce json
// @Success 200 {array} models.Exercise
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /exercises [get]
func (eh *ExerciseHandler) GetAllExercises(c *gin.Context) {
	exercises, err := eh.exerciseService.GetAllExercises(c.Request.Context())
	if err != nil {
		slog.Error("Failed to get exercises", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get exercises"})
		return
	}

	c.JSON(http.StatusOK, exercises)
}
