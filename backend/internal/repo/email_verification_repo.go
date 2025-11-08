package repo

import (
	"context"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type EmailVerificationRepo interface {
	Create(ctx context.Context, verification *model.EmailVerification) (*model.EmailVerification, error)
	GetByEmail(ctx context.Context, email string) (*model.EmailVerification, error)
	Update(ctx context.Context, verification *model.EmailVerification) (*model.EmailVerification, error)
	Delete(ctx context.Context, email string) error
}

type emailVerificationRepo struct {
	collection *mongo.Collection
}

func NewEmailVerificationRepo(db *mongo.Database) EmailVerificationRepo {
	return &emailVerificationRepo{collection: db.Collection(config.EmailVerificationColName)}
}

func (r *emailVerificationRepo) Create(ctx context.Context, verification *model.EmailVerification) (*model.EmailVerification, error) {
	result, err := r.collection.InsertOne(ctx, verification)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		verification.ID = oid
	}

	return verification, nil
}

func (r *emailVerificationRepo) GetByEmail(ctx context.Context, email string) (*model.EmailVerification, error) {
	var verification model.EmailVerification
	filter := bson.M{"email": email}

	err := r.collection.FindOne(ctx, filter).Decode(&verification)
	if err != nil {
		return nil, err
	}

	return &verification, nil
}

func (r *emailVerificationRepo) Update(ctx context.Context, verification *model.EmailVerification) (*model.EmailVerification, error) {
	filter := bson.M{"_id": verification.ID}
	update := bson.M{"$set": verification}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return verification, nil
}

func (r *emailVerificationRepo) Delete(ctx context.Context, email string) error {
	filter := bson.M{"email": email}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
