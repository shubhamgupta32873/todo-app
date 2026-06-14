package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/shubhamgupta32873/todo-app/config"
	"github.com/shubhamgupta32873/todo-app/models"
)

var reminderScheduler *cron.Cron

// CreateReminder creates a new reminder
func CreateReminder(req models.CreateReminderRequest) (*models.Reminder, error) {
	db := config.GetDB()

	reminder := models.Reminder{
		TaskID:       req.TaskID,
		ReminderTime: req.ReminderTime,
		ReminderType: req.ReminderType,
		Sent:         false,
	}

	if result := db.Create(&reminder); result.Error != nil {
		return nil, result.Error
	}

	return &reminder, nil
}

// GetReminderByID fetches a reminder by ID
func GetReminderByID(reminderID uint) (*models.Reminder, error) {
	db := config.GetDB()

	var reminder models.Reminder
	if result := db.Where("id = ?", reminderID).Preload("Task").First(&reminder); result.RowsAffected == 0 {
		return nil, errors.New("reminder not found")
	}

	return &reminder, nil
}

// GetTaskReminders fetches all reminders for a task
func GetTaskReminders(taskID uint) ([]models.Reminder, error) {
	db := config.GetDB()

	var reminders []models.Reminder
	if result := db.Where("task_id = ?", taskID).Order("reminder_time ASC").Find(&reminders); result.Error != nil {
		return nil, result.Error
	}

	return reminders, nil
}

// DeleteReminder deletes a reminder
func DeleteReminder(reminderID uint) error {
	db := config.GetDB()

	reminder, err := GetReminderByID(reminderID)
	if err != nil {
		return err
	}

	if result := db.Delete(reminder); result.Error != nil {
		return result.Error
	}

	return nil
}

// StartReminderScheduler starts the background job for reminders
func StartReminderScheduler() {
	reminderScheduler = cron.New()

	// Run every minute to check for due reminders
	reminderScheduler.AddFunc("* * * * *", func() {
		processDueReminders()
	})

	reminderScheduler.Start()
	log.Println("✅ Reminder scheduler started")
}

// processDueReminders checks for reminders that are due and sends them
func processDueReminders() {
	db := config.GetDB()

	// Find all unsent reminders that are due
	var reminders []models.Reminder
	now := time.Now()

	if result := db.Where("sent = ? AND reminder_time <= ?", false, now).
		Preload("Task").
		Preload("Task.User").
		Find(&reminders); result.Error != nil {
		log.Printf("Error fetching reminders: %v", result.Error)
		return
	}

	// Process each reminder
	for _, reminder := range reminders {
		log.Printf("Processing reminder #%d for task #%d", reminder.ID, reminder.TaskID)

		// Send notification based on type
		var err error
		switch reminder.ReminderType {
		case "email":
			err = sendEmailReminder(reminder)
		case "sms":
			err = sendSMSReminder(reminder)
		case "in_app":
			err = sendInAppReminder(reminder)
		default:
			err = sendInAppReminder(reminder) // Default to in-app
		}

		// Mark as sent if successful
		if err != nil {
			log.Printf("Error sending reminder #%d: %v", reminder.ID, err)
		} else {
			db.Model(&reminder).Update("sent", true)
			log.Printf("✅ Reminder #%d sent successfully", reminder.ID)
		}
	}
}

// sendEmailReminder sends email reminder
func sendEmailReminder(reminder models.Reminder) error {
	// TODO: Implement email sending using nodemailer equivalent in Go
	// For now, just log it
	log.Printf("📧 Email reminder: Task '%s' is due at %s", reminder.Task.Title, reminder.Task.DueTime)
	return nil
}

// sendSMSReminder sends SMS reminder
func sendSMSReminder(reminder models.Reminder) error {
	// TODO: Implement SMS sending using Twilio equivalent in Go
	// For now, just log it
	log.Printf("📱 SMS reminder: Task '%s' is due at %s", reminder.Task.Title, reminder.Task.DueTime)
	return nil
}

// sendInAppReminder stores in-app notification
func sendInAppReminder(reminder models.Reminder) error {
	log.Printf("🔔 In-app reminder: Task '%s' is due at %s", reminder.Task.Title, reminder.Task.DueTime)
	// In-app notifications could be stored in a separate notifications table
	return nil
}

// StopReminderScheduler stops the background scheduler
func StopReminderScheduler() {
	if reminderScheduler != nil {
		reminderScheduler.Stop()
		log.Println("⏸️ Reminder scheduler stopped")
	}
}
