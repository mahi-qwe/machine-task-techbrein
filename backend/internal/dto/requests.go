package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=120"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=admin developer"`
}

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=150"`
	Description string `json:"description" binding:"required,min=5,max=2000"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=3,max=150"`
	Description *string `json:"description" binding:"omitempty,min=5,max=2000"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=150"`
	Description string `json:"description" binding:"required,min=5,max=2000"`
	ProjectID   uint   `json:"project_id" binding:"required"`
	AssignedTo  *uint  `json:"assigned_to"`
	DueDate     string `json:"due_date"`
}

type AssignTaskRequest struct {
	AssignedTo uint `json:"assigned_to" binding:"required"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=todo in_progress done"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=3,max=150"`
	Description *string `json:"description" binding:"omitempty,min=5,max=2000"`
	ProjectID   *uint   `json:"project_id"`
	AssignedTo  *uint   `json:"assigned_to"`
	DueDate     *string `json:"due_date"`
	Status      *string `json:"status" binding:"omitempty,oneof=todo in_progress done"`
}
