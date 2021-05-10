package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloudProjectContainerRegistry_importBasic(t *testing.T) {
	regName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
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
			},
			{
				ResourceName:            "ovh_cloud_project_containerregistry.reg",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"plan_id"},
				ImportStateIdFunc:       testAccCloudProjectContainerRegistryImportId("ovh_cloud_project_containerregistry.reg"),
			},
		},
	})
}

func testAccCloudProjectContainerRegistryImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		testFarmServer, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_cloud_project_containerregistry not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			testFarmServer.Primary.Attributes["service_name"],
			testFarmServer.Primary.Attributes["id"],
		), nil
	}
}
