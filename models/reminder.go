package models

import (
	"time"
)

// Reminder represents a reminder for a task
type Reminder struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	TaskID       uint      `gorm:"not null;index" json:"task_id"`
	ReminderTime time.Time `gorm:"index" json:"reminder_time"` // When to send the reminder
	ReminderType string    `gorm:"default:'in_app'" json:"reminder_type"` // in_app, email, sms
	Sent         bool      `gorm:"default:false" json:"sent"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Task Task `gorm:"foreignKey:TaskID" json:"task,omitempty"`
}

// TableName specifies the table name for Reminder model
func (Reminder) TableName() string {
	return "reminders"
}

// CreateReminderRequest is the request body for creating a reminder
type CreateReminderRequest struct {
	TaskID       uint      `json:"task_id" binding:"required"`
	ReminderTime time.Time `json:"reminder_time" binding:"required"`
	ReminderType string    `json:"reminder_type"` // in_app, email, sms
}
