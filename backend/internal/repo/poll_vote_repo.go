package repo

import (
	"context"
	"time"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PollVoteRepo defines the data access layer for post polls.
type PollVoteRepo interface {
	Vote(ctx context.Context, userID, postID, optionID string) error
	RemoveVote(ctx context.Context, userID, postID string) error
	GetUserVoteIDs(ctx context.Context, userID, postID string) ([]string, error)
	FindUserVotes(ctx context.Context, userID string, postIDs []string) (map[string][]string, error)
}

type pollVoteRepo struct {
	client             *mongo.Client
	postCollection     *mongo.Collection
	pollVoteCollection *mongo.Collection
}

// NewPollVoteRepo creates a new instance of PollVoteRepo.
func NewPollVoteRepo(client *mongo.Client, db *mongo.Database) PollVoteRepo {
	return &pollVoteRepo{
		client:             client,
		postCollection:     db.Collection(config.PostColName),
		pollVoteCollection: db.Collection(config.PollVoteColName),
	}
}

func (r *pollVoteRepo) GetUserVoteIDs(ctx context.Context, userID, postID string) ([]string, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	var votes []*model.PollVote
	filter := bson.M{"post_id": postObjID, "user_id": userObjID}
	cursor, err := r.pollVoteCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &votes); err != nil {
		return nil, err
	}

	optionIDs := make([]string, len(votes))
	for i, v := range votes {
		optionIDs[i] = v.OptionID
	}
	return optionIDs, nil
}

func (r *pollVoteRepo) FindUserVotes(ctx context.Context, userID string, postIDs []string) (map[string][]string, error) {
	if len(postIDs) == 0 {
		return make(map[string][]string), nil
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	postObjIDs := make([]primitive.ObjectID, 0, len(postIDs))
	for _, pid := range postIDs {
		if objID, err := primitive.ObjectIDFromHex(pid); err == nil {
			postObjIDs = append(postObjIDs, objID)
		}
	}

	filter := bson.M{"user_id": userObjID, "post_id": bson.M{"$in": postObjIDs}}

	cursor, err := r.pollVoteCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	userVotesMap := make(map[string][]string)
	for cursor.Next(ctx) {
		var vote model.PollVote
		if err := cursor.Decode(&vote); err != nil {
			continue
		}
		postIDStr := vote.PostID.Hex()
		userVotesMap[postIDStr] = append(userVotesMap[postIDStr], vote.OptionID)
	}

	return userVotesMap, cursor.Err()
}

func (r *pollVoteRepo) Vote(ctx context.Context, userID, postID, optionID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperror.ErrInvalidID
	}
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		var post model.Post
		if err := r.postCollection.FindOne(sessCtx, bson.M{"_id": postObjID}).Decode(&post); err != nil {
			return nil, apperror.ErrPostNotFound
		}

		if post.Type != model.PostTypePoll || post.Content == nil || post.Content.Poll == nil {
			return nil, apperror.ErrBadRequest
		}

		userVoteIDs, err := r.getUserVoteIDsInTx(sessCtx, userObjID, postObjID)
		if err != nil {
			return nil, err
		}

		// Check if user already voted this option → Toggle off (remove vote)
		isAlreadyVoted := false
		for _, id := range userVoteIDs {
			if id == optionID {
				isAlreadyVoted = true
				break
			}
		}
		if isAlreadyVoted {
			// Toggle off: Remove this specific vote
			if err := r.removeSingleVoteInTransaction(sessCtx, userObjID, postObjID, optionID); err != nil {
				return nil, err
			}
			return nil, nil
		}

		// For single-choice polls: remove existing votes before adding new one
		if !post.Content.Poll.AllowMultiple && len(userVoteIDs) > 0 {
			if err := r.removeVotesInTransaction(sessCtx, userObjID, postObjID); err != nil {
				return nil, err
			}
		}

		newVote := &model.PollVote{
			UserID:    userObjID,
			PostID:    postObjID,
			OptionID:  optionID,
			CreatedAt: time.Now(),
		}
		if _, err := r.pollVoteCollection.InsertOne(sessCtx, newVote); err != nil {
			return nil, err
		}

		filter := bson.M{"_id": postObjID, "content.poll.options.id": optionID}
		update := bson.M{"$inc": bson.M{"content.poll.options.$.votes": 1, "content.poll.total_votes": 1}}
		_, err = r.postCollection.UpdateOne(sessCtx, filter, update)
		return nil, err
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

func (r *pollVoteRepo) RemoveVote(ctx context.Context, userID, postID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperror.ErrInvalidID
	}
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		return nil, r.removeVotesInTransaction(sessCtx, userObjID, postObjID)
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

func (r *pollVoteRepo) getUserVoteIDsInTx(ctx context.Context, userID, postID primitive.ObjectID) ([]string, error) {
	var votes []*model.PollVote
	filter := bson.M{"post_id": postID, "user_id": userID}
	cursor, err := r.pollVoteCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &votes); err != nil {
		return nil, err
	}

	optionIDs := make([]string, len(votes))
	for i, v := range votes {
		optionIDs[i] = v.OptionID
	}
	return optionIDs, nil
}

func (r *pollVoteRepo) removeSingleVoteInTransaction(sessCtx context.Context, userID, postID primitive.ObjectID, optionID string) error {
	// Decrement vote count for the specific option
	filter := bson.M{"_id": postID, "content.poll.options.id": optionID}
	update := bson.M{"$inc": bson.M{"content.poll.options.$.votes": -1, "content.poll.total_votes": -1}}
	if _, err := r.postCollection.UpdateOne(sessCtx, filter, update); err != nil {
		return err
	}

	// Delete the specific poll vote
	_, err := r.pollVoteCollection.DeleteOne(sessCtx, bson.M{
		"post_id":   postID,
		"user_id":   userID,
		"option_id": optionID,
	})
	return err
}

func (r *pollVoteRepo) removeVotesInTransaction(sessCtx context.Context, userID, postID primitive.ObjectID) error {
	votes, err := r.getUserVoteIDsInTx(sessCtx, userID, postID)
	if err != nil || len(votes) == 0 {
		return err
	}

	// Decrement vote count for each option separately
	// Note: Must loop because MongoDB's $ operator only updates the first matched element
	for _, optionID := range votes {
		filter := bson.M{"_id": postID, "content.poll.options.id": optionID}
		update := bson.M{"$inc": bson.M{"content.poll.options.$.votes": -1}}
		if _, err := r.postCollection.UpdateOne(sessCtx, filter, update); err != nil {
			return err
		}
	}

	// Decrement total_votes
	totalUpdate := bson.M{"$inc": bson.M{"content.poll.total_votes": -len(votes)}}
	if _, err := r.postCollection.UpdateOne(sessCtx, bson.M{"_id": postID}, totalUpdate); err != nil {
		return err
	}

	// Delete all poll votes
	_, err = r.pollVoteCollection.DeleteMany(sessCtx, bson.M{"post_id": postID, "user_id": userID})
	return err
}
