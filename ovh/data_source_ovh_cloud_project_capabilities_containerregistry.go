package ovh

import (
	"fmt"
	"log"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers/hashcode"
)

func dataSourceCloudProjectCapabilitiesContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloudProjectCapabilitiesContainerRegistryRead,
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"result": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of container registry capability for a single region",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_name": {
							Type:        schema.TypeString,
							Description: "The region name",
							Computed:    true,
						},
						"plans": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Available plans in the region",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": {
										Type:        schema.TypeString,
										Description: "Plan code from catalog",
										Computed:    true,
									},
									"created_at": {
										Type:        schema.TypeString,
										Description: "Plan creation date",
										Computed:    true,
									},
									"features": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Features of the plan",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vulnerability": {
													Type:        schema.TypeBool,
													Description: "Vulnerability scanning",
													Computed:    true,
												},
											},
										},
									},
									"id": {
										Type:        schema.TypeString,
										Description: "Plan ID",
										Computed:    true,
									},
									"name": {
										Type:        schema.TypeString,
										Description: "Plan name",
										Computed:    true,
									},
									"registry_limits": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Container registry limits",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"image_storage": {
													Type:        schema.TypeInt,
													Description: "Docker image storage limits in bytes",

													Computed: true,
												},
												"parallel_request": {
													Type:        schema.TypeInt,
													Description: "Parallel requests on Docker image API (/v2 Docker registry API)",
													Computed:    true,
												},
											},
										},
									},
									"updated_at": {
										Type:        schema.TypeString,
										Description: "Plan last update date",
										Computed:    true,
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

func dataSourceCloudProjectCapabilitiesContainerRegistryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will read cloud project capabilities container registry for project: %s", serviceName)

	capregs := []CloudProjectCapabilitiesContainerRegistry{}

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/capabilities/containerRegistry",
		url.PathEscape(serviceName),
	)
	err := config.OVHClient.Get(endpoint, &capregs)
	if err != nil {
		return fmt.Errorf("Error calling GET %s:\n\t %q", endpoint, err)
	}

	mapcapregs := make([]map[string]interface{}, len(capregs))
	ids := make([]string, len(capregs))

	for i, capreg := range capregs {
		mapcapregs[i] = capreg.ToMap()
		for _, plan := range capreg.Plans {
			ids = append(ids, plan.Id)
		}
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(ids)

	d.SetId(hashcode.Strings(ids))
	d.Set("result", mapcapregs)

	return nil
}
