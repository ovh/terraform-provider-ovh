package ovh

import (
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/ovh/go-ovh/ovh"
)

func resourceCloudKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudKubernetesClusterCreate,
		Read:   resourceCloudKubernetesClusterRead,
		Delete: resourceCloudKubernetesClusterDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				err := resourceCloudKubernetesClusterRead(d, meta)
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Default:  "GRA5",
				Optional: true,
				ForceNew: true,
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
				ForceNew: true,
			},
			"client_certificate": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"client_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"cluster_ca_certificate": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"kubeconfig": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

type PublicCloudKubernetesClusterCreateOpts struct {
	Name    string `json:"name"`
	Region  string `json:"region"`
	Version string `json:"version"`
}

func resourceCloudKubernetesClusterCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	params := &PublicCloudKubernetesClusterCreateOpts{
		Name: d.Get("name").(string),
		// TODO: Write check to ensure selected region is available for Kubernetes cluster
		Region: d.Get("region").(string),
		// TODO: Write check to ensure version is supported
		Version: d.Get("version").(string),
	}

	r := &PublicCloudKubernetesClusterResponse{}

	log.Printf("[DEBUG] Will create public cloud kubernetes cluster: %s", params)

	d.Partial(true)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube", projectId)

	err := config.OVHClient.Post(endpoint, params, r)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %s:\n\t %q", endpoint, params, err)
	}

	// This is a fix for a weird bug where the cluster is not immediately available on API
	bugFixWait := &resource.StateChangeConf{
		Pending:    []string{"NOT_FOUND"},
		Target:     []string{"FOUND"},
		Refresh:    waitForCloudKubernetesClusterToBeReal(config.OVHClient, projectId, r.Id),
		Timeout:    30 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	log.Printf("[DEBUG] Waiting for Cluster %s:", r)
	_, err = bugFixWait.WaitForState()

	log.Printf("[DEBUG] Waiting for Cluster %s:", r)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INSTALLING"},
		Target:     []string{"READY"},
		Refresh:    waitForCloudKubernetesClusterActive(config.OVHClient, projectId, r.Id),
		Timeout:    20 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("waiting for user (%s): %s", params, err)
	}
	log.Printf("[DEBUG] Created User %s", r)

	d.SetId(r.Id)
	err = resourceCloudKubernetesClusterRead(d, meta)

	if err != nil {
		return fmt.Errorf("error while reading cloud config: %s", err)
	}

	d.Partial(false)

	return nil
}

func resourceCloudKubernetesClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)

	d.Partial(true)
	r := &PublicCloudKubernetesClusterResponse{}

	log.Printf("[DEBUG] Will read public cloud kubernetes cluster %s from project: %s", d.Id(), projectId)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", projectId, d.Id())

	err := config.OVHClient.Get(endpoint, r)
	if err != nil {
		return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
	}

	err = readCloudKubernetesCluster(projectId, config, d, r)
	if err != nil {
		return fmt.Errorf("error while reading cluster data %s:\n\t %q", d.Id(), err)
	}
	d.Partial(false)

	log.Printf("[DEBUG] Read Public Cloud Kubernetes Cluster %s", r)
	return nil
}

func resourceCloudKubernetesClusterDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will delete public cloud kubernetes cluster %s from project: %s", id, projectId)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", projectId, id)

	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return fmt.Errorf("calling Delete %s:\n\t %q", endpoint, err)
	}

	log.Printf("[DEBUG] Deleting Public Cloud Kubernetes Cluster %s from project %s:", id, projectId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForCloudKubernetesClusterDelete(config.OVHClient, projectId, id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("deleting Public Cloud Kubernetes Cluster %s from project %s", id, projectId)
	}
	log.Printf("[DEBUG] Deleted Public Cloud Kubernetes Cluster %s from project %s", id, projectId)

	d.SetId("")

	return nil
}

func cloudKubernetesClusterExists(projectId, id string, c *ovh.Client) error {
	r := &PublicCloudKubernetesClusterResponse{}

	log.Printf("[DEBUG] Will read public cloud Kubernetes Cluster for project: %s, id: %s", projectId, id)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", projectId, id)

	err := c.Get(endpoint, r)
	if err != nil {
		return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
	}
	log.Printf("[DEBUG] Read public cloud Kubernetes Cluster: %s", r)

	return nil
}

func readCloudKubernetesCluster(projectId string, config *Config, d *schema.ResourceData, cluster *PublicCloudKubernetesClusterResponse) (err error) {
	_ = d.Set("control_plane_is_up_to_date", cluster.ControlPlaneIsUpToDate)
	_ = d.Set("is_up_to_date", cluster.IsUpToDate)
	_ = d.Set("name", cluster.Name)
	_ = d.Set("next_upgrade_versions", cluster.NextUpgradeVersions)
	_ = d.Set("nodes_url", cluster.NodesUrl)
	_ = d.Set("region", cluster.Region)
	_ = d.Set("status", cluster.Status)
	_ = d.Set("update_policy", cluster.UpdatePolicy)
	_ = d.Set("url", cluster.Url)
	_ = d.Set("kubernetes_version", cluster.Version)

	kubeconfigRaw := PublicCloudKubernetesKubeConfigResponse{}
	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/kubeconfig", projectId, cluster.Id)
	err = config.OVHClient.Post(endpoint, nil, &kubeconfigRaw)

	if err != nil {
		return err
	}
	_ = d.Set("kubeconfig", kubeconfigRaw.Content)

	kubeconfig, err := clientcmd.Load([]byte(kubeconfigRaw.Content))
	if err != nil {
		return err
	}
	currentContext := kubeconfig.CurrentContext
	currentUser := kubeconfig.Contexts[currentContext].AuthInfo
	currentCluster := kubeconfig.Contexts[currentContext].Cluster
	_ = d.Set("client_certificate", string(kubeconfig.AuthInfos[currentUser].ClientCertificateData))
	_ = d.Set("client_key", string(kubeconfig.AuthInfos[currentUser].ClientKeyData))
	_ = d.Set("cluster_ca_certificate", string(kubeconfig.Clusters[currentCluster].CertificateAuthorityData))

	err = nil
	return
}

func waitForCloudKubernetesClusterActive(c *ovh.Client, projectId, cloudKubernetesClusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &PublicCloudKubernetesClusterResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", projectId, cloudKubernetesClusterId)
		err := c.Get(endpoint, r)
		if err != nil {
			return r, "", err
		}

		log.Printf("[DEBUG] Pending User: %s", r)
		return r, r.Status, nil
	}
}

func waitForCloudKubernetesClusterDelete(c *ovh.Client, projectId, CloudKubernetesClusterId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &PublicCloudKubernetesClusterResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", projectId, CloudKubernetesClusterId)
		err := c.Get(endpoint, r)
		if err != nil {
			if err.(*ovh.APIError).Code == 404 {
				log.Printf("[DEBUG] kubernetes cluster %s on project %s deleted", CloudKubernetesClusterId, projectId)
				return r, "DELETED", nil
			} else {
				return r, "", err
			}
		}

		log.Printf("[DEBUG] Pending Kubernetes Cluster: %s", r)
		return r, r.Status, nil
	}
}

func waitForCloudKubernetesClusterToBeReal(client *ovh.Client, projectId string, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &PublicCloudKubernetesClusterResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", projectId, id)
		err := client.Get(endpoint, r)
		if err != nil {
			if err.(*ovh.APIError).Code == 404 {
				log.Printf("[DEBUG] kubernetes cluster %s on project %s deleted", id, projectId)
				return r, "NOT_FOUND", nil
			} else {
				return r, "", err
			}
		}

		log.Printf("[DEBUG] Pending Kubernetes Cluster: %s", r)
		return r, "FOUND", nil
	}
}
