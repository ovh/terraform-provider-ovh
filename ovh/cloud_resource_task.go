package ovh

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

// CloudResourceTaskError is a single error reported on an asynchronous task.
type CloudResourceTaskError struct {
	Message string `json:"message,omitempty"`
}

// CloudResourceTask mirrors apiv2 model.ResourceTask: an asynchronous operation
// reconciling a Public Cloud API v2 resource towards its target spec. Used for
// wait/diagnostics only, never stored in Terraform state.
type CloudResourceTask struct {
	Id     string                   `json:"id,omitempty"`
	Type   string                   `json:"type,omitempty"`
	Status string                   `json:"status,omitempty"`
	Link   string                   `json:"link,omitempty"`
	Errors []CloudResourceTaskError `json:"errors,omitempty"`
}

// cloudResourceErrorFromTasks builds the error surfaced when a resource reaches
// the terminal ERROR status, joining the reasons reported by its failed tasks.
func cloudResourceErrorFromTasks(resourceKind, resourceId string, tasks []CloudResourceTask) error {
	var messages []string
	for _, task := range tasks {
		for _, taskErr := range task.Errors {
			if taskErr.Message != "" {
				messages = append(messages, taskErr.Message)
			}
		}
	}
	if len(messages) == 0 {
		return fmt.Errorf("%s %s entered ERROR state (no error detail reported by the API)", resourceKind, resourceId)
	}
	return fmt.Errorf("%s %s entered ERROR state: %s", resourceKind, resourceId, strings.Join(messages, "; "))
}

// cloudResourceStatus is the slice of the envelope the pollers decode; callers
// re-GET the full resource once the wait is over.
type cloudResourceStatus struct {
	ResourceStatus string              `json:"resourceStatus"`
	CurrentTasks   []CloudResourceTask `json:"currentTasks,omitempty"`
}

// cloudResourceBusy lists the statuses of a resource still reconciling.
var cloudResourceBusy = []string{"CREATION_PLANNED", "CREATING", "UPDATING", "OUT_OF_SYNC"}

// waitForCloudResourceStatus polls endpoint until resourceStatus reaches one of target.
// A 404 reads as DELETED; ERROR is terminal and surfaces the failed tasks' reasons.
func waitForCloudResourceStatus(ctx context.Context, client *ovhwrap.Client, endpoint, kind, id string, pending, target []string) error {
	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     target,
		Refresh:    cloudResourceStatusRefresh(ctx, client, endpoint, kind, id),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// waitForCloudResourceUpdate waits for the update to start, then for READY: the status can lag behind the PUT.
func waitForCloudResourceUpdate(ctx context.Context, client *ovhwrap.Client, endpoint, kind, id string) error {
	start := &retry.StateChangeConf{
		Pending:    []string{"READY"},
		Target:     cloudResourceBusy,
		Refresh:    cloudResourceStatusRefresh(ctx, client, endpoint, kind, id),
		Timeout:    time.Minute,
		MinTimeout: 2 * time.Second,
	}
	if _, err := start.WaitForStateContext(ctx); err != nil {
		if _, ok := err.(*retry.TimeoutError); !ok {
			return err
		}
	}
	return waitForCloudResourceStatus(ctx, client, endpoint, kind, id, cloudResourceBusy, []string{"READY"})
}

// cloudResourceStatusRefresh reads a 404 as DELETED.
func cloudResourceStatusRefresh(ctx context.Context, client *ovhwrap.Client, endpoint, kind, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res := &cloudResourceStatus{}
		if err := client.GetWithContext(ctx, endpoint, res); err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
				return res, "DELETED", nil
			}
			return res, "", err
		}
		if res.ResourceStatus == "ERROR" {
			return res, res.ResourceStatus, cloudResourceErrorFromTasks(kind, id, res.CurrentTasks)
		}
		return res, res.ResourceStatus, nil
	}
}
