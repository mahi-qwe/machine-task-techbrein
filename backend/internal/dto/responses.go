package dto

type ErrorResponse struct {
	Message string `json:"message"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type PaginatedResponse struct {
	Items any            `json:"items"`
	Meta  PaginationMeta `json:"meta"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}
