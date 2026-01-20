package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReportRepo interface {
	Create(ctx context.Context, report *model.Report) error
	GetByID(ctx context.Context, id string) (*model.Report, error)
	GetFilter(
		ctx context.Context,
		reporterID *string,
		targetID *string,
		targetType model.ReportTargetType,
		reason *string,
		startDate *time.Time,
		endDate *time.Time,
		page int, pageSize int,
	) ([]model.Report, int64, error)
	Delete(ctx context.Context, reportID string) error
	DeleteBatch(ctx context.Context, reportIDs []string) error

	IsReporter(ctx context.Context, userID string, reportID string) (bool, error)
	IsReporterOfAllReports(ctx context.Context, userID string, reportIDs []string) (bool, error)
	
	// Stats methods
	CountPending(ctx context.Context) (int64, error)
	CountCreatedAfter(ctx context.Context, since time.Time) (int64, error)
}

type reportRepo struct {
	reportCollection *mongo.Collection
}

func NewReportRepo(db *mongo.Database) ReportRepo {
	return &reportRepo{
		reportCollection: db.Collection(config.ReportColName),
	}
}

func (r *reportRepo) Create(ctx context.Context, report *model.Report) error {
	_, err := r.reportCollection.InsertOne(ctx, report)
	return err
}

func (r *reportRepo) GetByID(ctx context.Context, id string) (*model.Report, error) {
	reportObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	var report model.Report
	err = r.reportCollection.FindOne(ctx, bson.M{"_id": reportObjectID}).Decode(&report)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (r *reportRepo) GetFilter(
	ctx context.Context,
	reporterID *string,
	targetID *string,
	targetType model.ReportTargetType,
	reason *string,
	startDate *time.Time,
	endDate *time.Time,
	page int, pageSize int,
) ([]model.Report, int64, error) {
	filter := bson.M{}

	if reporterID != nil {
		filter["reporter_id"] = *reporterID
	}

	if targetID != nil {
		filter["target_id"] = *targetID
	}

	if targetType != "" {
		filter["target_type"] = targetType
	}

	if reason != nil {
		filter["reason"] = bson.M{"$regex": *reason, "$options": "i"}
	}

	if startDate != nil || endDate != nil {
		dateFilter := bson.M{}

		if startDate != nil {
			dateFilter["$gte"] = *startDate
		}
		if endDate != nil {
			dateFilter["$lte"] = *endDate
		}

		filter["created_at"] = dateFilter
	}

	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.M{"createAt": -1})

	cursor, err := r.reportCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	var results []model.Report
	if err := cursor.All(ctx, &results); err != nil {
		return nil, 0, err
	}

	total, err := r.reportCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *reportRepo) Delete(ctx context.Context, reportID string) error {
	reportObjectID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		},
	}

	result, err := r.reportCollection.UpdateOne(ctx, bson.M{"_id": reportObjectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return apperror.ErrReportNotFound
	}

	return nil
}

func (r *reportRepo) DeleteBatch(ctx context.Context, reportIDs []string) error {
	objectIDs := make([]primitive.ObjectID, 0, len(reportIDs))
	for _, id := range reportIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return apperror.ErrInvalidID
		}
		objectIDs = append(objectIDs, objID)
	}

	// Soft delete
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		},
	}

	filter := bson.M{"_id": bson.M{"$in": objectIDs}}

	result, err := r.reportCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return apperror.ErrReportNotFound
	}

	return nil
}

func (r *reportRepo) IsReporter(ctx context.Context, userID string, reportID string) (bool, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	reportObjID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		return false, apperror.ErrInvalidID
	}

	filter := bson.M{
		"_id":         reportObjID,
		"reporter_id": userObjID,
	}

	count, err := r.reportCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *reportRepo) IsReporterOfAllReports(ctx context.Context, userID string, reportIDs []string) (bool, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, fmt.Errorf("invalid userID: %w", err)
	}

	reportObjIDs := make([]primitive.ObjectID, 0, len(reportIDs))
	for _, id := range reportIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return false, fmt.Errorf("invalid reportID '%s': %w", id, err)
		}
		reportObjIDs = append(reportObjIDs, objID)
	}

	filter := bson.M{
		"_id":         bson.M{"$in": reportObjIDs},
		"reporter_id": userObjID,
	}

	count, err := r.reportCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count == int64(len(reportIDs)), nil
}

// Stats methods implementations
func (r *reportRepo) CountPending(ctx context.Context) (int64, error) {
filter := bson.M{
"is_deleted": bson.M{"": true},
}
return r.reportCollection.CountDocuments(ctx, filter)
}

func (r *reportRepo) CountCreatedAfter(ctx context.Context, since time.Time) (int64, error) {
filter := bson.M{
"created_at": bson.M{"": since},
}
return r.reportCollection.CountDocuments(ctx, filter)
}
