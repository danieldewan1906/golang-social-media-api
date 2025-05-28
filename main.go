package main

import (
	"golang-restful-api/internal/api"
	"golang-restful-api/internal/config"
	"golang-restful-api/internal/connection"
	"golang-restful-api/internal/repository"
	"golang-restful-api/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cnf := config.Get()
	dbConnection := connection.GetDatabase(cnf.Database)
	redisConnection := connection.GetRedisConnection(cnf.Redis)

	app := fiber.New()

	// middleware
	app.Use(cors.New())

	// repository
	userRepository := repository.NewUser(dbConnection)
	userImageRepository := repository.NewUserImages(dbConnection)
	postRepository := repository.NewPostRepository(dbConnection)
	userDetailRepository := repository.NewUserDetail(dbConnection)
	likeRepository := repository.NewLikeRepository(dbConnection)
	commentRepository := repository.NewCommentRepository(dbConnection)
	followRepository := repository.NewFollowRepository(dbConnection)
	redisRepository := repository.NewRedisRepository(redisConnection)

	// service
	authService := service.NewAuth(dbConnection, cnf, userRepository, userDetailRepository)
	userImageService := service.NewUserImageService(dbConnection, userImageRepository)
	postService := service.NewPostService(dbConnection, postRepository, likeRepository, commentRepository)
	userDetailService := service.NewUserDetail(dbConnection, userDetailRepository, userImageService, followRepository, postService, redisRepository)
	likeService := service.NewLikeService(dbConnection, likeRepository, postService)
	commentService := service.NewCommentService(dbConnection, commentRepository, postService)
	followService := service.NewFollowService(dbConnection, followRepository, userDetailService)

	// define path api
	endpoint := app.Group("/api/v1")
	api.NewAuth(endpoint, authService)
	api.NewUserImages(endpoint, cnf, userImageService)
	api.NewUserDetail(endpoint, userDetailService)
	api.NewPostApi(endpoint, cnf, postService)
	api.NewLikeApi(endpoint, likeService)
	api.NewCommentApi(endpoint, commentService)
	api.NewFollowApi(endpoint, followService)

	// Log all registered routes
	log.Println("📋 Registered Routes:")
	for _, routes := range app.Stack() {
		for _, route := range routes {
			if route.Path != "/" {
				log.Printf("%s\t%s\n", route.Method, route.Path)
			}
		}
	}

	log.Printf("application running on port %s", cnf.Server.Port)
	// fmt.Println("application running on port", cnf.Server.Port)
	_ = app.Listen(cnf.Server.Host + ":" + cnf.Server.Port)
}
