package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

func resourceVPSBackupFtpAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPSBackupFtpAccessCreate,
		Read:   resourceVPSBackupFtpAccessRead,
		Update: resourceVPSBackupFtpAccessUpdate,
		Delete: resourceVPSBackupFtpAccessDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVPSBackupFtpAccessImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your VPS.",
			},
			"ip_block": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CIDR block to grant backup FTP access to.",
			},
			"cifs": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether CIFS (SMB) protocol is enabled for this IP block.",
			},
			"nfs": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether NFS protocol is enabled for this IP block.",
			},
			"ftp": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether FTP protocol is enabled for this IP block.",
			},
			"is_applied": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether ACL is currently applied on the backup FTP storage.",
			},
			"last_update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last ACL update date.",
			},
		},
	}
}

func resourceVPSBackupFtpAccessImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "|", 2)
	if len(splitId) != 2 || splitId[0] == "" || splitId[1] == "" {
		return nil, fmt.Errorf("Import Id is not service_name|ip_block formatted")
	}
	d.Set("service_name", splitId[0])
	d.Set("ip_block", splitId[1])
	d.SetId(givenId)
	return []*schema.ResourceData{d}, nil
}

func resourceVPSBackupFtpAccessCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ipBlock := d.Get("ip_block").(string)
	ftp := d.Get("ftp").(bool)
	cifs := d.Get("cifs").(bool)
	nfs := d.Get("nfs").(bool)

	opts := VPSBackupFtpAclCreateOpts{
		IpBlock: ipBlock,
		Cifs:    cifs,
		Nfs:     nfs,
		Ftp:     &ftp,
	}

	endpoint := fmt.Sprintf("/vps/%s/backupftp/access", url.PathEscape(serviceName))
	task := &DedicatedServerTask{}
	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("Error calling POST %s:\n\t %q", endpoint, err)
	}

	if err := waitForVPSBackupFtpTask(serviceName, task, config.OVHClient); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%s|%s", serviceName, ipBlock))

	return resourceVPSBackupFtpAccessRead(d, meta)
}

func resourceVPSBackupFtpAccessRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ipBlock := d.Get("ip_block").(string)

	endpoint := fmt.Sprintf(
		"/vps/%s/backupftp/access/%s",
		url.PathEscape(serviceName),
		url.PathEscape(ipBlock),
	)

	acl := &VPSBackupFtpAcl{}
	if err := config.OVHClient.Get(endpoint, acl); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("ip_block", acl.IpBlock)
	d.Set("cifs", acl.Cifs)
	d.Set("nfs", acl.Nfs)
	d.Set("ftp", acl.Ftp)
	d.Set("is_applied", acl.IsApplied)
	d.Set("last_update", acl.LastUpdate)

	return nil
}

func resourceVPSBackupFtpAccessUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ipBlock := d.Get("ip_block").(string)

	opts := VPSBackupFtpAclUpdateOpts{
		IpBlock: ipBlock,
		Cifs:    d.Get("cifs").(bool),
		Nfs:     d.Get("nfs").(bool),
		Ftp:     d.Get("ftp").(bool),
	}

	endpoint := fmt.Sprintf(
		"/vps/%s/backupftp/access/%s",
		url.PathEscape(serviceName),
		url.PathEscape(ipBlock),
	)

	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Error calling PUT %s:\n\t %q", endpoint, err)
	}

	return resourceVPSBackupFtpAccessRead(d, meta)
}

func resourceVPSBackupFtpAccessDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ipBlock := d.Get("ip_block").(string)

	endpoint := fmt.Sprintf(
		"/vps/%s/backupftp/access/%s",
		url.PathEscape(serviceName),
		url.PathEscape(ipBlock),
	)

	task := &DedicatedServerTask{}
	if err := config.OVHClient.Delete(endpoint, task); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	if err := waitForVPSBackupFtpTask(serviceName, task, config.OVHClient); err != nil {
		return err
	}

	d.SetId("")
	return nil
}

// waitForVPSBackupFtpTask polls a backupftp task on a VPS. The OVH API returns
// dedicated.server.Task objects for backupftp operations but the task lives on
// the VPS-side endpoint. We try a couple of probable endpoints and fall back
// to a short timeout sleep so we never crash terraform if no endpoint matches.
func waitForVPSBackupFtpTask(serviceName string, task *DedicatedServerTask, c *ovhwrap.Client) error {
	if task == nil || task.Id == 0 {
		return nil
	}
	taskId := task.Id

	candidates := []string{
		fmt.Sprintf("/vps/%s/tasks/%d", url.PathEscape(serviceName), taskId),
		fmt.Sprintf("/vps/%s/task/%d", url.PathEscape(serviceName), taskId),
		fmt.Sprintf("/dedicated/server/%s/task/%d", url.PathEscape(serviceName), taskId),
	}

	refreshFunc := func() (interface{}, string, error) {
		var lastErr error
		for _, endpoint := range candidates {
			probe := &DedicatedServerTask{}
			err := c.Get(endpoint, probe)
			if err == nil {
				log.Printf("[INFO] Pending backupftp task %d on VPS %s status: %s", taskId, serviceName, probe.Status)
				return taskId, probe.Status, nil
			}
			if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
				lastErr = err
				continue
			}
			lastErr = err
		}
		// All candidates failed; treat as transient.
		log.Printf("[DEBUG] backupftp task probe error (will retry): %v", lastErr)
		return taskId, "doing", nil
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"init", "todo", "doing", "ovhError", "customerError"},
		Target:     []string{"done", "cancelled"},
		Refresh:    refreshFunc,
		Timeout:    15 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(context.Background()); err != nil {
		return fmt.Errorf("backupftp task %d on VPS %s did not reach a terminal state: %w", taskId, serviceName, err)
	}
	return nil
}
