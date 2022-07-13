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
		Create: resourceCloudProjectKubeIpRestrictionsCreateOrUpdate,
		Update: resourceCloudProjectKubeIpRestrictionsCreateOrUpdate,
		Delete: resourceCloudProjectKubeIpRestrictionsDelete,
		Read:   resourceCloudProjectKubeIpRestrictionsRead,

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

func resourceCloudProjectKubeIpRestrictionsCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	params := (&CloudProjectKubeIpRestrictionsCreateOrUpdateOpts{}).FromResource(d)

	err := resourceCloudProjectKubeIpRestrictionsUpdate(config, serviceName, kubeId, params)
	if err != nil {
		return err
	}

	output := resourceCloudProjectKubeIpRestrictionsRead(d, meta)
	d.SetId(kubeId)

	return output
}

func resourceCloudProjectKubeIpRestrictionsDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	return resourceCloudProjectKubeIpRestrictionsUpdate(config, serviceName, kubeId, &CloudProjectKubeIpRestrictionsCreateOrUpdateOpts{
		Ips: []string{},
	})
}

func resourceCloudProjectKubeIpRestrictionsUpdate(config *Config, serviceName string, kubeId string, params *CloudProjectKubeIpRestrictionsCreateOrUpdateOpts) error {
	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/ipRestrictions", url.PathEscape(serviceName), url.PathEscape(kubeId))
	res := make(CloudProjectKubeIpRestrictionsResponse, 0)

	log.Printf("[DEBUG] Will update iprestrictions: %+v", params)
	err := config.OVHClient.Put(endpoint, params, &res)
	if err != nil {
		return fmt.Errorf("calling Put %s with params %s:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for kube %s to be READY", kubeId)
	err = waitForCloudProjectKubeReady(config.OVHClient, serviceName, kubeId, []string{"REDEPLOYING", "RESETTING"}, []string{"READY"})
	if err != nil {
		return fmt.Errorf("timeout while waiting kube %s to be READY: %v", kubeId, err)
	}
	log.Printf("[DEBUG] kube %s is READY", kubeId)

	return nil
}
