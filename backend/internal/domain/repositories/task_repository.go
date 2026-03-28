package repositories

import "techbrein-project-management/internal/domain/entities"

type TaskFilter struct {
	ProjectID  *uint
	Status     *string
	AssignedTo *uint
}

type TaskRepository interface {
	Create(task *entities.Task) error
	List(limit, offset int, filter TaskFilter) ([]entities.Task, int64, error)
	GetByID(id uint) (*entities.Task, error)
	Assign(taskID uint, assignedTo uint) error
	Update(task *entities.Task) error
}
