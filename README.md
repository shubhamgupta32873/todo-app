# Todo App with Reminders

A web-based task management application with automated reminder system built with Go, Gin, and PostgreSQL.

## Features

✅ User Authentication (JWT)  
✅ Create, Read, Update, Delete Tasks  
✅ Organize tasks by dates  
✅ Set reminders for tasks (email, SMS, in-app)  
✅ Automated reminder scheduler  
✅ RESTful API  

## Tech Stack

- **Backend**: Go + Gin Framework
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (jsonwebtoken)
- **Scheduler**: robfig/cron
- **Password Hashing**: bcrypt

## Project Structure

```
todo-app/
├── config/           # Database configuration
├── database/         # Migrations
├── handlers/         # HTTP request handlers
├── middleware/       # Authentication middleware
├── models/           # Data models (User, Task, Reminder)
├── routes/           # API routes
├── services/         # Business logic
├── main.go           # Application entry point
├── .env.example      # Environment variables template
├── go.mod            # Go dependencies
└── README.md         # This file
```

## Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Git

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/shubhamgupta32873/todo-app.git
cd todo-app
```

### 2. Set up environment variables

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=todo_app

# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# JWT Configuration
JWT_SECRET=your_secret_key_change_in_production
```

### 3. Create PostgreSQL database

```bash
createdb todo_app
```

Or using psql:

```sql
CREATE DATABASE todo_app;
```

### 4. Install dependencies

```bash
go mod download
go mod tidy
```

### 5. Run the application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Using Docker

### Prerequisites
- Docker
- Docker Compose

### Run with Docker

```bash
docker-compose up
```

This will start both PostgreSQL and the Go application.

## API Endpoints

### Authentication

**Register**
```
POST /api/auth/register
Body: {
  "email": "user@example.com",
  "password": "password123"
}
```

**Login**
```
POST /api/auth/login
Body: {
  "email": "user@example.com",
  "password": "password123"
}
Response: {
  "token": "jwt_token_here"
}
```

### Tasks

All task endpoints require JWT authentication in header:
```
Authorization: Bearer <jwt_token>
```

**Create Task**
```
POST /api/tasks
Body: {
  "title": "Complete project",
  "description": "Finish the Go backend",
  "due_date": "2024-12-31T00:00:00Z",
  "due_time": "14:30",
  "priority": "high"
}
```

**Get All Tasks**
```
GET /api/tasks
```

**Get Tasks by Date**
```
GET /api/tasks/date?date=2024-12-31
```

**Get Single Task**
```
GET /api/tasks/:id
```

**Update Task**
```
PUT /api/tasks/:id
Body: {
  "completed": true,
  "title": "Updated title"
}
```

**Delete Task**
```
DELETE /api/tasks/:id
```

### Reminders

**Create Reminder**
```
POST /api/reminders
Body: {
  "task_id": 1,
  "reminder_time": "2024-12-31T14:00:00Z",
  "reminder_type": "in_app"  // in_app, email, sms
}
```

**Get Task Reminders**
```
GET /api/tasks/:id/reminders
```

**Get Single Reminder**
```
GET /api/reminders/:id
```

**Delete Reminder**
```
DELETE /api/reminders/:id
```

### Health Check

```
GET /health
```

## Reminder Types

- **in_app**: In-application notification (logged to console)
- **email**: Email notification (framework ready for implementation)
- **sms**: SMS notification (framework ready for implementation)

## How Reminders Work

1. The application starts a background scheduler using robfig/cron
2. Every minute, the scheduler checks the database for due reminders
3. For each reminder where `reminder_time <= now()` and `sent = false`:
   - A notification is sent based on the reminder type
   - The reminder is marked as sent in the database
4. The process repeats every minute

## Database Schema

### Users Table
```sql
- id (PRIMARY KEY)
- email (UNIQUE, NOT NULL)
- password (NOT NULL, bcrypt hashed)
- created_at
- updated_at
```

### Tasks Table
```sql
- id (PRIMARY KEY)
- user_id (FOREIGN KEY -> users)
- title (NOT NULL)
- description
- due_date
- due_time
- priority (low, medium, high)
- completed (boolean, default: false)
- created_at
- updated_at
```

### Reminders Table
```sql
- id (PRIMARY KEY)
- task_id (FOREIGN KEY -> tasks)
- reminder_time
- reminder_type (in_app, email, sms)
- sent (boolean, default: false)
- created_at
- updated_at
```

## Development

### Running tests

```bash
go test ./...
```

### Database migrations

Migrations run automatically on application startup using GORM's AutoMigrate feature.

## Build Commands

```bash
# Install dependencies
make install

# Run the application
make run

# Build binary
make build

# Run tests
make test

# Start Docker containers
make docker-up

# Stop Docker containers
make docker-down

# View Docker logs
make docker-logs

# Clean build artifacts
make clean
```

## Common Issues

### Database connection error

Make sure PostgreSQL is running and credentials in `.env` are correct.

```bash
psql -h localhost -U postgres -d todo_app
```

### Port already in use

Change `SERVER_PORT` in `.env` or kill the process using the port.

### JWT token errors

- Ensure `JWT_SECRET` in `.env` is set
- Check that the token is sent in the `Authorization` header as `Bearer <token>`
- Tokens expire after 24 hours

## Future Enhancements

- [ ] Email notification implementation (SMTP)
- [ ] SMS notification implementation (Twilio)
- [ ] React/Vue frontend
- [ ] Unit and integration tests
- [ ] API documentation (Swagger/OpenAPI)
- [ ] Task categories/tags
- [ ] Recurring tasks
- [ ] Task sharing
- [ ] WebSocket for real-time updates
- [ ] Rate limiting

## License

MIT

## Author

Shubham Gupta (shubhamgupta32873)

## Support

For issues and questions, please create an issue on GitHub.
