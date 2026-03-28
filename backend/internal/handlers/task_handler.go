package handlers

import (
	"net/http"
	"strconv"

	"techbrein-project-management/internal/domain/repositories"
	"techbrein-project-management/internal/domain/services"
	"techbrein-project-management/internal/dto"
	"techbrein-project-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var request dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	task, err := h.service.Create(request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) Assign(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request dto.AssignTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	task, err := h.service.Assign(uint(id), request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request dto.UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	role, _ := c.Get("role")
	task, err := h.service.UpdateStatus(uint(id), c.GetUint("user_id"), role.(string), request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	task, err := h.service.Update(uint(id), request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) List(c *gin.Context) {
	page, pageSize := parsePagination(c)
	filter := repositories.TaskFilter{}

	if projectID := c.Query("project_id"); projectID != "" {
		parsed, _ := strconv.Atoi(projectID)
		value := uint(parsed)
		filter.ProjectID = &value
	}
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}
	if assignedTo := c.Query("assigned_to"); assignedTo != "" {
		parsed, _ := strconv.Atoi(assignedTo)
		value := uint(parsed)
		filter.AssignedTo = &value
	}

	tasks, total, err := h.service.List(page, pageSize, filter)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items: tasks,
		Meta:  utils.BuildPagination(page, pageSize, total),
	})
}
