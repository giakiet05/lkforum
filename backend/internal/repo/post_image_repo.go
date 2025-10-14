package repo

import (
	"context"
	"errors"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PostImageRepo interface {
	AddImages(ctx context.Context, postID primitive.ObjectID, images []model.Image) error
	RemoveImages(ctx context.Context, postID primitive.ObjectID, imageIDs []primitive.ObjectID) error
	GetPostImages(ctx context.Context, postID primitive.ObjectID) ([]model.Image, error)
}

type postImageRepo struct {
	postCollection *mongo.Collection
}

func NewPostImageRepo(db *mongo.Database) PostImageRepo {
	return &postImageRepo{postCollection: db.Collection(config.PostColName)}
}

// AddImages thêm một hoặc nhiều ảnh vào mảng 'images' của một bài đăng.
func (r *postImageRepo) AddImages(ctx context.Context, postID primitive.ObjectID, images []model.Image) error {
	filter := bson.M{"_id": postID}

	// SỬA: Đường dẫn được cập nhật thành "content.images"
	update := bson.M{
		"$push": bson.M{
			"content.images": bson.M{"$each": images},
		},
	}

	result, err := r.postCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrPostNotFound
	}

	return nil
}

// RemoveImages xóa một hoặc nhiều ảnh khỏi mảng 'content.images' dựa trên ID của chúng.
func (r *postImageRepo) RemoveImages(ctx context.Context, postID primitive.ObjectID, imageIDs []primitive.ObjectID) error {
	filter := bson.M{"_id": postID}

	// SỬA: Đường dẫn được cập nhật thành "content.images"
	update := bson.M{
		"$pull": bson.M{
			"content.images": bson.M{"_id": bson.M{"$in": imageIDs}},
		},
	}

	result, err := r.postCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrPostNotFound
	}

	return nil
}

// GetPostImages lấy về danh sách tất cả các ảnh của một bài đăng.
func (r *postImageRepo) GetPostImages(ctx context.Context, postID primitive.ObjectID) ([]model.Image, error) {
	var post model.Post

	filter := bson.M{"_id": postID}

	// SỬA: Projection được cập nhật để chỉ lấy về trường "content.images"
	projection := options.FindOne().SetProjection(bson.M{"content.images": 1, "_id": 0})

	err := r.postCollection.FindOne(ctx, filter, projection).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	// SỬA: Kiểm tra xem Content có nil không trước khi truy cập Images
	if post.Content == nil || post.Content.Images == nil {
		return []model.Image{}, nil
	}

	return post.Content.Images, nil
}
