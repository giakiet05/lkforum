package service

import (
	"log"

	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
)

type MigrationService interface {
	RecalculateAllPostCommentCounts() error
}

type migrationService struct {
	postRepo    repo.PostRepo
	commentRepo repo.CommentRepo
}

func NewMigrationService(postRepo repo.PostRepo, commentRepo repo.CommentRepo) MigrationService {
	return &migrationService{
		postRepo:    postRepo,
		commentRepo: commentRepo,
	}
}

func (s *migrationService) RecalculateAllPostCommentCounts() error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	// Get all posts
	posts, _, err := s.postRepo.Find(ctx, repo.Filter{}, nil)
	if err != nil {
		return err
	}

	log.Printf("Starting comment count recalculation for %d posts", len(posts))

	for _, post := range posts {
		postID := post.ID.Hex()

		// Count comments for this post (only top-level, non-deleted)
		comments, _, err := s.commentRepo.GetCommentsFilterPaginated(
			ctx,
			&postID,
			nil,  // parentID = nil (only top-level)
			nil,  // userID
			nil,  // content
			nil,  // currentUserID (get all approved)
			1,    // page
			1000, // pageSize (get all)
		)

		if err != nil {
			log.Printf("Error counting comments for post %s: %v", postID, err)
			continue
		}

		count := len(comments)

		// Update post's comment count
		err = s.postRepo.UpdateByID(ctx, postID, map[string]interface{}{
			"$set": map[string]interface{}{"comments_count": count},
		})

		if err != nil {
			log.Printf("Error updating comment count for post %s: %v", postID, err)
			continue
		}

		log.Printf("Updated post %s: %d comments", postID, count)
	}

	log.Println("Comment count recalculation completed")
	return nil
}
