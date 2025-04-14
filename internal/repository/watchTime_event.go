package repository

import (
	"context"
	"fmt"
	apperror "ranking_video/pkg/app_error"
	"time"

	"gorm.io/gorm"
)
type WatchEventRepository[T any] struct {
	*EventRepository[T]
}

func NewWatchEventRepository[T any](db *gorm.DB) *WatchEventRepository[T] {
	return &WatchEventRepository[T]{
		EventRepository: NewEventRepository[T](db),
	}
}

func (r *WatchEventRepository[T]) CalculateWatchTime(
	ctx context.Context,
	videoID string,
	conditions map[string]interface{},
) (int, *apperror.AppError)) {
	var totalWatchTime int

	query := r.db.WithContext(ctx).
		Model(new(T)).
		Select("COALESCE(SUM(duration), 0) as total_watch_time").
		Where("video_id = ?", videoID)

	for field, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	if err := query.Scan(&totalWatchTime).Error; err != nil {
		return 0, apperror.ErrDBInternal(err)
	}

	return totalWatchTime, nil
}

func (r *WatchEventRepository[T]) CalculateWatchTimeByTimeRange(
	ctx context.Context,
	videoID string,
	startTime time.Time,
	endTime time.Time,
	additionalConditions map[string]interface{},
) (int, *apperror.AppError) {
	var totalWatchTime int

	query := r.db.WithContext(ctx).
		Model(new(T)).
		Select("COALESCE(SUM(duration), 0) as total_watch_time").
		Where("video_id = ?", videoID).
		Where("created_at BETWEEN ? AND ?", startTime, endTime)

	for field, value := range additionalConditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	if err := query.Scan(&totalWatchTime).Error; err != nil {
		return 0, apperror.ErrDBInternal(err)
	}

	return totalWatchTime, nil
}

func (r *WatchEventRepository[T]) GetWatchTimeByUser(
	ctx context.Context,
	videoID string,
) (map[string]int, *apperror.AppError) {
	type Result struct {
		UserID         string
		TotalWatchTime int
	}

	var results []Result

	query := r.db.WithContext(ctx).
		Model(new(T)).
		Select("user_id, COALESCE(SUM(duration), 0) as total_watch_time").
		Where("video_id = ?", videoID).
		Group("user_id")

	if err := query.Scan(&results).Error; err != nil {
		return nil, apperror.ErrDBInternal(err)
	}

	watchTimeByUser := make(map[string]int)
	for _, result := range results {
		watchTimeByUser[result.UserID] = result.TotalWatchTime
	}

	return watchTimeByUser, nil
}