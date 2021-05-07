package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_cloud_project_containerregistry", &resource.Sweeper{
		Name: "ovh_cloud_project_containerregistry",
		F:    testSweepCloudProjectContainerRegistry,
	})
}

func testSweepCloudProjectContainerRegistry(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_CLOUD_PROJECT_SERVICE_TEST is not set. No container registry to sweep")
		return nil
	}

	regs := []CloudProjectContainerRegistry{}

	endpoint := fmt.Sprintf(
		"/cloud/project/%s/containerRegistry",
		url.PathEscape(serviceName),
	)
	if err := client.Get(endpoint, &regs); err != nil {
		return fmt.Errorf("Error calling %s:\n\t %q", endpoint, err)
	}

	if len(regs) == 0 {
		log.Print("[DEBUG] No container registry to sweep")
		return nil
	}

	for _, reg := range regs {
		if !strings.HasPrefix(reg.Name, test_prefix) {
			continue
		}

		log.Printf("[DEBUG] container registry found %v", reg.Name)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting container registry %s/%s", reg.Name, reg.Id)
			endpoint := fmt.Sprintf(
				"/cloud/project/%s/containerRegistry/%s",
				url.PathEscape(serviceName),
				url.PathEscape(reg.Id),
			)
			if err := client.Delete(endpoint, nil); err != nil {
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

func TestAccCloudProjectContainerRegistry_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	regName := acctest.RandomWithPrefix(test_prefix)
	region := "GRA"
	config := fmt.Sprintf(
		testAccCloudProjectContainerRegistryConfig,
		serviceName,
		region,
		regName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_containerregistry.reg", "name", regName),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_containerregistry.reg", "region", region),
					resource.TestCheckResourceAttr(
						"ovh_cloud_project_containerregistry.reg", "plan.0.name", "SMALL"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_containerregistry.reg", "url"),
				),
			},
		},
	})
}

const testAccCloudProjectContainerRegistryConfig = `
data "ovh_cloud_project_capabilities_containerregistry_filter" "regcap" {
	service_name = "%s"
    plan_name    = "SMALL"
    region       = "%s"
}

resource "ovh_cloud_project_containerregistry" "reg" {
	service_name = data.ovh_cloud_project_capabilities_containerregistry_filter.regcap.service_name
    plan_id      = data.ovh_cloud_project_capabilities_containerregistry_filter.regcap.id
	name         = "%s"
    region       = data.ovh_cloud_project_capabilities_containerregistry_filter.regcap.region
}
`
