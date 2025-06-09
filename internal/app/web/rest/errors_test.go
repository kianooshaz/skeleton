package rest

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/labstack/echo/v4"
)

func Test_errorResponse(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantBody   string
	}{
		{
			name:       "known error maps to status",
			err:        derror.ErrUserNotFound,
			wantStatus: http.StatusBadRequest,
			wantBody:   `{"error":"100100"}`,
		},
		{
			name:       "unknown error maps to 500",
			err:        errors.New("some unknown error"),
			wantStatus: http.StatusInternalServerError,
			wantBody:   `{"error":"100000"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			errorResponse(tt.err, c)

			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tt.wantStatus)
			}
			if rec.Body.String() != tt.wantBody+"\n" {
				t.Errorf("got body %q, want %q", rec.Body.String(), tt.wantBody+"\n")
			}
		})
	}
}
