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
	Attachments     string                 `bson:"document_url" json:"document_url"`
	Tags            []string               `bson:"tags" json:"tags"`
	Status          string                 `bson:"status" json:"status"`
	IsDeleted       bool                   `bson:"is_deleted" json:"-"`
	CreatedAt       time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time              `bson:"updated_at" json:"updated_at"`
}
