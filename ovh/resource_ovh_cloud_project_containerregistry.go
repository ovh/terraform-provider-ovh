package ovh

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"

	"github.com/ovh/go-ovh/ovh"
)

func resourceCloudProjectContainerRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectContainerRegistryCreate,
		Read:   resourceCloudProjectContainerRegistryRead,
		Delete: resourceCloudProjectContainerRegistryDelete,
		Update: resourceCloudProjectContainerRegistryUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceCloudProjectContainerRegistryImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Registry name",
				Required:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "Region of the registry.",
				Required:    true,
			},
			"plan_id": {
				Type:        schema.TypeString,
				Description: "Plan ID of the registry.",
				Optional:    true,
				Computed:    true,
			},

			// Computed
			"created_at": {
				Type:        schema.TypeString,
				Description: "Registry creation date",
				Computed:    true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Description: "Project ID of your registry",
				Computed:    true,
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "Current size of the registry (bytes)",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Registry status",
				Computed:    true,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Description: "Registry last update date",
				Computed:    true,
			},
			"url": {
				Type:        schema.TypeString,
				Description: "Access url of the registry",
				Computed:    true,
			},
			"version": {
				Type:        schema.TypeString,
				Description: "Version of your registry",
				Computed:    true,
			},
			"plan": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Plan of the registry",
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
	}
}

func resourceCloudProjectContainerRegistryImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not service_name/id id formatted")
	}
	serviceName := splitId[0]
	id := splitId[1]
	d.SetId(id)
	d.Set("service_name", serviceName)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceCloudProjectContainerRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will create cloud project registry for project: %s", serviceName)

	opts := (&CloudProjectContainerRegistryCreateOpts{}).FromResource(d)
	reg := &CloudProjectContainerRegistry{}

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry",
		url.PathEscape(serviceName),
	)

	if err := config.OVHClient.Post(endpoint, opts, reg); err != nil {
		return fmt.Errorf("Error calling post %s:\n\t %q", endpoint, err)
	}

	d.SetId(reg.Id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"INSTALLING"},
		Target:     []string{"READY"},
		Refresh:    waitForCloudProjectContainerRegistry(config.OVHClient, serviceName, d.Id()),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for registry (%s): %s", d.Id(), err)
	}

	return resourceCloudProjectContainerRegistryRead(d, meta)
}

func resourceCloudProjectContainerRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Will update cloud project registry for project: %s", serviceName)

	opts := (&CloudProjectContainerRegistryUpdateOpts{}).FromResource(d)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s",
		url.PathEscape(serviceName),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Error calling put %s:\n\t %q", endpoint, err)
	}

	if err := cloudProjectContainerRegistryPlanUpdate(d, meta); err != nil {
		return err
	}

	return resourceCloudProjectContainerRegistryRead(d, meta)
}

func resourceCloudProjectContainerRegistryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will read cloud project registry %s for project: %s", id, serviceName)

	reg := &CloudProjectContainerRegistry{}

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s",
		url.PathEscape(serviceName),
		url.PathEscape(id),
	)

	if err := config.OVHClient.Get(endpoint, reg); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	for k, v := range reg.ToMap() {
		if k != "id" {
			d.Set(k, v)
		}
	}

	// OVH API Bug: the api doesn't set the region attribute value.
	// As a temp workaround, if the API sets an empty string for the region attr
	// we override it by extracting the region from the URL
	if d.Get("region").(string) == "" {
		urlRegionRx := regexp.MustCompile(`^https://[[:alnum:]]+\.([[:alpha:]]+)[0-9]+\.container-registry\.ovh\.net`)
		matches := urlRegionRx.FindStringSubmatch(reg.Url)
		if len(matches) > 1 {
			d.Set("region", strings.ToUpper(matches[1]))
		}
	}

	return cloudProjectContainerRegistryPlanRead(d, meta)
}

func cloudProjectContainerRegistryPlanRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will read cloud project registry plan %s for project: %s", id, serviceName)

	plan := &CloudProjectCapabilitiesContainerRegistryPlan{}

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/plan",
		url.PathEscape(serviceName),
		url.PathEscape(id),
	)

	if err := config.OVHClient.Get(endpoint, plan); err != nil {
		return fmt.Errorf("Error calling get %s:\n\t %q", endpoint, err)
	}

	d.Set("plan", []interface{}{plan.ToMap()})
	d.Set("plan_id", plan.Id)

	return nil
}

func cloudProjectContainerRegistryPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	if !d.HasChange("plan_id") {
		log.Printf("[DEBUG] cloud project registry plan for project : %s hasnt changed.", serviceName)
		return nil
	}

	log.Printf("[DEBUG] Will update cloud project registry plan for project: %s", serviceName)

	opts := (&CloudProjectContainerRegistryPlanUpdateOpts{}).FromResource(d)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s/plan",
		url.PathEscape(serviceName),
		url.PathEscape(d.Id()),
	)

	if err := config.OVHClient.Put(endpoint, opts, nil); err != nil {
		return fmt.Errorf("Error calling put %s:\n\t %q", endpoint, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"UPDATING"},
		Target:     []string{"READY"},
		Refresh:    waitForCloudProjectContainerRegistry(config.OVHClient, serviceName, d.Id()),
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for registry (%s): %s", d.Id(), err)
	}

	return nil
}

func resourceCloudProjectContainerRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	id := d.Id()

	log.Printf("[DEBUG] Will delete cloud project registry %s for project: %s", id, serviceName)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry/%s",
		url.PathEscape(serviceName),
		url.PathEscape(id),
	)

	if err := config.OVHClient.Delete(endpoint, nil); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED", "deleted"},
		Refresh:    waitForCloudProjectContainerRegistry(config.OVHClient, serviceName, d.Id()),
		Timeout:    30 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Deleting container registry %s from project %s: %s", id, serviceName, err)
	}

	d.SetId("")

	return nil
}

func waitForCloudProjectContainerRegistry(c *ovh.Client, serviceName, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudProjectContainerRegistry{}
		endpoint := fmt.Sprintf(
			"/cloud/project/%s/containerRegistry/%s",
			url.PathEscape(serviceName),
			url.PathEscape(id),
		)
		err := c.Get(endpoint, r)
		if err != nil {
			if err.(*ovh.APIError).Code == 404 {
				log.Printf("[DEBUG] container registry id %s on project %s deleted", id, serviceName)
				return r, "deleted", nil
			} else {
				return r, "", err
			}
		}

		log.Printf("[DEBUG] Pending container registry: %s", r)
		return r, r.Status, nil
	}
}
