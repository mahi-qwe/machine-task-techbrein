package tests

import (
	"testing"

	"techbrein-project-management/internal/config"
	"techbrein-project-management/internal/domain/entities"
	"techbrein-project-management/internal/domain/services"
	"techbrein-project-management/internal/dto"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type stubUserRepo struct {
	user *entities.User
}

func (s *stubUserRepo) Create(user *entities.User) error { return nil }
func (s *stubUserRepo) List(limit, offset int) ([]entities.User, int64, error) {
	return nil, 0, nil
}
func (s *stubUserRepo) GetByID(id uint) (*entities.User, error) { return s.user, nil }
func (s *stubUserRepo) GetByEmail(email string) (*entities.User, error) {
	if s.user != nil && s.user.Email == email {
		return s.user, nil
	}
	return nil, nil
}

func TestAuthServiceLoginSuccess(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("Password@123"), bcrypt.DefaultCost)
	assert.NoError(t, err)

	repo := &stubUserRepo{
		user: &entities.User{
			ID:           1,
			Name:         "Admin",
			Email:        "admin@example.com",
			PasswordHash: string(hash),
			Role:         entities.RoleAdmin,
		},
	}

	service := services.NewAuthService(repo, config.Config{
		JWTSecret:       "test-secret",
		TokenTTLMinutes: 60,
	})

	response, err := service.Login(dto.LoginRequest{
		Email:    "admin@example.com",
		Password: "Password@123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, response.AccessToken)
	assert.Equal(t, "bearer", response.TokenType)
}
