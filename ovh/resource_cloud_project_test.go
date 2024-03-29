package ovh

import (
	"fmt"
	"log"
	"net/url"
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

func TestAccResourceCloudProject_basic(t *testing.T) {
	desc := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(
		testAccCloudProjectBasic,
		desc,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckOrderCloudProject(t) },

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
				),
			},
		},
	})
}
