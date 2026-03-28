package services

import (
	"net/http"

	"techbrein-project-management/internal/config"
	"techbrein-project-management/internal/domain/repositories"
	"techbrein-project-management/internal/dto"
	"techbrein-project-management/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repositories.UserRepository
	cfg      config.Config
}

func NewAuthService(userRepo repositories.UserRepository, cfg config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

func (s *AuthService) Login(input dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.NewAppError(http.StatusUnauthorized, "invalid email or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, utils.NewAppError(http.StatusUnauthorized, "invalid email or password")
	}

	token, err := utils.GenerateToken(s.cfg.JWTSecret, s.cfg.TokenTTLMinutes, user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: token,
		TokenType:   "bearer",
	}, nil
}
