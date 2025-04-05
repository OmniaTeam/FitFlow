package services

import (
	"context"
	"gym-core/internal/models"
	"gym-core/internal/repositories"
)

type ProgramService interface {
	GetProgramByUserId(ctx context.Context, userId int) (*models.Program, error)
}

type programService struct {
	programRepository repositories.ProgramRepository
}

func NewProgramService(programRepository repositories.ProgramRepository) ProgramService {
	return &programService{
		programRepository: programRepository,
	}
}

func (s *programService) GetProgramByUserId(ctx context.Context, userId int) (*models.Program, error) {
	return s.programRepository.GetProgramByUserId(ctx, userId)
}
