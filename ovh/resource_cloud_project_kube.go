package ovh

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

const (
	kubeClusterLoadBalancersSubnetIdKey       = "load_balancers_subnet_id"
	kubeClusterNodesSubnetIdKey               = "nodes_subnet_id"
	kubeClusterNameKey                        = "name"
	kubeClusterPrivateNetworkIDKey            = "private_network_id"
	kubeClusterPrivateNetworkConfigurationKey = "private_network_configuration"
	kubeClusterUpdatePolicyKey                = "update_policy"
	kubeClusterVersionKey                     = "version"

	kubeClusterProxyModeKey = "kube_proxy_mode"

	kubeClusterCustomization             = "customization" // Deprecated
	kubeClusterCustomizationApiServerKey = "customization_apiserver"
	kubeClusterCustomizationKubeProxyKey = "customization_kube_proxy"
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
			Update:  schema.DefaultTimeout(10 * time.Minute),
			Delete:  schema.DefaultTimeout(10 * time.Minute),
			Read:    schema.DefaultTimeout(5 * time.Minute),
			Default: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			kubeClusterNameKey: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			kubeClusterVersionKey: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: false,
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
						"admissionplugins": {
							Type:     schema.TypeSet,
							Computed: true,
							Optional: true,
							ForceNew: false,
							Set:      CustomApiServerAdmissionPluginsSchemaSetFunc(),
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
										Set:      CustomApiServerAdmissionPluginsSchemaSetFunc(),
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

			kubeClusterPrivateNetworkIDKey: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// private_network_id is required when load_balancers_subnet_id is set
				RequiredWith: []string{kubeClusterLoadBalancersSubnetIdKey},
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
						"default_vrack_gateway": {
							Required:    true,
							Type:        schema.TypeString,
							Description: "If defined, all egress traffic will be routed towards this IP address, which should belong to the private network. Empty string means disabled.",
						},
						"private_network_routing_as_default": {
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
			},
			kubeClusterNodesSubnetIdKey: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// Computed
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
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			kubeClusterUpdatePolicyKey: {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"url": {
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
	d.Set("service_name", serviceName)

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
	serviceName := d.Get("service_name").(string)

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
	if err := waitForCloudProjectKubeReady(config.OVHClient, serviceName, res.Id, []string{"INSTALLING"}, []string{"READY"}, d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("timeout while waiting kube %s to be READY: %w", res.Id, err)
	}

	log.Printf("[DEBUG] kube %s is READY", res.Id)
	d.SetId(res.Id)

	return resourceCloudProjectKubeRead(d, meta)
}

func resourceCloudProjectKubeRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

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

	if d.IsNewResource() || d.Get("kubeconfig") == "" || len(d.Get("kubeconfig_attributes").([]interface{})) == 0 {
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
	serviceName := d.Get("service_name").(string)

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
	serviceName := d.Get("service_name").(string)

	// if customization has changed, update it
	if d.HasChange(kubeClusterCustomizationApiServerKey) || d.HasChange(kubeClusterCustomization) || d.HasChange(kubeClusterCustomizationKubeProxyKey) {
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

		params := &CloudProjectKubeUpdateCustomizationOpts{
			APIServer: customization.APIServer,
			KubeProxy: customization.KubeProxy,
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

	if d.HasChange(kubeClusterVersionKey) {
		oldValueI, newValueI := d.GetChange(kubeClusterVersionKey)

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

		endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s/updateLoadBalancersSubnetId", serviceName, d.Id())
		err := config.OVHClient.Put(endpoint, CloudProjectKubeUpdateLoadBalancersSubnetIdOpts{
			LoadBalancersSubnetId: value,
		}, nil)
		if err != nil {
			return err
		}
	}

	if d.HasChange(kubeClusterNameKey) {
		_, newValue := d.GetChange(kubeClusterNameKey)
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
			pncOutput.DefaultVrackGateway = mapping["default_vrack_gateway"].(string)
			pncOutput.PrivateNetworkRoutingAsDefault = mapping["private_network_routing_as_default"].(bool)
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

func cloudProjectKubeExists(serviceName, id string, client *ovh.Client) error {
	res := &CloudProjectKubeResponse{}

	endpoint := fmt.Sprintf("/cloud/project/%s/kube/%s", serviceName, id)
	return client.Get(endpoint, res)
}

func waitForCloudProjectKubeReady(client *ovh.Client, serviceName, kubeId string, pending []string, target []string, timeout time.Duration) error {
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
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

func waitForCloudProjectKubeDeleted(d *schema.ResourceData, client *ovh.Client, serviceName, kubeId string) error {
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
	serviceName := d.Get("service_name").(string)
	kubeConfig, err := getKubeconfig(meta.(*Config), serviceName, d.Id())
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
	_ = d.Set("kubeconfig_attributes", []map[string]interface{}{kubeconf})

	return nil
}
