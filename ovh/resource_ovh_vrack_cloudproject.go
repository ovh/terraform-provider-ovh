package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceVrackCloudProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceVrackCloudProjectCreate,
		Read:   resourceVrackCloudProjectRead,
		Delete: resourceVrackCloudProjectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVrackCloudProjectImportState,
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_VRACK_SERVICE", nil),
				Description: "Service name of the resource representing the id of the cloud project.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", ""),
			},
		},
	}
}

func resourceVrackCloudProjectImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not SERVICE_NAME/PROJECT_ID formatted")
	}
	serviceName := splitId[0]
	projectId := splitId[1]
	d.SetId(fmt.Sprintf("vrack_%s-cloudproject_%s", serviceName, projectId))
	d.Set("service_name", serviceName)
	d.Set("project_id", projectId)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVrackCloudProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)

	opts := (&VrackCloudProjectCreateOpts{}).FromResource(d)
	task := &VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/cloudProject", serviceName)

	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to attach cloud project %v: %s", serviceName, opts, err)
	}

	//set id
	d.SetId(fmt.Sprintf("vrack_%s-cloudproject_%s", serviceName, opts.Project))

	return resourceVrackCloudProjectRead(d, meta)
}

func resourceVrackCloudProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vcp := &VrackCloudProject{}
	serviceName := d.Get("service_name").(string)
	projectId := d.Get("project_id").(string)

	endpoint := fmt.Sprintf("/vrack/%s/cloudProject/%s",
		url.PathEscape(serviceName),
		url.PathEscape(projectId),
	)

	if err := config.OVHClient.Get(endpoint, vcp); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	d.Set("service_name", vcp.Vrack)
	d.Set("project_id", vcp.Project)

	return nil
}

func resourceVrackCloudProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	serviceName := d.Get("service_name").(string)
	projectId := d.Get("project_id").(string)

	task := &VrackTask{}
	endpoint := fmt.Sprintf("/vrack/%s/cloudProject/%s",
		url.PathEscape(serviceName),
		url.PathEscape(projectId),
	)

	if err := config.OVHClient.Delete(endpoint, task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, serviceName, projectId, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach cloud project (%s): %s", serviceName, projectId, err)
	}

	d.SetId("")
	return nil
}
