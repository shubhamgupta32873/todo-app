package models

import (
	"time"

	"gorm.io/gorm"
)

// Task represents a to-do task
type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `gorm:"index" json:"due_date"`
	DueTime     string    `json:"due_time"` // Format: "14:30"
	Priority    string    `gorm:"default:'medium'" json:"priority"` // low, medium, high
	Completed   bool      `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	User      User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Reminders []Reminder `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE" json:"reminders,omitempty"`
}

// TableName specifies the table name for Task model
func (Task) TableName() string {
	return "tasks"
}

// CreateTaskRequest is the request body for creating a task
type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	DueTime     string    `json:"due_time"` // Format: "14:30"
	Priority    string    `json:"priority"` // low, medium, high
}

// UpdateTaskRequest is the request body for updating a task
type UpdateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	DueTime     string    `json:"due_time"`
	Priority    string    `json:"priority"`
	Completed   bool      `json:"completed"`
}
