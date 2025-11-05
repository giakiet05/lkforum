package service

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrPermissionDenied = errors.New("user does not have permission to perform this action")
	ErrInvalidInput     = errors.New("invalid input provided")
	ErrPostNotFound     = repo.ErrPostNotFound
	ErrPollCannotEdit   = repo.ErrPollCannotEdit
)

type PostService interface {
	CreatePost(ctx context.Context, userID primitive.ObjectID, req *dto.CreatePostRequest) (*dto.PostResponse, error)
	GetPostByID(ctx context.Context, postID, userID primitive.ObjectID) (*dto.PostResponse, error)
	UpdatePost(ctx context.Context, postID, userID primitive.ObjectID, req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	DeletePost(ctx context.Context, postID, userID primitive.ObjectID) error
	GetPosts(ctx context.Context, userID primitive.ObjectID, query *dto.GetPostsQuery) (*PaginatedPostsResponse, error)
	VoteOnPost(ctx context.Context, userID, postID primitive.ObjectID, voteValue bool) (*dto.VotesCountResponse, error)
	RemovePostVote(ctx context.Context, userID, postID primitive.ObjectID) (*dto.VotesCountResponse, error)
	VoteOnPoll(ctx context.Context, userID, postID, optionID primitive.ObjectID) (*dto.PollResponse, error)
	RemovePollVote(ctx context.Context, userID, postID primitive.ObjectID) (*dto.PollResponse, error)
	AddImagesToPost(ctx context.Context, userID, postID primitive.ObjectID, req *dto.AddImageRequest) ([]dto.ImageResponse, error)
	RemoveImagesFromPost(ctx context.Context, userID, postID primitive.ObjectID, req *dto.RemoveImageRequest) error
	UpdatePollDetails(ctx context.Context, postID, userID primitive.ObjectID, req *dto.UpdatePollRequest) (*dto.PollResponse, error)
	AddPollOptions(ctx context.Context, userID, postID primitive.ObjectID, req *dto.AddPollOptionRequest) (*dto.PollResponse, error)
	UpdatePollOption(ctx context.Context, userID, postID, optionID primitive.ObjectID, newText string) (*dto.PollResponse, error)
	RemovePollOptions(ctx context.Context, userID, postID primitive.ObjectID, req *dto.RemovePollOptionRequest) (*dto.PollResponse, error)
	BookmarkPost(ctx context.Context, userID, postID primitive.ObjectID) error
	RemoveBookmark(ctx context.Context, userID, postID primitive.ObjectID) error
}

type postService struct {
	postRepo      repo.PostRepo
	postVoteRepo  repo.PostVoteRepo
	postPollRepo  repo.PostPollRepo
	postImageRepo repo.PostImageRepo
	bus           *bus.EventBus
}

func NewPostService(
	postRepo repo.PostRepo,
	postVoteRepo repo.PostVoteRepo,
	postPollRepo repo.PostPollRepo,
	postImageRepo repo.PostImageRepo,
	bus *bus.EventBus,
) PostService {
	return &postService{
		postRepo:      postRepo,
		postVoteRepo:  postVoteRepo,
		postPollRepo:  postPollRepo,
		postImageRepo: postImageRepo,
		bus:           bus,
	}
}

type PaginatedPostsResponse struct {
	CurrentPage int                 `json:"currentPage"`
	TotalPages  int64               `json:"totalPages"`
	TotalPosts  int64               `json:"totalPosts"`
	HasNextPage bool                `json:"hasNextPage"`
	Posts       []*dto.PostResponse `json:"posts"`
}

func (s *postService) CreatePost(ctx context.Context, userID primitive.ObjectID, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	postModel, err := mapCreateRequestToPostModel(req, userID)
	if err != nil {
		return nil, ErrInvalidInput
	}

	createdPost, err := s.postRepo.Create(ctx, postModel)
	if err != nil {
		return nil, err
	}

	// Publish event for reputation system
	s.bus.Publish(bus.PostCreatedEvent{AuthorID: userID.Hex()})

	return mapPostModelToResponse(createdPost, "", nil), nil
}

func (s *postService) VoteOnPost(ctx context.Context, userID, postID primitive.ObjectID, voteValue bool) (*dto.VotesCountResponse, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Prevent users from voting on their own posts
	if post.AuthorID == userID {
		return mapVotesToResponse(post.VotesCount), nil
	}

	vote := &model.Vote{
		UserID:     userID,
		TargetID:   postID,
		TargetType: model.VoteTargetPost,
		Value:      voteValue,
	}
	if err := s.postVoteRepo.Vote(ctx, vote); err != nil {
		return nil, err
	}

	// Publish event for reputation and notification systems
	if voteValue {
		s.bus.Publish(bus.PostUpvotedEvent{
			AuthorID: post.AuthorID.Hex(),
			VoterID:  userID.Hex(),
			PostID:   postID.Hex(),
		})
	} else {
		s.bus.Publish(bus.PostDownvotedEvent{
			AuthorID: post.AuthorID.Hex(),
			VoterID:  userID.Hex(),
			PostID:   postID.Hex(),
		})
	}

	updatedPost, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	return mapVotesToResponse(updatedPost.VotesCount), nil
}

// ... other methods ...

func (s *postService) GetPostByID(ctx context.Context, postID, userID primitive.ObjectID) (*dto.PostResponse, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

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
	filter := s.buildFilter(query)
	totalPosts, err := s.postRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	if totalPosts == 0 {
		return &PaginatedPostsResponse{
			CurrentPage: query.Page,
			TotalPosts:  0,
			TotalPages:  0,
			HasNextPage: false,
			Posts:       []*dto.PostResponse{},
		}, nil
	}

	findOptions := s.buildFindOptions(query)
	posts, err := s.postRepo.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	postIDs := make([]primitive.ObjectID, len(posts))
	var pollPostIDs []primitive.ObjectID
	for i, p := range posts {
		postIDs[i] = p.ID
		if p.Type == model.PostTypePoll {
			pollPostIDs = append(pollPostIDs, p.ID)
		}
	}

	userVotes := make(map[primitive.ObjectID]string)
	if !userID.IsZero() {
		userVotes, err = s.postVoteRepo.FindUserVotesOnPosts(ctx, userID, postIDs)
		if err != nil {
			// Ghi log lỗi nhưng không dừng request
		}
	}

	userPollVotesMap := make(map[primitive.ObjectID][]*model.PollVote)
	if len(pollPostIDs) > 0 && !userID.IsZero() {
		userPollVotesMap, err = s.postPollRepo.FindUserVotesOnPolls(ctx, userID, pollPostIDs)
		if err != nil {
			// Ghi log lỗi nhưng không dừng request
		}
	}

	responses := make([]*dto.PostResponse, len(posts))
	for i, post := range posts {
		userVoteStr := userVotes[post.ID]
		userPollVotes := userPollVotesMap[post.ID]
		responses[i] = mapPostModelToResponse(post, userVoteStr, userPollVotes)
	}

	totalPages := int64(math.Ceil(float64(totalPosts) / float64(query.Limit)))

	return &PaginatedPostsResponse{
		CurrentPage: query.Page,
		TotalPosts:  totalPosts,
		TotalPages:  totalPages,
		HasNextPage: query.Page < int(totalPages),
		Posts:       responses,
	}, nil
}

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

	if query.Sort == "top" || query.Sort == "controversial" {
		if timeFilter := createTimeFilter(query.TimeFrame); timeFilter != nil {
			for key, value := range timeFilter {
				filter[key] = value
			}
		}
	}

	return filter
}

func (s *postService) buildFindOptions(query *dto.GetPostsQuery) *options.FindOptions {
	opts := options.Find()

	limit := query.Limit
	offset := (query.Page - 1) * limit
	opts.SetLimit(int64(limit))
	opts.SetSkip(int64(offset))

	var sortDoc bson.D
	switch query.Sort {
	case "top":
		sortDoc = bson.D{{Key: "votes_count.up", Value: -1}, {Key: "created_at", Value: -1}}
	case "new":
		fallthrough
	default:
		sortDoc = bson.D{{Key: "created_at", Value: -1}}
	}
	opts.SetSort(sortDoc)

	return opts
}

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

func (s *postService) UpdatePost(ctx context.Context, postID, userID primitive.ObjectID, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post.AuthorID != userID {
		return nil, ErrPermissionDenied
	}

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

func (s *postService) RemovePostVote(ctx context.Context, userID, postID primitive.ObjectID) (*dto.VotesCountResponse, error) {
	if err := s.postVoteRepo.RemoveVote(ctx, postID, userID); err != nil {
		return nil, err
	}

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post.VotesCount == nil {
		return &dto.VotesCountResponse{Up: 0, Down: 0, Score: 0}, nil
	}

	response := &dto.VotesCountResponse{
		Up:    post.VotesCount.Up,
		Down:  post.VotesCount.Down,
		Score: post.VotesCount.Up - post.VotesCount.Down,
	}

	return response, nil
}

func (s *postService) VoteOnPoll(ctx context.Context, userID, postID, optionID primitive.ObjectID) (*dto.PollResponse, error) {
	pollVote := &model.PollVote{
		PostID:   postID,
		UserID:   userID,
		OptionID: optionID,
	}
	if err := s.postPollRepo.VotePoll(ctx, pollVote); err != nil {
		return nil, err
	}

	updatedPost, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	userPollVotes, _ := s.postPollRepo.GetUserPollVotes(ctx, postID, userID)

	return mapPollToResponse(updatedPost.Content.Poll, userPollVotes), nil
}

func (s *postService) RemovePollVote(ctx context.Context, userID, postID primitive.ObjectID) (*dto.PollResponse, error) {
	if err := s.postPollRepo.RemovePollVote(ctx, postID, userID); err != nil {
		return nil, err
	}

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	if post.Type != model.PostTypePoll || post.Content == nil || post.Content.Poll == nil {
		return nil, errors.New("post is not a valid poll")
	}

	response := mapPollToResponse(post.Content.Poll, []*model.PollVote{})

	return response, nil
}

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
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID != userID {
		return nil, ErrPermissionDenied
	}

	if err := s.postPollRepo.UpdatePollOption(ctx, postID, optionID, newText); err != nil {
		return nil, err
	}

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

func (s *postService) BookmarkPost(ctx context.Context, userID, postID primitive.ObjectID) error {
	return errors.New("bookmark not implemented")
}

func (s *postService) RemoveBookmark(ctx context.Context, userID, postID primitive.ObjectID) error {
	return errors.New("bookmark not implemented")
}

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
