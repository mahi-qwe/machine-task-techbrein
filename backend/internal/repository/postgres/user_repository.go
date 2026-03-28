package postgres

import (
	"errors"

	"techbrein-project-management/internal/domain/entities"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *entities.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) List(limit, offset int) ([]entities.User, int64, error) {
	var users []entities.User
	var total int64

	if err := r.db.Model(&entities.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Order("id desc").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepository) GetByID(id uint) (*entities.User, error) {
	var user entities.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*entities.User, error) {
	var user entities.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
