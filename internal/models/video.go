package models

import (
	"time"

	"gorm.io/gorm"
)

// Entity represents a user, creator, or channel in the system
type Entity struct {
	ID          string `gorm:"primaryKey;type:varchar(36)"`
	Type        string `gorm:"type:enum('user','creator','channel');not null"`
	Name        string `gorm:"type:varchar(255);not null"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relationships
	Videos []Video `gorm:"foreignKey:EntityID"`
}

// Video represents a video in the system
type Video struct {
	ID           string  `gorm:"primaryKey;type:varchar(36)"`
	EntityID     string  `gorm:"type:varchar(36);index;not null"`
	Title        string  `gorm:"type:varchar(255);not null"`
	Description  string  `gorm:"type:text"`
	URL          string  `gorm:"type:varchar(512);not null"`
	Duration     int     `gorm:"not null"`           // Duration in seconds
	Score        float64 `gorm:"not null;default:0"` // Current calculated score
	ViewCount    int     `gorm:"not null;default:0"`
	LikeCount    int     `gorm:"not null;default:0"`
	CommentCount int     `gorm:"not null;default:0"`
	ShareCount   int     `gorm:"not null;default:0"`
	WatchTime    int     `gorm:"not null;default:0"` // Total watch time in seconds
	Status       string  `gorm:"type:enum('active','inactive','deleted');default:'active'"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Foreign Key Relationship
	Entity        Entity         `gorm:"foreignKey:EntityID"`
	ViewEvents    []ViewEvent    `gorm:"foreignKey:VideoID"`
	LikeEvents    []LikeEvent    `gorm:"foreignKey:VideoID"`
	CommentEvents []CommentEvent `gorm:"foreignKey:VideoID"`
	ShareEvents   []ShareEvent   `gorm:"foreignKey:VideoID"`
	WatchEvents   []WatchEvent   `gorm:"foreignKey:VideoID"`
}

// ViewEvent represents a user view on a video
type ViewEvent struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UserID    string `gorm:"type:varchar(36);index;not null"`
	VideoID   string `gorm:"type:varchar(36);index;not null"`
	IP        string `gorm:"type:varchar(45)"`
	UserAgent string `gorm:"type:varchar(512)"`
	CreatedAt time.Time

	// Foreign Key Relationship
	User  Entity `gorm:"foreignKey:UserID"`
	Video Video  `gorm:"foreignKey:VideoID"`
}

// LikeEvent represents a user like on a video
type LikeEvent struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UserID    string `gorm:"type:varchar(36);index;not null"`
	VideoID   string `gorm:"type:varchar(36);index;not null"`
	Value     int    `gorm:"not null;default:1"` // Usually 1, could be -1 for dislike
	IP        string `gorm:"type:varchar(45)"`
	CreatedAt time.Time

	// Foreign Key Relationship
	User  Entity `gorm:"foreignKey:UserID"`
	Video Video  `gorm:"foreignKey:VideoID"`
}

// CommentEvent represents a user comment on a video
type CommentEvent struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UserID    string `gorm:"type:varchar(36);index;not null"`
	VideoID   string `gorm:"type:varchar(36);index;not null"`
	Content   string `gorm:"type:text;not null"`
	IP        string `gorm:"type:varchar(45)"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Foreign Key Relationship
	User  Entity `gorm:"foreignKey:UserID"`
	Video Video  `gorm:"foreignKey:VideoID"`
}

// ShareEvent represents a video share by a user
type ShareEvent struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	UserID      string `gorm:"type:varchar(36);index;not null"`
	VideoID     string `gorm:"type:varchar(36);index;not null"`
	Platform    string `gorm:"type:varchar(50);not null"` // e.g., Facebook, Twitter
	ReferrerURL string `gorm:"type:varchar(512)"`
	IP          string `gorm:"type:varchar(45)"`
	CreatedAt   time.Time

	// Foreign Key Relationship
	User  Entity `gorm:"foreignKey:UserID"`
	Video Video  `gorm:"foreignKey:VideoID"`
}

// WatchEvent represents a user's watch time on a video
type WatchEvent struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	UserID    string `gorm:"type:varchar(36);index;not null"`
	VideoID   string `gorm:"type:varchar(36);index;not null"`
	Duration  int    `gorm:"not null"` // Watch time in seconds
	Completed bool   `gorm:"not null;default:false"`
	IP        string `gorm:"type:varchar(45)"`
	CreatedAt time.Time

	// Foreign Key Relationship
	User  Entity `gorm:"foreignKey:UserID"`
	Video Video  `gorm:"foreignKey:VideoID"`
}

// UserPreference for personalized ranking
type UserPreference struct {
	UserID     string `gorm:"primaryKey;type:varchar(36)"`
	Categories string `gorm:"type:json"` // JSON stored as string in MySQL
	Entities   string `gorm:"type:json"` // JSON stored as string in MySQL
	UpdatedAt  time.Time

	// Foreign Key Relationship
	User Entity `gorm:"foreignKey:UserID"`
}
