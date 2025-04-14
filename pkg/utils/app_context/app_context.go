package app_context

import (
	dbs "ranking_video/pkg/db"
	"ranking_video/pkg/localredis"

	"gorm.io/gorm"

	"github.com/go-playground/validator/v10"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetSecretKey() string
	GetPubSub() pubsub.PubSub
	GetCache() *localredis.RedisWRealStore
	GetConfig() *config.Schema
	GetValidatetor() *common.Validator
	GetProducer() *producers.OrderProducer
	GetConsumer() *consumerlocal.SagaConsumer
}

type AppCtx struct {
	Dbs       *dbs.Database
	Pb        pubsub.PubSub
	Cfg       *config.Schema
	Cache     *localredis.RedisWRealStore
	Validator *validator.Validate
	Producer  *producers.OrderProducer
	Consumer  *consumerlocal.SagaConsumer
}

func NewAppContext(dbs *dbs.Database,
	pb pubsub.PubSub,
	cache *localredis.RedisWRealStore,
	producer *producers.OrderProducer,
	consumer *consumerlocal.SagaConsumer) *AppCtx {
	return &AppCtx{
		Dbs:       dbs,
		Pb:        pb,
		Cfg:       config.GetConfig(),
		Cache:     cache,
		Validator: validator.New(),
		Producer:  producer,
		Consumer:  consumer,
	}
}

func (ctx *AppCtx) GetMainDBConnection() *gorm.DB {
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

func (ctx *AppCtx) GetProducer() *producers.OrderProducer {
	return ctx.Producer
}

func (ctx *AppCtx) GetConsumer() *consumerlocal.SagaConsumer {
	return ctx.Consumer
}

func (ctx *AppCtx) GetValidatetor() *validator.Validate {
	return ctx.Validator
}
