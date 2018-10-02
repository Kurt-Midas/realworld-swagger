package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/kurt-midas/realworld-swagger/routes/articles"
	"github.com/kurt-midas/realworld-swagger/routes/comments"
	"github.com/kurt-midas/realworld-swagger/routes/favorites"
	"github.com/kurt-midas/realworld-swagger/routes/profile"
	"github.com/kurt-midas/realworld-swagger/routes/tags"
	"github.com/kurt-midas/realworld-swagger/routes/user"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type MyRouter struct {
	U *user.UserApiEngine
}

func (r MyRouter) NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range r.getroutes() {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func (r MyRouter) getroutes() []Route {
	return []Route{
		{
			"Index",
			http.MethodGet,
			"/api/",
			Index,
		},

		{
			"CreateArticle",
			http.MethodPost,
			"/api/articles",
			articles.CreateArticle,
		},

		{
			"DeleteArticle",
			http.MethodDelete,
			"/api/articles/{slug}",
			articles.DeleteArticle,
		},

		{
			"GetArticle",
			http.MethodGet,
			"/api/articles/{slug}",
			articles.GetArticle,
		},

		{
			"GetArticles",
			http.MethodGet,
			"/api/articles",
			articles.GetArticles,
		},

		{
			"GetArticlesFeed",
			http.MethodGet,
			"/api/articles/feed",
			articles.GetArticlesFeed,
		},

		{
			"UpdateArticle",
			http.MethodPut,
			"/api/articles/{slug}",
			articles.UpdateArticle,
		},

		{
			"CreateArticleComment",
			http.MethodPost,
			"/api/articles/{slug}/comments",
			comments.CreateArticleComment,
		},

		{
			"DeleteArticleComment",
			http.MethodDelete,
			"/api/articles/{slug}/comments/{id}",
			comments.DeleteArticleComment,
		},

		{
			"GetArticleComments",
			http.MethodGet,
			"/api/articles/{slug}/comments",
			comments.GetArticleComments,
		},

		{
			"TagsGet",
			http.MethodGet,
			"/api/tags",
			tags.TagsGet,
		},

		{
			"CreateArticleFavorite",
			http.MethodPost,
			"/api/articles/{slug}/favorite",
			favorites.CreateArticleFavorite,
		},

		{
			"DeleteArticleFavorite",
			http.MethodDelete,
			"/api/articles/{slug}/favorite",
			favorites.DeleteArticleFavorite,
		},

		{
			"FollowUserByUsername",
			http.MethodPost,
			"/api/profiles/{username}/follow",
			profile.FollowUserByUsername,
		},

		{
			"GetProfileByUsername",
			http.MethodGet,
			"/api/profiles/{username}",
			profile.GetProfileByUsername,
		},

		{
			"UnfollowUserByUsername",
			http.MethodDelete,
			"/api/profiles/{username}/follow",
			profile.UnfollowUserByUsername,
		},

		{
			"CreateUser",
			http.MethodPost,
			"/api/users",
			r.U.CreateUser,
		},

		{
			"GetCurrentUser",
			http.MethodGet,
			"/api/users",
			r.U.GetCurrentUser,
		},

		{
			"Login",
			http.MethodPost,
			"/api/users/login",
			r.U.Login,
		},

		{
			"UpdateCurrentUser",
			http.MethodPut,
			"/api/users",
			r.U.UpdateCurrentUser,
		},
	}
}
