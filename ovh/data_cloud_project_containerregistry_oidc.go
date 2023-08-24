package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudProjectContainerRegistryOIDC() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectContainerRegistryOIDCRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"registry_id": {
				Type:        schema.TypeString,
				Description: "Registry ID",
				Required:    true,
				ForceNew:    true,
			},
			"oidc_name": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"oidc_endpoint": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"oidc_client_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"oidc_scope": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"oidc_groups_claim": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"oidc_admin_group": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"oidc_verify_cert": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
			},
			"oidc_auto_onboard": {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
			},
			"oidc_user_claim": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
		},
	}
}

func dataSourceCloudProjectContainerRegistryOIDCRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	registryID := d.Get("registry_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/containerRegistry/%s/openIdConnect", serviceName, registryID)
	res := &CloudProjectContainerRegistryOIDCResponse{}

	log.Printf("[DEBUG] Will read OIDC from registry %s and project: %s", registryID, serviceName)
	err := config.OVHClient.Get(endpoint, res)
	if err != nil {
		return fmt.Errorf("calling get %s %w", endpoint, err)
	}
	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		}
	}
	d.SetId(registryID + "-" + res.ClientID + "-" + res.Endpoint)

	log.Printf("[DEBUG] Read OIDC %+v", res)
	return nil
}
