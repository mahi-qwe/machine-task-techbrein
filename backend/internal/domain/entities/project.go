package entities

import "time"

type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:150;not null" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	CreatedBy   uint      `gorm:"not null" json:"created_by"`
	Creator     User      `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
