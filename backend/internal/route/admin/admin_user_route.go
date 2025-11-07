package adminroute

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAdminUserRoutes(rg *gin.RouterGroup, c *controller.UserController) {
	// All routes in this group are for admin management of users
	// It is expected that the parent router group `rg` is already prefixed with "/admin"
	// and has the appropriate admin middleware.
	users := rg.Group("/users")
	users.Use(middleware.RequireAuth(), middleware.RequireAdmin())
	users.GET("", c.GetUsers)
	users.DELETE("/:id", c.DeleteUser)

	// Other admin actions on users can be added here, e.g.:
	// PUT /:id/ban
	// PUT /:id/role
}
