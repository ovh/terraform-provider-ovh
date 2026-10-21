package ovh

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestNormalizeSavingsPlanFlavor(t *testing.T) {
	tests := []struct {
		flavor         string
		wantFlavor     string
		wantDeployment string
	}{
		{"b3-8", "b3-8", "1AZ"},
		{"B3-8", "b3-8", "1AZ"},
		{"B3-8 3AZ", "b3-8", "3AZ"},
		{"R3-16", "r3-16", "1AZ"},
		{"rancher", "rancher_standard", "1AZ"},
		{"Rancher", "rancher_standard", "1AZ"},
		{"rancher_standard", "rancher_standard", "1AZ"},
		{"Rancher Standard", "rancher_standard", "1AZ"},
		{"rancher_ovhcloud_edition", "rancher_ovhcloud_edition", "1AZ"},
		{"Rancher OVHcloud Edition", "rancher_ovhcloud_edition", "1AZ"},
	}

	for _, tt := range tests {
		flavor, deploymentType := normalizeSavingsPlanFlavor(tt.flavor)
		if flavor != tt.wantFlavor || deploymentType != tt.wantDeployment {
			t.Errorf("normalizeSavingsPlanFlavor(%q) = (%q, %q), want (%q, %q)",
				tt.flavor, flavor, deploymentType, tt.wantFlavor, tt.wantDeployment)
		}
	}
}

func TestEqualSavingsPlanPeriods(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
	}{
		{"P1Y", "P1Y", true},
		{"P1Y", "P12M", true},
		{"P36M", "P3Y", true},
		{"p3y", "P3Y", true},
		{"P1Y", "P3Y", false},
		{"P1M", "P1Y", false},
		{"P1Y6M", "P18M", true},
		{"", "P1Y", false},
		{"P", "P", true},
		{"P", "P0M", false},
		{"1Y", "P1Y", false},
	}

	for _, tt := range tests {
		if got := equalSavingsPlanPeriods(tt.a, tt.b); got != tt.want {
			t.Errorf("equalSavingsPlanPeriods(%q, %q) = %t, want %t", tt.a, tt.b, got, tt.want)
		}
	}
}

// Savings plans are subscribed for one or three years and cannot be deleted,
// so this test imports an existing Savings Plan instead of creating one.
// It renames the plan during the test and restores its display name at the end.
func TestAccResourceSavingsPlan_import(t *testing.T) {
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}
	testAccPreCheckCloud(t)
	checkEnvOrSkip(t, "OVH_SAVINGS_PLAN_ID_TEST")

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	savingsPlanID := os.Getenv("OVH_SAVINGS_PLAN_ID_TEST")

	serviceID, err := serviceIdFromResourceName(testAccOVHClient, serviceName)
	if err != nil {
		t.Fatal(err)
	}

	var plan savingsPlansSimulateResponse
	endpoint := fmt.Sprintf("/services/%d/savingsPlans/subscribed/%s", serviceID, url.PathEscape(savingsPlanID))
	if err := testAccOVHClient.Get(endpoint, &plan); err != nil {
		t.Fatalf("error calling GET %s: %s", endpoint, err)
	}

	flavor, deploymentType := normalizeSavingsPlanFlavor(plan.Flavor)
	newDisplayName := acctest.RandomWithPrefix(test_prefix)

	// Destroying the resource doesn't change the Savings Plan, so its display name is restored
	// through the API, even when a step fails.
	t.Cleanup(func() {
		var current savingsPlansSimulateResponse
		if err := testAccOVHClient.Get(endpoint, &current); err != nil {
			t.Errorf("error calling GET %s: %s", endpoint, err)
			return
		}
		if current.DisplayName == plan.DisplayName {
			return
		}
		if err := testAccOVHClient.Put(endpoint, map[string]string{"displayName": plan.DisplayName}, nil); err != nil {
			t.Errorf("failed to restore display name %q of Savings Plan %s: %s", plan.DisplayName, savingsPlanID, err)
		}
	})

	// The flavor is uppercased to check that it doesn't plan a replacement. The display name is passed
	// as a variable, so that any name of the existing plan is valid without escaping it in HCL.
	planConfig := fmt.Sprintf(`
		variable "display_name" {
			type = string
		}

		resource "ovh_savings_plan" "sp" {
			service_name    = "%s"
			flavor          = "%s"
			deployment_type = "%s"
			period          = "%s"
			size            = %d
			display_name    = var.display_name
			auto_renewal    = %t
		}`,
		serviceName,
		strings.ToUpper(flavor),
		deploymentType,
		plan.Period,
		plan.Size,
		plan.PeriodEndAction == "REACTIVATE",
	)
	displayName := func(name string) config.Variables {
		return config.Variables{"display_name": config.StringVariable(name)}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:             planConfig,
				ConfigVariables:    displayName(plan.DisplayName),
				ResourceName:       "ovh_savings_plan.sp",
				ImportState:        true,
				ImportStateId:      serviceName + "/" + savingsPlanID,
				ImportStatePersist: true,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						return fmt.Errorf("expected 1 imported state, got %d", len(states))
					}
					want := map[string]string{
						"service_name":    serviceName,
						"service_id":      strconv.Itoa(serviceID),
						"flavor":          flavor,
						"deployment_type": deploymentType,
						"period":          plan.Period,
						"size":            strconv.Itoa(plan.Size),
						"display_name":    plan.DisplayName,
					}
					for k, v := range want {
						if got := states[0].Attributes[k]; got != v {
							return fmt.Errorf("imported %s = %q, want %q", k, got, v)
						}
					}
					return nil
				},
			},
			{
				// The imported state must match the configuration
				Config:          planConfig,
				ConfigVariables: displayName(plan.DisplayName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config:          planConfig,
				ConfigVariables: displayName(newDisplayName),
				// Never apply a replacement, it would subscribe a new Savings Plan
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ovh_savings_plan.sp", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_savings_plan.sp", "display_name", newDisplayName),
					resource.TestCheckResourceAttr("ovh_savings_plan.sp", "id", savingsPlanID),
				),
			},
			{
				Config:          planConfig,
				ConfigVariables: displayName(plan.DisplayName),
				// Never apply a replacement, it would subscribe a new Savings Plan
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ovh_savings_plan.sp", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_savings_plan.sp", "display_name", plan.DisplayName),
					resource.TestCheckResourceAttr("ovh_savings_plan.sp", "id", savingsPlanID),
				),
			},
		},
	})
}
