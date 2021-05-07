package ovh

import (
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceCloudProjectContainerRegistryUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectContainerRegistryUsersRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"registry_id": {
				Type:        schema.TypeString,
				Description: "RegistryID",
				Required:    true,
			},

			"result": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:        schema.TypeString,
							Description: "User email",
							Computed:    true,
						},
						"id": {
							Type:        schema.TypeString,
							Description: "User ID",
							Computed:    true,
						},
						"user": {
							Type:        schema.TypeString,
							Description: "User name",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudProjectContainerRegistryUsersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	regId := d.Get("registry_id").(string)

	log.Printf("[DEBUG] Will read cloud project registry %s for project: %s", regId, serviceName)

	users := []CloudProjectContainerRegistryUser{}

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/users",
		url.PathEscape(serviceName),
		url.PathEscape(regId),
	)

	err := config.OVHClient.Get(endpoint, &users)
	if err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	mapusers := make([]map[string]interface{}, len(users))
	ids := make([]string, len(users))

	for i, user := range users {
		mapusers[i] = user.ToMapWithKeys([]string{
			"email",
			"id",
			"user",
		})
		ids = append(ids, user.Id)
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(ids)

	d.SetId(hashcode.Strings(ids))
	d.Set("result", mapusers)

	return nil
}
