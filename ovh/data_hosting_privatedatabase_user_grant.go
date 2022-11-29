package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHostingPrivateDatabaseUserGrant() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostingPrivateDatabaseUserGrantRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your private database",
			},
			"database_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database name",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "User name used to connect to your databases",
			},

			// Computed
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date of the user",
			},
			"grant": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Grant name",
			},
		},
	}
}

func dataSourceHostingPrivateDatabaseUserGrantRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userName := d.Get("user_name").(string)
	databaseName := d.Get("database_name").(string)

	ds := &DataSourceHostingPrivateDatabaseUserGrant{}

	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/hosting/privateDatabase/%s/user/%s/grant/%s",
			url.PathEscape(serviceName),
			url.PathEscape(userName),
			url.PathEscape(databaseName),
		),
		&ds,
	)

	if err != nil {
		return fmt.Errorf(
			"Error calling hosting/privateDatabase/%s/user/%s/grant/%s:\n\t %q",
			serviceName,
			userName,
			databaseName,
			err,
		)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s/%s", serviceName, userName, databaseName, ds.Grant))
	for k, v := range ds.ToMap() {
		d.Set(k, v)
	}

	return nil
}
