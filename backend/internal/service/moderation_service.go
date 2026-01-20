package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/platform/gemini"
	"github.com/giakiet05/lkforum/internal/repo"
	"go.mongodb.org/mongo-driver/bson"
)

type ModerationService interface {
	Start()
	ModeratePost(ctx context.Context, postID string) error
	ModerateComment(ctx context.Context, commentID string) error
}

type moderationService struct {
	postRepo      repo.PostRepo
	commentRepo   repo.CommentRepo
	userRepo      repo.UserRepo
	communityRepo repo.CommunityRepo
	geminiClient  *gemini.GeminiClient
	bus           bus.EventBus
	config        *config.GeminiConfig
}

func NewModerationService(
	postRepo repo.PostRepo,
	commentRepo repo.CommentRepo,
	userRepo repo.UserRepo,
	communityRepo repo.CommunityRepo,
	geminiClient *gemini.GeminiClient,
	bus bus.EventBus,
	config *config.GeminiConfig,
) ModerationService {
	return &moderationService{
		postRepo:      postRepo,
		commentRepo:   commentRepo,
		userRepo:      userRepo,
		communityRepo: communityRepo,
		geminiClient:  geminiClient,
		bus:           bus,
		config:        config,
	}
}

func (s *moderationService) Start() {
	eventChannel := make(bus.EventListener, 100)

	s.bus.Subscribe(bus.TopicPostCreated, eventChannel)
	s.bus.Subscribe(bus.TopicPostUpdated, eventChannel)
	s.bus.Subscribe(bus.TopicCommentCreated, eventChannel)

	log.Println("ModerationService started and subscribed to events.")

	go s.processEvents(eventChannel)
}

func (s *moderationService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicPostCreated:
			s.handlePostCreated(event)
		case bus.TopicPostUpdated:
			s.handlePostUpdated(event)
		case bus.TopicCommentCreated:
			s.handleCommentCreated(event)
		}
	}
}

func (s *moderationService) handlePostCreated(event bus.Event) {
	payload := event.Payload()
	postID, ok := payload["post_id"].(string)
	if !ok {
		postID, ok = payload["PostID"].(string)
		if !ok {
			log.Printf("Invalid PostCreatedEvent payload: %v", payload)
			return
		}
	}

	log.Printf("📬 ModerationService received PostCreated event for post: %s", postID)

	// Run moderation in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.ModeratePost(ctx, postID); err != nil {
			log.Printf("❌ Post moderation failed for %s: %v", postID, err)
		}
	}()
}

func (s *moderationService) handlePostUpdated(event bus.Event) {
	payload := event.Payload()
	postID, ok := payload["post_id"].(string)
	if !ok {
		postID, ok = payload["PostID"].(string)
		if !ok {
			log.Printf("Invalid PostUpdatedEvent payload: %v", payload)
			return
		}
	}

	// Run moderation in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Get post to check moderation status
		post, err := s.postRepo.GetByID(ctx, postID)
		if err != nil {
			log.Printf("Failed to get post %s: %v", postID, err)
			return
		}

		// Only moderate if status is pending (normal users)
		// Admin/moderator edits won't have pending status
		if post.ModerationStatus != model.ModerationPending {
			log.Printf("Skipping moderation for post %s (status: %s)", postID, post.ModerationStatus)
			return
		}

		if err := s.ModeratePost(ctx, postID); err != nil {
			log.Printf("Post moderation failed for %s: %v", postID, err)
		}
	}()
}

func (s *moderationService) handleCommentCreated(event bus.Event) {
	payload := event.Payload()
	commentID, ok := payload["comment_id"].(string)
	if !ok {
		commentID, ok = payload["CommentID"].(string)
		if !ok {
			log.Printf("Invalid CommentCreatedEvent payload: %v", payload)
			return
		}
	}

	// Run moderation in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.ModerateComment(ctx, commentID); err != nil {
			log.Printf("Comment moderation failed for %s: %v", commentID, err)
		}
	}()
}

func (s *moderationService) ModeratePost(ctx context.Context, postID string) error {
	// Get post
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("failed to get post: %w", err)
	}

	// Skip if already skipped (admin/mod posts)
	if post.ModerationStatus == model.ModerationSkipped {
		log.Printf("Post %s moderation skipped (admin/moderator post)", postID)
		return s.approvePost(ctx, postID, model.ModerationSkipped, nil)
	}

	// Get community to check settings
	community, err := s.communityRepo.GetByID(ctx, post.CommunityID.Hex())
	if err != nil {
		return fmt.Errorf("failed to get community: %w", err)
	}

	log.Printf("🤖 ModeratePost - Post: %s, Community: %s", postID, community.Name)
	log.Printf("   - Current Status: %s", post.ModerationStatus)
	log.Printf("   - PostRequireApproval: %v", community.Setting.PostRequireApproval)

	// If community requires manual approval, skip automated moderation
	if community.Setting.PostRequireApproval {
		log.Printf("⏸️  Post %s requires manual approval in community %s, waiting for moderator", postID, community.Name)
		return nil
	}

	// Community uses AI moderation - proceed with AI check
	log.Printf("🔍 Post %s using AI moderation in community %s", postID, community.Name)

	// Build moderation request
	req := &gemini.ContentCheckRequest{
		Title: post.Title,
	}

	if post.Content != nil {
		req.Text = post.Content.Text

		// Add image URLs
		for _, img := range post.Content.Images {
			req.ImageURLs = append(req.ImageURLs, img.URL)
		}

		// Note: Video moderation would require downloading video or using video API
		// For now, we just note that videos exist
		if len(post.Content.Videos) > 0 {
			log.Printf("Post %s contains %d video(s) - video content not analyzed", postID, len(post.Content.Videos))
		}
	}

	// Call Gemini
	result, err := s.geminiClient.CheckContent(ctx, req)
	if err != nil {
		log.Printf("Gemini API error for post %s: %v", postID, err)
		// Fallback: approve to avoid blocking user
		return s.approvePost(ctx, postID, model.ModerationApproved, nil)
	}

	// Build moderation result
	moderationResult := &model.ModerationResult{
		IsViolation:  result.IsViolation,
		Confidence:   result.Confidence,
		Categories:   result.Categories,
		Reason:       result.Reason,
		CheckedText:  req.Text != "" || req.Title != "",
		CheckedMedia: len(req.ImageURLs) > 0,
	}

	// Apply confidence threshold
	if result.IsViolation && result.Confidence >= s.config.ConfidenceThreshold {
		return s.rejectPost(ctx, postID, post.AuthorID.Hex(), moderationResult)
	}

	return s.approvePost(ctx, postID, model.ModerationApproved, moderationResult)
}

func (s *moderationService) ModerateComment(ctx context.Context, commentID string) error {
	// Get comment
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return fmt.Errorf("failed to get comment: %w", err)
	}

	// Check if moderation should be skipped
	author, _ := s.userRepo.GetByID(ctx, comment.Author.ID.Hex())
	if s.shouldSkipModeration(author) {
		return s.approveComment(ctx, commentID, model.ModerationSkipped, nil)
	}

	// Skip very short comments (likely not harmful)
	if len(strings.TrimSpace(comment.Content)) < 10 {
		return s.approveComment(ctx, commentID, model.ModerationSkipped, nil)
	}

	// Build moderation request
	req := &gemini.ContentCheckRequest{
		Text: comment.Content,
	}

	// Call Gemini
	result, err := s.geminiClient.CheckContent(ctx, req)
	if err != nil {
		log.Printf("Gemini API error for comment %s: %v", commentID, err)
		// Fallback: approve to avoid blocking user
		return s.approveComment(ctx, commentID, model.ModerationApproved, nil)
	}

	// Build moderation result
	moderationResult := &model.ModerationResult{
		IsViolation:  result.IsViolation,
		Confidence:   result.Confidence,
		Categories:   result.Categories,
		Reason:       result.Reason,
		CheckedText:  true,
		CheckedMedia: false,
	}

	// Apply confidence threshold
	if result.IsViolation && result.Confidence >= s.config.ConfidenceThreshold {
		return s.rejectComment(ctx, commentID, comment.Author.ID.Hex(), moderationResult)
	}

	return s.approveComment(ctx, commentID, model.ModerationApproved, moderationResult)
}

func (s *moderationService) shouldSkipModeration(user *model.User) bool {
	if user == nil {
		return false
	}
	// Skip moderation for admins and moderators
	return user.Role == "admin" || user.Role == "moderator"
}

func (s *moderationService) approvePost(ctx context.Context, postID string, status model.ModerationStatus, result *model.ModerationResult) error {
	// Get post to retrieve author ID
	post, err := s.postRepo.GetByID(ctx, postID)
	if err != nil {
		return err
	}

	now := time.Now()
	update := repo.UpdateDocument{
		"$set": bson.M{
			"moderation_status": status,
			"moderated_at":      &now,
		},
	}

	if result != nil {
		update["$set"].(bson.M)["moderation_result"] = result
	}

	err = s.postRepo.UpdateByID(ctx, postID, update)
	if err != nil {
		return err
	}

	// Publish approved event with author ID
	s.bus.Publish(&bus.PostApprovedEvent{
		PostID:   postID,
		AuthorID: post.AuthorID.Hex(),
	})

	log.Printf("Post %s approved (status: %s)", postID, status)
	return nil
}

func (s *moderationService) rejectPost(ctx context.Context, postID, authorID string, result *model.ModerationResult) error {
	now := time.Now()
	update := repo.UpdateDocument{
		"$set": bson.M{
			"moderation_status": model.ModerationRejected,
			"moderation_result": result,
			"moderated_at":      &now,
		},
	}

	err := s.postRepo.UpdateByID(ctx, postID, update)
	if err != nil {
		return err
	}

	// Publish rejected event for notifications
	s.bus.Publish(&bus.PostRejectedEvent{
		PostID:   postID,
		AuthorID: authorID,
		Reason:   result.Reason,
	})

	log.Printf("Post %s rejected: %s (confidence: %.2f)", postID, result.Reason, result.Confidence)
	return nil
}

func (s *moderationService) approveComment(ctx context.Context, commentID string, status model.ModerationStatus, result *model.ModerationResult) error {
	// Get comment to retrieve full details for event
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}

	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"moderation_status": status,
			"moderated_at":      &now,
		},
	}

	if result != nil {
		update["$set"].(bson.M)["moderation_result"] = result
	}

	err = s.commentRepo.UpdateByID(ctx, commentID, update)
	if err != nil {
		return err
	}

	// Get parent author ID if exists
	var parentAuthorID *string
	if comment.ParentID != nil {
		if parentComment, err := s.commentRepo.GetByID(ctx, comment.ParentID.Hex()); err == nil {
			parentAuthorIDStr := parentComment.Author.ID.Hex()
			parentAuthorID = &parentAuthorIDStr
		}
	}

	// Publish approved event with full details
	s.bus.Publish(&bus.CommentApprovedEvent{
		CommentID:      commentID,
		PostID:         comment.PostID.Hex(),
		AuthorID:       comment.Author.ID.Hex(),
		ParentAuthorID: parentAuthorID,
	})

	log.Printf("Comment %s approved (status: %s)", commentID, status)
	return nil
}

func (s *moderationService) rejectComment(ctx context.Context, commentID, authorID string, result *model.ModerationResult) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"moderation_status": model.ModerationRejected,
			"moderation_result": result,
			"moderated_at":      &now,
		},
	}

	err := s.commentRepo.UpdateByID(ctx, commentID, update)
	if err != nil {
		return err
	}

	// Publish rejected event for notifications
	s.bus.Publish(&bus.CommentRejectedEvent{
		CommentID: commentID,
		AuthorID:  authorID,
		Reason:    result.Reason,
	})

	log.Printf("Comment %s rejected: %s (confidence: %.2f)", commentID, result.Reason, result.Confidence)
	return nil
}
