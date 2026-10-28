package ovh

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/go-ovh/ovh"
)

func testAccPreCheckCloudStorageSwiftContainer(t *testing.T) {
	testAccPreCheckCloud(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_STORAGE_REGION_TEST")
}

func testAccCloudStorageSwiftContainerImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}

func testAccCheckCloudStorageSwiftContainerDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_cloud_storage_swift_container" {
			continue
		}

		endpoint := storageSwiftContainerEndpoint(rs.Primary.Attributes["service_name"]) + "/" + url.PathEscape(rs.Primary.Attributes["id"])
		err := testAccOVHClient.Get(endpoint, &CloudStorageSwiftContainerAPIResponse{})
		if err == nil {
			return fmt.Errorf("container %s still exists", rs.Primary.Attributes["id"])
		}
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			continue
		}
		return fmt.Errorf("error checking container %s was destroyed: %w", rs.Primary.Attributes["id"], err)
	}
	return nil
}

func TestAccCloudStorageSwiftContainer_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))

	containerName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_swift_container" "container" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
}
`, serviceName, containerName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageSwiftContainer(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageSwiftContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "name", containerName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "id", region+"_"+containerName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "archive", "false"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "type", "PRIVATE"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "cors.#", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_swift_container.container", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_swift_container.container", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_swift_container.container", "updated_at"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.name", containerName),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.location.region", region),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.archive", "false"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.type", "PRIVATE"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.objects_count", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.objects_size", "0"),
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.static_url"),
				),
			},
			// Guards against an Optional+Computed attribute re-planning an update when left unset.
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:      "ovh_cloud_storage_swift_container.container",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageSwiftContainerImportStateIdFunc("ovh_cloud_storage_swift_container.container"),
			},
		},
	})
}

func TestAccCloudStorageSwiftContainer_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))

	containerName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_swift_container" "container" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  type         = "PRIVATE"
}
`, serviceName, containerName, region)

	publicConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_swift_container" "container" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  type         = "PUBLIC"
  cors         = ["https://www.example.com", "https://app.example.com"]
}
`, serviceName, containerName, region)

	staticConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_swift_container" "container" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  type         = "STATIC"
  cors         = ["https://www.example.com"]
}
`, serviceName, containerName, region)

	prunedConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_swift_container" "container" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  type         = "PRIVATE"
  cors         = []
}
`, serviceName, containerName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageSwiftContainer(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageSwiftContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "type", "PRIVATE"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "cors.#", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.type", "PRIVATE"),
				),
			},
			{
				Config: publicConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ovh_cloud_storage_swift_container.container", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "type", "PUBLIC"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "cors.#", "2"),
					resource.TestCheckTypeSetElemAttr("ovh_cloud_storage_swift_container.container", "cors.*", "https://www.example.com"),
					resource.TestCheckTypeSetElemAttr("ovh_cloud_storage_swift_container.container", "cors.*", "https://app.example.com"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.type", "PUBLIC"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.cors.#", "2"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_swift_container.container", "current_state.static_url"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "resource_status", "READY"),
				),
			},
			{
				Config: staticConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "type", "STATIC"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "cors.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.type", "STATIC"),
					resource.TestCheckResourceAttrSet("ovh_cloud_storage_swift_container.container", "current_state.static_url"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "resource_status", "READY"),
				),
			},
			{
				Config: prunedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "type", "PRIVATE"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "cors.#", "0"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.type", "PRIVATE"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.cors.#", "0"),
					resource.TestCheckNoResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.static_url"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "resource_status", "READY"),
				),
			},
			{
				ResourceName:      "ovh_cloud_storage_swift_container.container",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageSwiftContainerImportStateIdFunc("ovh_cloud_storage_swift_container.container"),
			},
		},
	})
}

func TestAccCloudStorageSwiftContainer_archive(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))

	containerName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_storage_swift_container" "container" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  archive      = true
}
`, serviceName, containerName, region)

	standardConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_swift_container" "container" {
  service_name = "%s"
  name         = "%s"
  region       = "%s"
  archive      = false
}
`, serviceName, containerName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageSwiftContainer(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageSwiftContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "archive", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.archive", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "resource_status", "READY"),
				),
			},
			{
				Config: config,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			// archive is immutable: switching the policy recreates the container.
			{
				Config: standardConfig,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("ovh_cloud_storage_swift_container.container", plancheck.ResourceActionReplace),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "archive", "false"),
					resource.TestCheckResourceAttr("ovh_cloud_storage_swift_container.container", "current_state.archive", "false"),
				),
			},
			{
				ResourceName:      "ovh_cloud_storage_swift_container.container",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudStorageSwiftContainerImportStateIdFunc("ovh_cloud_storage_swift_container.container"),
			},
		},
	})
}
