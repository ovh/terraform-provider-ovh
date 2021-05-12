package ovh

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectKube() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectKubeCreate,
		Read:   resourceCloudProjectKubeRead,
		Delete: resourceCloudProjectKubeDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectKubeNodePoolImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Computed
			"control_plane_is_up_to_date": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_up_to_date": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"next_upgrade_versions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"nodes_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kubeconfig": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceCloudProjectKubeImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 3 {
		return nil, fmt.Errorf("Import Id is not service_name/kubeid formatted")
	}
	serviceName := splitId[0]
	id := splitId[1]
	d.SetId(id)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectKubeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube", serviceName)
	params := (&CloudProjectKubeCreateOpts{}).FromResource(d)
	res := &CloudProjectKubeResponse{}

	log.Printf("[DEBUG] Will create kube: %+v", params)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %s:\n\t %q", endpoint, params, err)
	}

	// This is a fix for a weird bug where the kube is not immediately available on API
	log.Printf("[DEBUG] Waiting for kube %s to be available", res.Id)
	endpoint = fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, res.Id)
	err = helpers.WaitAvailable(config.OVHClient, endpoint, 2*time.Minute)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Waiting for kube %s to be READY", res.Id)
	err = waitForCloudProjectKubeReady(config.OVHClient, serviceName, res.Id)
	if err != nil {
		return fmt.Errorf("timeout while waiting kube %s to be READY: %v", res.Id, err)
	}
	log.Printf("[DEBUG] kube %s is READY", res.Id)

	d.SetId(res.Id)

	return resourceCloudProjectKubeRead(d, meta)
}

func resourceCloudProjectKubeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, d.Id())
	res := &CloudProjectKubeResponse{}

	log.Printf("[DEBUG] Will read kube %s from project: %s", d.Id(), serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}
	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	if d.IsNewResource() {
		kubeconfigRaw := CloudProjectKubeKubeConfigResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/kubeconfig", serviceName, res.Id)
		err := config.OVHClient.Post(endpoint, nil, &kubeconfigRaw)

		if err != nil {
			return err
		}
		d.Set("kubeconfig", kubeconfigRaw.Content)
	}

	log.Printf("[DEBUG] Read kube %+v", res)
	return nil
}

func resourceCloudProjectKubeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, d.Id())

	log.Printf("[DEBUG] Will delete kube %s from project: %s", d.Id(), serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	log.Printf("[DEBUG] Waiting for kube %s to be DELETED", d.Id())
	err = waitForCloudProjectKubeDeleted(config.OVHClient, serviceName, d.Id())
	if err != nil {
		return fmt.Errorf("timeout while waiting kube %s to be DELETED: %v", d.Id(), err)
	}
	log.Printf("[DEBUG] kube %s is DELETED", d.Id())

	d.SetId("")

	return nil
}

func cloudProjectKubeExists(serviceName, id string, client *ovh.Client) error {
	res := &CloudProjectKubeResponse{}

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, id)
	return client.Get(endpoint, res)
}

func waitForCloudProjectKubeReady(client *ovh.Client, serviceName, kubeId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INSTALLING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectKubeResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, kubeId)
			err := client.Get(endpoint, res)
			if err != nil {
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectKubeDeleted(client *ovh.Client, serviceName, kubeId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectKubeResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, kubeId)
			err := client.Get(endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				} else {
					return res, "", err
				}
			}

			return res, res.Status, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}
