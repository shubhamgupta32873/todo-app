package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shubhamgupta32873/todo-app/handlers"
	"github.com/shubhamgupta32873/todo-app/middleware"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine) {
	// Public routes (no authentication required)
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Protected routes (authentication required)
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Task routes
		tasks := api.Group("/tasks")
		{
			tasks.POST("", handlers.CreateTask)
			tasks.GET("", handlers.GetUserTasks)
			tasks.GET("/date", handlers.GetTasksByDate)
			tasks.GET("/:id", handlers.GetTask)
			tasks.PUT("/:id", handlers.UpdateTask)
			tasks.DELETE("/:id", handlers.DeleteTask)
		}

		// Reminder routes
		reminders := api.Group("/reminders")
		{
			reminders.POST("", handlers.CreateReminder)
			reminders.GET("/:id", handlers.GetReminder)
			reminders.DELETE("/:id", handlers.DeleteReminder)
		}

		// Task reminders
		taskReminders := api.Group("/tasks/:id/reminders")
		{
			taskReminders.GET("", handlers.GetTaskReminders)
		}
	}

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "API is running",
		})
	})
}
