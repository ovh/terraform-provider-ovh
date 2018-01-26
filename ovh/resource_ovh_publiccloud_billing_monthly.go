package ovh

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

func resourcePublicCloudBillingMonthly() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicCloudBillingMonthlyCreate,
		Read:   resourcePublicCloudBillingMonthlyRead,
		Delete: resourcePublicCloudBillingMonthlyDelete,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourcePublicCloudBillingMonthlyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	instanceId := d.Get("instance_id").(string)

	params := &BillingMonthlyOpts{Project: projectId, InstanceId: instanceId}
	r := PublicCloudInstanceDetail{}

	log.Printf("[DEBUG] Will update billing to monthly on instance %s -> PublicCloud %s", instanceId, params.Project)
	endpoint := fmt.Sprintf("/cloud/project/%s/instance/%s/activeMonthlyBilling", projectId, instanceId)

	err := config.OVHClient.Post(endpoint, params, &r)
	if err != nil {
		// Code 461: Monthly billing is already subsribed for this instance
		if err.(*ovh.APIError).Code == 461 {
			log.Printf("[DEBUG] The api return error code %d, message: %s, Instance %s ->  PublicCloud %s", err.(*ovh.APIError).Code, err.(*ovh.APIError).Message, instanceId, params.Project)
			log.Printf("[DEBUG] Creating the ressource since it already exists. (Inconsistent state)")
		} else {
			return fmt.Errorf("Error calling %s with params %s:\n\t %q", endpoint, params, err)
		}
	}

	log.Printf("[DEBUG] Switched instance to monthly billing Task id %s: Instance %s ->  PublicCloud %s", r.Id, instanceId, params.Project)

	//set id
	d.SetId(fmt.Sprintf("instance_%s-cloudproject_%s-billing_monthly", instanceId, params.Project))

	return nil
}

func resourcePublicCloudBillingMonthlyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	projectId := d.Get("project_id").(string)
	instanceId := d.Get("instance_id").(string)

	params := BillingMonthlyOpts{Project: projectId, InstanceId: instanceId}
	r := PublicCloudInstanceDetail{}

	endpoint := fmt.Sprintf("/cloud/project/%s/instance/%s", projectId, instanceId)

	err := config.OVHClient.Get(endpoint, &r)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Read billing interval for instance %s ->  PublicCloud %s", instanceId, params.Project)

	if r.MonthlyBilling == nil {
		d.SetId("")
		return nil
	}

	return nil
}

func resourcePublicCloudBillingMonthlyDelete(d *schema.ResourceData, meta interface{}) error {
	projectId := d.Get("project_id").(string)
	instanceId := d.Get("instance_id").(string)

	log.Printf("[DEBUG] Cannot actually reverse monthly billing back to hourly billing until the instance %s is recreated -> PublicCloud %s", instanceId, projectId)

	d.SetId("")
	return nil
}
