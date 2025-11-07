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
)

var ErrPostNotFound = errors.New("post not found")

// PostRepo defines the data access layer for posts, independent of the database implementation.
type PostRepo interface {
	Create(ctx context.Context, post *model.Post) (*model.Post, error)
	GetByID(ctx context.Context, id string) (*model.Post, error)
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
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	return &post, nil
}

// Find uses an aggregation pipeline to fetch paginated data and total count in a single query.
func (r *postRepo) Find(ctx context.Context, filter Filter, opts *FindOptions) ([]*model.Post, int64, error) {
	pipeline := mongo.Pipeline{}

	// Match stage for filtering
	pipeline = append(pipeline, bson.D{{Key: "$match", Value: bson.M(filter)}})

	// Facet stage for getting both metadata (total count) and paginated data
	facetStage := bson.D{{
		Key: "$facet",
		Value: bson.D{
			{"metadata", bson.A{bson.D{{Key: "$count", Value: "total"}}}},
			{"data", bson.A{}},
		},
	}}

	// Build the data pipeline inside the facet
	dataPipeline := bson.A{}
	if opts != nil {
		if opts.Sort != nil {
			sortDoc := bson.D{}
			for key, value := range opts.Sort {
				sortDoc = append(sortDoc, bson.E{Key: key, Value: value})
			}
			dataPipeline = append(dataPipeline, bson.D{{Key: "$sort", Value: sortDoc}})
		}
		if opts.Skip > 0 {
			dataPipeline = append(dataPipeline, bson.D{{Key: "$skip", Value: opts.Skip}})
		}
		if opts.Limit > 0 {
			dataPipeline = append(dataPipeline, bson.D{{Key: "$limit", Value: opts.Limit}})
		}
	}
	facetStage[0].Value.(bson.D)[1].Value.(bson.A)[0] = dataPipeline[0]

	pipeline = append(pipeline, facetStage)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var result []struct {
		Metadata []struct {
			Total int64 `bson:"total"`
		} `bson:"metadata"`
		Data []*model.Post `bson:"data"`
	}

	if !cursor.Next(ctx) {
		// No documents found, return empty slice and 0 count
		return []*model.Post{}, 0, nil
	}
	if err := cursor.Decode(&result); err != nil {
		return nil, 0, err
	}

	if len(result) == 0 || len(result[0].Metadata) == 0 {
		return []*model.Post{}, 0, nil
	}

	return result[0].Data, result[0].Metadata[0].Total, nil
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
		return ErrPostNotFound
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
		return ErrPostNotFound
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
