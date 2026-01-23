package bootstrap

import (
	"log"

	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/platform/email"
	"github.com/giakiet05/lkforum/internal/platform/gemini"
	"github.com/giakiet05/lkforum/internal/platform/ws"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/route"
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
	repo.PollVoteRepo
	repo.VoteRepo
	repo.CommentRepo
	repo.NotificationRepo
	repo.ChannelRepo
	repo.MessageRepo
	repo.PostHistoryRepo
	repo.EmailVerificationRepo
	repo.PasswordResetRepo
	repo.SavedPostRepo
	repo.ReportRepo
	repo.DraftRepo
}

type Services struct {
	service.AuthService
	service.UserService
	service.CommunityService
	service.MembershipService
	service.PostService
	service.CommentService
	service.VoteService
	service.ReputationService
	service.NotificationService
	service.ChannelService
	service.MessageService
	service.PostHistoryService
	service.DraftService
	service.ReportService
	service.ModerationService
	service.AdminUserService
	service.AdminCommunityService
	service.AdminStatsService
	service.AdminAuthService
}

type Controllers struct {
	controller.AuthController
	controller.UserController
	controller.CommunityController
	controller.MembershipController
	controller.PostController
	controller.VoteController // Added VoteController
	controller.CommentController
	controller.NotificationController
	controller.WebSocketController
	controller.ChannelController
	controller.MessageController
	controller.PostHistoryController
	controller.DraftController
	controller.ReportController
	controller.AdminUserController
	controller.AdminCommunityController
	controller.AdminStatsController
	controller.AdminAuthController
	controller.DebugController
}

func initRepos(client *mongo.Client, db *mongo.Database) *Repos {
	return &Repos{
		UserRepo:              repo.NewUserRepo(db),
		CommunityRepo:         repo.NewCommunityRepo(db),
		MembershipRepo:        repo.NewMembershipRepo(db),
		PostRepo:              repo.NewPostRepo(db),
		VoteRepo:              repo.NewVoteRepo(client, db),
		PollVoteRepo:          repo.NewPollVoteRepo(client, db),
		CommentRepo:           repo.NewCommentRepo(db),
		NotificationRepo:      repo.NewNotificationRepo(db),
		ChannelRepo:           repo.NewChannelRepo(db),
		MessageRepo:           repo.NewMessageRepo(db),
		PostHistoryRepo:       repo.NewPostHistoryRepo(db),
		EmailVerificationRepo: repo.NewEmailVerificationRepo(db),
		PasswordResetRepo:     repo.NewPasswordResetRepo(db),
		SavedPostRepo:         repo.NewSavedPostRepo(db),
		ReportRepo:            repo.NewReportRepo(db),
		DraftRepo:             repo.NewDraftRepo(db),
	}
}

func initServices(repos *Repos, redisClient *redis.Client, emailSender email.Sender, eventBus bus.EventBus, geminiClient *gemini.GeminiClient, tokenService *auth.TokenService) *Services {
	services := &Services{
		AuthService:         service.NewAuthService(repos.UserRepo, repos.EmailVerificationRepo, repos.PasswordResetRepo, emailSender, redisClient, tokenService),
		UserService:         service.NewUserService(repos.UserRepo, eventBus, redisClient),
		MembershipService:   service.NewMembershipService(repos.MembershipRepo, repos.CommunityRepo, redisClient),
		ReputationService:   service.NewReputationService(repos.UserRepo, eventBus),
		NotificationService: service.NewNotificationService(repos.NotificationRepo, repos.UserRepo, repos.PostRepo, repos.CommentRepo, repos.CommunityRepo, eventBus, redisClient),
		ChannelService:      service.NewChannelService(repos.ChannelRepo, eventBus),
		MessageService:      service.NewMessageService(repos.MessageRepo, repos.ChannelRepo, repos.UserRepo, eventBus, redisClient),
		PostHistoryService:  service.NewPostHistoryService(repos.PostHistoryRepo),
		ReportService:       service.NewReportService(repos.ReportRepo),
		CommentService:      service.NewCommentService(repos.CommentRepo, repos.UserRepo, repos.CommunityRepo, repos.PostRepo, eventBus),
		CommunityService:    service.NewCommunityService(repos.CommunityRepo, repos.MembershipRepo, repos.PostRepo, repos.UserRepo, eventBus),
	}

	// Set MembershipService in CommunityService to get real-time member count
	if communitySvc, ok := services.CommunityService.(interface {
		SetMembershipService(service.MembershipService)
	}); ok {
		communitySvc.SetMembershipService(services.MembershipService)
	}

	// VoteService needs to be created first as PostService and CommentService depend on it
	services.VoteService = service.NewVoteService(repos.VoteRepo, repos.PostRepo, repos.CommentRepo, eventBus)

	// PostService and CommentService need VoteService
	services.PostService = service.NewPostService(repos.PostRepo, services.VoteService, repos.PollVoteRepo, repos.UserRepo, repos.CommunityRepo, repos.MembershipRepo, repos.SavedPostRepo, repos.ReportRepo, eventBus)

	// DraftService needs PostService
	services.DraftService = service.NewDraftService(repos.DraftRepo, repos.PostRepo, services.PostService)

	// ModerationService needs CommunityRepo for checking PostRequireApproval
	services.ModerationService = service.NewModerationService(repos.PostRepo, repos.CommentRepo, repos.UserRepo, repos.CommunityRepo, geminiClient, eventBus, &config.Cfg.Gemini)

	// AdminUserService for admin operations
	services.AdminUserService = service.NewAdminUserService(repos.UserRepo)
	services.AdminCommunityService = service.NewAdminCommunityService(repos.CommunityRepo)
	services.AdminStatsService = service.NewAdminStatsService(repos.UserRepo, repos.CommunityRepo, repos.PostRepo, repos.CommentRepo, repos.ReportRepo)
	services.AdminAuthService = service.NewAdminAuthService(repos.UserRepo, redisClient, tokenService)

	return services
}

func initControllers(services *Services, wsHub *ws.Hub, db *mongo.Database) *Controllers {
	return &Controllers{
		AuthController:           *controller.NewAuthController(services.AuthService),
		UserController:           *controller.NewUserController(services.UserService),
		CommunityController:      *controller.NewCommunityController(services.CommunityService),
		MembershipController:     *controller.NewMembershipController(services.MembershipService),
		PostController:           *controller.NewPostController(services.PostService),
		VoteController:           *controller.NewVoteController(services.VoteService), // Added VoteController
		CommentController:        *controller.NewCommentController(services.CommentService),
		NotificationController:   *controller.NewNotificationController(services.NotificationService),
		WebSocketController:      *controller.NewWebSocketController(wsHub),
		ChannelController:        *controller.NewChannelController(services.ChannelService),
		MessageController:        *controller.NewMessageController(services.MessageService),
		PostHistoryController:    *controller.NewPostHistoryController(services.PostHistoryService),
		DraftController:          *controller.NewDraftController(services.DraftService),
		ReportController:         *controller.NewReportController(services.ReportService),
		AdminUserController:      *controller.NewAdminUserController(services.AdminUserService),
		AdminCommunityController: *controller.NewAdminCommunityController(services.AdminCommunityService),
		AdminStatsController:     *controller.NewAdminStatsController(services.AdminStatsService),
		AdminAuthController:      *controller.NewAdminAuthController(services.AdminAuthService),
		DebugController:          *controller.NewDebugController(db),
	}
}

func initRoutes(controllers *Controllers, r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// ⚠️ DEBUG ONLY - Remove in production
	r.POST("/debug/create-admin", controllers.DebugController.CreateAdminUser)

	api := r.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to LKForum API!"})
	})

	route.RegisterAuthRoutes(api, &controllers.AuthController, &controllers.UserController)
	route.RegisterUserRoutes(api, &controllers.UserController)
	route.RegisterCommunityRoutes(api, &controllers.CommunityController)
	route.RegisterMembershipRoutes(api, &controllers.MembershipController)
	route.RegisterPostRoutes(api, &controllers.PostController)
	route.RegisterVoteRoutes(api, &controllers.VoteController) // Added VoteRoutes
	route.RegisterCommentRoutes(api, &controllers.CommentController)
	route.RegisterNotificationRoutes(api, &controllers.NotificationController)
	route.RegisterWebSocketRoutes(api, &controllers.WebSocketController)
	route.RegisterChannelRoutes(api, &controllers.ChannelController)
	route.RegisterMessageRoutes(api, &controllers.MessageController)
	route.RegisterPostHistoryRoutes(api, &controllers.PostHistoryController)
	route.RegisterDraftRoutes(api, &controllers.DraftController)
	route.RegisterReportRoutes(api, &controllers.ReportController)
	route.RegisterAdminAuthRoutes(api, &controllers.AdminAuthController)
	route.RegisterAdminUserRoutes(api, &controllers.AdminUserController)
	route.RegisterAdminCommunityRoutes(api, &controllers.CommunityController, &controllers.AdminCommunityController)
	route.RegisterAdminReportRoutes(api, &controllers.ReportController)
	route.RegisterAdminStatsRoutes(api, &controllers.AdminStatsController)
}

func Init() (*gin.Engine, error) {
	config.LoadConfig()
	auth.InitGoogleOAuthConfig()

	redisClient := config.NewRedisClient()

	tokenService, err := InitializeTokenService(redisClient)
	if err != nil {
		log.Printf("Warning: Token invalidation service not available: %v\n", err)
	}

	client := config.NewMongoClient()
	db := client.Database(config.Cfg.DBName)
	router := gin.Default()
	router.MaxMultipartMemory = 100 << 20 // 100 MB

	router.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		for _, allowedOrigin := range config.Cfg.AllowedOrigins {
			if origin == allowedOrigin {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	eventBus := bus.NewEventBus()
	wsHub := ws.NewHub(eventBus)
	emailSender := email.NewSMTPSender()

	// Initialize Gemini client for content moderation
	geminiClient, err := gemini.NewGeminiClient(&config.Cfg.Gemini)
	if err != nil {
		log.Printf("Warning: Gemini client initialization failed: %v. Content moderation will be disabled.", err)
	}

	repos := initRepos(client, db)
	services := initServices(repos, redisClient, emailSender, eventBus, geminiClient, tokenService)
	controllers := initControllers(services, wsHub, db)

	// Inject userRepo into middleware for settings caching
	middleware.SetUserRepo(repos.UserRepo)

	initRoutes(controllers, router)

	// Start background services
	go wsHub.Start()
	services.ReputationService.Start()
	services.NotificationService.Start()
	services.MessageService.Start()
	services.ChannelService.Start()
	services.CommunityService.Start()
	services.ModerationService.Start() // Start content moderation service

	return router, nil
}
