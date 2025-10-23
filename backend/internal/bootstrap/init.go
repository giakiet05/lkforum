package bootstrap

import (
	"log"

	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/email"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/route/user"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repos struct {
	repo.UserRepo
	repo.CommunityRepo
	repo.MembershipRepo
	repo.PostRepo
	repo.PostPollRepo
	repo.PostImageRepo
	repo.PostVoteRepo
	repo.CommentRepo
}

type Services struct {
	service.AuthService
	service.UserService
	service.CommunityService
	service.MembershipService
	service.PostService
	service.CommentService
}

type Controllers struct {
	controller.AuthController
	controller.UserController
	controller.CommunityController
	controller.MembershipController
	controller.PostController
	controller.CommentController
}

// initRepos initializes repositories with the given database
func initRepos(client *mongo.Client, db *mongo.Database) *Repos {
	return &Repos{
		UserRepo:       repo.NewUserRepo(db),
		CommunityRepo:  repo.NewCommunityRepo(db),
		MembershipRepo: repo.NewMembershipRepo(db),
		PostRepo:       repo.NewPostRepo(db),
		PostVoteRepo:   repo.NewPostVoteRepo(client, db),
		PostImageRepo:  repo.NewPostImageRepo(db),
		PostPollRepo:   repo.NewPostPollRepo(client, db),
		CommentRepo:    repo.NewCommentRepo(db),
	}
}

// initServices Initialize services with the given repositories
func initServices(repos *Repos, redisClient *redis.Client, emailSender email.Sender) *Services {
	return &Services{
		AuthService:       service.NewAuthService(repos.UserRepo, emailSender),
		UserService:       service.NewUserService(repos.UserRepo),
		CommunityService:  service.NewCommunityService(repos.CommunityRepo),
		MembershipService: service.NewMembershipService(repos.MembershipRepo, redisClient),
		PostService:       service.NewPostService(repos.PostRepo, repos.PostVoteRepo, repos.PostPollRepo, repos.PostImageRepo),
		CommentService:    service.NewCommentService(repos.CommentRepo),
	}
}

// initControllers Initialize controllers with the given services
func initControllers(services *Services) *Controllers {
	return &Controllers{
		AuthController:       *controller.NewAuthController(services.AuthService),
		UserController:       *controller.NewUserController(services.UserService),
		CommunityController:  *controller.NewCommunityController(services.CommunityService),
		MembershipController: *controller.NewMembershipController(services.MembershipService),
		PostController:       *controller.NewPostController(services.PostService),
		CommentController:    *controller.NewCommentController(services.CommentService),
	}
}

// initRoutes sets up the routes for the Gin engine
func initRoutes(controllers *Controllers, r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := r.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to LKForum API!"})
	})

	// Register routes
	userroute.RegisterAuthRoutes(api, &controllers.AuthController)
	userroute.RegisterUserRoutes(api, &controllers.UserController)
	userroute.RegisterCommunityRoutes(api, &controllers.CommunityController)
	userroute.RegisterMembershipRoutes(api, &controllers.MembershipController)
	userroute.RegisterPostRoutes(api, &controllers.PostController)
	userroute.RegisterCommentRoutes(api, &controllers.CommentController)
}

// Init initializes all application components
func Init() (*gin.Engine, error) {
	config.LoadConfig()
	auth.InitGoogleOAuthConfig()

	redisClient := config.NewRedisClient()

	if err := InitializeTokenService(redisClient); err != nil {
		log.Printf("Warning: Token invalidation service not available: %v\n", err)
	}

	client := config.NewMongoClient()
	db := client.Database(config.Cfg.DBName)
	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", config.Cfg.FrontendURL)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	emailSender := email.NewSMTPSender()
	repos := initRepos(client, db)
	services := initServices(repos, redisClient, emailSender)
	controllers := initControllers(services)
	initRoutes(controllers, router)

	return router, nil
}
