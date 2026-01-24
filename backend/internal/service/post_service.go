package service

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	GetMyPosts(userID string, query *dto.GetPostsQuery) (*dto.PaginatedPostsResponse, error)
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
	UnhidePost(userID, postID string) error
	GetHiddenPosts(userID string, query *dto.GetPostsQuery) (*dto.PaginatedPostsResponse, error)

	BanPost(postID string, reason *string) error
	UnbanPost(postID string) error
	GetBanPosts(query *dto.GetBanPostsQuery, requesterID string) (*dto.PaginatedPostsResponse, error)

	VoteOnPoll(userID, postID, optionID string) (*dto.PollResponse, error)
	RemovePollVote(userID, postID string) (*dto.PollResponse, error)
	AddPollOptions(userID, postID string, options []string) (*dto.PollResponse, error)
	RemovePollOptions(userID, postID string, optionIDs []string) (*dto.PollResponse, error)
	UpdatePoll(userID, postID string, req *dto.UpdatePollRequest) (*dto.PollResponse, error)
	UpdatePollOption(userID, postID, optionID string, req *dto.UpdatePollOptionRequest) (*dto.PollResponse, error)
}

type postService struct {
	postRepo       repo.PostRepo
	voteService    VoteService
	pollVoteRepo   repo.PollVoteRepo
	userRepo       repo.UserRepo
	communityRepo  repo.CommunityRepo
	membershipRepo repo.MembershipRepo
	savedPostRepo  repo.SavedPostRepo
	reportRepo     repo.ReportRepo
	bus            bus.EventBus
}

// NewPostService creates a new instance of PostService.
func NewPostService(
	postRepo repo.PostRepo,
	voteService VoteService,
	pollVoteRepo repo.PollVoteRepo,
	userRepo repo.UserRepo,
	communityRepo repo.CommunityRepo,
	membershipRepo repo.MembershipRepo,
	savedPostRepo repo.SavedPostRepo,
	reportRepo repo.ReportRepo,
	bus bus.EventBus,
) PostService {
	return &postService{
		postRepo:       postRepo,
		voteService:    voteService,
		pollVoteRepo:   pollVoteRepo,
		userRepo:       userRepo,
		communityRepo:  communityRepo,
		membershipRepo: membershipRepo,
		savedPostRepo:  savedPostRepo,
		reportRepo:     reportRepo,
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

	// Get community settings and author info
	community, err := s.communityRepo.GetByID(ctx, req.CommunityID)
	if err != nil {
		return nil, err
	}

	author, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if user is banned from this community
	isBanned, err := s.communityRepo.IsUserBanned(ctx, userID, model.Banned, req.CommunityID)
	if err != nil {
		return nil, err
	}
	if isBanned {
		return nil, apperror.ErrUserIsBannedFromCommunity
	}

	// Determine initial moderation status
	var initialStatus model.ModerationStatus

	// Check if user is admin, creator, or moderator
	isAdminOrMod := author.Role == model.AdminRole
	log.Printf("🔍 User %s - Role: %s, IsAdmin: %v", author.Username, author.Role, author.Role == model.AdminRole)

	if !isAdminOrMod {
		// Check if user is creator of this community
		if community.CreateByID.Hex() == userID {
			isAdminOrMod = true
			log.Printf("✅ User is creator of community %s", community.Name)
		}
	}

	if !isAdminOrMod {
		// Check if user is moderator in this community
		for _, mod := range community.Moderators {
			if mod.UserID.Hex() == userID {
				isAdminOrMod = true
				log.Printf("✅ User is moderator in community %s", community.Name)
				break
			}
		}
	}

	log.Printf("📋 Community %s - PostRequireApproval: %v", community.Name, community.Setting.PostRequireApproval)
	log.Printf("👤 User %s - IsAdminOrMod: %v", author.Username, isAdminOrMod)

	if isAdminOrMod {
		// Admin/Moderator posts are always approved
		initialStatus = model.ModerationSkipped
		log.Printf("✅ Post will be SKIPPED (admin/mod)")
	} else if community.Setting.PostRequireApproval {
		// Community requires manual approval -> pending for mod review
		initialStatus = model.ModerationPending
		log.Printf("⏳ Post will be PENDING (manual approval)")
	} else {
		// Community uses AI moderation -> pending for AI check
		initialStatus = model.ModerationPending
		log.Printf("🤖 Post will be PENDING (AI check)")
	}

	createdAt := time.Now()

	post := &model.Post{
		AuthorID:         authorID,
		CommunityID:      communityID,
		Title:            req.Title,
		Type:             req.Type,
		Content:          &model.PostContent{Text: req.Text},
		VotesCount:       &model.VotesCount{Up: 0, Down: 0},
		CommentsCount:    0,
		HotScore:         model.CalculateHotScore(0, 0, createdAt), // Initial hot score based on creation time
		IsDeleted:        false,
		IsHidden:         false,
		IsDraft:          false,
		IsBan:            false,
		CreatedAt:        createdAt,
		Tags:             req.Tags,
		ModerationStatus: initialStatus,
	}

	fmt.Printf("📝 Creating post with:\n")
	fmt.Printf("  - Title: %s\n", post.Title)
	fmt.Printf("  - Type: %s\n", post.Type)
	fmt.Printf("  - CommunityID: %s\n", post.CommunityID.Hex())
	fmt.Printf("  - AuthorID: %s\n", post.AuthorID.Hex())
	fmt.Printf("  - IsDeleted: %v\n", post.IsDeleted)
	fmt.Printf("  - IsHidden: %v\n", post.IsHidden)
	fmt.Printf("  - IsDraft: %v\n", post.IsDraft)
	fmt.Printf("  - IsBan: %v\n", post.IsBan)
	fmt.Printf("  - ModerationStatus: %s\n", post.ModerationStatus)

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

	// Increment user's post count
	if err := s.userRepo.IncrementPostCount(ctx, userID, 1); err != nil {
		log.Printf("⚠️ Failed to increment post count for user %s: %v", userID, err)
	}

	log.Printf("✅ Post created in DB - ID: %s, Status: %s, IsDeleted: %v", createdPost.ID.Hex(), createdPost.ModerationStatus, createdPost.IsDeleted)

	// Auto-upvote for post creator (use VoteOnTargetByAuthor to skip moderation check)
	if _, err := s.voteService.VoteOnTargetByAuthor(userID, createdPost.ID.Hex(), model.VoteTargetPost, true); err != nil {
		log.Printf("⚠️ Failed to auto-upvote post %s: %v", createdPost.ID.Hex(), err)
	}

	// Publish event for moderation
	s.bus.Publish(&bus.PostCreatedEvent{
		PostID:   createdPost.ID.Hex(),
		AuthorID: userID,
	})

	return s.GetPostByID(createdPost.ID.Hex(), userID)
}

func (s *postService) GetPostByID(postID string, userID string) (*dto.PostResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// Prevent non-authors from seeing hidden posts
	if post.IsHidden && post.AuthorID.Hex() != userID {
		return nil, apperror.ErrPostNotFound
	}

	author, _ := s.userRepo.GetByID(ctx, post.AuthorID.Hex())
	community, _ := s.communityRepo.GetByID(ctx, post.CommunityID.Hex())

	var userVoteStr string
	var userPollVoteIDs []string

	if userID != "" {
		vote, _ := s.voteService.GetUserVote(userID, postID, model.VoteTargetPost)
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

	// Handle feed type filtering
	if query.FeedType == "home" {
		// Home feed: only posts from communities user joined
		if userID == "" {
			// If not logged in, return empty or redirect to explore
			return &dto.PaginatedPostsResponse{Posts: []*dto.PostResponse{}}, nil
		}

		// Get user's joined communities
		memberships, err := s.membershipRepo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}

		if len(memberships) == 0 {
			// User hasn't joined any communities yet
			return &dto.PaginatedPostsResponse{Posts: []*dto.PostResponse{}}, nil
		}

		// Extract community IDs
		var communityIDs []primitive.ObjectID
		for _, m := range memberships {
			communityIDs = append(communityIDs, m.CommunityID)
		}

		// Filter posts from these communities only
		filter := s.buildFilter(query)
		filter["community_id"] = bson.M{"$in": communityIDs}
		findOptions := s.buildFindOptions(query)

		posts, totalPosts, err := s.postRepo.Find(ctx, filter, findOptions)
		if err != nil {
			return nil, err
		}

		if totalPosts == 0 {
			return &dto.PaginatedPostsResponse{Posts: []*dto.PostResponse{}}, nil
		}

		authorIDs, communityIDs2 := s.extractIDs(posts)
		authors, _ := s.userRepo.GetByIDs(ctx, authorIDs)
		communities, _ := s.communityRepo.GetByIDs(ctx, communityIDs2)

		authorsMap := s.mapUsers(authors)
		communitiesMap := s.mapCommunities(communities)

		postIDs := make([]string, len(posts))
		for i, p := range posts {
			postIDs[i] = p.ID.Hex()
		}

		var userVotes map[string]string
		var userPollVotes map[string][]string
		if userID != "" {
			userVotes, _ = s.voteService.FindUserVotes(userID, postIDs, model.VoteTargetPost)
			userPollVotes, _ = s.pollVoteRepo.FindUserVotes(ctx, userID, postIDs)
		}

		postResponses := dto.FromPosts(posts, authorsMap, communitiesMap, userVotes, userPollVotes)

		limit := query.Limit
		if limit <= 0 {
			limit = 10
		}

		return &dto.PaginatedPostsResponse{
			Posts: postResponses,
			Pagination: dto.Pagination{
				Page:     query.Page,
				PageSize: limit,
				Total:    totalPosts,
			},
		}, nil
	} else if query.FeedType == "popular" {
		// Popular feed: sort by score (up - down votes)
		query.Sort = "top"
	}
	// explore and all: no special filtering, show all public posts

	// Check if filtering by community and if it's private
	if query.CommunityID != "" {
		community, err := s.communityRepo.GetByID(ctx, query.CommunityID)
		if err != nil {
			return nil, err
		}

		// Check if community is banned
		if community.IsBanned {
			return nil, apperror.ErrCommunityNotFound // Or create a specific ErrCommunityBanned
		}

		// If community is private, check membership
		if community.Setting.IsPrivate {
			// User must be logged in
			if userID == "" {
				return nil, apperror.ErrForbidden
			}

			// Check if user is a member
			isMember, err := s.membershipRepo.IsMember(ctx, userID, query.CommunityID)
			if err != nil {
				return nil, err
			}
			if !isMember {
				return nil, apperror.ErrForbidden
			}
		}
	}

	filter := s.buildFilter(query)
	findOptions := s.buildFindOptions(query)

	fmt.Printf("GetPosts filter: %+v\n", filter)

	posts, totalPosts, err := s.postRepo.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Found %d posts\n", totalPosts)

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
		userVotes, _ = s.voteService.FindUserVotes(userID, postIDs, model.VoteTargetPost)
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

func (s *postService) GetMyPosts(userID string, query *dto.GetPostsQuery) (*dto.PaginatedPostsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	filter := repo.Filter{
		"author_id": userObjID,
		"is_hidden": bson.M{"$ne": true},
		"is_draft":  bson.M{"$ne": true},
	}
	findOptions := s.buildFindOptions(query)

	posts, total, err := s.postRepo.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	limit := query.Limit
	if limit <= 0 {
		limit = 10
	}
	page := query.Page
	if page <= 0 {
		page = 1
	}

	if total == 0 {
		return &dto.PaginatedPostsResponse{
			Posts: []*dto.PostResponse{},
			Pagination: dto.Pagination{
				Page:     page,
				PageSize: limit,
				Total:    0,
			},
		}, nil
	}

	author, _ := s.userRepo.GetByID(ctx, userID)
	authorsMap := map[string]*model.User{userID: author}

	_, communityIDs := s.extractIDs(posts)
	communities, _ := s.communityRepo.GetByIDs(ctx, communityIDs)
	communitiesMap := s.mapCommunities(communities)

	postIDs := make([]string, len(posts))
	for i, p := range posts {
		postIDs[i] = p.ID.Hex()
	}

	userVotes, _ := s.voteService.FindUserVotes(userID, postIDs, model.VoteTargetPost)
	userPollVotes, _ := s.pollVoteRepo.FindUserVotes(ctx, userID, postIDs)

	responses := dto.FromPostsWithModeration(posts, authorsMap, communitiesMap, userVotes, userPollVotes)

	return &dto.PaginatedPostsResponse{
		Posts: responses,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: limit,
			Total:    total,
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
	if req.Tags != nil { // Check if tags are provided in the request
		update["$set"].(bson.M)["tags"] = req.Tags
		changed = true
	}

	if !changed {
		return s.GetPostByID(postID, userID)
	}

	// Mark post as edited if it was previously approved
	if post.ModerationStatus == model.ModerationApproved {
		update["$set"].(bson.M)["is_edited"] = true
	}

	update["$set"].(bson.M)["updated_at"] = time.Now()
	if err := s.postRepo.UpdateByID(ctx, postID, update); err != nil {
		return nil, err
	}

	// Publish PostUpdatedEvent for moderation
	s.bus.Publish(&bus.PostUpdatedEvent{
		PostID:   postID,
		AuthorID: userID,
	})

	return s.GetPostByID(postID, userID)
}

func (s *postService) DeletePost(postID string, userID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	// Check if user is the post author
	isAuthor := post.AuthorID.Hex() == userID

	// If not the author, check if user has permission (admin/creator/moderator)
	if !isAuthor {
		// Get user to check if admin
		user, err := s.userRepo.GetByID(ctx, userID)
		if err != nil {
			return err
		}

		isAdmin := user.Role == model.AdminRole
		if !isAdmin {
			// Get community to check if user is creator or moderator
			community, err := s.communityRepo.GetByID(ctx, post.CommunityID.Hex())
			if err != nil {
				return err
			}

			// Check if user is community creator
			isCreator := community.CreateByID.Hex() == userID

			// Check if user is moderator
			isModerator := false
			for _, mod := range community.Moderators {
				if mod.UserID.Hex() == userID {
					isModerator = true
					break
				}
			}

			// If not creator and not moderator, deny access
			if !isCreator && !isModerator {
				return apperror.ErrForbidden
			}
		}
	}

	// Decrement user's post count (only for the author)
	if isAuthor {
		if err := s.userRepo.IncrementPostCount(ctx, userID, -1); err != nil {
			log.Printf("⚠️ Failed to decrement post count for user %s: %v", userID, err)
		}
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
	if post.Content == nil || len(post.Content.Videos) == 0 {
		return apperror.NewError(nil, apperror.ErrBadRequest.Code, "post does not contain any videos")
	}

	// Verify that the publicIDs exist in the post's videos
	var foundPublicIDs []string
	for _, reqID := range publicIDs {
		for _, postVideo := range post.Content.Videos {
			if postVideo.PublicID == reqID {
				foundPublicIDs = append(foundPublicIDs, reqID)
				break
			}
		}
	}

	if len(foundPublicIDs) == 0 {
		return apperror.NewError(nil, apperror.ErrBadRequest.Code, "no matching videos found in post to remove")
	}

	update := repo.UpdateDocument{"$pull": bson.M{"content.videos": bson.M{"public_id": bson.M{"$in": foundPublicIDs}}}}
	if err := s.postRepo.UpdateByID(ctx, postID, update); err != nil {
		return err
	}

	for _, pid := range foundPublicIDs {
		go cloudinary.Delete(pid)
	}

	return nil
}

func (s *postService) VoteOnPost(userID, postID string, voteValue bool) (*dto.VotesCountResponse, error) {
	return s.voteService.VoteOnTarget(userID, postID, model.VoteTargetPost, voteValue)
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
		userVotes, _ = s.voteService.FindUserVotes(userID, postIDs, model.VoteTargetPost)
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
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	// Check if the reporter is the author of the post
	if post.AuthorID == reporterObjID {
		return apperror.NewError(nil, apperror.ErrForbidden.Code, "You cannot report your own post")
	}

	// Check if user has already reported this post
	reports, _, err := s.reportRepo.GetFilter(ctx, &reporterID, &postID, model.ReportTypePost, nil, nil, nil, 1, 10)
	if err != nil {
		return err
	}
	if len(reports) > 0 {
		return apperror.ErrAlreadyReported
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

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	if post.AuthorID.Hex() != userID {
		return apperror.ErrForbidden
	}

	update := repo.UpdateDocument{"$set": bson.M{"is_hidden": true}}
	return s.postRepo.UpdateByID(ctx, postID, update)
}

func (s *postService) UnhidePost(userID, postID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	if post.AuthorID.Hex() != userID {
		return apperror.ErrForbidden
	}

	update := repo.UpdateDocument{"$set": bson.M{"is_hidden": false}}
	return s.postRepo.UpdateByID(ctx, postID, update)
}

func (s *postService) BanPost(postID string, reason *string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	_, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return apperror.ErrPostNotFound
	}

	return s.postRepo.BanPost(ctx, postID, reason)
}

func (s *postService) UnbanPost(postID string) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	_, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return apperror.ErrPostNotFound
	}

	return s.postRepo.UnbanPost(ctx, postID)
}

func (s *postService) GetBanPosts(query *dto.GetBanPostsQuery, requesterID string) (*dto.PaginatedPostsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	ok, err := s.communityRepo.IsModerator(ctx, query.CommunityID, requesterID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.ErrForbidden
	}

	posts, total, err := s.postRepo.GetBannedPosts(ctx, query.CommunityID, query.Page, query.PageSize)
	if err != nil {
		return nil, err
	}

	if total == 0 {
		return &dto.PaginatedPostsResponse{Posts: []*dto.PostResponse{}}, nil
	}

	// Since the author is the user, we can fetch their details once
	author, _ := s.userRepo.GetByID(ctx, requesterID)
	authorsMap := map[string]*model.User{requesterID: author}

	// Fetch communities
	_, communityIDs := s.extractIDs(posts)
	communities, _ := s.communityRepo.GetByIDs(ctx, communityIDs)
	communitiesMap := s.mapCommunities(communities)

	// No need to check for votes on one's own hidden posts
	userVotes := make(map[string]string)
	userPollVotes := make(map[string][]string)

	responses := dto.FromPosts(posts, authorsMap, communitiesMap, userVotes, userPollVotes)

	return &dto.PaginatedPostsResponse{
		Posts: responses,
		Pagination: dto.Pagination{
			Page:     query.Page,
			PageSize: query.PageSize,
			Total:    total,
		},
	}, nil
}

func (s *postService) GetHiddenPosts(userID string, query *dto.GetPostsQuery) (*dto.PaginatedPostsResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	filter := repo.Filter{
		"author_id": userObjID,
		"is_hidden": true,
	}
	findOptions := s.buildFindOptions(query)

	posts, total, err := s.postRepo.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &dto.PaginatedPostsResponse{Posts: []*dto.PostResponse{}}, nil
	}

	// Since the author is the user, we can fetch their details once
	author, _ := s.userRepo.GetByID(ctx, userID)
	authorsMap := map[string]*model.User{userID: author}

	// Fetch communities
	_, communityIDs := s.extractIDs(posts)
	communities, _ := s.communityRepo.GetByIDs(ctx, communityIDs)
	communitiesMap := s.mapCommunities(communities)

	// No need to check for votes on one's own hidden posts
	userVotes := make(map[string]string)
	userPollVotes := make(map[string][]string)

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

// --- Helper methods ---

func (s *postService) buildFilter(query *dto.GetPostsQuery) repo.Filter {
	filter := repo.Filter{
		"is_hidden":         bson.M{"$in": []interface{}{false, nil}},
		"is_draft":          bson.M{"$in": []interface{}{false, nil}},
		"is_ban":            bson.M{"$in": []interface{}{false, nil}},
		"moderation_status": bson.M{"$in": []interface{}{model.ModerationApproved, model.ModerationSkipped}},
	}

	fmt.Printf("🔍 Building filter with base conditions:\n")
	fmt.Printf("  - is_hidden: false or nil\n")
	fmt.Printf("  - is_draft: false or nil\n")
	fmt.Printf("  - is_ban: false or nil\n")
	fmt.Printf("  - moderation_status: approved OR skipped\n")

	if query.CommunityID != "" {
		if id, err := primitive.ObjectIDFromHex(query.CommunityID); err == nil {
			filter["community_id"] = id
			fmt.Printf("  - community_id: %s\n", query.CommunityID)
		}
	}
	if query.AuthorID != "" {
		if id, err := primitive.ObjectIDFromHex(query.AuthorID); err == nil {
			filter["author_id"] = id
			fmt.Printf("  - author_id: %s\n", query.AuthorID)
		}
	}
	if query.Type != "" {
		filter["type"] = query.Type
		fmt.Printf("  - type: %s\n", query.Type)
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
		// Top: sort by upvotes count
		opts.Sort = map[string]int{"votes_count.up": -1}
	case constant.SortTypeNew:
		// New: sort by creation time (newest first)
		opts.Sort = map[string]int{"created_at": -1}
	case constant.SortTypeHot:
		// Hot: sort by comment count and recent activity
		opts.Sort = map[string]int{"comment_count": -1, "created_at": -1}
	case constant.SortTypeBest:
		// Best: Use hot_score for ranking (calculated field that combines votes + time decay)
		// This provides variety as newer posts with decent engagement can compete with older popular posts
		opts.Sort = map[string]int{"hot_score": -1, "created_at": -1}
	case constant.SortTypeRising:
		// Rising: recent posts with high engagement (votes + comments)
		opts.Sort = map[string]int{"votes_count.up": -1, "comment_count": -1, "created_at": -1}
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
		// Only include non-banned communities
		if !c.IsBanned {
			communitiesMap[c.ID.Hex()] = c
		}
	}
	return communitiesMap
}
