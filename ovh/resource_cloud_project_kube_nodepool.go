package ovh

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

func resourceCloudProjectKubeNodePool() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectKubeNodePoolCreate,
		Read:   resourceCloudProjectKubeNodePoolRead,
		Delete: resourceCloudProjectKubeNodePoolDelete,
		Update: resourceCloudProjectKubeNodePoolUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectKubeNodePoolImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(time.Hour),
			Update:  schema.DefaultTimeout(time.Hour),
			Delete:  schema.DefaultTimeout(time.Hour),
			Read:    schema.DefaultTimeout(5 * time.Minute),
			Default: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			kubeServiceNameKey: {
				Type:        schema.TypeString,
				Description: "Service name",
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			kubeKubeIdKey: {
				Type:        schema.TypeString,
				Description: "Kube ID",
				Required:    true,
				ForceNew:    true,
			},
			kubeNodePoolAttachFloatingIpsKey: {
				Description: "Floating IPs attachment configuration for pool nodes",
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeSet,
				MaxItems:    1,
				Set:         CustomSchemaSetFunc(),
				Default:     nil,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeClusterCustomizationEnabledKey: {
							Type:        schema.TypeBool,
							Description: "Whether floating IPs attachment is enabled on nodes of this pool",
							Optional:    true,
							Computed:    true,
							ForceNew:    false,
							Default:     nil,
						},
					},
				},
			},
			kubeNodePoolAutoscaleKey: {
				Type:        schema.TypeBool,
				Description: "Enable auto-scaling for the pool",
				Optional:    true,
				Computed:    true,
				ForceNew:    false,
			},
			kubeNodePoolAutoscalingScaleDownUnneededTimeSecondsKey: {
				Description: "scaleDownUnneededTimeSeconds for autoscaling",
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
			},
			kubeNodePoolAutoscalingScaleDownUnreadyTimeSecondsKey: {
				Description: "scaleDownUnreadyTimeSeconds for autoscaling",
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
			},
			kubeNodePoolAutoscalingScaleDownUtilizationThresholdKey: {
				Description: "scaleDownUtilizationThreshold for autoscaling",
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeFloat,
			},
			kubeNodePoolAntiAffinityKey: {
				Type:        schema.TypeBool,
				Description: "Enable anti affinity groups for nodes in the pool",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			kubeNodePoolFlavorNameKey: {
				Type:        schema.TypeString,
				Description: "Flavor name",
				Required:    true,
				ForceNew:    true,
			},
			kubeNodePoolDesiredNodesKey: {
				Type:        schema.TypeInt,
				Description: "Number of nodes you desire in the pool",
				Optional:    true,
				Computed:    true,
			},
			kubeNameKey: {
				Type:        schema.TypeString,
				Description: "NodePool resource name",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},
			kubeNodePoolMaxNodesKey: {
				Type:        schema.TypeInt,
				Description: "Number of nodes you desire in the pool",
				Computed:    true,
				Optional:    true,
			},
			kubeNodePoolMinNodesKey: {
				Type:        schema.TypeInt,
				Description: "Number of nodes you desire in the pool",
				Computed:    true,
				Optional:    true,
			},
			kubeNodePoolMonthlyBilledKey: {
				Type:        schema.TypeBool,
				Description: "Enable monthly billing on all nodes in the pool",
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
			},

			// computed
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
			kubeNodePoolTemplateKey: {
				Description: "Node pool template",
				Optional:    true,
				Type:        schema.TypeSet,
				MaxItems:    1,
				Set:         CustomSchemaSetFunc(),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeNodePoolTemplateMetadataKey: {
							Description: "metadata",
							Required:    true,
							Type:        schema.TypeSet,
							MaxItems:    1,
							Set:         CustomSchemaSetFunc(),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									kubeNodePoolTemplateFinalizersKey: {
										Description: "finalizers",
										Required:    true,
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									kubeNodePoolTemplateLabelsKey: {
										Description: "labels",
										Required:    true,
										Type:        schema.TypeMap,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
									},
									kubeNodePoolTemplateAnnotationsKey: {
										Description: "annotations",
										Required:    true,
										Type:        schema.TypeMap,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Set:         schema.HashString,
									},
								},
							},
						},
						kubeNodePoolTemplateSpecKey: {
							Description: "spec",
							Required:    true,
							Type:        schema.TypeSet,
							MaxItems:    1,
							Set:         CustomSchemaSetFunc(),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									kubeNodePoolTemplateUnschedulableKey: {
										Description: "unschedulable",
										Required:    true,
										Type:        schema.TypeBool,
									},
									kubeNodePoolTemplateTaintsKey: {
										Description: "taints",
										Required:    true,
										Type:        schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeMap,
											Set:  schema.HashString,
											ValidateFunc: func(taintInterface interface{}, path string) (warning []string, errorList []error) {
												taint := taintInterface.(map[string]interface{})

												if taint[kubeNodePoolTemplateTaintKeyKey] == nil {
													return nil, []error{fmt.Errorf("key attribute is mandatory for taint: %s", path)}
												}

												if taint[kubeNodePoolTemplateTaintEffectKey] == nil {
													return nil, []error{fmt.Errorf("effect attribute is mandatory for taint: %s", path)}
												}

												effectString := taint[kubeNodePoolTemplateTaintEffectKey].(string)
												effect := TaintEffecTypeToID[effectString]
												if effect == NotATaint {
													return nil, []error{fmt.Errorf("effect: %s is not a allowable taint %#v", effectString, TaintEffecTypeToID)}
												}

												return nil, nil
											},
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
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceCloudProjectKubeNodePoolImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 3)
	if len(splitId) != 3 {
		return nil, fmt.Errorf("import Id is not service_name/kubeid/poolid formatted")
	}
	serviceName := splitId[0]
	kubeId := splitId[1]
	id := splitId[2]
	d.SetId(id)
	d.Set(kubeKubeIdKey, kubeId)
	d.Set(kubeServiceNameKey, serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectKubeNodePoolCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeId := d.Get(kubeKubeIdKey).(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool", serviceName, kubeId)
	params, err := (&CloudProjectKubeNodePoolCreateOpts{}).FromResource(d)
	if err != nil {
		return err
	}
	res := &CloudProjectKubeNodePoolResponse{}

	log.Printf("[DEBUG] Will create nodepool: %+v", params)
	err = config.OVHClient.Post(endpoint, params, res)
	if err != nil {
		return fmt.Errorf("calling Post %s with params %s:\n\t %w", endpoint, params, err)
	}

	// This is a fix for a weird bug where the nodepool is not immediately available on API
	log.Printf("[DEBUG] Waiting for nodepool %s to be available", res.Id)
	endpoint = fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, res.Id)
	err = helpers.WaitAvailable(config.OVHClient, endpoint, 2*time.Minute)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Waiting for nodepool %s to be READY or ERROR", res.Id)
	err = waitForCloudProjectKubeNodePoolWithStateTarget(config.OVHClient, serviceName, kubeId, res.Id, d.Timeout(schema.TimeoutCreate), []string{"READY", "ERROR"})
	if err != nil {
		return fmt.Errorf("timeout while waiting nodepool %s to be READY: %w", res.Id, err)
	}
	log.Printf("[DEBUG] nodepool %s is READY", res.Id)

	d.SetId(res.Id)

	return resourceCloudProjectKubeNodePoolRead(d, meta)
}

func resourceCloudProjectKubeNodePoolRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeId := d.Get(kubeKubeIdKey).(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, d.Id())
	res := &CloudProjectKubeNodePoolResponse{}

	log.Printf("[DEBUG] Will read nodepool %s from cluster %s in project %s", d.Id(), kubeId, serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	log.Printf("[DEBUG] Read nodepool: %+v", res)
	return nil
}

func resourceCloudProjectKubeNodePoolUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeId := d.Get(kubeKubeIdKey).(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, d.Id())
	params, err := (&CloudProjectKubeNodePoolUpdateOpts{}).FromResource(d)
	if err != nil {
		return err
	}

	if params.MaxNodes != nil && params.DesiredNodes == nil {
		current := &CloudProjectKubeNodePoolResponse{}
		if err := config.OVHClient.Get(endpoint, current); err == nil {
			if current.DesiredNodes > *params.MaxNodes {
				log.Printf("[DEBUG] desired_nodes (%d) exceeds new max_nodes (%d), capping desired_nodes to max_nodes", current.DesiredNodes, *params.MaxNodes)
				params.DesiredNodes = params.MaxNodes
			}
		} else {
			log.Printf("[WARN] failed to GET existing nodepool %s before update (skipping desired_nodes capping): %v", d.Id(), err)
			if apiErr, ok := err.(*ovh.APIError); ok && apiErr.Code == 404 {
				// Nodepool not found when attempting pre-update GET; proceed without capping.
			} else {
				return fmt.Errorf("failed to GET existing nodepool %s before update: %w", d.Id(), err)
			}
		}
	}

	log.Printf("[DEBUG] Will update nodepool: %#v", *params)
	err = config.OVHClient.Put(endpoint, params, nil)
	if err != nil {
		return fmt.Errorf("calling Put %s with params %v:\n\t %w", endpoint, *params, err)
	}

	log.Printf("[DEBUG] Waiting for nodepool %s to be READY", d.Id())
	err = waitForCloudProjectKubeNodePoolWithStateTarget(config.OVHClient, serviceName, kubeId, d.Id(), d.Timeout(schema.TimeoutUpdate), []string{"READY"})
	if err != nil {
		return fmt.Errorf("timeout while waiting nodepool %s to be READY: %w", d.Id(), err)
	}
	log.Printf("[DEBUG] nodepool %s is READY", d.Id())

	return resourceCloudProjectKubeNodePoolRead(d, meta)
}

func resourceCloudProjectKubeNodePoolDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeId := d.Get(kubeKubeIdKey).(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, d.Id())

	log.Printf("[DEBUG] Will delete nodepool %s from cluster %s in project %s", d.Id(), kubeId, serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	log.Printf("[DEBUG] Waiting for nodepool %s to be DELETED", d.Id())
	err = waitForCloudProjectKubeNodePoolDeleted(config.OVHClient, serviceName, kubeId, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmt.Errorf("timeout while waiting nodepool %s to be DELETED: %v", d.Id(), err)
	}
	log.Printf("[DEBUG] nodepool %s is DELETED", d.Id())

	d.SetId("")

	return nil
}

func cloudProjectKubeNodePoolExists(serviceName, kubeId, id string, client *ovhwrap.Client) error {
	res := &CloudProjectKubeNodePoolResponse{}

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, id)
	return client.Get(endpoint, res)
}

func waitForCloudProjectKubeNodePoolWithStateTarget(client *ovhwrap.Client, serviceName, kubeId, id string, timeout time.Duration, stateTargets []string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"INSTALLING", "UPDATING", "REDEPLOYING", "RESIZING", "DOWNSCALING", "UPSCALING", "UNKNOWN"},
		Target:  stateTargets,
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectKubeNodePoolResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/nodepool/%s", serviceName, kubeId, id)
			err := client.Get(endpoint, res)
			if err != nil {
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectKubeNodePoolDeleted(client *ovhwrap.Client, serviceName, kubeId, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectKubeNodePoolResponse{}
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
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}
