package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Achievement struct {
	ID              primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	StudentID       string                 `bson:"student_id" json:"student_id"`
	Title           string                 `bson:"title" json:"title"`
	AchievementType string                 `bson:"achievement_type" json:"achievement_type"`
	Description     string                 `bson:"description" json:"description"`
	Details         map[string]interface{} `bson:"details" json:"details"`
	Attachments     []Attachment           `bson:"attachments" json:"attachments"`
	Tags            []string               `bson:"tags" json:"tags"`
	Status          string                 `bson:"status" json:"status"`
	IsDeleted       bool                   `bson:"is_deleted" json:"-"`
	CreatedAt       time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `bson:"updated_at" json:"updated_at"`
}

type Attachment struct {
	FileName   string    `bson:"file_name" json:"file_name"`
	FileURL    string    `bson:"file_url" json:"file_url"`
	FileType   string    `bson:"file_type" json:"file_type"`
	UploadedAt time.Time `bson:"uploaded_at" json:"uploaded_at"`
}
