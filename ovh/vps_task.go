package ovh

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

// getVPSTask fetches a single VPS task by its id.
func getVPSTask(serviceName string, taskId int64, c *ovhwrap.Client) (*VPSTask, error) {
	task := &VPSTask{}
	endpoint := fmt.Sprintf(
		"/vps/%s/tasks/%d",
		url.PathEscape(serviceName),
		taskId,
	)
	if err := c.Get(endpoint, task); err != nil {
		return nil, err
	}
	return task, nil
}

// waitForVPSTask polls a vps.Task until it reaches a terminal state.
// Terminal success: "done". Terminal failures: "error", "cancelled", "blocked".
func waitForVPSTask(serviceName string, task *VPSTask, c *ovhwrap.Client) error {
	taskId := task.Id

	refresh := func() (interface{}, string, error) {
		var t *VPSTask
		var fetchErr error

		retryErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			t, err = getVPSTask(serviceName, taskId, c)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 500) {
					return resource.RetryableError(err)
				}
				fetchErr = err
			}
			return nil
		})

		if retryErr != nil {
			return taskId, "", retryErr
		}
		if fetchErr != nil {
			return taskId, "", fetchErr
		}

		log.Printf("[INFO] Pending VPS Task id %d on %s state: %s", taskId, serviceName, t.State)
		return taskId, t.State, nil
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"init", "todo", "doing", "waitingAck", "running", "paused"},
		Target:     []string{"done"},
		Refresh:    refresh,
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	res, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("error waiting for VPS task %s/%d to complete: %s", serviceName, taskId, err)
	}

	// Re-fetch terminal state for explicit failure reporting.
	if final, ferr := getVPSTask(serviceName, taskId, c); ferr == nil {
		switch final.State {
		case "error", "cancelled", "blocked":
			return fmt.Errorf("VPS task %s/%d ended in non-success state %q", serviceName, taskId, final.State)
		}
	}

	_ = res
	return nil
}
