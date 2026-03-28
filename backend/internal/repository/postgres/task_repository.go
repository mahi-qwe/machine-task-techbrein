package postgres

import (
	"errors"

	"techbrein-project-management/internal/domain/entities"
	"techbrein-project-management/internal/domain/repositories"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *entities.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) List(limit, offset int, filter repositories.TaskFilter) ([]entities.Task, int64, error) {
	var tasks []entities.Task
	var total int64

	query := r.db.Model(&entities.Task{}).Preload("Project").Preload("Assignee")
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.AssignedTo != nil {
		query = query.Where("assigned_to = ?", *filter.AssignedTo)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("id desc").Limit(limit).Offset(offset).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TaskRepository) GetByID(id uint) (*entities.Task, error) {
	var task entities.Task
	if err := r.db.Preload("Project").Preload("Assignee").First(&task, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) Assign(taskID uint, assignedTo uint) error {
	return r.db.Model(&entities.Task{}).Where("id = ?", taskID).Update("assigned_to", assignedTo).Error
}

func (r *TaskRepository) Update(task *entities.Task) error {
	return r.db.Save(task).Error
}
