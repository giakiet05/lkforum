package repo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo interface {
	Create(ctx context.Context, user *model.User) (*model.User, error)
	Update(ctx context.Context, user *model.User) (*model.User, error)
	Delete(ctx context.Context, id string) error
	UpdateReputation(ctx context.Context, userID string, points int) error

	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Find(ctx context.Context, filter Filter, opts *FindOptions) ([]*model.User, int64, error)

	// User activity stats methods
	IncrementPostCount(ctx context.Context, userID string, delta int) error
	IncrementCommentCount(ctx context.Context, userID string, delta int) error
	IncrementTotalUpvotes(ctx context.Context, userID string, delta int) error

	// Stats methods
	CountTotal(ctx context.Context) (int64, error)
	CountActiveAfter(ctx context.Context, since time.Time) (int64, error)
	CountCreatedAfter(ctx context.Context, since time.Time) (int64, error)
	CountBanned(ctx context.Context) (int64, error)
	CountVerified(ctx context.Context) (int64, error)
}

type userRepo struct {
	userCollection *mongo.Collection
}

func NewUserRepo(db *mongo.Database) UserRepo {
	return &userRepo{userCollection: db.Collection(config.UserColName)}
}

func (r *userRepo) GetByIDs(ctx context.Context, ids []string) ([]*model.User, error) {
	if len(ids) == 0 {
		return []*model.User{}, nil
	}

	objIDs := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if objID, err := primitive.ObjectIDFromHex(id); err == nil {
			objIDs = append(objIDs, objID)
		}
	}

	filter := bson.M{"_id": bson.M{"$in": objIDs}}
	cursor, err := r.userCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

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
		return apperror.ErrInvalidID
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

func (r *userRepo) UpdateReputation(ctx context.Context, userID string, points int) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$inc": bson.M{"reputation": points},
	}

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
		return nil, apperror.ErrInvalidID
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

// Find fetches users with filter and pagination options
func (r *userRepo) Find(ctx context.Context, filter Filter, opts *FindOptions) ([]*model.User, int64, error) {
	// Get total count
	countPipeline := mongo.Pipeline{
		{{"$match", bson.M(filter)}},
		{{"$count", "total"}},
	}
	cursor, err := r.userCollection.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}

	var countResult []struct {
		Total int64 `bson:"total"`
	}
	if err = cursor.All(ctx, &countResult); err != nil {
		return nil, 0, err
	}

	var total int64
	if len(countResult) > 0 {
		total = countResult[0].Total
	}

	// If no documents, return early
	if total == 0 {
		return []*model.User{}, 0, nil
	}

	// Get paginated data
	findOptions := options.Find()
	if opts != nil {
		if opts.Sort != nil {
			sortDoc := bson.D{}
			for key, value := range opts.Sort {
				sortDoc = append(sortDoc, bson.E{Key: key, Value: value})
			}
			findOptions.SetSort(sortDoc)
		}
		if opts.Skip > 0 {
			findOptions.SetSkip(opts.Skip)
		}
		if opts.Limit > 0 {
			findOptions.SetLimit(opts.Limit)
		}
	}

	cursor, err = r.userCollection.Find(ctx, bson.M(filter), findOptions)
	if err != nil {
		return nil, total, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Stats methods implementations
func (r *userRepo) CountTotal(ctx context.Context) (int64, error) {
	return r.userCollection.CountDocuments(ctx, bson.M{})
}

func (r *userRepo) CountActiveAfter(ctx context.Context, since time.Time) (int64, error) {
	filter := bson.M{
		"last_login": bson.M{"$gte": since},
	}
	return r.userCollection.CountDocuments(ctx, filter)
}

func (r *userRepo) CountCreatedAfter(ctx context.Context, since time.Time) (int64, error) {
	filter := bson.M{
		"created_at": bson.M{"$gte": since},
	}
	return r.userCollection.CountDocuments(ctx, filter)
}

func (r *userRepo) CountBanned(ctx context.Context) (int64, error) {
	filter := bson.M{
		"is_banned": true,
	}
	return r.userCollection.CountDocuments(ctx, filter)
}

func (r *userRepo) CountVerified(ctx context.Context) (int64, error) {
	filter := bson.M{
		"is_verified": true,
	}
	return r.userCollection.CountDocuments(ctx, filter)
}

// IncrementPostCount increments or decrements user's post count
func (r *userRepo) IncrementPostCount(ctx context.Context, userID string, delta int) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$inc": bson.M{"role_content.as_user.stats.post_count": delta},
	}

	_, err = r.userCollection.UpdateOne(ctx, filter, update)
	return err
}

// IncrementCommentCount increments or decrements user's comment count
func (r *userRepo) IncrementCommentCount(ctx context.Context, userID string, delta int) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$inc": bson.M{"role_content.as_user.stats.comment_count": delta},
	}

	_, err = r.userCollection.UpdateOne(ctx, filter, update)
	return err
}

// IncrementTotalUpvotes increments or decrements user's total upvotes
func (r *userRepo) IncrementTotalUpvotes(ctx context.Context, userID string, delta int) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$inc": bson.M{"role_content.as_user.stats.total_upvotes": delta},
	}

	_, err = r.userCollection.UpdateOne(ctx, filter, update)
	return err
}
