package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shubhamgupta32873/todo-app/models"
	"github.com/shubhamgupta32873/todo-app/services"
)

// CreateReminder handles reminder creation
func CreateReminder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req models.CreateReminderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the task belongs to the user
	task, err := services.GetTaskByID(req.TaskID, uint(userID.(float64)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != uint(userID.(float64)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to create reminder for this task"})
		return
	}

	reminder, err := services.CreateReminder(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reminder)
}

// GetReminder handles fetching a single reminder
func GetReminder(c *gin.Context) {
	reminderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reminder id"})
		return
	}

	reminder, err := services.GetReminderByID(uint(reminderID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminder)
}

// GetTaskReminders handles fetching all reminders for a task
func GetTaskReminders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	// Verify task ownership
	task, err := services.GetTaskByID(uint(taskID), uint(userID.(float64)))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if task.UserID != uint(userID.(float64)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to view reminders for this task"})
		return
	}

	reminders, err := services.GetTaskReminders(uint(taskID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reminders)
}

// DeleteReminder handles reminder deletion
func DeleteReminder(c *gin.Context) {
	reminderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reminder id"})
		return
	}

	err = services.DeleteReminder(uint(reminderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reminder deleted successfully"})
}
