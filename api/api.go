// Package api configures an http server for administration and application resources.
package api

import (
	"time"

	// _ "github.com/ykaseng/articles-library/auth/pwdless"
	"github.com/ykaseng/articles-library/api/app"
	"github.com/ykaseng/articles-library/database"
	"github.com/ykaseng/articles-library/logging"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

// New configures application resources and routes.
func New() (*chi.Mux, error) {
	logger := logging.NewLogger()

	db, err := database.DBConn()
	if err != nil {
		logger.WithField("module", "database").Error(err)
		return nil, err
	}

	// authStore := database.NewAuthStore(db)
	// authResource, err := pwdless.NewResource(authStore)
	// if err != nil {
	// 	logger.WithField("module", "auth").Error(err)
	// 	return nil, err
	// }

	appAPI, err := app.NewAPI(db)
	if err != nil {
		logger.WithField("module", "app").Error(err)
		return nil, err
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Timeout(15 * time.Second))

	r.Use(logging.NewStructuredLogger(logger))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// r.Mount("/auth", authResource.Router())
	r.Group(func(r chi.Router) {
		// r.Use(authResource.TokenAuth.Verifier())
		// r.Use(jwt.Authenticator)
		r.Mount("/api/v1", appAPI.Router())
	})

	return r, nil
}
