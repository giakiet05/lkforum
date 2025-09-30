package repo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id string) error

	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	GetPaginated(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
}

type userRepo struct {
	userCollection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) UserRepo {
	return &userRepo{userCollection: db.Collection(config.UserColName)}
}

func (r *userRepo) GetAll(ctx context.Context) ([]*model.User, error) {
	cursor, err := r.userCollection.Find(ctx, bson.M{"deleted_at": bson.M{"$exists": false}})
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepo) Create(ctx context.Context, user *model.User) (*model.User, error) {
	result, err := r.userCollection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid
	}

	return user, nil
}

func (r *userRepo) Update(ctx context.Context, user *model.User) (*model.User, error) {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	result, err := r.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return user, nil
}

func (r *userRepo) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objectID, "deleted_at": bson.M{"$exists": false}}
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	result, err := r.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objectID, "deleted_at": bson.M{"$exists": false}}
	var user model.User
	err = r.userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	filter := bson.M{"username": username, "deleted_at": bson.M{"$exists": false}}
	var user model.User
	err := r.userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	filter := bson.M{"email": email, "deleted_at": bson.M{"$exists": false}}
	var user model.User
	err := r.userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetPaginated(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	skip := (page - 1) * pageSize
	filter := bson.M{"deleted_at": bson.M{"$exists": false}}
	cursor, err := r.userCollection.Find(ctx, filter, options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)))
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	count, err := r.userCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return users, count, nil
}
