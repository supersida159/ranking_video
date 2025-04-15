package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"ranking_video/internal/kafka"
	"ranking_video/internal/models"
	config "ranking_video/pkg/config/env_config"
	kafkaconfig "ranking_video/pkg/config/kafka"
	"ranking_video/pkg/database"
	"ranking_video/pkg/localredis"
	"ranking_video/pkg/pubsub"
	"ranking_video/pkg/pubsub/pubsublocal"
	"ranking_video/pkg/utils/app_context"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting Server Setup...")

	// --- Configuration Loading ---
	cfg := config.LoadConfig()

	// --- Database Connection ---
	db, err := database.NewDatabase(cfg.DatabaseURI)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	err = db.AutoMigrate(
		models.Video{},
		models.VideoCountDaily{},
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// --- PubSub Initialization ---
	newPubsub := pubsublocal.NewPubSub()

	// --- Redis Connection ---
	redisConfig := localredis.Config{
		Address:  cfg.Redis.Address,
		Password: cfg.Redis.Password,
	}
	redisClient := localredis.NewRedis(redisConfig)

	// --- Kafka Connection ---
	kafkaConfig := kafkaconfig.NewKafkaConfig(cfg.Kafka.Broker)

	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Kafka producer")
	}
	defer producer.Close()

	consumer, err := kafka.NewConsumer(kafkaConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Kafka consumer")
	}
	defer consumer.Close()

	// --- Application Context ---
	ctx := context.Background()
	appCtx := app_context.NewAppContext(db, newPubsub, cfg, redisClient, producer, consumer)
	fmt.Println("App context created successfully :", appCtx)

	msgHandler := func(msg kafka.Event) error {
		// Process the message
		log.Info().Msgf("Received message: %s", string(msg.EventType))
		newMsg := pubsub.NewMessage(msg)
		err = newPubsub.Publish(ctx, pubsub.TopicNewEvent, newMsg)
		if err != nil {
			log.Error().Err(err).Msg("Failed to publish message")
		}
		return nil
	}

	// Start the consumer in a non-blocking way
	if err := consumer.Consume(msgHandler); err != nil {
		log.Fatal().Err(err).Msg("Failed to start Kafka consumer")
	}

	log.Info().Msg("Application started successfully")

	// Set up a channel to listen for termination signals
	// This keeps the application running until it receives a termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a termination signal
	sig := <-sigChan
	log.Info().Msgf("Received signal %s, shutting down...", sig)
}
