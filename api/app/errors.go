package app

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Status
	Data *interface{} `json:"data"`
}

// Render sets the application-specific error code in AppCode.
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Status.Code)
	return nil
}

// ErrBadRequest returns status 400 Bad Request returns status 400 Bad Request for malformed request body including error message.
func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Status: Status{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		},
	}
}

// ErrUnprocessableEntity returns status 422 Unprocessable Entity rendering response error.
func ErrUnprocessableEntity(err error) render.Renderer {
	return &ErrResponse{
		Status: Status{
			Code:    http.StatusUnprocessableEntity,
			Message: err.Error(),
		},
	}
}

var (
	// ErrUnauthorized returns 401 Unauthorized.
	ErrUnauthorized = &ErrResponse{Status: Status{Code: http.StatusUnauthorized, Message: http.StatusText(http.StatusUnauthorized)}}

	// ErrForbidden returns status 403 Forbidden for unauthorized request.
	ErrForbidden = &ErrResponse{Status: Status{Code: http.StatusForbidden, Message: http.StatusText(http.StatusForbidden)}}

	// ErrNotFound returns status 404 Not Found for invalid resource request.
	ErrNotFound = &ErrResponse{Status: Status{Code: http.StatusNotFound, Message: http.StatusText(http.StatusNotFound)}}

	// ErrInternalServerError returns status 500 Internal Server Error.
	ErrInternalServerError = &ErrResponse{Status: Status{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}}
)
