package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHostingPrivateDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostingPrivateDatabaseUserRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The internal name of your private database",
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
			"databases": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Databases granted for this user",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name": {
							Type:        schema.TypeString,
							Description: "Database's name linked to this user",
							Computed:    true,
						},
						"grant_type": {
							Type:        schema.TypeString,
							Description: "Grant of this user for this database",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceHostingPrivateDatabaseUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	userName := d.Get("user_name").(string)

	ds := &DataSourceHostingPrivateDatabaseUser{}

	err := config.OVHClient.Get(
		fmt.Sprintf(
			"/hosting/privateDatabase/%s/user/%s",
			url.PathEscape(serviceName),
			url.PathEscape(userName),
		),
		&ds,
	)

	if err != nil {
		return fmt.Errorf(
			"Error calling hosting/privateDatabase/%s/user/%s:\n\t %q",
			serviceName,
			userName,
			err,
		)
	}

	d.SetId(fmt.Sprintf("%s/%s", serviceName, userName))
	for k, v := range ds.ToMap() {
		d.Set(k, v)
	}

	return nil
}
