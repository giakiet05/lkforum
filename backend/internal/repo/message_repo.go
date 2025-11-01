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
)

type MessageRepo interface {
	Create(ctx context.Context, message *model.Message) (*model.Message, error)
	GetByID(ctx context.Context, messageID string) (*model.Message, error)
	GetFilter(ctx context.Context,
		channelID string, senderID string, searchContent string,
		isRead bool, isSent bool, isMedia bool,
		page int, pageSize int,
	) ([]model.Message, int64, error)
	Delete(ctx context.Context, messageID string) error
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
	channelID string, senderID string, searchContent string,
	isRead bool, isSent bool, isMedia bool,
	page int, pageSize int,
) ([]model.Message, int64, error) {
	//TODO implement me
	panic("implement me")
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
