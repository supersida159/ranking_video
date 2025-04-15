package services_processor

import (
	"context"
	"ranking_video/internal/models"
	apperror "ranking_video/pkg/app_error"
)

type VideoCountServiceInterface interface {
	UpdateCount(ctx context.Context, video *models.Video) *apperror.AppError
}

type VideoCountService struct {
	storage VideoCountServiceInterface
}

func NewVideoCountService(storage VideoCountServiceInterface) *VideoCountService {
	return &VideoCountService{
		storage: storage,
	}
}
func (s *VideoCountService) UpdateCount(ctx context.Context, video *models.Video) *apperror.AppError {
	return s.storage.UpdateCount(ctx, video)
}
