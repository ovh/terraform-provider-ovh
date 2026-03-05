package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccCloudInstance_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorId := os.Getenv("OVH_CLOUD_PROJECT_FLAVOR_ID_TEST")
	imageId := os.Getenv("OVH_CLOUD_PROJECT_IMAGE_ID_TEST")

	instanceName := acctest.RandomWithPrefix(testAccResourceCloudInstanceNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_instance" "instance" {
  service_name = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"
  region       = "%s"
  networks     = [
	{
		public = true
	}
  ]
}
`, serviceName, instanceName, flavorId, imageId, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckCloudInstanceV2(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.instance", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_instance.instance", "name", instanceName),
					resource.TestCheckResourceAttr("ovh_cloud_instance.instance", "flavor_id", flavorId),
					resource.TestCheckResourceAttr("ovh_cloud_instance.instance", "image_id", imageId),
					resource.TestCheckResourceAttr("ovh_cloud_instance.instance", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance.instance", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance.instance", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance.instance", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance.instance", "resource_status"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudInstanceImportStateIdFunc("ovh_cloud_instance.instance"),
			},
		},
	})
}

func TestAccCloudInstance_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorId := os.Getenv("OVH_CLOUD_PROJECT_FLAVOR_ID_TEST")
	imageId := os.Getenv("OVH_CLOUD_PROJECT_IMAGE_ID_TEST")

	instanceName := acctest.RandomWithPrefix(testAccResourceCloudInstanceNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudInstanceNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_instance" "instance" {
  service_name = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"
  region       = "%s"
}
`, serviceName, instanceName, flavorId, imageId, region)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_instance" "instance" {
  service_name = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"
  region       = "%s"
}
`, serviceName, updatedName, flavorId, imageId, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckCloudInstanceV2(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.instance", "name", instanceName),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.instance", "name", updatedName),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance.instance", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_instance.instance", "checksum"),
				),
			},
		},
	})
}

func TestAccCloudInstance_withNetworks(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavorId := os.Getenv("OVH_CLOUD_PROJECT_FLAVOR_ID_TEST")
	imageId := os.Getenv("OVH_CLOUD_PROJECT_IMAGE_ID_TEST")
	networkId := os.Getenv("OVH_CLOUD_PROJECT_NETWORK_ID_TEST")

	instanceName := acctest.RandomWithPrefix(testAccResourceCloudInstanceNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_instance" "instance" {
  service_name = "%s"
  name         = "%s"
  flavor_id    = "%s"
  image_id     = "%s"
  region       = "%s"

  networks {
    id = "%s"
  }
}
`, serviceName, instanceName, flavorId, imageId, region, networkId)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			testAccPreCheckCloudInstanceV2(t)
			testAccPreCheckCloudInstanceV2Network(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_instance.instance", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_instance.instance", "name", instanceName),
					resource.TestCheckResourceAttr("ovh_cloud_instance.instance", "networks.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_instance.instance", "networks.0.id", networkId),
				),
			},
		},
	})
}

const testAccResourceCloudInstanceNamePrefix = "tf-test-instance-v2-"

func testAccPreCheckCloudInstanceV2(t *testing.T) {
	if os.Getenv("OVH_CLOUD_PROJECT_FLAVOR_ID_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_FLAVOR_ID_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_IMAGE_ID_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_IMAGE_ID_TEST must be set for acceptance tests")
	}
}

func testAccPreCheckCloudInstanceV2Network(t *testing.T) {
	if os.Getenv("OVH_CLOUD_PROJECT_NETWORK_ID_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_NETWORK_ID_TEST must be set for acceptance tests")
	}
}

func testAccCloudInstanceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}
