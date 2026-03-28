package services

import (
	"net/http"

	"techbrein-project-management/internal/domain/entities"
	"techbrein-project-management/internal/domain/repositories"
	"techbrein-project-management/internal/dto"
	"techbrein-project-management/internal/utils"
)

type ProjectService struct {
	projectRepo repositories.ProjectRepository
}

func NewProjectService(projectRepo repositories.ProjectRepository) *ProjectService {
	return &ProjectService{projectRepo: projectRepo}
}

func (s *ProjectService) Create(input dto.CreateProjectRequest, createdBy uint) (*entities.Project, error) {
	project := &entities.Project{
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   createdBy,
	}
	if err := s.projectRepo.Create(project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) List(page, pageSize int) ([]entities.Project, int64, error) {
	return s.projectRepo.List(pageSize, (page-1)*pageSize)
}

func (s *ProjectService) Update(projectID uint, input dto.UpdateProjectRequest) (*entities.Project, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, utils.NewAppError(http.StatusNotFound, "project not found")
	}

	if input.Name != nil {
		project.Name = *input.Name
	}
	if input.Description != nil {
		project.Description = *input.Description
	}

	if err := s.projectRepo.Update(project); err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) Delete(projectID uint) error {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return utils.NewAppError(http.StatusNotFound, "project not found")
	}
	return s.projectRepo.Delete(project)
}
