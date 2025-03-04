package ovh

import (
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers/hashcode"
)

func dataSourceCloudProjectKubeNodepoolNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectKubeNodePoolNodesRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"kube_id": {
				Type:        schema.TypeString,
				Description: "Kube ID",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "NodePool resource name",
				Required:    true,
			},

			// Computed
			"nodes": {
				Type:        schema.TypeList,
				Description: "Nodes composing the node pool",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": {
							Type:        schema.TypeString,
							Description: "Creation date",
							Computed:    true,
						},
						"deployed_at": {
							Type:        schema.TypeString,
							Description: "Node deployment date",
							Computed:    true,
						},
						"flavor": {
							Type:        schema.TypeString,
							Description: "Flavor name",
							Computed:    true,
						},
						"id": {
							Type:        schema.TypeString,
							Description: "Node ID",
							Computed:    true,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Description: "Public Cloud instance ID",
							Computed:    true,
						},
						"is_up_to_date": {
							Type:        schema.TypeBool,
							Description: "True if the node is up to date",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Node name",
							Computed:    true,
						},
						"node_pool_id": {
							Type:        schema.TypeString,
							Description: "NodePool parent ID",
							Computed:    true,
						},
						"project_id": {
							Type:        schema.TypeString,
							Description: "Project ID",
							Computed:    true,
						},
						"status": {
							Type:        schema.TypeString,
							Description: "Current status",
							Computed:    true,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Description: "Last update date",
							Computed:    true,
						},
						"version": {
							Type:        schema.TypeString,
							Description: "Node version",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudProjectKubeNodePoolNodesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)
	name := d.Get("name").(string)

	endpointNodepool := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool", serviceName, kubeId)
	var nodePoolRes []CloudProjectKubeNodePoolResponse
	log.Printf("[DEBUG] Will read nodepools from cluster %s in project %s", kubeId, serviceName)
	if err := config.OVHClient.Get(endpointNodepool, &nodePoolRes); err != nil {
		return helpers.CheckDeleted(d, err, endpointNodepool)
	}

	var nodepoolTarget *CloudProjectKubeNodePoolResponse

	for _, nodepool := range nodePoolRes {
		if nodepool.Name == name {
			nodepoolTarget = &nodepool
			break
		}
	}

	if nodepoolTarget == nil {
		return fmt.Errorf("The nodepool named %s cannot be found for cluster %s in project %s", name, kubeId, serviceName)
	}

	endpointNodepoolNodes := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s/nodes",
		url.PathEscape(serviceName),
		url.PathEscape(kubeId),
		url.PathEscape(nodepoolTarget.Id))
	var resNodepoolNodes []CloudProjectKubeNodeResponse

	log.Printf("[DEBUG] Will read nodes from node pool %s and cluster %s in project %s", name, kubeId, serviceName)
	if err := config.OVHClient.Get(endpointNodepoolNodes, &resNodepoolNodes); err != nil {
		return helpers.CheckDeleted(d, err, endpointNodepoolNodes)
	}

	nodes := make([]map[string]interface{}, len(resNodepoolNodes))
	ids := make([]string, len(resNodepoolNodes))

	for i, node := range resNodepoolNodes {
		nodes[i] = node.ToMap()
		ids = append(ids, node.Id)
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(ids)

	d.SetId(hashcode.Strings(ids))
	d.Set("nodes", nodes)

	log.Printf("[DEBUG] Read nodepool nodes: %+v", resNodepoolNodes)
	return nil
}
