package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
)

const testAccCloudProjectBasic = `
data "ovh_me" "myaccount" {}

data "ovh_order_cart" "mycart" {
	ovh_subsidiary = data.ovh_me.myaccount.ovh_subsidiary
}

data "ovh_order_cart_product_plan" "cloud" {
	cart_id        = data.ovh_order_cart.mycart.id
	price_capacity = "renew"
	product        = "cloud"
	plan_code      = "project.2018"
}

resource "ovh_cloud_project" "cloud" {
	ovh_subsidiary = data.ovh_order_cart.mycart.ovh_subsidiary
	description    = "%s"
 
	plan {
		duration     = data.ovh_order_cart_product_plan.cloud.selected_price.0.duration
		plan_code    = data.ovh_order_cart_product_plan.cloud.plan_code
		pricing_mode = data.ovh_order_cart_product_plan.cloud.selected_price.0.pricing_mode

		configuration {
			label = "vrack"
			value = "%s"
		}
	}
}
`

func init() {
	resource.AddTestSweepers("ovh_cloud_project", &resource.Sweeper{
		Name: "ovh_cloudProject",
		F:    testSweepCloudProject,
	})
}

func testSweepCloudProject(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceNames := make([]string, 0)
	if err := config.OVHClient.Get("/cloud/project", &serviceNames); err != nil {
		return fmt.Errorf("Error calling GET /cloud/project:\n\t %q", err)
	}

	if len(serviceNames) == 0 {
		log.Print("[DEBUG] No cloudProject to sweep")
		return nil
	}

	for _, serviceName := range serviceNames {
		r := &CloudProject{}
		log.Printf("[DEBUG] Will get cloudProject: %v", serviceName)
		endpoint := fmt.Sprintf(
			"/cloud/project/%s",
			url.PathEscape(serviceName),
		)

		if err := config.OVHClient.Get(endpoint, r); err != nil {
			return fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
		}

		if r.Description == nil || !strings.HasPrefix(*r.Description, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] Will delete cloudProject: %v", serviceName)

		terminate := func() (string, error) {
			log.Printf("[DEBUG] Will terminate cloudProject %s", serviceName)
			endpoint := fmt.Sprintf(
				"/cloud/project/%s/terminate",
				url.PathEscape(serviceName),
			)
			if err := config.OVHClient.Post(endpoint, nil, nil); err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && (errOvh.Code == 404 || errOvh.Code == 460) {
					return "", nil
				}
				return "", fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
			}
			return serviceName, nil
		}

		confirmTerminate := func(token string) error {
			log.Printf("[DEBUG] Will confirm termination of cloudProject %s", serviceName)
			endpoint := fmt.Sprintf(
				"/cloud/project/%s/confirmTermination",
				url.PathEscape(serviceName),
			)
			if err := config.OVHClient.Post(endpoint, &CloudProjectConfirmTerminationOpts{Token: token}, nil); err != nil {
				return fmt.Errorf("calling Post %s:\n\t %q", endpoint, err)
			}
			return nil
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := orderDeleteFromResource(nil, config, terminate, confirmTerminate); err != nil {
				return resource.RetryableError(err)
			}

			// Successful delete
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func TestResourceCloudProjectDelete_deletionProtection(t *testing.T) {
	resourceDef := resourceCloudProject()
	d := resourceDef.TestResourceData()
	d.SetId("test-project-id")
	d.Set("project_id", "test-project-id")
	d.Set("deletion_protection", true)

	err := resourceCloudProjectDelete(d, &Config{})
	if err == nil {
		t.Fatal("expected error when deletion_protection is true, got nil")
	}
	if !strings.Contains(err.Error(), "protected from deletion") {
		t.Fatalf("expected deletion protection error, got: %s", err)
	}

	// Verify that deletion_protection=false does not trigger the guard.
	// The delete will panic downstream due to nil API client, but the
	// important thing is that the deletion protection guard did not fire.
	d2 := resourceDef.TestResourceData()
	d2.SetId("test-project-id")
	d2.Set("project_id", "test-project-id")
	d2.Set("deletion_protection", false)

	func() {
		defer func() {
			if r := recover(); r != nil {
				// Expected: nil client panic means we passed the guard
			}
		}()
		err = resourceCloudProjectDelete(d2, &Config{})
		if err != nil && strings.Contains(err.Error(), "protected from deletion") {
			t.Fatalf("should not get deletion protection error when set to false, got: %s", err)
		}
	}()
}

func TestAccResourceCloudProject_basic(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)

	// Link a vrack to the new cloud project at order time
	vrackServiceName := os.Getenv("OVH_VRACK_SERVICE_TEST")

	config := fmt.Sprintf(
		testAccCloudProjectBasic,
		desc,
		vrackServiceName,
	)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckOrderCloudProject(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project.cloud", "description", desc),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project.cloud", "project_id"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project.cloud", "project_name"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project.cloud", "urn"),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project.cloud", "plan.0.plan_code", "project.2018"),
				),
			},
			{
				ResourceName:            "ovh_cloud_project.cloud",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"plan", "ovh_subsidiary", "order"},
			},
		},
	})
}
