package model

import "time"

type AchievementReference struct {
	ID                 string    `json:"id"`
	StudentID          string    `json:"student_id"`
	MongoAchievementID string    `json:"mongo_achievement_id"`
	Status             string    `json:"status"` // draft, submitted, verified, rejected, deleted
	SubmittedAt        time.Time `json:"submitted_at"`
	VerifiedAt         time.Time `json:"verified_at"`
	VerifiedBy         string    `json:"verified_by"`
	RejectionNote      string    `json:"rejection_note"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type AchievementAdminView struct {
	ID              string    `json:"id"` // mongo id
	StudentID       string    `json:"student_id"`
	Status          string    `json:"status"`
	SubmittedAt     time.Time `json:"submitted_at"`
	Title           string    `json:"title"`
	AchievementType string    `json:"achievement_type"`
	Description     string    `json:"description"`
	Tags            []string  `json:"tags"`
	CreatedAt       time.Time `json:"created_at"`
}
