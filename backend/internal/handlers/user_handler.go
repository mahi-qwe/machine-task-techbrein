package handlers

import (
	"net/http"

	"techbrein-project-management/internal/domain/services"
	"techbrein-project-management/internal/dto"
	"techbrein-project-management/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(c *gin.Context) {
	var request dto.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := h.service.Create(request)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) List(c *gin.Context) {
	page, pageSize := parsePagination(c)
	users, total, err := h.service.List(page, pageSize)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Items: users,
		Meta:  utils.BuildPagination(page, pageSize, total),
	})
}
