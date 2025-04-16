package repository

import (
	"context"
	"ranking_video/internal/models"
	apperror "ranking_video/pkg/app_error"

	"go.mongodb.org/mongo-driver/mongo"
)

type VideoScoreDailyRepository struct {
	collection *mongo.Collection
}

func NewVideoScoreDailyRepository(db *mongo.Database) *VideoScoreDailyRepository {
	return &VideoScoreDailyRepository{
		collection: db.Collection("video_count_daily"),
	}
}

func (r *VideoScoreDailyRepository) Create(ctx context.Context, video *models.VideoCountDaily) *apperror.AppError {
	_, err := r.collection.InsertOne(ctx, video)
	if err != nil {
		return apperror.ErrDBInternal(err)
	}
	return nil
}
