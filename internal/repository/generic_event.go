package repository

import (
	"context"
	"fmt"
	helper "ranking_video/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// EventRepository is a generic repository for event entities
type EventRepository[T any] struct {
	db *gorm.DB
}

// NewEventRepository creates a new event repository instance
func NewEventRepository[T any](db *gorm.DB) *EventRepository[T] {
	return &EventRepository[T]{
		db: db,
	}
}

// Add creates a new event record
func (r *EventRepository[T]) Add(ctx context.Context, event T) error {
	return r.db.WithContext(ctx).Create(&event).Error
}

// SoftDelete performs a soft delete on the event record
func (r *EventRepository[T]) SoftDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&id).Error
}

// Count returns the count of events matching the conditions
func (r *EventRepository[T]) Count(ctx context.Context, conditions map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(new(T))

	for field, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetList retrieves a list of events with sorting, paging and conditions
func (r *EventRepository[T]) GetList(
	ctx context.Context,
	paging *helper.Paging,
	conditions map[string]interface{},
) ([]T, error) {
	var events []T

	// Ensure paging is properly set
	paging.Fullfill()

	query := r.db.WithContext(ctx).Model(new(T))

	// Apply conditions
	for field, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// Get total count
	if err := query.Count(&paging.Total).Error; err != nil {
		return nil, err
	}

	// Apply cursor-based pagination if cursor is provided
	if !paging.CurrentCursor.IsZero() {
		query = query.Where("created_at < ?", paging.CurrentCursor)
	}

	// Apply limit and offset
	offset := (paging.Page - 1) * paging.Limit
	query = query.Offset(offset).Limit(paging.Limit)

	// Apply default sorting for cursor-based pagination
	query = query.Order("created_at DESC")

	// Execute query
	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}

	// Set next cursor if there are results
	if len(events) > 0 {
		// This assumes the events have a CreatedAt field
		// You might need to adjust this according to your actual model structure
		var lastEvent struct {
			CreatedAt time.Time
		}
		if err := query.Last(&lastEvent).Error; err == nil {
			paging.NextCursor = lastEvent.CreatedAt
		}
	}

	return events, nil
}

// GetListWithPreload retrieves a list of events with sorting, paging, and conditions using PreloadPagination
func (r *EventRepository[T]) GetListWithPreload(
	ctx context.Context,
	pagination *helper.PreloadPagination,
	conditions map[string]interface{},
) ([]T, error) {
	var events []T

	// Ensure pagination is properly set
	if pagination.Page <= 0 {
		pagination.Page = 1
	}
	if pagination.Limit <= 0 {
		pagination.Limit = 50
	}

	query := r.db.WithContext(ctx).Model(new(T))

	// Apply conditions
	for field, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", field), value)
	}

	// Get total count
	if err := query.Count(&pagination.Total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	query = query.Offset(offset).Limit(pagination.Limit)

	// Apply sorting
	if len(pagination.Sort) > 0 {
		for _, sort := range pagination.Sort {
			query = query.Order(sort)
		}
	} else {
		// Default sorting
		query = query.Order("created_at DESC")
	}

	// Execute query
	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}
