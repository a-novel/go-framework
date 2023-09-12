package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/a-novel/go-framework/errors"
	"github.com/samber/lo"
	"io"
	"net/http"
	"net/url"
)

type HTTPCallConfig struct {
	Path            *url.URL
	Method          string
	Body            interface{}
	Headers         map[string]string
	SuccessStatuses []int
	Client          *http.Client
}

func MakeHTTPCall(ctx context.Context, cfg HTTPCallConfig, output interface{}) error {
	var reqBody io.Reader

	if cfg.Body != nil {
		bodyBytes, err := json.Marshal(cfg.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		reqBody = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, cfg.Method, cfg.Path.String(), reqBody)
	if err != nil {
		return fmt.Errorf("failed to initiate request: %w", err)
	}

	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	res, err := cfg.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	defer res.Body.Close()

	if _, ok := lo.Find(cfg.SuccessStatuses, func(item int) bool {
		return res.StatusCode == item
	}); !ok {
		return errors.NewHTTPClientErr(res, cfg.SuccessStatuses...)
	}

	if output == nil {
		return nil
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, output); err != nil {
		return fmt.Errorf("failed to parse response body: %w", err)
	}

	return nil
}
