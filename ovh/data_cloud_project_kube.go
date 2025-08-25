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
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"kube_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plan": {
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
				// Required: true,
				ForceNew: false,
				// MaxItems: 1,
				Set: CustomSchemaSetFunc(),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admissionplugins": {
							Type:     schema.TypeSet,
							Computed: true,
							Optional: true,
							// Required: true,
							ForceNew: false,
							// MaxItems: 1,
							Set: CustomSchemaSetFunc(),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeList,
										Computed: true,
										Optional: true,
										// Required: true,
										ForceNew: false,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"disabled": {
										Type:     schema.TypeList,
										Computed: true,
										Optional: true,
										// Required: true,
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
						"apiserver": {
							Type:       schema.TypeSet,
							Computed:   true,
							Optional:   true,
							ForceNew:   false,
							Set:        CustomSchemaSetFunc(),
							Deprecated: fmt.Sprintf("Use %s instead", kubeClusterCustomizationApiServerKey),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"admissionplugins": {
										Type:     schema.TypeSet,
										Computed: true,
										Optional: true,
										ForceNew: false,
										Set:      CustomSchemaSetFunc(),
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:     schema.TypeList,
													Computed: true,
													Optional: true,
													ForceNew: false,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"disabled": {
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
						"iptables": {
							Type:     schema.TypeSet,
							Computed: false,
							Optional: true,
							ForceNew: false,
							MaxItems: 1,
							Set:      CustomIPVSIPTablesSchemaSetFunc(),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_sync_period": {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
									"sync_period": {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
								},
							},
						},
						"ipvs": {
							Type:     schema.TypeSet,
							Computed: false,
							Optional: true,
							ForceNew: false,
							MaxItems: 1,
							Set:      CustomIPVSIPTablesSchemaSetFunc(),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_sync_period": {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
									"sync_period": {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
									"scheduler": {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateEnum([]string{"rr", "lc", "dh", "sh", "sed", "nq"}),
									},
									"tcp_fin_timeout": {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
									"tcp_timeout": {
										Type:         schema.TypeString,
										Computed:     false,
										Optional:     true,
										ForceNew:     false,
										ValidateFunc: helpers.ValidateRFC3339Duration,
									},
									"udp_timeout": {
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

			"private_network_id": {
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
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"url": {
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
			"kubeconfig": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"kubeconfig_attributes": {
				Type:        schema.TypeList,
				Computed:    true,
				Sensitive:   true,
				Description: "The kubeconfig configuration file of the Kubernetes cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_ca_certificate": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
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
					},
				},
			},
		},
	}
}

func dataSourceCloudProjectKubeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	kubeId := d.Get("kube_id").(string)

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
	serviceName := d.Get("service_name").(string)
	var kubeId string

	// For data source, use kube_id instead of d.Id()
	if id := d.Get("kube_id"); id != nil {
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
	d.Set("kubeconfig", kubeConfig.Raw)

	// kubeconfig attributes
	kubeconf := map[string]interface{}{}
	kubeconf["host"] = kubeConfig.Clusters[0].Cluster.Server
	kubeconf["cluster_ca_certificate"] = kubeConfig.Clusters[0].Cluster.CertificateAuthorityData
	kubeconf["client_certificate"] = kubeConfig.Users[0].User.ClientCertificateData
	kubeconf["client_key"] = kubeConfig.Users[0].User.ClientKeyData
	d.Set("kubeconfig_attributes", []map[string]interface{}{kubeconf})

	return nil
}
