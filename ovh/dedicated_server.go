package ovh

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/ovh/go-ovh/ovh"
)

func dedicatedServerReboot(serviceName string, c *ovh.Client) error {
	endpoint := fmt.Sprintf(
		"/dedicated/server/%s/reboot",
		url.PathEscape(serviceName),
	)

	task := &DedicatedServerTask{}

	if err := c.Post(endpoint, nil, task); err != nil {
		return fmt.Errorf("Error calling PUT %s:\n\t %q", endpoint, err)
	}

	if err := waitForDedicatedServerTask(serviceName, task, c); err != nil {
		return err
	}

	return nil
}

func waitForDedicatedServerTask(serviceName string, task *DedicatedServerTask, c *ovh.Client) error {
	taskId := task.Id

	refreshFunc := func() (interface{}, string, error) {
		task := &DedicatedServerTask{}
		endpoint := fmt.Sprintf(
			"/dedicated/server/%s/task/%d",
			url.PathEscape(serviceName),
			taskId,
		)

		if err := c.Get(endpoint, task); err != nil {
			return taskId, "", fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
		}

		log.Printf("[INFO] Pending Task id %d on Dedicated %s status: %s", task.Id, serviceName, task.Status)
		return taskId, task.Status, nil
	}

	log.Printf("[INFO] Waiting for Dedicated Server Task id %s/%d", serviceName, taskId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"init", "todo", "doing"},
		Target:     []string{"done"},
		Refresh:    refreshFunc,
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Dedicated Server task %s/%d to complete: %s", serviceName, taskId, err)
	}

	return nil
}
