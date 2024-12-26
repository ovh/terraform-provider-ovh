package ovh

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSavingsPlan() *schema.Resource {
	return &schema.Resource{
		Create: resourceSavingsPlanCreate,
		Read:   resourceSavingsPlanRead,
		Update: resourceSavingsPlanUpdate,
		Delete: resourceSavingsPlanDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSavingsPlanImport,
		},
		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:        schema.TypeString,
				Description: "ID of the public cloud project",
				ForceNew:    true,
				Required:    true,
			},
			"flavor": {
				Type:        schema.TypeString,
				Description: "Savings Plan flavor (e.g. Rancher, C3-4, any instance flavor, ...)",
				ForceNew:    true,
				Required:    true,
			},
			"period": {
				Type:        schema.TypeString,
				Description: "Periodicity of the Savings Plan",
				ForceNew:    true,
				Required:    true,
			},
			"size": {
				Type:        schema.TypeInt,
				Description: "Size of the Savings Plan",
				Required:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "Custom display name, used in invoices",
				Required:    true,
			},
			"auto_renewal": {
				Type:        schema.TypeBool,
				Description: "Whether Savings Plan should be renewed at the end of the period (defaults to false)",
				Optional:    true,
				Computed:    true,
			},

			// computed
			"service_id": {
				Type:        schema.TypeInt,
				Description: "ID of the service",
				Computed:    true,
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the Savings Plan",
				Computed:    true,
			},
			"start_date": {
				Type:        schema.TypeString,
				Description: "Start date of the Savings Plan",
				Computed:    true,
			},
			"end_date": {
				Type:        schema.TypeString,
				Description: "End date of the Savings Plan",
				Computed:    true,
			},
			"period_end_action": {
				Type:        schema.TypeString,
				Description: "Action performed when reaching the end of the period",
				Computed:    true,
			},
			"period_start_date": {
				Type:        schema.TypeString,
				Description: "Start date of the current period",
				Computed:    true,
			},
			"period_end_date": {
				Type:        schema.TypeString,
				Description: "End date of the current period",
				Computed:    true,
			},
		},
	}
}

func resourceSavingsPlanImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	importID := d.Id()
	parts := strings.Split(importID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("import ID is not correctly formatted, expected 'serviceName/savingsPlanID'")
	}

	serviceName := parts[0]
	savingsPlanID := parts[1]

	// Retrieve service ID
	serviceId, err := serviceIdFromResourceName(config.OVHClient, serviceName)
	if err != nil {
		return nil, err
	}
	d.Set("service_id", serviceId)
	d.SetId(savingsPlanID)

	return []*schema.ResourceData{d}, nil
}

func resourceSavingsPlanCreate(d *schema.ResourceData, meta interface{}) error {
	serviceName := d.Get("service_name").(string)
	config := meta.(*Config)

	// Retrieve service ID
	serviceId, err := serviceIdFromResourceName(config.OVHClient, serviceName)
	if err != nil {
		return err
	}
	d.Set("service_id", serviceId)

	// Get subscribables savings plans
	log.Print("[DEBUG] Will fetch subscribables savings plans")
	endpoint := fmt.Sprintf("/services/%d/savingsPlans/subscribable", serviceId)
	subscribables := []savingsPlansSubscribable{}
	if err := config.OVHClient.Get(endpoint, &subscribables); err != nil {
		return fmt.Errorf("error calling GET %s:\n\t %q", endpoint, err)
	}

	// Search for a savings plan corresponding to the given parameters
	endpoint = fmt.Sprintf("/services/%d/savingsPlans/subscribe/simulate", serviceId)
	for _, subscribable := range subscribables {
		var (
			req = savingsPlansSimulateRequest{
				DisplayName: d.Get("display_name").(string),
				OfferID:     subscribable.OfferID,
				Size:        d.Get("size").(int),
			}
			resp savingsPlansSimulateResponse
		)

		if err := config.OVHClient.Post(endpoint, req, &resp); err != nil {
			return fmt.Errorf("error calling POST %s:\n\t %q", endpoint, err)
		}

		if d.Get("flavor").(string) == resp.Flavor &&
			d.Get("period").(string) == resp.Period &&
			d.Get("size").(int) == resp.Size {
			// We found the right savings plan, execute subscription
			endpoint = fmt.Sprintf("/services/%d/savingsPlans/subscribe/execute", serviceId)
			if err := config.OVHClient.Post(endpoint, req, &resp); err != nil {
				return fmt.Errorf("error calling POST %s:\n\t %q", endpoint, err)
			}

			// Then update the action at end of period if renewal was asked
			autoRenewalConfig := d.Get("auto_renewal").(bool)
			if autoRenewalConfig {
				endpoint = fmt.Sprintf("/services/%d/savingsPlans/subscribed/%s/changePeriodEndAction", serviceId, url.PathEscape(resp.ID))
				if err := config.OVHClient.Post(endpoint, savingsPlanPeriodEndActionRequest{
					PeriodEndAction: "REACTIVATE",
				}, nil); err != nil {
					return fmt.Errorf("error calling POST %s:\n\t %q", endpoint, err)
				}
			}

			d.SetId(resp.ID)
			d.Set("status", resp.Status)
			d.Set("start_date", resp.StartDate)
			d.Set("end_date", resp.EndDate)
			d.Set("period_end_action", resp.PeriodEndAction)
			d.Set("period_start_date", resp.PeriodStartDate)
			d.Set("period_end_date", resp.PeriodEndDate)

			return nil
		}
	}

	return errors.New("no savings plan available with the given parameters")
}

func resourceSavingsPlanRead(d *schema.ResourceData, meta interface{}) error {
	serviceID := d.Get("service_id").(int)
	config := meta.(*Config)

	endpoint := fmt.Sprintf("/services/%d/savingsPlans/subscribed/%s", serviceID, url.PathEscape(d.Id()))
	var resp savingsPlansSimulateResponse
	if err := config.OVHClient.Get(endpoint, &resp); err != nil {
		return fmt.Errorf("error calling GET %s:\n\t %q", endpoint, err)
	}

	d.Set("status", resp.Status)
	d.Set("start_date", resp.StartDate)
	d.Set("end_date", resp.EndDate)
	d.Set("period_end_action", resp.PeriodEndAction)
	d.Set("period_start_date", resp.PeriodStartDate)
	d.Set("period_end_date", resp.PeriodEndDate)
	d.Set("auto_renewal", resp.PeriodEndAction == "REACTIVATE")

	return nil
}

func resourceSavingsPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	serviceID := d.Get("service_id").(int)
	config := meta.(*Config)

	// Update display name if needed
	if d.HasChange("display_name") {
		endpoint := fmt.Sprintf("/services/%d/savingsPlans/subscribed/%s", serviceID, url.PathEscape(d.Id()))
		if err := config.OVHClient.Put(endpoint, map[string]string{
			"displayName": d.Get("display_name").(string),
		}, nil); err != nil {
			return fmt.Errorf("error calling PUT %s:\n\t %q", endpoint, err)
		}
	}

	// Update size if needed
	if d.HasChange("size") {
		endpoint := fmt.Sprintf("/services/%d/savingsPlans/subscribed/%s/changeSize", serviceID, url.PathEscape(d.Id()))
		if err := config.OVHClient.Post(endpoint, map[string]int{
			"size": d.Get("size").(int),
		}, nil); err != nil {
			return fmt.Errorf("error calling POST %s:\n\t %q", endpoint, err)
		}
	}

	// Update auto renewal if needed
	if d.HasChange("auto_renewal") {
		newValue := d.Get("auto_renewal").(bool)
		endAction := "TERMINATE"
		if newValue {
			endAction = "REACTIVATE"
		}

		endpoint := fmt.Sprintf("/services/%d/savingsPlans/subscribed/%s/changePeriodEndAction", serviceID, url.PathEscape(d.Id()))
		if err := config.OVHClient.Post(endpoint, map[string]string{
			"periodEndAction": endAction,
		}, nil); err != nil {
			return fmt.Errorf("error calling POST %s:\n\t %q", endpoint, err)
		}
	}

	return nil
}

func resourceSavingsPlanDelete(d *schema.ResourceData, meta interface{}) error {
	// Does nothing, savings plans cannot be deleted
	d.SetId("")
	return nil
}
