package service

import (
	"context"
	"ranking_video/internal/kafka"
	apperror "ranking_video/pkg/app_error"
	"time"
)

// WatchEventServiceInterface defines methods specific to watch events
type WatchEventServiceInterface[T any] interface {
	EventServiceInterface[T]
	CalculateWatchTime(ctx context.Context, videoID string, conditions map[string]interface{}) (int, *apperror.AppError)
	CalculateWatchTimeByTimeRange(ctx context.Context,
		videoID string,
		startTime time.Time,
		endTime time.Time,
		additionalConditions map[string]interface{}) (int, *apperror.AppError)
	GetWatchTimeByUser(ctx context.Context, videoID string) (map[string]int, *apperror.AppError)
}

// WatchEventService handles business logic for watch events
type WatchEventService[T any] struct {
	*EventService[T]
	watchStorage WatchEventServiceInterface[T]
}

// NewWatchEventService creates a new watch event service
func NewWatchEventService[T any](evenStorage WatchEventService[T], producer *kafka.Producer) *WatchEventService[T] {
	return &WatchEventService[T]{
		EventService: NewEventService(evenStorage.storage, producer),
		watchStorage: evenStorage.watchStorage,
	}
}

// CalculateWatchTime calculates the total watch time for a video with optional conditions
func (s *WatchEventService[T]) CalculateWatchTime(
	ctx context.Context,
	videoID string,
	conditions map[string]interface{},
) (int, error) {
	return s.watchStorage.CalculateWatchTime(ctx, videoID, conditions)
}

// CalculateWatchTimeByTimeRange calculates watch time within a specific time range
func (s *WatchEventService[T]) CalculateWatchTimeByTimeRange(
	ctx context.Context,
	videoID string,
	startTime time.Time,
	endTime time.Time,
	additionalConditions map[string]interface{},
) (int, error) {
	return s.watchStorage.CalculateWatchTimeByTimeRange(ctx, videoID, startTime, endTime, additionalConditions)
}

// GetWatchTimeByUser calculates watch time grouped by user
func (s *WatchEventService[T]) GetWatchTimeByUser(
	ctx context.Context,
	videoID string,
) (map[string]int, error) {
	return s.watchStorage.GetWatchTimeByUser(ctx, videoID)
}
