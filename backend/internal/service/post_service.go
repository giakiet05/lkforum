// internal/service/post_service.go

package service

import (
	"context"
	"errors"
	"time"

	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Định nghĩa Lỗi ---
var (
	ErrPermissionDenied = errors.New("user does not have permission to perform this action")
	ErrInvalidInput     = errors.New("invalid input provided")
	ErrPostNotFound     = repo.ErrPostNotFound
	ErrPollCannotEdit   = repo.ErrPollCannotEdit
)

// --- PostService Interface (Đầy Đủ) ---
type PostService interface {
	// CRUD cơ bản
	CreatePost(ctx context.Context, userID primitive.ObjectID, req *dto.CreatePostRequest) (*dto.PostResponse, error)
	GetPostByID(ctx context.Context, postID, userID primitive.ObjectID) (*dto.PostResponse, error)
	UpdatePost(ctx context.Context, postID, userID primitive.ObjectID, req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	DeletePost(ctx context.Context, postID, userID primitive.ObjectID) error

	// Lấy danh sách (Feed)
	GetPosts(ctx context.Context, userID primitive.ObjectID, query *dto.GetPostsQuery) ([]*dto.PostResponse, error)

	// Tương tác (Vote)
	VoteOnPost(ctx context.Context, userID, postID primitive.ObjectID, voteValue bool) (*dto.VotesCountResponse, error)
	RemovePostVote(ctx context.Context, userID, postID primitive.ObjectID) (*dto.VotesCountResponse, error)
	VoteOnPoll(ctx context.Context, userID, postID, optionID primitive.ObjectID) (*dto.PollResponse, error)
	RemovePollVote(ctx context.Context, userID, postID primitive.ObjectID) (*dto.PollResponse, error)

	// Quản lý Image (chi tiết)
	AddImagesToPost(ctx context.Context, userID, postID primitive.ObjectID, req *dto.AddImageRequest) ([]dto.ImageResponse, error)
	RemoveImagesFromPost(ctx context.Context, userID, postID primitive.ObjectID, req *dto.RemoveImageRequest) error

	// Quản lý Poll (chi tiết)
	UpdatePollDetails(ctx context.Context, postID, userID primitive.ObjectID, req *dto.UpdatePollRequest) (*dto.PollResponse, error)
	AddPollOptions(ctx context.Context, userID, postID primitive.ObjectID, req *dto.AddPollOptionRequest) (*dto.PollResponse, error)
	UpdatePollOption(ctx context.Context, userID, postID, optionID primitive.ObjectID, newText string) (*dto.PollResponse, error)
	RemovePollOptions(ctx context.Context, userID, postID primitive.ObjectID, req *dto.RemovePollOptionRequest) (*dto.PollResponse, error)

	// Tính năng cho người dùng
	BookmarkPost(ctx context.Context, userID, postID primitive.ObjectID) error
	RemoveBookmark(ctx context.Context, userID, postID primitive.ObjectID) error

	//GetMembersCount(communityID string) (int64, error)
	//increaseMembersCount(communityID string) error
	//decreaseMembersCount(communityID string) error
	//ensureMembersCountExists(communityID string) (string, error)
	//
	//StartRedisToMongoMembershipSync()
	//syncMemberCounts() error
}

// --- postService Implementation ---
type postService struct {
	postRepo      repo.PostRepo
	postVoteRepo  repo.PostVoteRepo
	postPollRepo  repo.PostPollRepo
	postImageRepo repo.PostImageRepo
	communityRepo repo.CommunityRepo
	// --- Placeholder Repositories ---
	// communityRepo repo.CommunityRepo
	// bookmarkRepo  repo.BookmarkRepo
}

func NewPostService(
	postRepo repo.PostRepo,
	postVoteRepo repo.PostVoteRepo,
	postPollRepo repo.PostPollRepo,
	postImageRepo repo.PostImageRepo,
) PostService {
	return &postService{
		postRepo:      postRepo,
		postVoteRepo:  postVoteRepo,
		postPollRepo:  postPollRepo,
		postImageRepo: postImageRepo,
	}
}

// CreatePost tạo bài đăng mới
func (s *postService) CreatePost(ctx context.Context, userID primitive.ObjectID, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	postModel, err := mapCreateRequestToPostModel(req, userID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	createdPost, err := s.postRepo.Create(ctx, postModel)
	if err != nil {
		return nil, err
	}

	// Khi mới tạo, không có thông tin vote của user
	return mapPostModelToResponse(createdPost, nil, nil), nil
}

// GetPostByID lấy chi tiết một bài đăng
func (s *postService) GetPostByID(ctx context.Context, postID, userID primitive.ObjectID) (*dto.PostResponse, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Lấy thông tin phụ để làm giàu DTO
	userVote, _ := s.postVoteRepo.GetUserVote(ctx, postID, userID)
	userPollVotes, _ := s.postPollRepo.GetUserPollVotes(ctx, postID, userID)

	return mapPostModelToResponse(post, userVote, userPollVotes), nil
}

// GetPosts lấy danh sách bài đăng theo nhiều tiêu chí
func (s *postService) GetPosts(ctx context.Context, userID primitive.ObjectID, query *dto.GetPostsQuery) ([]*dto.PostResponse, error) {
	var posts []*model.Post
	var err error

	// TODO: Cần có logic để lấy danh sách postID mà user đã vote để truyền xuống mapper

	if query.CommunityID != "" {
		communityID, _ := primitive.ObjectIDFromHex(query.CommunityID)
		posts, err = s.postRepo.GetByCommunityID(ctx, communityID, query.Sort, query.TimeFrame, query.Limit, (query.Page-1)*query.Limit)
	} else if query.AuthorID != "" {
		authorID, _ := primitive.ObjectIDFromHex(query.AuthorID)
		posts, err = s.postRepo.GetByAuthorID(ctx, authorID, query.Limit, (query.Page-1)*query.Limit)
	} else {
		posts, err = s.postRepo.GetFeed(ctx, query.Sort, query.TimeFrame, query.Limit, (query.Page-1)*query.Limit)
	}

	if err != nil {
		return nil, err
	}

	// Chuyển đổi hàng loạt sang DTO
	responses := make([]*dto.PostResponse, len(posts))
	for i, post := range posts {
		// Trong danh sách, ta có thể không cần lấy chi tiết vote của user để tối ưu hiệu năng
		responses[i] = mapPostModelToResponse(post, nil, nil)
	}
	return responses, nil
}

// UpdatePost cập nhật bài đăng
func (s *postService) UpdatePost(ctx context.Context, postID, userID primitive.ObjectID, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post.AuthorID != userID {
		return nil, ErrPermissionDenied
	}

	// Cập nhật các trường
	if req.Title != "" {
		post.Title = req.Title
	}

	if post.Content != nil {
		if req.Text != "" {
			post.Content.Text = req.Text
		}
	}

	if err := s.postRepo.Update(ctx, post); err != nil {
		return nil, err
	}

	return mapPostModelToResponse(post, nil, nil), nil
}

// DeletePost xóa một bài đăng
func (s *postService) DeletePost(ctx context.Context, postID, userID primitive.ObjectID) error {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}
	if post.AuthorID != userID {
		return ErrPermissionDenied
	}
	return s.postRepo.SoftDelete(ctx, postID)
}

// VoteOnPost xử lý vote up/down
func (s *postService) VoteOnPost(ctx context.Context, userID, postID primitive.ObjectID, voteValue bool) (*dto.VotesCountResponse, error) {
	vote := &model.Vote{
		UserID:     userID,
		TargetID:   postID,
		TargetType: model.VoteTargetPost,
		Value:      voteValue,
	}
	if err := s.postVoteRepo.Vote(ctx, vote); err != nil {
		return nil, err
	}

	// Lấy lại thông tin post để có số vote mới nhất
	updatedPost, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	return mapVotesToResponse(updatedPost.VotesCount), nil
}
func (s *postService) RemovePostVote(ctx context.Context, userID, postID primitive.ObjectID) (*dto.VotesCountResponse, error) {

	// 1. Gọi PostVoteRepo để xóa vote. Logic transaction được xử lý ở tầng Repo.
	if err := s.postVoteRepo.RemoveVote(ctx, postID, userID); err != nil {
		return nil, err
	}

	// 2. Lấy lại thông tin post để có số vote đã cập nhật
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post.VotesCount == nil {
		return &dto.VotesCountResponse{Up: 0, Down: 0, Score: 0}, nil
	}

	// 3. Tạo và trả về DTO response
	response := &dto.VotesCountResponse{
		Up:    post.VotesCount.Up,
		Down:  post.VotesCount.Down,
		Score: post.VotesCount.Up - post.VotesCount.Down,
	}

	return response, nil
}

// VoteOnPoll xử lý vote cho poll
func (s *postService) VoteOnPoll(ctx context.Context, userID, postID, optionID primitive.ObjectID) (*dto.PollResponse, error) {
	pollVote := &model.PollVote{
		PostID:   postID,
		UserID:   userID,
		OptionID: optionID,
	}
	if err := s.postPollRepo.VotePoll(ctx, pollVote); err != nil {
		return nil, err
	}

	// Lấy lại thông tin post và vote của user để trả về poll response mới nhất
	updatedPost, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	userPollVotes, _ := s.postPollRepo.GetUserPollVotes(ctx, postID, userID)

	return mapPollToResponse(updatedPost.Content.Poll, userPollVotes), nil
}

func (s *postService) RemovePollVote(ctx context.Context, userID, postID primitive.ObjectID) (*dto.PollResponse, error) {
	// 1. Gọi PostPollRepo để xóa tất cả các vote của user cho poll này.
	if err := s.postPollRepo.RemovePollVote(ctx, postID, userID); err != nil {
		return nil, err
	}

	// 2. Lấy lại thông tin post để có trạng thái poll đã cập nhật
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// 3. Kiểm tra xem đây có thực sự là một poll hợp lệ không
	if post.Type != model.PostTypePoll || post.Content == nil || post.Content.Poll == nil {
		return nil, errors.New("post is not a valid poll")
	}

	// 4. Tạo và trả về DTO response (sử dụng lại hàm helper `mapPollModelToResponse`)
	response := mapPollToResponse(post.Content.Poll, []*model.PollVote{})

	return response, nil
}

// === Quản lý Image ===
func (s *postService) AddImagesToPost(ctx context.Context, userID, postID primitive.ObjectID, req *dto.AddImageRequest) ([]dto.ImageResponse, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID != userID {
		return nil, ErrPermissionDenied
	}

	newImages := make([]model.Image, len(req.Images))
	for i, imgReq := range req.Images {
		newImages[i] = model.Image{
			ID:         primitive.NewObjectID(),
			URL:        imgReq.URL,
			UploadedAt: time.Now(),
		}
	}

	if err := s.postImageRepo.AddImages(ctx, postID, newImages); err != nil {
		return nil, err
	}

	// Chuyển đổi kết quả trả về
	res := make([]dto.ImageResponse, len(newImages))
	for i, img := range newImages {
		res[i] = dto.ImageResponse{ID: img.ID.Hex(), URL: img.URL}
	}
	return res, nil
}

func (s *postService) RemoveImagesFromPost(ctx context.Context, userID, postID primitive.ObjectID, req *dto.RemoveImageRequest) error {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}
	if post.AuthorID != userID {
		return ErrPermissionDenied
	}

	imageObjectIDs := make([]primitive.ObjectID, len(req.ImageIDs))
	for i, idStr := range req.ImageIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return ErrInvalidInput
		}
		imageObjectIDs[i] = id
	}

	return s.postImageRepo.RemoveImages(ctx, postID, imageObjectIDs)
}

// === Quản lý Poll ===
func (s *postService) UpdatePollDetails(ctx context.Context, postID, userID primitive.ObjectID, req *dto.UpdatePollRequest) (*dto.PollResponse, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID != userID {
		return nil, ErrPermissionDenied
	}

	if err := s.postPollRepo.UpdatePoll(ctx, postID, req.Question, req.ExpiresAt, &req.AllowMultiple); err != nil {
		return nil, err
	}

	return s.getUpdatedPollResponse(ctx, postID, userID)
}

func (s *postService) AddPollOptions(ctx context.Context, userID, postID primitive.ObjectID, req *dto.AddPollOptionRequest) (*dto.PollResponse, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID != userID {
		return nil, ErrPermissionDenied
	}

	newOptions := make([]model.PollOption, len(req.Options))
	for i, text := range req.Options {
		newOptions[i] = model.PollOption{ID: primitive.NewObjectID(), Text: text, Votes: 0}
	}

	if err := s.postPollRepo.AddPollOptions(ctx, postID, newOptions); err != nil {
		return nil, err
	}
	return s.getUpdatedPollResponse(ctx, postID, userID)
}

func (s *postService) UpdatePollOption(ctx context.Context, userID, postID, optionID primitive.ObjectID, newText string) (*dto.PollResponse, error) {
	// 1. Kiểm tra quyền hạn của người dùng
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID != userID {
		return nil, ErrPermissionDenied
	}

	// 2. Gọi đến repository để cập nhật
	if err := s.postPollRepo.UpdatePollOption(ctx, postID, optionID, newText); err != nil {
		return nil, err
	}

	// 3. Lấy lại thông tin poll mới nhất và trả về
	return s.getUpdatedPollResponse(ctx, postID, userID)
}
func (s *postService) RemovePollOptions(ctx context.Context, userID, postID primitive.ObjectID, req *dto.RemovePollOptionRequest) (*dto.PollResponse, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID != userID {
		return nil, ErrPermissionDenied
	}

	optionObjectIDs := make([]primitive.ObjectID, len(req.OptionIDs))
	for i, idStr := range req.OptionIDs {
		id, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			return nil, ErrInvalidInput
		}
		optionObjectIDs[i] = id
	}

	if err := s.postPollRepo.RemovePollOptions(ctx, postID, optionObjectIDs); err != nil {
		return nil, err
	}
	return s.getUpdatedPollResponse(ctx, postID, userID)
}

// === Tính năng người dùng ===
func (s *postService) BookmarkPost(ctx context.Context, userID, postID primitive.ObjectID) error {
	// TODO: Implement logic với bookmarkRepo
	// _, err := s.postRepo.GetByID(ctx, postID) // Kiểm tra post có tồn tại
	// if err != nil { return err }
	// return s.bookmarkRepo.Create(ctx, userID, postID)
	return errors.New("bookmark not implemented")
}

func (s *postService) RemoveBookmark(ctx context.Context, userID, postID primitive.ObjectID) error {
	// TODO: Implement logic với bookmarkRepo
	// return s.bookmarkRepo.Delete(ctx, userID, postID)
	return errors.New("bookmark not implemented")
}

// === Tính năng Admin/Mod ===
func (s *postService) PinPost(ctx context.Context, moderatorID, postID primitive.ObjectID) error {
	// TODO: Dùng communityRepo để kiểm tra quyền moderator
	// post, err := s.postRepo.GetByID(ctx, postID)
	// if err != nil { return err }
	// isMod, err := s.communityRepo.IsModerator(ctx, moderatorID, post.CommunityID)
	// if err != nil { return err }
	// if !isMod { return ErrPermissionDenied }
	// return s.postRepo.SetPinned(ctx, postID, true)
	return errors.New("pin not implemented")
}

func (s *postService) UnpinPost(ctx context.Context, moderatorID, postID primitive.ObjectID) error {
	// TODO: Tương tự PinPost
	// return s.postRepo.SetPinned(ctx, postID, false)
	return errors.New("unpin not implemented")
}

// --- Hàm helper ---
// getUpdatedPollResponse là hàm helper để lấy lại thông tin poll mới nhất và map sang DTO.
func (s *postService) getUpdatedPollResponse(ctx context.Context, postID, userID primitive.ObjectID) (*dto.PollResponse, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	userPollVotes, err := s.postPollRepo.GetUserPollVotes(ctx, postID, userID)
	if err != nil {
		return nil, err
	}
	return mapPollToResponse(post.Content.Poll, userPollVotes), nil
}

// (File mapper.go và các hàm mapper khác giữ nguyên như câu trả lời trước)
