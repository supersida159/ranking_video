package app_context

import (
	"ranking_video/internal/kafka"
	config "ranking_video/pkg/config/env_config"
	dbs "ranking_video/pkg/database"
	"ranking_video/pkg/localredis"
	"ranking_video/pkg/pubsub"

	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
)

type AppContext interface {
	GetDBConnection() *gorm.DB
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

func (ctx *AppCtx) GetDBConnection() *gorm.DB {
	return ctx.Dbs.GetDB()
}

func (ctx *AppCtx) GetSecretKey() string {
	return ctx.Cfg.AuthSecret
}

func (ctx *AppCtx) GetPubSub() pubsub.PubSub {
	return ctx.Pb
}

func (ctx *AppCtx) GetCache() *localredis.RedisWRealStore {
	return ctx.Cache
}

func (ctx *AppCtx) GetConfig() *config.Schema {
	return ctx.Cfg
}

func (ctx *AppCtx) GetProducer() *kafka.Producer {
	return ctx.Producer
}

func (ctx *AppCtx) GetConsumer() *kafka.Consumer {
	return ctx.Consumer
}

func (ctx *AppCtx) GetValidatetor() *validator.Validate {
	return ctx.Validator
}
