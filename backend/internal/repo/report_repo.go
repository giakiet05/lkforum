package repo

import (
	"context"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReportRepo interface {
	Create(ctx context.Context, report *model.Report) error
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
