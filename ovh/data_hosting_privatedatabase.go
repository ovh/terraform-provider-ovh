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
				Description: "Number of cpu on your private database",
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
				Description: "Private database ftp hostname",
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
				Description: "Private database ftp port",
			},
			"quota_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Space allowed (in MB) on your private database",
			},
			"quota_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Sapce used (in MB) on your private database",
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

	d.SetId(ds.ServiceName)
	d.Set("cpu", ds.Cpu)
	d.Set("datacenter", ds.Datacenter)
	d.Set("display_name", ds.DisplayName)
	d.Set("hostname", ds.Hostname)
	d.Set("hostname_ftp", ds.HostnameFtp)
	d.Set("infrastructure", ds.Infrastructure)
	d.Set("offer", ds.Offer)
	d.Set("port", ds.Port)
	d.Set("port_ftp", ds.PortFtp)
	d.Set("quota_size", ds.QuotaSize.Value)
	d.Set("quota_used", ds.QuotaUsed.Value)
	d.Set("ram", ds.Ram.Value)
	d.Set("server", ds.Server)
	d.Set("service_name", ds.ServiceName)
	d.Set("state", ds.State)
	d.Set("type", ds.Type)
	d.Set("version", ds.Version)
	d.Set("version_label", ds.VersionLabel)
	d.Set("version_number", ds.VersionNumber)

	return nil
}
