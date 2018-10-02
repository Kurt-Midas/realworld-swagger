package models

type MultipleArticlesResponse struct {

	Articles []Article `json:"articles"`

	ArticlesCount int32 `json:"articlesCount"`
}
