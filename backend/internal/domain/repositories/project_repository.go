package repositories

import "techbrein-project-management/internal/domain/entities"

type ProjectRepository interface {
	Create(project *entities.Project) error
	List(limit, offset int) ([]entities.Project, int64, error)
	GetByID(id uint) (*entities.Project, error)
	Update(project *entities.Project) error
	Delete(project *entities.Project) error
}
