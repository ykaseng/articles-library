package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	tt := []struct {
		name     string
		errResp  ErrResponse
		expected int
	}{
		{
			name: "render with teapot status",
			errResp: ErrResponse{
				Status: Status{
					Code:    http.StatusTeapot,
					Message: http.StatusText(http.StatusTeapot),
				},
			},
			expected: http.StatusTeapot,
		},
		{
			name: "render with unavailable for legal reasons",
			errResp: ErrResponse{
				Status: Status{
					Code:    http.StatusUnavailableForLegalReasons,
					Message: http.StatusText(http.StatusUnavailableForLegalReasons),
				},
			},
			expected: http.StatusUnavailableForLegalReasons,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "localhost:8080/api/v1/articles", nil)
			if err != nil {
				t.Errorf("create request failed: %v", err)
			}

			rec := httptest.NewRecorder()
			if err := (&tc.errResp).Render(rec, req); err != nil {
				t.Errorf("render failed: %v", err)
			}

			ctx := req.Context()
			actual := ctx.Value(render.StatusCtxKey)

			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestErrBadRequest(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		expected ErrResponse
	}{
		{
			name: "render bad request with error message",
			err:  fmt.Errorf("malformed request"),
			expected: ErrResponse{
				Status: Status{
					Code:    http.StatusBadRequest,
					Message: "malformed request",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := ErrBadRequest(tc.err)
			assert.Equal(t, tc.expected, *actual.(*ErrResponse))
		})
	}
}

func TestErrUnprocessableEntity(t *testing.T) {
	tt := []struct {
		name     string
		err      error
		expected ErrResponse
	}{
		{
			name: "render unprocessable entity request with error message",
			err:  fmt.Errorf("semantically errornous request"),
			expected: ErrResponse{
				Status: Status{
					Code:    http.StatusUnprocessableEntity,
					Message: "semantically errornous request",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actual := ErrUnprocessableEntity(tc.err)
			assert.Equal(t, tc.expected, *actual.(*ErrResponse))
		})
	}
}
