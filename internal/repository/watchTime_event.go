package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// WatchEventRepository handles CRUD operations for watch events
type WatchEventRepository[T any] struct {
	*EventRepository[T]
}

// NewWatchEventRepository creates a new watch event repository
func NewWatchEventRepository[T any](db *gorm.DB) *WatchEventRepository[T] {
	return &WatchEventRepository[T]{
		EventRepository: NewEventRepository[T](db),
	}
}

// CalculateWatchTime calculates the total watch time for a video with optional conditions
func (r *WatchEventRepository[T]) CalculateWatchTime(
	ctx context.Context,
	videoID string,
	conditions map[string]interface{},
) (int, error) {
	var totalWatchTime int

	query := r.db.WithContext(ctx).
		Model(new(T)).
		Select("COALESCE(SUM(duration), 0) as total_watch_time").
		Where("video_id = ?", videoID)

	// Apply additional conditions
	for field, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	if err := query.Scan(&totalWatchTime).Error; err != nil {
		return 0, err
	}

	return totalWatchTime, nil
}

// CalculateWatchTimeByTimeRange calculates watch time within a specific time range
func (r *WatchEventRepository[T]) CalculateWatchTimeByTimeRange(
	ctx context.Context,
	videoID string,
	startTime time.Time,
	endTime time.Time,
	additionalConditions map[string]interface{},
) (int, error) {
	var totalWatchTime int

	query := r.db.WithContext(ctx).
		Model(new(T)).
		Select("COALESCE(SUM(duration), 0) as total_watch_time").
		Where("video_id = ?", videoID).
		Where("created_at BETWEEN ? AND ?", startTime, endTime)

	// Apply additional conditions
	for field, value := range additionalConditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	if err := query.Scan(&totalWatchTime).Error; err != nil {
		return 0, err
	}

	return totalWatchTime, nil
}

// GetWatchTimeByUser calculates watch time grouped by user
func (r *WatchEventRepository[T]) GetWatchTimeByUser(
	ctx context.Context,
	videoID string,
) (map[string]int, error) {
	type Result struct {
		UserID         string
		TotalWatchTime int
	}

	var results []Result

	if err := r.db.WithContext(ctx).
		Model(new(T)).
		Select("user_id, COALESCE(SUM(duration), 0) as total_watch_time").
		Where("video_id = ?", videoID).
		Group("user_id").
		Scan(&results).Error; err != nil {
		return nil, err
	}

	watchTimeByUser := make(map[string]int)
	for _, result := range results {
		watchTimeByUser[result.UserID] = result.TotalWatchTime
	}

	return watchTimeByUser, nil
}
