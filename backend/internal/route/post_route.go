package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterPostRoutes(rg *gin.RouterGroup, c *controller.PostController) {
	posts := rg.Group("/posts")

	// --- Public routes that can be enriched with auth info ---
	// Anyone can access these. If a valid token is provided, extra info (like user's vote) is added.
	posts.GET("", middleware.LoadUserIfAuthenticated(), c.GetPosts)
	posts.GET("/:id", middleware.LoadUserIfAuthenticated(), c.GetPostByID)

	// --- Private routes that require authentication ---
	private := posts.Group("/")
	private.Use(middleware.RequireAuth())
	{
		// Basic CRUD
		private.POST("", c.CreatePost)
		private.PUT("/:id", c.UpdatePost)
		private.DELETE("/:id", c.DeletePost)

		// My Posts
		private.GET("/me", c.GetMyPosts)

		// Save, Hide & Report
		private.GET("/saved", c.GetSavedPosts)
		private.GET("/hidden", c.GetHiddenPosts)
		private.POST("/:id/save", c.SavePost)
		private.DELETE("/:id/save", c.UnsavePost)
		private.POST("/:id/report", c.ReportPost)
		private.POST("/:id/hide", c.HidePost)
		private.POST("/:id/unhide", c.UnhidePost)

		// Get Ban Post Community
		private.GET("/banned", c.GetBanPosts)

		// Image Management
		private.POST("/:id/images", c.AddImagesToPost)
		private.DELETE("/:id/images", c.RemoveImagesFromPost) // Body: { "public_ids": [...] }

		// Video Management
		private.POST("/:id/videos", c.AddVideosToPost)
		private.DELETE("/:id/videos", c.RemoveVideosFromPost) // Body: { "public_ids": [...] }

		// Poll Management
		private.PUT("/:id/poll", c.UpdatePoll)
		private.POST("/:id/poll/options", c.AddPollOptions)
		private.DELETE("/:id/poll/options", c.RemovePollOptions) // Body: { "option_ids": [...] }
		private.PUT("/:id/poll/options/:optionID", c.UpdatePollOption)

		// Poll voting (kept in post controller since it's poll-specific)
		private.POST("/:id/poll/vote", c.VoteOnPoll)
		private.DELETE("/:id/poll/vote", c.RemovePollVote)

		// Note: Post voting moved to /api/votes/post/:id (use VoteController)
	}
}
