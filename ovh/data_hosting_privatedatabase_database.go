package ovh

import (
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceHostingPrivateDatabaseDatabase() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHostingPrivateDatabaseDatabaseRead,
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

			// Computed
			"backup_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time of the next backup (every day)",
			},
			"creation_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation date of the database",
			},
			"quota_used": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Space used by the database (in MB)",
			},
			"users": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Users granted to this database",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": {
							Type:        schema.TypeString,
							Description: "User's name granted on this database",
							Computed:    true,
						},
						"grant_type": {
							Type:        schema.TypeString,
							Description: "User's rights on this database",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceHostingPrivateDatabaseDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	databaseName := d.Get("database_name").(string)

	ds := &DataSourceHostingPrivateDatabaseDatabase{}

	err := config.OVHClient.Get(fmt.Sprintf("/hosting/privateDatabase/%s/database/%s", url.PathEscape(serviceName), url.PathEscape(databaseName)), &ds)

	if err != nil {
		return fmt.Errorf("error calling hosting/privateDatabase/%s/database/%s:\n\t %q", serviceName, databaseName, err)
	}

	d.SetId(serviceName)
	for k, v := range ds.ToMap() {
		d.Set(k, v)
	}

	return nil
}
