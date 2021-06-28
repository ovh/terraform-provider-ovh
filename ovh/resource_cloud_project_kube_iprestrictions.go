package ovh

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectKubeIpRestrictions() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectKubeIpRestrictionsCreate,
		Read:   resourceCloudProjectKubeIpRestrictionsRead,
		Delete: resourceCloudProjectKubeIpRestrictionsDelete,
		Update: resourceCloudProjectKubeIpRestrictionsUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectKubeIpRestrictionsImportState,
		},

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
			"ips": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of ip restrictions for the cluster",
				Required:    true,
			},
		},
	}
}

func resourceCloudProjectKubeIpRestrictionsImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not service_name/kubeid formatted")
	}
	serviceName := splitId[0]
	kubeId := splitId[1]
	d.SetId(kubeId)
	d.Set("kube_id", kubeId)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectKubeIpRestrictionsCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceCloudProjectKubeIpRestrictionsUpdate(d, meta)
}

func resourceCloudProjectKubeIpRestrictionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/ipRestrictions", url.PathEscape(serviceName), url.PathEscape(kubeId))
	res := make(CloudProjectKubeIpRestrictionsResponse, 0)

	log.Printf("[DEBUG] Will read iprestrictions from cluster %s in project %s", kubeId, serviceName)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(kubeId)
	d.Set("ips", res)

	log.Printf("[DEBUG] Read iprestrictions: %+v", res)
	return nil
}

func resourceCloudProjectKubeIpRestrictionsUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/ipRestrictions", url.PathEscape(serviceName), url.PathEscape(kubeId))
	params := (&CloudProjectKubeIpRestrictionsCreateOrUpdateOpts{}).FromResource(d)
	res := make(CloudProjectKubeIpRestrictionsResponse, 0)

	log.Printf("[DEBUG] Will update iprestrictions: %+v", params)
	err := config.OVHClient.Put(endpoint, params, &res)
	if err != nil {
		return fmt.Errorf("calling Put %s with params %s:\n\t %q", endpoint, params, err)
	}

	d.SetId(kubeId)

	return resourceCloudProjectKubeIpRestrictionsRead(d, meta)
}

func resourceCloudProjectKubeIpRestrictionsDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	ips, _ := helpers.StringsFromSchema(d, "ips")

	for _, ip := range ips {
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/ipRestrictions/%s", url.PathEscape(serviceName), url.PathEscape(kubeId), url.PathEscape(ip))

		log.Printf("[DEBUG] Will delete iprestrictions ip %s from cluster %s in project %s", ip, kubeId, serviceName)
		err := config.OVHClient.Delete(endpoint, nil)
		if err != nil {
			return helpers.CheckDeleted(d, err, endpoint)
		}
	}

	return nil
}
