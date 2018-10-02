package models

import (
	"time"
)

type Article struct {
	Slug string `json:"slug"`

	Title string `json:"title"`

	Description string `json:"description"`

	Body string `json:"body"`

	TagList []string `json:"tagList"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`

	Favorited bool `json:"favorited"`

	FavoritesCount int32 `json:"favoritesCount"`

	Author *Profile `json:"author"`
}
