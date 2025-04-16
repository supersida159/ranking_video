package repository

import (
	"context"
	"errors"
	apperror "ranking_video/pkg/app_error"
	helper "ranking_video/pkg/utils"
	"reflect"
	"strings"
	"time"
	"unicode"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository[T any] struct {
	collection *mongo.Collection
}

func NewEventRepository[T any](db *mongo.Database) *EventRepository[T] {
	// Get the type name of T and convert it to a collection name
	var model T
	modelType := reflect.TypeOf(model)

	// Handle pointer types
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// Get the type name and convert it to snake_case for collection name
	typeName := modelType.Name()
	collectionName := toSnakeCase(typeName)

	return &EventRepository[T]{
		collection: db.Collection(collectionName),
	}
}

func (r *EventRepository[T]) Add(ctx context.Context, event T) *apperror.AppError {
	_, err := r.collection.InsertOne(ctx, event)
	if err != nil {
		return apperror.ErrDBInternal(err)
	}
	return nil
}

func (r *EventRepository[T]) SoftDelete(ctx context.Context, id primitive.ObjectID) *apperror.AppError {
	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id, "deleted_at": nil},
		update,
	)

	if err != nil {
		return apperror.ErrDBInternal(err)
	}

	if result.MatchedCount == 0 {
		return apperror.ErrRecordNotFound(errors.New("record not found"))
	}

	return nil
}

func (r *EventRepository[T]) Count(ctx context.Context, conditions map[string]interface{}) (int64, *apperror.AppError) {
	filter := bson.M{"deleted_at": nil}

	for field, value := range conditions {
		filter[field] = value
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, apperror.ErrDBInternal(err)
	}

	return count, nil
}

func (r *EventRepository[T]) GetList(
	ctx context.Context,
	paging *helper.Paging,
	orderClauses []string,
	conditions map[string]interface{},
) ([]T, *apperror.AppError) {
	var events []T

	paging.Fullfill()

	// Build filter
	filter := bson.M{"deleted_at": nil}
	for field, value := range conditions {
		filter[field] = value
	}

	// Count total documents for pagination
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, apperror.ErrDBInternal(err)
	}
	paging.Total = count

	// Build options (sorting and pagination)
	findOptions := options.Find()
	findOptions = r.applySorting(findOptions, orderClauses)
	findOptions = r.applyPagination(findOptions, paging)

	// Execute query
	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, apperror.ErrDBInternal(err)
	}
	defer cursor.Close(ctx)

	// Decode results
	if err := cursor.All(ctx, &events); err != nil {
		return nil, apperror.ErrDBInternal(err)
	}

	return events, nil
}

func (r *EventRepository[T]) applyPagination(findOptions *options.FindOptions, paging *helper.Paging) *options.FindOptions {
	skip := int64((paging.Page - 1) * paging.Limit)

	// If using cursor-based pagination
	if !paging.CurrentCursor.IsZero() {
		// For cursor-based pagination, we need to modify the filter instead
		// This is handled outside this function in the main query
		// No need to set Skip since we're filtering by created_at
	} else {
		// For page-based pagination
		findOptions.SetSkip(skip)
	}

	findOptions.SetLimit(int64(paging.Limit))
	return findOptions
}

func (r *EventRepository[T]) applySorting(findOptions *options.FindOptions, orderClauses []string) *options.FindOptions {
	if len(orderClauses) > 0 {
		sort := bson.D{}

		for _, orderClause := range orderClauses {
			// Parse order clause (e.g., "created_at desc")
			field := orderClause
			order := 1 // default ascending

			// Check if it has desc or asc suffix
			if len(field) > 5 && field[len(field)-5:] == " desc" {
				field = field[:len(field)-5]
				order = -1
			} else if len(field) > 4 && field[len(field)-4:] == " asc" {
				field = field[:len(field)-4]
			}

			sort = append(sort, bson.E{Key: field, Value: order})
		}

		findOptions.SetSort(sort)
	} else {
		// Default sort by created_at descending
		findOptions.SetSort(bson.D{{Key: "created_at", Value: -1}})
	}

	return findOptions
} // toSnakeCase converts a CamelCase string to snake_case
func toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
