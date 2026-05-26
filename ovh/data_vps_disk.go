package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSDisk() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSDiskRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disk_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			// Computed
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bandwidth_limit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"monitoring": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"low_free_space_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVPSDiskRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	diskID := int64(d.Get("disk_id").(int))

	endpoint := fmt.Sprintf("/vps/%s/disks/%d", url.PathEscape(serviceName), diskID)
	resp := &VPSDisk{}
	if err := config.OVHClient.Get(endpoint, resp); err != nil {
		return fmt.Errorf("calling GET %s:\n\t %s", endpoint, err.Error())
	}

	d.SetId(vpsDiskID(serviceName, diskID))
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
