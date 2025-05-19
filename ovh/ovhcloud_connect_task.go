package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
)

type OccTask struct {
	Function   string `json:"function"`
	TaskID     int    `json:"id"`
	ResourceID int64  `json:"resourceId"`
	Status     string `json:"status"`
}

func waitForOccTask(ctx context.Context, client *ovh.Client, serviceName string, taskId int) error {
	endpoint := fmt.Sprintf("/ovhCloudConnect/%s/task/%d", url.PathEscape(serviceName), taskId)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"todo", "doing"},
		Target:  []string{"done"},
		Refresh: func() (result interface{}, state string, err error) {
			var task OccTask

			if err := client.GetWithContext(ctx, endpoint, &task); err != nil {
				log.Printf("[ERROR] couldn't fetch task %d for occ %s:\n\t%s\n", taskId, serviceName, err.Error())
				return nil, "error", err
			}

			return task, task.Status, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      30 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	if err != nil {
		return fmt.Errorf("error waiting for occ: %s task: %d to complete:\n\t%s", serviceName, taskId, err.Error())
	}

	return nil
}
