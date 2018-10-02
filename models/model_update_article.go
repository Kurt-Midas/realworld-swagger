package models

type UpdateArticle struct {

	Title string `json:"title,omitempty"`

	Description string `json:"description,omitempty"`

	Body string `json:"body,omitempty"`
}
