package repository

// import (
// 	"context"
// 	"ranking_video/internal/models"
// 	apperror "ranking_video/pkg/app_error"

// 	"gorm.io/gorm"
// )

// type VideoScoreDailyRepository struct {
// 	db *gorm.DB
// }

// func NewVideoScoreDailyRepository(db *gorm.DB) *VideoScoreDailyRepository {
// 	return &VideoScoreDailyRepository{
// 		db: db,
// 	}
// }
// func (r *VideoScoreDailyRepository) Create(ctx context.Context, video *models.VideoCountDaily) *apperror.AppError {
// 	err := r.db.WithContext(ctx).Create(video).Error
// 	if err != nil {
// 		return apperror.ErrDBInternal(err)
// 	}
// 	return nil
// }
