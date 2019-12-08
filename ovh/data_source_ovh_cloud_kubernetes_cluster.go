package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	clientcmd "k8s.io/client-go/tools/clientcmd"
)

func dataSourceCloudKubernetesCluster() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudKubernetesClusterRead,
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
			"id": {
				Type:     schema.TypeString,
				Computed: true,
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
				Computed: true,
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
			"version": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func dataSourceCloudKubernetesClusterRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	projectId := d.Get("project_id").(string)
	name := d.Get("name").(string)

	log.Printf("[DEBUG] Will read public cloud region %s for project: %s", name, projectId)
	d.Partial(true)

	cluster, err := findKubernetesClusterByName(config, projectId, name)

	if err != nil {
		return err
	}

	_ = d.Set("control_plane_is_up_to_date", cluster.ControlPlaneIsUpToDate)
	_ = d.Set("is_up_to_date", cluster.IsUpToDate)
	_ = d.Set("name", cluster.Name)
	_ = d.Set("next_upgrade_versions", cluster.NextUpgradeVersions)
	_ = d.Set("nodes_url", cluster.NodesUrl)
	_ = d.Set("region", cluster.Region)
	_ = d.Set("status", cluster.Status)
	_ = d.Set("update_policy", cluster.UpdatePolicy)
	_ = d.Set("url", cluster.Url)
	_ = d.Set("version", cluster.Version)

	kubeconfig_raw := PublicCloudKubernetesKubeConfigResponse{}
	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/kubeconfig", projectId, cluster.Id)
	err = config.OVHClient.Post(endpoint, nil, &kubeconfig_raw)

	if err != nil {
		return err
	}

	kubeconfig, err := clientcmd.Load([]byte(kubeconfig_raw.Content))
	if err != nil {
		return err
	}
	currentContext := kubeconfig.CurrentContext
	currentUser := kubeconfig.Contexts[currentContext].AuthInfo
	currentCluster := kubeconfig.Contexts[currentContext].Cluster
	_ = d.Set("client_certificate", string(kubeconfig.AuthInfos[currentUser].ClientCertificateData))
	_ = d.Set("client_key", string(kubeconfig.AuthInfos[currentUser].ClientKeyData))
	_ = d.Set("cluster_ca_certificate", string(kubeconfig.Clusters[currentCluster].CertificateAuthorityData))

	d.Partial(false)
	d.SetId(cluster.Id)

	return nil
}

func findKubernetesClusterByName(config *Config, projectId string, name string) (cluster PublicCloudKubernetesClusterResponse, err error) {
	cluster = PublicCloudKubernetesClusterResponse{}
	response := []string{}
	endpoint := fmt.Sprintf("/cloud/project/%s/kube", projectId)
	err = config.OVHClient.Get(endpoint, &response)

	if err != nil {
		return
	}

	for i := 0; i < len(response); i++ {
		id := response[i]
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", projectId, id)
		err = config.OVHClient.Get(endpoint, &cluster)

		if err != nil {
			return
		}

		if cluster.Name == name {
			err = nil
			return
		}
	}
	err = fmt.Errorf("unable to find cluster named %s", name)
	return
}
