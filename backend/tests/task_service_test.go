package tests

import (
	"testing"

	"techbrein-project-management/internal/domain/entities"
	"techbrein-project-management/internal/domain/repositories"
	"techbrein-project-management/internal/domain/services"
	"techbrein-project-management/internal/dto"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type stubProjectRepo struct {
	project *entities.Project
}

func (s *stubProjectRepo) Create(project *entities.Project) error { return nil }
func (s *stubProjectRepo) List(limit, offset int) ([]entities.Project, int64, error) {
	return nil, 0, nil
}
func (s *stubProjectRepo) GetByID(id uint) (*entities.Project, error) { return s.project, nil }
func (s *stubProjectRepo) Update(project *entities.Project) error     { return nil }
func (s *stubProjectRepo) Delete(project *entities.Project) error     { return nil }

type stubTaskRepo struct {
	task *entities.Task
}

func (s *stubTaskRepo) Create(task *entities.Task) error { s.task = task; return nil }
func (s *stubTaskRepo) List(limit, offset int, filter repositories.TaskFilter) ([]entities.Task, int64, error) {
	return nil, 0, nil
}
func (s *stubTaskRepo) GetByID(id uint) (*entities.Task, error) { return s.task, nil }
func (s *stubTaskRepo) Assign(taskID uint, assignedTo uint) error {
	if s.task != nil && s.task.ID == taskID {
		s.task.AssignedTo = &assignedTo
	}
	return nil
}
func (s *stubTaskRepo) Update(task *entities.Task) error { s.task = task; return nil }

func TestTaskStatusBlockedForNonAssignee(t *testing.T) {
	taskRepo := &stubTaskRepo{
		task: &entities.Task{
			ID:         10,
			ProjectID:  1,
			AssignedTo: ptr(uint(5)),
			Status:     entities.TaskStatusTodo,
		},
	}

	service := services.NewTaskService(
		&gorm.DB{},
		taskRepo,
		&stubProjectRepo{project: &entities.Project{ID: 1}},
		&stubUserRepo{user: &entities.User{ID: 5}},
	)

	_, err := service.UpdateStatus(10, 9, "developer", dto.UpdateTaskStatusRequest{Status: "done"})
	assert.Error(t, err)
	assert.Equal(t, "only admins or the assignee can update task status", err.Error())
}

func ptr[T any](value T) *T {
	return &value
}
