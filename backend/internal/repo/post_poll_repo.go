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

var ErrPollVoted = errors.New("user has already voted for this option")
var ErrPollCannotEdit = errors.New("poll cannot be edited after votes have been cast")

type PostPollRepo interface {
	VotePoll(ctx context.Context, pollVote *model.PollVote) error
	RemovePollVote(ctx context.Context, postID, userID primitive.ObjectID) error
	GetUserPollVotes(ctx context.Context, postID, userID primitive.ObjectID) ([]*model.PollVote, error)
	UpdatePoll(ctx context.Context, postID primitive.ObjectID, question string, expiresAt *time.Time, allowMultiple *bool) error
	AddPollOptions(ctx context.Context, postID primitive.ObjectID, options []model.PollOption) error
	RemovePollOptions(ctx context.Context, postID primitive.ObjectID, optionIDs []primitive.ObjectID) error
	UpdatePollOption(ctx context.Context, postID, optionID primitive.ObjectID, text string) error
	CanEditPoll(ctx context.Context, postID primitive.ObjectID) (bool, error)
}
type postPollRepo struct {
	client             *mongo.Client // Cần client để tạo session cho transaction
	postCollection     *mongo.Collection
	pollVoteCollection *mongo.Collection
}

func NewPostPollRepo(client *mongo.Client, db *mongo.Database) PostPollRepo {
	return &postPollRepo{
		client:             client,
		postCollection:     db.Collection(config.PostColName),
		pollVoteCollection: db.Collection(config.PollVoteColName), // Thêm collection cho poll votes
	}
}

// GetUserPollVotes lấy tất cả các phiếu bầu của một người dùng cho một bài đăng.
func (r *postPollRepo) GetUserPollVotes(ctx context.Context, postID, userID primitive.ObjectID) ([]*model.PollVote, error) {
	var votes []*model.PollVote
	filter := bson.M{"post_id": postID, "user_id": userID}
	cursor, err := r.pollVoteCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &votes); err != nil {
		return nil, err
	}
	return votes, nil
}

// UpdatePoll cập nhật các thuộc tính chính của poll.
func (r *postPollRepo) UpdatePoll(ctx context.Context, postID primitive.ObjectID, question string, expiresAt *time.Time, allowMultiple *bool) error {
	canEdit, err := r.CanEditPoll(ctx, postID)
	if err != nil {
		return err
	}
	if !canEdit {
		return ErrPollCannotEdit
	}

	setFields := bson.M{}
	if question != "" {
		setFields["content.poll.question"] = question
	}
	if expiresAt != nil {
		setFields["content.poll.expires_at"] = expiresAt
	}
	if allowMultiple != nil {
		setFields["content.poll.allow_multiple"] = allowMultiple
	}

	if len(setFields) == 0 {
		return nil // Không có gì để cập nhật
	}

	filter := bson.M{"_id": postID}
	update := bson.M{"$set": setFields}

	result, err := r.postCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrPostNotFound
	}
	return nil
}

// AddPollOptions thêm các lựa chọn mới vào poll.
func (r *postPollRepo) AddPollOptions(ctx context.Context, postID primitive.ObjectID, options []model.PollOption) error {
	canEdit, err := r.CanEditPoll(ctx, postID)
	if err != nil {
		return err
	}
	if !canEdit {
		return ErrPollCannotEdit
	}

	// Gán ID mới cho các option chưa có ID
	for i := range options {
		if options[i].ID.IsZero() {
			options[i].ID = primitive.NewObjectID()
		}
	}

	filter := bson.M{"_id": postID}
	update := bson.M{"$push": bson.M{"content.poll.options": bson.M{"$each": options}}}

	result, err := r.postCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrPostNotFound
	}
	return nil
}

// RemovePollOptions xóa các lựa chọn khỏi poll.
func (r *postPollRepo) RemovePollOptions(ctx context.Context, postID primitive.ObjectID, optionIDs []primitive.ObjectID) error {
	// Chú ý: Thao tác này có thể làm cho các phiếu bầu cũ trỏ đến option không còn tồn tại.
	// Một giải pháp tốt hơn có thể là "soft delete" các option.
	canEdit, err := r.CanEditPoll(ctx, postID)
	if err != nil {
		return err
	}
	if !canEdit {
		return ErrPollCannotEdit
	}

	filter := bson.M{"_id": postID}
	update := bson.M{"$pull": bson.M{"content.poll.options": bson.M{"_id": bson.M{"$in": optionIDs}}}}

	result, err := r.postCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrPostNotFound
	}
	return nil
}

// UpdatePollOption cập nhật nội dung text của một lựa chọn.
func (r *postPollRepo) UpdatePollOption(ctx context.Context, postID, optionID primitive.ObjectID, text string) error {
	filter := bson.M{"_id": postID, "content.poll.options._id": optionID}
	update := bson.M{"$set": bson.M{"content.poll.options.$.text": text}}

	result, err := r.postCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrPostNotFound
	}
	return nil
}

// CanEditPoll kiểm tra xem poll có thể được chỉnh sửa hay không (khi chưa có ai vote).
func (r *postPollRepo) CanEditPoll(ctx context.Context, postID primitive.ObjectID) (bool, error) {
	var post model.Post
	filter := bson.M{"_id": postID}
	projection := options.FindOne().SetProjection(bson.M{"content.poll.total_votes": 1})

	err := r.postCollection.FindOne(ctx, filter, projection).Decode(&post)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, ErrPostNotFound
		}
		return false, err
	}

	canEdit := post.Content != nil && post.Content.Poll != nil && post.Content.Poll.TotalVotes != nil && *post.Content.Poll.TotalVotes == 0
	return canEdit, nil
}

// VotePoll xử lý việc người dùng bỏ phiếu cho một lựa chọn.
// Hàm này sử dụng transaction để đảm bảo tính toàn vẹn dữ liệu.
func (r *postPollRepo) VotePoll(ctx context.Context, pollVote *model.PollVote) error {
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Lấy thông tin poll và các vote hiện tại của user
		var post model.Post
		err := r.postCollection.FindOne(sessCtx, bson.M{"_id": pollVote.PostID}).Decode(&post)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return nil, ErrPostNotFound
			}
			return nil, err
		}

		if post.Type != model.PostTypePoll || post.Content == nil || post.Content.Poll == nil {
			return nil, errors.New("post is not a poll")
		}

		userVotes, err := r.GetUserPollVotes(sessCtx, pollVote.PostID, pollVote.UserID)
		if err != nil {
			return nil, err
		}

		// 2. Kiểm tra logic vote
		isAlreadyVoted := false
		for _, v := range userVotes {
			if v.OptionID == pollVote.OptionID {
				isAlreadyVoted = true
				break
			}
		}
		if isAlreadyVoted {
			return nil, ErrPollVoted
		}

		// Nếu poll không cho phép nhiều lựa chọn và user đã vote rồi, thì phải xóa vote cũ.
		if !post.Content.Poll.AllowMultiple && len(userVotes) > 0 {
			err = r.removeVotesInTransaction(sessCtx, pollVote.PostID, pollVote.UserID, userVotes)
			if err != nil {
				return nil, err
			}
		}

		// 3. Thêm vote mới
		pollVote.CreatedAt = time.Now()
		_, err = r.pollVoteCollection.InsertOne(sessCtx, pollVote)
		if err != nil {
			return nil, err
		}

		// 4. Cập nhật số lượng vote trong document Post
		filter := bson.M{"_id": pollVote.PostID, "content.poll.options._id": pollVote.OptionID}
		update := bson.M{
			"$inc": bson.M{
				"content.poll.options.$.votes": 1,
				"content.poll.total_votes":     1,
			},
		}
		_, err = r.postCollection.UpdateOne(sessCtx, filter, update)
		return nil, err
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

// RemovePollVote xóa tất cả các phiếu bầu của một người dùng cho một bài đăng.
func (r *postPollRepo) RemovePollVote(ctx context.Context, postID, userID primitive.ObjectID) error {
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		userVotes, err := r.GetUserPollVotes(sessCtx, postID, userID)
		if err != nil {
			return nil, err
		}
		if len(userVotes) == 0 {
			return nil, nil // Không có gì để xóa
		}

		return nil, r.removeVotesInTransaction(sessCtx, postID, userID, userVotes)
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

// removeVotesInTransaction là hàm helper để xóa vote và cập nhật số lượng bên trong một transaction.
func (r *postPollRepo) removeVotesInTransaction(sessCtx mongo.SessionContext, postID, userID primitive.ObjectID, votes []*model.PollVote) error {
	// 1. Cập nhật (giảm) số vote trong document Post
	for _, vote := range votes {
		filter := bson.M{"_id": postID, "content.poll.options._id": vote.OptionID}
		update := bson.M{
			"$inc": bson.M{
				"content.poll.options.$.votes": -1,
				"content.poll.total_votes":     -1,
			},
		}
		_, err := r.postCollection.UpdateOne(sessCtx, filter, update)
		if err != nil {
			return err // Hủy transaction nếu có lỗi
		}
	}

	// 2. Xóa các bản ghi vote trong collection poll_votes
	_, err := r.pollVoteCollection.DeleteMany(sessCtx, bson.M{"post_id": postID, "user_id": userID})
	return err
}
