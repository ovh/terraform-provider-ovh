package ovh

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

type DedicatedCloudTask struct {
	TaskId int    `json:"taskId"`
	State  string `json:"state"`
	UserId int64  `json:"userId"`
}

func waitForDedicatedCloudTask(ctx context.Context, client *ovhwrap.Client, serviceName string, taskId int) (*DedicatedCloudTask, error) {
	endpoint := fmt.Sprintf("/dedicatedCloud/%s/task/%d", url.PathEscape(serviceName), taskId)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"todo", "doing", "waitingTodo", "waitingForChilds", "toCreate", "toCancel", "fixing"},
		Target:  []string{"done"},
		Refresh: func() (result interface{}, state string, err error) {
			var task DedicatedCloudTask

			if err := client.GetWithContext(ctx, endpoint, &task); err != nil {
				return nil, "error", err
			}

			return &task, task.State, nil
		},
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for dedicatedCloud %s task %d to complete:\n\t%s", serviceName, taskId, err.Error())
	}

	return result.(*DedicatedCloudTask), nil
}
