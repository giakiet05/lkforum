package repo

import (
	"context"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DraftRepo interface {
	Create(ctx context.Context, draft *model.Draft) (*model.Draft, error)
	Update(ctx context.Context, draft *model.Draft) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Draft, error)
	GetByAuthor(ctx context.Context, authorID primitive.ObjectID, opts *FindOptions) ([]*model.Draft, int64, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type draftRepo struct {
	collection *mongo.Collection
}

func NewDraftRepo(db *mongo.Database) DraftRepo {
	return &draftRepo{
		collection: db.Collection(config.DraftColName),
	}
}

func (r *draftRepo) Create(ctx context.Context, draft *model.Draft) (*model.Draft, error) {
	res, err := r.collection.InsertOne(ctx, draft)
	if err != nil {
		return nil, err
	}
	draft.ID = res.InsertedID.(primitive.ObjectID)
	return draft, nil
}

func (r *draftRepo) Update(ctx context.Context, draft *model.Draft) error {
	filter := bson.M{"_id": draft.ID, "author_id": draft.AuthorID}
	update := bson.M{"$set": draft}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return apperror.ErrDraftNotFound
	}
	return nil
}

func (r *draftRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Draft, error) {
	var draft model.Draft
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&draft)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrDraftNotFound
		}
		return nil, err
	}
	return &draft, nil
}

func (r *draftRepo) GetByAuthor(ctx context.Context, authorID primitive.ObjectID, opts *FindOptions) ([]*model.Draft, int64, error) {
	filter := bson.M{"author_id": authorID}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	findOptions := options.Find()
	if opts != nil {
		if opts.Sort != nil {
			findOptions.SetSort(opts.Sort)
		}
		if opts.Skip != 0 {
			findOptions.SetSkip(opts.Skip)
		}
		if opts.Limit != 0 {
			findOptions.SetLimit(opts.Limit)
		}
	}

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var drafts []*model.Draft
	if err = cursor.All(ctx, &drafts); err != nil {
		return nil, 0, err
	}

	return drafts, total, nil
}

func (r *draftRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return apperror.ErrDraftNotFound
	}
	return nil
}
