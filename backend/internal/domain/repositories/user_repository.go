package repositories

import "techbrein-project-management/internal/domain/entities"

type UserRepository interface {
	Create(user *entities.User) error
	List(limit, offset int) ([]entities.User, int64, error)
	GetByID(id uint) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
}
