package ovh

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ovh/go-ovh/ovh"
)

func resourceCloudKubernetesNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudKubernetesNodeCreate,
		Read:   resourceCloudKubernetesNodeRead,
		Delete: resourceCloudKubernetesNodeDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				return []*schema.ResourceData{d}, nil
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
				Required: false,
				Optional: true,
				Default:  "",
				ForceNew: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_up_to_date": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kubernetes_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCloudKubernetesNodeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	clusterId := d.Get("cluster_id").(string)
	params := &CloudKubernetesNodeCreationRequest{
		Name: d.Get("name").(string),
		// TODO: Write check to ensure flavor is supported
		FlavorName: d.Get("flavor").(string),
	}

	r := &CloudKubernetesNodeResponse{}

	log.Printf("[DEBUG] Will create public cloud kubernetes cluster: %s", params)

	d.Partial(true)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/node", projectId, clusterId)

	err := config.OVHClient.Post(endpoint, params, r)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %s:\n\t %q", endpoint, params, err)
	}

	// This is a fix for a weird bug where the cluster is not immediately available on API
	bugFixWait := &resource.StateChangeConf{
		Pending:    []string{"NOT_FOUND"},
		Target:     []string{"FOUND"},
		Refresh:    waitForCloudKubernetesNodeToBeReal(config.OVHClient, projectId, clusterId, r.Id),
		Timeout:    2 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	log.Printf("[DEBUG] Waiting for Node %s:", r)
	_, err = bugFixWait.WaitForState()

	log.Printf("[DEBUG] Waiting for Node %s:", r)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INSTALLING"},
		Target:     []string{"READY"},
		Refresh:    waitForCloudKubernetesNodeActive(config.OVHClient, projectId, clusterId, r.Id),
		Timeout:    10 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("waiting for user (%s): %s", params, err)
	}
	log.Printf("[DEBUG] Created User %s", r)

	d.SetId(r.Id)
	err = resourceCloudKubernetesNodeRead(d, meta)

	if err != nil {
		return fmt.Errorf("error while reading cloud config: %s", err)
	}

	d.Partial(false)

	return nil
}

func resourceCloudKubernetesNodeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	clusterId := d.Get("cluster_id").(string)

	d.Partial(true)
	r := &CloudKubernetesNodeResponse{}

	log.Printf("[DEBUG] Will read public cloud kubernetes cluster %s from project: %s", d.Id(), projectId)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/node/%s", projectId, clusterId, d.Id())

	err := config.OVHClient.Get(endpoint, r)
	if err != nil {
		return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
	}

	err = readCloudKubernetesNode(projectId, config, d, r)
	if err != nil {
		return fmt.Errorf("error while reading cluster data %s:\n\t %q", d.Id(), err)
	}
	d.Partial(false)

	log.Printf("[DEBUG] Read Public Cloud Kubernetes Node %s", r)
	return nil
}

func resourceCloudKubernetesNodeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	clusterId := d.Get("cluster_id").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will delete public cloud kubernetes cluster %s from project: %s", id, projectId)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/node/%s", projectId, clusterId, id)

	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return fmt.Errorf("calling Delete %s:\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] Deleting Public Cloud Kubernetes Node %s from project %s:", id, projectId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForCloudKubernetesNodeDelete(config.OVHClient, projectId, clusterId, id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("deleting Public Cloud Kubernetes Node %s from project %s", id, projectId)
	}
	log.Printf("[DEBUG] Deleted Public Cloud Kubernetes Node %s from project %s", id, projectId)

	d.SetId("")

	return nil
}

func cloudKubernetesNodeExists(projectId string, clusterId string, id string, c *ovh.Client) error {
	r := &CloudKubernetesNodeResponse{}

	log.Printf("[DEBUG] Will read public cloud Kubernetes Node for project: %s, id: %s", projectId, id)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/node/%s", projectId, clusterId, id)

	err := c.Get(endpoint, r)
	if err != nil {
		return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
	}
	log.Printf("[DEBUG] Read public cloud Kubernetes Node: %s", r)

	return nil
}

func readCloudKubernetesNode(projectId string, config *Config, d *schema.ResourceData, cluster *CloudKubernetesNodeResponse) (err error) {
	_ = d.Set("is_up_to_date", cluster.IsUpToDate)
	_ = d.Set("name", cluster.Name)
	_ = d.Set("status", cluster.Status)
	_ = d.Set("kubernetes_version", cluster.Version)
	err = nil
	return
}

func waitForCloudKubernetesNodeActive(c *ovh.Client, projectId, clusterId string, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudKubernetesNodeResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/node/%s", projectId, clusterId, id)
		err := c.Get(endpoint, r)
		if err != nil {
			return r, "", err
		}

		log.Printf("[DEBUG] Pending User: %s", r)
		return r, r.Status, nil
	}
}

func waitForCloudKubernetesNodeDelete(c *ovh.Client, projectId, clusterId string, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudKubernetesNodeResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/node/%s", projectId, clusterId, id)
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

func waitForCloudKubernetesNodeToBeReal(client *ovh.Client, projectId string, clusterId string, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudKubernetesNodeResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/node/%s", projectId, clusterId, id)
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
