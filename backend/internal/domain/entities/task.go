package entities

import "time"

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:150;not null" json:"title"`
	Description string     `gorm:"type:text;not null" json:"description"`
	Status      TaskStatus `gorm:"type:varchar(30);default:'todo';not null" json:"status"`
	ProjectID   uint       `gorm:"not null;index" json:"project_id"`
	Project     Project    `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	AssignedTo  *uint      `gorm:"index" json:"assigned_to"`
	Assignee    *User      `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
