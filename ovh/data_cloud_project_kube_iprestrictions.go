package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func dataSourceCloudProjectKubeIPRestrictions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectKubeIpRestrictionsRead,
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
			kubeIpRestrictionsIpsKey: {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of IP restrictions for the cluster",
				Computed:    true,
			},
		},
	}
}

func dataSourceCloudProjectKubeIpRestrictionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeId := d.Get(kubeKubeIdKey).(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/ipRestrictions", url.PathEscape(serviceName), url.PathEscape(kubeId))
	var res CloudProjectKubeIpRestrictionsResponse

	log.Printf("[DEBUG] Will read iprestrictions from cluster %s in project %s", kubeId, serviceName)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(kubeId)
	d.Set(kubeIpRestrictionsIpsKey, res)

	log.Printf("[DEBUG] Read iprestrictions: %+v", res)
	return nil
}
