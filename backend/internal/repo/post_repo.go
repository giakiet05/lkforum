package repo

import (
	"context"
	"errors"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PostRepo defines the data access layer for posts, independent of the database implementation.
type PostRepo interface {
	Create(ctx context.Context, post *model.Post) (*model.Post, error)
	GetByID(ctx context.Context, id string) (*model.Post, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.Post, error)
	Find(ctx context.Context, filter Filter, opts *FindOptions) ([]*model.Post, int64, error)
	UpdateByID(ctx context.Context, id string, update UpdateDocument) error
	Update(ctx context.Context, filter Filter, update UpdateDocument) error
	Delete(ctx context.Context, id string) error
	SoftDelete(ctx context.Context, id string) error
	Increment(ctx context.Context, id string, field string, value int) error
}

type postRepo struct {
	collection *mongo.Collection
}

// NewPostRepo creates a new instance of PostRepo.
func NewPostRepo(db *mongo.Database) PostRepo {
	return &postRepo{collection: db.Collection(config.PostColName)}
}

func (r *postRepo) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	result, err := r.collection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}
	post.ID = result.InsertedID.(primitive.ObjectID)
	return post, nil
}

func (r *postRepo) GetByID(ctx context.Context, id string) (*model.Post, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	var post model.Post
	filter := bson.M{"_id": objectID, "is_deleted": bson.M{"$ne": true}}

	if err := r.collection.FindOne(ctx, filter).Decode(&post); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperror.ErrPostNotFound
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepo) GetByIDs(ctx context.Context, ids []string) ([]*model.Post, error) {
	if len(ids) == 0 {
		return []*model.Post{}, nil
	}

	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if objID, err := primitive.ObjectIDFromHex(id); err == nil {
			objIDs = append(objIDs, objID)
		}
	}

	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*model.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}


// Find fetches paginated data and total count using two separate queries for simplicity and robustness.
func (r *postRepo) Find(ctx context.Context, filter Filter, opts *FindOptions) ([]*model.Post, int64, error) {
	// 1. Get total count using an aggregation pipeline to avoid potential bugs in CountDocuments.
	countPipeline := mongo.Pipeline{
		{{"$match", bson.M(filter)}},
		{{"$count", "total"}},
	}
	cursor, err := r.collection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}

	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err = cursor.All(ctx, &countResult); err != nil {
		return nil, 0, err
	}

	var total int64
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// If there are no documents, return early
	if total == 0 {
		return []*model.Post{}, 0, nil
	}

	// 2. Get the paginated data
	findOptions := options.Find()
	if opts != nil {
		if opts.Sort != nil {
			sortDoc := bson.D{}
			for key, value := range opts.Sort {
				sortDoc = append(sortDoc, bson.E{Key: key, Value: value})
			}
			findOptions.SetSort(sortDoc)
		}
		if opts.Skip > 0 {
			findOptions.SetSkip(opts.Skip)
		}
		if opts.Limit > 0 {
			findOptions.SetLimit(opts.Limit)
		}
	}

	cursor, err = r.collection.Find(ctx, bson.M(filter), findOptions)
	if err != nil {
		// If the find operation fails, we already have the count, but we should return the error.
		return nil, total, err
	}
	defer cursor.Close(ctx)

	var posts []*model.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
func (r *postRepo) UpdateByID(ctx context.Context, id string, update UpdateDocument) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrInvalidID
	}
	filter := bson.M{"_id": objectID}
	result, err := r.collection.UpdateOne(ctx, filter, bson.M(update))
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return apperror.ErrPostNotFound
	}
	return nil
}

func (r *postRepo) Update(ctx context.Context, filter Filter, update UpdateDocument) error {
	_, err := r.collection.UpdateMany(ctx, bson.M(filter), bson.M(update))
	return err
}

func (r *postRepo) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrInvalidID
	}
	filter := bson.M{"_id": objectID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return apperror.ErrPostNotFound
	}
	return nil
}

func (r *postRepo) SoftDelete(ctx context.Context, id string) error {
	update := UpdateDocument{"$set": bson.M{"is_deleted": true}}
	return r.UpdateByID(ctx, id, update)
}

func (r *postRepo) Increment(ctx context.Context, id string, field string, value int) error {
	update := UpdateDocument{"$inc": bson.M{field: value}}
	return r.UpdateByID(ctx, id, update)
}
