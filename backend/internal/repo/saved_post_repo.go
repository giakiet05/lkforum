package repo

import (
	"context"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SavedPostRepo interface {
	Save(ctx context.Context, userID, postID primitive.ObjectID) error
	Unsave(ctx context.Context, userID, postID primitive.ObjectID) error
	GetByUserID(ctx context.Context, userID primitive.ObjectID, findOptions *FindOptions) ([]*model.SavedPost, int64, error)
	IsSaved(ctx context.Context, userID, postID primitive.ObjectID) (bool, error)
}

type savedPostRepo struct {
	collection *mongo.Collection
}

func NewSavedPostRepo(db *mongo.Database) SavedPostRepo {
	return &savedPostRepo{
		collection: db.Collection(config.SavedPostColName),
	}
}

func (r *savedPostRepo) Save(ctx context.Context, userID, postID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID, "post_id": postID}
	update := bson.M{
		"$setOnInsert": bson.M{
			"user_id":  userID,
			"post_id":  postID,
			"saved_at": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *savedPostRepo) Unsave(ctx context.Context, userID, postID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID, "post_id": postID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *savedPostRepo) IsSaved(ctx context.Context, userID, postID primitive.ObjectID) (bool, error) {
	filter := bson.M{"user_id": userID, "post_id": postID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *savedPostRepo) GetByUserID(ctx context.Context, userID primitive.ObjectID, findOptions *FindOptions) ([]*model.SavedPost, int64, error) {
	filter := bson.M{"user_id": userID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find()
	if findOptions != nil {
		if findOptions.Sort != nil {
			opts.SetSort(findOptions.Sort)
		}
		if findOptions.Skip != 0 {
			opts.SetSkip(findOptions.Skip)
		}
		if findOptions.Limit != 0 {
			opts.SetLimit(findOptions.Limit)
		}
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var posts []*model.SavedPost
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
