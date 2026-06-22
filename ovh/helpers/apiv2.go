package helpers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

type GenericAPIv2Resource struct {
	ResourceStatus string `json:"resourceStatus"`
}

const (
	// paginationCursorHeader is the request header used to ask for a specific page
	// of a cursor-paginated OVHcloud API v2 list endpoint.
	paginationCursorHeader = "X-Pagination-Cursor"
	// paginationCursorNextHeader is the response header carrying the cursor of the
	// next page. Its absence means the last page has been reached.
	paginationCursorNextHeader = "X-Pagination-Cursor-Next"
)

// GetAllPagesV2 retrieves every page of a cursor-paginated OVHcloud API v2 list endpoint.
//
// OVHcloud API v2 list endpoints paginate using cursors: each response may include an
// "X-Pagination-Cursor-Next" header whose value must be sent back as the "X-Pagination-Cursor"
// request header to fetch the following page. When that header is absent from a response, the
// last page has been reached.
//
// The standard client Get helper only unmarshals the response body and does not expose response
// headers, so this helper drives the request manually through NewRequest/Do/UnmarshalResponse
// while preserving the client rate limiter. Items from every page are accumulated and returned.
func GetAllPagesV2[T any](ctx context.Context, c *ovhwrap.Client, path string) ([]T, error) {
	if c == nil {
		return nil, fmt.Errorf("OVH API client is not initialized, check your provider credentials configuration")
	}

	var all []T
	cursor := ""

	for {
		req, err := c.NewRequest(http.MethodGet, path, nil, true)
		if err != nil {
			return nil, fmt.Errorf("failed to build request for %q: %w", path, err)
		}
		req = req.WithContext(ctx)
		if cursor != "" {
			req.Header.Set(paginationCursorHeader, cursor)
		}

		c.RateLimiter.Take()
		resp, err := c.Do(req)
		if err != nil {
			return nil, fmt.Errorf("call to %q failed: %w", path, err)
		}

		var page []T
		if err := c.UnmarshalResponse(resp, &page); err != nil {
			return nil, err
		}
		all = append(all, page...)

		cursor = resp.Header.Get(paginationCursorNextHeader)
		if cursor == "" {
			break
		}
	}

	return all, nil
}

// WaitForAPIv2ResourceStatusReady retries a GET on the given URL until the fetched resource
// is in state "READY". It expects the given URL to target a route that fetches an asynchronous
// resource on APIv2.
func WaitForAPIv2ResourceStatusReady(ctx context.Context, c *ovhwrap.Client, url string) error {
	return retry.RetryContext(ctx, time.Hour, func() *retry.RetryError {
		var resource GenericAPIv2Resource

		if err := c.GetWithContext(ctx, url, &resource); err != nil {
			if ovhErr, ok := err.(ovh.APIError); ok && ovhErr.Code < http.StatusInternalServerError {
				return retry.NonRetryableError(fmt.Errorf("failed to fetch %q : %s", url, ovhErr))
			}
			return retry.RetryableError(fmt.Errorf("call to %q failed, retrying (error: %s)", url, err))
		}

		switch resource.ResourceStatus {
		case "READY":
			return nil
		case "ERROR":
			return retry.NonRetryableError(errors.New("resource is in status ERROR"))
		default:
			return retry.RetryableError(fmt.Errorf("resource not ready, retrying (current status: %s)", resource.ResourceStatus))
		}
	})
}
