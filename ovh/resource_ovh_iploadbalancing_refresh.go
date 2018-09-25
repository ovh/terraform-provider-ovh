package ovh

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

type IPLoadbalancingRefreshTask struct {
	CreationDate string   `json:"creationDate"`
	Status       string   `json:"status"`
	Progress     int      `json:"progress"`
	Action       string   `json:"action"`
	ID           int      `json:"id"`
	DoneDate     string   `json:"doneDate"`
	Zones        []string `json:"zones"`
}

type IPLoadbalancingRefreshPending struct {
	Number int    `json:"number"`
	Zone   string `json:"zone"`
}

type IPLoadbalancingRefreshPendings []IPLoadbalancingRefreshPending

func resourceIPLoadbalancingRefresh() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPLoadbalancingRefreshCreate,
		Read:   resourceIPLoadbalancingRefreshRead,
		Delete: resourceIPLoadbalancingRefreshDelete,

		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"keepers": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceIPLoadbalancingRefreshCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	service := d.Get("service_name").(string)

	// verify if there are no active tasks for the loadbalancer
	// at the moment and wait till finished if there are any

	stateConf := &resource.StateChangeConf{
		Target: []string{"empty"},
		Refresh: func() (interface{}, string, error) {
			for _, state := range []string{"todo", "doing"} {
				taskResp := &[]int{}
				endpoint := fmt.Sprintf("/ipLoadbalancing/%s/task?action=refreshIplb&status=%s", service, state)
				err := config.OVHClient.Get(endpoint, taskResp)
				if err != nil {
					return d, "error", fmt.Errorf("calling GET %s :\n\t %s", endpoint, err.Error())
				}
				if len(*taskResp) > 0 {
					return d, "exists", nil
				}
			}
			return d, "empty", nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for IPLoadbalancer tasks to finish: %s", err)
	}

	// verify if there are any outstanding changes to refresh

	checkResp := &IPLoadbalancingRefreshPendings{}
	endpoint := fmt.Sprintf("/ipLoadbalancing/%s/pendingChanges", service)

	err = config.OVHClient.Get(endpoint, checkResp)
	if err != nil {
		return fmt.Errorf("calling GET %s :\n\t %s", endpoint, err.Error())
	}

	// no changes detected, return successfull creation/refresh
	if len(*checkResp) == 0 {
		d.SetId(service)
		return nil
	}

	// proceed with refresh

	resp := &IPLoadbalancingRefreshTask{}
	endpoint = fmt.Sprintf("/ipLoadbalancing/%s/refresh", service)

	err = config.OVHClient.Post(endpoint, nil, resp)
	if err != nil {
		return fmt.Errorf("calling POST %s :\n\t %s", endpoint, err.Error())
	}

	stateConf = &resource.StateChangeConf{
		Target: []string{"done"},
		Refresh: func() (interface{}, string, error) {
			endpoint := fmt.Sprintf("/ipLoadbalancing/%s/task/%d", service, resp.ID)
			stateResp := &IPLoadbalancingRefreshTask{}
			err := config.OVHClient.Get(endpoint, stateResp)
			if err != nil {
				return nil, "", err
			}
			return d, stateResp.Status, nil
		},
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for IPLoadbalancer refresh: %s", err)
	}

	d.SetId(service)

	return nil
}

func resourceIPLoadbalancingRefreshRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIPLoadbalancingRefreshDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
