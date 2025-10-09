package bootstrap

import (
	"log"
	"os"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/repo"
	route "github.com/giakiet05/lkforum/internal/route/user"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repos struct {
	repo.UserRepo
	repo.CommunityRepo
	repo.MembershipRepo
}

type Services struct {
	service.UserService
	service.CommunityService
	service.MembershipService
}

type Controllers struct {
	controller.UserController
	controller.CommunityController
	controller.MembershipController
}

// initRepos initializes repositories with the given database
func initRepos(db *mongo.Database) *Repos {
	return &Repos{
		UserRepo:       repo.NewUserRepo(db),
		CommunityRepo:  repo.NewCommunityRepo(db),
		MembershipRepo: repo.NewMembershipRepo(db),
	}
}

// initServices Initialize services with the given repositories
func initServices(repos *Repos, redisClient *redis.Client) *Services {
	return &Services{
		UserService:       service.NewUserService(repos.UserRepo),
		CommunityService:  service.NewCommunityService(repos.CommunityRepo),
		MembershipService: service.NewMembershipService(repos.MembershipRepo, redisClient),
	}
}

// initControllers Initialize controllers with the given services
func initControllers(services *Services) *Controllers {
	return &Controllers{
		UserController:       *controller.NewUserController(services.UserService),
		CommunityController:  *controller.NewCommunityController(services.CommunityService),
		MembershipController: *controller.NewMembershipController(services.MembershipService),
	}
}

// initRoutes sets up the routes for the Gin engine
func initRoutes(controllers *Controllers, r *gin.Engine) {
	//Test route
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	//Test API group
	api := r.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome to LKForum API!"})
	})

	//Register more routes here
	route.RegisterAuthRoutes(api, &controllers.UserController)
	route.RegisterUserRoutes(api, &controllers.UserController)
	route.RegisterCommunityRoutes(api, &controllers.CommunityController)
	route.RegisterMembershipRoutes(api, &controllers.MembershipController)
}

// Init initializes all application components
func Init() (*gin.Engine, error) {
	// Create Redis client once
	redisClient := config.NewRedisClient()

	// Initialize token service first for JWT blacklisting
	if err := InitializeTokenService(redisClient); err != nil {
		// Log error but continue - the system will work without Redis, just without token invalidation
		log.Printf("Warning: Token invalidation service not available: %v\n", err)
	}

	// Connect to MongoDB
	client := config.NewMongoClient()
	db := client.Database(os.Getenv("DB_NAME"))
	router := gin.Default()

	// Register CORS middleware before any routes or other middleware
	allowOrigin := os.Getenv("FRONTEND_URL")
	if allowOrigin == "" {
		allowOrigin = "http://localhost:5173"
	}
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Initialize other components
	repos := initRepos(db)
	services := initServices(repos, redisClient)
	controllers := initControllers(services)
	initRoutes(controllers, router)

	// Setup router

	return router, nil
}
