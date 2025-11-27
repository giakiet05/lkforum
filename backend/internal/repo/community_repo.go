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

type CommunityRepo interface {
	Create(ctx context.Context, community *model.Community) (*model.Community, error)
	GetByID(ctx context.Context, id string) (*model.Community, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.Community, error)
	GetFilter(ctx context.Context, name string, description string, is18Plus bool, createFrom time.Time, page int, pageSize int) ([]model.Community, int64, error)
	GetByModeratorIDPaginated(ctx context.Context, moderatorID string, page int, pageSize int) ([]model.Community, int64, error)
	GetAllPaginated(ctx context.Context, page int, pageSize int) ([]model.Community, int64, error)
	Update(ctx context.Context, communityID string, updates bson.M) (*model.Community, error)
	UpdateUserAvatar(ctx context.Context, userID string, newAvatar string) error
	Replace(ctx context.Context, community *model.Community) error
	Delete(ctx context.Context, communityID string) error
	IsUserExist(ctx context.Context, userID string) (bool, error)

	GetBannedUsers(ctx context.Context, communityID string) ([]*model.User, error)
	GetMutedUsers(ctx context.Context, communityID string) ([]*model.User, error)
	BanUser(ctx context.Context, ban *model.CommunityBan) error
	UnmuteUser(ctx context.Context, userID string, communityID string) error
	UnbanUser(ctx context.Context, userID string, communityID string) error
}

type communityRepo struct {
	communityCollection    *mongo.Collection
	userCollection         *mongo.Collection
	communityBanCollection *mongo.Collection
}

func NewCommunityRepo(db *mongo.Database) CommunityRepo {
	return &communityRepo{
		communityCollection:    db.Collection(config.CommunityColName),
		userCollection:         db.Collection(config.UserColName),
		communityBanCollection: db.Collection(config.CommunityBanColName),
	}
}

func (c *communityRepo) GetByIDs(ctx context.Context, ids []string) ([]*model.Community, error) {
	if len(ids) == 0 {
		return []*model.Community{}, nil
	}

	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if objID, err := primitive.ObjectIDFromHex(id); err == nil {
			objIDs = append(objIDs, objID)
		}
	}

	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	cursor, err := c.communityCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var communities []*model.Community
	if err := cursor.All(ctx, &communities); err != nil {
		return nil, err
	}
	return communities, nil
}

func (c *communityRepo) Create(ctx context.Context, community *model.Community) (*model.Community, error) {
	result, err := c.communityCollection.InsertOne(ctx, community)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		community.ID = oid
	}

	return community, nil
}

func (c *communityRepo) GetByID(ctx context.Context, id string) (*model.Community, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	var community model.Community
	err = c.communityCollection.FindOne(ctx, bson.M{"_id": communityObjectID}).Decode(&community)
	if err != nil {
		return nil, err
	}

	return &community, nil
}

func (c *communityRepo) GetFilter(
	ctx context.Context,
	name string,
	description string,
	is18Plus bool,
	createFrom time.Time,
	page int,
	pageSize int,
) ([]model.Community, int64, error) {
	filter := bson.M{}
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if description != "" {
		filter["description"] = bson.M{"$regex": description, "$options": "i"}
	}
	if is18Plus {
		filter["is18_plus"] = true
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
	opt := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize))

	cursor, err := c.communityCollection.Find(ctx, filter, opt)
	if err != nil {
		return nil, -1, err
	}
	defer cursor.Close(ctx)

	var communities []model.Community
	if err := cursor.All(ctx, &communities); err != nil {
		return nil, -1, err
	}

	count, err := c.communityCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, -1, err
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
		return nil, -1, err
	}
	defer cursor.Close(ctx)

	var communities []model.Community
	if err := cursor.All(ctx, &communities); err != nil {
		return nil, -1, err
	}

	count, err := c.communityCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, -1, err
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
		return nil, err
	}

	return &updated, nil
}

func (c *communityRepo) UpdateUserAvatar(ctx context.Context, userID string, newAvatar string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	// Update moderator avatars
	filterModerators := bson.M{
		"moderators.user_id": userObjectID,
	}

	updateModerators := bson.M{
		"$set": bson.M{
			"moderators.$[elem].avatar": newAvatar,
			"updated_at":                time.Now(),
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"elem.user_id": userObjectID},
		},
	}

	opts := options.Update().SetArrayFilters(arrayFilters)

	if _, err := c.communityCollection.UpdateMany(ctx, filterModerators, updateModerators, opts); err != nil {
		return err
	}

	// Update communities creator avatar
	filterCreator := bson.M{
		"create_by_id": userObjectID,
	}

	updateCreator := bson.M{
		"$set": bson.M{
			"create_by_avatar": newAvatar,
			"updated_at":       time.Now(),
		},
	}

	if _, err := c.communityCollection.UpdateMany(ctx, filterCreator, updateCreator); err != nil {
		return err
	}

	return nil
}

func (c *communityRepo) Replace(ctx context.Context, community *model.Community) error {
	res, err := c.communityCollection.ReplaceOne(ctx, bson.M{"_id": community.ID}, community)
	if err != nil {
		return err
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
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("no community found with id %v", communityID)
	}

	return nil
}

func (c *communityRepo) IsUserExist(ctx context.Context, userID string) (bool, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	count, err := c.userCollection.CountDocuments(ctx, bson.M{"_id": userObjectID})
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (c *communityRepo) GetBannedUsers(ctx context.Context, communityID string) ([]*model.User, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return nil, err
	}

	cursor, err := c.communityBanCollection.Find(ctx, bson.M{
		"community_id": communityObjectID,
		"type":         model.Banned,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bans []model.CommunityBan
	if err := cursor.All(ctx, &bans); err != nil {
		return nil, err
	}

	if len(bans) == 0 {
		return []*model.User{}, nil
	}

	// 2. Extract user IDs
	userIDs := make([]primitive.ObjectID, 0, len(bans))
	for _, b := range bans {
		userIDs = append(userIDs, b.UserID)
	}

	// 3. Query users
	userCursor, err := c.userCollection.Find(ctx, bson.M{
		"_id": bson.M{"$in": userIDs},
	})
	if err != nil {
		return nil, err
	}
	defer userCursor.Close(ctx)

	var users []*model.User
	if err := userCursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (c *communityRepo) GetMutedUsers(ctx context.Context, communityID string) ([]*model.User, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return nil, err
	}

	// 1. Find mute records
	cursor, err := c.communityBanCollection.Find(ctx, bson.M{
		"community_id": communityObjectID,
		"type":         model.Muted,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bans []model.CommunityBan
	if err := cursor.All(ctx, &bans); err != nil {
		return nil, err
	}

	if len(bans) == 0 {
		return []*model.User{}, nil
	}

	// 2. Extract user IDs
	userIDs := make([]primitive.ObjectID, 0, len(bans))
	for _, b := range bans {
		userIDs = append(userIDs, b.UserID)
	}

	// 3. Query users
	userCursor, err := c.userCollection.Find(ctx, bson.M{
		"_id": bson.M{"$in": userIDs},
	})
	if err != nil {
		return nil, err
	}
	defer userCursor.Close(ctx)

	var users []*model.User
	if err := userCursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (c *communityRepo) BanUser(ctx context.Context, ban *model.CommunityBan) error {
	_, err := c.communityBanCollection.DeleteMany(ctx, bson.M{
		"user_id":      ban.UserID,
		"community_id": ban.CommunityID,
	})
	if err != nil {
		return err
	}

	_, err = c.communityBanCollection.InsertOne(ctx, ban)
	return err
}

func (c *communityRepo) UnmuteUser(ctx context.Context, userID string, communityID string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return err
	}

	// Remove ONLY the mute record
	_, err = c.communityBanCollection.DeleteOne(ctx, bson.M{
		"user_id":      userObjectID,
		"community_id": communityObjectID,
		"type":         model.Muted,
	})

	return err
}

func (c *communityRepo) UnbanUser(ctx context.Context, userID string, communityID string) error {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return err
	}

	// Remove ONLY the ban record
	_, err = c.communityBanCollection.DeleteOne(ctx, bson.M{
		"user_id":      userObjectID,
		"community_id": communityObjectID,
		"type":         model.Banned,
	})

	return err
}
