package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

func resourceVPSSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPSSnapshotCreate,
		ReadContext:   resourceVPSSnapshotRead,
		UpdateContext: resourceVPSSnapshotUpdate,
		DeleteContext: resourceVPSSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your VPS. The VPS must have the 'snapshot' option subscribed.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable description of the snapshot.",
			},
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "RFC3339 creation date of the snapshot.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "OVHcloud region where the snapshot is stored.",
			},
		},
	}
}

// surfaceOptionError tries to detect "snapshot option not subscribed"
// API errors and reformulate them as something users can act on.
func surfaceOptionError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if errOvh, ok := err.(*ovh.APIError); ok {
		if errOvh.Code == 460 ||
			strings.Contains(strings.ToLower(errOvh.Message), "option") ||
			strings.Contains(strings.ToLower(errOvh.Message), "snapshot") &&
				strings.Contains(strings.ToLower(errOvh.Message), "not") {
			return fmt.Errorf(
				"VPS snapshot option does not appear to be subscribed for this VPS "+
					"(API said %q). Subscribe to the 'snapshot' option on the VPS "+
					"before creating an ovh_vps_snapshot resource.",
				errOvh.Message,
			)
		}
	}
	if strings.Contains(strings.ToLower(msg), "option") &&
		strings.Contains(strings.ToLower(msg), "snapshot") {
		return fmt.Errorf(
			"VPS snapshot option does not appear to be subscribed for this VPS "+
				"(API said %q). Subscribe to the 'snapshot' option on the VPS "+
				"before creating an ovh_vps_snapshot resource.",
			msg,
		)
	}
	return err
}

func resourceVPSSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	description := d.Get("description").(string)

	body := &VPSCreateSnapshotOpts{Description: description}
	task := &VPSTask{}

	endpoint := fmt.Sprintf("/vps/%s/createSnapshot", url.PathEscape(serviceName))
	if err := config.OVHClient.Post(endpoint, body, task); err != nil {
		return diag.FromErr(surfaceOptionError(fmt.Errorf("calling Post %s: %w", endpoint, err)))
	}

	stateConf := vpsSnapshotTaskStateChangeConf(serviceName, task.Id, meta)
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("waiting for VPS snapshot create task %d to complete: %s", task.Id, err)
	}

	// After task completes, fetch the snapshot to populate id, creation_date, region.
	snap := &VPSSnapshot{}
	getEndpoint := fmt.Sprintf("/vps/%s/snapshot", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(getEndpoint, snap); err != nil {
		return diag.Errorf("calling Get %s after createSnapshot: %s", getEndpoint, err)
	}

	d.SetId(serviceName)

	// If a description was requested and the API didn't persist it via
	// createSnapshot, ensure it sticks via a Put.
	if description != "" && snap.Description != description {
		putEndpoint := fmt.Sprintf("/vps/%s/snapshot", url.PathEscape(serviceName))
		body := &VPSSnapshot{Description: description}
		if err := config.OVHClient.Put(putEndpoint, body, nil); err != nil {
			return diag.Errorf("calling Put %s: %s", putEndpoint, err)
		}
		// re-read
		if err := config.OVHClient.Get(getEndpoint, snap); err != nil {
			return diag.Errorf("calling Get %s after Put: %s", getEndpoint, err)
		}
	}

	d.Set("description", snap.Description)
	if !snap.CreationDate.IsZero() {
		d.Set("creation_date", snap.CreationDate.Format("2006-01-02T15:04:05Z07:00"))
	}
	d.Set("region", snap.Region)

	log.Printf("[DEBUG] Created VPS snapshot for %s", serviceName)
	return nil
}

func resourceVPSSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	if serviceName == "" {
		// Importing — id is the service_name.
		serviceName = d.Id()
		d.Set("service_name", serviceName)
	}

	snap := &VPSSnapshot{}
	endpoint := fmt.Sprintf("/vps/%s/snapshot", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(endpoint, snap); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("calling Get %s: %s", endpoint, err)
	}

	d.Set("description", snap.Description)
	if !snap.CreationDate.IsZero() {
		d.Set("creation_date", snap.CreationDate.Format("2006-01-02T15:04:05Z07:00"))
	}
	d.Set("region", snap.Region)
	return nil
}

func resourceVPSSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	body := &VPSSnapshot{
		Description: d.Get("description").(string),
	}
	endpoint := fmt.Sprintf("/vps/%s/snapshot", url.PathEscape(serviceName))
	if err := config.OVHClient.Put(endpoint, body, nil); err != nil {
		return diag.Errorf("calling Put %s: %s", endpoint, err)
	}
	return resourceVPSSnapshotRead(ctx, d, meta)
}

func resourceVPSSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	task := &VPSTask{}
	endpoint := fmt.Sprintf("/vps/%s/snapshot", url.PathEscape(serviceName))
	if err := config.OVHClient.Delete(endpoint, task); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return nil
		}
		return diag.Errorf("calling Delete %s: %s", endpoint, err)
	}

	if task.Id != 0 {
		stateConf := vpsSnapshotTaskStateChangeConf(serviceName, task.Id, meta)
		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return diag.Errorf("waiting for VPS snapshot delete task %d to complete: %s", task.Id, err)
		}
	}
	return nil
}
