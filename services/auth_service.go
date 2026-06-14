package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shubhamgupta32873/todo-app/config"
	"github.com/shubhamgupta32873/todo-app/models"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser creates a new user
func RegisterUser(email, password string) (*models.User, error) {
	db := config.GetDB()

	// Check if user already exists
	var existingUser models.User
	if result := db.Where("email = ?", email).First(&existingUser); result.RowsAffected > 0 {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create user
	user := models.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if result := db.Create(&user); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// LoginUser authenticates a user and returns JWT token
func LoginUser(email, password string) (string, error) {
	db := config.GetDB()

	// Find user by email
	var user models.User
	if result := db.Where("email = ?", email).First(&user); result.RowsAffected == 0 {
		return "", errors.New("user not found")
	}

	// Compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns user ID
func ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("user_id not found in token")
	}

	return uint(userID), nil
}
