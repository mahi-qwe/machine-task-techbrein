package postgres

import (
	"errors"

	"techbrein-project-management/internal/domain/entities"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(project *entities.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepository) List(limit, offset int) ([]entities.Project, int64, error) {
	var projects []entities.Project
	var total int64

	if err := r.db.Model(&entities.Project{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Order("id desc").Limit(limit).Offset(offset).Find(&projects).Error; err != nil {
		return nil, 0, err
	}
	return projects, total, nil
}

func (r *ProjectRepository) GetByID(id uint) (*entities.Project, error) {
	var project entities.Project
	if err := r.db.First(&project, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) Update(project *entities.Project) error {
	return r.db.Save(project).Error
}

func (r *ProjectRepository) Delete(project *entities.Project) error {
	return r.db.Delete(project).Error
}
