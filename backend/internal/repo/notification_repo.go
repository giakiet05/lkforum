package repo

import (
	"context"
	"log"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NotificationRepo interface {
	Create(ctx context.Context, notification *model.Notification) (*model.Notification, error)
	GetByRecipientID(ctx context.Context, recipientID string, page, pageSize int) ([]*model.Notification, int64, error)
	MarkAsRead(ctx context.Context, notificationID primitive.ObjectID, recipientID string) error
	MarkAllAsRead(ctx context.Context, recipientID string) (int64, error)
	DeleteNotification(ctx context.Context, notificationID primitive.ObjectID, recipientID string) error
	CountUnread(ctx context.Context, recipientID string) (int64, error)
}

type notificationRepo struct {
	notificationCollection *mongo.Collection
}

func NewNotificationRepo(db *mongo.Database) NotificationRepo {
	return &notificationRepo{notificationCollection: db.Collection(config.NotificationColName)}
}

func (r *notificationRepo) Create(ctx context.Context, notification *model.Notification) (*model.Notification, error) {
	result, err := r.notificationCollection.InsertOne(ctx, notification)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		notification.ID = oid
	}

	return notification, nil
}

func (r *notificationRepo) GetByRecipientID(ctx context.Context, recipientID string, page, pageSize int) ([]*model.Notification, int64, error) {
	recipientObjID, err := primitive.ObjectIDFromHex(recipientID)
	if err != nil {
		return nil, 0, err
	}

	filter := bson.M{"recipient_id": recipientObjID}
	skip := (page - 1) * pageSize

	findOptions := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := r.notificationCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var notifications []*model.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, 0, err
	}

	// Debug log
	unreadCount := 0
	for _, n := range notifications {
		if !n.IsRead {
			unreadCount++
		}
	}
	log.Printf("📋 [GetByRecipientID] Found %d notifications, %d unread for user %s", len(notifications), unreadCount, recipientID)

	total, err := r.notificationCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *notificationRepo) MarkAsRead(ctx context.Context, notificationID primitive.ObjectID, recipientID string) error {
	recipientObjID, err := primitive.ObjectIDFromHex(recipientID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":          notificationID,
		"recipient_id": recipientObjID,
	}
	update := bson.M{
		"$set": bson.M{"is_read": true},
	}

	result, err := r.notificationCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments // Or a custom error like ErrNotificationNotFound
	}

	return nil
}

func (r *notificationRepo) MarkAllAsRead(ctx context.Context, recipientID string) (int64, error) {
	recipientObjID, err := primitive.ObjectIDFromHex(recipientID)
	if err != nil {
		log.Printf("❌ [MarkAllAsRead] Invalid recipientID: %s, error: %v", recipientID, err)
		return 0, err
	}

	// Match notifications where is_read is false OR is_read doesn't exist (nil)
	filter := bson.M{
		"recipient_id": recipientObjID,
		"$or": []bson.M{
			{"is_read": false},
			{"is_read": bson.M{"$exists": false}},
			{"is_read": nil},
		},
	}
	update := bson.M{
		"$set": bson.M{"is_read": true},
	}

	log.Printf("🔔 [MarkAllAsRead] Updating notifications for user %s", recipientID)

	result, err := r.notificationCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		log.Printf("❌ [MarkAllAsRead] UpdateMany error: %v", err)
		return 0, err
	}

	log.Printf("✅ [MarkAllAsRead] MatchedCount: %d, ModifiedCount: %d", result.MatchedCount, result.ModifiedCount)
	return result.ModifiedCount, nil
}

func (r *notificationRepo) CountUnread(ctx context.Context, recipientID string) (int64, error) {
	recipientObjID, err := primitive.ObjectIDFromHex(recipientID)
	if err != nil {
		return 0, err
	}

	filter := bson.M{
		"recipient_id": recipientObjID,
		"is_read":      false,
	}

	return r.notificationCollection.CountDocuments(ctx, filter)
}

func (r *notificationRepo) DeleteNotification(ctx context.Context, notificationID primitive.ObjectID, recipientID string) error {
	recipientObjID, err := primitive.ObjectIDFromHex(recipientID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":          notificationID,
		"recipient_id": recipientObjID,
	}

	result, err := r.notificationCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
