package main

import (
	"fmt" // Added for port configuration
	route "ranking_video/internal/api/route/inetract"
	"ranking_video/internal/kafka"
	"ranking_video/internal/models"
	config "ranking_video/pkg/config/env_config"
	kafkaconfig "ranking_video/pkg/config/kafka"
	"ranking_video/pkg/database"
	"ranking_video/pkg/localredis"
	"ranking_video/pkg/pubsub/pubsublocal"
	"ranking_video/pkg/utils/app_context"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ranking_video/docs" // This is for Swagger documentation generation

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	// _ "ranking_video/docs/swagger" // Uncomment if you have a specific swagger package
	// "ranking_video/pkg/middleware" // Uncomment if you have middleware to use
)

// @title           Ranking Video API
// @version         1.0
// @description     API Server for Ranking Video application
// @termsOfService  http://swagger.io/terms/

// @contact.name
// @contact.url    http://www.example.com/support
// @contact.email  96duongtung@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8089
// @BasePath  /api/v1

func main() {

	// If flags indicate actions that shouldn't run the server, you might want to os.Exit(0) here.

	log.Info().Msg("Starting Server Setup...")

	// --- Configuration Loading ---
	cfg := config.LoadConfig() // Renamed for clarity

	// --- Database Connection ---
	db, err := database.NewDatabase(cfg.DatabaseURI)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	err = db.AutoMigrate(
		models.Entity{},
		models.CommentEvent{},
		models.LikeEvent{},
		models.ShareEvent{},
		models.ViewEvent{},
		models.Video{},
	)
	// Consider running migrations if isMigrate is true
	// database.Migrate(db) // Example

	// --- PubSub Initialization ---
	pubsub := pubsublocal.NewPubSub()

	// --- Redis Connection ---
	redisConfig := localredis.Config{
		Address:  cfg.Redis.Address,
		Password: cfg.Redis.Password,
	}
	redisClient := localredis.NewRedis(redisConfig) // Renamed for clarity

	// --- Kafka Connection ---
	// It seems you intended to create both a producer and a consumer.
	// Let's assume kafka.NewProducer creates a producer and kafka.NewConsumer creates a consumer.
	kafkaConfig := kafkaconfig.NewKafkaConfig(cfg.Kafka.Broker)

	producer, err := kafka.NewProducer(kafkaConfig) // Assuming this creates a producer
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Kafka producer")
	}
	// defer producer.Close() // Important to close resources

	consumer, err := kafka.NewConsumer(kafkaConfig) // Assuming this creates a consumer
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Kafka consumer")
	}
	// defer consumer.Close() // Important to close resources
	// You likely need to start the consumer listening in a separate goroutine.
	// go consumer.StartListening() // Example

	// --- Application Context ---
	// Pass the correctly initialized components
	appCtx := app_context.NewAppContext(db, pubsub, cfg, redisClient, producer, consumer)

	// --- Gin Router Setup ---
	router := gin.Default()
	// Add middleware if needed (e.g., CORS, logging, recovery)
	// router.Use(cors.Default())
	// router.Use(gin.Recovery())

	// --- API Routing ---
	v1 := router.Group("/api/v1")
	{ // Using braces for visual grouping of v1 routes
		route.CommentRoute(v1, appCtx)
		route.LikeRoute(v1, appCtx)
		route.ShareRoute(v1, appCtx)
		route.ViewRoute(v1, appCtx)
	}

	// Add a simple health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --- Start Server ---
	serverPort := "8089"                                     // Define the port
	serverAddress := fmt.Sprintf("localhost:%s", serverPort) // Format for ListenAndServe/Run

	log.Info().Msgf("Server starting on port %s", serverPort)

	// This is the crucial line to start the server
	err = router.Run(serverAddress)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to start server on port %s", serverPort)
	}

	// Code here will only run if router.Run returns (which it normally doesn't unless there's an error)
	log.Info().Msg("Server gracefully shut down") // Or "Server encountered an error" if err != nil
}
