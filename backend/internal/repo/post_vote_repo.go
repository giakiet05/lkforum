package repo

import (
	"context"
	"errors"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ErrVoteNotFound = errors.New("vote not found")

type PostVoteRepo interface {
	Vote(ctx context.Context, vote *model.Vote) error
	RemoveVote(ctx context.Context, postID, userID primitive.ObjectID) error
	GetUserVote(ctx context.Context, postID, userID primitive.ObjectID) (*model.Vote, error)
}
type postVoteRepo struct {
	client         *mongo.Client
	postCollection *mongo.Collection
	voteCollection *mongo.Collection
}

func NewPostVoteRepo(client *mongo.Client, db *mongo.Database) PostVoteRepo {
	return &postVoteRepo{
		client:         client,
		postCollection: db.Collection(config.PostColName),
		voteCollection: db.Collection(config.VoteColName),
	}
}

// GetUserVote lấy phiếu bầu hiện tại của một người dùng cho một bài đăng.
// Trả về (nil, nil) nếu người dùng chưa vote.
func (r *postVoteRepo) GetUserVote(ctx context.Context, postID, userID primitive.ObjectID) (*model.Vote, error) {
	var vote model.Vote
	filter := bson.M{
		"target_id":   postID,
		"user_id":     userID,
		"target_type": model.VoteTargetPost,
	}

	err := r.voteCollection.FindOne(ctx, filter).Decode(&vote)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Người dùng chưa vote, đây không phải là lỗi
		}
		return nil, err
	}
	return &vote, nil
}

// Vote xử lý việc bỏ phiếu, thay đổi phiếu hoặc hủy phiếu.
// Đây là hàm phức tạp nhất vì nó xử lý 3 trường hợp:
// 1. Vote mới.
// 2. Thay đổi vote (up -> down hoặc down -> up).
// 3. Hủy vote (up -> up hoặc down -> down).
func (r *postVoteRepo) Vote(ctx context.Context, vote *model.Vote) error {
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Lấy vote trước đó của người dùng (nếu có)
		previousVote, err := r.GetUserVote(sessCtx, vote.TargetID, vote.UserID)
		if err != nil {
			return nil, err
		}

		// Khởi tạo các biến để cập nhật counter
		upInc, downInc := 0, 0

		if previousVote == nil {
			// TRƯỜNG HỢP 1: VOTE MỚI
			// Ghi vote mới vào collection 'votes'
			if _, err := r.voteCollection.InsertOne(sessCtx, vote); err != nil {
				return nil, err
			}
			// Cập nhật counter
			if vote.Value { // true = upvote
				upInc = 1
			} else {
				downInc = 1
			}
		} else {
			// Người dùng đã vote trước đó
			if previousVote.Value == vote.Value {
				// TRƯỜNG HỢP 2: HỦY VOTE (bấm lại nút cũ)
				if err := r.removeVoteInTransaction(sessCtx, previousVote); err != nil {
					return nil, err
				}
				// Không cần làm gì thêm, removeVoteInTransaction đã xử lý counter
				return nil, nil
			} else {
				// TRƯỜNG HỢP 3: THAY ĐỔI VOTE (up -> down hoặc down -> up)
				// Cập nhật vote trong collection 'votes'
				filter := bson.M{"_id": previousVote.ID}
				update := bson.M{"$set": bson.M{"value": vote.Value}}
				if _, err := r.voteCollection.UpdateOne(sessCtx, filter, update); err != nil {
					return nil, err
				}

				// Cập nhật counter
				if vote.Value { // Chuyển sang upvote (trước đó là downvote)
					upInc = 1
					downInc = -1
				} else { // Chuyển sang downvote (trước đó là upvote)
					upInc = -1
					downInc = 1
				}
			}
		}

		// Áp dụng thay đổi counter vào collection 'posts'
		postFilter := bson.M{"_id": vote.TargetID}
		postUpdate := bson.M{"$inc": bson.M{"votes_count.up": upInc, "votes_count.down": downInc}}
		result, err := r.postCollection.UpdateOne(sessCtx, postFilter, postUpdate)
		if err != nil {
			return nil, err
		}
		if result.MatchedCount == 0 {
			return nil, ErrPostNotFound
		}
		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

// RemoveVote xóa hoàn toàn một phiếu bầu.
func (r *postVoteRepo) RemoveVote(ctx context.Context, postID, userID primitive.ObjectID) error {
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		voteToRemove, err := r.GetUserVote(sessCtx, postID, userID)
		if err != nil {
			return nil, err
		}
		if voteToRemove == nil {
			return nil, ErrVoteNotFound // Không có vote để xóa
		}

		return nil, r.removeVoteInTransaction(sessCtx, voteToRemove)
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

// removeVoteInTransaction là hàm helper để xóa vote và cập nhật counter bên trong một transaction.
func (r *postVoteRepo) removeVoteInTransaction(sessCtx mongo.SessionContext, vote *model.Vote) error {
	// 1. Xóa vote khỏi collection 'votes'
	filter := bson.M{"_id": vote.ID}
	_, err := r.voteCollection.DeleteOne(sessCtx, filter)
	if err != nil {
		return err
	}

	// 2. Cập nhật (giảm) counter trong collection 'posts'
	upInc, downInc := 0, 0
	if vote.Value { // true = upvote
		upInc = -1
	} else {
		downInc = -1
	}

	postFilter := bson.M{"_id": vote.TargetID}
	postUpdate := bson.M{"$inc": bson.M{"votes_count.up": upInc, "votes_count.down": downInc}}
	result, err := r.postCollection.UpdateOne(sessCtx, postFilter, postUpdate)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrPostNotFound
	}
	return nil
}
