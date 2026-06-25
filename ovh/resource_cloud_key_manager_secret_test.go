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
	resource.AddTestSweepers("ovh_cloud_key_manager_secret", &resource.Sweeper{
		Name: "ovh_cloud_key_manager_secret",
		F:    testSweepCloudKeyManagerSecret,
	})
}

func testSweepCloudKeyManagerSecret(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if serviceName == "" {
		return nil
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/keyManager/secret"
	var secrets []CloudKeyManagerSecretAPIResponse
	if err := client.Get(endpoint, &secrets); err != nil {
		return fmt.Errorf("error listing key manager secrets: %s", err)
	}

	for _, s := range secrets {
		name := ""
		if s.TargetSpec != nil {
			name = s.TargetSpec.Name
		}
		if len(name) < len(test_prefix) || name[:len(test_prefix)] != test_prefix {
			continue
		}
		deleteEndpoint := endpoint + "/" + url.PathEscape(s.Id)
		if err := client.Delete(deleteEndpoint, nil); err != nil {
			return fmt.Errorf("error deleting key manager secret %s: %s", s.Id, err)
		}
	}

	return nil
}

func testAccPreCheckCloudKeyManager(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_REGION_TEST")
}

func testAccCloudKeyManagerSecretImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["service_name"], rs.Primary.Attributes["id"]), nil
	}
}

func TestAccCloudKeyManagerSecret_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_key_manager_secret" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  secret_type  = "OPAQUE"
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
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "region", region),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "secret_type", "OPAQUE"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_secret.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_secret.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_secret.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "current_state.location.region", region),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_key_manager_secret.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudKeyManagerSecretImportStateIdFunc("ovh_cloud_key_manager_secret.test"),
			},
		},
	})
}

func TestAccCloudKeyManagerSecret_withPayload(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_key_manager_secret" "test" {
  service_name         = "%s"
  region               = "%s"
  name                 = "%s"
  secret_type          = "OPAQUE"
  payload              = base64encode("my-secret-payload")
  payload_content_type = "APPLICATION_OCTET_STREAM"
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
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "secret_type", "OPAQUE"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "payload_content_type", "APPLICATION_OCTET_STREAM"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_secret.test", "id"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "resource_status", "READY"),
				),
			},
		},
	})
}

func TestAccCloudKeyManagerSecret_updateMetadata(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	name := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_key_manager_secret" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  secret_type  = "OPAQUE"
  metadata = {
    env = "test"
  }
}
`, serviceName, region, name)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_key_manager_secret" "test" {
  service_name = "%s"
  region       = "%s"
  name         = "%s"
  secret_type  = "OPAQUE"
  metadata = {
    env   = "production"
    owner = "team-a"
  }
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
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "metadata.env", "test"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_secret.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_secret.test", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "resource_status", "READY"),
				),
			},
			// Mutate only metadata: the secret must update in place (same id) and
			// return to READY once the OUT_OF_SYNC re-sync flow completes.
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "name", name),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "metadata.env", "production"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "metadata.owner", "team-a"),
					resource.TestCheckResourceAttrSet("ovh_cloud_key_manager_secret.test", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_key_manager_secret.test", "resource_status", "READY"),
				),
			},
		},
	})
}
