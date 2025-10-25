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
}

func NewReputationService(userRepo repo.UserRepo) ReputationService {
	return &reputationService{userRepo: userRepo}
}

// Start subscribes to relevant events and starts the reputation processing goroutine.
func (s *reputationService) Start() {
	eventChannel := make(bus.EventListener, 100)

	bus.Bus.Subscribe(bus.TopicPostCreated, eventChannel)
	bus.Bus.Subscribe(bus.TopicPostUpvoted, eventChannel)
	bus.Bus.Subscribe(bus.TopicPostDownvoted, eventChannel)
	bus.Bus.Subscribe(bus.TopicCommentCreated, eventChannel)

	log.Println("ReputationService started and subscribed to events.")

	go s.processEvents(eventChannel)
}

// processEvents runs in a separate goroutine, listening for and handling events.
func (s *reputationService) processEvents(ch bus.EventListener) {
	for event := range ch {
		switch event.Topic() {
		case bus.TopicPostCreated:
			s.handleReputationUpdate(event, "authorId", PointsPostCreated)
		case bus.TopicPostUpvoted:
			s.handleReputationUpdate(event, "authorId", PointsPostUpvoted)
		case bus.TopicPostDownvoted:
			s.handleReputationUpdate(event, "authorId", PointsPostDownvoted)
			s.handleReputationUpdate(event, "voterId", PointsDownvoteAction)
		case bus.TopicCommentCreated:
			s.handleReputationUpdate(event, "authorId", PointsCommentCreated)
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
