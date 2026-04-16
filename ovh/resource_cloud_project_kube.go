package ovh

import (
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

func resourceCloudProjectKube() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectKubeCreate,
		Read:   resourceCloudProjectKubeRead,
		Delete: resourceCloudProjectKubeDelete,
		Update: resourceCloudProjectKubeUpdate,

		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectKubeImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(15 * time.Minute),
			Update:  schema.DefaultTimeout(time.Hour),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
			Read:    schema.DefaultTimeout(5 * time.Minute),
			Default: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			kubeServiceNameKey: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			kubeNameKey: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			kubeVersionKey: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
			kubeClusterPlanKey: {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ForceNew:     false,
				ValidateFunc: helpers.ValidateEnum([]string{"standard", "free"}),
			},
			kubeClusterCustomizationApiServerKey: {
				Type:          schema.TypeSet,
				Computed:      true,
				Optional:      true,
				ForceNew:      false,
				Set:           CustomSchemaSetFunc(),
				ConflictsWith: []string{kubeClusterCustomization},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeClusterCustomizationAdmissionPluginsKey: {
							Type:     schema.TypeSet,
							Computed: true,
							Optional: true,
							ForceNew: false,
							Set:      CustomApiServerAdmissionPluginsSchemaSetFunc(),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									kubeClusterCustomizationEnabledKey: {
										Type:     schema.TypeList,
										Computed: true,
										Optional: true,
										ForceNew: false,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									kubeClusterCustomizationDisabledKey: {
										Type:     schema.TypeList,
										Computed: true,
										Optional: true,
										ForceNew: false,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			kubeClusterCustomization: {
				Type:          schema.TypeSet,
				Computed:      true,
				Optional:      true,
				ForceNew:      false,
				Set:           CustomSchemaSetFunc(),
				ConflictsWith: []string{kubeClusterCustomizationApiServerKey},
				Deprecated:    fmt.Sprintf("Use %s instead", kubeClusterCustomizationApiServerKey),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeClusterCustomizationApiServerNestedKey: {
							Type:       schema.TypeSet,
							Computed:   true,
							Optional:   true,
							ForceNew:   false,
							Set:        CustomSchemaSetFunc(),
							Deprecated: fmt.Sprintf("Use %s instead", kubeClusterCustomizationApiServerKey),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									kubeClusterCustomizationAdmissionPluginsKey: {
										Type:     schema.TypeSet,
										Computed: true,
										Optional: true,
										ForceNew: false,
										Set:      CustomApiServerAdmissionPluginsSchemaSetFunc(),
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												kubeClusterCustomizationEnabledKey: {
													Type:     schema.TypeList,
													Computed: true,
													Optional: true,
													ForceNew: false,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												kubeClusterCustomizationDisabledKey: {
													Type:     schema.TypeList,
													Computed: true,
													Optional: true,
													ForceNew: false,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			kubeClusterCustomizationKubeProxyKey: {
				Type:     schema.TypeSet,
				Computed: false,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeClusterCustomizationIptablesKey: {
							Type:     schema.TypeSet,
							Computed: false,
							Optional: true,
							ForceNew: false,
							MaxItems: 1,
							Set:      CustomIPVSIPTablesSchemaSetFunc(),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									kubeClusterCustomizationMinSyncPeriodKey: {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
									kubeClusterCustomizationSyncPeriodKey: {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
								},
							},
						},
						kubeClusterCustomizationIpvsKey: {
							Type:     schema.TypeSet,
							Computed: false,
							Optional: true,
							ForceNew: false,
							MaxItems: 1,
							Set:      CustomIPVSIPTablesSchemaSetFunc(),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									kubeClusterCustomizationMinSyncPeriodKey: {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
									kubeClusterCustomizationSyncPeriodKey: {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
									kubeClusterCustomizationSchedulerKey: {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateEnum([]string{"rr", "lc", "dh", "sh", "sed", "nq"}),
									},
									kubeClusterCustomizationTcpFinTimeoutKey: {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
									kubeClusterCustomizationTcpTimeoutKey: {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
									kubeClusterCustomizationUdpTimeoutKey: {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
								},
							},
						},
					},
				},
			},
			kubeClusterCustomizationCiliumKey: {
				Type:        schema.TypeSet,
				Description: "Allow the customization of the Cilium deployment",
				Computed:    true,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeClusterCiliumClusterMeshKey: {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Customize Cilium's cluster mesh feature",
							Optional:    true,
							MaxItems:    1,
							Set:         CustomSchemaSetFunc(),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									kubeClusterCustomizationEnabledKey: {
										Type:        schema.TypeBool,
										Description: "Enable or disable the Cilium's Cluster mesh feature",
										Computed:    true,
										Optional:    true,
									},
									kubeClusterCiliumApiServerKey: {
										Type:        schema.TypeSet,
										Description: "Define how the cluster mesh will be exposed",
										Computed:    true,
										Optional:    true,
										MaxItems:    1,
										Set:         CustomSchemaSetFunc(),
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												kubeClusterCiliumNodePortKey: {
													Type:        schema.TypeInt,
													Description: "If the ServiceType is \"NodePort\", define on which port the service will be exposed",
													Computed:    true,
													Optional:    true,
												},
												kubeClusterCiliumServiceTypeKey: {
													Type:        schema.TypeString,
													Description: "Define if the cluster mesh service is will be exposed by a K8s Service of type NodePort or LoadBalancer",
													Computed:    true,
													Optional:    true,
												},
											},
										},
									},
								},
							},
						},
						kubeClusterCiliumHubbleKey: {
							Type:        schema.TypeSet,
							Description: "Allow the customization of the Hubble deployment",
							Computed:    true,
							Optional:    true,
							MaxItems:    1,
							Set:         CustomSchemaSetFunc(),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									kubeClusterCustomizationEnabledKey: {
										Type:        schema.TypeBool,
										Description: "Enable or disable the Hubble deployment",
										Computed:    true,
										Optional:    true,
									},
									kubeClusterCiliumRelayKey: {
										Type:        schema.TypeSet,
										Description: "Allow the customization of the Relay deployment",
										Computed:    true,
										Optional:    true,
										MaxItems:    1,
										Set:         CustomSchemaSetFunc(),
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												kubeClusterCustomizationEnabledKey: {
													Type:        schema.TypeBool,
													Description: "Enable or disable the Relay deployment",
													Computed:    true,
													Optional:    true,
												},
											},
										},
									},
									kubeClusterCiliumUiKey: {
										Type:        schema.TypeSet,
										Computed:    true,
										Description: "Allow the customization of the Hubble's UI deployment",
										Optional:    true,
										MaxItems:    1,
										Set:         CustomSchemaSetFunc(),
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												kubeClusterCustomizationEnabledKey: {
													Type:        schema.TypeBool,
													Description: "Enable or disable the Hubble's UI deployment",
													Computed:    true,
													Optional:    true,
												},
												kubeClusterCiliumBackendResources: {
													Type:        schema.TypeSet,
													Description: "Allow the customization of the Hubble UI Backend",
													Computed:    true,
													Optional:    true,
													MaxItems:    1,
													Set:         CustomSchemaSetFunc(),
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															kubeClusterCiliumLimitsKey: {
																Type:        schema.TypeSet,
																Description: "Define the limits of the Hubble UI Backend",
																Computed:    true,
																Optional:    true,
																MaxItems:    1,
																Set:         CustomSchemaSetFunc(),
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		kubeClusterCiliumCpuKey: {
																			Type:     schema.TypeString,
																			Computed: true,
																			Optional: true,
																		},
																		kubeClusterCiliumMemoryKey: {
																			Type:     schema.TypeString,
																			Computed: true,
																			Optional: true,
																		},
																	},
																},
															},
															kubeClusterCiliumRequestsKey: {
																Type:        schema.TypeSet,
																Description: "Define the requests of the Hubble UI Backend",
																Computed:    true,
																Optional:    true,
																MaxItems:    1,
																Set:         CustomSchemaSetFunc(),
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		kubeClusterCiliumCpuKey: {
																			Type:     schema.TypeString,
																			Computed: true,
																			Optional: true,
																		},
																		kubeClusterCiliumMemoryKey: {
																			Type:     schema.TypeString,
																			Computed: true,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
												kubeClusterCiliumFrontendResources: {
													Type:        schema.TypeSet,
													Description: "Allow the customization of the Hubble UI Frontend",
													Computed:    true,
													Optional:    true,
													MaxItems:    1,
													Set:         CustomSchemaSetFunc(),
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															kubeClusterCiliumLimitsKey: {
																Type:        schema.TypeSet,
																Description: "Define the limits of the Hubble UI Frontend",
																Computed:    true,
																Optional:    true,
																MaxItems:    1,
																Set:         CustomSchemaSetFunc(),
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		kubeClusterCiliumCpuKey: {
																			Type:     schema.TypeString,
																			Computed: true,
																			Optional: true,
																		},
																		kubeClusterCiliumMemoryKey: {
																			Type:     schema.TypeString,
																			Computed: true,
																			Optional: true,
																		},
																	},
																},
															},
															kubeClusterCiliumRequestsKey: {
																Type:        schema.TypeSet,
																Description: "Define the requests of the Hubble UI Frontend",
																Computed:    true,
																Optional:    true,
																MaxItems:    1,
																Set:         CustomSchemaSetFunc(),
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		kubeClusterCiliumCpuKey: {
																			Type:     schema.TypeString,
																			Computed: true,
																			Optional: true,
																		},
																		kubeClusterCiliumMemoryKey: {
																			Type:     schema.TypeString,
																			Computed: true,
																			Optional: true,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			kubeClusterPrivateNetworkIDKey: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			kubeClusterProxyModeKey: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: helpers.ValidateEnum([]string{"iptables", "ipvs"}),
			},
			kubeClusterPrivateNetworkConfigurationKey: {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeClusterDefaultVrackGatewayKey: {
							Required:    true,
							Type:        schema.TypeString,
							Description: "If defined, all egress traffic will be routed towards this IP address, which should belong to the private network. Empty string means disabled.",
						},
						kubeClusterPrivateNetworkRoutingAsDefault: {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Defines whether routing should default to using the nodes' private interface, instead of their public interface. Default is false.",
						},
					},
				},
			},
			kubeClusterLoadBalancersSubnetIdKey: {
				Type:     schema.TypeString,
				Optional: true,
				// private_network_id is required when load_balancers_subnet_id is set
				RequiredWith: []string{kubeClusterPrivateNetworkIDKey},
			},
			kubeClusterNodesSubnetIdKey: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				// private_network_id is required when nodes_subnet_id is set
				RequiredWith: []string{kubeClusterPrivateNetworkIDKey},
			},
			kubeRegionKey: {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Computed
			kubeClusterControlPlaneIsUpToDateKey: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			kubeClusterIsUpToDateKey: {
				Type:     schema.TypeBool,
				Computed: true,
			},
			kubeClusterNextUpgradeVersionsKey: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			kubeClusterNodesUrlKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			kubeStatusKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			kubeClusterUpdatePolicyKey: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			kubeClusterUrlKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			kubeClusterKubeconfigKey: {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			kubeClusterKubeconfigAttributesKey: {
				Type:        schema.TypeList,
				Computed:    true,
				Sensitive:   true,
				Description: "The kubeconfig configuration file of the Kubernetes cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeClusterKubeconfigHostKey: {
							Type:     schema.TypeString,
							Computed: true,
						},
						kubeClusterKubeconfigClusterCaCertificateKey: {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						kubeClusterKubeconfigClientCertificateKey: {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
						kubeClusterKubeconfigClientKeyKey: {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
			kubeClusterIpAllocationPolicyKey: {
				Description: "IP Allocation policy for the MKS cluster",
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeSet,
				MaxItems:    1,
				Set:         CustomSchemaSetFunc(),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeClusterPodsIpv4CidrKey: {
							Type:        schema.TypeString,
							Description: "The Kubernetes cluster's pods CIDR",
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Default:     nil,
						},
						kubeClusterServicesIpv4CidrKey: {
							Type:        schema.TypeString,
							Description: "The Kubernetes cluster's services CIDR",
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Default:     nil,
						},
					},
				},
			},
		},
	}
}

// CustomIPVSIPTablesSchemaSetFunc is a custom schema.SchemaSetFunc for IPVS and IPTables
// block configuration.
//
// Even if setting in the API `PT0S`, it returns `P0D` which is exactly the same duration but
// induce issue when calculating hashset.
//
// Moreover, we cannot use DiffSuppressFunc because even if the diff is removed the hashset is still different.
//
// Using schema.StateFunc does not help because of internal terraform execution diff calculation
// order.
func CustomIPVSIPTablesSchemaSetFunc() schema.SchemaSetFunc {
	return func(i interface{}) int {
		for k, v := range i.(map[string]interface{}) {
			if v == "P0D" {
				i.(map[string]interface{})[k] = "PT0S"
			}
		}

		out := fmt.Sprintf("%#v", i)
		return schema.HashString(out)
	}
}

func CustomSchemaSetFunc() schema.SchemaSetFunc {
	return func(i interface{}) int {
		out := fmt.Sprintf("%#v", i)
		return schema.HashString(out)
	}
}

// CustomApiServerAdmissionPluginsSchemaSetFunc is a custom schema.SchemaSetFunc for api_server.admission_plugins
// It orders plugins by alphabetical order to avoid hashset diff
func CustomApiServerAdmissionPluginsSchemaSetFunc() schema.SchemaSetFunc {
	return func(i interface{}) int {
		enabled := i.(map[string]interface{})["enabled"].([]interface{})
		disabled := i.(map[string]interface{})["disabled"].([]interface{})

		orderSliceByAlphabeticalOrder := func(s []interface{}) {
			sort.Slice(s, func(i, j int) bool {
				return s[i].(string) < s[j].(string)
			})
		}

		orderSliceByAlphabeticalOrder(enabled)
		orderSliceByAlphabeticalOrder(disabled)

		i.(map[string]interface{})["enabled"] = enabled
		i.(map[string]interface{})["disabled"] = disabled

		out := fmt.Sprintf("%#v", i)
		return schema.HashString(out)
	}
}

func resourceCloudProjectKubeImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("import Id is not service_name/kubeid formatted")
	}
	serviceName := splitId[0]
	id := splitId[1]
	d.SetId(id)
	d.Set(kubeServiceNameKey, serviceName)

	// add kubeconfig in state
	if err := setKubeconfig(d, meta); err != nil {
		return nil, err
	}

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectKubeCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)

	params := new(CloudProjectKubeCreateOpts)
	params.FromResource(d)

	res := &CloudProjectKubeResponse{}

	log.Printf("[DEBUG] Will create kube: %s", params)
	endpoint := fmt.Sprintf("/cloud/project/%s/kube", serviceName)
	if err := config.OVHClient.Post(endpoint, params, res); err != nil {
		return fmt.Errorf("calling Post %s with params %s:\n\t %w", endpoint, params, err)
	}

	log.Printf("[DEBUG] Waiting for kube %s to be available", res.Id)
	endpoint = fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, res.Id)
	if err := helpers.WaitAvailable(config.OVHClient, endpoint, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	log.Printf("[DEBUG] Waiting for kube %s to be READY", res.Id)
	if err := waitForCloudProjectKubeReady(config.OVHClient, serviceName, res.Id, []string{"INSTALLING", "UNKNOWN"}, []string{"READY"}, d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("timeout while waiting kube %s to be READY: %w", res.Id, err)
	}

	log.Printf("[DEBUG] kube %s is READY", res.Id)
	d.SetId(res.Id)

	return resourceCloudProjectKubeRead(d, meta)
}

func resourceCloudProjectKubeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, d.Id())
	res := &CloudProjectKubeResponse{}

	log.Printf("[DEBUG] Will read kube %s from project: %s", d.Id(), serviceName)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}
	for k, v := range res.ToMap(d) {
		log.Printf("[DEBUG] Will set %s to %v", k, v)
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	if d.IsNewResource() || d.Get(kubeClusterKubeconfigKey) == "" || len(d.Get(kubeClusterKubeconfigAttributesKey).([]interface{})) == 0 {
		// add kubeconfig in state
		if err := setKubeconfig(d, meta); err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Read kube %+v", res)
	return nil
}

func resourceCloudProjectKubeDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, d.Id())

	log.Printf("[DEBUG] Will delete kube %s from project: %s", d.Id(), serviceName)
	err := config.OVHClient.Delete(endpoint, nil)
	if err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	log.Printf("[DEBUG] Waiting for kube %s to be DELETED", d.Id())
	err = waitForCloudProjectKubeDeleted(d, config.OVHClient, serviceName, d.Id())
	if err != nil {
		return fmt.Errorf("timeout while waiting kube %s to be DELETED: %w", d.Id(), err)
	}
	log.Printf("[DEBUG] kube %s is DELETED", d.Id())

	d.SetId("")

	return nil
}

func resourceCloudProjectKubeUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)

	// if customization has changed, update it
	if d.HasChange(kubeClusterCustomizationApiServerKey) || d.HasChange(kubeClusterCustomization) || d.HasChange(kubeClusterCustomizationKubeProxyKey) || d.HasChange(kubeClusterCustomizationCiliumKey) {
		customization := new(Customization)

		if d.HasChange(kubeClusterCustomizationKubeProxyKey) {
			customization.KubeProxy = loadKubeProxyCustomization(d.Get(kubeClusterCustomizationKubeProxyKey))
		}

		if d.HasChange(kubeClusterCustomizationApiServerKey) {
			_, apiServerCustomization := d.GetChange(kubeClusterCustomizationApiServerKey)
			customization.APIServer = loadApiServerCustomization(apiServerCustomization)
		}

		// deprecated api server customization
		if d.HasChange(kubeClusterCustomization) {
			_, oldApiServerCustomization := d.GetChange(kubeClusterCustomization)
			customization.APIServer = loadDeprecatedApiServerCustomization(oldApiServerCustomization)
		}

		if d.HasChange(kubeClusterCustomizationCiliumKey) {
			_, newCiliumCustomization := d.GetChange(kubeClusterCustomizationCiliumKey)
			customization.Cilium = loadCiliumCustomizationFromResource(newCiliumCustomization)
		}

		params := &CloudProjectKubeUpdateCustomizationOpts{
			APIServer: customization.APIServer,
			KubeProxy: customization.KubeProxy,
			Cilium:    customization.Cilium,
		}

		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/customization", serviceName, d.Id())
		if err := config.OVHClient.Put(endpoint, params, nil); err != nil {
			return err
		}

		log.Printf("[DEBUG] Waiting for kube %s to be READY", d.Id())
		if err := waitForCloudProjectKubeReady(config.OVHClient, serviceName, d.Id(), []string{"REDEPLOYING", "RESETTING"}, []string{"READY"}, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return fmt.Errorf("timeout while waiting kube %s to be READY: %w", d.Id(), err)
		}

		log.Printf("[DEBUG] kube %s is READY", d.Id())
	}

	if d.HasChange(kubeVersionKey) {
		oldValueI, newValueI := d.GetChange(kubeVersionKey)

		oldValue := oldValueI.(string)
		newValue := newValueI.(string)

		log.Printf("[DEBUG] cluster version change from %s to %s", oldValue, newValue)

		oldVersion, err := version.NewVersion(oldValueI.(string))
		if err != nil {
			return fmt.Errorf("version %s does not match a semver", oldValue)
		}
		newVersion, err := version.NewVersion(newValueI.(string))
		if err != nil {
			return fmt.Errorf("version %s does not match a semver", newValue)
		}

		oldVersionSegments := oldVersion.Segments()
		newVersionSegments := newVersion.Segments()

		if oldVersionSegments[0] != 1 || newVersionSegments[0] != 1 {
			return fmt.Errorf("the only supported major version is 1")
		}
		if len(oldVersionSegments) < 2 || len(newVersionSegments) < 2 {
			log.Printf("[DEBUG] old version segments: %#v new version segments: %#v", oldVersionSegments, newVersionSegments)
			return fmt.Errorf("the version should only specify the major and minor versions (e.g. \\\"1.20\\\")")
		}

		if newVersion.LessThan(oldVersion) {
			return fmt.Errorf("cannot downgrade cluster from %s to %s", oldValue, newValue)
		}

		if oldVersionSegments[1]+1 != newVersionSegments[1] {
			return fmt.Errorf("cannot upgrade cluster from %s to %s, only next minor version is authorized", oldValue, newValue)
		}

		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/update", serviceName, d.Id())
		err = config.OVHClient.Post(endpoint, CloudProjectKubeUpdateOpts{
			Strategy: "NEXT_MINOR",
		}, nil)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Waiting for kube %s to be READY", d.Id())
		err = waitForCloudProjectKubeReady(config.OVHClient, serviceName, d.Id(), []string{"UPDATING", "REDEPLOYING", "RESETTING"}, []string{"READY"}, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("timeout while waiting kube %s to be READY: %w", d.Id(), err)
		}
		log.Printf("[DEBUG] kube %s is READY", d.Id())
	}

	if d.HasChanges(kubeClusterPlanKey) {
		oldValue, newValue := d.GetChange(kubeClusterPlanKey)
		newPlan := newValue.(string)
		oldPlan := oldValue.(string)
		if oldPlan == "standard" {
			return fmt.Errorf("you cannot migrate from %s to %s", oldPlan, newPlan)
		}

		return fmt.Errorf("migrate from %s to %s is not available yet", oldPlan, newPlan)
	}

	if d.HasChange(kubeClusterUpdatePolicyKey) {
		_, newValue := d.GetChange(kubeClusterUpdatePolicyKey)
		value := newValue.(string)

		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/updatePolicy", serviceName, d.Id())
		err := config.OVHClient.Put(endpoint, CloudProjectKubeUpdatePolicyOpts{
			UpdatePolicy: value,
		}, nil)
		if err != nil {
			return err
		}
	}

	if d.HasChange(kubeClusterLoadBalancersSubnetIdKey) {
		_, newValue := d.GetChange(kubeClusterLoadBalancersSubnetIdKey)
		value := newValue.(string)
		id := d.Id()
		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/updateLoadBalancersSubnetId",
			url.PathEscape(serviceName),
			url.PathEscape(id))
		err := config.OVHClient.Put(endpoint, CloudProjectKubeUpdateLoadBalancersSubnetIdOpts{
			LoadBalancersSubnetId: value,
		}, nil)
		if err != nil {
			return err
		}
		err = waitForCloudProjectKubeReady(config.OVHClient, serviceName, d.Id(), []string{"REDEPLOYING", "RESETTING"}, []string{"READY"}, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("timeout while waiting kube %s to be READY: %w", d.Id(), err)
		}
	}

	if d.HasChange(kubeNameKey) {
		_, newValue := d.GetChange(kubeNameKey)
		value := newValue.(string)

		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, d.Id())
		err := config.OVHClient.Put(endpoint, CloudProjectKubePutOpts{
			Name: &value,
		}, nil)
		if err != nil {
			return err
		}
	}

	if d.HasChange(kubeClusterPrivateNetworkConfigurationKey) {
		_, newValue := d.GetChange(kubeClusterPrivateNetworkConfigurationKey)
		pncOutput := privateNetworkConfiguration{}

		pncSet := newValue.(*schema.Set).List()
		for _, pnc := range pncSet {
			mapping := pnc.(map[string]interface{})
			pncOutput.DefaultVrackGateway = mapping[kubeClusterDefaultVrackGatewayKey].(string)
			pncOutput.PrivateNetworkRoutingAsDefault = mapping[kubeClusterPrivateNetworkRoutingAsDefault].(bool)
		}

		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/privateNetworkConfiguration", serviceName, d.Id())
		err := config.OVHClient.Put(endpoint, CloudProjectKubeUpdatePNCOpts{
			DefaultVrackGateway:            pncOutput.DefaultVrackGateway,
			PrivateNetworkRoutingAsDefault: pncOutput.PrivateNetworkRoutingAsDefault,
		}, nil)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Waiting for kube %s to be READY", d.Id())
		err = waitForCloudProjectKubeReady(config.OVHClient, serviceName, d.Id(), []string{"REDEPLOYING", "RESETTING"}, []string{"READY"}, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("timeout while waiting kube %s to be READY: %w", d.Id(), err)
		}
		log.Printf("[DEBUG] kube %s is READY", d.Id())
	}

	return nil
}

func cloudProjectKubeExists(serviceName, id string, client *ovhwrap.Client) error {
	res := &CloudProjectKubeResponse{}

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, id)
	return client.Get(endpoint, res)
}

func waitForCloudProjectKubeReady(client *ovhwrap.Client, serviceName, kubeId string, pending []string, target []string, timeout time.Duration) error {
	return waitForCloudProjectKubeReadyWithDelay(client, serviceName, kubeId, pending, target, timeout, 5*time.Second)
}

func waitForCloudProjectKubeReadyWithDelay(client *ovhwrap.Client, serviceName, kubeId string, pending []string, target []string, timeout time.Duration, delay time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: pending,
		Target:  target,
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectKubeResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, kubeId)
			err := client.Get(endpoint, res)
			if err != nil {
				return res, "", err
			}

			return res, res.Status, nil
		},
		Timeout:    timeout,
		Delay:      delay,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectKubeDeleted(d *schema.ResourceData, client *ovhwrap.Client, serviceName, kubeId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudProjectKubeResponse{}
			endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, kubeId)
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
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func setKubeconfig(d *schema.ResourceData, meta interface{}) error {
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeConfig, err := getKubeconfig(meta.(*Config), serviceName, d.Id())
	if err != nil {
		return err
	}

	if len(kubeConfig.Clusters) == 0 || len(kubeConfig.Users) == 0 {
		return fmt.Errorf("kubeconfig is invalid")
	}

	// raw kubeconfig
	d.Set(kubeClusterKubeconfigKey, kubeConfig.Raw)

	// kubeconfig attributes
	kubeconf := map[string]interface{}{}
	kubeconf[kubeClusterKubeconfigHostKey] = kubeConfig.Clusters[0].Cluster.Server
	kubeconf[kubeClusterKubeconfigClusterCaCertificateKey] = kubeConfig.Clusters[0].Cluster.CertificateAuthorityData
	kubeconf[kubeClusterKubeconfigClientCertificateKey] = kubeConfig.Users[0].User.ClientCertificateData
	kubeconf[kubeClusterKubeconfigClientKeyKey] = kubeConfig.Users[0].User.ClientKeyData
	_ = d.Set(kubeClusterKubeconfigAttributesKey, []map[string]interface{}{kubeconf})

	return nil
}
