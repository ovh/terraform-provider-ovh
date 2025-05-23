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
