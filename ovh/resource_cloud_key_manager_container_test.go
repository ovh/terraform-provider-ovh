package ovh

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("ovh_cloud_key_manager_container", &resource.Sweeper{
		Name:         "ovh_cloud_key_manager_container",
		F:            testSweepCloudKeyManagerContainer,
		Dependencies: []string{"ovh_cloud_key_manager_secret"},
	})
}

func testSweepCloudKeyManagerContainer(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if serviceName == "" {
		return nil
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/keyManager/container"
	var containers []CloudKeyManagerContainerAPIResponse
	if err := client.Get(endpoint, &containers); err != nil {
		return fmt.Errorf("error listing key manager containers: %s", err)
	}

	for _, c := range containers {
		name := ""
		if c.TargetSpec != nil {
			name = c.TargetSpec.Name
		}
		if len(name) < len(test_prefix) || name[:len(test_prefix)] != test_prefix {
			continue
		}
		deleteEndpoint := endpoint + "/" + url.PathEscape(c.Id)
		if err := client.Delete(deleteEndpoint, nil); err != nil {
			return fmt.Errorf("error deleting key manager container %s: %s", c.Id, err)
		}
	}

	return nil
}

func testAccCloudKeyManagerContainerImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}

func TestAccCloudKeyManagerContainer_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_key_manager_container" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  type         = "GENERIC"
}
`, serviceName, region, name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudKeyManager(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "type", "GENERIC"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_container.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_container.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_container.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "current_state.location.region", region),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_key_manager_container.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudKeyManagerContainerImportStateIdFunc("ovh_cloud_key_manager_container.test"),
			},
		},
	})
}

func TestAccCloudKeyManagerContainer_withSecretRefs(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	secretName := acctest.RandomWithPrefix(test_prefix)
	containerName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_key_manager_secret" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  secret_type  = "OPAQUE"
}

resource "ovh_cloud_key_manager_container" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  type         = "GENERIC"

  secret_refs = [
		{
			name      = "my-secret"
			secret_id = ovh_cloud_key_manager_secret.test.id
		}
	]
}
`, serviceName, region, secretName, serviceName, region, containerName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudKeyManager(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "name", containerName),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "type", "GENERIC"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "secret_refs.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "secret_refs.0.name", "my-secret"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_container.test", "secret_refs.0.secret_id"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "resource_status", "READY"),
				),
			},
		},
	})
}

func TestAccCloudKeyManagerContainer_updateSecretRefs(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	secretName1 := acctest.RandomWithPrefix(test_prefix)
	secretName2 := acctest.RandomWithPrefix(test_prefix)
	containerName := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_key_manager_secret" "secret1" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  secret_type  = "OPAQUE"
}

resource "ovh_cloud_key_manager_secret" "secret2" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  secret_type  = "OPAQUE"
}

resource "ovh_cloud_key_manager_container" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  type         = "GENERIC"

  secret_refs = [
    {
      name      = "first"
      secret_id = ovh_cloud_key_manager_secret.secret1.id
    }
  ]
}
`, serviceName, region, secretName1, serviceName, region, secretName2, serviceName, region, containerName)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_key_manager_secret" "secret1" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  secret_type  = "OPAQUE"
}

resource "ovh_cloud_key_manager_secret" "secret2" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  secret_type  = "OPAQUE"
}

resource "ovh_cloud_key_manager_container" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  type         = "GENERIC"

  secret_refs = [
    {
      name      = "first"
      secret_id = ovh_cloud_key_manager_secret.secret1.id
    },
	{
	  name      = "second"
	  secret_id = ovh_cloud_key_manager_secret.secret2.id
	}
  ]
}
`, serviceName, region, secretName1, serviceName, region, secretName2, serviceName, region, containerName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudKeyManager(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "secret_refs.#", "1"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_container.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_container.test", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "resource_status", "READY"),
				),
			},
			// Mutate only secret_refs: the container must update in place and
			// return to READY once the OUT_OF_SYNC re-sync flow completes.
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "secret_refs.#", "2"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_container.test", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_container.test", "resource_status", "READY"),
				),
			},
		},
	})
}
