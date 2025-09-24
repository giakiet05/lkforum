package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Report struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ReporterID  primitive.ObjectID `bson:"reporter_id,omitempty" json:"reporter_id,omitempty"`
	TargetID    primitive.ObjectID `bson:"target_id,omitempty" json:"target_id,omitempty"`
	TargetType  string             `bson:"target_type,omitempty" json:"target_type,omitempty"`
	Reason      string             `bson:"reason,omitempty" json:"reason,omitempty"`
	Description *string            `bson:"description,omitempty" json:"description,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
}
