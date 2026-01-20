package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EmailVerification stores temporary email verification data before user registration
type EmailVerification struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email        string             `bson:"email" json:"email"`
	OTP          string             `bson:"otp" json:"-"` // Hidden from JSON
	OTPExpiresAt time.Time          `bson:"otp_expires_at" json:"otp_expires_at"`
	IsVerified   bool               `bson:"is_verified" json:"is_verified"` // true after OTP verified
	Nonce        string             `bson:"nonce" json:"-"`                 // Used in verification token to prevent replay
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}
