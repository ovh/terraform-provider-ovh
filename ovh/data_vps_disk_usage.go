package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSDiskUsage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSDiskUsageRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disk_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "used",
				Description: "Type of usage to query: max or used.",
				ValidateFunc: func(v interface{}, k string) (warns []string, errs []error) {
					switch v.(string) {
					case "max", "used":
					default:
						errs = append(errs, fmt.Errorf("%s must be one of max,used", k))
					}
					return
				},
			},
			// Computed
			"unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func dataSourceVPSDiskUsageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	diskID := int64(d.Get("disk_id").(int))
	usageType := d.Get("type").(string)

	endpoint := fmt.Sprintf("/vps/%s/disks/%d/use?type=%s",
		url.PathEscape(serviceName), diskID, url.QueryEscape(usageType))

	resp := &VPSDiskUsage{}
	if err := config.OVHClient.Get(endpoint, resp); err != nil {
		return fmt.Errorf("calling GET %s:\n\t %s", endpoint, err.Error())
	}

	d.SetId(fmt.Sprintf("%s|%d|%s", serviceName, diskID, usageType))
	d.Set("unit", resp.Unit)
	d.Set("value", resp.Value)
	return nil
}
