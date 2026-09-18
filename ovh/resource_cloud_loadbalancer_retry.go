package ovh

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

// The v2 load balancer API serializes writes per load balancer: while an
// operation is being applied, the load balancer leaves the ACTIVE state and any
// other write on it or on one of its children is rejected with a 409. Terraform
// creates, updates and deletes sibling resources (listeners, pools, members, L7
// policies) of the same load balancer in parallel, so these collisions happen
// routinely and are purely transient — the API itself asks the caller to retry.
const (
	// cloudLoadbalancerBusyTimeout bounds how long a write is retried while the
	// load balancer applies a previous operation. It matches the timeout of the
	// waitFor*Ready helpers, since that is how long a single operation may take.
	cloudLoadbalancerBusyTimeout = 20 * time.Minute

	// cloudLoadbalancerBusyClass is the error class returned when the targeted
	// resource, or its parent, is not in a state allowing the write.
	cloudLoadbalancerBusyClass = "Client::Conflict::ResourceInInvalidState"

	// cloudLoadbalancerBusyHint narrows that class down to the transient case
	// ("The parent load balancer is not in ACTIVE state. Wait for ongoing
	// operations to complete and retry."), so a permanent conflict reported
	// under the same class is still surfaced right away.
	cloudLoadbalancerBusyHint = "ACTIVE state"
)

// isCloudLoadbalancerBusy reports whether err is the API telling us the load
// balancer is still applying a previous operation.
func isCloudLoadbalancerBusy(err error) bool {
	var apiErr *ovh.APIError
	if !errors.As(err, &apiErr) {
		return false
	}

	return apiErr.Code == http.StatusConflict &&
		apiErr.Class == cloudLoadbalancerBusyClass &&
		strings.Contains(apiErr.Message, cloudLoadbalancerBusyHint)
}

// retryWhileCloudLoadbalancerBusy calls call, retrying for as long as the API
// reports the load balancer is busy. Any other error, including the one left by
// the last attempt when the timeout is reached, is returned unwrapped so that
// callers can keep inspecting it as an *ovh.APIError.
func retryWhileCloudLoadbalancerBusy(ctx context.Context, call func() error) error {
	return retry.RetryContext(ctx, cloudLoadbalancerBusyTimeout, func() *retry.RetryError {
		err := call()
		switch {
		case err == nil:
			return nil
		case isCloudLoadbalancerBusy(err):
			return retry.RetryableError(err)
		default:
			return retry.NonRetryableError(err)
		}
	})
}

// cloudLoadbalancerPost performs a POST on a load balancer endpoint, retrying
// while the load balancer is busy.
func cloudLoadbalancerPost(ctx context.Context, client *ovhwrap.Client, endpoint string, reqBody, resType interface{}) error {
	return retryWhileCloudLoadbalancerBusy(ctx, func() error {
		return client.PostWithContext(ctx, endpoint, reqBody, resType)
	})
}

// cloudLoadbalancerPut performs a PUT on a load balancer endpoint, retrying
// while the load balancer is busy.
func cloudLoadbalancerPut(ctx context.Context, client *ovhwrap.Client, endpoint string, reqBody, resType interface{}) error {
	return retryWhileCloudLoadbalancerBusy(ctx, func() error {
		return client.PutWithContext(ctx, endpoint, reqBody, resType)
	})
}

// cloudLoadbalancerDelete performs a DELETE on a load balancer endpoint,
// retrying while the load balancer is busy.
func cloudLoadbalancerDelete(ctx context.Context, client *ovhwrap.Client, endpoint string, resType interface{}) error {
	return retryWhileCloudLoadbalancerBusy(ctx, func() error {
		return client.DeleteWithContext(ctx, endpoint, resType)
	})
}
