package ovh

import (
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
)

// WaitArchivedHostingPrivateDabaseTask wait for a task to become archived in the API (aka 404)
func WaitArchivedHostingPrivateDabaseTask(client *OVHClient, endpoint string, timeout time.Duration) error {
	return resource.Retry(timeout, func() *resource.RetryError {
		err := client.Get(endpoint, nil)
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code != 404 {
			return resource.NonRetryableError(err)
		}
		if err == nil {
			return resource.RetryableError(errors.New("not archived yet"))
		}
		return nil
	})
}
