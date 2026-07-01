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
			kubeServiceNameKey: {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			kubeKubeIdKey: {
				Type:        schema.TypeString,
				Description: "Kube ID",
				Required:    true,
			},

			// Computed
			kubeNodesKey: {
				Type:        schema.TypeList,
				Description: "Nodes composing the cluster",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeCreatedAtKey: {
							Type:        schema.TypeString,
							Description: "Creation date",
							Computed:    true,
						},
						kubeNodeDeployedAtKey: {
							Type:        schema.TypeString,
							Description: "Node deployment date",
							Computed:    true,
						},
						kubeFlavorKey: {
							Type:        schema.TypeString,
							Description: "Flavor name",
							Computed:    true,
						},
						kubeNodeIdKey: {
							Type:        schema.TypeString,
							Description: "Node ID",
							Computed:    true,
						},
						kubeNodeInstanceIdKey: {
							Type:        schema.TypeString,
							Description: "Public Cloud instance ID",
							Computed:    true,
						},
						kubeNodeIsUpToDateKey: {
							Type:        schema.TypeBool,
							Description: "True if the node is up to date",
							Computed:    true,
						},
						kubeNameKey: {
							Type:        schema.TypeString,
							Description: "Node name",
							Computed:    true,
						},
						kubeNodePoolIdKey: {
							Type:        schema.TypeString,
							Description: "NodePool parent ID",
							Computed:    true,
						},
						kubeProjectIdKey: {
							Type:        schema.TypeString,
							Description: "Project ID",
							Computed:    true,
						},
						kubeStatusKey: {
							Type:        schema.TypeString,
							Description: "Current status",
							Computed:    true,
						},
						kubeUpdatedAtKey: {
							Type:        schema.TypeString,
							Description: "Last update date",
							Computed:    true,
						},
						kubeVersionKey: {
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
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeId := d.Get(kubeKubeIdKey).(string)

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
	d.Set(kubeNodesKey, nodes)

	log.Printf("[DEBUG] Read nodes: %+v", res)
	return nil
}
