package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"techbrein-project-management/internal/utils"

	"github.com/gin-gonic/gin"
)

func respondError(c *gin.Context, err error) {
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.StatusCode, gin.H{"message": appErr.Message})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
}

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
