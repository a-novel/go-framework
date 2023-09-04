package errors_test

import (
	"fmt"
	"github.com/a-novel/go-framework/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestErrorToHTTPCode(t *testing.T) {
	data := []struct {
		name string

		err          error
		errorsToCode []errors.HTTPError

		expectStatus int
	}{
		{
			name: "Success",
			err:  errors.ErrInvalidEntity,
			errorsToCode: []errors.HTTPError{
				{
					Err:  errors.ErrNotFound,
					Code: http.StatusNotFound,
				},
				{
					Err:  errors.ErrInvalidEntity,
					Code: http.StatusBadRequest,
				},
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Success/NotFound",
			err:  errors.ErrInvalidEntity,
			errorsToCode: []errors.HTTPError{
				{
					Err:  errors.ErrNotFound,
					Code: http.StatusNotFound,
				},
			},
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			errors.ErrorToHTTPCode(c, d.err, d.errorsToCode)
			require.Equal(t, d.expectStatus, w.Code)
		})
	}
}

func TestNewHTTPClientErr(t *testing.T) {
	data := []struct {
		name string

		res            *http.Response
		expectedStatus []int

		expect *errors.HTTPClientErr
	}{
		{
			name: "Success",
			res: &http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
			},
			expect: &errors.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("got unexpected status code 400"),
			},
		},
		{
			name: "Success/WithBody",
			res: &http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("response body content")),
			},
			expect: &errors.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("got unexpected status code 400: %w", fmt.Errorf("response body content")),
			},
		},
		{
			name: "Success/WithExpectedStatus",
			res: &http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
			},
			expectedStatus: []int{http.StatusOK},
			expect: &errors.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("expected status 200, got 400"),
			},
		},
		{
			name: "Success/WithExpectedStatuses",
			res: &http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
			},
			expectedStatus: []int{http.StatusOK, http.StatusNotFound},
			expect: &errors.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("expected any status in [200 404], got 400"),
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := errors.NewHTTPClientErr(d.res, d.expectedStatus...)
			require.Equal(t, d.expect, err)
		})
	}
}

func TestAsHTTPClientErr(t *testing.T) {
	data := []struct {
		name string

		err error

		expect *errors.HTTPClientErr
	}{
		{
			name: "Success",
			err: errors.NewHTTPClientErr(&http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
			}, http.StatusOK),
			expect: &errors.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("expected status 200, got 400"),
			},
		},
		{
			name: "Success/WithBody",
			err: errors.NewHTTPClientErr(&http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("response body content")),
			}),
			expect: &errors.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("got unexpected status code 400: %w", fmt.Errorf("response body content")),
			},
		},
		{
			name: "Success/Nested",
			err: fmt.Errorf("error wrapper message: %w", errors.NewHTTPClientErr(&http.Response{
				Status:     "400 Bad Request",
				StatusCode: http.StatusBadRequest,
			}, http.StatusOK)),
			expect: &errors.HTTPClientErr{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("expected status 200, got 400"),
			},
		},
		{
			name: "Success/NotAnHTTPClientErr",
			err:  fmt.Errorf("not an HTTPClientErr"),
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			err := errors.AsHTTPClientErr(d.err)
			require.Equal(t, d.expect, err)
		})
	}
}
