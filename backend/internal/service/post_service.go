package service

import (
	"context"
	"errors"
	"mime/multipart"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/constant"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/platform/cloudinary"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PostService defines the business logic for post-related operations.
type PostService interface {
	CreatePost(userID string, req *dto.CreatePostRequest) (*dto.PostResponse, error)
	GetPostByID(postID string, userID string) (*dto.PostResponse, error)
	GetPosts(userID string, query *dto.GetPostsQuery) (*dto.PaginatedPostsResponse, error)
	UpdatePost(postID string, userID string, req *dto.UpdatePostRequest) (*dto.PostResponse, error)
	DeletePost(postID string, userID string) error

	AddImagesToPost(userID, postID string, form *multipart.Form) ([]*model.Image, error)
	RemoveImagesFromPost(userID, postID string, publicIDs []string) error
	AddVideosToPost(userID, postID string, form *multipart.Form) ([]*model.Video, error)
	RemoveVideosFromPost(userID, postID string, publicIDs []string) error

	VoteOnPost(userID, postID string, voteValue bool) (*dto.VotesCountResponse, error)

	SavePost(userID, postID string) error
	UnsavePost(userID, postID string) error
	GetSavedPosts(userID string, query *dto.GetPostsQuery) (*dto.PaginatedPostsResponse, error)

	ReportPost(reporterID, postID, reason, description string) error

	HidePost(userID, postID string) error

	VoteOnPoll(userID, postID, optionID string) (*dto.PollResponse, error)
	RemovePollVote(userID, postID string) (*dto.PollResponse, error)
	AddPollOptions(userID, postID string, options []string) (*dto.PollResponse, error)
	RemovePollOptions(userID, postID string, optionIDs []string) (*dto.PollResponse, error)
	UpdatePoll(userID, postID string, req *dto.UpdatePollRequest) (*dto.PollResponse, error)
	UpdatePollOption(userID, postID, optionID string, req *dto.UpdatePollOptionRequest) (*dto.PollResponse, error)
}

type postService struct {
	postRepo       repo.PostRepo
	postVoteRepo   repo.PostVoteRepo
	pollVoteRepo   repo.PollVoteRepo
	userRepo       repo.UserRepo
	communityRepo  repo.CommunityRepo
	savedPostRepo  repo.SavedPostRepo
	reportRepo     repo.ReportRepo
	hiddenPostRepo repo.HiddenPostRepo
	bus            bus.EventBus
}

// NewPostService creates a new instance of PostService.
func NewPostService(
	postRepo repo.PostRepo,
	postVoteRepo repo.PostVoteRepo,
	pollVoteRepo repo.PollVoteRepo,
	userRepo repo.UserRepo,
	communityRepo repo.CommunityRepo,
	savedPostRepo repo.SavedPostRepo,
	reportRepo repo.ReportRepo,
	hiddenPostRepo repo.HiddenPostRepo,
	bus bus.EventBus,
) PostService {
	return &postService{
		postRepo:       postRepo,
		postVoteRepo:   postVoteRepo,
		pollVoteRepo:   pollVoteRepo,
		userRepo:       userRepo,
		communityRepo:  communityRepo,
		savedPostRepo:  savedPostRepo,
		reportRepo:     reportRepo,
		hiddenPostRepo: hiddenPostRepo,
		bus:            bus,
	}
}

func (s *postService) CreatePost(userID string, req *dto.CreatePostRequest) (*dto.PostResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	authorID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}
	communityID, err := primitive.ObjectIDFromHex(req.CommunityID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	post := &model.Post{
		AuthorID:      authorID,
		CommunityID:   communityID,
		Title:         req.Title,
		Type:          req.Type,
		Content:       &model.PostContent{Text: req.Text},
		VotesCount:    &model.VotesCount{Up: 1, Down: 0},
		CommentsCount: 0,
		IsDeleted:     false,
		CreatedAt:     time.Now(),
	}

	if req.Type == model.PostTypePoll && req.Poll != nil {
		pollOptions := make([]model.PollOption, len(req.Poll.Options))
		for i, optText := range req.Poll.Options {
			pollOptions[i] = model.PollOption{
				ID:    util.GenerateRandomString(8),
				Text:  optText,
				Votes: 0,
			}
		}
		post.Content.Poll = &model.Poll{
			Question:      req.Poll.Question,
			Options:       pollOptions,
			ExpiresAt:     req.Poll.ExpiresAt,
			AllowMultiple: req.Poll.AllowMultiple,
		}
	}

	createdPost, err := s.postRepo.Create(ctx, post)
	if err != nil {
		return nil, err
	}

	go s.postVoteRepo.Vote(context.Background(), userID, createdPost.ID.Hex(), true)
	s.bus.Publish(bus.PostCreatedEvent{AuthorID: userID})

	return s.GetPostByID(createdPost.ID.Hex(), userID)
}

func (s *postService) GetPostByID(postID string, userID string) (*dto.PostResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	author, _ := s.userRepo.GetByID(ctx, post.AuthorID.Hex())
	community, _ := s.communityRepo.GetByID(ctx, post.CommunityID.Hex())

	var userVoteStr string
	var userPollVoteIDs []string

	if userID != "" {
		vote, _ := s.postVoteRepo.GetUserVote(ctx, userID, postID)
		if vote != nil {
			if vote.Value {
				userVoteStr = "up"
			} else {
				userVoteStr = "down"
			}
		}
		if post.Type == model.PostTypePoll {
			userPollVoteIDs, _ = s.pollVoteRepo.GetUserVoteIDs(ctx, userID, postID)
		}
	}

	return dto.FromPost(post, author, community, userVoteStr, userPollVoteIDs), nil
}

func (s *postService) GetPosts(userID string, query *dto.GetPostsQuery) (*dto.PaginatedPostsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	var hiddenPostIDs []primitive.ObjectID
	if userID != "" {
		userObjID, err := primitive.ObjectIDFromHex(userID)
		if err == nil {
			hiddenPostIDs, _ = s.hiddenPostRepo.GetHiddenPostIDs(ctx, userObjID)
		}
	}

	filter := s.buildFilter(query, hiddenPostIDs)
	findOptions := s.buildFindOptions(query)

	posts, totalPosts, err := s.postRepo.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	if totalPosts == 0 {
		return &dto.PaginatedPostsResponse{Posts: []*dto.PostResponse{}}, nil
	}

	authorIDs, communityIDs := s.extractIDs(posts)
	authors, _ := s.userRepo.GetByIDs(ctx, authorIDs)
	communities, _ := s.communityRepo.GetByIDs(ctx, communityIDs)

	authorsMap := s.mapUsers(authors)
	communitiesMap := s.mapCommunities(communities)

	postIDs := make([]string, len(posts))
	for i, p := range posts {
		postIDs[i] = p.ID.Hex()
	}

	var userVotes map[string]string
	var userPollVotes map[string][]string
	if userID != "" {
		userVotes, _ = s.postVoteRepo.FindUserVotes(ctx, userID, postIDs)
		userPollVotes, _ = s.pollVoteRepo.FindUserVotes(ctx, userID, postIDs)
	}

	responses := dto.FromPosts(posts, authorsMap, communitiesMap, userVotes, userPollVotes)

	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}

	return &dto.PaginatedPostsResponse{
		Posts: responses,
		Pagination: dto.Pagination{
			Page:     query.Page,
			PageSize: limit,
			Total:    totalPosts,
		},
	}, nil
}

func (s *postService) UpdatePost(postID string, userID string, req *dto.UpdatePostRequest) (*dto.PostResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID.Hex() != userID {
		return nil, apperror.ErrForbidden
	}

	update := repo.UpdateDocument{"$set": bson.M{}}
	changed := false
	if req.Title != nil {
		update["$set"].(bson.M)["title"] = *req.Title
		changed = true
	}
	if req.Text != nil {
		update["$set"].(bson.M)["content.text"] = *req.Text
		changed = true
	}

	if !changed {
		return s.GetPostByID(postID, userID)
	}

	update["$set"].(bson.M)["updated_at"] = time.Now()
	if err := s.postRepo.UpdateByID(ctx, postID, update); err != nil {
		return nil, err
	}

	return s.GetPostByID(postID, userID)
}

func (s *postService) DeletePost(postID string, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}
	if post.AuthorID.Hex() != userID {
		return apperror.ErrForbidden
	}
	return s.postRepo.Delete(ctx, postID)
}

func (s *postService) AddImagesToPost(userID, postID string, form *multipart.Form) ([]*model.Image, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID.Hex() != userID {
		return nil, apperror.ErrForbidden
	}

	files := form.File["images"]
	if len(files) == 0 {
		return nil, apperror.ErrBadRequest
	}

	//Use upload images function to handle both single and multiple image uploads
	uploadedImages, err := cloudinary.UploadImages(files)

	if err != nil {
		return nil, err
	}

	if len(uploadedImages) == 0 {
		return nil, apperror.ErrInternal
	}

	update := repo.UpdateDocument{"$push": bson.M{"content.images": bson.M{"$each": uploadedImages}}}
	if err := s.postRepo.UpdateByID(ctx, postID, update); err != nil {
		return nil, err
	}

	return uploadedImages, nil
}

func (s *postService) RemoveImagesFromPost(userID, postID string, publicIDs []string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}
	if post.AuthorID.Hex() != userID {
		return apperror.ErrForbidden
	}

	update := repo.UpdateDocument{"$pull": bson.M{"content.images": bson.M{"public_id": bson.M{"$in": publicIDs}}}}
	if err := s.postRepo.UpdateByID(ctx, postID, update); err != nil {
		return err
	}

	for _, pid := range publicIDs {
		go cloudinary.Delete(pid)
	}

	return nil
}

func (s *postService) AddVideosToPost(userID, postID string, form *multipart.Form) ([]*model.Video, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID.Hex() != userID {
		return nil, apperror.ErrForbidden
	}

	files := form.File["videos"]
	if len(files) == 0 {
		return nil, apperror.ErrBadRequest
	}

	uploadedVideos, err := cloudinary.UploadVideos(files)
	if err != nil {
		return nil, err
	}
	if len(uploadedVideos) == 0 {
		return nil, apperror.ErrInternal
	}

	update := repo.UpdateDocument{"$push": bson.M{"content.videos": bson.M{"$each": uploadedVideos}}, "$set": bson.M{"type": model.PostTypeVideo}}
	if err := s.postRepo.UpdateByID(ctx, postID, update); err != nil {
		// Optionally, try to delete the just-uploaded videos from Cloudinary
		for _, v := range uploadedVideos {
			go cloudinary.Delete(v.PublicID)
		}
		return nil, err
	}

	return uploadedVideos, nil
}

func (s *postService) RemoveVideosFromPost(userID, postID string, publicIDs []string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}
	if post.AuthorID.Hex() != userID {
		return apperror.ErrForbidden
	}

	update := repo.UpdateDocument{"$pull": bson.M{"content.videos": bson.M{"public_id": bson.M{"$in": publicIDs}}}}
	if err := s.postRepo.UpdateByID(ctx, postID, update); err != nil {
		return err
	}

	for _, pid := range publicIDs {
		go cloudinary.Delete(pid)
	}

	return nil
}

func (s *postService) VoteOnPost(userID, postID string, voteValue bool) (*dto.VotesCountResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID.Hex() == userID {
		return nil, apperror.ErrForbidden
	}

	prevVote, _ := s.postVoteRepo.GetUserVote(ctx, userID, postID)

	if err := s.postVoteRepo.Vote(ctx, userID, postID, voteValue); err != nil {
		return nil, err
	}

	s.publishVoteEvents(post.AuthorID.Hex(), userID, postID, prevVote, voteValue)

	updatedPost, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	return &dto.VotesCountResponse{Up: updatedPost.VotesCount.Up, Down: updatedPost.VotesCount.Down, Score: updatedPost.VotesCount.Up - updatedPost.VotesCount.Down}, nil
}

func (s *postService) SavePost(userID, postID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperror.ErrInvalidID
	}
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	// Check if post exists
	_, err = s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	return s.savedPostRepo.Save(ctx, userObjID, postObjID)
}

func (s *postService) UnsavePost(userID, postID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperror.ErrInvalidID
	}
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	err = s.savedPostRepo.Unsave(ctx, userObjID, postObjID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apperror.ErrPostNotFound
		}
		return err
	}
	return nil
}

func (s *postService) GetSavedPosts(userID string, query *dto.GetPostsQuery) (*dto.PaginatedPostsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	findOptions := s.buildFindOptions(query)
	// Sort by when the post was saved, not when it was created
	findOptions.Sort = map[string]int{"saved_at": -1}

	savedPostEntries, total, err := s.savedPostRepo.GetByUserID(ctx, userObjID, findOptions)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &dto.PaginatedPostsResponse{Posts: []*dto.PostResponse{}}, nil
	}

	postIDs := make([]string, len(savedPostEntries))
	for i, entry := range savedPostEntries {
		postIDs[i] = entry.PostID.Hex()
	}

	posts, err := s.postRepo.GetByIDs(ctx, postIDs)
	if err != nil {
		return nil, err
	}

	authorIDs, communityIDs := s.extractIDs(posts)
	authors, _ := s.userRepo.GetByIDs(ctx, authorIDs)
	communities, _ := s.communityRepo.GetByIDs(ctx, communityIDs)

	authorsMap := s.mapUsers(authors)
	communitiesMap := s.mapCommunities(communities)

	var userVotes map[string]string
	var userPollVotes map[string][]string
	if userID != "" {
		userVotes, _ = s.postVoteRepo.FindUserVotes(ctx, userID, postIDs)
		userPollVotes, _ = s.pollVoteRepo.FindUserVotes(ctx, userID, postIDs)
	}

	responses := dto.FromPosts(posts, authorsMap, communitiesMap, userVotes, userPollVotes)

	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}

	return &dto.PaginatedPostsResponse{
		Posts: responses,
		Pagination: dto.Pagination{
			Page:     query.Page,
			PageSize: limit,
			Total:    total,
		},
	}, nil
}

func (s *postService) ReportPost(reporterID, postID, reason, description string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	reporterObjID, err := primitive.ObjectIDFromHex(reporterID)
	if err != nil {
		return apperror.ErrInvalidID
	}
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	// Check if post exists
	_, err = s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	report := &model.Report{
		ReporterID:  reporterObjID,
		TargetID:    postObjID,
		TargetType:  model.ReportTypePost,
		Reason:      reason,
		Description: &description,
		CreatedAt:   time.Now(),
	}

	return s.reportRepo.Create(ctx, report)
}

func (s *postService) HidePost(userID, postID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperror.ErrInvalidID
	}
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	// Check if post exists
	_, err = s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	return s.hiddenPostRepo.Hide(ctx, userObjID, postObjID)
}

func (s *postService) VoteOnPoll(userID, postID, optionID string) (*dto.PollResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if err := s.pollVoteRepo.Vote(ctx, userID, postID, optionID); err != nil {
		return nil, err
	}

	return s.getPollResponse(ctx, postID, userID)
}

func (s *postService) RemovePollVote(userID, postID string) (*dto.PollResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if err := s.pollVoteRepo.RemoveVote(ctx, userID, postID); err != nil {
		return nil, err
	}

	return s.getPollResponse(ctx, postID, userID)
}

func (s *postService) AddPollOptions(userID, postID string, options []string) (*dto.PollResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID.Hex() != userID {
		return nil, apperror.ErrForbidden
	}

	newOptions := make([]model.PollOption, len(options))
	for i, text := range options {
		newOptions[i] = model.PollOption{ID: util.GenerateRandomString(8), Text: text}
	}

	update := repo.UpdateDocument{"$push": bson.M{"content.poll.options": bson.M{"$each": newOptions}}}
	if err := s.postRepo.UpdateByID(ctx, postID, update); err != nil {
		return nil, err
	}

	return s.getPollResponse(ctx, postID, userID)
}

func (s *postService) RemovePollOptions(userID, postID string, optionIDs []string) (*dto.PollResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID.Hex() != userID {
		return nil, apperror.ErrForbidden
	}

	update := repo.UpdateDocument{"$pull": bson.M{"content.poll.options": bson.M{"id": bson.M{"$in": optionIDs}}}}
	if err := s.postRepo.UpdateByID(ctx, postID, update); err != nil {
		return nil, err
	}

	return s.getPollResponse(ctx, postID, userID)
}

func (s *postService) UpdatePoll(userID, postID string, req *dto.UpdatePollRequest) (*dto.PollResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID.Hex() != userID {
		return nil, apperror.ErrForbidden
	}

	updateData := bson.M{}
	if req.Question != nil {
		updateData["content.poll.question"] = *req.Question
	}
	if req.ExpiresAt != nil {
		updateData["content.poll.expires_at"] = *req.ExpiresAt
	}
	if req.AllowMultiple != nil {
		updateData["content.poll.allow_multiple"] = *req.AllowMultiple
	}

	if len(updateData) == 0 {
		return s.getPollResponse(ctx, postID, userID)
	}

	update := repo.UpdateDocument{"$set": updateData}
	if err := s.postRepo.UpdateByID(ctx, postID, update); err != nil {
		return nil, err
	}

	return s.getPollResponse(ctx, postID, userID)
}

func (s *postService) UpdatePollOption(userID, postID, optionID string, req *dto.UpdatePollOptionRequest) (*dto.PollResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	if post.AuthorID.Hex() != userID {
		return nil, apperror.ErrForbidden
	}

	filter := repo.Filter{"_id": post.ID, "content.poll.options.id": optionID}
	update := repo.UpdateDocument{"$set": bson.M{"content.poll.options.$.text": req.Text}}

	if err := s.postRepo.Update(ctx, filter, update); err != nil {
		return nil, err
	}

	return s.getPollResponse(ctx, postID, userID)
}

func (s *postService) getPollResponse(ctx context.Context, postID, userID string) (*dto.PollResponse, error) {
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}
	userPollVoteIDs, _ := s.pollVoteRepo.GetUserVoteIDs(ctx, userID, postID)
	return dto.FromPoll(post.Content.Poll, userPollVoteIDs), nil
}

func (s *postService) publishVoteEvents(authorID, voterID, postID string, prevVote *model.Vote, newVoteValue bool) {
	if prevVote == nil {
		if newVoteValue {
			s.bus.Publish(bus.PostUpvotedEvent{AuthorID: authorID, VoterID: voterID, PostID: postID})
		} else {
			s.bus.Publish(bus.PostDownvotedEvent{AuthorID: authorID, VoterID: voterID, PostID: postID})
		}
	} else if prevVote.Value != newVoteValue {
		if newVoteValue {
			s.bus.Publish(bus.PostDownvoteRemovedEvent{AuthorID: authorID, VoterID: voterID, PostID: postID})
			s.bus.Publish(bus.PostUpvotedEvent{AuthorID: authorID, VoterID: voterID, PostID: postID})
		} else {
			s.bus.Publish(bus.PostUpvoteRemovedEvent{AuthorID: authorID, VoterID: voterID, PostID: postID})
			s.bus.Publish(bus.PostDownvotedEvent{AuthorID: authorID, VoterID: voterID, PostID: postID})
		}
	}
}

// --- Helper methods ---

func (s *postService) buildFilter(query *dto.GetPostsQuery, hiddenPostIDs []primitive.ObjectID) repo.Filter {
	filter := repo.Filter{}
	if query.CommunityID != "" {
		if id, err := primitive.ObjectIDFromHex(query.CommunityID); err == nil {
			filter["community_id"] = id
		}
	}
	if query.AuthorID != "" {
		if id, err := primitive.ObjectIDFromHex(query.AuthorID); err == nil {
			filter["author_id"] = id
		}
	}
	if query.Type != "" {
		filter["type"] = query.Type
	}

	if len(hiddenPostIDs) > 0 {
		filter["_id"] = bson.M{"$nin": hiddenPostIDs}
	}

	return filter
}

func (s *postService) buildFindOptions(query *dto.GetPostsQuery) *repo.FindOptions {
	page := query.Page
	if page <= 0 {
		page = 1
	}
	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	opts := &repo.FindOptions{
		Skip:  int64((page - 1) * limit),
		Limit: int64(limit),
		Sort:  map[string]int{"created_at": -1},
	}

	switch constant.SortType(query.Sort) {
	case constant.SortTypeTop:
		opts.Sort = map[string]int{"votes_count.up": -1}
	case constant.SortTypeNew:
		opts.Sort = map[string]int{"created_at": -1}
	case constant.SortTypeHot:
		opts.Sort = map[string]int{"created_at": -1} // Simplified hot sort
	}

	return opts
}

func (s *postService) extractIDs(posts []*model.Post) ([]string, []string) {
	authorIDMap := make(map[string]bool)
	communityIDMap := make(map[string]bool)
	for _, post := range posts {
		authorIDMap[post.AuthorID.Hex()] = true
		communityIDMap[post.CommunityID.Hex()] = true
	}

	authorIDs := make([]string, 0, len(authorIDMap))
	for id := range authorIDMap {
		authorIDs = append(authorIDs, id)
	}
	communityIDs := make([]string, 0, len(communityIDMap))
	for id := range communityIDMap {
		communityIDs = append(communityIDs, id)
	}
	return authorIDs, communityIDs
}

func (s *postService) mapUsers(users []*model.User) map[string]*model.User {
	authorsMap := make(map[string]*model.User, len(users))
	for _, u := range users {
		authorsMap[u.ID.Hex()] = u
	}
	return authorsMap
}

func (s *postService) mapCommunities(communities []*model.Community) map[string]*model.Community {
	communitiesMap := make(map[string]*model.Community, len(communities))
	for _, c := range communities {
		communitiesMap[c.ID.Hex()] = c
	}
	return communitiesMap
}
