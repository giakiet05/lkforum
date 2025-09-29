package repo

import (
	"context"
	"fmt"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MembershipRepo interface {
	Create(ctx context.Context, membership *model.Membership) (*model.Membership, error)
	//GetAll(ctx context.Context) ([]*model.Membership, error)
	GetByID(ctx context.Context, id string) (*model.Membership, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.Membership, error)
	GetByCommunityID(ctx context.Context, communityID string) ([]*model.Membership, error)
	GetByRole(ctx context.Context, role model.CommunityRole) ([]*model.Membership, error)
	Update(ctx context.Context, membership *model.Membership) (*model.Membership, error)
	Delete(ctx context.Context, id string) error
}

type membershipRepo struct {
	membershipCollection *mongo.Collection
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
	cursor, err := m.membershipCollection.Find(ctx, bson.M{"_id": id})
	if err != nil {
		return nil, err
	}

	var membership *model.Membership
	if err := cursor.Decode(&membership); err != nil {
		return nil, err
	}

	return membership, nil
}

func (m *membershipRepo) GetByUserID(ctx context.Context, userID string) ([]*model.Membership, error) {
	cursor, err := m.membershipCollection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}

	var memberships []*model.Membership
	if err := cursor.Decode(&memberships); err != nil {
		return nil, err
	}

	return memberships, nil
}

func (m *membershipRepo) GetByCommunityID(ctx context.Context, communityID string) ([]*model.Membership, error) {
	//TODO implement me
	panic("implement me")
}

func (m *membershipRepo) GetByRole(ctx context.Context, role model.CommunityRole) ([]*model.Membership, error) {
	//TODO implement me
	panic("implement me")
}

func (m *membershipRepo) Update(ctx context.Context, membership *model.Membership) (*model.Membership, error) {
	//TODO implement me
	panic("implement me")
}

func (m *membershipRepo) Delete(ctx context.Context, membershipID string) error {
	result, err := m.membershipCollection.DeleteOne(ctx, bson.M{"_id": membershipID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no membership found with id %v", membershipID)
	}

	return nil
}

func NewMembershipRepo(db *mongo.Database) MembershipRepo {
	return &membershipRepo{membershipCollection: db.Collection(config.MembershipColName)}
}
