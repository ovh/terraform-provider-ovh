package ovh

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

func resourceCloudProjectInstanceActiveMonthlyBilling() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudProjectInstanceActiveMonthlyBillingCreate,
		Update: resourceCloudProjectInstanceActiveMonthlyBillingUpdate,
		Read:   resourceCloudProjectInstanceActiveMonthlyBillingRead,
		Delete: resourceCloudProjectInstanceActiveMonthlyBillingDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The internal name of your dedicated server",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Public Cloud instance ID",
			},
			"wait_activation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Wait for monthly billing activation",
			},

			// Computed
			"monthly_billing_since": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Monthly billing activated since",
			},

			"monthly_billing_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Monthly billing status",
			},
		},
	}
}

func waitForCloudProjectInstanceActiveMonthlyBillingDone(client *ovh.Client, serviceName string, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r := &CloudProjectInstanceActiveMonthlyBillingResponse{}
		endpoint := fmt.Sprintf("/cloud/project/%s/instance/%s", url.PathEscape(serviceName), url.PathEscape(instanceId))
		if err := client.Get(endpoint, r); err != nil {
			return r, "", err
		}

		log.Printf("[DEBUG] Pending active monthly billing: %s", r)
		return r, r.MonthlyBilling.Status, nil
	}
}

func resourceCloudProjectInstanceActiveMonthlyBillingCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	instanceId := d.Get("instance_id").(string)
	waitActivation := d.Get("wait_activation").(bool)

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/instance/%s/activeMonthlyBilling", url.PathEscape(serviceName), url.PathEscape(instanceId),
	)

	params := (&CloudProjectInstanceActiveMonthlyBillingCreateOpts{}).FromResource(d)

	r := &CloudProjectInstanceActiveMonthlyBillingResponse{}

	log.Printf("[DEBUG] Will install active monthly billing: %s", params)

	if err := config.OVHClient.Post(endpoint, params, r); err != nil {
		return fmt.Errorf("Error calling POST %s:\n\t %q", endpoint, err)
	}

	if waitActivation {
		log.Printf("[DEBUG] Waiting for active monthly billing %s:", r)

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"activationPending"},
			Target:     []string{"ok"},
			Refresh:    waitForCloudProjectInstanceActiveMonthlyBillingDone(config.OVHClient, serviceName, instanceId),
			Timeout:    45 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("waiting for active monthly billing (%s): %s", params, err)
		}
	}

	log.Printf("[DEBUG] Created active monthly billing %s", r)

	return CloudProjectInstanceActiveMonthlyBillingRead(d, meta)
}

func CloudProjectInstanceActiveMonthlyBillingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)
	instanceId := d.Get("instance_id").(string)

	r := &CloudProjectInstanceActiveMonthlyBillingResponse{}

	log.Printf("[DEBUG] Will read active monthly billing: %s", serviceName)

	endpoint := fmt.Sprintf("/cloud/project/%s/instance/%s", url.PathEscape(serviceName), url.PathEscape(instanceId))

	if err := config.OVHClient.Get(endpoint, r); err != nil {
		return helpers.CheckDeleted(d, err, endpoint)
	}

	if r.MonthlyBilling != nil {
		d.Set("monthly_billing_since", r.MonthlyBilling.Since)
		d.Set("monthly_billing_status", r.MonthlyBilling.Status)
	}

	log.Printf("[DEBUG] Read active monthly billing %s", r)
	return nil
}

func resourceCloudProjectInstanceActiveMonthlyBillingUpdate(d *schema.ResourceData, meta interface{}) error {
	// nothing to do on update
	return resourceCloudProjectInstanceActiveMonthlyBillingRead(d, meta)
}

func resourceCloudProjectInstanceActiveMonthlyBillingRead(d *schema.ResourceData, meta interface{}) error {
	return CloudProjectInstanceActiveMonthlyBillingRead(d, meta)
}

func resourceCloudProjectInstanceActiveMonthlyBillingDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to do on DELETE
	return nil
}
