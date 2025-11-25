package repo

import (
	"context"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HiddenPostRepo interface {
	Hide(ctx context.Context, userID, postID primitive.ObjectID) error
	GetHiddenPostIDs(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error)
}

type hiddenPostRepo struct {
	collection *mongo.Collection
}

func NewHiddenPostRepo(db *mongo.Database) HiddenPostRepo {
	return &hiddenPostRepo{
		collection: db.Collection(config.HiddenPostColName),
	}
}

func (r *hiddenPostRepo) Hide(ctx context.Context, userID, postID primitive.ObjectID) error {
	filter := bson.M{"user_id": userID, "post_id": postID}
	update := bson.M{
		"$setOnInsert": bson.M{
			"user_id":   userID,
			"post_id":   postID,
			"hidden_at": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func (r *hiddenPostRepo) GetHiddenPostIDs(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	filter := bson.M{"user_id": userID}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetProjection(bson.M{"post_id": 1}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		PostID primitive.ObjectID `bson:"post_id"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	ids := make([]primitive.ObjectID, len(results))
	for i, result := range results {
		ids[i] = result.PostID
	}

	return ids, nil
}
