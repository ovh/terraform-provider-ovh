package ovh

import (
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceCloudProjectDatabases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectDatabasesRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"engine": {
				Type:        schema.TypeString,
				Description: "Name of the engine of the service",
				Required:    true,
			},

			//Computed
			"cluster_ids": {
				Type:        schema.TypeList,
				Description: "List of database clusters uuids",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceCloudProjectDatabasesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	engine := d.Get("engine").(string)

	serviceEndpoint := fmt.Sprintf("/cloud/project/%s/database/%s",
		url.PathEscape(serviceName),
		url.PathEscape(engine),
	)
	res := make([]string, 0)

	log.Printf("[DEBUG] Will list databases from project: %s", serviceName)
	if err := config.OVHClient.Get(serviceEndpoint, &res); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", serviceEndpoint, err)
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(res)

	d.SetId(hashcode.Strings(res))
	d.Set("cluster_ids", res)

	return nil
}
