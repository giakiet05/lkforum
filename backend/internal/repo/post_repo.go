// Filename: internal/repo/post_repo.go

package repo

import (
	"context"
	"errors"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrPostNotFound là lỗi được trả về khi không tìm thấy bài đăng.
var ErrPostNotFound = errors.New("post not found")

// PostRepo định nghĩa các phương thức để tương tác với dữ liệu bài đăng.
type PostRepo interface {
	Create(ctx context.Context, post *model.Post) (*model.Post, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*model.Post, error)
	GetByCommunityID(ctx context.Context, communityID primitive.ObjectID, sort string, timeFrame string, limit, offset int) ([]*model.Post, error)
	GetByAuthorID(ctx context.Context, authorID primitive.ObjectID, limit, offset int) ([]*model.Post, error)
	GetFeed(ctx context.Context, sort string, timeFrame string, limit, offset int) ([]*model.Post, error)
	GetPostsByType(ctx context.Context, postType model.PostType, limit, offset int) ([]*model.Post, error)
	Update(ctx context.Context, post *model.Post) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	SoftDelete(ctx context.Context, id primitive.ObjectID) error
	IncrementCommentsCount(ctx context.Context, postID primitive.ObjectID) error
	DecrementCommentsCount(ctx context.Context, postID primitive.ObjectID) error
}

type postRepo struct {
	postCollection *mongo.Collection
}

// NewPostRepo khởi tạo một implementation mới của PostRepo.
func NewPostRepo(db *mongo.Database) PostRepo {
	return &postRepo{postCollection: db.Collection(config.PostColName)}
}

// Create tạo một bài đăng mới trong cơ sở dữ liệu.
func (r *postRepo) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	post.CreatedAt = time.Now()
	post.VotesCount = &model.VotesCount{Up: 0, Down: 0}
	post.CommentsCount = 0
	post.IsDeleted = false

	result, err := r.postCollection.InsertOne(ctx, post)
	if err != nil {
		return nil, err
	}

	post.ID = result.InsertedID.(primitive.ObjectID)
	return post, nil
}

// GetByID lấy một bài đăng bằng ID, không bao gồm các bài đã bị xóa mềm.
func (r *postRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*model.Post, error) {
	var post model.Post
	filter := bson.M{
		"_id":        id,
		"is_deleted": bson.M{"$ne": true},
	}

	err := r.postCollection.FindOne(ctx, filter).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return &post, nil
}

// GetByCommunityID lấy danh sách bài đăng thuộc một cộng đồng cụ thể.
func (r *postRepo) GetByCommunityID(ctx context.Context, communityID primitive.ObjectID, sort string, timeFrame string, limit, offset int) ([]*model.Post, error) {
	filter := bson.M{
		"community_id": communityID,
		"is_deleted":   bson.M{"$ne": true},
	}

	// Áp dụng bộ lọc thời gian
	for k, v := range buildTimeFilter(timeFrame) {
		filter[k] = v
	}

	return r.findPosts(ctx, filter, sort, limit, offset)
}

// GetByAuthorID lấy danh sách bài đăng của một tác giả cụ thể.
func (r *postRepo) GetByAuthorID(ctx context.Context, authorID primitive.ObjectID, limit, offset int) ([]*model.Post, error) {
	filter := bson.M{
		"author_id":  authorID,
		"is_deleted": bson.M{"$ne": true},
	}
	// Mặc định sắp xếp theo bài mới nhất cho trang cá nhân
	return r.findPosts(ctx, filter, "new", limit, offset)
}

// GetFeed lấy danh sách bài đăng cho trang chủ hoặc feed chung.
func (r *postRepo) GetFeed(ctx context.Context, sort string, timeFrame string, limit, offset int) ([]*model.Post, error) {
	filter := bson.M{"is_deleted": bson.M{"$ne": true}}

	// Áp dụng bộ lọc thời gian
	for k, v := range buildTimeFilter(timeFrame) {
		filter[k] = v
	}

	return r.findPosts(ctx, filter, sort, limit, offset)
}

// GetPostsByType lấy danh sách bài đăng theo một loại cụ thể (text, image, ...).
func (r *postRepo) GetPostsByType(ctx context.Context, postType model.PostType, limit, offset int) ([]*model.Post, error) {
	filter := bson.M{
		"type":       postType,
		"is_deleted": bson.M{"$ne": true},
	}
	return r.findPosts(ctx, filter, "new", limit, offset)
}

// Update cập nhật thông tin của một bài đăng.
func (r *postRepo) Update(ctx context.Context, post *model.Post) error {
	filter := bson.M{"_id": post.ID}
	update := bson.M{
		"$set": bson.M{
			"title":      post.Title,
			"content":    post.Content,
			"updated_at": time.Now(),
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

// Delete xóa vĩnh viễn một bài đăng khỏi cơ sở dữ liệu.
func (r *postRepo) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	result, err := r.postCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrPostNotFound
	}

	return nil
}

// SoftDelete đánh dấu một bài đăng là đã bị xóa.
func (r *postRepo) SoftDelete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id, "is_deleted": bson.M{"$ne": true}}
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"updated_at": time.Now(),
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

// IncrementCommentsCount tăng số lượng bình luận của bài đăng lên 1.
func (r *postRepo) IncrementCommentsCount(ctx context.Context, postID primitive.ObjectID) error {
	filter := bson.M{"_id": postID}
	update := bson.M{"$inc": bson.M{"comment_count": 1}}
	result, err := r.postCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrPostNotFound
	}
	return nil
}

// DecrementCommentsCount giảm số lượng bình luận của bài đăng đi 1.
// Hàm này sẽ không giảm nếu số lượng bình luận đã là 0.
func (r *postRepo) DecrementCommentsCount(ctx context.Context, postID primitive.ObjectID) error {
	filter := bson.M{"_id": postID, "comment_count": bson.M{"$gt": 0}}
	update := bson.M{"$inc": bson.M{"comment_count": -1}}
	_, err := r.postCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	// Không trả về lỗi nếu MatchedCount = 0, vì có thể do comment_count đã bằng 0
	// Service layer có thể cần biết điều này, nhưng ở repo layer thì đây không phải là một lỗi.
	return nil
}

// --- Helper Functions ---

// findPosts là hàm helper chung để thực hiện truy vấn tìm kiếm và phân trang.
func (r *postRepo) findPosts(ctx context.Context, filter bson.M, sort string, limit, offset int) ([]*model.Post, error) {
	var posts []*model.Post
	sortOptions := buildSortOptions(sort)

	findOptions := options.Find().
		SetSort(sortOptions).
		SetSkip(int64(offset)).
		SetLimit(int64(limit))

	cursor, err := r.postCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
	}

	// Nếu không có kết quả, trả về một slice rỗng thay vì nil
	if posts == nil {
		posts = []*model.Post{}
	}

	return posts, nil
}

// buildTimeFilter tạo bộ lọc thời gian dựa trên chuỗi đầu vào.
func buildTimeFilter(timeFrame string) bson.M {
	if timeFrame == "" || timeFrame == "all" {
		return bson.M{}
	}

	var duration time.Duration
	now := time.Now()
	switch timeFrame {
	case "day":
		duration = 24 * time.Hour
	case "week":
		duration = 7 * 24 * time.Hour
	case "month":
		duration = 30 * 24 * time.Hour // ~1 tháng
	case "year":
		duration = 365 * 24 * time.Hour // ~1 năm
	default:
		return bson.M{}
	}

	return bson.M{"created_at": bson.M{"$gte": now.Add(-duration)}}
}

// buildSortOptions tạo các tùy chọn sắp xếp dựa trên chuỗi đầu vào.
func buildSortOptions(sort string) bson.D {
	switch sort {
	case "hot":
		// Ghi chú: Thuật toán "hot" đơn giản, có thể được cải thiện sau.
		return bson.D{{"votes_count.up", -1}, {"created_at", -1}}
	case "top":
		return bson.D{{"votes_count.up", -1}}
	case "controversial":
		return bson.D{{"votes_count.down", -1}, {"votes_count.up", -1}}
	case "new":
		return bson.D{{"created_at", -1}}
	default:
		// Mặc định luôn là bài đăng mới nhất
		return bson.D{{"created_at", -1}}
	}
}
