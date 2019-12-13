package ovh

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/ovh/go-ovh/ovh"
)

func waitForDedicatedServerTask(serviceName string, task *DedicatedServerTask, c *ovh.Client) error {
	taskId := task.Id

	refreshFunc := func() (interface{}, string, error) {
		task, err := getDedicatedServerTask(serviceName, taskId, c)
		if err != nil {
			return taskId, "", err
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
