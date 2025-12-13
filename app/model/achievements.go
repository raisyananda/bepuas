package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Achievement struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty" json:"id,omitempty"`
	StudentID  string                 `bson:"student_id" json:"student_id"`
	Title      string                 `bson:"title" json:"title"`
	AType      string                 `bson:"achievementType" json:"achievementType"`
	Dscription string                 `bson:"description" json:"description"`
	Details    map[string]interface{} `bson:"details" json:"details"`
	Tags       []string               `bson:"tags" json:"tags"`
	CreatedAt  time.Time              `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time              `bson:"updated_at" json:"updated_at"`
}
