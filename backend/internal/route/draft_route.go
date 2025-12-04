package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterDraftRoutes(rg *gin.RouterGroup, c *controller.DraftController) {
	drafts := rg.Group("/drafts")
	drafts.Use(middleware.RequireAuth())
	{
		drafts.POST("", c.CreateDraft)
		drafts.GET("", c.GetDrafts)
		drafts.GET("/:id", c.GetDraftByID)
		drafts.PUT("/:id", c.UpdateDraft)
		drafts.DELETE("/:id", c.DeleteDraft)
		drafts.POST("/:id/publish", c.PublishDraft)
	}
}
