package handler

import (
	"net/http"

	service "ranking_video/internal/interact_bussiness"
	"ranking_video/internal/models"
	apperror "ranking_video/pkg/app_error"
	helper "ranking_video/pkg/utils"
	"ranking_video/pkg/utils/app_context"

	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// BaseEventHandler provides common handler functionality for events
type BaseEventHandler[T any] struct {
	service service.EventServiceInterface[T]
	appCtx  app_context.AppContext
}

// Add creates a new event record
func (h *BaseEventHandler[T]) Add(ctx *gin.Context) {
	log.Info().Msg("Add event")

	var event T
	if err := ctx.ShouldBindJSON(&event); err != nil {
		helper.BuildErrorGinResponse(ctx, apperror.ErrUsernameTooShort())
		return
	}

	// Validate using app context validator
	if err := h.appCtx.GetValidator().ValidateStruct(event); err != nil {
		helper.BuildErrorGinResponse(ctx, err)
		return
	}

	if err := h.service.Add(ctx, event); err != nil {
		helper.BuildErrorGinResponse(ctx, err)
		return
	}

	helper.BuildSuccessGinResponse(ctx, event)
}

// SoftDelete performs a soft delete on an event record
func (h *BaseEventHandler[T]) SoftDelete(ctx *gin.Context) {
	log.Info().Msg("SoftDelete event")

	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helper.BuildErrorGinResponse(ctx, apperror.ErrInvalidRequest(err))
		return
	}

	if appErr := h.service.SoftDelete(ctx, uint(id)); err != nil {
		helper.BuildErrorGinResponse(ctx, appErr)
		return
	}

	helper.BuildSuccessGinResponse(ctx, gin.H{"message": "Event deleted successfully"})
}

// List retrieves a list of events with filtering, sorting, and pagination
func (h *BaseEventHandler[T]) List(ctx *gin.Context) {
	log.Info().Msg("List events")

	// Parse query parameters for conditions
	conditions := parseQueryConditions(ctx)

	// Parse pagination parameters
	var paging helper.Paging
	paging.Limit, _ = strconv.Atoi(ctx.Query("limit"))
	paging.Page, _ = strconv.Atoi(ctx.Query("page"))
	paging.Fullfill() // Set defaults for empty fields

	// Parse sort parameters
	orderClauses := parseSortParameters(ctx)

	// Get count for pagination
	count, err := h.service.Count(ctx, conditions)
	if err != nil {
		helper.BuildErrorGinResponse(ctx, err)
		return
	}
	paging.Total = count

	// Get events list
	events, err := h.service.GetList(ctx, &paging, orderClauses, conditions)
	if err != nil {
		helper.BuildErrorGinResponse(ctx, err)
		return
	}

	data := map[string]interface{}{
		"data":   events,
		"paging": paging,
	}

	helper.BuildSuccessGinResponse(ctx, data)
}

// ViewEventHandler handles HTTP requests for view events
type ViewEventHandler struct {
	BaseEventHandler[models.ViewEvent]
}

// NewViewEventHandler creates a new view event handler
func NewViewEventHandler(service service.EventServiceInterface[models.ViewEvent], appCtx app_context.AppContext) *ViewEventHandler {
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
func NewLikeEventHandler(service service.EventServiceInterface[models.LikeEvent], appCtx app_context.AppContext) *LikeEventHandler {
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
func NewCommentEventHandler(service service.EventServiceInterface[models.CommentEvent], appCtx app_context.AppContext) *CommentEventHandler {
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
func NewShareEventHandler(service service.EventServiceInterface[models.ShareEvent], appCtx app_context.AppContext) *ShareEventHandler {
	return &ShareEventHandler{
		BaseEventHandler: BaseEventHandler[models.ShareEvent]{
			service: service,
			appCtx:  appCtx,
		},
	}
}

// WatchEventHandler handles HTTP requests for watch events
type WatchEventHandler struct {
	BaseEventHandler[models.WatchEvent]
	watchService service.WatchEventServiceInterface[models.WatchEvent]
}

// NewWatchEventHandler creates a new watch event handler
func NewWatchEventHandler(service service.WatchEventServiceInterface[models.WatchEvent], appCtx app_context.AppContext) *WatchEventHandler {
	return &WatchEventHandler{
		BaseEventHandler: BaseEventHandler[models.WatchEvent]{
			service: service,
			appCtx:  appCtx,
		},
		watchService: service,
	}
}

// CalculateWatchTime calculates the total watch time for a video
func (h *WatchEventHandler) CalculateWatchTime(ctx *gin.Context) {
	log.Info().Msg("Calculate watch time")

	videoID := ctx.Param("videoId")
	if videoID == "" {
		helper.BuildErrorGinResponse(ctx, helper.NewError("video_id required", http.StatusBadRequest))
		return
	}

	// Parse additional conditions
	conditions := parseQueryConditions(ctx)

	watchTime, err := h.watchService.CalculateWatchTime(ctx, videoID, conditions)
	if err != nil {
		helper.BuildErrorGinResponse(ctx, err)
		return
	}

	data := map[string]interface{}{
		"video_id":   videoID,
		"watch_time": watchTime,
	}

	helper.BuildSuccessGinResponse(ctx, data)
}

// CalculateWatchTimeByTimeRange calculates watch time within a specific time range
func (h *WatchEventHandler) CalculateWatchTimeByTimeRange(ctx *gin.Context) {
	log.Info().Msg("Calculate watch time by time range")

	videoID := ctx.Param("videoId")
	if videoID == "" {
		helper.BuildErrorGinResponse(ctx, apperror.ErrInvalidRequest(nil))
		return
	}

	// Parse time range parameters
	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")

	var startTime, endTime time.Time
	var err error

	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			helper.BuildErrorGinResponse(ctx, apperror.ErrInvalidTimeFormat(err))
			return
		}
	} else {
		startTime = time.Now().AddDate(0, -1, 0) // Default to 1 month ago
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			helper.BuildErrorGinResponse(ctx, apperror.ErrInvalidTimeFormat(err))
			return
		}
	} else {
		endTime = time.Now() // Default to now
	}

	// Parse additional conditions
	conditions := parseQueryConditions(ctx)

	watchTime, appErr := h.watchService.CalculateWatchTimeByTimeRange(ctx, videoID, startTime, endTime, conditions)
	if appErr != nil {
		helper.BuildErrorGinResponse(ctx, appErr)
		return
	}

	data := map[string]interface{}{
		"video_id":   videoID,
		"start_time": startTime,
		"end_time":   endTime,
		"watch_time": watchTime,
	}

	helper.BuildSuccessGinResponse(ctx, data)
}

// GetWatchTimeByUser retrieves watch time grouped by user
func (h *WatchEventHandler) GetWatchTimeByUser(ctx *gin.Context) {
	log.Info().Msg("Get watch time by user")

	videoID := ctx.Param("videoId")
	if videoID == "" {
		helper.BuildErrorGinResponse(ctx, helper.NewError("video_id required", http.StatusBadRequest))
		return
	}

	watchTimeByUser, err := h.watchService.GetWatchTimeByUser(ctx, videoID)
	if err != nil {
		helper.BuildErrorGinResponse(ctx, err)
		return
	}

	data := map[string]interface{}{
		"video_id":         videoID,
		"watch_time_users": watchTimeByUser,
	}

	helper.BuildSuccessGinResponse(ctx, data)
}

// Helper functions

// parseQueryConditions parses query parameters into a conditions map
func parseQueryConditions(ctx *gin.Context) map[string]interface{} {
	queryParams := ctx.Request.URL.Query()
	conditions := make(map[string]interface{})

	for key, values := range queryParams {
		if len(values) > 0 {
			value := values[0]

			// Skip pagination, sort, and preload parameters
			if key == "page" || key == "limit" || key == "sort" || key == "preload" ||
				key == "start_time" || key == "end_time" {
				continue
			}

			// Check for comparison operators
			switch {
			case strings.HasPrefix(value, "<="):
				conditions[key+" <="] = strings.TrimPrefix(value, "<=")
			case strings.HasPrefix(value, ">="):
				conditions[key+" >="] = strings.TrimPrefix(value, ">=")
			case strings.HasPrefix(value, "<"):
				conditions[key+" <"] = strings.TrimPrefix(value, "<")
			case strings.HasPrefix(value, ">"):
				conditions[key+" >"] = strings.TrimPrefix(value, ">")
			case strings.Contains(value, ","):
				// Handle IN clause
				conditions[key] = strings.Split(value, ",")
			default:
				// Default to equality
				conditions[key] = value
			}
		}
	}

	return conditions
}

// parseSortParameters parses sort parameters into order clauses
func parseSortParameters(ctx *gin.Context) []string {
	sortParams := ctx.QueryArray("sort")
	var orderClauses []string

	for _, sortParam := range sortParams {
		fields := strings.Split(sortParam, ",")
		for _, field := range fields {
			field = strings.TrimSpace(field)
			if field == "" {
				continue
			}
			var sortStr string
			switch {
			case strings.HasPrefix(field, "-"):
				sortStr = strings.TrimPrefix(field, "-") + " DESC"
			case strings.HasPrefix(field, "+"):
				sortStr = strings.TrimPrefix(field, "+") + " ASC"
			default:
				sortStr = field + " ASC"
			}
			orderClauses = append(orderClauses, sortStr)
		}
	}

	return orderClauses
}
