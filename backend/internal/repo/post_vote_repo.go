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

type PostVoteRepo interface {
	Vote(ctx context.Context, userID, postID string, voteValue bool) error
	RemoveVote(ctx context.Context, userID, postID string) error
	GetUserVote(ctx context.Context, userID, postID string) (*model.Vote, error)
	FindUserVotes(ctx context.Context, userID string, postIDs []string) (map[string]string, error)
}

type postVoteRepo struct {
	client         *mongo.Client
	postCollection *mongo.Collection
	voteCollection *mongo.Collection
}

func NewPostVoteRepo(client *mongo.Client, db *mongo.Database) PostVoteRepo {
	return &postVoteRepo{
		client:         client,
		postCollection: db.Collection(config.PostColName),
		voteCollection: db.Collection(config.VoteColName),
	}
}

func (r *postVoteRepo) GetUserVote(ctx context.Context, userID, postID string) (*model.Vote, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}
	postObjID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	var vote model.Vote
	filter := bson.M{
		"target_id":   postObjID,
		"user_id":     userObjID,
		"target_type": model.VoteTargetPost,
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

func (r *postVoteRepo) FindUserVotes(ctx context.Context, userID string, postIDs []string) (map[string]string, error) {
	if len(postIDs) == 0 {
		return make(map[string]string), nil
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, apperror.ErrInvalidID
	}

	postObjIDs := make([]primitive.ObjectID, len(postIDs))
	for i, pid := range postIDs {
		objID, err := primitive.ObjectIDFromHex(pid)
		if err != nil {
			continue // Skip invalid IDs
		}
		postObjIDs[i] = objID
	}

	filter := bson.M{
		"user_id":     userObjID,
		"target_type": model.VoteTargetPost,
		"target_id":   bson.M{"$in": postObjIDs},
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

func (r *postVoteRepo) Vote(ctx context.Context, userID, postID string, voteValue bool) error {
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
		previousVote, err := r.getUserVoteInTx(sessCtx, userObjID, postObjID)
		if err != nil {
			return nil, err
		}

		upInc, downInc := 0, 0

		if previousVote == nil {
			// Case 1: New vote
			newVote := &model.Vote{
				UserID:     userObjID,
				TargetID:   postObjID,
				TargetType: model.VoteTargetPost,
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
				return nil, r.removeVoteInTransaction(sessCtx, previousVote)
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

		// Apply counter update to the posts collection
		postFilter := bson.M{"_id": postObjID}
		postUpdate := bson.M{"$inc": bson.M{"votes_count.up": upInc, "votes_count.down": downInc}}
		if _, err := r.postCollection.UpdateOne(sessCtx, postFilter, postUpdate); err != nil {
			return nil, err
		}
		return nil, nil
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

func (r *postVoteRepo) RemoveVote(ctx context.Context, userID, postID string) error {
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
		voteToRemove, err := r.getUserVoteInTx(sessCtx, userObjID, postObjID)
		if err != nil {
			return nil, err
		}
		if voteToRemove == nil {
			return nil, apperror.ErrVoteNotFound
		}
		return nil, r.removeVoteInTransaction(sessCtx, voteToRemove)
	}

	_, err = session.WithTransaction(ctx, callback)
	return err
}

// getUserVoteInTx is a helper to get a vote within a transaction context.
func (r *postVoteRepo) getUserVoteInTx(ctx context.Context, userID, postID primitive.ObjectID) (*model.Vote, error) {
	var vote model.Vote
	filter := bson.M{"target_id": postID, "user_id": userID, "target_type": model.VoteTargetPost}
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
func (r *postVoteRepo) removeVoteInTransaction(sessCtx mongo.SessionContext, vote *model.Vote) error {
	if _, err := r.voteCollection.DeleteOne(sessCtx, bson.M{"_id": vote.ID}); err != nil {
		return err
	}

	upInc, downInc := 0, 0
	if vote.Value {
		upInc = -1
	} else {
		downInc = -1
	}

	postFilter := bson.M{"_id": vote.TargetID}
	postUpdate := bson.M{"$inc": bson.M{"votes_count.up": upInc, "votes_count.down": downInc}}
	_, err := r.postCollection.UpdateOne(sessCtx, postFilter, postUpdate)
	return err
}
