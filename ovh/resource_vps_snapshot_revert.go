package ovh

import (
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVPSSnapshotRevert() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPSSnapshotRevertCreate,
		Read:   resourceVPSSnapshotRevertRead,
		Delete: resourceVPSSnapshotRevertDelete,

		Schema: map[string]*schema.Schema{
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
				Description: "Map of arbitrary string values; changing any value re-runs the revert.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"task_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"task_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reverted_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "RFC3339 timestamp of when the revert was issued.",
			},
		},
	}
}

func resourceVPSSnapshotRevertCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/snapshot/revert", url.PathEscape(serviceName))
	task := &VPSTask{}
	if err := config.OVHClient.Post(endpoint, nil, task); err != nil {
		return fmt.Errorf("calling POST %s: %w", endpoint, err)
	}

	if err := waitForVPSTask(serviceName, task, config.OVHClient); err != nil {
		return fmt.Errorf("waiting for VPS snapshot revert task %d: %w", task.Id, err)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	d.SetId(fmt.Sprintf("%s/snapshot-revert/%d", serviceName, task.Id))
	_ = d.Set("task_id", task.Id)
	_ = d.Set("task_state", task.State)
	_ = d.Set("reverted_at", now)
	return nil
}

func resourceVPSSnapshotRevertRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVPSSnapshotRevertDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
