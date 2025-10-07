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

type MembershipRepo interface {
	Create(ctx context.Context, membership *model.Membership) (*model.Membership, error)
	GetByID(ctx context.Context, id string) (*model.Membership, error)
	GetByUserID(ctx context.Context, userID string) ([]model.Membership, error)
	GetAllPaginated(ctx context.Context, page int, pageSize int) ([]model.Membership, int64, error)
	GetByCommunityIDPaginated(ctx context.Context, communityID string, page int, pageSize int) ([]model.Membership, int64, error)
	Delete(ctx context.Context, id string) error

	CountMembersByCommunityID(ctx context.Context, communityID string) (int64, error)
	UpdateCommunityMemberCount(ctx context.Context, communityID string, count int64) error
}

type membershipRepo struct {
	membershipCollection *mongo.Collection
	communityCollection  *mongo.Collection
}

func NewMembershipRepo(db *mongo.Database) MembershipRepo {
	return &membershipRepo{
		membershipCollection: db.Collection(config.MembershipColName),
		communityCollection:  db.Collection(config.CommunityColName),
	}
}

func (m *membershipRepo) Create(ctx context.Context, membership *model.Membership) (*model.Membership, error) {
	result, err := m.membershipCollection.InsertOne(ctx, membership)
	if err != nil {
		return nil, fmt.Errorf("failed to create membership: %w", err)
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

	cursor, err := m.membershipCollection.Find(ctx, bson.M{"_id": membershipObjectID})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf(`membership with id %s not found`, id)
		}

		return nil, fmt.Errorf("failed to get membership by id: %w", err)
	}
	defer cursor.Close(ctx)

	var membership *model.Membership
	if err := cursor.Decode(&membership); err != nil {
		return nil, err
	}

	return membership, nil
}

func (m *membershipRepo) GetByUserID(ctx context.Context, userID string) ([]model.Membership, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	cursor, err := m.membershipCollection.Find(ctx, bson.M{"user_id": userObjectID})
	if err != nil {
		return nil, fmt.Errorf("failed to get membership by user id: %w", err)
	}
	defer cursor.Close(ctx)

	var memberships []model.Membership
	if err := cursor.All(ctx, &memberships); err != nil {
		return nil, err
	}

	return memberships, nil
}

func (m *membershipRepo) GetAllPaginated(ctx context.Context, page int, pageSize int) ([]model.Membership, int64, error) {
	skip := (page - 1) * pageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize))

	cursor, err := m.membershipCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all memberships: %w", err)
	}
	defer cursor.Close(ctx)

	var memberships []model.Membership
	if err := cursor.All(ctx, &memberships); err != nil {
		return nil, 0, fmt.Errorf("failed to decode memberships: %w", err)
	}

	count, err := m.membershipCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, -1, fmt.Errorf("failed to count memberships: %w", err)
	}

	return memberships, count, nil
}

func (m *membershipRepo) GetByCommunityIDPaginated(ctx context.Context, communityID string, page int, pageSize int) ([]model.Membership, int64, error) {
	communityObjectID, err := primitive.ObjectIDFromHex(communityID)
	if err != nil {
		return nil, 0, err
	}

	skip := (page - 1) * pageSize
	filter := bson.M{"community_id": communityObjectID}
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize))

	cursor, err := m.membershipCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get memberships by community id: %w", err)
	}
	defer cursor.Close(ctx)

	var memberships []model.Membership
	if err := cursor.All(ctx, &memberships); err != nil {
		return nil, 0, err
	}

	count, err := m.membershipCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, -1, fmt.Errorf("failed to query communities: %w", err)
	}

	return memberships, count, nil
}

func (m *membershipRepo) Delete(ctx context.Context, membershipID string) error {
	membershipObjectID, err := primitive.ObjectIDFromHex(membershipID)
	if err != nil {
		return fmt.Errorf("failed to parse id: %w", err)
	}

	result, err := m.membershipCollection.DeleteOne(ctx, bson.M{"_id": membershipObjectID})
	if err != nil {
		return fmt.Errorf("failed to delete membership: %w", err)
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
		return fmt.Errorf("invalid community id: %w", err)
	}

	filter := bson.M{"_id": communityObjectID}
	update := bson.M{"$set": bson.M{"member_count": count}}

	res, err := m.communityCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update community %s: %w", communityID, err)
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("community not found: %s", communityID)
	}

	return nil
}
