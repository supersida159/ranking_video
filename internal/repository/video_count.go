package repository

// import (
// 	"context"
// 	"ranking_video/internal/models"
// 	apperror "ranking_video/pkg/app_error"

// 	"gorm.io/gorm"
// )

// type VideoCountRepository struct {
// 	db *gorm.DB
// }

// func NewVideoCountRepository(db *gorm.DB) *VideoCountRepository {
// 	return &VideoCountRepository{
// 		db: db,
// 	}
// }
// func (r *VideoCountRepository) UpdateCount(ctx context.Context, video *models.Video) *apperror.AppError {
// 	err := r.db.Model(&models.Video{}).Where("id = ?", video.ID).Updates(map[string]interface{}{
// 		"like_count":    gorm.Expr("GREATEST(like_count + ?, 0)", video.LikeCount),
// 		"comment_count": gorm.Expr("GREATEST(comment_count + ?,0)", video.CommentCount),
// 		"share_count":   gorm.Expr("GREATEST(share_count + ?,0)", video.ShareCount),
// 		"view_count":    gorm.Expr("GREATEST(view_count + ?,0)", video.ViewCount),
// 		"score":         gorm.Expr("GREATEST(score + ?,0)", video.Score),
// 	}).Error

// 	if err != nil {
// 		return apperror.ErrDBInternal(err)
// 	}
// 	return nil
// }
