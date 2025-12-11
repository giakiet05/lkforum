package repo

import (
	"context"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PasswordResetRepo interface {
	Create(ctx context.Context, reset *model.PasswordReset) (*model.PasswordReset, error)
	GetByEmail(ctx context.Context, email string) (*model.PasswordReset, error)
	Update(ctx context.Context, reset *model.PasswordReset) (*model.PasswordReset, error)
	Delete(ctx context.Context, email string) error
}

type passwordResetRepo struct {
	collection *mongo.Collection
}

func NewPasswordResetRepo(db *mongo.Database) PasswordResetRepo {
	return &passwordResetRepo{collection: db.Collection(config.PasswordResetColName)}
}

func (r *passwordResetRepo) Create(ctx context.Context, reset *model.PasswordReset) (*model.PasswordReset, error) {
	result, err := r.collection.InsertOne(ctx, reset)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		reset.ID = oid
	}

	return reset, nil
}

func (r *passwordResetRepo) GetByEmail(ctx context.Context, email string) (*model.PasswordReset, error) {
	var reset model.PasswordReset
	filter := bson.M{"email": email}

	err := r.collection.FindOne(ctx, filter).Decode(&reset)
	if err != nil {
		return nil, err
	}

	return &reset, nil
}

func (r *passwordResetRepo) Update(ctx context.Context, reset *model.PasswordReset) (*model.PasswordReset, error) {
	filter := bson.M{"_id": reset.ID}
	update := bson.M{"$set": reset}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return reset, nil
}

func (r *passwordResetRepo) Delete(ctx context.Context, email string) error {
	filter := bson.M{"email": email}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
