package service

import (
	"log"

	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
)

// Reputation Action Points
const (
	PointsPostCreated      = 2
	PointsPostUpvoted      = 10
	PointsPostDownvoted    = -2
	PointsCommentCreated   = 1
	PointsCommentUpvoted   = 5
	PointsCommentDownvoted = -1
	PointsDownvoteAction   = -1 // Penalty for the user who downvotes
)

type ReputationService interface {
	Start()
}

type reputationService struct {
	userRepo repo.UserRepo
	eventBus bus.EventBus
}

func NewReputationService(userRepo repo.UserRepo, bus bus.EventBus) ReputationService {
	return &reputationService{userRepo: userRepo, eventBus: bus}
}

// Start subscribes to relevant events and starts the reputation processing goroutine.
func (s *reputationService) Start() {
	eventChannel := make(bus.EventListener, 100)

	s.eventBus.Subscribe(bus.TopicPostApproved, eventChannel)
	s.eventBus.Subscribe(bus.TopicPostUpvoted, eventChannel)
	s.eventBus.Subscribe(bus.TopicPostDownvoted, eventChannel)
	s.eventBus.Subscribe(bus.TopicPostUpvoteRemoved, eventChannel)
	s.eventBus.Subscribe(bus.TopicPostDownvoteRemoved, eventChannel)
	s.eventBus.Subscribe(bus.TopicCommentApproved, eventChannel)
	s.eventBus.Subscribe(bus.TopicCommentUpvoted, eventChannel)
	s.eventBus.Subscribe(bus.TopicCommentDownvoted, eventChannel)
	s.eventBus.Subscribe(bus.TopicCommentUpvoteRemoved, eventChannel)
	s.eventBus.Subscribe(bus.TopicCommentDownvoteRemoved, eventChannel)

	log.Println("ReputationService started and subscribed to events.")

	go s.processEvents(eventChannel)
}

// processEvents runs in a separate goroutine, listening for and handling events.
func (s *reputationService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicPostApproved:
			s.handleReputationUpdate(event, "author_id", PointsPostCreated)
		case bus.TopicPostUpvoted:
			s.handleReputationUpdate(event, "author_id", PointsPostUpvoted)
		case bus.TopicPostDownvoted:
			s.handleReputationUpdate(event, "author_id", PointsPostDownvoted)
			s.handleReputationUpdate(event, "voter_id", PointsDownvoteAction)
		case bus.TopicPostUpvoteRemoved:
			// Revert the upvote: remove +10 points from author
			s.handleReputationUpdate(event, "author_id", -PointsPostUpvoted)
		case bus.TopicPostDownvoteRemoved:
			// Revert the downvote: add back +2 points to author, +1 to voter
			s.handleReputationUpdate(event, "author_id", -PointsPostDownvoted)
			s.handleReputationUpdate(event, "voter_id", -PointsDownvoteAction)
		case bus.TopicCommentApproved:
			s.handleReputationUpdate(event, "author_id", PointsCommentCreated)
		case bus.TopicCommentUpvoted:
			s.handleReputationUpdate(event, "author_id", PointsCommentUpvoted)
		case bus.TopicCommentDownvoted:
			s.handleReputationUpdate(event, "author_id", PointsCommentDownvoted)
			s.handleReputationUpdate(event, "voter_id", PointsDownvoteAction)
		case bus.TopicCommentUpvoteRemoved:
			// Revert the upvote: remove +5 points from author
			s.handleReputationUpdate(event, "author_id", -PointsCommentUpvoted)
		case bus.TopicCommentDownvoteRemoved:
			// Revert the downvote: add back +1 point to author, +1 to voter
			s.handleReputationUpdate(event, "author_id", -PointsCommentDownvoted)
			s.handleReputationUpdate(event, "voter_id", -PointsDownvoteAction)
		default:
			log.Printf("ReputationService: Received unknown event topic: %s", event.Topic())
		}
	}
}

// handleReputationUpdate is a helper to process a single reputation update.
func (s *reputationService) handleReputationUpdate(event bus.Event, userIdKey string, points int) {
	payload := event.Payload()
	userID, ok := payload[userIdKey].(string)
	if !ok || userID == "" {
		log.Printf("ERROR: ReputationService: could not get '%s' from event payload for topic %s", userIdKey, event.Topic())
		return
	}

	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if err := s.userRepo.UpdateReputation(ctx, userID, points); err != nil {
		log.Printf("ERROR: ReputationService: failed to update reputation for user %s: %v", userID, err)
	}
}
