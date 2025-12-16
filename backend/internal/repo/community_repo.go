package repo

import (
	"context"
	"errors"
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
	GetByName(ctx context.Context, name string) (*model.Community, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.Community, error)
	GetFilter(
		ctx context.Context,
		name string,
		description string,
		is18Plus bool,
		createFrom time.Time,
		page int, pageSize int,
	) ([]*model.Community, int64, error)
	GetByModeratorIDPaginated(ctx context.Context, moderatorID string, page int, pageSize int) ([]*model.Community, int64, error)
	GetAllPaginated(ctx context.Context, page int, pageSize int) ([]*model.Community, int64, error)
	FindCommunities(ctx context.Context, filter Filter, findOptions *FindOptions) ([]*model.Community, int64, error)

	// Stats methods
	CountTotal(ctx context.Context) (int64, error)
	CountActive(ctx context.Context) (int64, error)
	CountCreatedAfter(ctx context.Context, since time.Time) (int64, error)
	CountBanned(ctx context.Context) (int64, error)
	CountPrivate(ctx context.Context) (int64, error)
	Update(ctx context.Context, communityID string, updates bson.M) (*model.Community, error)
	UpdateUserAvatar(ctx context.Context, userID string, newAvatar string) error
	Replace(ctx context.Context, community *model.Community) error
	Delete(ctx context.Context, communityID string) error
	IsUserExist(ctx context.Context, userID string) (bool, error)

	GetBannedUsers(ctx context.Context, communityID string, expired bool) ([]*model.User, error)
	GetMutedUsers(ctx context.Context, communityID string, expired bool) ([]*model.User, error)
	GetBannedCommunityIDs(ctx context.Context, userID string, banType model.CommunityBanType, communityIDs []string) ([]string, error)
	BanUser(ctx context.Context, ban *model.CommunityBan) error
	IsUserBanned(ctx context.Context, userID string, banType model.CommunityBanType, communityID string) (bool, error)
	UnmuteUser(ctx context.Context, userID string, communityID string) error
	UnbanUser(ctx context.Context, userID string, communityID string) error

	ActivateModerator(ctx context.Context, communityID string, userID string) error
	IsModerator(ctx context.Context, communityID string, userID string) (bool, error)
	IsCreator(ctx context.Context, communityID string, userID string) (bool, error)
}

type communityRepo struct {
	communityCollection    *mongo.Collection
	userCollection         *mongo.Collection
	communityBanCollection *mongo.Collection
}

func NewCommunityRepo(db *mongo.Database) CommunityRepo {
	communityCollection := db.Collection(config.CommunityColName)

	// Create unique index on community name
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := communityCollection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		// Log error but don't fail - index might already exist
		fmt.Printf("Warning: Failed to create unique index on community.name: %v\n", err)
	}

	return &communityRepo{
		communityCollection:    communityCollection,
		userCollection:         db.Collection(config.UserColName),
		communityBanCollection: db.Collection(config.CommunityBanColName),
	}
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

func (c *communityRepo) GetByName(ctx context.Context, name string) (*model.Community, error) {
	var community model.Community
	err := c.communityCollection.FindOne(ctx, bson.M{"name": name}).Decode(&community)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, apperror.ErrCommunityNotFound
		}
		return nil, err
	}

	return &community, nil
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

func (c *communityRepo) GetFilter(
	ctx context.Context,
	name string,
	description string,
	is18Plus bool,
	createFrom time.Time,
	page int,
	pageSize int,
) ([]*model.Community, int64, error) {
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

	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.M{"createAt": -1})

	cursor, err := c.communityCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var communities []*model.Community
	err = cursor.All(ctx, &communities)
	if err != nil {
		return nil, 0, err
	}

	total, err := c.communityCollection.CountDocuments(ctx, filter)
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
) ([]*model.Community, int64, error) {
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

	var communities []*model.Community
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
) ([]*model.Community, int64, error) {
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

	var communities []*model.Community
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
		return apperror.ErrCommunityNotFound
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

func (c *communityRepo) GetBannedUsers(ctx context.Context, communityID string, expired bool) ([]*model.User, error) {
	communityOID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return nil, err
	}

	filter := buildBanFilter(communityOID, model.Banned, expired)

	cursor, err := c.communityBanCollection.Find(ctx, filter)
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

	userIDs := make([]primitive.ObjectID, 0, len(bans))
	for _, b := range bans {
		userIDs = append(userIDs, b.UserID)
	}

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

func (c *communityRepo) GetBannedCommunityIDs(
	ctx context.Context,
	userID string,
	banType model.CommunityBanType,
	communityIDs []string,
) ([]string, error) {

	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Convert communityIDs to ObjectIDs
	communityOIDs := make([]primitive.ObjectID, 0, len(communityIDs))
	for _, id := range communityIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		communityOIDs = append(communityOIDs, oid)
	}

	// Query all active bans for these communities
	filter := bson.M{
		"user_id":    userOID,
		"type":       banType,
		"is_deleted": false,
		"expires_at": bson.M{"$gt": time.Now()},
		"community_id": bson.M{
			"$in": communityOIDs,
		},
	}

	cursor, err := c.communityBanCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bannedIDs []string
	for cursor.Next(ctx) {
		var ban model.CommunityBan
		if err := cursor.Decode(&ban); err != nil {
			return nil, err
		}
		bannedIDs = append(bannedIDs, ban.CommunityID.Hex())
	}

	return bannedIDs, nil
}

func (c *communityRepo) GetMutedUsers(ctx context.Context, communityID string, expired bool) ([]*model.User, error) {
	communityOID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return nil, err
	}

	filter := buildBanFilter(communityOID, model.Muted, expired)

	cursor, err := c.communityBanCollection.Find(ctx, filter)
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

	userIDs := make([]primitive.ObjectID, 0, len(bans))
	for _, b := range bans {
		userIDs = append(userIDs, b.UserID)
	}

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

func (c *communityRepo) IsUserBanned(ctx context.Context, userID string, banType model.CommunityBanType, communityID string) (bool, error) {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	communityOID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return false, err
	}

	filter := bson.M{
		"user_id":      userOID,
		"community_id": communityOID,
		"type":         banType,
		"is_deleted":   false,
		"expires_at": bson.M{
			"$gt": time.Now(),
		},
	}

	err = c.communityBanCollection.FindOne(ctx, filter).Err()
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (c *communityRepo) UnmuteUser(ctx context.Context, userID string, communityID string) error {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	communityOID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"user_id":      userOID,
		"community_id": communityOID,
		"type":         model.Muted,
		"is_deleted":   false,
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		},
	}

	_, err = c.communityBanCollection.UpdateOne(ctx, filter, update)
	return err
}

func (c *communityRepo) UnbanUser(ctx context.Context, userID string, communityID string) error {
	userOID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	communityOID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"user_id":      userOID,
		"community_id": communityOID,
		"type":         model.Banned,
		"is_deleted":   false,
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_at": time.Now(),
		},
	}

	_, err = c.communityBanCollection.UpdateOne(ctx, filter, update)
	return err
}

func buildBanFilter(communityID primitive.ObjectID, banType model.CommunityBanType, expired bool) bson.M {
	now := time.Now()

	// active records:
	activeFilter := bson.M{
		"$or": []bson.M{
			{"expires_at": bson.M{"$gt": now}}, // expires in future
			{"expires_at": time.Time{}},        // permanent
		},
	}

	// expired records:
	expiredFilter := bson.M{
		"expires_at": bson.M{"$lte": now}, // expired
	}

	return bson.M{
		"community_id": communityID,
		"type":         banType,
		"is_deleted":   false,
		"$and": []bson.M{
			func() bson.M {
				if expired {
					return expiredFilter
				}
				return activeFilter
			}(),
		},
	}
}

// FindCommunities finds communities with generic filter and options
func (c *communityRepo) FindCommunities(ctx context.Context, filter Filter, findOptions *FindOptions) ([]*model.Community, int64, error) {
	// Convert Filter to bson.M
	bsonFilter := bson.M(filter)

	// Count total documents
	total, err := c.communityCollection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, err
	}

	// Prepare MongoDB find options
	opts := options.Find()
	if findOptions != nil {
		if findOptions.Limit > 0 {
			opts.SetLimit(findOptions.Limit)
		}
		if findOptions.Skip > 0 {
			opts.SetSkip(findOptions.Skip)
		}
		if len(findOptions.Sort) > 0 {
			opts.SetSort(bson.D(func() []bson.E {
				var sortFields []bson.E
				for field, order := range findOptions.Sort {
					sortFields = append(sortFields, bson.E{Key: field, Value: order})
				}
				return sortFields
			}()))
		}
	}

	// Execute find
	cursor, err := c.communityCollection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var communities []*model.Community
	for cursor.Next(ctx) {
		var community model.Community
		if err := cursor.Decode(&community); err != nil {
			return nil, 0, err
		}
		communities = append(communities, &community)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return communities, total, nil
}

// Stats methods implementations
func (c *communityRepo) CountTotal(ctx context.Context) (int64, error) {
	return c.communityCollection.CountDocuments(ctx, bson.M{})
}

func (c *communityRepo) CountActive(ctx context.Context) (int64, error) {
	filter := bson.M{
		"is_deleted": bson.M{"$ne": true},
		"is_banned":  bson.M{"$ne": true},
	}
	return c.communityCollection.CountDocuments(ctx, filter)
}

func (c *communityRepo) CountCreatedAfter(ctx context.Context, since time.Time) (int64, error) {
	filter := bson.M{
		"create_at": bson.M{"$gte": since},
	}
	return c.communityCollection.CountDocuments(ctx, filter)
}

func (c *communityRepo) CountBanned(ctx context.Context) (int64, error) {
	filter := bson.M{
		"is_banned": true,
	}
	return c.communityCollection.CountDocuments(ctx, filter)
}

func (c *communityRepo) CountPrivate(ctx context.Context) (int64, error) {
	filter := bson.M{
		"setting.is_private": true,
	}
	return c.communityCollection.CountDocuments(ctx, filter)
}

func (c *communityRepo) ActivateModerator(ctx context.Context, communityID string, userID string) error {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return err
	}

	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id":                communityObjectID,
		"moderators.user_id": userObjectID,
	}

	update := bson.M{
		"$set": bson.M{
			"moderators.$[m].is_active":   true,
			"moderators.$[m].assigned_at": time.Now(),
		},
	}

	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"m.user_id": userObjectID},
		},
	})

	res, err := c.communityCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return apperror.ErrBadRequest
	}

	return nil
}

func (c *communityRepo) IsModerator(ctx context.Context, communityID string, userID string) (bool, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	community, err := c.GetByID(ctx, communityID)
	if err != nil {
		return false, err
	}

	if community.CreateByID == userObjectID {
		return true, nil
	}
	for _, m := range community.Moderators {
		if m.UserID == userObjectID {
			if !m.IsActive {
				return false, nil
			}

			return true, nil
		}
	}
	return false, nil
}

func (c *communityRepo) IsCreator(ctx context.Context, communityID string, userID string) (bool, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	community, err := c.GetByID(ctx, communityID)
	if err != nil {
		return false, err
	}

	return community.CreateByID == userObjectID, nil
}
