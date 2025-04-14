package helper

import "time"

type Paging struct {
	Page  int   `json:"page" form:"page"`
	Limit int   `json:"limit" form:"limit"`
	Total int64 `json:"total" form:"total"`
	// support cusor with UID
	CurrentCursor time.Time `json:"cursor" form:"cursor"`
	NextCursor    time.Time `json:"next_cursor" form:"next_cursor"`
}

func (p *Paging) Fullfill() {
	if p.Page <= 0 {
		p.Page = 1
	}

	if p.Limit <= 0 {
		p.Limit = 50
	}
}

// PreloadPagination defines the structure for pagination and sorting options
type PreloadPagination struct {
	Limit int      `json:"limit"`
	Page  int      `json:"page"`
	Sort  []string `json:"sort"` // Support multiple sort fields
	Total int64    `json:"total"`
}
