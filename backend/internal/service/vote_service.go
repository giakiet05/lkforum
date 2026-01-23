package service

import (
	"context"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
)

type VoteService interface {
	VoteOnTarget(userID, targetID string, targetType model.VoteTargetType, voteValue bool) (*dto.VotesCountResponse, error)
	VoteOnTargetByAuthor(userID, targetID string, targetType model.VoteTargetType, voteValue bool) (*dto.VotesCountResponse, error)
	GetUserVote(userID, targetID string, targetType model.VoteTargetType) (*model.Vote, error)
	FindUserVotes(userID string, targetIDs []string, targetType model.VoteTargetType) (map[string]string, error)
}

type voteService struct {
	voteRepo    repo.VoteRepo
	postRepo    repo.PostRepo
	commentRepo repo.CommentRepo
	eventBus    bus.EventBus
}

func NewVoteService(
	voteRepo repo.VoteRepo,
	postRepo repo.PostRepo,
	commentRepo repo.CommentRepo,
	eventBus bus.EventBus,
) VoteService {
	return &voteService{
		voteRepo:    voteRepo,
		postRepo:    postRepo,
		commentRepo: commentRepo,
		eventBus:    eventBus,
	}
}

func (s *voteService) VoteOnTarget(userID, targetID string, targetType model.VoteTargetType, voteValue bool) (*dto.VotesCountResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Check moderation status before allowing vote
	if err := s.checkModerationStatus(ctx, targetID, targetType); err != nil {
		return nil, err
	}

	// Get author ID based on target type
	authorID, err := s.getAuthorID(ctx, targetID, targetType)
	if err != nil {
		return nil, err
	}

	// Get previous vote to determine events to publish
	prevVote, _ := s.voteRepo.GetUserVote(ctx, userID, targetID, targetType)

	// Perform vote
	if err := s.voteRepo.Vote(ctx, userID, targetID, targetType, voteValue); err != nil {
		return nil, err
	}

	// Publish events
	s.publishVoteEvents(authorID, userID, targetID, targetType, prevVote, voteValue)

	// Get updated vote counts
	return s.getUpdatedVoteCounts(ctx, targetID, targetType)
}

// VoteOnTargetByAuthor allows the author to vote on their own content without moderation check
// Used for auto-upvote when creating a post
func (s *voteService) VoteOnTargetByAuthor(userID, targetID string, targetType model.VoteTargetType, voteValue bool) (*dto.VotesCountResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Skip moderation check for author's own vote

	// Get author ID based on target type
	authorID, err := s.getAuthorID(ctx, targetID, targetType)
	if err != nil {
		return nil, err
	}

	// Get previous vote to determine events to publish
	prevVote, _ := s.voteRepo.GetUserVote(ctx, userID, targetID, targetType)

	// Perform vote
	if err := s.voteRepo.Vote(ctx, userID, targetID, targetType, voteValue); err != nil {
		return nil, err
	}

	// Publish events
	s.publishVoteEvents(authorID, userID, targetID, targetType, prevVote, voteValue)

	// Get updated vote counts
	return s.getUpdatedVoteCounts(ctx, targetID, targetType)
}

func (s *voteService) GetUserVote(userID, targetID string, targetType model.VoteTargetType) (*model.Vote, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return s.voteRepo.GetUserVote(ctx, userID, targetID, targetType)
}

func (s *voteService) FindUserVotes(userID string, targetIDs []string, targetType model.VoteTargetType) (map[string]string, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	return s.voteRepo.FindUserVotes(ctx, userID, targetIDs, targetType)
}

// checkModerationStatus checks if target is approved before allowing vote
func (s *voteService) checkModerationStatus(ctx context.Context, targetID string, targetType model.VoteTargetType) error {
	switch targetType {
	case model.VoteTargetPost:
		post, err := s.postRepo.GetByID(ctx, targetID)
		if err != nil {
			return err
		}
		if post.ModerationStatus != model.ModerationApproved && post.ModerationStatus != model.ModerationSkipped {
			return apperror.ErrForbidden
		}
	case model.VoteTargetComment:
		_, err := s.commentRepo.GetByID(ctx, targetID)
		if err != nil {
			return err
		}
		// No moderation check - all comments are visible
	}
	return nil
}

// getAuthorID retrieves the author ID based on target type
func (s *voteService) getAuthorID(ctx context.Context, targetID string, targetType model.VoteTargetType) (string, error) {
	switch targetType {
	case model.VoteTargetPost:
		post, err := s.postRepo.GetByID(ctx, targetID)
		if err != nil {
			return "", err
		}
		return post.AuthorID.Hex(), nil
	case model.VoteTargetComment:
		comment, err := s.commentRepo.GetByID(ctx, targetID)
		if err != nil {
			return "", err
		}
		return comment.Author.ID.Hex(), nil
	default:
		return "", nil
	}
}

// getUpdatedVoteCounts retrieves the updated vote counts after voting
func (s *voteService) getUpdatedVoteCounts(ctx context.Context, targetID string, targetType model.VoteTargetType) (*dto.VotesCountResponse, error) {
	var votesCount *model.VotesCount

	switch targetType {
	case model.VoteTargetPost:
		post, err := s.postRepo.GetByID(ctx, targetID)
		if err != nil {
			return nil, err
		}
		votesCount = post.VotesCount
	case model.VoteTargetComment:
		comment, err := s.commentRepo.GetByID(ctx, targetID)
		if err != nil {
			return nil, err
		}
		votesCount = comment.VotesCount
	}

	if votesCount == nil {
		return &dto.VotesCountResponse{Up: 0, Down: 0, Score: 0}, nil
	}

	return &dto.VotesCountResponse{
		Up:    votesCount.Up,
		Down:  votesCount.Down,
		Score: votesCount.Up - votesCount.Down,
	}, nil
}

// publishVoteEvents publishes appropriate events based on vote type and previous vote
func (s *voteService) publishVoteEvents(authorID, voterID, targetID string, targetType model.VoteTargetType, prevVote *model.Vote, newVoteValue bool) {
	if prevVote == nil {
		// New vote
		s.publishNewVoteEvent(authorID, voterID, targetID, targetType, newVoteValue)
	} else if prevVote.Value != newVoteValue {
		// Change vote (up to down or down to up)
		s.publishVoteChangeEvents(authorID, voterID, targetID, targetType, newVoteValue)
	}
	// If prevVote.Value == newVoteValue, it's an un-vote, handled by repo (no event needed here)
}

// publishNewVoteEvent publishes event for a new vote
func (s *voteService) publishNewVoteEvent(authorID, voterID, targetID string, targetType model.VoteTargetType, isUpvote bool) {
	switch targetType {
	case model.VoteTargetPost:
		if isUpvote {
			s.eventBus.Publish(bus.PostUpvotedEvent{AuthorID: authorID, VoterID: voterID, PostID: targetID})
		} else {
			s.eventBus.Publish(bus.PostDownvotedEvent{AuthorID: authorID, VoterID: voterID, PostID: targetID})
		}
	case model.VoteTargetComment:
		if isUpvote {
			s.eventBus.Publish(bus.CommentUpvotedEvent{AuthorID: authorID, VoterID: voterID, CommentID: targetID})
		} else {
			s.eventBus.Publish(bus.CommentDownvotedEvent{AuthorID: authorID, VoterID: voterID, CommentID: targetID})
		}
	}
}

// publishVoteChangeEvents publishes events for changing vote from up to down or vice versa
func (s *voteService) publishVoteChangeEvents(authorID, voterID, targetID string, targetType model.VoteTargetType, newIsUpvote bool) {
	switch targetType {
	case model.VoteTargetPost:
		if newIsUpvote {
			// Changed from down to up
			s.eventBus.Publish(bus.PostDownvoteRemovedEvent{AuthorID: authorID, VoterID: voterID, PostID: targetID})
			s.eventBus.Publish(bus.PostUpvotedEvent{AuthorID: authorID, VoterID: voterID, PostID: targetID})
		} else {
			// Changed from up to down
			s.eventBus.Publish(bus.PostUpvoteRemovedEvent{AuthorID: authorID, VoterID: voterID, PostID: targetID})
			s.eventBus.Publish(bus.PostDownvotedEvent{AuthorID: authorID, VoterID: voterID, PostID: targetID})
		}
	case model.VoteTargetComment:
		if newIsUpvote {
			// Changed from down to up
			s.eventBus.Publish(bus.CommentDownvoteRemovedEvent{AuthorID: authorID, VoterID: voterID, CommentID: targetID})
			s.eventBus.Publish(bus.CommentUpvotedEvent{AuthorID: authorID, VoterID: voterID, CommentID: targetID})
		} else {
			// Changed from up to down
			s.eventBus.Publish(bus.CommentUpvoteRemovedEvent{AuthorID: authorID, VoterID: voterID, CommentID: targetID})
			s.eventBus.Publish(bus.CommentDownvotedEvent{AuthorID: authorID, VoterID: voterID, CommentID: targetID})
		}
	}
}
