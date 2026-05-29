package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVPSBackupFtpPassword() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPSBackupFtpPasswordCreate,
		Read:   resourceVPSBackupFtpPasswordRead,
		Delete: resourceVPSBackupFtpPasswordDelete,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your VPS.",
			},
			"triggers": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Arbitrary map of values that, when changed, will trigger a new password rotation.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// Computed
			"task_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Id of the dedicated.server.Task returned by the API.",
			},
			"task_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Final state of the rotation task (done, error, etc.).",
			},
		},
	}
}

func resourceVPSBackupFtpPasswordCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf(
		"/vps/%s/backupftp/password",
		url.PathEscape(serviceName),
	)

	task := &DedicatedServerTask{}
	if err := config.OVHClient.Post(endpoint, nil, task); err != nil {
		return fmt.Errorf("Error calling POST %s:\n\t %q", endpoint, err)
	}

	d.SetId(fmt.Sprintf("%d", task.Id))
	d.Set("task_id", task.Id)

	if err := waitForVPSBackupFtpTask(serviceName, task, config.OVHClient); err != nil {
		return err
	}
	d.Set("task_state", task.Status)
	return nil
}

func resourceVPSBackupFtpPasswordRead(d *schema.ResourceData, meta interface{}) error {
	// One-shot resource: there is nothing to refresh once the rotation has
	// happened. Treat as a no-op so Terraform does not try to recreate it
	// when the task entry is purged from the API.
	return nil
}

func resourceVPSBackupFtpPasswordDelete(d *schema.ResourceData, meta interface{}) error {
	// Cannot un-rotate a password. Just drop the id from state.
	d.SetId("")
	return nil
}

// waitForVPSBackupFtpTask polls the OVH API for a backup FTP task. The
// /vps/{sn}/backupftp/password endpoint returns a dedicated.server.Task,
