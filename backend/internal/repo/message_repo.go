package repo

import (
	"context"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepo interface {
	Create(ctx context.Context, message *model.Message) (*model.Message, error)
	GetByID(ctx context.Context, messageID string) (*model.Message, error)
	GetFilter(ctx context.Context,
		channelID string, senderID *string, searchContent *string,
		isRead *bool, isSent *bool, isMedia *bool,
		page int, pageSize int,
	) ([]model.Message, int64, error)
	Delete(ctx context.Context, messageID string) error
	IsSendByUser(ctx context.Context, messageID string, userID string) (bool, error)
	MarkChannelAsRead(ctx context.Context, channelID string, userID string) error
}

type messageRepo struct {
	messageCollection *mongo.Collection
}

func NewMessageRepo(db *mongo.Database) MessageRepo {
	return &messageRepo{messageCollection: db.Collection(config.MessageColName)}
}

func (m *messageRepo) Create(ctx context.Context, message *model.Message) (*model.Message, error) {
	result, err := m.messageCollection.InsertOne(ctx, message)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		message.ID = oid
	}

	return message, nil
}

func (m *messageRepo) GetByID(ctx context.Context, messageID string) (*model.Message, error) {
	messageObjectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return nil, err
	}

	var message model.Message
	err = m.messageCollection.FindOne(ctx, bson.M{"_id": messageObjectID}).Decode(&message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func (m *messageRepo) GetFilter(
	ctx context.Context,
	channelID string, senderID *string,
	searchContent *string,
	isRead *bool, isSent *bool, isMedia *bool,
	page int, pageSize int,
) ([]model.Message, int64, error) {
	channelObjectID, err := primitive.ObjectIDFromHex(channelID)
	if err != nil {
		return nil, 0, err
	}

	var senderObjectID *primitive.ObjectID
	if senderID != nil && *senderID != "" {
		oid, err := primitive.ObjectIDFromHex(*senderID)
		if err != nil {
			return nil, 0, err
		}
		senderObjectID = &oid
	}

	// Base filter
	filter := bson.M{
		"channel_id": channelObjectID,
		"is_deleted": false,
	}

	// Conditionally include filters
	if senderObjectID != nil {
		filter["sender_id"] = senderObjectID
	}

	if searchContent != nil && *searchContent != "" {
		filter["content"] = bson.M{"$regex": searchContent, "$options": "i"}
	}

	if isRead != nil {
		filter["is_read"] = *isRead
	}

	if isSent != nil {
		filter["is_send"] = *isSent
	}

	if isMedia != nil && *isMedia {
		filter["content"] = bson.M{"$regex": `https?://|\.jpg|\.png|\.mp4`, "$options": "i"}
	}

	skip := (page - 1) * pageSize

	// Count total matching documents
	total, err := m.messageCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: 1}}).
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize))

	cursor, err := m.messageCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var messages []model.Message
	if err := cursor.All(ctx, &messages); err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (m *messageRepo) Delete(ctx context.Context, messageID string) error {
	messageObjectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"content":    "[deleted]",
			"deleted_at": time.Now(),
		},
	}

	result, err := m.messageCollection.UpdateOne(ctx, bson.M{"_id": messageObjectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return apperror.ErrCommentNotFound
	}

	return nil
}

func (m *messageRepo) IsSendByUser(ctx context.Context, messageID string, userID string) (bool, error) {
	messageObjectID, err := primitive.ObjectIDFromHex(messageID)
	if err != nil {
		return false, err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	filter := bson.M{
		"_id":       messageObjectID,
		"sender_id": userObjectID,
	}

	count, err := m.messageCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *messageRepo) MarkChannelAsRead(ctx context.Context, channelID string, userID string) error {
	channelObjectID, err := primitive.ObjectIDFromHex(channelID)
	if err != nil {
		return err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	// Mark all messages in this channel as read, except those sent by the user
	filter := bson.M{
		"channel_id": channelObjectID,
		"sender_id":  bson.M{"$ne": userObjectID},
		"is_read":    false,
	}

	update := bson.M{
		"$set": bson.M{
			"is_read": true,
		},
	}

	_, err = m.messageCollection.UpdateMany(ctx, filter, update)
	return err
}
