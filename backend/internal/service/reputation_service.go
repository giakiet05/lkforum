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

// Event Topics
const (
	TopicPostCreated      = "post.created"
	TopicPostUpvoted      = "post.upvoted"
	TopicPostDownvoted    = "post.downvoted"
	TopicCommentCreated   = "comment.created"
	TopicCommentUpvoted   = "comment.upvoted"
	TopicCommentDownvoted = "comment.downvoted"
)

type ReputationService interface {
	Start()
}

type reputationService struct {
	userRepo repo.UserRepo
}

func NewReputationService(userRepo repo.UserRepo) ReputationService {
	return &reputationService{userRepo: userRepo}
}

// Start subscribes to relevant events and starts the reputation processing goroutine.
func (s *reputationService) Start() {
	// Create a single channel to listen to all reputation-related events.
	// A buffered channel can help prevent blocking the event bus if this service is slow.
	eventChannel := make(bus.EventListener, 100)

	bus.Bus.Subscribe(TopicPostCreated, eventChannel)
	bus.Bus.Subscribe(TopicPostUpvoted, eventChannel)
	bus.Bus.Subscribe(TopicPostDownvoted, eventChannel)
	bus.Bus.Subscribe(TopicCommentCreated, eventChannel)
	// Add more subscriptions here as needed (e.g., for comments)

	log.Println("ReputationService started and subscribed to events.")

	// Start a background worker to process events from the channel.
	go s.processEvents(eventChannel)
}

// processEvents runs in a separate goroutine, listening for and handling events.
func (s *reputationService) processEvents(ch bus.EventListener) {
	for event := range ch {
		// Use a switch to handle different event topics.
		switch event.Topic() {
		case TopicPostCreated:
			s.handleReputationUpdate(event, "authorId", PointsPostCreated)
		case TopicPostUpvoted:
			s.handleReputationUpdate(event, "authorId", PointsPostUpvoted)
		case TopicPostDownvoted:
			s.handleReputationUpdate(event, "authorId", PointsPostDownvoted)
			s.handleReputationUpdate(event, "voterId", PointsDownvoteAction)
		case TopicCommentCreated:
			s.handleReputationUpdate(event, "authorId", PointsCommentCreated)
		// Add more cases here for other events
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
