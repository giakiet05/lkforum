package repo

import (
	"context"
	"fmt"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostHistoryRepo interface {
	Create(ctx context.Context, postHistory *model.PostHistory) (*model.PostHistory, error)
	CreateBatch(ctx context.Context, histories []*model.PostHistory) ([]*model.PostHistory, error)
	GetByID(ctx context.Context, id string) (*model.PostHistory, error)
	GetByUserID(ctx context.Context, userID string, page int, pageSize int) ([]*model.PostHistory, int64, error)
	DeleteByID(ctx context.Context, id string) error
}

type postHistoryRepo struct {
	postHistoryCollection *mongo.Collection
}

func NewPostHistoryRepo(db *mongo.Database) PostHistoryRepo {
	return &postHistoryRepo{postHistoryCollection: db.Collection(config.UserPostHistoryColName)}
}

func (h *postHistoryRepo) Create(ctx context.Context, postHistory *model.PostHistory) (*model.PostHistory, error) {
	result, err := h.postHistoryCollection.InsertOne(ctx, postHistory)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		postHistory.ID = oid
	}

	return postHistory, nil
}

func (h *postHistoryRepo) CreateBatch(ctx context.Context, histories []*model.PostHistory) ([]*model.PostHistory, error) {
	docs := make([]interface{}, len(histories))
	for i, history := range histories {
		docs[i] = history
	}

	result, err := h.postHistoryCollection.InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}

	// Assign inserted IDs back to the objects
	for i, id := range result.InsertedIDs {
		if oid, ok := id.(primitive.ObjectID); ok {
			histories[i].ID = oid
		}
	}

	return histories, nil
}

func (h *postHistoryRepo) GetByID(ctx context.Context, id string) (*model.PostHistory, error) {
	postHistoryObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var postHistory model.PostHistory
	err = h.postHistoryCollection.FindOne(ctx, bson.M{"_id": postHistoryObjectID}).Decode(&postHistory)
	if err != nil {
		return nil, err
	}

	return &postHistory, nil
}

func (h *postHistoryRepo) GetByUserID(ctx context.Context, userID string, page int, pageSize int) ([]*model.PostHistory, int64, error) {
	filter := bson.M{"user_id": userID}
	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.M{"createAt": -1})

	cursor, err := h.postHistoryCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var postHistorys []*model.PostHistory
	err = cursor.All(ctx, &postHistorys)
	if err != nil {
		return nil, 0, err
	}

	total, err := h.postHistoryCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return postHistorys, total, nil
}

func (h *postHistoryRepo) DeleteByID(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	res, err := h.postHistoryCollection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("no community found with id %v", id)
	}

	return nil
}
