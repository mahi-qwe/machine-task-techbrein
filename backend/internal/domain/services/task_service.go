package services

import (
	"net/http"
	"time"

	"techbrein-project-management/internal/domain/entities"
	"techbrein-project-management/internal/domain/repositories"
	"techbrein-project-management/internal/dto"
	"techbrein-project-management/internal/utils"

	"gorm.io/gorm"
)

type TaskService struct {
	db          *gorm.DB
	taskRepo    repositories.TaskRepository
	projectRepo repositories.ProjectRepository
	userRepo    repositories.UserRepository
}

func NewTaskService(db *gorm.DB, taskRepo repositories.TaskRepository, projectRepo repositories.ProjectRepository, userRepo repositories.UserRepository) *TaskService {
	return &TaskService{
		db:          db,
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

func (s *TaskService) Create(input dto.CreateTaskRequest) (*entities.Task, error) {
	project, err := s.projectRepo.GetByID(input.ProjectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, utils.NewAppError(http.StatusNotFound, "project not found")
	}

	if input.AssignedTo != nil {
		user, err := s.userRepo.GetByID(*input.AssignedTo)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, utils.NewAppError(http.StatusNotFound, "assignee not found")
		}
	}

	var dueDate *time.Time
	if input.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", input.DueDate)
		if err != nil {
			return nil, utils.NewAppError(http.StatusBadRequest, "due_date must be in YYYY-MM-DD format")
		}
		dueDate = &parsed
	}

	task := &entities.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      entities.TaskStatusTodo,
		ProjectID:   input.ProjectID,
		AssignedTo:  input.AssignedTo,
		DueDate:     dueDate,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(task).Error
	})
	if err != nil {
		return nil, err
	}

	return s.taskRepo.GetByID(task.ID)
}

func (s *TaskService) Assign(taskID uint, input dto.AssignTaskRequest) (*entities.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, utils.NewAppError(http.StatusNotFound, "task not found")
	}

	user, err := s.userRepo.GetByID(input.AssignedTo)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.NewAppError(http.StatusNotFound, "assignee not found")
	}

	if err := s.taskRepo.Assign(taskID, input.AssignedTo); err != nil {
		return nil, err
	}
	return s.taskRepo.GetByID(taskID)
}

func (s *TaskService) UpdateStatus(taskID uint, userID uint, role string, input dto.UpdateTaskStatusRequest) (*entities.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, utils.NewAppError(http.StatusNotFound, "task not found")
	}
	if role != string(entities.RoleAdmin) && (task.AssignedTo == nil || *task.AssignedTo != userID) {
		return nil, utils.NewAppError(http.StatusForbidden, "only admins or the assignee can update task status")
	}

	task.Status = entities.TaskStatus(input.Status)
	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}
	return s.taskRepo.GetByID(taskID)
}

func (s *TaskService) Update(taskID uint, input dto.UpdateTaskRequest) (*entities.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, utils.NewAppError(http.StatusNotFound, "task not found")
	}

	if input.ProjectID != nil {
		project, err := s.projectRepo.GetByID(*input.ProjectID)
		if err != nil {
			return nil, err
		}
		if project == nil {
			return nil, utils.NewAppError(http.StatusNotFound, "project not found")
		}
		task.ProjectID = *input.ProjectID
	}

	if input.AssignedTo != nil {
		user, err := s.userRepo.GetByID(*input.AssignedTo)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, utils.NewAppError(http.StatusNotFound, "assignee not found")
		}
		task.AssignedTo = input.AssignedTo
	}

	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Description = *input.Description
	}
	if input.Status != nil {
		task.Status = entities.TaskStatus(*input.Status)
	}
	if input.DueDate != nil {
		parsed, err := time.Parse("2006-01-02", *input.DueDate)
		if err != nil {
			return nil, utils.NewAppError(http.StatusBadRequest, "due_date must be in YYYY-MM-DD format")
		}
		task.DueDate = &parsed
	}

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}
	return s.taskRepo.GetByID(taskID)
}

func (s *TaskService) List(page, pageSize int, filter repositories.TaskFilter) ([]entities.Task, int64, error) {
	return s.taskRepo.List(pageSize, (page-1)*pageSize, filter)
}
