package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func dataSourceCloudProjectKubeNodepool() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectKubeNodePoolRead,
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
			kubeNameKey: {
				Type:        schema.TypeString,
				Description: "NodePool resource name",
				Required:    true,
			},

			// computed
			kubeNodePoolAttachFloatingIpsKey: {
				Description: "Floating IPs attachment configuration for pool nodes",
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeSet,
				MaxItems:    1,
				Set:         CustomSchemaSetFunc(),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeClusterCustomizationEnabledKey: {
							Type:        schema.TypeBool,
							Description: "Whether floating IPs attachment is enabled on nodes of this pool",
							Optional:    true,
							Computed:    true,
							ForceNew:    false,
						},
					},
				},
			},
			kubeNodePoolAutoscaleKey: {
				Type:        schema.TypeBool,
				Description: "Enable auto-scaling for the pool",
				Computed:    true,
			},
			kubeNodePoolAntiAffinityKey: {
				Type:        schema.TypeBool,
				Description: "Enable anti affinity groups for nodes in the pool",
				Computed:    true,
			},
			kubeNodePoolFlavorNameKey: {
				Type:        schema.TypeString,
				Description: "Flavor name",
				Computed:    true,
			},
			kubeNodePoolDesiredNodesKey: {
				Type:        schema.TypeInt,
				Description: "Number of nodes you desire in the pool",
				Computed:    true,
			},
			kubeNodePoolMaxNodesKey: {
				Type:        schema.TypeInt,
				Description: "Number of nodes you desire in the pool",
				Computed:    true,
			},
			kubeNodePoolMinNodesKey: {
				Type:        schema.TypeInt,
				Description: "Number of nodes you desire in the pool",
				Computed:    true,
			},
			kubeNodePoolMonthlyBilledKey: {
				Type:        schema.TypeBool,
				Description: "Enable monthly billing on all nodes in the pool",
				Computed:    true,
			},
			kubeNodePoolAvailableNodesKey: {
				Type:        schema.TypeInt,
				Description: "Number of nodes which are actually ready in the pool",
				Computed:    true,
			},
			kubeCreatedAtKey: {
				Type:        schema.TypeString,
				Description: "Creation date",
				Computed:    true,
			},
			kubeNodePoolCurrentNodesKey: {
				Type:        schema.TypeInt,
				Description: "Number of nodes present in the pool",
				Computed:    true,
			},
			kubeFlavorKey: {
				Type:        schema.TypeString,
				Description: "Flavor name",
				Computed:    true,
			},
			kubeProjectIdKey: {
				Type:        schema.TypeString,
				Description: "Project id",
				Computed:    true,
			},
			kubeNodePoolSizeStatusKey: {
				Type:        schema.TypeString,
				Description: "Status describing the state between number of nodes wanted and available ones",
				Computed:    true,
			},
			kubeStatusKey: {
				Type:        schema.TypeString,
				Description: "Current status",
				Computed:    true,
			},
			kubeNodePoolUpToDateNodesKey: {
				Type:        schema.TypeInt,
				Description: "Number of nodes with latest version installed in the pool",
				Computed:    true,
			},
			kubeUpdatedAtKey: {
				Type:        schema.TypeString,
				Description: "Last update date",
				Computed:    true,
			},
			kubeNodePoolAutoscalingScaleDownUnneededTimeSecondsKey: {
				Type:        schema.TypeInt,
				Description: "scaleDownUnneededTimeSeconds for autoscaling",
				Computed:    true,
			},
			kubeNodePoolAutoscalingScaleDownUnreadyTimeSecondsKey: {
				Type:        schema.TypeInt,
				Description: "scaleDownUnreadyTimeSeconds for autoscaling",
				Computed:    true,
			},
			kubeNodePoolAutoscalingScaleDownUtilizationThresholdKey: {
				Type:        schema.TypeFloat,
				Description: "scaleDownUtilizationThreshold for autoscaling",
				Computed:    true,
			},
			kubeNodePoolTemplateKey: {
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
						kubeNodePoolTemplateMetadataKey: {
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
									kubeNodePoolTemplateFinalizersKey: {
										Description: "finalizers",
										Optional:    true,
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									kubeNodePoolTemplateLabelsKey: {
										Description: "labels",
										Optional:    true,
										Type:        schema.TypeMap,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
									},
									kubeNodePoolTemplateAnnotationsKey: {
										Description: "annotations",
										Optional:    true,
										Type:        schema.TypeMap,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
									},
								},
							},
						},
						kubeNodePoolTemplateSpecKey: {
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
									kubeNodePoolTemplateUnschedulableKey: {
										Description: "unschedulable",
										Optional:    true,
										Type:        schema.TypeBool,
									},
									kubeNodePoolTemplateTaintsKey: {
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
			kubeNodePoolAvailabilityZonesKey: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceCloudProjectKubeNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeId := d.Get(kubeKubeIdKey).(string)
	nodepoolName := d.Get(kubeNameKey).(string)

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
