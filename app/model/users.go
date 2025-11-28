package model

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	FullName     string    `json:"full_name"`
	RoleID       string    `json:"role_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User        User     `json:"user"`
	Token       string   `json:"token"`
	Permissions []string `json:"permissions"`
}

// JWT claims includes permissions
type JWTClaims struct {
	UserID      string   `json:"user_id"`
	Username    string   `json:"username"`
	RoleID      string   `json:"role_id"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type Achievement struct {
	ID              string                 `json:"id"`
	StudentUserID   string                 `json:"student_user_id"`
	AchievementType string                 `json:"achievement_type"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	Details         map[string]interface{} `json:"details"`
	Status          string                 `json:"status"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}
