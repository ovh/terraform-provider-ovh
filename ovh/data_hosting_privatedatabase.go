package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHostingPrivateDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostingPrivateDatabaseRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// Computed
			"cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of CPU on your private database",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Datacenter where this private database is located",
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name displayed in customer panel for your private database",
			},
			"hostname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private database hostname",
			},
			"hostname_ftp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private database FTP hostname",
			},
			"infrastructure": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Infrastructure where service was stored",
			},
			"offer": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the private database offer",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Private database service port",
			},
			"port_ftp": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Private database FTP port",
			},
			"quota_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Space allowed (in MB) on your private database",
			},
			"quota_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Space used (in MB) on your private database",
			},
			"ram": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Amount of ram (in MB) on your private database",
			},
			"server": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private database server name",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private database state",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private database type",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private database available versions",
			},
			"version_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private database version label",
			},
			"version_number": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Private database version number",
			},
		},
	}
}

func dataSourceHostingPrivateDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	ds := &HostingPrivateDatabase{}
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/hosting/privateDatabase/%s",
			url.PathEscape(serviceName),
		),
		&ds,
	)

	if err != nil {
		return fmt.Errorf(
			"Error calling hosting/privateDatabase/%s:\n\t %q",
			serviceName,
			err,
		)
	}

	for k, v := range ds.ToMap() {
		if k != "service_name" {
			d.Set(k, v)
		}
	}
	d.SetId(ds.ServiceName)

	return nil
}
