package service

import (
	"context"
	apperror "ranking_video/pkg/app_error"
	helper "ranking_video/pkg/utils"
)

// EventServiceInterface defines the common methods for event services
type EventServiceInterface[T any] interface {
	Add(ctx context.Context, event T) *apperror.AppError
	SoftDelete(ctx context.Context, id uint) *apperror.AppError
	Count(ctx context.Context, conditions map[string]interface{}) (int64, *apperror.AppError)
	GetList(ctx context.Context, paging *helper.Paging, orderClauses []string, conditions map[string]interface{}) ([]T, *apperror.AppError)
}

// EventService is a generic service for event entities
type EventService[T any] struct {
	storage EventServiceInterface[T]
}

// NewEventService creates a new event service instance
func NewEventService[T any](storage EventServiceInterface[T]) *EventService[T] {
	return &EventService[T]{
		storage: storage,
	}
}

// Add creates a new event record
func (s *EventService[T]) Add(ctx context.Context, event T) *apperror.AppError {
	return s.storage.Add(ctx, event)
}

// SoftDelete performs a soft delete on the event record
func (s *EventService[T]) SoftDelete(ctx context.Context, id uint) *apperror.AppError {
	return s.storage.SoftDelete(ctx, id)
}

// Count returns the count of events matching the conditions
func (s *EventService[T]) Count(ctx context.Context, conditions map[string]interface{}) (int64, *apperror.AppError) {
	return s.storage.Count(ctx, conditions)
}

// GetList retrieves a list of events with sorting, paging and conditions
func (s *EventService[T]) GetList(
	ctx context.Context,
	paging *helper.Paging,
	orderClauses []string,
	conditions map[string]interface{},
) ([]T, *apperror.AppError) {
	return s.storage.GetList(ctx, paging, orderClauses, conditions)
}
