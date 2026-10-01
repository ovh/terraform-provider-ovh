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
			kubeServiceNameKey: {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			kubeKubeIdKey: {
				Type:        schema.TypeString,
				Description: "Kube ID",
				Required:    true,
				ForceNew:    true,
			},
			kubeOidcClientIdKey: {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			kubeOidcIssuerUrlKey: {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			kubeOidcUsernameClaimKey: {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			kubeOidcUsernamePrefixKey: {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			kubeOidcGroupsClaimKey: {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			kubeOidcGroupsPrefixKey: {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			kubeOidcRequiredClaimKey: {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			kubeOidcSigningAlgsKey: {
				Type:     schema.TypeList,
				Required: false,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			kubeOidcCaContentKey: {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
		},
	}
}

func dataSourceCloudProjectKubeOIDCRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeId := d.Get(kubeKubeIdKey).(string)

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
