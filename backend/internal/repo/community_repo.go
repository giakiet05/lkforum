package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommunityRepo interface {
	Create(ctx context.Context, community *model.Community) (*model.Community, error)
	GetAll(ctx context.Context) ([]*model.Community, error)
	GetById(ctx context.Context, id primitive.ObjectID) (*model.Community, error)
	GetByModeratorId(ctx context.Context, moderatorId primitive.ObjectID) ([]*model.Community, error)
	UpdateOne(ctx context.Context, communityID primitive.ObjectID, updates bson.M) (*model.Community, error)
	Delete(ctx context.Context, community *model.Community) error
}

type communityRepo struct {
	communityCollection *mongo.Collection
}

func NewCommunityRepo(db *mongo.Database) CommunityRepo {
	return &communityRepo{communityCollection: db.Collection(config.CommunityColName)}
}

func (c *communityRepo) Create(ctx context.Context, community *model.Community) (*model.Community, error) {
	result, err := c.communityCollection.InsertOne(ctx, community)
	if err != nil {
		return nil, fmt.Errorf("failed to create community: %w", err)
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		community.ID = oid
	}

	return community, nil
}

func (c *communityRepo) GetAll(ctx context.Context) ([]*model.Community, error) {
	cursor, err := c.communityCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, fmt.Errorf("failed to query communities: %w", err)
	}
	defer cursor.Close(ctx)

	var communities []*model.Community
	if err := cursor.All(ctx, &communities); err != nil {
		return nil, fmt.Errorf("failed to decode communities: %w", err)
	}

	return communities, nil
}

func (c *communityRepo) GetByModeratorId(ctx context.Context, moderatorId primitive.ObjectID) ([]*model.Community, error) {
	filter := bson.M{"moderators.user_id": moderatorId}

	cursor, err := c.communityCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to query communities: %w", err)
	}
	defer cursor.Close(ctx)

	var communities []*model.Community
	if err := cursor.All(ctx, &communities); err != nil {
		return nil, fmt.Errorf("failed to decode communities: %w", err)
	}

	return communities, nil
}

func (c *communityRepo) GetById(ctx context.Context, id primitive.ObjectID) (*model.Community, error) {
	var community model.Community
	err := c.communityCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&community)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("failed to query community: %w", err)
	}

	return &community, nil
}

func (c *communityRepo) UpdateOne(ctx context.Context, communityID primitive.ObjectID, updates bson.M) (*model.Community, error) {
	filter := bson.M{"_id": communityID}
	update := bson.M{"$set": updates}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated model.Community
	err := c.communityCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update community: %w", err)
	}

	return &updated, nil
}

func (c *communityRepo) Delete(ctx context.Context, community *model.Community) error {
	_, err := c.communityCollection.DeleteOne(ctx, bson.M{"_id": community.ID})
	if err != nil {
		return fmt.Errorf("failed to delete community: %w", err)
	}
	return nil
}
