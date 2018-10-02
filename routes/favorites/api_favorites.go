package favorites

import (
	"net/http"
)

// CreateArticleFavorite - Favorite an article
func CreateArticleFavorite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// DeleteArticleFavorite - Unfavorite an article
func DeleteArticleFavorite(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}
