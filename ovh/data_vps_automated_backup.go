package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVPSAutomatedBackup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVPSAutomatedBackupRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your VPS",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Automated backup state (enabled/disabled)",
			},
			"schedule": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Scheduled backup time (HH:MM:SS)",
			},
			"rotation": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of retained automated backups",
			},
			"service_resource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource name of the automated backup service",
			},
		},
	}
}

func dataSourceVPSAutomatedBackupRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	ab := &VPSAutomatedBackup{}
	endpoint := fmt.Sprintf("/vps/%s/automatedBackup", url.PathEscape(serviceName))
	if err := config.OVHClient.Get(endpoint, ab); err != nil {
		return fmt.Errorf("error calling GET %s: %w", endpoint, err)
	}

	d.SetId(serviceName)
	d.Set("state", ab.State)
	d.Set("schedule", ab.Schedule)
	d.Set("rotation", ab.Rotation)
	d.Set("service_resource_name", ab.ServiceResourceName)
	return nil
}
