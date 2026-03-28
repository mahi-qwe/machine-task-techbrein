package utils

import "techbrein-project-management/internal/dto"

func BuildPagination(page, pageSize int, total int64) dto.PaginationMeta {
	totalPages := total / int64(pageSize)
	if total%int64(pageSize) != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	return dto.PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}
