package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudProjectContainerRegistryIPRestrictionsManagement() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectContainerRegistryIPRestrictionsManagementRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"registry_id": {
				Type:        schema.TypeString,
				Description: "Registry ID",
				Required:    true,
			},
			"ip_restrictions": {
				Type:        schema.TypeList,
				Description: "List your IP restrictions applied on artifact manager component",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Set:  schema.HashString,
				},
			},
		},
	}
}

func dataSourceCloudProjectContainerRegistryIPRestrictionsManagementRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/ipRestrictions/management",
		url.PathEscape(serviceName),
		url.PathEscape(registryID),
	)
	ipRestrictions := []CloudProjectContainerRegistryIPRestriction{}

	log.Printf("[DEBUG] Will read Management IP Restrictions from registry %s and project: %s", registryID, serviceName)
	err := config.OVHClient.Get(endpoint, &ipRestrictions)
	if err != nil {
		return fmt.Errorf("calling get %s %w", endpoint, err)
	}

	mapIPRestrictions := make([]map[string]interface{}, len(ipRestrictions))
	for i, ipRestriction := range ipRestrictions {
		mapIPRestrictions[i] = ipRestriction.ToMap()
	}
	d.Set("ip_restrictions", mapIPRestrictions)
	d.SetId(serviceName + "/" + registryID)

	log.Printf("[DEBUG] Read Management IP Restrictions %+v", mapIPRestrictions)

	return nil
}
