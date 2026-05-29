package ovh

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func resourceVPSVeeamRestore() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPSVeeamRestoreCreate,
		Read:   resourceVPSVeeamRestoreRead,
		Delete: resourceVPSVeeamRestoreDelete,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your VPS",
			},
			"restore_point_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the Veeam restore point to mount",
			},
			"full": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: "Whether to perform a full restore (true) or to only expose the backup (false)",
			},
			"export": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "Export protocol used to expose the restore point (nfs|smb)",
				ValidateFunc: helpers.ValidateEnum([]string{"nfs", "smb"}),
			},
			"change_password": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Whether to change the VPS administration password (only relevant when full=true)",
			},

			// Computed
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the restored backup",
			},
			"access_nfs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "NFS access information",
			},
			"access_smb": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SMB access information",
			},
			"task_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the restore task",
			},
		},
	}
}

func resourceVPSVeeamRestoreCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	restorePointId := int64(d.Get("restore_point_id").(int))
	full := d.Get("full").(bool)
	export := d.Get("export").(string)

	opts := &VpsVeeamRestoreOpts{
		Full:   full,
		Export: export,
	}

	if v, ok := d.GetOkExists("change_password"); ok {
		cp := v.(bool)
		if !full && cp {
			return errors.New("change_password can only be set when full = true")
		}
		// Only forward change_password when full is true; the OVH API
		// rejects the field otherwise.
		if full {
			opts.ChangePassword = &cp
		}
	}

	endpoint := fmt.Sprintf(
		"/vps/%s/veeam/restorePoints/%d/restore",
		url.PathEscape(serviceName),
		restorePointId,
	)

	task := &VPSTask{}
	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("error calling POST %s: %s", endpoint, err)
	}

	d.Set("task_id", task.Id)

	if err := waitForVPSTask(serviceName, task, config.OVHClient); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s/%d", serviceName, restorePointId))

	return resourceVPSVeeamRestoreRead(d, meta)
}

func resourceVPSVeeamRestoreRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/veeam/restoredBackup", url.PathEscape(serviceName))
	rb := &VpsVeeamRestoredBackup{}
	if err := config.OVHClient.Get(endpoint, rb); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("restore_point_id", rb.RestorePointId)
	d.Set("state", rb.State)
	d.Set("access_nfs", rb.AccessInfos.Nfs)
	d.Set("access_smb", rb.AccessInfos.Smb)
	return nil
}

func resourceVPSVeeamRestoreDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/vps/%s/veeam/restoredBackup", url.PathEscape(serviceName))
	task := &VPSTask{}
	if err := config.OVHClient.Delete(endpoint, task); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error calling DELETE %s: %s", endpoint, err)
	}

	if task.Id != 0 {
		if err := waitForVPSTask(serviceName, task, config.OVHClient); err != nil {
			return err
		}
	}

	// As an extra safety net, wait for the restoredBackup endpoint to
	// return a 404, which indicates the restore has actually been removed.
	err := resource.Retry(15*time.Minute, func() *resource.RetryError {
		readErr := config.OVHClient.Get(endpoint, nil)
		if readErr != nil {
			if errOvh, ok := readErr.(*ovh.APIError); ok && errOvh.Code == 404 {
				return nil
			}
			return resource.NonRetryableError(readErr)
		}
		return resource.RetryableError(errors.New("waiting for restored backup to be unmounted"))
	})

	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
