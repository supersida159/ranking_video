package service

import (
	"context"
	"ranking_video/internal/repository"
	helper "ranking_video/pkg/utils"
)

// EventServiceInterface defines the common methods for event services
type EventServiceInterface[T any] interface {
	Add(ctx context.Context, event T) error
	SoftDelete(ctx context.Context, id uint) error
	Count(ctx context.Context, conditions map[string]interface{}) (int64, error)
	GetList(ctx context.Context, paging *helper.Paging, orderClauses []string, conditions map[string]interface{}) ([]T, error)
}

// EventService is a generic service for event entities
type EventService[T any] struct {
	repo *repository.EventRepository[T]
}

// NewEventService creates a new event service instance
func NewEventService[T any](repo *repository.EventRepository[T]) *EventService[T] {
	return &EventService[T]{
		repo: repo,
	}
}

// Add creates a new event record
func (s *EventService[T]) Add(ctx context.Context, event T) error {
	return s.repo.Add(ctx, event)
}

// SoftDelete performs a soft delete on the event record
func (s *EventService[T]) SoftDelete(ctx context.Context, id uint) error {
	return s.repo.SoftDelete(ctx, id)
}

// Count returns the count of events matching the conditions
func (s *EventService[T]) Count(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	return s.repo.Count(ctx, conditions)
}

// GetList retrieves a list of events with sorting, paging and conditions
func (s *EventService[T]) GetList(
	ctx context.Context,
	paging *helper.Paging,
	orderClauses []string,
	conditions map[string]interface{},
) ([]T, error) {
	return s.repo.GetList(ctx, paging, orderClauses, conditions)
}
