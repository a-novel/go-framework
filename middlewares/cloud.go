package middlewares

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	"google.golang.org/api/oauth2/v2"
	"strings"
)

const (
	GoogleCloudSchedulerUserAgent = "Google-Cloud-Scheduler"
)

// ProtectCloud prevents access from public context, and only allows cloud environment to access the endpoint.
func ProtectCloud(ctx context.Context, authorization, userAgent string, allowedUserAgents []string) error {
	// https://jackcuthbert.dev/blog/verifying-google-cloud-scheduler-requests-in-cloud-run-with-typescript
	if _, ok := lo.Find(allowedUserAgents, func(item string) bool {
		return item == userAgent
	}); !ok {
		return fmt.Errorf("bad user agent: expected %q, got %q", "Google-Cloud-Scheduler", userAgent)
	}

	// https://stackoverflow.com/questions/53181297/verify-http-request-from-google-cloud-scheduler
	if authorization == "" {
		return fmt.Errorf("missing authorization header")
	}

	idToken := strings.Split(authorization, "Bearer ")[0]

	authenticator, err := oauth2.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire authenticator: %w", err)
	}

	_, err = authenticator.Tokeninfo().IdToken(idToken).Do()
	if err != nil {
		return fmt.Errorf("failed to retrieve token information: %w", err)
	}

	return nil
}
