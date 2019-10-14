package app

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/ykaseng/articles-library/models"
)

// The list of error types returned from article resource.
var (
	ErrEmptyRequest = errors.New("request cannot be empty")
)

// ArticleStore defines database operations for article.
type ArticleStore interface {
	Get(id int) (*[]models.Article, error)
	GetAll() (*[]models.Article, error)
	Post(*models.Article) (*models.ArticleID, error)
}

// ArticleResource implements article management handler.
type ArticleResource struct {
	Store ArticleStore
}

// NewArticleResource creates and returns an article resource.
func NewArticleResource(store ArticleStore) *ArticleResource {
	return &ArticleResource{
		Store: store,
	}
}

func (rs *ArticleResource) router() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", rs.post)
	r.Get("/", rs.getAll)
	r.Route("/{articleID}", func(r chi.Router) {
		r.Get("/", rs.get)
	})
	return r
}

func (rs *ArticleResource) get(w http.ResponseWriter, r *http.Request) {
	type getArticleResponse struct {
		Status
		Data *[]models.Article `json:"data"`
	}

	id, err := strconv.Atoi(chi.URLParam(r, "articleID"))
	if err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	article, err := rs.Store.Get(id)
	if err != nil {
		render.Render(w, r, ErrUnprocessableEntity(err))
		return
	}

	render.Respond(w, r, &getArticleResponse{
		Status: Status{
			Code: http.StatusOK,
			Message: "SUCCESS",
		},
		Data: article,
	})
}

func (rs *ArticleResource) getAll(w http.ResponseWriter, r *http.Request) {
	type getAllArticlesResponse struct {
		Status
		Data *[]models.Article `json:"data"`
	}

	articles, err := rs.Store.GetAll()
	if err != nil {
		render.Render(w, r, ErrUnprocessableEntity(err))
		return
	}

	render.Respond(w, r, &getAllArticlesResponse{
		Status: Status{
			Code: http.StatusOK,
			Message: "SUCCESS",
		},
		Data: articles,
	})
}

func (rs *ArticleResource) post(w http.ResponseWriter, r *http.Request) {
	type postArticleRequest struct { *models.Article }
	type postArticleResponse struct {
		Status
		Data *models.ArticleID `json:"data"`
	}

	data := &postArticleRequest{}
	if err := render.DecodeJSON(r.Body, data); err != nil {
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	if *data == (postArticleRequest{}) { 
		render.Render(w, r, ErrBadRequest(ErrEmptyRequest)) 
		return 
	}
	
	if err := data.Validate(); err != nil { 
		render.Render(w, r, ErrBadRequest(err)) 
		return 
	}

	articleID, err := rs.Store.Post(data.Article)
	if err != nil {
		render.Render(w, r, ErrUnprocessableEntity(err))
		return
	}

	render.Respond(w, r, &postArticleResponse{
		Status: Status{
			Code:  http.StatusCreated,
			Message: "SUCCESS",
		},
		Data: articleID,
	})
}
