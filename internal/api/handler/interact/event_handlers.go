package handler

import (
	"ranking_video/internal/models"
	"ranking_video/internal/repository"
	service "ranking_video/internal/service_interact"
	"ranking_video/pkg/utils/app_context"
)

// ViewEventHandler handles HTTP requests for view events
type ViewEventHandler struct {
	BaseEventHandler[models.ViewEvent]
}

// NewViewEventHandler creates a new view event handler
func NewViewEventHandler(appCtx app_context.AppContext) *ViewEventHandler {
	dbs := appCtx.GetMongoDatabase()
	store := repository.NewEventRepository[models.ViewEvent](dbs)
	service := service.NewEventService[models.ViewEvent](store, appCtx.GetProducer())
	return &ViewEventHandler{
		BaseEventHandler: BaseEventHandler[models.ViewEvent]{
			service: service,
			appCtx:  appCtx,
		},
	}
}

// LikeEventHandler handles HTTP requests for like events
type LikeEventHandler struct {
	BaseEventHandler[models.LikeEvent]
}

// NewLikeEventHandler creates a new like event handler
func NewLikeEventHandler(appCtx app_context.AppContext) *LikeEventHandler {
	dbs := appCtx.GetMongoDatabase()
	store := repository.NewEventRepository[models.LikeEvent](dbs)
	service := service.NewEventService[models.LikeEvent](store, appCtx.GetProducer())
	return &LikeEventHandler{
		BaseEventHandler: BaseEventHandler[models.LikeEvent]{
			service: service,
			appCtx:  appCtx,
		},
	}
}

// CommentEventHandler handles HTTP requests for comment events
type CommentEventHandler struct {
	BaseEventHandler[models.CommentEvent]
}

// NewCommentEventHandler creates a new comment event handler
func NewCommentEventHandler(appCtx app_context.AppContext) *CommentEventHandler {
	dbs := appCtx.GetMongoDatabase()
	store := repository.NewEventRepository[models.CommentEvent](dbs)
	service := service.NewEventService[models.CommentEvent](store, appCtx.GetProducer())
	return &CommentEventHandler{
		BaseEventHandler: BaseEventHandler[models.CommentEvent]{
			service: service,
			appCtx:  appCtx,
		},
	}
}

// ShareEventHandler handles HTTP requests for share events
type ShareEventHandler struct {
	BaseEventHandler[models.ShareEvent]
}

// NewShareEventHandler creates a new share event handler
func NewShareEventHandler(appCtx app_context.AppContext) *ShareEventHandler {
	dbs := appCtx.GetMongoDatabase()
	store := repository.NewEventRepository[models.ShareEvent](dbs)
	service := service.NewEventService[models.ShareEvent](store, appCtx.GetProducer())
	return &ShareEventHandler{
		BaseEventHandler: BaseEventHandler[models.ShareEvent]{
			service: service,
			appCtx:  appCtx,
		},
	}
}
