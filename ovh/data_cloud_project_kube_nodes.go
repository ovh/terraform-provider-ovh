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

func dataSourceCloudProjectKubeNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectKubeNodesRead,
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

			// Computed
			"nodes": {
				Type:        schema.TypeList,
				Description: "Nodes composing the cluster",
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

func dataSourceCloudProjectKubeNodesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/node",
		url.PathEscape(serviceName),
		url.PathEscape(kubeId))
	var res []CloudProjectKubeNodeResponse

	log.Printf("[DEBUG] Will read nodes from cluster %s in project %s", kubeId, serviceName)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	nodes := make([]map[string]interface{}, len(res))
	ids := make([]string, len(res))

	for i, node := range res {
		nodes[i] = node.ToMap()
		ids = append(ids, node.Id)
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(ids)

	d.SetId(hashcode.Strings(ids))
	d.Set("nodes", nodes)

	log.Printf("[DEBUG] Read nodes: %+v", res)
	return nil
}
