package services

import (
	"errors"
	"time"

	"github.com/shubhamgupta32873/todo-app/config"
	"github.com/shubhamgupta32873/todo-app/models"
)

// CreateTask creates a new task
func CreateTask(userID uint, req models.CreateTaskRequest) (*models.Task, error) {
	db := config.GetDB()

	task := models.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		DueTime:     req.DueTime,
		Priority:    req.Priority,
		Completed:   false,
	}

	if result := db.Create(&task); result.Error != nil {
		return nil, result.Error
	}

	return &task, nil
}

// GetTaskByID fetches a task by ID
func GetTaskByID(taskID, userID uint) (*models.Task, error) {
	db := config.GetDB()

	var task models.Task
	if result := db.Where("id = ? AND user_id = ?", taskID, userID).First(&task); result.RowsAffected == 0 {
		return nil, errors.New("task not found")
	}

	return &task, nil
}

// GetUserTasks fetches all tasks for a user
func GetUserTasks(userID uint) ([]models.Task, error) {
	db := config.GetDB()

	var tasks []models.Task
	if result := db.Where("user_id = ?", userID).Order("due_date ASC").Find(&tasks); result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

// GetTasksByDate fetches tasks for a specific date
func GetTasksByDate(userID uint, date time.Time) ([]models.Task, error) {
	db := config.GetDB()

	// Start and end of the day
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)

	var tasks []models.Task
	if result := db.Where("user_id = ? AND due_date BETWEEN ? AND ?", userID, startOfDay, endOfDay).
		Order("due_time ASC").Find(&tasks); result.Error != nil {
		return nil, result.Error
	}

	return tasks, nil
}

// UpdateTask updates a task
func UpdateTask(taskID, userID uint, req models.UpdateTaskRequest) (*models.Task, error) {
	db := config.GetDB()

	task, err := GetTaskByID(taskID, userID)
	if err != nil {
		return nil, err
	}

	// Update only provided fields
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if !req.DueDate.IsZero() {
		task.DueDate = req.DueDate
	}
	if req.DueTime != "" {
		task.DueTime = req.DueTime
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	task.Completed = req.Completed

	if result := db.Save(task); result.Error != nil {
		return nil, result.Error
	}

	return task, nil
}

// DeleteTask deletes a task
func DeleteTask(taskID, userID uint) error {
	db := config.GetDB()

	task, err := GetTaskByID(taskID, userID)
	if err != nil {
		return err
	}

	if result := db.Delete(task); result.Error != nil {
		return result.Error
	}

	return nil
}
