package ovh

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/ovh/go-ovh/ovh"
)

func waitForDedicatedServerTask(serviceName string, task *DedicatedServerTask, c *ovh.Client) error {
	taskId := task.Id

	refreshFunc := func() (interface{}, string, error) {
		var taskErr error
		var task *DedicatedServerTask

		// The Dedicated Server API often returns 500/404 errors
		// in such case we retry to retrieve task status
		// 404 may happen because of some inconsistency between the
		// api endpoint call and the target region executing the task
		retryErr := resource.Retry(5*time.Minute, func() *resource.RetryError {
			var err error
			task, err = getDedicatedServerTask(serviceName, taskId, c)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 500) {
					return resource.RetryableError(err)
				}
				// other error dont retry and fail
				taskErr = err
			}
			return nil
		})

		if retryErr != nil {
			return taskId, "", retryErr
		}

		if taskErr != nil {
			return taskId, "", taskErr
		}

		log.Printf("[INFO] Pending Task id %d on Dedicated %s status: %s", taskId, serviceName, task.Status)
		return taskId, task.Status, nil
	}

	log.Printf("[INFO] Waiting for Dedicated Server Task id %s/%d", serviceName, taskId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"init", "todo", "doing"},
		Target:     []string{"done"},
		Refresh:    refreshFunc,
		Timeout:    45 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Dedicated Server task %s/%d to complete: %s", serviceName, taskId, err)
	}

	return nil
}

func getDedicatedServerTask(serviceName string, taskId int64, c *ovh.Client) (*DedicatedServerTask, error) {
	task := &DedicatedServerTask{}
	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/task/%d",
		url.PathEscape(serviceName),
		taskId,
	)

	if err := c.Get(endpoint, task); err != nil {
		return nil, err
	}

	return task, nil
}
