package ovh

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ovh/go-ovh/ovh"
)

func resourceCloudKubernetesNodePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudKubernetesNodePoolCreate,
		Read:          resourceCloudKubernetesNodePoolRead,
		DeleteContext: resourceCloudKubernetesNodePoolDelete,
		UpdateContext: resourceCloudKubernetesNodePoolUpdate,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceCloudKubernetesNodePoolRead(d, meta)
				return []*schema.ResourceData{d}, err
			},
		},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_PROJECT_ID", nil),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Optional: false,
				ForceNew: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"desired_nodes": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"max_nodes": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"min_nodes": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"monthly_billed": {
				Type:     schema.TypeBool,
				Required: false,
				Default:  false,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudKubernetesNodePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	clusterId := d.Get("cluster_id").(string)
	poolId := d.Id()

	params := &CloudKubernetesNodePoolUpdateRequest{
		DesiredNodes: d.Get("desired_nodes").(int),
		MaxNodes:     d.Get("max_nodes").(int),
		MinNodes:     d.Get("min_nodes").(int),
	}

	d.Partial(true)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", projectId, clusterId, poolId)

	r := &CloudKubernetesNodePoolResponse{}

	err := config.OVHClient.Put(endpoint, params, r)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Failed to update node pool",
			Detail:        err.Error(),
			AttributePath: nil,
		})
		return
	}

	log.Printf("[DEBUG] Waiting for Node %s:", r)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INSTALLING", "UPDATING", "REDEPLOYING", "RESIZING"},
		Target:     []string{"READY"},
		Refresh:    waitForCloudKubernetesNodePoolActive(config.OVHClient, projectId, clusterId, r.Id),
		Timeout:    45 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		{
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Failed to get node pool ready",
				Detail:        err.Error(),
				AttributePath: nil,
			})
			return
		}
	}

	err = resourceCloudKubernetesNodePoolRead(d, meta)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Failed to read node pool",
			Detail:        err.Error(),
			AttributePath: nil,
		})
		return
	}

	d.Partial(false)

	return nil

}

func resourceCloudKubernetesNodePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	clusterId := d.Get("cluster_id").(string)
	params := &CloudKubernetesNodePoolCreationRequest{
		Name:          d.Get("name").(string),
		FlavorName:    d.Get("flavor").(string),
		DesiredNodes:  d.Get("desired_nodes").(int),
		MaxNodes:      d.Get("max_nodes").(int),
		MinNodes:      d.Get("min_nodes").(int),
		MonthlyBilled: d.Get("monthly_billed").(bool),
	}

	r := &CloudKubernetesNodePoolResponse{}

	log.Printf("[DEBUG] Will create public cloud kubernetes cluster: %s", params)

	d.Partial(true)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool", projectId, clusterId)

	err := config.OVHClient.Post(endpoint, params, r)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Failed to create node pool",
			Detail:        err.Error(),
			AttributePath: nil,
		})
		return
	}

	d.SetId(r.Id)

	// This is a fix for a weird bug where the cluster is not immediately available on API
	bugFixWait := &resource.StateChangeConf{
		Pending:    []string{"NOT_FOUND"},
		Target:     []string{"FOUND"},
		Refresh:    waitForCloudKubernetesNodePoolToBeReal(config.OVHClient, projectId, clusterId, r.Id),
		Timeout:    2 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	log.Printf("[DEBUG] Waiting for Node %s:", r)
	_, err = bugFixWait.WaitForStateContext(ctx)

	log.Printf("[DEBUG] Waiting for Node %s:", r)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INSTALLING", "UPDATING", "REDEPLOYING", "RESIZING"},
		Target:     []string{"READY"},
		Refresh:    waitForCloudKubernetesNodePoolActive(config.OVHClient, projectId, clusterId, r.Id),
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		{
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Failed to get node pool ready",
				Detail:        err.Error(),
				AttributePath: nil,
			})
			return
		}
	}
	log.Printf("[DEBUG] Created User %s", r)

	d.SetId(r.Id)
	err = resourceCloudKubernetesNodePoolRead(d, meta)

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       "Failed to read node pool",
			Detail:        err.Error(),
			AttributePath: nil,
		})
		return
	}

	d.Partial(false)

	return nil
}

func resourceCloudKubernetesNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	clusterId := d.Get("cluster_id").(string)

	d.Partial(true)
	r := &CloudKubernetesNodePoolResponse{}

	log.Printf("[DEBUG] Will read public cloud kubernetes cluster %s from project: %s", d.Id(), projectId)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", projectId, clusterId, d.Id())

	err := config.OVHClient.Get(endpoint, r)
	if err != nil {
		return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
	}

	err = readCloudKubernetesNodePool(projectId, config, d, r)
	if err != nil {
		return fmt.Errorf("error while reading cluster data %s:\n\t %q", d.Id(), err)
	}
	d.Partial(false)

	log.Printf("[DEBUG] Read Public Cloud Kubernetes Node %s", r)
	return nil
}

func resourceCloudKubernetesNodePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will delete public cloud kubernetes cluster %s from project: %s", id, projectId)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", projectId, clusterId, id)

	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		diags = append(diags, diag.Errorf("calling Delete %s:\n\t %q", endpoint, err)...)
		return
	}

	log.Printf("[DEBUG] Deleting Public Cloud Kubernetes Node %s from project %s:", id, projectId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForCloudKubernetesNodePoolDelete(config.OVHClient, projectId, clusterId, id),
		Timeout:    45 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		diags = append(diags, diag.Errorf("deleting Public Cloud Kubernetes Node %s from project %s", id, projectId)...)
		return
	}
	log.Printf("[DEBUG] Deleted Public Cloud Kubernetes Node %s from project %s", id, projectId)

	d.SetId("")

	return
}

func cloudKubernetesNodePoolExists(projectId string, clusterId string, id string, c *ovh.Client) error {
	r := &CloudKubernetesNodePoolResponse{}

	log.Printf("[DEBUG] Will read public cloud Kubernetes Node for project: %s, id: %s", projectId, id)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", projectId, clusterId, id)

	err := c.Get(endpoint, r)
	if err != nil {
		return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
	}
	log.Printf("[DEBUG] Read public cloud Kubernetes Node: %s", r)

	return nil
}

func readCloudKubernetesNodePool(projectId string, config *Config, d *schema.ResourceData, cluster *CloudKubernetesNodePoolResponse) (err error) {
	_ = d.Set("name", cluster.Name)
	_ = d.Set("flavor", cluster.Flavor)
	_ = d.Set("desired_nodes", cluster.DesiredNodes)
	_ = d.Set("max_nodes", cluster.MaxNodes)
	_ = d.Set("min_nodes", cluster.MinNodes)
	_ = d.Set("monthly_billed", cluster.MonthlyBilled)
	_ = d.Set("status", cluster.Status)
	d.SetId(cluster.Id)
	err = nil
	return
}

func waitForCloudKubernetesNodePoolActive(c *ovh.Client, projectId, clusterId string, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudKubernetesNodeResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", projectId, clusterId, id)
		err := c.Get(endpoint, r)
		if err != nil {
			return r, "", err
		}

		log.Printf("[DEBUG] Pending User: %s", r)
		return r, r.Status, nil
	}
}

func waitForCloudKubernetesNodePoolDelete(c *ovh.Client, projectId, clusterId string, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudKubernetesNodeResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", projectId, clusterId, id)
		err := c.Get(endpoint, r)
		if err != nil {
			if err.(*ovh.APIError).Code == 404 {
				log.Printf("[DEBUG] kubernetes cluster %s on project %s deleted", clusterId, projectId)
				return r, "DELETED", nil
			} else {
				return r, "", err
			}
		}

		log.Printf("[DEBUG] Pending Kubernetes Node: %s", r)
		return r, r.Status, nil
	}
}

func waitForCloudKubernetesNodePoolToBeReal(client *ovh.Client, projectId string, clusterId string, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudKubernetesNodeResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", projectId, clusterId, id)
		err := client.Get(endpoint, r)
		if err != nil {
			if err.(*ovh.APIError).Code == 404 {
				log.Printf("[DEBUG] kubernetes cluster %s on project %s deleted", id, projectId)
				return r, "NOT_FOUND", nil
			} else {
				return r, "", err
			}
		}

		log.Printf("[DEBUG] Pending Kubernetes Node: %s", r)
		return r, "FOUND", nil
	}
}
