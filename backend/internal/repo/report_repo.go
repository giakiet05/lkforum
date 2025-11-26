package repo

import (
	"context"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportRepo interface {
	Create(ctx context.Context, report *model.Report) error
	GetByReporterAndTarget(ctx context.Context, reporterID, targetID primitive.ObjectID, targetType model.ReportTargetType) (bool, error)
}

type reportRepo struct {
	collection *mongo.Collection
}

func NewReportRepo(db *mongo.Database) ReportRepo {
	return &reportRepo{
		collection: db.Collection(config.ReportColName),
	}
}

func (r *reportRepo) Create(ctx context.Context, report *model.Report) error {
	_, err := r.collection.InsertOne(ctx, report)
	return err
}

func (r *reportRepo) GetByReporterAndTarget(ctx context.Context, reporterID, targetID primitive.ObjectID, targetType model.ReportTargetType) (bool, error) {
	filter := bson.M{
		"reporter_id": reporterID,
		"target_id":   targetID,
		"target_type": targetType,
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
