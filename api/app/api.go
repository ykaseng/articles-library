// Package app ties together application resources and handlers.
package app

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-pg/pg"
	"github.com/sirupsen/logrus"

	"github.com/ykaseng/articles-library/database"
	"github.com/ykaseng/articles-library/logging"
)

type ctxKey int

const (
	ctxAccount ctxKey = iota
	ctxProfile
)

// API provides application resources and handlers.
type API struct {
	Article *ArticleResource
}

// NewAPI configures and returns application API.
func NewAPI(db *pg.DB) (*API, error) {
	articleStore := database.NewArticleStore(db)
	article := NewArticleResource(articleStore)

	api := &API{
		Article: article,
	}

	return api, nil
}

// Router provides application routes.
func (a *API) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Mount("/articles", a.Article.router())

	return r
}

func log(r *http.Request) logrus.FieldLogger {
	return logging.GetLogEntry(r)
}
