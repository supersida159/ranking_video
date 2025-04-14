package repository

import (
	"context"
	"fmt"
	helper "ranking_video/pkg/utils"

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
	orderClauses []string,
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

	// Apply limit and offset
	query = r.applyPagination(query, paging)

	// Apply default sorting for cursor-based pagination
	query = r.applySorting(query, orderClauses)

	// Execute query
	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}

	// // Set next cursor if there are results
	// if len(events) > 0 {
	// 	// This assumes the events have a CreatedAt field
	// 	// You might need to adjust this according to your actual model structure
	// 	var lastEvent struct {
	// 		CreatedAt time.Time
	// 	}
	// 	if err := query.Last(&lastEvent).Error; err == nil {
	// 		paging.NextCursor = lastEvent.CreatedAt
	// 	}
	// }

	return events, nil
}

// applyPagination applies pagination parameters to a query
func (s *EventRepository[T]) applyPagination(query *gorm.DB, paging *helper.Paging) *gorm.DB {
	if !paging.CurrentCursor.IsZero() {
		query = query.Where("created_at < ?", paging.CurrentCursor)
	} else {
		query = query.Offset((paging.Page - 1) * paging.Limit)
	}
	return query.Limit(paging.Limit)
}

// applySorting applies sorting parameters to a query
func (s *EventRepository[T]) applySorting(query *gorm.DB, orderClauses []string) *gorm.DB {
	if len(orderClauses) > 0 {
		for _, orderClause := range orderClauses {
			query = query.Order(orderClause)
		}
	} else {
		query = query.Order("%s.created_at desc")
	}
	return query
}
