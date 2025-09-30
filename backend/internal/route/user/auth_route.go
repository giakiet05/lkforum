package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(rg *gin.RouterGroup, c *controller.UserController) {
	auth := rg.Group("/auth")
	auth.POST("/register", c.RegisterUser)
	auth.POST("/login", c.Login)
	auth.POST("/refresh", c.RefreshToken)
}
