package ovh

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var validSavingsPlanFlavors = []string{
	"rancher", "rancher_standard", "rancher_ovhcloud_edition",
	"b3-8", "b3-16", "b3-32", "b3-64", "b3-128", "b3-256",
	"c3-4", "c3-8", "c3-16", "c3-32", "c3-64", "c3-128",
	"r3-16", "r3-32", "r3-64", "r3-128", "r3-256", "r3-512",
}

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
				DefaultFunc: schema.EnvDefaultFunc("OVH_CLOUD_PROJECT_SERVICE", nil),
			},
			"flavor": {
				Type:        schema.TypeString,
				Description: "Savings Plan flavor",
				ForceNew:    true,
				Required:    true,
				ValidateFunc: func(v interface{}, s string) ([]string, []error) {
					value := strings.ToLower(v.(string))
					if !slices.Contains(validSavingsPlanFlavors, value) {
						return nil, []error{fmt.Errorf("invalid flavor %q, valid values are %s", value, validSavingsPlanFlavors)}
					}
					return nil, nil
				},
			},
			"deployment_type": {
				Type:        schema.TypeString,
				Description: "Deployment type of the Savings Plan (1AZ / 3AZ)",
				Optional:    true,
				ValidateFunc: func(v any, s string) ([]string, []error) {
					value := strings.ToUpper(v.(string))
					if value != "1AZ" && value != "3AZ" {
						return nil, []error{fmt.Errorf("invalid deployment_type %q, valid values are 1AZ or 3AZ", value)}
					}
					return nil, nil
				},
				Default: "1AZ",
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

func resourceSavingsPlanImport(d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
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

func fetchSavingsPlanOffers(config *Config, d *schema.ResourceData, serviceID int) ([]savingsPlansSubscribable, error) {
	flavor := strings.ReplaceAll(d.Get("flavor").(string), "_", " ")
	if flavor == "rancher" {
		flavor = "rancher standard"
	}

	deploymentType := strings.ToUpper(d.Get("deployment_type").(string))
	if strings.HasPrefix(flavor, "rancher") && deploymentType != "1AZ" {
		return nil, fmt.Errorf("invalid deployment_type %q for flavor %q, only 1AZ is supported", deploymentType, flavor)

	}

	fullFlavor := flavor
	if deploymentType == "3AZ" {
		fullFlavor += " 3AZ"
	}

	endpoint := fmt.Sprintf("/services/%d/savingsPlans/subscribable?productCode=%q", serviceID, url.QueryEscape(fullFlavor))
	subscribables := []savingsPlansSubscribable{}
	if err := config.OVHClient.Get(endpoint, &subscribables); err != nil {
		return nil, fmt.Errorf("error calling GET %s:\n\t %q", endpoint, err)
	}

	if strings.HasPrefix(flavor, "rancher") || deploymentType == "3AZ" {
		return subscribables, nil
	}

	// Fetch 3AZ flavors to be able to find the 1AZ ones
	threeAZPlans := []savingsPlansSubscribable{}
	endpoint = fmt.Sprintf("/services/%d/savingsPlans/subscribable?productCode=%q", serviceID, url.QueryEscape(flavor+" 3AZ"))
	if err := config.OVHClient.Get(endpoint, &threeAZPlans); err != nil {
		return nil, fmt.Errorf("error calling GET %s:\n\t %q", endpoint, err)
	}

	// Extract offer IDs of 3AZ plans
	threeAZOfferIDs := make(map[string]struct{})
	for _, plan := range threeAZPlans {
		threeAZOfferIDs[plan.OfferID] = struct{}{}
	}

	// Filter out 3AZ plans to keep only 1AZ ones
	var oneAZPlans []savingsPlansSubscribable
	for _, plan := range subscribables {
		if _, ok := threeAZOfferIDs[plan.OfferID]; !ok {
			oneAZPlans = append(oneAZPlans, plan)
		}
	}

	return oneAZPlans, nil
}

func resourceSavingsPlanCreate(d *schema.ResourceData, meta any) error {
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
	subscribables, err := fetchSavingsPlanOffers(config, d, serviceId)
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Retrieved %d subscribable savings plans", len(subscribables))

	// Search for a savings plan corresponding to the given parameters
	endpoint := fmt.Sprintf("/services/%d/savingsPlans/subscribe/simulate", serviceId)
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

		if d.Get("period").(string) == resp.Period &&
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

func resourceSavingsPlanRead(d *schema.ResourceData, meta any) error {
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

func resourceSavingsPlanUpdate(d *schema.ResourceData, meta any) error {
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

func resourceSavingsPlanDelete(d *schema.ResourceData, meta any) error {
	// Does nothing, savings plans cannot be deleted
	d.SetId("")
	return nil
}
