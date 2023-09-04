package errors

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type HTTPError struct {
	Err  error
	Code int
}

func ErrorToHTTPCode(c *gin.Context, err error, errorsToCode []HTTPError) {
	for _, e := range errorsToCode {
		if errors.Is(err, e.Err) {
			_ = c.AbortWithError(e.Code, err)
			return
		}
	}

	_ = c.AbortWithError(http.StatusInternalServerError, err)
}

type HTTPClientErr struct {
	Code int   `json:"code"`
	Err  error `json:"body"`
}

func (h HTTPClientErr) Error() string {
	return fmt.Sprintf("request failed with status %d: %s", h.Code, h.Err)
}

func NewHTTPClientErr(res *http.Response, expect ...int) error {
	var (
		body       error
		errMessage error
	)

	if res.Body != nil {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			body = fmt.Errorf("failed to read response body: %w", err)
		} else {
			body = fmt.Errorf(string(bodyBytes))
		}
	}

	if len(expect) == 0 {
		errMessage = fmt.Errorf("got unexpected status code %d", res.StatusCode)
	} else if len(expect) == 1 {
		errMessage = fmt.Errorf("expected status %d, got %d", expect[0], res.StatusCode)
	} else {
		errMessage = fmt.Errorf("expected any status in %v, got %d", expect, res.StatusCode)
	}

	if body != nil {
		errMessage = fmt.Errorf("%s: %w", errMessage, body)
	}

	return &HTTPClientErr{
		Code: res.StatusCode,
		Err:  errMessage,
	}
}

func AsHTTPClientErr(err error) *HTTPClientErr {
	httpClientErr := new(HTTPClientErr)

	if errors.As(err, &httpClientErr) {
		return httpClientErr
	}

	return nil
}
