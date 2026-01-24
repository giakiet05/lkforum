package repo

import (
	"context"
	"errors"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VoteRepo interface {
	Vote(ctx context.Context, userID, targetID string, targetType model.VoteTargetType, voteValue bool) error
	RemoveVote(ctx context.Context, userID, targetID string, targetType model.VoteTargetType) error
	GetUserVote(ctx context.Context, userID, targetID string, targetType model.VoteTargetType) (*model.Vote, error)
	FindUserVotes(ctx context.Context, userID string, targetIDs []string, targetType model.VoteTargetType) (map[string]string, error)
}

type voteRepo struct {
	client            *mongo.Client
	postCollection    *mongo.Collection
	commentCollection *mongo.Collection
	voteCollection    *mongo.Collection
}

func NewVoteRepo(client *mongo.Client, db *mongo.Database) VoteRepo {
	return &voteRepo{
		client:            client,
		postCollection:    db.Collection(config.PostColName),
		commentCollection: db.Collection(config.CommentColName),
		voteCollection:    db.Collection(config.VoteColName),
	}
}

func (r *voteRepo) GetUserVote(ctx context.Context, userID, targetID string, targetType model.VoteTargetType) (*model.Vote, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}
	targetObjID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	var vote model.Vote
	filter := bson.M{
		"target_id":   targetObjID,
		"user_id":     userObjID,
		"target_type": targetType,
	}

	err = r.voteCollection.FindOne(ctx, filter).Decode(&vote)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Not an error, just means no vote exists
		}
		return nil, err
	}
	return &vote, nil
}

func (r *voteRepo) FindUserVotes(ctx context.Context, userID string, targetIDs []string, targetType model.VoteTargetType) (map[string]string, error) {
	if len(targetIDs) == 0 {
		return make(map[string]string), nil
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	targetObjIDs := make([]primitive.ObjectID, len(targetIDs))
	for i, tid := range targetIDs {
		objID, err := primitive.ObjectIDFromHex(tid)
		if err != nil {
			continue // Skip invalid IDs
		}
		targetObjIDs[i] = objID
	}

	filter := bson.M{
		"user_id":     userObjID,
		"target_type": targetType,
		"target_id":   bson.M{"$in": targetObjIDs},
	}

	cursor, err := r.voteCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	userVotesMap := make(map[string]string)
	for cursor.Next(ctx) {
		var vote model.Vote
		if err := cursor.Decode(&vote); err != nil {
			continue
		}
		if vote.Value {
			userVotesMap[vote.TargetID.Hex()] = "up"
		} else {
			userVotesMap[vote.TargetID.Hex()] = "down"
		}
	}

	return userVotesMap, cursor.Err()
}

func (r *voteRepo) Vote(ctx context.Context, userID, targetID string, targetType model.VoteTargetType, voteValue bool) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperror.ErrInvalidID
	}
	targetObjID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		previousVote, err := r.getUserVoteInTx(sessCtx, userObjID, targetObjID, targetType)
		if err != nil {
			return nil, err
		}

		upInc, downInc := 0, 0

		if previousVote == nil {
			// Case 1: New vote
			newVote := &model.Vote{
				UserID:     userObjID,
				TargetID:   targetObjID,
				TargetType: targetType,
				Value:      voteValue,
			}
			if _, err := r.voteCollection.InsertOne(sessCtx, newVote); err != nil {
				return nil, err
			}
			if voteValue {
				upInc = 1
			} else {
				downInc = 1
			}
		} else {
			if previousVote.Value == voteValue {
				// Case 2: Un-voting (e.g., clicking upvote again)
				return nil, r.removeVoteInTransaction(sessCtx, previousVote, targetType)
			} else {
				// Case 3: Changing vote (e.g., from up to down)
				update := bson.M{"$set": bson.M{"value": voteValue}}
				if _, err := r.voteCollection.UpdateByID(sessCtx, previousVote.ID, update); err != nil {
					return nil, err
				}
				if voteValue {
					upInc = 1
					downInc = -1
				} else {
					upInc = -1
					downInc = 1
				}
			}
		}

		// Apply counter update to the appropriate collection
		if err := r.updateTargetVoteCount(sessCtx, targetObjID, targetType, upInc, downInc); err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

func (r *voteRepo) RemoveVote(ctx context.Context, userID, targetID string, targetType model.VoteTargetType) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apperror.ErrInvalidID
	}
	targetObjID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		return apperror.ErrInvalidID
	}

	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	callback := func(sessCtx mongo.SessionContext) (interface{}, error) {
		voteToRemove, err := r.getUserVoteInTx(sessCtx, userObjID, targetObjID, targetType)
		if err != nil {
			return nil, err
		}
		if voteToRemove == nil {
			return nil, apperror.ErrVoteNotFound
		}
		return nil, r.removeVoteInTransaction(sessCtx, voteToRemove, targetType)
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

// getUserVoteInTx is a helper to get a vote within a transaction context.
func (r *voteRepo) getUserVoteInTx(ctx context.Context, userID, targetID primitive.ObjectID, targetType model.VoteTargetType) (*model.Vote, error) {
	var vote model.Vote
	filter := bson.M{"target_id": targetID, "user_id": userID, "target_type": targetType}
	err := r.voteCollection.FindOne(ctx, filter).Decode(&vote)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &vote, nil
}

// removeVoteInTransaction is a helper to delete a vote and update the counter within a transaction.
func (r *voteRepo) removeVoteInTransaction(sessCtx mongo.SessionContext, vote *model.Vote, targetType model.VoteTargetType) error {
	if _, err := r.voteCollection.DeleteOne(sessCtx, bson.M{"_id": vote.ID}); err != nil {
		return err
	}

	upInc, downInc := 0, 0
	if vote.Value {
		upInc = -1
	} else {
		downInc = -1
	}

	return r.updateTargetVoteCount(sessCtx, vote.TargetID, targetType, upInc, downInc)
}

// updateTargetVoteCount updates the vote count for either post or comment
func (r *voteRepo) updateTargetVoteCount(ctx context.Context, targetID primitive.ObjectID, targetType model.VoteTargetType, upInc, downInc int) error {
	filter := bson.M{"_id": targetID}

	var collection *mongo.Collection
	switch targetType {
	case model.VoteTargetPost:
		collection = r.postCollection

		// For posts, also update hot_score
		// First get the current post to calculate new hot_score
		var post model.Post
		if err := collection.FindOne(ctx, filter).Decode(&post); err != nil {
			return err
		}

		newUpvotes := 0
		newDownvotes := 0
		if post.VotesCount != nil {
			newUpvotes = post.VotesCount.Up + upInc
			newDownvotes = post.VotesCount.Down + downInc
		} else {
			newUpvotes = upInc
			newDownvotes = downInc
		}

		hotScore := model.CalculateHotScore(newUpvotes, newDownvotes, post.CreatedAt)

		update := bson.M{
			"$inc": bson.M{"votes_count.up": upInc, "votes_count.down": downInc},
			"$set": bson.M{"hot_score": hotScore},
		}
		_, err := collection.UpdateOne(ctx, filter, update)
		return err

	case model.VoteTargetComment:
		collection = r.commentCollection
		update := bson.M{"$inc": bson.M{"votes_count.up": upInc, "votes_count.down": downInc}}
		_, err := collection.UpdateOne(ctx, filter, update)
		return err

	default:
		return apperror.ErrInvalidID
	}
}
