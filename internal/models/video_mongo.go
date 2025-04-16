package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entity represents a user, creator, or channel in the system
type Entity struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	StringID    string             `bson:"string_id,omitempty"` // If you want to keep your string IDs
	Type        string             `bson:"type"`
	Name        string             `bson:"name"`
	Description string             `bson:"description,omitempty"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	DeletedAt   *time.Time         `bson:"deleted_at,omitempty"`

	// In MongoDB, you don't typically include relationships directly in the struct
	// They are retrieved through separate queries
}

// Video represents a video in the system
type Video struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	StringID     string             `bson:"string_id,omitempty"` // If you want to keep your string IDs
	EntityID     string             `bson:"entity_id"`
	Title        string             `bson:"title"`
	Description  string             `bson:"description,omitempty"`
	URL          string             `bson:"url"`
	Score        uint               `bson:"score,omitempty"`
	ViewCount    int                `bson:"view_count,omitempty"`
	LikeCount    int                `bson:"like_count,omitempty"`
	CommentCount int                `bson:"comment_count,omitempty"`
	ShareCount   int                `bson:"share_count,omitempty"`
	Status       string             `bson:"status,omitempty"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty"`
}

type VideoCountDaily struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	VideoID   string             `bson:"video_id"`
	Score     uint               `bson:"score"`
	CreatedAt time.Time          `bson:"created_at"`
}

// ViewEvent represents a user view on a video
type ViewEvent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	VideoID   string             `bson:"video_id"`
	IP        string             `bson:"ip,omitempty"`
	UserAgent string             `bson:"user_agent,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
}

func (e ViewEvent) GetVideoID() string {
	return e.VideoID
}

// LikeEvent represents a user like on a video
type LikeEvent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	VideoID   string             `bson:"video_id"`
	Value     int                `bson:"value"`
	IP        string             `bson:"ip,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
}

func (e LikeEvent) GetVideoID() string {
	return e.VideoID
}

// CommentEvent represents a user comment on a video
type CommentEvent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    string             `bson:"user_id"`
	VideoID   string             `bson:"video_id"`
	Content   string             `bson:"content"`
	IP        string             `bson:"ip,omitempty"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt *time.Time         `bson:"deleted_at,omitempty"`
}

func (e CommentEvent) GetVideoID() string {
	return e.VideoID
}

// ShareEvent represents a video share by a user
type ShareEvent struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"user_id"`
	VideoID     string             `bson:"video_id"`
	Platform    string             `bson:"platform"`
	ReferrerURL string             `bson:"referrer_url,omitempty"`
	IP          string             `bson:"ip,omitempty"`
	CreatedAt   time.Time          `bson:"created_at"`
}

func (e ShareEvent) GetVideoID() string {
	return e.VideoID
}

// UserPreference for personalized ranking
type UserPreference struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     string             `bson:"user_id"`
	Categories []string           `bson:"categories,omitempty"` // In MongoDB, store as actual array
	Entities   []string           `bson:"entities,omitempty"`   // In MongoDB, store as actual array
	UpdatedAt  time.Time          `bson:"updated_at"`
}
