package models

import "time"

type Material struct {
	Id          uint64    `json:"id"`
	CategoryId  *uint32   `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	PreviewMeta string    `json:"preview_meta"`
	VideoMeta   string    `json:"video_meta"`
	Deleted     bool      `json:"deleted"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MaterialSearch struct {
	Id          string  `json:"-"`
	CategoryId  *uint32 `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
}

type Category struct {
	Id         uint32  `json:"id"`
	CategoryId *uint32 `json:"category_id"`
	Name       string  `json:"name"`
}
