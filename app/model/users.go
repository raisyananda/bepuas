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

type UserResponse struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	FullName string  `json:"full_name"`
	RoleID   *string `json:"role_id"`
	RoleName *string `json:"role_name"`
}

type UpdateUserRequest struct {
	Lecturer *struct {
		LecturerID string `json:"lecturer_id"`
		Department string `json:"department"`
	} `json:"lecturer"`

	Student *struct {
		StudentID    string `json:"student_id"`
		ProgramStudy string `json:"program_study"`
		AcademicYear string `json:"academic_year"`
	} `json:"student"`
}

type LoginRequest struct {
	// Username string `json:"username"`
	Identifier string `json:"identifier"` // username/email
	Password   string `json:"password"`
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
