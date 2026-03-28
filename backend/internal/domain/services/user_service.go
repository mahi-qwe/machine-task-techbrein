package services

import (
	"net/http"

	"techbrein-project-management/internal/domain/entities"
	"techbrein-project-management/internal/domain/repositories"
	"techbrein-project-management/internal/dto"
	"techbrein-project-management/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) Create(input dto.CreateUserRequest) (*entities.User, error) {
	existing, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, utils.NewAppError(http.StatusConflict, "a user with this email already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hash),
		Role:         entities.Role(input.Role),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) List(page, pageSize int) ([]entities.User, int64, error) {
	return s.userRepo.List(pageSize, (page-1)*pageSize)
}
