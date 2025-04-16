package service_interact

// import (
// 	"context"
// 	"ranking_video/internal/kafka"
// 	apperror "ranking_video/pkg/app_error"
// 	helper "ranking_video/pkg/utils"
// 	"reflect"
// )

// // EventEntity defines the required fields for event entities
// type EventEntity interface {
// 	GetVideoID() string // Method to retrieve VideoID
// }

// // EventServiceInterface defines the common methods for event services
// type EventServiceInterface[T EventEntity] interface {
// 	Add(ctx context.Context, event T) *apperror.AppError
// 	SoftDelete(ctx context.Context, id uint) *apperror.AppError
// 	Count(ctx context.Context, conditions map[string]interface{}) (int64, *apperror.AppError)
// 	GetList(ctx context.Context, paging *helper.Paging, orderClauses []string, conditions map[string]interface{}) ([]T, *apperror.AppError)
// }

// // EventService is a generic service for event entities
// type EventService[T EventEntity] struct {
// 	storage  EventServiceInterface[T]
// 	producer *kafka.Producer
// }

// // NewEventService creates a new event service instance
// func NewEventService[T EventEntity](storage EventServiceInterface[T], producer *kafka.Producer) *EventService[T] {
// 	return &EventService[T]{
// 		storage:  storage,
// 		producer: producer,
// 	}
// }

// // Add creates a new event record and sends a message to Kafka
// func (s *EventService[T]) Add(ctx context.Context, event T) *apperror.AppError {
// 	// Get the type name for the event
// 	eventType := reflect.TypeOf(event).Name()

// 	// Create a Kafka event
// 	kafkaEvent, err := kafka.NewEvent(eventType, event.GetVideoID())
// 	if err != nil {
// 		return apperror.ErrFailedCreateEvent(err)
// 	}

// 	// Send the event to Kafka
// 	if err := s.producer.SendEvent(kafkaEvent); err != nil {
// 		return apperror.ErrFailedSendEvent(err)
// 	}

// 	// Continue with the original storage operation
// 	return s.storage.Add(ctx, event)
// }

// // SoftDelete performs a soft delete on the event record
// func (s *EventService[T]) SoftDelete(ctx context.Context, id uint) *apperror.AppError {
// 	return s.storage.SoftDelete(ctx, id)
// }

// // Count returns the count of events matching the conditions
// func (s *EventService[T]) Count(ctx context.Context, conditions map[string]interface{}) (int64, *apperror.AppError) {
// 	return s.storage.Count(ctx, conditions)
// }

// // GetList retrieves a list of events with sorting, paging and conditions
// func (s *EventService[T]) GetList(
// 	ctx context.Context,
// 	paging *helper.Paging,
// 	orderClauses []string,
// 	conditions map[string]interface{},
// ) ([]T, *apperror.AppError) {
// 	return s.storage.GetList(ctx, paging, orderClauses, conditions)
// }
