package ovh

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceVPSAutomatedBackupRestore() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPSAutomatedBackupRestoreCreate,
		Read:   resourceVPSAutomatedBackupRestoreRead,
		Delete: resourceVPSAutomatedBackupRestoreDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Hour),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your VPS",
			},
			"restore_point": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Restore point (RFC3339 datetime) to restore from",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Restore type: file (mount as additional access) or full",
				ValidateFunc: helpers.ValidateEnum([]string{"file", "full"}),
			},
			"change_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to change the root password after a full restore",
			},
			"task_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the underlying VPS task that performed the restore",
			},
		},
	}
}

func resourceVPSAutomatedBackupRestoreCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	opts := &VPSAutomatedBackupRestoreOpts{
		RestorePoint: d.Get("restore_point").(string),
		Type:         d.Get("type").(string),
	}
	if v, ok := d.GetOkExists("change_password"); ok {
		b := v.(bool)
		opts.ChangePassword = &b
	}

	endpoint := fmt.Sprintf("/vps/%s/automatedBackup/restore", url.PathEscape(serviceName))
	task := &VPSTask{}
	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("error calling POST %s: %w", endpoint, err)
	}

	if err := waitForVPSTask(serviceName, task, config.OVHClient); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, opts.RestorePoint))
	d.Set("task_id", task.Id)
	return resourceVPSAutomatedBackupRestoreRead(d, meta)
}

func resourceVPSAutomatedBackupRestoreRead(d *schema.ResourceData, meta interface{}) error {
	// The restore action is a fire-and-forget task; we keep state stable as
	// long as the user doesn't change the (ForceNew) restore_point/type.
	return nil
}

func resourceVPSAutomatedBackupRestoreDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	restorePoint := d.Get("restore_point").(string)

	// Only "file" restores create attached backups that can be detached.
	// For "full" restores there is nothing to detach: simply drop state.
	if d.Get("type").(string) != "file" {
		d.SetId("")
		return nil
	}

	endpoint := fmt.Sprintf("/vps/%s/automatedBackup/detachBackup", url.PathEscape(serviceName))
	opts := &VPSAutomatedBackupDetachOpts{RestorePoint: restorePoint}
	task := &VPSTask{}
	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			// Already detached.
			log.Printf("[INFO] automated backup %s already detached from %s", restorePoint, serviceName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error calling POST %s: %w", endpoint, err)
	}

	if task.Id != 0 {
		if err := waitForVPSTask(serviceName, task, config.OVHClient); err != nil {
			return err
		}
	}

	d.SetId("")
	return nil
}
