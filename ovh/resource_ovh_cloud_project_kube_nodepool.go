package ovh

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectKubeNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectKubeNodePoolCreate,
		Read:   resourceCloudProjectKubeNodePoolRead,
		Delete: resourceCloudProjectKubeNodePoolDelete,
		Update: resourceCloudProjectKubeNodePoolUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceCloudProjectKubeNodePoolRead(d, meta)
				return []*schema.ResourceData{d}, err
			},
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"flavor_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"desired_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"monthly_billed": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  "false",
			},
			"anti_affinity": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  "false",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudProjectKubeNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool", serviceName, kubeId)
	params := &CloudProjectKubeNodePoolCreateOpts{
		Name:          d.Get("name").(string),
		FlavorName:    d.Get("flavor_name").(string),
		DesiredNodes:  d.Get("desired_nodes").(int),
		MaxNodes:      d.Get("max_nodes").(int),
		MinNodes:      d.Get("min_nodes").(int),
		MonthlyBilled: d.Get("monthly_billed").(bool),
		AntiAffinity:  d.Get("anti_affinity").(bool),
	}
	res := &CloudProjectKubeNodePoolResponse{}

	log.Printf("[DEBUG] Will create nodepool: %+v", params)
	err := config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %s:\n\t %q", endpoint, params, err)
	}

	// This is a fix for a weird bug where the nodepool is not immediately available on API
	log.Printf("[DEBUG] Waiting for nodepool %s to be available", res.Id)
	endpoint = fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, res.Id)
	err = helpers.WaitAvailable(config.OVHClient, endpoint, 2*time.Minute)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Waiting for nodepool %s to be READY", res.Id)
	err = waitForCloudProjectKubeNodePoolReady(config.OVHClient, serviceName, kubeId, res.Id)
	if err != nil {
		return fmt.Errorf("timeout while waiting nodepool %s to be READY: %v", res.Id, err)
	}
	log.Printf("[DEBUG] nodepool %s is READY", res.Id)

	d.SetId(res.Id)

	return resourceCloudProjectKubeNodePoolRead(d, meta)
}

func resourceCloudProjectKubeNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, d.Id())
	res := &CloudProjectKubeNodePoolResponse{}

	log.Printf("[DEBUG] Will read nodepool %s from cluster %s in project %s", d.Id(), kubeId, serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.SetId(res.Id)
	d.Set("name", res.Name)
	d.Set("flavor_name", res.Flavor)
	d.Set("desired_nodes", res.DesiredNodes)
	d.Set("max_nodes", res.MaxNodes)
	d.Set("min_nodes", res.MinNodes)
	d.Set("monthly_billed", res.MonthlyBilled)
	d.Set("anti_affinity", res.AntiAffinity)
	d.Set("status", res.Status)

	log.Printf("[DEBUG] Read nodepool: %+v", res)
	return nil
}

func resourceCloudProjectKubeNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, d.Id())
	params := &CloudProjectKubeNodePoolUpdateOpts{
		DesiredNodes: d.Get("desired_nodes").(int),
		MaxNodes:     d.Get("max_nodes").(int),
		MinNodes:     d.Get("min_nodes").(int),
	}

	log.Printf("[DEBUG] Will update nodepool: %+v", params)
	err := config.OVHClient.Put(endpoint, params, nil)
	if err != nil {
		return fmt.Errorf("calling Put %s with params %s:\n\t %q", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for nodepool %s to be READY", d.Id())
	err = waitForCloudProjectKubeNodePoolReady(config.OVHClient, serviceName, kubeId, d.Id())
	if err != nil {
		return fmt.Errorf("timeout while waiting nodepool %s to be READY: %v", d.Id(), err)
	}
	log.Printf("[DEBUG] nodepool %s is READY", d.Id())

	return resourceCloudProjectKubeNodePoolRead(d, meta)
}

func resourceCloudProjectKubeNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, d.Id())

	log.Printf("[DEBUG] Will delete nodepool %s from cluster %s in project %s", d.Id(), kubeId, serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	log.Printf("[DEBUG] Waiting for nodepool %s to be DELETED", d.Id())
	err = waitForCloudProjectKubeNodePoolDeleted(config.OVHClient, serviceName, kubeId, d.Id())
	if err != nil {
		return fmt.Errorf("timeout while waiting nodepool %s to be DELETED: %v", d.Id(), err)
	}
	log.Printf("[DEBUG] nodepool %s is DELETED", d.Id())

	d.SetId("")

	return nil
}

func cloudProjectKubeNodePoolExists(serviceName, kubeId, id string, client *ovh.Client) error {
	res := &CloudProjectKubeResponse{}

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, id)
	return client.Get(endpoint, res)
}

func waitForCloudProjectKubeNodePoolReady(client *ovh.Client, serviceName, kubeId, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INSTALLING", "UPDATING", "REDEPLOYING", "RESIZING"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectKubeResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, id)
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

func waitForCloudProjectKubeNodePoolDeleted(client *ovh.Client, serviceName, kubeId, id string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectKubeResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, id)
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
