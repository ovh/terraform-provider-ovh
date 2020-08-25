package ovh

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			"vrack_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// Deprecated:
				// this default value based on env vars is kept for retro compatibility
				// but should be removed in a future release
				DefaultFunc: schema.EnvDefaultFunc("OVH_VRACK_ID", ""),
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// Deprecated:
				// this default value based on env vars is kept for retro compatibility
				// but should be removed in a future release
				DefaultFunc: schema.EnvDefaultFunc("OVH_PROJECT_ID", ""),
			},
		},
	}
}

func resourceVrackCloudProjectImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	givenId := d.Id()
	splitId := strings.SplitN(givenId, "/", 2)
	if len(splitId) != 2 {
		return nil, fmt.Errorf("Import Id is not VRACK_ID/PROJECT_ID formatted")
	}
	vrackId := splitId[0]
	projectId := splitId[1]
	d.SetId(fmt.Sprintf("vrack_%s-cloudproject_%s", vrackId, projectId))
	d.Set("vrack_id", vrackId)
	d.Set("project_id", projectId)

	results := make([]*schema.ResourceData, 1)
	results[0] = d
	return results, nil
}

func resourceVrackCloudProjectCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vrackId := d.Get("vrack_id").(string)
	opts := (&VrackCloudProjectCreateOpts{}).FromResource(d)
	task := &VrackTask{}

	endpoint := fmt.Sprintf("/vrack/%s/cloudProject", vrackId)

	if err := config.OVHClient.Post(endpoint, opts, task); err != nil {
		return fmt.Errorf("Error calling POST %s with opts %v:\n\t %q", endpoint, opts, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to attach cloud project %v: %s", vrackId, opts, err)
	}

	//set id
	d.SetId(fmt.Sprintf("vrack_%s-cloudproject_%s", vrackId, opts.Project))

	return resourceVrackCloudProjectRead(d, meta)
}

func resourceVrackCloudProjectRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vcp := &VrackCloudProject{}

	vrackId := d.Get("vrack_id").(string)
	projectId := d.Get("project_id").(string)

	endpoint := fmt.Sprintf("/vrack/%s/cloudProject/%s",
		url.PathEscape(vrackId),
		url.PathEscape(projectId),
	)

	err := config.OVHClient.Get(endpoint, vcp)
	if err != nil {
		return err
	}

	d.Set("vrack_id", vcp.Vrack)
	d.Set("project_id", vcp.Project)

	return nil
}

func resourceVrackCloudProjectDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	vrackId := d.Get("vrack_id").(string)
	projectId := d.Get("project_id").(string)

	task := &VrackTask{}
	endpoint := fmt.Sprintf("/vrack/%s/cloudProject/%s",
		url.PathEscape(vrackId),
		url.PathEscape(projectId),
	)

	if err := config.OVHClient.Delete(endpoint, task); err != nil {
		return fmt.Errorf("Error calling DELETE %s with %s/%s:\n\t %q", endpoint, vrackId, projectId, err)
	}

	if err := waitForVrackTask(task, config.OVHClient); err != nil {
		return fmt.Errorf("Error waiting for vrack (%s) to detach cloud project (%s): %s", vrackId, projectId, err)
	}

	d.SetId("")
	return nil
}
