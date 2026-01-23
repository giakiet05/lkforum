package repo

import (
	"context"
	"fmt"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MembershipRepo interface {
	Create(ctx context.Context, membership *model.Membership) (*model.Membership, error)
	GetByID(ctx context.Context, id string) (*model.Membership, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.Membership, error)
	GetByUserIDAndCommunityID(ctx context.Context, userID string, communityID string) (*model.Membership, error)
	GetAllPaginated(ctx context.Context, page int, pageSize int) ([]*model.Membership, int64, error)
	GetByCommunityIDPaginated(ctx context.Context, communityID string, page int, pageSize int) ([]*model.Membership, int64, error)
	GetPendingByCommunityID(ctx context.Context, communityID string, page int, pageSize int) ([]*model.Membership, int64, error)
	GetApprovedByCommunityID(ctx context.Context, communityID string, page int, pageSize int) ([]*model.Membership, int64, error)
	UpdateStatus(ctx context.Context, membershipID string, status model.MembershipStatus) error
	Delete(ctx context.Context, id string) error

	CountMembersByCommunityID(ctx context.Context, communityID string) (int64, error)
	UpdateCommunityMemberCount(ctx context.Context, communityID string, count int64) error

	IsMember(ctx context.Context, userID string, communityID string) (bool, error)
	IsUserExist(ctx context.Context, userID string) (bool, error)
	IsCommunityExist(ctx context.Context, communityID string) (bool, error)
}

type membershipRepo struct {
	membershipCollection *mongo.Collection
	communityCollection  *mongo.Collection
	userCollection       *mongo.Collection
}

func NewMembershipRepo(db *mongo.Database) MembershipRepo {
	return &membershipRepo{
		membershipCollection: db.Collection(config.MembershipColName),
		communityCollection:  db.Collection(config.CommunityColName),
		userCollection:       db.Collection(config.UserColName),
	}
}

func (m *membershipRepo) Create(ctx context.Context, membership *model.Membership) (*model.Membership, error) {
	result, err := m.membershipCollection.InsertOne(ctx, membership)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		membership.ID = oid
	}
	return membership, nil
}

func (m *membershipRepo) GetByID(ctx context.Context, id string) (*model.Membership, error) {
	membershipObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var membership model.Membership
	err = m.membershipCollection.FindOne(ctx, bson.M{"_id": membershipObjectID}).Decode(&membership)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &membership, nil
}

func (m *membershipRepo) GetByUserID(ctx context.Context, userID string) ([]*model.Membership, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := m.membershipCollection.Find(ctx, bson.M{"user_id": userObjectID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var memberships []*model.Membership
	if err := cursor.All(ctx, &memberships); err != nil {
		return nil, err
	}

	return memberships, nil
}

func (m *membershipRepo) GetAllPaginated(ctx context.Context, page int, pageSize int) ([]*model.Membership, int64, error) {
	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize))

	cursor, err := m.membershipCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var memberships []*model.Membership
	if err := cursor.All(ctx, &memberships); err != nil {
		return nil, 0, err
	}

	count, err := m.membershipCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, -1, err
	}

	return memberships, count, nil
}

func (m *membershipRepo) GetByCommunityIDPaginated(ctx context.Context, communityID string, page int, pageSize int) ([]*model.Membership, int64, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	filter := bson.M{"community_id": communityObjectID}

	// Count total documents
	count, err := m.membershipCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, -1, err
	}

	// Use aggregation to join with users collection
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$skip", Value: int64(skip)}},
		{{Key: "$limit", Value: int64(pageSize)}},
		{{Key: "$lookup", Value: bson.M{
			"from":         config.UserColName,
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user_data",
		}}},
		{{Key: "$addFields", Value: bson.M{
			"user": bson.M{
				"_id":      bson.M{"$arrayElemAt": bson.A{"$user_data._id", 0}},
				"username": bson.M{"$arrayElemAt": bson.A{"$user_data.username", 0}},
				"avatar":   bson.M{"$arrayElemAt": bson.A{"$user_data.role_content.as_user.avatar", 0}},
			},
		}}},
		{{Key: "$project", Value: bson.M{
			"user_data": 0, // Remove the temporary array field
		}}},
	}

	cursor, err := m.membershipCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var memberships []*model.Membership
	if err := cursor.All(ctx, &memberships); err != nil {
		return nil, 0, err
	}

	return memberships, count, nil
}

func (m *membershipRepo) Delete(ctx context.Context, membershipID string) error {
	membershipObjectID, err := primitive.ObjectIDFromHex(membershipID)
	if err != nil {
		return err
	}

	result, err := m.membershipCollection.DeleteOne(ctx, bson.M{"_id": membershipObjectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no membership found with id %v", membershipID)
	}

	return nil
}

func (m *membershipRepo) CountMembersByCommunityID(ctx context.Context, communityID string) (int64, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return -1, err
	}

	filter := bson.M{"community_id": communityObjectID}
	count, err := m.membershipCollection.CountDocuments(ctx, filter)
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (m *membershipRepo) UpdateCommunityMemberCount(ctx context.Context, communityID string, count int64) error {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": communityObjectID}
	update := bson.M{"$set": bson.M{"member_count": count}}

	res, err := m.communityCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("community not found: %s", communityID)
	}

	return nil
}

func (m *membershipRepo) IsMember(ctx context.Context, userID string, communityID string) (bool, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	communityObjID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return false, err
	}

	filter := bson.M{
		"user_id":      userObjID,
		"community_id": communityObjID,
	}

	count, err := m.membershipCollection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m *membershipRepo) IsUserExist(ctx context.Context, userID string) (bool, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return false, err
	}

	count, err := m.userCollection.CountDocuments(ctx, bson.M{"_id": userObjectID})
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (m *membershipRepo) IsCommunityExist(ctx context.Context, communityID string) (bool, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return false, err
	}

	count, err := m.communityCollection.CountDocuments(ctx, bson.M{"_id": communityObjectID})
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (m *membershipRepo) GetByUserIDAndCommunityID(ctx context.Context, userID string, communityID string) (*model.Membership, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	communityObjID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{
		"user_id":      userObjID,
		"community_id": communityObjID,
	}

	var membership model.Membership
	err = m.membershipCollection.FindOne(ctx, filter).Decode(&membership)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &membership, nil
}

func (m *membershipRepo) GetPendingByCommunityID(ctx context.Context, communityID string, page int, pageSize int) ([]*model.Membership, int64, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	filter := bson.M{
		"community_id": communityObjectID,
		"status":       model.MembershipStatusPending,
	}

	// Count total documents
	count, err := m.membershipCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, -1, err
	}

	// Use aggregation to join with users collection
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.M{"created_at": -1}}}, // Newest first
		{{Key: "$skip", Value: int64(skip)}},
		{{Key: "$limit", Value: int64(pageSize)}},
		{{Key: "$lookup", Value: bson.M{
			"from":         config.UserColName,
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user_data",
		}}},
		{{Key: "$addFields", Value: bson.M{
			"user": bson.M{
				"_id":      bson.M{"$arrayElemAt": bson.A{"$user_data._id", 0}},
				"username": bson.M{"$arrayElemAt": bson.A{"$user_data.username", 0}},
				"avatar":   bson.M{"$arrayElemAt": bson.A{"$user_data.role_content.as_user.avatar", 0}},
			},
		}}},
		{{Key: "$project", Value: bson.M{
			"user_data": 0, // Remove the temporary array field
		}}},
	}

	cursor, err := m.membershipCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var memberships []*model.Membership
	if err := cursor.All(ctx, &memberships); err != nil {
		return nil, 0, err
	}

	return memberships, count, nil
}

func (m *membershipRepo) GetApprovedByCommunityID(ctx context.Context, communityID string, page int, pageSize int) ([]*model.Membership, int64, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	filter := bson.M{
		"community_id": communityObjectID,
		"status":       model.MembershipStatusApproved,
	}

	// Count total documents
	count, err := m.membershipCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, -1, err
	}

	// Use aggregation to join with users collection
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$sort", Value: bson.M{"created_at": -1}}}, // Newest first
		{{Key: "$skip", Value: int64(skip)}},
		{{Key: "$limit", Value: int64(pageSize)}},
		{{Key: "$lookup", Value: bson.M{
			"from":         config.UserColName,
			"localField":   "user_id",
			"foreignField": "_id",
			"as":           "user_data",
		}}},
		{{Key: "$addFields", Value: bson.M{
			"user": bson.M{
				"_id":      bson.M{"$arrayElemAt": bson.A{"$user_data._id", 0}},
				"username": bson.M{"$arrayElemAt": bson.A{"$user_data.username", 0}},
				"avatar":   bson.M{"$arrayElemAt": bson.A{"$user_data.role_content.as_user.avatar", 0}},
			},
		}}},
		{{Key: "$project", Value: bson.M{
			"user_data": 0, // Remove the temporary array field
		}}},
	}

	cursor, err := m.membershipCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var memberships []*model.Membership
	if err := cursor.All(ctx, &memberships); err != nil {
		return nil, 0, err
	}

	return memberships, count, nil
}

func (m *membershipRepo) UpdateStatus(ctx context.Context, membershipID string, status model.MembershipStatus) error {
	membershipObjectID, err := primitive.ObjectIDFromHex(membershipID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": membershipObjectID}
	update := bson.M{"$set": bson.M{"status": status}}

	res, err := m.membershipCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("membership not found: %s", membershipID)
	}

	return nil
}
