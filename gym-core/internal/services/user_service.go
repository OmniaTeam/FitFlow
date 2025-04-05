package services

import (
	"context"
	"gym-core/internal/models"
	"gym-core/internal/repositories"
)

type UserService interface {
	GetUserById(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, id int, req models.UserUpdateRequest) (*models.User, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) GetUserById(ctx context.Context, id int) (*models.User, error) {
	return s.userRepository.GetUserById(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id int, req models.UserUpdateRequest) (*models.User, error) {
	return s.userRepository.UpdateUser(ctx, id, req)
}
