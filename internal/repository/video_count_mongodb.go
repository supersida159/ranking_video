package repository

import (
	"context"
	"ranking_video/internal/models"
	apperror "ranking_video/pkg/app_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VideoCountRepository struct {
	collection *mongo.Collection
}

func NewVideoCountRepository(db *mongo.Database) *VideoCountRepository {
	return &VideoCountRepository{
		collection: db.Collection("videos"),
	}
}

func (r *VideoCountRepository) UpdateCount(ctx context.Context, video *models.Video) *apperror.AppError {
	var videoID primitive.ObjectID
	var err error

	// Convert string ID to ObjectID if needed
	videoID = video.ID

	// In MongoDB, we use $inc to increment values
	// And $max to ensure values don't go below 0
	update := bson.M{
		"$set": bson.M{
			"like_count": bson.M{
				"$max": []interface{}{
					bson.M{"$add": []interface{}{"$like_count", video.LikeCount}},
					0,
				},
			},
			"comment_count": bson.M{
				"$max": []interface{}{
					bson.M{"$add": []interface{}{"$comment_count", video.CommentCount}},
					0,
				},
			},
			"share_count": bson.M{
				"$max": []interface{}{
					bson.M{"$add": []interface{}{"$share_count", video.ShareCount}},
					0,
				},
			},
			"view_count": bson.M{
				"$max": []interface{}{
					bson.M{"$add": []interface{}{"$view_count", video.ViewCount}},
					0,
				},
			},
			"score": bson.M{
				"$max": []interface{}{
					bson.M{"$add": []interface{}{"$score", video.Score}},
					0,
				},
			},
		},
	}

	_, err = r.collection.UpdateOne(
		ctx,
		bson.M{"_id": videoID},
		update,
	)

	if err != nil {
		return apperror.ErrDBInternal(err)
	}
	return nil
}
