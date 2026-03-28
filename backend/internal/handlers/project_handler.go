package handlers

import (
	"net/http"
	"strconv"

	"techbrein-project-management/internal/domain/services"
	"techbrein-project-management/internal/dto"
	"techbrein-project-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var request dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	userID := c.GetUint("user_id")
	project, err := h.service.Create(request, userID)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) List(c *gin.Context) {
	page, pageSize := parsePagination(c)
	projects, total, err := h.service.List(page, pageSize)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items: projects,
		Meta:  utils.BuildPagination(page, pageSize, total),
	})
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var request dto.UpdateProjectRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	project, err := h.service.Update(uint(id), request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.Delete(uint(id)); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
