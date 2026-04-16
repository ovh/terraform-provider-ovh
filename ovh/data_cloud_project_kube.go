package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
)

func dataSourceCloudProjectKube() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectKubeRead,
		Schema: map[string]*schema.Schema{
			kubeServiceNameKey: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			kubeKubeIdKey: {
				Type:     schema.TypeString,
				Required: true,
			},
			kubeNameKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			kubeVersionKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			kubeClusterPlanKey: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     false,
				Default:      "free",
				ValidateFunc: helpers.ValidateEnum([]string{"standard", "free"}),
			},
			kubeClusterProxyModeKey: {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: helpers.ValidateEnum([]string{"iptables", "ipvs"}),
			},
			kubeClusterCustomizationApiServerKey: {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				ForceNew: false,
				Set:      CustomSchemaSetFunc(),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						kubeClusterCustomizationAdmissionPluginsKey: {
							Type:     schema.TypeSet,
							Computed: true,
							Optional: true,
							ForceNew: false,
							Set:      CustomSchemaSetFunc(),
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
				Type:       schema.TypeSet,
				Computed:   true,
				Optional:   true,
				ForceNew:   false,
				Set:        CustomSchemaSetFunc(),
				Deprecated: fmt.Sprintf("Use %s instead", kubeClusterCustomizationApiServerKey),
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
										Set:      CustomSchemaSetFunc(),
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
				Computed: true,
			},
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
			kubeRegionKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			kubeStatusKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			kubeClusterUpdatePolicyKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			kubeClusterUrlKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			kubeClusterLoadBalancersSubnetIdKey: {
				Type:     schema.TypeString,
				Computed: true,
			},
			kubeClusterNodesSubnetIdKey: {
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
							Default:     nil,
						},
						kubeClusterServicesIpv4CidrKey: {
							Type:        schema.TypeString,
							Description: "The Kubernetes cluster's services CIDR",
							Optional:    true,
							Computed:    true,
							Default:     nil,
						},
					},
				},
			},
		},
	}
}

func dataSourceCloudProjectKubeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get(kubeServiceNameKey).(string)
	kubeId := d.Get(kubeKubeIdKey).(string)

	log.Printf("[DEBUG] Will read public cloud kube %s for project: %s", kubeId, serviceName)

	res := &CloudProjectKubeResponse{}
	endpoint := fmt.Sprintf(
		"/cloud/project/%s/kube/%s",
		url.PathEscape(serviceName),
		url.PathEscape(kubeId),
	)
	if err := config.OVHClient.Get(endpoint, res); err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	for k, v := range res.ToMap(d) {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	// add kubeconfig in state
	if err := dataSourceKubeconfig(d, meta); err != nil {
		return err
	}

	return nil
}

func dataSourceKubeconfig(d *schema.ResourceData, meta interface{}) error {
	serviceName := d.Get(kubeServiceNameKey).(string)
	var kubeId string

	// For data source, use kube_id instead of d.Id()
	if id := d.Get(kubeKubeIdKey); id != nil {
		kubeId = id.(string)
	} else {
		kubeId = d.Id()
	}

	kubeConfig, err := getKubeconfig(meta.(*Config), serviceName, kubeId)
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
	d.Set(kubeClusterKubeconfigAttributesKey, []map[string]interface{}{kubeconf})

	return nil
}
