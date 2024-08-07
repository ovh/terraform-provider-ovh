package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHostingPrivateDatabaseWhitelist() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostingPrivateDatabaseWhitelistRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your private database",
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The whitelisted IP in your instance",
			},

			// Computed
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date of this whitelist",
			},
			"last_update": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last update date of this whitelist",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Custom name for your Whitelisted IP",
			},
			"service": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Authorize this IP to access service port",
			},
			"sftp": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Authorize this IP to access SFTP port",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Whitelist status",
			},
		},
	}
}

func dataSourceHostingPrivateDatabaseWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	ip := HostingPrivateDatabaseWhitelistefaultNetmask(d.Get("ip").(string))

	ds := &HostingPrivateDatabaseWhitelist{}
	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/hosting/privateDatabase/%s/whitelist/%s",
			url.PathEscape(serviceName),
			url.PathEscape(ip),
		),
		&ds,
	)

	if err != nil {
		return fmt.Errorf(
			"error calling hosting/privateDatabase/%s/whitelist/%s:\n\t %q",
			serviceName,
			ip,
			err,
		)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, ip))
	for k, v := range ds.DataSourceToMap() {
		d.Set(k, v)
	}

	return nil
}
