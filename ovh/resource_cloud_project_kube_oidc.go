package ovh

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloudProjectKubeOIDC() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectKubeOIDCCreate,
		Read:   resourceCloudProjectKubeOIDCRead,
		Delete: resourceCloudProjectKubeOIDCDelete,
		Update: resourceCloudProjectKubeOIDCUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectKubeOIDCImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(10 * time.Minute),
			Update:  schema.DefaultTimeout(10 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
			Read:    schema.DefaultTimeout(5 * time.Minute),
			Default: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"kube_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer_url": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceCloudProjectKubeOIDCImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	log.Printf("[DEBUG] Importing cloud project kube OIDC %s", givenId)

	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("import Id is not service_name/kubeid formatted")
	}
	serviceName := splitId[0]
	kubeID := splitId[1]
	d.SetId(serviceName + "/" + kubeID)
	d.Set("kube_id", kubeID)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectKubeOIDCCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	kubeID := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/openIdConnect", serviceName, kubeID)
	params := (&CloudProjectKubeOIDCCreateOpts{}).FromResource(d)
	res := &CloudProjectKubeOIDCResponse{}

	log.Printf("[DEBUG] Will create OIDC: %+v", params)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %s:\n\t %w", endpoint, params, err)
	}

	d.SetId(serviceName + "/" + kubeID)

	log.Printf("[DEBUG] Waiting for kube %s to be READY", kubeID)
	err = waitForCloudProjectKubeReady(config.OVHClient, serviceName, kubeID, []string{"REDEPLOYING"}, []string{"READY"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("timeout while waiting kube %s to be READY: %w", kubeID, err)
	}
	log.Printf("[DEBUG] kube %s is READY", kubeID)

	return resourceCloudProjectKubeOIDCRead(d, meta)
}

func resourceCloudProjectKubeOIDCRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	kubeID := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/openIdConnect", serviceName, kubeID)
	res := &CloudProjectKubeOIDCResponse{}

	log.Printf("[DEBUG] Will read oidc from kube %s and project: %s", kubeID, serviceName)
	err := config.OVHClient.Get(endpoint, res)
	if err != nil {
		return fmt.Errorf("calling get %s %w", endpoint, err)
	}
	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(serviceName + "/" + kubeID)
		}
	}

	log.Printf("[DEBUG] Read kube %+v", res)
	return nil
}

func resourceCloudProjectKubeOIDCUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	kubeID := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/openIdConnect", serviceName, kubeID)
	params := (&CloudProjectKubeOIDCUpdateOpts{}).FromResource(d)
	res := &CloudProjectKubeOIDCResponse{}

	log.Printf("[DEBUG] Will update OIDC: %+v", params)
	err := config.OVHClient.Put(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Put %s with params %s:\n\t %w", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for kube %s to be READY", kubeID)
	err = waitForCloudProjectKubeReady(config.OVHClient, serviceName, kubeID, []string{"REDEPLOYING"}, []string{"READY"}, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("timeout while waiting kube %s to be READY: %w", kubeID, err)
	}
	log.Printf("[DEBUG] kube %s is READY", kubeID)

	return nil
}

func resourceCloudProjectKubeOIDCDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	kubeID := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/openIdConnect", serviceName, kubeID)

	log.Printf("[DEBUG] Will delete OIDC")
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return fmt.Errorf("calling delete %s %w", endpoint, err)
	}

	log.Printf("[DEBUG] Waiting for kube %s to be READY", kubeID)
	err = waitForCloudProjectKubeReady(config.OVHClient, serviceName, kubeID, []string{"REDEPLOYING"}, []string{"READY"}, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("timeout while waiting kube %s to be READY: %w", kubeID, err)
	}
	log.Printf("[DEBUG] kube %s is READY", kubeID)

	return nil
}
