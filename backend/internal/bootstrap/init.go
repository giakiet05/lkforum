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
	"go.mongodb.org/mongo-driver/mongo"
)

type Repos struct {
	repo.UserRepo
	repo.CommunityRepo
}

type Services struct {
	service.UserService
	service.CommunityService
}

type Controllers struct {
	controller.UserController
}

// initRepos initializes repositories with the given database
func initRepos(db *mongo.Database) *Repos {
	return &Repos{
		UserRepo:      repo.NewUserRepo(db),
		CommunityRepo: repo.NewCommunityRepo(db),
	}
}

// initServices Initialize services with the given repositories
func initServices(repos *Repos) *Services {
	return &Services{
		UserService: service.NewUserService(repos.UserRepo),
	}
}

// initControllers Initialize controllers with the given services
func initControllers(services *Services) *Controllers {
	return &Controllers{
		UserController: *controller.NewUserController(services.UserService),
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

	// Initialize other components
	repos := initRepos(db)
	services := initServices(repos)
	controllers := initControllers(services)
	initRoutes(controllers, router)

	// Setup router

	return router, nil
}
