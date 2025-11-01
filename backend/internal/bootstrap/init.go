package bootstrap

import (
	"log"

	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/email"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/platform/ws"
	"github.com/giakiet05/lkforum/internal/repo"
	userroute "github.com/giakiet05/lkforum/internal/route/user"
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
	repo.NotificationRepo
	repo.ChannelRepo
}

type Services struct {
	service.AuthService
	service.UserService
	service.CommunityService
	service.MembershipService
	service.PostService
	service.CommentService
	service.ReputationService
	service.NotificationService
	service.ChannelService
}

type Controllers struct {
	controller.AuthController
	controller.UserController
	controller.CommunityController
	controller.MembershipController
	controller.PostController
	controller.CommentController
	controller.NotificationController
	controller.WebSocketController
	controller.ChannelController
}

func initRepos(client *mongo.Client, db *mongo.Database) *Repos {
	return &Repos{
		UserRepo:         repo.NewUserRepo(db),
		CommunityRepo:    repo.NewCommunityRepo(db),
		MembershipRepo:   repo.NewMembershipRepo(db),
		PostRepo:         repo.NewPostRepo(db),
		PostVoteRepo:     repo.NewPostVoteRepo(client, db),
		PostImageRepo:    repo.NewPostImageRepo(db),
		PostPollRepo:     repo.NewPostPollRepo(client, db),
		CommentRepo:      repo.NewCommentRepo(db),
		NotificationRepo: repo.NewNotificationRepo(db),
		ChannelRepo:      repo.NewChannelRepo(db),
	}
}

func initServices(repos *Repos, redisClient *redis.Client, emailSender email.Sender) *Services {
	return &Services{
		AuthService:         service.NewAuthService(repos.UserRepo, emailSender),
		UserService:         service.NewUserService(repos.UserRepo),
		CommunityService:    service.NewCommunityService(repos.CommunityRepo),
		MembershipService:   service.NewMembershipService(repos.MembershipRepo, redisClient),
		PostService:         service.NewPostService(repos.PostRepo, repos.PostVoteRepo, repos.PostPollRepo, repos.PostImageRepo, bus.Bus),
		CommentService:      service.NewCommentService(repos.CommentRepo, bus.Bus),
		ReputationService:   service.NewReputationService(repos.UserRepo),
		NotificationService: service.NewNotificationService(repos.NotificationRepo, repos.UserRepo, repos.PostRepo, repos.CommentRepo, bus.Bus),
		ChannelService:      service.NewChannelService(repos.ChannelRepo),
	}
}

func initControllers(services *Services) *Controllers {
	return &Controllers{
		AuthController:         *controller.NewAuthController(services.AuthService),
		UserController:         *controller.NewUserController(services.UserService),
		CommunityController:    *controller.NewCommunityController(services.CommunityService),
		MembershipController:   *controller.NewMembershipController(services.MembershipService),
		PostController:         *controller.NewPostController(services.PostService),
		CommentController:      *controller.NewCommentController(services.CommentService),
		NotificationController: *controller.NewNotificationController(services.NotificationService),
		WebSocketController:    *controller.NewWebSocketController(),
		ChannelController:      *controller.NewChannelController(services.ChannelService),
	}
}

func initRoutes(controllers *Controllers, r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := r.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to LKForum API!"})
	})

	userroute.RegisterAuthRoutes(api, &controllers.AuthController)
	userroute.RegisterUserRoutes(api, &controllers.UserController)
	userroute.RegisterCommunityRoutes(api, &controllers.CommunityController)
	userroute.RegisterMembershipRoutes(api, &controllers.MembershipController)
	userroute.RegisterPostRoutes(api, &controllers.PostController)
	userroute.RegisterCommentRoutes(api, &controllers.CommentController)
	userroute.RegisterNotificationRoutes(api, &controllers.NotificationController)
	userroute.RegisterWebSocketRoutes(api, &controllers.WebSocketController)
	userroute.RegisterChannelRoutes(api, &controllers.ChannelController)
}

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

	// Start background services
	go ws.WSHub.Start()
	services.ReputationService.Start()
	services.NotificationService.Start()

	return router, nil
}
