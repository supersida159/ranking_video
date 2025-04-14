package dto

// Base DTOs for requests

// BaseEventDTO contains common fields for all event types
type BaseEventDTO struct {
	UserID  string `json:"user_id" validate:"required"`
	VideoID string `json:"video_id" validate:"required"`
	IP      string `json:"ip,omitempty"`
}

// ViewEventDTO represents request data for view events
type ViewEventDTO struct {
	BaseEventDTO
}

// LikeEventDTO represents request data for like events
type LikeEventDTO struct {
	BaseEventDTO
}

// CommentEventDTO represents request data for comment events
type CommentEventDTO struct {
	BaseEventDTO
	Content string `json:"content" validate:"required"`
}

// ShareEventDTO represents request data for share events
type ShareEventDTO struct {
	BaseEventDTO
	Platform    string `json:"platform" validate:"required"`
	ReferrerURL string `json:"referrer_url,omitempty"`
}

// WatchEventDTO represents request data for watch events
type WatchEventDTO struct {
	BaseEventDTO
	Duration  int  `json:"duration" validate:"required,min=1"`
	Completed bool `json:"completed"`
}

// Update DTOs

// // UpdateViewEventDTO for updating view events
// type UpdateViewEventDTO struct {
// 	UserAgent string `json:"user_agent,omitempty"`
// 	IP        string `json:"ip,omitempty"`
// }

// // UpdateLikeEventDTO for updating like events
// type UpdateLikeEventDTO struct {
// 	Value int    `json:"value" validate:"required,oneof=-1 1"`
// 	IP    string `json:"ip,omitempty"`
// }

// // UpdateCommentEventDTO for updating comment events
// type UpdateCommentEventDTO struct {
// 	Content string `json:"content" validate:"required"`
// 	IP      string `json:"ip,omitempty"`
// }

// // UpdateShareEventDTO for updating share events
// type UpdateShareEventDTO struct {
// 	Platform    string `json:"platform,omitempty"`
// 	ReferrerURL string `json:"referrer_url,omitempty"`
// 	IP          string `json:"ip,omitempty"`
// }

// // UpdateWatchEventDTO for updating watch events
// type UpdateWatchEventDTO struct {
// 	Duration  int    `json:"duration,omitempty" validate:"omitempty,min=1"`
// 	Completed bool   `json:"completed,omitempty"`
// 	IP        string `json:"ip,omitempty"`
// }

// // Response DTOs

// // BaseEventResponseDTO contains common fields for all event type responses
// type BaseEventResponseDTO struct {
// 	ID        uint      `json:"id"`
// 	UserID    string    `json:"user_id"`
// 	VideoID   string    `json:"video_id"`
// 	IP        string    `json:"ip,omitempty"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// // ViewEventResponseDTO represents response data for view events
// type ViewEventResponseDTO struct {
// 	BaseEventResponseDTO
// 	UserAgent string `json:"user_agent,omitempty"`
// }

// // LikeEventResponseDTO represents response data for like events
// type LikeEventResponseDTO struct {
// 	BaseEventResponseDTO
// 	Value int `json:"value"`
// }

// // CommentEventResponseDTO represents response data for comment events
// type CommentEventResponseDTO struct {
// 	BaseEventResponseDTO
// 	Content   string     `json:"content"`
// 	UpdatedAt time.Time  `json:"updated_at"`
// 	DeletedAt *time.Time `json:"deleted_at,omitempty"`
// }

// // ShareEventResponseDTO represents response data for share events
// type ShareEventResponseDTO struct {
// 	BaseEventResponseDTO
// 	Platform    string `json:"platform"`
// 	ReferrerURL string `json:"referrer_url,omitempty"`
// }

// // WatchEventResponseDTO represents response data for watch events
// type WatchEventResponseDTO struct {
// 	BaseEventResponseDTO
// 	Duration  int  `json:"duration"`
// 	Completed bool `json:"completed"`
// }
