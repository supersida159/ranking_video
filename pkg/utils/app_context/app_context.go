package app_context

import (
	"ranking_video/internal/kafka"
	config "ranking_video/pkg/config/env_config"
	dbs "ranking_video/pkg/database"
	"ranking_video/pkg/localredis"
	"ranking_video/pkg/pubsub"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-playground/validator/v10"
)

type AppContext interface {
	GetMongoDatabase() *mongo.Database
	GetMongoDB() *dbs.Database
	GetSecretKey() string
	GetPubSub() pubsub.PubSub
	GetCache() *localredis.RedisWRealStore
	GetConfig() *config.Schema
	GetValidatetor() *validator.Validate
	GetProducer() *kafka.Producer
	GetConsumer() *kafka.Consumer
}

type AppCtx struct {
	Dbs       *dbs.Database
	Pb        pubsub.PubSub
	Cfg       *config.Schema
	Cache     *localredis.RedisWRealStore
	Validator *validator.Validate
	Producer  *kafka.Producer
	Consumer  *kafka.Consumer
}

func NewAppContext(dbs *dbs.Database,
	pb pubsub.PubSub,
	config *config.Schema,
	cache *localredis.RedisWRealStore,
	producer *kafka.Producer,
	consumer *kafka.Consumer) *AppCtx {
	return &AppCtx{
		Dbs:       dbs,
		Pb:        pb,
		Cfg:       config,
		Cache:     cache,
		Validator: validator.New(),
		Producer:  producer,
		Consumer:  consumer,
	}
}

// GetMongoDatabase returns the MongoDB database instance
func (ctx *AppCtx) GetMongoDatabase() *mongo.Database {
	return ctx.Dbs.GetDatabase()
}

// GetMongoDB returns the Database wrapper
func (ctx *AppCtx) GetMongoDB() *dbs.Database {
	return ctx.Dbs
}

// GetSecretKey returns the authentication secret key
func (ctx *AppCtx) GetSecretKey() string {
	return ctx.Cfg.AuthSecret
}

// GetPubSub returns the pubsub instance
func (ctx *AppCtx) GetPubSub() pubsub.PubSub {
	return ctx.Pb
}

// GetCache returns the Redis cache instance
func (ctx *AppCtx) GetCache() *localredis.RedisWRealStore {
	return ctx.Cache
}

// GetConfig returns the application configuration
func (ctx *AppCtx) GetConfig() *config.Schema {
	return ctx.Cfg
}

// GetProducer returns the Kafka producer
func (ctx *AppCtx) GetProducer() *kafka.Producer {
	return ctx.Producer
}

// GetConsumer returns the Kafka consumer
func (ctx *AppCtx) GetConsumer() *kafka.Consumer {
	return ctx.Consumer
}

// GetValidatetor returns the validator instance
func (ctx *AppCtx) GetValidatetor() *validator.Validate {
	return ctx.Validator
}
