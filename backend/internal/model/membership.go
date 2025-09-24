package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Membership struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
	CommunityID primitive.ObjectID `bson:"community_id" json:"community_id"`
	Role        string             `bson:"role,omitempty" json:"role,omitempty"`
}
