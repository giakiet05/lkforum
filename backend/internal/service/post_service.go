// internal/service/post_service.go

package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	GetPosts(ctx context.Context, userID primitive.ObjectID, query *dto.GetPostsQuery) (*PaginatedPostsResponse, error)

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

type PaginatedPostsResponse struct {
	CurrentPage int                 `json:"currentPage"`
	TotalPages  int64               `json:"totalPages"`
	TotalPosts  int64               `json:"totalPosts"`
	HasNextPage bool                `json:"hasNextPage"`
	Posts       []*dto.PostResponse `json:"posts"`
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
	return mapPostModelToResponse(createdPost, "", nil), nil
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
	var userVoteStr string
	if userVote != nil {
		if userVote.Value {
			userVoteStr = "up"
		} else {
			userVoteStr = "down"
		}
	}
	return mapPostModelToResponse(post, userVoteStr, userPollVotes), nil
}

func (s *postService) GetPosts(ctx context.Context, userID primitive.ObjectID, query *dto.GetPostsQuery) (*PaginatedPostsResponse, error) {
	// 1. Service xây dựng bộ lọc động dựa trên query DTO.
	filter := s.buildFilter(query)

	// 2. Lấy tổng số bài đăng khớp với bộ lọc để tính toán phân trang.
	totalPosts, err := s.postRepo.Count(ctx, filter)
	if err != nil {
		return nil, err // Nếu không đếm được thì trả về lỗi
	}

	// Nếu không có bài đăng nào, trả về kết quả rỗng ngay lập tức để tiết kiệm chi phí.
	if totalPosts == 0 {
		return &PaginatedPostsResponse{
			CurrentPage: query.Page,
			TotalPosts:  0,
			TotalPages:  0,
			HasNextPage: false,
			Posts:       []*dto.PostResponse{},
		}, nil
	}

	// 3. Service xây dựng các tùy chọn truy vấn (sắp xếp, phân trang).
	findOptions := s.buildFindOptions(query)

	// 4. Gọi hàm Find của Repository để lấy dữ liệu.
	posts, err := s.postRepo.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	// Chuẩn bị các slice ID để truy vấn thông tin vote
	postIDs := make([]primitive.ObjectID, len(posts))
	var pollPostIDs []primitive.ObjectID
	for i, p := range posts {
		postIDs[i] = p.ID
		if p.Type == model.PostTypePoll {
			pollPostIDs = append(pollPostIDs, p.ID)
		}
	}

	// 5. Lấy thông tin upvote/downvote của người dùng.
	userVotes := make(map[primitive.ObjectID]string)
	if !userID.IsZero() {
		userVotes, err = s.postVoteRepo.FindUserVotesOnPosts(ctx, userID, postIDs)
		if err != nil {
			// Ghi log lỗi nhưng không dừng request
		}
	}

	// 6. Lấy thông tin poll vote của người dùng.
	userPollVotesMap := make(map[primitive.ObjectID][]*model.PollVote)
	// Chỉ gọi DB nếu có bài poll và user đã đăng nhập
	if len(pollPostIDs) > 0 && !userID.IsZero() {
		userPollVotesMap, err = s.postPollRepo.FindUserVotesOnPolls(ctx, userID, pollPostIDs)
		if err != nil {
			// Ghi log lỗi nhưng không dừng request
		}
	}

	// 7. Chuyển đổi (map) model sang DTO.
	responses := make([]*dto.PostResponse, len(posts))
	for i, post := range posts {
		userVoteStr := userVotes[post.ID]
		userPollVotes := userPollVotesMap[post.ID]
		responses[i] = mapPostModelToResponse(post, userVoteStr, userPollVotes)
	}

	// 8. Tính toán các thông tin phân trang cuối cùng và trả về.
	totalPages := int64(math.Ceil(float64(totalPosts) / float64(query.Limit)))

	return &PaginatedPostsResponse{
		CurrentPage: query.Page,
		TotalPosts:  totalPosts,
		TotalPages:  totalPages,
		HasNextPage: query.Page < int(totalPages),
		Posts:       responses,
	}, nil
}

// buildFilter là hàm helper để tạo bộ lọc động.
func (s *postService) buildFilter(query *dto.GetPostsQuery) bson.M {
	filter := bson.M{"is_deleted": bson.M{"$ne": true}}

	if query.CommunityID != "" {
		if communityID, err := primitive.ObjectIDFromHex(query.CommunityID); err == nil {
			filter["community_id"] = communityID
		}
	}

	if query.AuthorID != "" {
		if authorID, err := primitive.ObjectIDFromHex(query.AuthorID); err == nil {
			filter["author_id"] = authorID
		}
	}

	if query.Type != "" {
		filter["type"] = query.Type
	}

	// Logic nghiệp vụ: Chỉ áp dụng bộ lọc thời gian cho các kiểu sort "top" hoặc "controversial".
	if query.Sort == "top" || query.Sort == "controversial" {
		if timeFilter := createTimeFilter(query.TimeFrame); timeFilter != nil {
			for key, value := range timeFilter {
				filter[key] = value
			}
		}
	}

	return filter
}

// buildFindOptions là hàm helper để tạo các tùy chọn truy vấn (sắp xếp, phân trang).
func (s *postService) buildFindOptions(query *dto.GetPostsQuery) *options.FindOptions {
	opts := options.Find()

	// Phân trang
	limit := query.Limit
	offset := (query.Page - 1) * limit
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))

	// Sắp xếp
	var sortDoc bson.D
	switch query.Sort {
	case "top":
		sortDoc = bson.D{{Key: "votes_count.up", Value: -1}, {Key: "created_at", Value: -1}}
	// Thêm các case "hot", "controversial" ở đây nếu cần.
	// Để có sort "hot" chính xác, bạn cần dùng Aggregation Pipeline.
	case "new":
		fallthrough
	default:
		sortDoc = bson.D{{Key: "created_at", Value: -1}}
	}
	opts.SetSort(sortDoc)

	return opts
}

// --- HÀM TIỆN ÍCH (UTILITY) ---

// createTimeFilter là một hàm tiện ích tạo ra bộ lọc thời gian.
func createTimeFilter(timeFrame string) bson.M {
	if timeFrame == "" || timeFrame == "all" {
		return nil
	}

	now := time.Now()
	var startTime time.Time

	switch timeFrame {
	case "hour":
		startTime = now.Add(-1 * time.Hour)
	case "day":
		startTime = now.Add(-24 * time.Hour)
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	case "year":
		startTime = now.AddDate(-1, 0, 0)
	default:
		return nil
	}

	return bson.M{"created_at": bson.M{"$gte": startTime}}
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

	return mapPostModelToResponse(post, "", nil), nil
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
