package ovh

import (
	"fmt"
	"log"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			kubeClusterProxyModeKey: {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			kubeClusterCustomizationApiServerKey: {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				// Required: true,
				ForceNew: false,
				// MaxItems: 1,
				Set: CustomSchemaSetFunc(false),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admissionplugins": {
							Type:     schema.TypeSet,
							Computed: true,
							Optional: true,
							// Required: true,
							ForceNew: false,
							// MaxItems: 1,
							Set: CustomSchemaSetFunc(false),
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

			kubeClusterCustomizationKubeProxyKey: {
				Type:     schema.TypeSet,
				Computed: false,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				// Set:      CustomSchemaSetFunc(true),
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"iptables": {
							Type:     schema.TypeSet,
							Computed: false,
							Optional: true,
							ForceNew: false,
							MaxItems: 1,
							Set:      CustomSchemaSetFunc(false),

							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_sync_period": {
										Type:     schema.TypeString,
										Computed: false,
										Optional: true,
										ForceNew: false,
									},
									"sync_period": {
										Type:     schema.TypeString,
										Computed: false,
										Optional: true,
										ForceNew: false,
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
							Set:      CustomSchemaSetFunc(false),
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_sync_period": {
										Type:     schema.TypeString,
										Computed: false,
										Optional: true,
										ForceNew: false,
									},
									"sync_period": {
										Type:     schema.TypeString,
										Computed: false,
										Optional: true,
										ForceNew: false,
									},
									"scheduler": {
										Type:     schema.TypeString,
										Computed: false,
										Optional: true,
										ForceNew: false,
									},
									"tcp_fin_timeout": {
										Type:     schema.TypeString,
										Computed: false,
										Optional: true,
										ForceNew: false,
									},
									"tcp_timeout": {
										Type:     schema.TypeString,
										Computed: false,
										Optional: true,
										ForceNew: false,
									},
									"udp_timeout": {
										Type:     schema.TypeString,
										Computed: false,
										Optional: true,
										ForceNew: false,
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
	err := config.OVHClient.Get(endpoint, res)
	if err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	for k, v := range res.ToMap() {
		if k != "id" {
			d.Set(k, v)
		} else {
			d.SetId(fmt.Sprint(v))
		}
	}

	return nil
}
