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
	ErrArticleValidation = errors.New("account validation error")
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

type Status struct {
	Status  int    `json:"status"`
	Message string `json:"mesage"`
}

type getArticleResponse struct {
	Status
	Data *[]models.Article `json:"data"`
}

func (rs *ArticleResource) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "articleID"))
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	article, err := rs.Store.Get(id)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, &getArticleResponse{
		Status: Status{
			Status: http.StatusOK,
			Message: "SUCCESS",
		},
		Data: article,
	})
}

type getAllArticlesResponse struct {
	Status
	Data *[]models.Article `json:"data"`
}

func (rs *ArticleResource) getAll(w http.ResponseWriter, r *http.Request) {
	articles, err := rs.Store.GetAll()
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, &getAllArticlesResponse{
		Status: Status{
			Status: http.StatusOK,
			Message: "SUCCESS",
		},
		Data: articles,
	})
}

type postArticleRequest struct { *models.Article }
func (a *postArticleRequest) Bind(r *http.Request) error {
	return nil
}

type postArticleResponse struct {
	Status
	Data *models.ArticleID `json:"data"`
}

func (rs *ArticleResource) post(w http.ResponseWriter, r *http.Request) {
	data := &postArticleRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	articleID, err := rs.Store.Post(data.Article)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, &postArticleResponse{
		Status: Status{
			Status:  http.StatusCreated,
			Message: "SUCCESS",
		},
		Data: articleID,
	})
}
