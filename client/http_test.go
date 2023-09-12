package client_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/a-novel/go-framework/client"
	"github.com/a-novel/go-framework/errors"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestMakeHTTPCall(t *testing.T) {
	type foo struct {
		Foo string `json:"foo"`
	}

	data := []struct {
		name string

		apiMethod string
		apiPath   string
		apiResp   interface{}
		apiStatus int

		config client.HTTPCallConfig
		output interface{}

		expect    interface{}
		expectErr error
	}{
		{
			name:      "Success",
			apiMethod: http.MethodGet,
			apiPath:   "/test",
			apiResp:   `{"foo": "bar"}`,
			apiStatus: http.StatusOK,
			config: client.HTTPCallConfig{
				Method:          http.MethodGet,
				SuccessStatuses: []int{http.StatusOK},
			},
		},
		{
			name:      "Success/WithRequestBody",
			apiMethod: http.MethodGet,
			apiPath:   "/test",
			apiResp:   `{"foo": "bar"}`,
			apiStatus: http.StatusOK,
			config: client.HTTPCallConfig{
				Method:          http.MethodGet,
				Body:            map[string]interface{}{"message": "hello"},
				SuccessStatuses: []int{http.StatusOK},
			},
		},
		{
			name:      "Success/WithHeaders",
			apiMethod: http.MethodGet,
			apiPath:   "/test",
			apiResp:   `{"foo": "bar"}`,
			apiStatus: http.StatusOK,
			config: client.HTTPCallConfig{
				Method:          http.MethodGet,
				SuccessStatuses: []int{http.StatusOK},
				Headers: map[string]string{
					"X-Test": "test",
				},
			},
		},
		{
			name:      "Success/ParseResponse",
			apiMethod: http.MethodGet,
			apiPath:   "/test",
			apiResp:   `{"foo": "bar"}`,
			apiStatus: http.StatusOK,
			output:    new(foo),
			config: client.HTTPCallConfig{
				Method:          http.MethodGet,
				SuccessStatuses: []int{http.StatusOK},
			},
			expect: &foo{Foo: "bar"},
		},
		{
			name:      "Error/WrongStatus",
			apiMethod: http.MethodGet,
			apiPath:   "/test",
			apiResp:   "forbidden",
			apiStatus: http.StatusForbidden,
			config: client.HTTPCallConfig{
				Method:          http.MethodGet,
				SuccessStatuses: []int{http.StatusOK},
			},
			expectErr: &errors.HTTPClientErr{
				Code: http.StatusForbidden,
				Err:  fmt.Errorf("expected status 200, got 403: forbidden"),
			},
		},
	}

	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != d.apiMethod {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}

				if r.URL.Path != d.apiPath {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				if d.config.Body != nil {
					bodyBytes, err := io.ReadAll(r.Body)
					require.NoError(t, err)

					var body interface{}

					err = json.Unmarshal(bodyBytes, &body)
					require.NoError(t, err)

					require.Equal(t, d.config.Body, body)
				}

				if d.config.Headers != nil {
					for k, v := range d.config.Headers {
						require.Equal(t, v, r.Header.Get(k))
					}
				}

				if d.apiStatus != 0 {
					w.WriteHeader(d.apiStatus)
				}

				if d.apiResp != nil {
					_, _ = fmt.Fprint(w, d.apiResp)
				}
			}))

			defer srv.Close()

			srvURL, err := new(url.URL).Parse(srv.URL)
			require.NoError(t, err)

			err = client.MakeHTTPCall(context.Background(), client.HTTPCallConfig{
				Path:            srvURL.JoinPath(d.apiPath),
				Method:          d.apiMethod,
				Body:            d.config.Body,
				Headers:         d.config.Headers,
				SuccessStatuses: d.config.SuccessStatuses,
				Client:          srv.Client(),
			}, d.output)

			if d.expectErr != nil {
				require.Equal(t, d.expectErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, d.expect, d.output)
		})
	}
}
