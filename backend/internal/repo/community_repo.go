package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommunityRepo interface {
	Create(ctx context.Context, community *model.Community) (*model.Community, error)
	//GetAll(ctx context.Context) ([]*model.Community, error)
	GetByID(ctx context.Context, id string) (*model.Community, error)
	GetFilter(ctx context.Context, name string, description string, createFrom time.Time, page int, pageSize int) ([]model.Community, int64, error)
	GetByModeratorIDPaginated(ctx context.Context, moderatorID string, page int, pageSize int) ([]model.Community, int64, error)
	GetAllPaginated(ctx context.Context, page int, pageSize int) ([]model.Community, int64, error)
	Update(ctx context.Context, communityID string, updates bson.M) (*model.Community, error)
	Replace(ctx context.Context, community *model.Community) error
	Delete(ctx context.Context, communityID string) error
}

type communityRepo struct {
	communityCollection *mongo.Collection
}

func NewCommunityRepo(db *mongo.Database) CommunityRepo {
	return &communityRepo{
		communityCollection: db.Collection(config.CommunityColName),
	}
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

/*
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
*/

func (c *communityRepo) GetByID(ctx context.Context, id string) (*model.Community, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var community model.Community
	err = c.communityCollection.FindOne(ctx, bson.M{"_id": communityObjectID}).Decode(&community)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(`community with id %s not found`, communityObjectID)
		}
		return nil, fmt.Errorf("failed to query community: %w", err)
	}

	return &community, nil
}

func (c *communityRepo) GetFilter(
	ctx context.Context,
	name string,
	description string,
	createFrom time.Time,
	page int,
	pageSize int,
) ([]model.Community, int64, error) {
	filter := bson.M{}
	if name != "" {
		// case-insensitive regex match
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if description != "" {
		filter["description"] = bson.M{"$regex": description, "$options": "i"}
	}
	if !createFrom.IsZero() {
		filter["createdAt"] = bson.M{"$gte": createFrom}
	}

	total, err := c.communityCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.M{"createAt": -1})

	cursor, err := c.communityCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var communities []model.Community
	err = cursor.All(ctx, &communities)
	if err != nil {
		return nil, 0, err
	}

	return communities, total, nil
}

func (c *communityRepo) GetByModeratorIDPaginated(
	ctx context.Context,
	moderatorID string,
	page int,
	pageSize int,
) ([]model.Community, int64, error) {
	modObjectID, err := primitive.ObjectIDFromHex(moderatorID)
	if err != nil {
		return nil, -1, err
	}

	skip := (page - 1) * pageSize
	filter := bson.M{"moderators.user_id": modObjectID}

	cursor, err := c.communityCollection.Find(ctx, filter, options.Find().SetSkip(int64(skip)), options.Find().SetLimit(int64(pageSize)))
	if err != nil {
		return nil, -1, fmt.Errorf("failed to query communities: %w", err)
	}
	defer cursor.Close(ctx)

	var communities []model.Community
	if err := cursor.All(ctx, &communities); err != nil {
		return nil, -1, fmt.Errorf("failed to decode communities: %w", err)
	}

	count, err := c.communityCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, -1, fmt.Errorf("failed to query communities: %w", err)
	}

	return communities, count, nil
}

func (c *communityRepo) GetAllPaginated(
	ctx context.Context,
	page int,
	pageSize int,
) ([]model.Community, int64, error) {
	skip := (page - 1) * pageSize
	filter := bson.M{
		"is_deleted": false,
		"is_banned":  false,
	}

	cursor, err := c.communityCollection.Find(ctx, filter, options.Find().SetSkip(int64(skip)), options.Find().SetLimit(int64(pageSize)))
	if err != nil {
		return nil, -1, fmt.Errorf("failed to query communities: %w", err)
	}
	defer cursor.Close(ctx)

	var communities []model.Community
	if err := cursor.All(ctx, &communities); err != nil {
		return nil, -1, fmt.Errorf("failed to decode communities: %w", err)
	}

	count, err := c.communityCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, -1, fmt.Errorf("failed to query communities: %w", err)
	}

	return communities, count, nil
}

func (c *communityRepo) Update(ctx context.Context, communityID string, updates bson.M) (*model.Community, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": communityObjectID}
	update := bson.M{"$set": updates}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated model.Community
	err = c.communityCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		return nil, fmt.Errorf("failed to update community: %w", err)
	}

	return &updated, nil
}

func (c *communityRepo) Replace(ctx context.Context, community *model.Community) error {
	res, err := c.communityCollection.ReplaceOne(ctx, bson.M{"_id": community.ID}, community)
	if err != nil {
		return fmt.Errorf("failed to replace community: %w", err)
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("no document found with id %v", community.ID)
	}

	return nil
}

func (c *communityRepo) Delete(ctx context.Context, communityID string) error {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return err
	}

	res, err := c.communityCollection.DeleteOne(ctx, bson.M{"_id": communityObjectID})
	if err != nil {
		return fmt.Errorf("failed to delete community: %w", err)
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("no community found with id %v", communityID)
	}

	return nil
}

/*
func (c *communityRepo) CountPostsByCommunityID(ctx context.Context, communityID string) (int64, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return -1, err
	}

	filter := bson.M{"community_id": communityObjectID}
	count, err := c.postCollection.CountDocuments(ctx, filter)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (c *communityRepo) CountMembersByCommunityID(ctx context.Context, moderatorID string) (int64, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(moderatorID)
	if err != nil {
		return -1, err
	}

	filter := bson.M{"community_id": communityObjectID}
	count, err := c.membershipCollection.CountDocuments(ctx, filter)
	if err != nil {
		return -1, err
	}

	return count, nil
}
*/
