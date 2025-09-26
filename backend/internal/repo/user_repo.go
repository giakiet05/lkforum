package repo

import (
	"context"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	GetAll(ctx context.Context) ([]*model.User, error)
}

type userRepo struct {
	userCollection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) UserRepo {
	return &userRepo{userCollection: db.Collection(config.UserColName)}
}

func (r *userRepo) GetAll(ctx context.Context) ([]*model.User, error) {
	cursor, err := r.userCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
