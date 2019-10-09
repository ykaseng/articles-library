package app

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/ykaseng/articles-library/models"
)

// The list of error types returned from account resource.
var (
	ErrArticleValidation = errors.New("account validation error")
)

// ArticleStore defines database operations for account.
type ArticleStore interface {
	Get(id int) (*models.Article, error)
	GetAll() (*[]models.Article, error)
	Post(*models.Article) (int, error)
}

// ArticleResource implements account management handler.
type ArticleResource struct {
	Store ArticleStore
}

// NewArticleResource creates and returns an account resource.
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

}

func (rs *ArticleResource) getAll(w http.ResponseWriter, r *http.Request) {

}

type articleRequest struct {
	*models.Article
}

func (a *articleRequest) Bind(r *http.Request) error {
	return nil
}

type articleResponse struct {
	Status  int    `json:"status"`
	Message string `json:"mesage"`
	Data    struct {
		ID int `json:"id"`
	} `json:"data,omitempty"`
}

func (rs *ArticleResource) post(w http.ResponseWriter, r *http.Request) {
	data := &articleRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	articleID, err := rs.Store.Post(data.Article)
	if err != nil {
		render.Render(w, r, ErrRender(err))
		return
	}

	render.Respond(w, r, &articleResponse{
		Status:  http.StatusCreated,
		Message: "SUCCESS",
		Data: struct {
			ID int `json:"id"`
		}{
			ID: articleID,
		},
	})
}
