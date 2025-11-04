package repo

import (
	"context"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChannelRepo interface {
	Create(ctx context.Context, channel *model.Channel) (*model.Channel, error)
	GetByID(ctx context.Context, channelID string) (*model.Channel, error)
	GetByUserID(ctx context.Context, userID string, page int, pageSize int) ([]model.Channel, int64, error)
	GetByBothUserID(ctx context.Context, user1ID string, user2ID string) (*model.Channel, error)
	Update(ctx context.Context, channel *model.Channel) (*model.Channel, error)
	UpdateUserAvatar(ctx context.Context, userID string, newAvatar string) error
	Delete(ctx context.Context, channelID string, userID string) error
	IsMember(ctx context.Context, channelID string, userID string) (bool, error)
}

type channelRepo struct {
	channelCollection *mongo.Collection
}

func NewChannelRepo(db *mongo.Database) ChannelRepo {
	return &channelRepo{channelCollection: db.Collection(config.ChannelColName)}
}

func (c *channelRepo) Create(ctx context.Context, channel *model.Channel) (*model.Channel, error) {
	result, err := c.channelCollection.InsertOne(ctx, channel)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		channel.ID = oid
	}

	return channel, nil
}

func (c *channelRepo) GetByID(ctx context.Context, channelID string) (*model.Channel, error) {
	channelObjectID, err := primitive.ObjectIDFromHex(channelID)
	if err != nil {
		return nil, err
	}

	var channel model.Channel
	err = c.channelCollection.FindOne(ctx, bson.M{"_id": channelObjectID}).Decode(&channel)
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (c *channelRepo) GetByUserID(ctx context.Context, userID string, page int, pageSize int) ([]model.Channel, int64, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	filter := bson.M{
		"members.user_id": userObjectID,
		"settings": bson.M{
			"$elemMatch": bson.M{
				"user_id":    userObjectID,
				"is_deleted": false,
			},
		},
	}
	opt := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize))

	cursor, err := c.channelCollection.Find(ctx, filter, opt)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var channels []model.Channel
	if err = cursor.All(ctx, &channels); err != nil {
		return nil, 0, err
	}

	count, err := c.channelCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return channels, count, nil
}

func (c *channelRepo) GetByBothUserID(ctx context.Context, user1ID string, user2ID string) (*model.Channel, error) {
	user1ObjectID, err := primitive.ObjectIDFromHex(user1ID)
	if err != nil {
		return nil, err
	}

	user2ObjectID, err := primitive.ObjectIDFromHex(user2ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"members.user_id": bson.M{
			"$all": bson.A{user1ObjectID, user2ObjectID},
		},
		"settings": bson.M{
			"$elemMatch": bson.M{
				"user_id":    user1ObjectID,
				"is_deleted": false,
			},
		},
	}

	var channel model.Channel
	err = c.channelCollection.FindOne(ctx, filter).Decode(&channel)
	if err != nil {
		return nil, err
	}

	return &channel, nil
}

func (c *channelRepo) Update(ctx context.Context, channel *model.Channel) (*model.Channel, error) {
	update := bson.M{
		"$set": bson.M{
			"settings":   channel.Settings,
			"status":     channel.Status,
			"updated_at": time.Now(),
		},
	}

	_, err := c.channelCollection.UpdateByID(ctx, channel.ID, update)
	if err != nil {
		return nil, err
	}

	var updated model.Channel
	err = c.channelCollection.FindOne(ctx, bson.M{"_id": channel.ID}).Decode(&updated)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (c *channelRepo) UpdateUserAvatar(ctx context.Context, userID string, newAvatar string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"members.user_id": userObjectID,
	}

	update := bson.M{
		"$set": bson.M{
			"members.$[elem].avatar": newAvatar,
			"updated_at":             time.Now(),
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.user_id": userObjectID},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters).SetUpsert(false)

	_, err = c.channelCollection.UpdateMany(ctx, filter, update, opts)
	return err
}

func (c *channelRepo) Delete(ctx context.Context, channelID string, userID string) error {
	channelObjectID, err := primitive.ObjectIDFromHex(channelID)
	if err != nil {
		return err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"settings.$[elem].is_deleted": true,
			"updated_at":                  time.Now(),
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"elem.user_id": userObjectID}},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)

	_, err = c.channelCollection.UpdateOne(ctx, bson.M{"_id": channelObjectID}, update, opts)
	return err
}

func (c *channelRepo) IsMember(ctx context.Context, channelID string, userID string) (bool, error) {
	channel, err := c.GetByID(ctx, channelID)
	if err != nil {
		return false, err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	for _, m := range channel.Members {
		if m.UserID == userObjectID {
			return true, nil
		}
	}

	return false, nil
}
