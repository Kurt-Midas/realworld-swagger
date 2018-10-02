package models

import (
	"time"
)

type Comment struct {
	Id int32 `json:"id"`

	CreatedAt time.Time `json:"createdAt"`

	UpdatedAt time.Time `json:"updatedAt"`

	Body string `json:"body"`

	Author *Profile `json:"author"`
}
