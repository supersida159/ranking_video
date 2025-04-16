package repository

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	apperror "ranking_video/pkg/app_error"
// 	helper "ranking_video/pkg/utils"

// 	"gorm.io/gorm"
// )

// type EventRepository[T any] struct {
// 	db *gorm.DB
// }

// func NewEventRepository[T any](db *gorm.DB) *EventRepository[T] {
// 	return &EventRepository[T]{
// 		db: db,
// 	}
// }

// func (r *EventRepository[T]) Add(ctx context.Context, event T) *apperror.AppError {
// 	if err := r.db.WithContext(ctx).Create(&event).Error; err != nil {
// 		return apperror.ErrDBInternal(err)
// 	}
// 	return nil
// }

// func (r *EventRepository[T]) SoftDelete(ctx context.Context, id uint) *apperror.AppError {
// 	err := r.db.WithContext(ctx).Delete(new(T), id).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return apperror.ErrRecordNotFound(err)
// 		}
// 		return apperror.ErrDBInternal(err)
// 	}
// 	return nil
// }

// func (r *EventRepository[T]) Count(ctx context.Context, conditions map[string]interface{}) (int64, *apperror.AppError) {
// 	var count int64
// 	query := r.db.WithContext(ctx).Model(new(T))

// 	for field, value := range conditions {
// 		query = query.Where(fmt.Sprintf("%s = ?", field), value)
// 	}

// 	if err := query.Count(&count).Error; err != nil {
// 		return 0, apperror.ErrDBInternal(err)
// 	}

// 	return count, nil
// }

// func (r *EventRepository[T]) GetList(
// 	ctx context.Context,
// 	paging *helper.Paging,
// 	orderClauses []string,
// 	conditions map[string]interface{},
// ) ([]T, *apperror.AppError) {
// 	var events []T

// 	paging.Fullfill()

// 	query := r.db.WithContext(ctx).Model(new(T))

// 	for field, value := range conditions {
// 		query = query.Where(fmt.Sprintf("%s = ?", field), value)
// 	}

// 	if err := query.Count(&paging.Total).Error; err != nil {
// 		return nil, apperror.ErrDBInternal(err)
// 	}

// 	query = r.applyPagination(query, paging)
// 	query = r.applySorting(query, orderClauses)

// 	if err := query.Find(&events).Error; err != nil {
// 		return nil, apperror.ErrDBInternal(err)
// 	}

// 	return events, nil
// }

// func (r *EventRepository[T]) applyPagination(query *gorm.DB, paging *helper.Paging) *gorm.DB {
// 	if !paging.CurrentCursor.IsZero() {
// 		query = query.Where("created_at < ?", paging.CurrentCursor)
// 	} else {
// 		query = query.Offset((paging.Page - 1) * paging.Limit)
// 	}
// 	return query.Limit(paging.Limit)
// }

// func (r *EventRepository[T]) applySorting(query *gorm.DB, orderClauses []string) *gorm.DB {
// 	if len(orderClauses) > 0 {
// 		for _, orderClause := range orderClauses {
// 			query = query.Order(orderClause)
// 		}
// 	} else {
// 		query = query.Order("created_at desc")
// 	}
// 	return query
// }
