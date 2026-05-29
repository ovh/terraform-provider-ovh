package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSVeeamRestorePoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSVeeamRestorePointRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS",
			},
			"restore_point_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the Veeam restore point",
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the restore point",
			},
		},
	}
}

func dataSourceVPSVeeamRestorePointRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	id := int64(d.Get("restore_point_id").(int))

	endpoint := fmt.Sprintf("/vps/%s/veeam/restorePoints/%d", url.PathEscape(serviceName), id)
	point := &VpsVeeamRestorePoint{}
	if err := config.OVHClient.Get(endpoint, point); err != nil {
		return fmt.Errorf("error calling GET %s: %s", endpoint, err)
	}

	d.SetId(fmt.Sprintf("%s/%d", serviceName, point.Id))
	d.Set("creation_time", point.CreationTime)
	return nil
}
