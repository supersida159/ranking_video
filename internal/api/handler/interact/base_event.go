package handler

// import (
// 	_ "ranking_video/internal/models"
// 	service "ranking_video/internal/service_interact"
// 	apperror "ranking_video/pkg/app_error"
// 	helper "ranking_video/pkg/utils"
// 	"ranking_video/pkg/utils/app_context"
// 	"strconv"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"github.com/rs/zerolog/log"
// )

// // BaseEventHandler provides common handler functionality for events
// type BaseEventHandler[T service.EventEntity] struct {
// 	service service.EventServiceInterface[T]
// 	appCtx  app_context.AppContext
// }

// // Add handles adding a new event
// // @Summary      Add a new event(like, comment, share, view)
// // @Description  Receive event data from request body and publish it to kafka then store it to mysqlDB
// // @Tags         Events
// // @Accept       json
// // @Produce      json
// // @Param        event  body      object  true  "Event data"
// // @Success      200   {object}   helper.Response
// // @Failure      400   {object}   helper.Response
// // @Failure      500   {object}   helper.Response
// // @Router       /events [post]
// func (h *BaseEventHandler[T]) Add(ctx *gin.Context) {
// 	log.Info().Msg("Add event")

// 	var event T
// 	if err := ctx.ShouldBindJSON(&event); err != nil {
// 		helper.BuildErrorGinResponse(ctx, apperror.ErrInvalidRequest(err))
// 		return
// 	}
// 	validator := h.appCtx.GetValidatetor()
// 	// Validate using app context validator
// 	if err := validator.Struct(event); err != nil {
// 		helper.BuildErrorGinResponse(ctx, apperror.ErrInvalidRequest(err))
// 		return
// 	}

// 	if err := h.service.Add(ctx, event); err != nil {
// 		helper.BuildErrorGinResponse(ctx, err)
// 		return
// 	}

// 	helper.BuildSuccessGinResponse(ctx, event)
// }

// // SoftDelete performs a soft delete on an event record

// // SoftDelete performs a soft delete on an event record
// // @Summary      Soft delete an event
// // @Description  Marks an event as deleted without removing it from the database
// // @Tags         Events
// // @Accept       json
// // @Produce      json
// // @Param        id   path        uint  true  "Event ID"
// // @Success      200  {object}    helper.Response
// // @Failure      400  {object}    helper.Response
// // @Failure      404  {object}    helper.Response
// // @Failure      500  {object}    helper.Response
// // @Router       /events/{id} [delete]
// func (h *BaseEventHandler[T]) SoftDelete(ctx *gin.Context) {
// 	log.Info().Msg("SoftDelete event")

// 	idStr := ctx.Param("id")
// 	id, err := strconv.ParseUint(idStr, 10, 32)
// 	if err != nil {
// 		helper.BuildErrorGinResponse(ctx, apperror.ErrInvalidRequest(err))
// 		return
// 	}

// 	if appErr := h.service.SoftDelete(ctx, uint(id)); appErr != nil {
// 		helper.BuildErrorGinResponse(ctx, appErr)
// 		return
// 	}

// 	helper.BuildSuccessGinResponse(ctx, gin.H{"message": "Event deleted successfully"})
// }

// // List retrieves a list of events with filtering, sorting, and pagination

// // List retrieves a list of events with filtering, sorting, and pagination
// // @Summary      List events
// // @Description  Get a list of events with filtering, sorting, and pagination
// // @Tags         Events
// // @Accept       json
// // @Produce      json
// // @Param        page        query   int     false  "Page number"
// // @Param        limit       query   int     false  "Number of items per page"
// // @Param        sort        query   string  false  "Sort fields (e.g. -created_at,+name for DESC and ASC)"
// // @Param        start_time  query   string  false  "Filter by start time"
// // @Param        end_time    query   string  false  "Filter by end time"
// // @Param        field_name  query   string  false  "Filter by any field name (can use >, <, >=, <= prefixes or CSV for IN clause)"
// // @Success      200  {object}  helper.Response{data=array,paging=helper.Paging}
// // @Failure      400  {object}  helper.Response
// // @Failure      500  {object}  helper.Response
// // @Router       /events [get]
// func (h *BaseEventHandler[T]) List(ctx *gin.Context) {
// 	log.Info().Msg("List events")

// 	// Parse query parameters for conditions
// 	conditions := parseQueryConditions(ctx)

// 	// Parse pagination parameters
// 	var paging helper.Paging
// 	paging.Limit, _ = strconv.Atoi(ctx.Query("limit"))
// 	paging.Page, _ = strconv.Atoi(ctx.Query("page"))
// 	paging.Fullfill() // Set defaults for empty fields

// 	// Parse sort parameters
// 	orderClauses := parseSortParameters(ctx)

// 	// Get count for pagination
// 	count, err := h.service.Count(ctx, conditions)
// 	if err != nil {
// 		helper.BuildErrorGinResponse(ctx, err)
// 		return
// 	}
// 	paging.Total = count

// 	// Get events list
// 	events, err := h.service.GetList(ctx, &paging, orderClauses, conditions)
// 	if err != nil {
// 		helper.BuildErrorGinResponse(ctx, err)
// 		return
// 	}

// 	data := map[string]interface{}{
// 		"data":   events,
// 		"paging": paging,
// 	}

// 	helper.BuildSuccessGinResponse(ctx, data)
// }

// // Helper functions

// // parseQueryConditions parses query parameters into a conditions map
// func parseQueryConditions(ctx *gin.Context) map[string]interface{} {
// 	queryParams := ctx.Request.URL.Query()
// 	conditions := make(map[string]interface{})

// 	for key, values := range queryParams {
// 		if len(values) > 0 {
// 			value := values[0]

// 			// Skip pagination, sort, and preload parameters
// 			if key == "page" || key == "limit" || key == "sort" || key == "preload" ||
// 				key == "start_time" || key == "end_time" {
// 				continue
// 			}

// 			// Check for comparison operators
// 			switch {
// 			case strings.HasPrefix(value, "<="):
// 				conditions[key+" <="] = strings.TrimPrefix(value, "<=")
// 			case strings.HasPrefix(value, ">="):
// 				conditions[key+" >="] = strings.TrimPrefix(value, ">=")
// 			case strings.HasPrefix(value, "<"):
// 				conditions[key+" <"] = strings.TrimPrefix(value, "<")
// 			case strings.HasPrefix(value, ">"):
// 				conditions[key+" >"] = strings.TrimPrefix(value, ">")
// 			case strings.Contains(value, ","):
// 				// Handle IN clause
// 				conditions[key] = strings.Split(value, ",")
// 			default:
// 				// Default to equality
// 				conditions[key] = value
// 			}
// 		}
// 	}

// 	return conditions
// }

// // parseSortParameters parses sort parameters into order clauses
// func parseSortParameters(ctx *gin.Context) []string {
// 	sortParams := ctx.QueryArray("sort")
// 	var orderClauses []string

// 	for _, sortParam := range sortParams {
// 		fields := strings.Split(sortParam, ",")
// 		for _, field := range fields {
// 			field = strings.TrimSpace(field)
// 			if field == "" {
// 				continue
// 			}
// 			var sortStr string
// 			switch {
// 			case strings.HasPrefix(field, "-"):
// 				sortStr = strings.TrimPrefix(field, "-") + " DESC"
// 			case strings.HasPrefix(field, "+"):
// 				sortStr = strings.TrimPrefix(field, "+") + " ASC"
// 			default:
// 				sortStr = field + " ASC"
// 			}
// 			orderClauses = append(orderClauses, sortStr)
// 		}
// 	}

// 	return orderClauses
// }
