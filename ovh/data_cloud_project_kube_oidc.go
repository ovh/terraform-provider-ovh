package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloudProjectKubeOIDC() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectKubeOIDCRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"kube_id": {
				Type:        schema.TypeString,
				Description: "Kube ID",
				Required:    true,
				ForceNew:    true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"issuer_url": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"oidc_username_claim": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"oidc_username_prefix": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"oidc_groups_claim": {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"oidc_groups_prefix": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			"oidc_required_claim": {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"oidc_signing_algs": {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"oidc_ca_content": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
		},
	}
}

func dataSourceCloudProjectKubeOIDCRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/openIdConnect", serviceName, kubeId)
	res := &CloudProjectKubeOIDCResponse{}

	log.Printf("[DEBUG] Will read OIDC from kube %s and project: %s", kubeId, serviceName)
	err := config.OVHClient.Get(endpoint, res)
	if err != nil {
		return fmt.Errorf("calling get %s %w", endpoint, err)
	}
	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		}
	}
	d.SetId(kubeId + "-" + res.ClientID + "-" + res.IssuerUrl)

	log.Printf("[DEBUG] Read OIDC %+v", res)
	return nil
}
