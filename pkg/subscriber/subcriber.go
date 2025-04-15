package subscriber

import (
	"context"
	"errors"
	"fmt"
	"log"
	"ranking_video/internal/kafka"
	"ranking_video/internal/models"
	services_processor "ranking_video/internal/processor"
	"ranking_video/internal/repository"
	apperror "ranking_video/pkg/app_error"
	"ranking_video/pkg/asyncjob"
	"ranking_video/pkg/pubsub"
	"ranking_video/pkg/utils/app_context"
	"strconv"
	"time"
)

const (
	// Redis key prefixes
	videoKeyPrefix = "video:%s"
	exVideoPrefix  = "exvideo:%s"
	// Default expiration time for video events
	defaultExpiration = 60 * time.Second
)

type consumerJob struct {
	Title string
	Hld   func(ctx context.Context, msg *pubsub.Message) *apperror.AppError
}

type consumerEngine struct {
	appCtx app_context.AppContext
}

// NewEngine initializes the consumer engine with the application context.
func NewEngine(appCtx app_context.AppContext) *consumerEngine {
	return &consumerEngine{
		appCtx: appCtx,
	}
}

// Start initializes all subscribers and starts the engine.
func (engine *consumerEngine) Start() *apperror.AppError {
	// Subscribe to message event topic
	if err := engine.startSubTopic(
		pubsub.TopicNewEvent,
		false,
		consumerJob{
			Title: "HandleMessageEvent",
			Hld:   engine.handleMessageEvent,
		},
	); err != nil {
		return apperror.ErrInvalidRequest(errors.New("failed to start topic subscription for new events"))
	}

	// Subscribe to expire event topic
	if err := engine.startSubTopic(
		pubsub.TopicExpireEvent,
		false,
		consumerJob{
			Title: "HandleExpiredVideoEvent",
			Hld:   engine.handleExpiredVideoEvent,
		},
	); err != nil {
		return apperror.ErrInvalidRequest(errors.New("failed to start topic subscription for expire events"))
	}

	// Set up Redis key event listener
	engine.ListenRedisEventKey(map[string]pubsub.Topic{
		"exvideo": pubsub.TopicExpireEvent,
	})

	log.Println("Consumer engine started successfully")
	return nil
}

// handleMessageEvent processes incoming event messages and updates Redis.
func (engine *consumerEngine) handleMessageEvent(ctx context.Context, msg *pubsub.Message) *apperror.AppError {
	log.Printf("Processing message event: %T, value: %v", msg.Data(), msg.Data())

	event, ok := msg.Data().(*kafka.Event)
	if !ok {
		return apperror.ErrInvalidRequest(errors.New("invalid event data type"))
	}

	videoID, ok := event.Data["video_id"].(string)
	if !ok {
		return apperror.ErrInvalidRequest(errors.New("invalid or missing video_id in event data"))
	}

	redisClient := engine.appCtx.GetCache()
	exVideoKey := fmt.Sprintf(exVideoPrefix, videoID)
	videoKey := fmt.Sprintf(videoKeyPrefix, videoID)

	// Check if the video expiration key exists
	exists, err := redisClient.Client.Exists(ctx, exVideoKey).Result()
	if err != nil {
		log.Printf("Error checking existence of key %s: %v", exVideoKey, err)
		return apperror.ErrInvalidRequest(err)
	}

	// If the key doesn't exist, create it with expiration
	if exists == 0 {
		log.Printf("Setting expiration for video %s", videoID)
		if err := redisClient.SetWithExpiration(exVideoKey, 1, defaultExpiration); err != nil {
			log.Printf("Error setting expiration for video %s: %v", videoID, err)
			return apperror.ErrInvalidRequest(err)
		}
	}

	// Increment event counter and score in a transaction
	pipe := redisClient.Client.Pipeline()
	pipe.HIncrBy(ctx, videoKey, event.EventType, 1)
	pipe.HIncrBy(ctx, videoKey, "score", 1)
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("Error incrementing counters for video %s: %v", videoID, err)
		return apperror.ErrInvalidRequest(err)
	}

	log.Printf("Successfully processed event %s for video %s", event.EventType, videoID)
	return nil
}

// handleExpiredVideoEvent processes expired video events and updates the database.
func (engine *consumerEngine) handleExpiredVideoEvent(ctx context.Context, msg *pubsub.Message) *apperror.AppError {
	log.Printf("Processing expired video event: %v", msg.Data())

	videoID, ok := msg.Data().(string)
	if !ok {
		return apperror.ErrInvalidRequest(errors.New("invalid video ID format"))
	}

	redisClient := engine.appCtx.GetCache()
	videoKey := fmt.Sprintf(videoKeyPrefix, videoID)

	// Get all counts from Redis hash
	countsMap, err := redisClient.Client.HGetAll(ctx, videoKey).Result()
	if err != nil {
		log.Printf("Error retrieving counts for video %s: %v", videoID, err)
		return apperror.ErrInvalidRequest(err)
	}

	// Skip processing if the hash doesn't exist
	if len(countsMap) == 0 {
		log.Printf("No data found for expired video: %s", videoID)
		return nil
	}

	// Create video model from the counts
	video := &models.Video{
		ID: videoID,
	}

	// Parse counts from Redis hash
	if err := parseCountsToVideo(countsMap, video); err != nil {
		log.Printf("Error parsing video counts: %v", err)
		// Continue processing to ensure Redis cleanup
	}

	// Delete the Redis hash
	if _, err := redisClient.Client.Del(ctx, videoKey).Result(); err != nil {
		log.Printf("Error deleting Redis key %s: %v", videoKey, err)
		// Continue to update DB anyway
	}

	// Update database with the video counts
	db := engine.appCtx.GetDBConnection()
	storage := repository.NewVideoCountRepository(db)
	videoCountService := services_processor.NewVideoCountService(storage)

	if err := videoCountService.UpdateCount(ctx, video); err != nil {
		log.Printf("Error updating video counts in database: %v", err)
		return err
	}

	log.Printf("Successfully processed expired video: %s", videoID)
	return nil
}

// parseCountsToVideo parses Redis hash values into a Video struct.
func parseCountsToVideo(countsMap map[string]string, video *models.Video) error {
	for field, value := range countsMap {
		count, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing value '%s' for field '%s': %w", value, field, err)
		}

		switch field {
		case "like":
			video.LikeCount = int(count)
		case "comment":
			video.CommentCount = int(count)
		case "share":
			video.ShareCount = int(count)
		case "view":
			video.ViewCount = int(count)
		case "score":
			video.Score = uint(count)
		}
	}
	return nil
}

// ListenRedisEventKey sets up a listener for Redis key expiration events.
func (engine *consumerEngine) ListenRedisEventKey(mapPrefixWithTopic map[string]pubsub.Topic) {
	redisClient := engine.appCtx.GetCache()
	pubsubClient := engine.appCtx.GetPubSub()

	// Configure Redis to notify on key expiration events
	if _, err := redisClient.Client.ConfigSet(context.Background(), "notify-keyspace-events", "Ex").Result(); err != nil {
		log.Printf("Error configuring Redis keyspace events: %v", err)
		return
	}

	// Start a goroutine to listen for Redis key expiration events
	go func() {
		ctx := context.Background()
		pubSubConn := redisClient.Client.Subscribe(ctx, "__keyevent@0__:expired")
		defer pubSubConn.Close()

		log.Println("Redis key expiration listener started")
		ch := pubSubConn.Channel()

		for msg := range ch {
			key := msg.Payload
			log.Printf("Received Redis expired key event: %s", key)

			for prefix, topic := range mapPrefixWithTopic {
				if len(key) > len(prefix) && key[:len(prefix)] == prefix {
					// Extract the video ID from the key (format: prefix:videoID)
					videoID := key[len(prefix)+1:] // +1 to skip the colon
					log.Printf("Publishing expired key event for video %s to topic %s", videoID, topic)

					// Create and publish a message with the videoID
					message := pubsub.NewMessage(videoID)
					if err := pubsubClient.Publish(ctx, topic, message); err != nil {
						log.Printf("Error publishing to topic %s: %v", topic, err)
					}
				}
			}
		}
	}()
}

// startSubTopic subscribes to a topic and processes messages.
func (engine *consumerEngine) startSubTopic(topic pubsub.Topic, isConcurrent bool, consumerJobs ...consumerJob) *apperror.AppError {
	ctx := context.Background()
	c, err := engine.appCtx.GetPubSub().Subscribe(ctx, topic)
	if err != nil {
		return apperror.ErrInvalidRequest(fmt.Errorf("failed to subscribe to topic %s: %w", topic, err))
	}

	for _, job := range consumerJobs {
		log.Printf("Set up consumer for: %s on topic: %s", job.Title, topic)
	}

	getJobHandler := func(job *consumerJob, msg *pubsub.Message) asyncjob.JobHandler {
		return func(ctx context.Context) *apperror.AppError {
			log.Printf("Running job %s for topic %s with message: %v", job.Title, topic, msg.Data())
			return job.Hld(ctx, msg)
		}
	}

	go func() {
		for msg := range c {
			jobHandlers := make([]asyncjob.Job, len(consumerJobs))
			for i := range consumerJobs {
				jobHandlers[i] = asyncjob.NewJob(getJobHandler(&consumerJobs[i], msg))
			}

			group := asyncjob.NewGroup(isConcurrent, jobHandlers...)
			if err := group.Run(ctx); err != nil {
				log.Printf("Error in async job group for topic %s: %v", topic, err)
			}
		}
	}()

	return nil
}
