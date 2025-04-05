package services

import (
	"context"
	"gym-core/internal/models"
	"gym-core/internal/repositories"
)

type ExerciseService interface {
	GetExerciseById(ctx context.Context, id int) (*models.Exercise, error)
	GetAllExercises(ctx context.Context) ([]models.Exercise, error)
}

type exerciseService struct {
	exerciseRepository repositories.ExerciseRepository
}

func NewExerciseService(exerciseRepository repositories.ExerciseRepository) ExerciseService {
	return &exerciseService{
		exerciseRepository: exerciseRepository,
	}
}

func (s *exerciseService) GetExerciseById(ctx context.Context, id int) (*models.Exercise, error) {
	return s.exerciseRepository.GetExerciseById(ctx, id)
}

func (s *exerciseService) GetAllExercises(ctx context.Context) ([]models.Exercise, error) {
	return s.exerciseRepository.GetAllExercises(ctx)
}
