package service

import (
	"context"
	"ranking_video/internal/repository"
	"time"
)

// WatchEventServiceInterface defines methods specific to watch events
type WatchEventServiceInterface[T any] interface {
	EventServiceInterface[T]
	CalculateWatchTime(ctx context.Context, videoID string, conditions map[string]interface{}) (int, error)
	CalculateWatchTimeByTimeRange(ctx context.Context, videoID string, startTime time.Time, endTime time.Time, additionalConditions map[string]interface{}) (int, error)
	GetWatchTimeByUser(ctx context.Context, videoID string) (map[string]int, error)
}

// WatchEventService handles business logic for watch events
type WatchEventService[T any] struct {
	*EventService[T]
	watchRepo *repository.WatchEventRepository[T]
}

// NewWatchEventService creates a new watch event service
func NewWatchEventService[T any](watchRepo *repository.WatchEventRepository[T]) *WatchEventService[T] {
	return &WatchEventService[T]{
		EventService: NewEventService[T](watchRepo.EventRepository),
		watchRepo:    watchRepo,
	}
}

// CalculateWatchTime calculates the total watch time for a video with optional conditions
func (s *WatchEventService[T]) CalculateWatchTime(
	ctx context.Context,
	videoID string,
	conditions map[string]interface{},
) (int, error) {
	return s.watchRepo.CalculateWatchTime(ctx, videoID, conditions)
}

// CalculateWatchTimeByTimeRange calculates watch time within a specific time range
func (s *WatchEventService[T]) CalculateWatchTimeByTimeRange(
	ctx context.Context,
	videoID string,
	startTime time.Time,
	endTime time.Time,
	additionalConditions map[string]interface{},
) (int, error) {
	return s.watchRepo.CalculateWatchTimeByTimeRange(ctx, videoID, startTime, endTime, additionalConditions)
}

// GetWatchTimeByUser calculates watch time grouped by user
func (s *WatchEventService[T]) GetWatchTimeByUser(
	ctx context.Context,
	videoID string,
) (map[string]int, error) {
	return s.watchRepo.GetWatchTimeByUser(ctx, videoID)
}
