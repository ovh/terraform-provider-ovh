package ovh

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// vpsPowerActionSchema returns the schema common to all VPS power-action
// one-shot resources (reboot/start/stop/setPassword).
//
// One-shot semantics:
//   - Create: POST the action endpoint and poll the returned task until
//     it reaches a terminal state.
//   - Read:   no-op (the action is an event, not a queryable resource).
//   - Update: never invoked — every field is ForceNew, so changes trigger
//     destroy+create, which re-runs the action.
//   - Delete: no-op — simply forgets the task id.
func vpsPowerActionSchema(extra map[string]*schema.Schema) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"service_name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "The internal name of your VPS service.",
		},
		"triggers": {
			Type:        schema.TypeMap,
			Optional:    true,
			ForceNew:    true,
			Description: "Map of arbitrary string values; changing any value re-runs the action.",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"task_id": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "ID of the OVH task created for this action.",
		},
		"task_state": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Final state of the OVH task (e.g. done).",
		},
	}
	for k, v := range extra {
		s[k] = v
	}
	return s
}

// runVpsPowerAction POSTs the given VPS sub-endpoint, waits for the task to
// reach a terminal state, then writes task_id / task_state into Terraform
// state. Returns the task so callers can read additional fields if needed.
func runVpsPowerAction(d *schema.ResourceData, meta interface{}, subPath string) (*VPSTask, error) {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf(
		"/vps/%s/%s",
		url.PathEscape(serviceName),
		subPath,
	)

	task := &VPSTask{}
	if err := config.OVHClient.Post(endpoint, nil, task); err != nil {
		return nil, fmt.Errorf("Error calling POST %s:\n\t %q", endpoint, err)
	}

	if err := waitForVPSTask(serviceName, task, config.OVHClient); err != nil {
		return task, err
	}

	d.SetId(strconv.FormatInt(task.Id, 10))
	_ = d.Set("task_id", task.Id)
	_ = d.Set("task_state", task.State)

	return task, nil
}

// resourceVpsPowerActionReadNoop is the shared Read implementation: VPS
// power actions are events, not resources, so Read is intentionally empty.
// OVH may purge the task entry after a while; we explicitly do NOT propagate
// 404s as resource deletion to avoid re-running the action on every plan.
func resourceVpsPowerActionReadNoop(d *schema.ResourceData, meta interface{}) error {
	return nil
}

// resourceVpsPowerActionDeleteNoop forgets the task id. We can't actually
// "undo" a reboot or password reset; the resource only existed to record
// that the action ran.
func resourceVpsPowerActionDeleteNoop(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
