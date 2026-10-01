package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

// Throwaway, well-formed SSH public keys used by the immutability tests. They
// are never deployed to any host; they only need to parse as valid public keys.
const (
	testAccCloudSshKeyPublicKeyA = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQC9xPpdqP3sx2H+gcBm65tJEaUbuifQ1uGkgrWtNY0PRKNNPdy+3yoVOtxk6Vjo4YZ0EU/JhmQfnrK7X7Q5vhqYxmozi0LiTRt0BxgqHJ+4hWTWMIOgr+C2jLx7ZsCReRk+fy5AHr6h0PHQEuXVLXeUy/TDyuY2JPtUZ5jcqvLYgQ== my-nuclear-power-plant"
	testAccCloudSshKeyPublicKeyB = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICVg+WcevAiCCkj9L6dKVqMv0yseIYOWvUWhnwWbzuhH throwaway-ed25519"
)

func TestAccCloudSshKey_basic(t *testing.T) {
	keyName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(`
	resource "ovh_cloud_ssh_key" "key" {
		service_name = "%s"
		name         = "%s"
		public_key   = "%s"
	}
	`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), keyName, testAccCloudSshKeyPublicKeyA)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_ssh_key.key", "name", keyName),
					resource.TestCheckResourceAttrSet("ovh_cloud_ssh_key.key", "public_key"),
					resource.TestCheckResourceAttrSet("ovh_cloud_ssh_key.key", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_cloud_ssh_key.key", "updated_at"),
				),
			},
			{
				// Re-planning the unchanged config must be a no-op (regression guard).
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudSshKey_noServiceName validates that when service_name is not set in
// the resource configuration, the provider falls back to the OVH_CLOUD_PROJECT_SERVICE
// environment variable at plan time (via the EnvDefaultString plan modifier) and that
// this does not produce a perpetual diff / phantom replace on subsequent plans.
func TestAccCloudSshKey_noServiceName(t *testing.T) {
	keyName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(`
	resource "ovh_cloud_ssh_key" "key" {
		name       = "%s"
		public_key = "%s"
	}
	`, keyName, testAccCloudSshKeyPublicKeyA)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_SERVICE")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_ssh_key.key", "name", keyName),
					resource.TestCheckResourceAttr("ovh_cloud_ssh_key.key", "service_name", os.Getenv("OVH_CLOUD_PROJECT_SERVICE")),
					resource.TestCheckResourceAttrSet("ovh_cloud_ssh_key.key", "public_key"),
				),
			},
			{
				// Re-planning with service_name still omitted must be a no-op:
				// the EnvDefaultString modifier injects the env value so the plan
				// matches state and RequiresReplace does not fire (regression guard).
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// TestAccCloudSshKey_renameForcesReplace validates that the resource is immutable:
// changing `name` forces a replace (name has RequiresReplace, no Update path).
func TestAccCloudSshKey_renameForcesReplace(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	keyName := acctest.RandomWithPrefix(test_prefix)
	renamedKeyName := acctest.RandomWithPrefix(test_prefix)

	config := func(name string) string {
		return fmt.Sprintf(`
	resource "ovh_cloud_ssh_key" "key" {
		service_name = "%s"
		name         = "%s"
		public_key   = "%s"
	}
	`, serviceName, name, testAccCloudSshKeyPublicKeyA)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(keyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_ssh_key.key", "name", keyName),
				),
			},
			{
				// Changing the name must replace the resource.
				Config: config(renamedKeyName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"ovh_cloud_ssh_key.key",
							plancheck.ResourceActionReplace,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_ssh_key.key", "name", renamedKeyName),
				),
			},
		},
	})
}

// TestAccCloudSshKey_publicKeyForcesReplace validates that the resource is
// immutable: changing `public_key` forces a replace (public_key has
// RequiresReplace, no Update path).
func TestAccCloudSshKey_publicKeyForcesReplace(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	keyName := acctest.RandomWithPrefix(test_prefix)

	config := func(publicKey string) string {
		return fmt.Sprintf(`
	resource "ovh_cloud_ssh_key" "key" {
		service_name = "%s"
		name         = "%s"
		public_key   = "%s"
	}
	`, serviceName, keyName, publicKey)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(testAccCloudSshKeyPublicKeyA),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_ssh_key.key", "public_key", testAccCloudSshKeyPublicKeyA),
				),
			},
			{
				// Changing the public key must replace the resource.
				Config: config(testAccCloudSshKeyPublicKeyB),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(
							"ovh_cloud_ssh_key.key",
							plancheck.ResourceActionReplace,
						),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_ssh_key.key", "public_key", testAccCloudSshKeyPublicKeyB),
				),
			},
		},
	})
}
