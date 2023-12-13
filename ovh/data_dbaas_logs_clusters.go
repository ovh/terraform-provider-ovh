package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceDbaasLogsClusters() *schema.Resource {
	return &schema.Resource{
		Read: func(d *schema.ResourceData, meta interface{}) error {
			return dataSourceDbaasLogsClustersRead(d, meta)
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "The service name",
				Required:    true,
			},
			// Computed
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"uuids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "UUID of clusters",
				Computed:    true,
			},
		},
	}
}

func dataSourceDbaasLogsClustersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will read dbaas logs clusters %s", serviceName)

	d.SetId(serviceName)
	d.Set("urn", helpers.ServiceURN(config.Plate, "ldp", serviceName))

	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/cluster",
		url.PathEscape(serviceName),
	)

	res := []string{}
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	d.Set("uuids", res)

	return nil
}
