package bootstrap

import (
	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/repo"
	route "github.com/giakiet05/lkforum/internal/route/user"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
)

type Repos struct {
	repo.UserRepo
}

type Services struct {
	service.UserService
}

type Controllers struct {
	controller.UserController
}

// initRepos initializes repositories with the given database
func initRepos(db *mongo.Database) *Repos {
	return &Repos{
		UserRepo: repo.NewUserRepo(db),
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
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := r.Group("/api")

	route.RegisterUserRoutes(api, &controllers.UserController)
}

// Init inits the app
func Init() *gin.Engine {
	//redisClient := config.NewRedisClient()
	mongoClient := config.NewMongoClient()           //Connect to MongoDB (contains many databases)
	db := mongoClient.Database(os.Getenv("DB_NAME")) //Choose a specific database

	repos := initRepos(db)
	services := initServices(repos)
	controllers := initControllers(services)

	r := gin.Default()
	initRoutes(controllers, r)

	return r
}
