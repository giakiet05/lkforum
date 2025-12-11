package repo

import (
	"context"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CommentRepo interface {
	Create(ctx context.Context, comment *model.Comment) (*model.Comment, error)
	GetByID(ctx context.Context, commentID string) (*model.Comment, error)
	GetCommentsFilterPaginated(
		ctx context.Context,
		postID *string, parentID *string, userID *string, content *string,
		currentUserID *string,
		page int, pageSize int,
	) ([]model.Comment, int64, error)
	GetByParentIDsPaginated(ctx context.Context, parentIDs []string, page int, pageSize int) ([]model.Comment, error)
	GetAllChildren(ctx context.Context, commentID string) ([]model.Comment, error)
	Delete(ctx context.Context, commentID string) error
	UpdateByID(ctx context.Context, commentID string, update bson.M) error

	// Stats methods
	CountTotal(ctx context.Context) (int64, error)
	CountCreatedAfter(ctx context.Context, since time.Time) (int64, error)
}

type commentRepo struct {
	commentCollection *mongo.Collection
}

func NewCommentRepo(db *mongo.Database) CommentRepo {
	return &commentRepo{
		commentCollection: db.Collection(config.CommentColName),
	}
}

func (c *commentRepo) Create(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	result, err := c.commentCollection.InsertOne(ctx, comment)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		comment.ID = oid
	}

	return comment, nil
}

func (c *commentRepo) GetByID(ctx context.Context, commentID string) (*model.Comment, error) {
	commentObjectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, err
	}

	result := c.commentCollection.FindOne(ctx, bson.M{"_id": commentObjectID})
	var comment model.Comment
	err = result.Decode(&comment)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (c *commentRepo) GetCommentsFilterPaginated(
	ctx context.Context,
	postID *string,
	parentID *string,
	userID *string,
	content *string,
	currentUserID *string,
	page int, pageSize int,
) ([]model.Comment, int64, error) {

	filter := bson.M{}

	// Build filter dynamically
	if postID != nil && *postID != "" {
		postObjectID, err := primitive.ObjectIDFromHex(*postID)
		if err != nil {
			return nil, 0, apperror.ErrInvalidID
		}
		filter["post_id"] = postObjectID
		filter["parent_id"] = nil
	}

	if parentID != nil && *parentID != "" {
		parentObjectID, err := primitive.ObjectIDFromHex(*parentID)
		if err != nil {
			return nil, 0, apperror.ErrInvalidID
		}
		filter["parent_id"] = parentObjectID
	}

	if userID != nil && *userID != "" {
		userObjectID, err := primitive.ObjectIDFromHex(*userID)
		if err != nil {
			return nil, 0, apperror.ErrInvalidID
		}
		filter["author.id"] = userObjectID
	}

	if content != nil && *content != "" {
		// Case-insensitive partial match
		filter["content"] = bson.M{"$regex": primitive.Regex{Pattern: *content, Options: "i"}}
	}

	// No moderation filter needed - all comments are visible

	// Pagination and sorting
	skip := (page - 1) * pageSize
	opt := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.M{"created_at": 1})

	cursor, err := c.commentCollection.Find(ctx, filter, opt)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var comments []model.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, 0, err
	}

	count, err := c.commentCollection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return comments, count, nil
}

func (c *commentRepo) GetByParentIDsPaginated(ctx context.Context, parentIDs []string, page int, pageSize int) ([]model.Comment, error) {
	if len(parentIDs) == 0 {
		return nil, nil
	}

	var allChildren []model.Comment
	skip := (page - 1) * pageSize
	opt := options.Find().SetSkip(int64(skip)).SetLimit(int64(pageSize)).SetSort(bson.M{"created_at": 1})
	for _, pid := range parentIDs {
		parentObjectID, err := primitive.ObjectIDFromHex(pid)
		if err != nil {
			return nil, apperror.ErrInvalidID
		}

		cursor, err := c.commentCollection.Find(ctx, bson.M{"parent_id": parentObjectID}, opt)
		if err != nil {
			return nil, err
		}

		var children []model.Comment
		if err := cursor.All(ctx, &children); err != nil {
			cursor.Close(ctx)
			return nil, err
		}
		cursor.Close(ctx)

		allChildren = append(allChildren, children...)
	}

	return allChildren, nil
}

func (c *commentRepo) GetAllChildren(ctx context.Context, commentID string) ([]model.Comment, error) {
	parentObjectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return nil, err
	}

	var allChildren []model.Comment
	// recursive helper
	var fetch func(parentID primitive.ObjectID) error
	fetch = func(parentID primitive.ObjectID) error {
		cursor, err := c.commentCollection.Find(ctx, bson.M{"parent_id": parentID})
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)

		var children []model.Comment
		if err = cursor.All(ctx, &children); err != nil {
			return err
		}

		allChildren = append(allChildren, children...)

		for _, child := range children {
			if err := fetch(child.ID); err != nil {
				return err
			}
		}

		return nil
	}

	if err := fetch(parentObjectID); err != nil {
		return nil, err
	}

	return allChildren, nil
}

func (c *commentRepo) Delete(ctx context.Context, commentID string) error {
	commentObjectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted":      true,
			"content":         "[deleted]",
			"author.id":       nil,
			"author.username": "[deleted]",
			"author.avatar":   "[deleted]",
			"deleted_at":      time.Now(),
		},
	}

	result, err := c.commentCollection.UpdateOne(ctx, bson.M{"_id": commentObjectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return apperror.ErrCommentNotFound
	}

	return nil
}

func (c *commentRepo) UpdateByID(ctx context.Context, commentID string, update bson.M) error {
	commentObjectID, err := primitive.ObjectIDFromHex(commentID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	result, err := c.commentCollection.UpdateOne(ctx, bson.M{"_id": commentObjectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return apperror.ErrCommentNotFound
	}

	return nil
}

// Stats methods implementations
func (c *commentRepo) CountTotal(ctx context.Context) (int64, error) {
	return c.commentCollection.CountDocuments(ctx, bson.M{})
}

func (c *commentRepo) CountCreatedAfter(ctx context.Context, since time.Time) (int64, error) {
	filter := bson.M{
		"created_at": bson.M{"": since},
	}
	return c.commentCollection.CountDocuments(ctx, filter)
}
