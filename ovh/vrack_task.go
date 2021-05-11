package ovh

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ovh/go-ovh/ovh"
)

func waitForVrackTask(task *VrackTask, c *ovh.Client) error {
	vrackId := task.ServiceName
	taskId := task.Id

	refreshFunc := func() (interface{}, string, error) {
		task := &VrackTask{}
		endpoint := fmt.Sprintf(
			"/vrack/%s/task/%d",
			url.PathEscape(vrackId),
			taskId,
		)

		err := c.Get(endpoint, task)
		if err != nil {
			if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
				log.Printf("[DEBUG] Task id %d on Vrack %s completed", taskId, vrackId)
				return taskId, "completed", nil
			} else {
				return taskId, "", err
			}
		}

		log.Printf("[DEBUG] Pending Task id %d on Vrack %s status: %s", task.Id, vrackId, task.Status)
		return taskId, task.Status, nil
	}

	log.Printf("[DEBUG] Waiting for Vrack Task id %d: Vrack %s ", taskId, vrackId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"init", "todo", "doing"},
		Target:     []string{"completed"},
		Refresh:    refreshFunc,
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for vrack task %s/%d to complete: %s", vrackId, taskId, err)
	}

	return nil
}
