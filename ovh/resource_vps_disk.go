package ovh

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

// resourceVPSDisk manages monitoring settings on an existing VPS disk.
//
// IMPORTANT: This resource does NOT provision an additional disk. Disks are
// provisioned through the cart options on the parent ovh_vps resource via
// plan_option entries (planCode "option-additional-disk-2025-XXXXg"). This
// resource only manages the writable fields (monitoring, lowFreeSpaceThreshold)
// of an already-existing disk. As a consequence, the Delete operation does
// not remove the disk; it restores default monitoring settings.
func resourceVPSDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPSDiskCreate,
		Read:   resourceVPSDiskRead,
		Update: resourceVPSDiskUpdate,
		Delete: resourceVPSDiskDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVPSDiskImportState,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of the VPS (e.g. vps-XXXXXX.vps.ovh.net).",
			},
			"disk_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Numeric identifier of the disk attached to the VPS.",
			},
			"monitoring": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Enable or disable OVHcloud monitoring on the disk.",
			},
			"low_free_space_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Threshold (in MiB) below which a low-free-space alert is raised.",
			},
			// Computed
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disk type (primary or additional).",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Disk state (connected, disconnected, pending).",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Disk size in GiB.",
			},
			"bandwidth_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Disk bandwidth limit.",
			},
		},
	}
}

func vpsDiskID(serviceName string, diskID int64) string {
	return fmt.Sprintf("%s|%d", serviceName, diskID)
}

func resourceVPSDiskImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	given := d.Id()
	parts := strings.SplitN(given, "|", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("import id must be formatted as service_name|disk_id")
	}
	serviceName := parts[0]
	diskID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid disk_id in import id %q: %s", given, err)
	}
	d.SetId(vpsDiskID(serviceName, diskID))
	d.Set("service_name", serviceName)
	d.Set("disk_id", diskID)
	return []*schema.ResourceData{d}, nil
}

func vpsDiskEndpoint(serviceName string, diskID int64) string {
	return fmt.Sprintf("/vps/%s/disks/%d", url.PathEscape(serviceName), diskID)
}

func resourceVPSDiskFetch(config *Config, serviceName string, diskID int64) (*VPSDisk, error) {
	endpoint := vpsDiskEndpoint(serviceName, diskID)
	resp := &VPSDisk{}
	if err := config.OVHClient.Get(endpoint, resp); err != nil {
		return nil, fmt.Errorf("calling GET %s:\n\t %s", endpoint, err.Error())
	}
	return resp, nil
}

func resourceVPSDiskPut(config *Config, serviceName string, diskID int64, opts *VPSDiskUpdateOpts) error {
	endpoint := vpsDiskEndpoint(serviceName, diskID)
	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("calling PUT %s:\n\t %s", endpoint, err.Error())
	}
	return nil
}

func resourceVPSDiskBuildOpts(d *schema.ResourceData, current *VPSDisk) *VPSDiskUpdateOpts {
	opts := &VPSDiskUpdateOpts{
		Monitoring: current.Monitoring,
	}
	if v, ok := d.GetOkExists("monitoring"); ok {
		opts.Monitoring = v.(bool)
	}

	if v, ok := d.GetOkExists("low_free_space_threshold"); ok {
		t := int64(v.(int))
		opts.LowFreeSpaceThreshold = &t
	} else if current.LowFreeSpaceThreshold != nil {
		opts.LowFreeSpaceThreshold = current.LowFreeSpaceThreshold
	}

	return opts
}

func resourceVPSDiskCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	diskID := int64(d.Get("disk_id").(int))

	current, err := resourceVPSDiskFetch(config, serviceName, diskID)
	if err != nil {
		return err
	}

	opts := resourceVPSDiskBuildOpts(d, current)
	if err := resourceVPSDiskPut(config, serviceName, diskID, opts); err != nil {
		return err
	}

	d.SetId(vpsDiskID(serviceName, diskID))
	return resourceVPSDiskRead(d, meta)
}

func resourceVPSDiskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	diskID := int64(d.Get("disk_id").(int))

	endpoint := vpsDiskEndpoint(serviceName, diskID)
	resp := &VPSDisk{}
	if err := config.OVHClient.Get(endpoint, resp); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("service_name", resp.ServiceName)
	d.Set("disk_id", resp.ID)
	d.Set("type", resp.Type)
	d.Set("state", resp.State)
	d.Set("size", resp.Size)
	d.Set("bandwidth_limit", resp.BandwidthLimit)
	d.Set("monitoring", resp.Monitoring)
	if resp.LowFreeSpaceThreshold != nil {
		d.Set("low_free_space_threshold", *resp.LowFreeSpaceThreshold)
	} else {
		d.Set("low_free_space_threshold", 0)
	}

	return nil
}

func resourceVPSDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	diskID := int64(d.Get("disk_id").(int))

	current, err := resourceVPSDiskFetch(config, serviceName, diskID)
	if err != nil {
		return err
	}
	opts := resourceVPSDiskBuildOpts(d, current)
	if err := resourceVPSDiskPut(config, serviceName, diskID, opts); err != nil {
		return err
	}

	return resourceVPSDiskRead(d, meta)
}

// resourceVPSDiskDelete does NOT actually delete the disk: the underlying
// OVH API does not expose a DELETE for /vps/{serviceName}/disks/{id}. Disks
// are removed by terminating their associated cart option through the parent
// ovh_vps resource. This Delete restores the writable fields to their
// defaults (monitoring=true, lowFreeSpaceThreshold cleared) and drops the
// resource from Terraform state.
func resourceVPSDiskDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	diskID := int64(d.Get("disk_id").(int))

	opts := &VPSDiskUpdateOpts{
		Monitoring:            true,
		LowFreeSpaceThreshold: nil,
	}
	if err := resourceVPSDiskPut(config, serviceName, diskID, opts); err != nil {
		// If the disk no longer exists, treat as deleted.
		return helpers.CheckDeleted(d, err, vpsDiskEndpoint(serviceName, diskID))
	}

	d.SetId("")
	return nil
}
