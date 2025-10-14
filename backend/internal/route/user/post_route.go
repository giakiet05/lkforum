// internal/route/post_route.go
package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterPostRoutes đăng ký các routes cho Post.
// Middleware AuthMiddleware() sẽ xác thực token và đưa thông tin user vào context.
// Middleware AuthOptional() sẽ làm tương tự nhưng không báo lỗi nếu không có token.
func RegisterPostRoutes(rg *gin.RouterGroup, c *controller.PostController) {
	posts := rg.Group("/posts")

	// --- 👁️ Public / Optional Auth Routes ---
	// Các route này ai cũng có thể truy cập.
	// Dùng AuthOptional để làm giàu dữ liệu (VD: trạng thái vote) cho người đã đăng nhập.
	{

	}

	// --- 🛡️ Private / Required Auth Routes ---
	// Các route này yêu cầu phải đăng nhập.
	posts.Use(middleware.AuthMiddleware())
	{
		posts.GET("", c.GetPosts)        // Xem danh sách bài đăng
		posts.GET("/:id", c.GetPostByID) // Xem chi tiết một bài đăng
		posts.POST("", c.CreatePost)
		posts.PUT("/:id", c.UpdatePost)
		posts.DELETE("/:id", c.DeletePost)

		// Tương tác
		posts.POST("/:id/vote", c.VoteOnPost)
		posts.POST("/:id/poll/vote", c.VoteOnPoll)

		// Quản lý Image
		posts.POST("/:id/images", c.AddImagesToPost)
		posts.DELETE("/:id/images", c.RemoveImagesFromPost)

		// Quản lý Poll
		posts.PUT("/:id/poll", c.UpdatePollDetails)
		posts.POST("/:id/poll/options", c.AddPollOptions)
		posts.DELETE("/:id/poll/options", c.RemovePollOptions)
		posts.PUT("/:id/poll/options/:optionID", c.UpdatePollOption)
	}
}
