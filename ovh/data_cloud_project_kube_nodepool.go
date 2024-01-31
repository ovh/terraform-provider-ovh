package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func dataSourceCloudProjectKubeNodepool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectKubeNodePoolRead,
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

			// computed
			"autoscale": {
				Type:        schema.TypeBool,
				Description: "Enable auto-scaling for the pool",
				Computed:    true,
			},
			"anti_affinity": {
				Type:        schema.TypeBool,
				Description: "Enable anti affinity groups for nodes in the pool",
				Computed:    true,
			},
			"flavor_name": {
				Type:        schema.TypeString,
				Description: "Flavor name",
				Computed:    true,
			},
			"desired_nodes": {
				Type:        schema.TypeInt,
				Description: "Number of nodes you desire in the pool",
				Computed:    true,
			},
			"max_nodes": {
				Type:        schema.TypeInt,
				Description: "Number of nodes you desire in the pool",
				Computed:    true,
			},
			"min_nodes": {
				Type:        schema.TypeInt,
				Description: "Number of nodes you desire in the pool",
				Computed:    true,
			},
			"monthly_billed": {
				Type:        schema.TypeBool,
				Description: "Enable monthly billing on all nodes in the pool",
				Computed:    true,
			},
			"available_nodes": {
				Type:        schema.TypeInt,
				Description: "Number of nodes which are actually ready in the pool",
				Computed:    true,
			},
			"created_at": {
				Type:        schema.TypeString,
				Description: "Creation date",
				Computed:    true,
			},
			"current_nodes": {
				Type:        schema.TypeInt,
				Description: "Number of nodes present in the pool",
				Computed:    true,
			},
			"flavor": {
				Type:        schema.TypeString,
				Description: "Flavor name",
				Computed:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "Project id",
				Computed:    true,
			},
			"size_status": {
				Type:        schema.TypeString,
				Description: "Status describing the state between number of nodes wanted and available ones",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Current status",
				Computed:    true,
			},
			"up_to_date_nodes": {
				Type:        schema.TypeInt,
				Description: "Number of nodes with latest version installed in the pool",
				Computed:    true,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "Last update date",
				Computed:    true,
			},
			"autoscaling_scale_down_unneeded_time_seconds": {
				Type:        schema.TypeInt,
				Description: "scaleDownUnneededTimeSeconds for autoscaling",
				Computed:    true,
			},
			"autoscaling_scale_down_unready_time_seconds": {
				Type:        schema.TypeInt,
				Description: "scaleDownUnreadyTimeSeconds for autoscaling",
				Computed:    true,
			},
			"autoscaling_scale_down_utilization_threshold": {
				Type:        schema.TypeFloat,
				Description: "scaleDownUtilizationThreshold for autoscaling",
				Computed:    true,
			},
			"template": {
				Description: "Node pool template",
				Optional:    true,
				Type:        schema.TypeSet,
				MaxItems:    1,
				Set: func(i interface{}) int {
					out := fmt.Sprintf("%#v", i)
					hash := int(schema.HashString(out))
					return hash
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata": {
							Description: "metadata",
							Optional:    true,
							Type:        schema.TypeSet,
							MaxItems:    1,
							Set: func(i interface{}) int {
								out := fmt.Sprintf("%#v", i)
								hash := int(schema.HashString(out))
								return hash
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"finalizers": {
										Description: "finalizers",
										Optional:    true,
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"labels": {
										Description: "labels",
										Optional:    true,
										Type:        schema.TypeMap,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
									},
									"annotations": {
										Description: "annotations",
										Optional:    true,
										Type:        schema.TypeMap,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
									},
								},
							},
						},
						"spec": {
							Description: "spec",
							Optional:    true,
							Type:        schema.TypeSet,
							MaxItems:    1,
							Set: func(i interface{}) int {
								out := fmt.Sprintf("%#v", i)
								hash := int(schema.HashString(out))
								return hash
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unschedulable": {
										Description: "unschedulable",
										Optional:    true,
										Type:        schema.TypeBool,
									},
									"taints": {
										Description: "taints",
										Optional:    true,
										Type:        schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeMap,
											Set:  schema.HashString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudProjectKubeNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)
	nodepoolName := d.Get("name").(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool", serviceName, kubeId)
	var res []CloudProjectKubeNodePoolResponse
	log.Printf("[DEBUG] Will read nodepools from cluster %s in project %s", kubeId, serviceName)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	var nodepoolTarget *CloudProjectKubeNodePoolResponse

	for _, nodepool := range res {
		if nodepool.Name == nodepoolName {
			nodepoolTarget = &nodepool
			break
		}
	}

	if nodepoolTarget == nil {
		return fmt.Errorf("the nodepool named %s cannot be found for cluster %s in project %s", nodepoolName, kubeId, serviceName)
	}

	for k, v := range nodepoolTarget.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read nodepool: %+v", res)
	return nil
}
