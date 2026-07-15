package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccCloudExtNetIPInstanceConfig returns the configuration of an
// instance attached to a public network, so that an Ext-Net IP exists in the
// project. It mirrors the fixture of resource_cloud_project_instance_test.go.
func testAccCloudExtNetIPInstanceConfig(serviceName, region, image, flavor string) string {
	return fmt.Sprintf(`
resource "ovh_cloud_project_instance" "instance" {
	service_name = "%s"
	region = "%s"
	billing_period = "hourly"
	boot_from {
		image_id = "%s"
	}
	flavor {
		flavor_id = "%s"
	}
	name = "TestInstance"
	ssh_key {
		name = "%s"
	}
	network {
		public = true
	}
}
`,
		serviceName,
		region,
		image,
		flavor,
		os.Getenv("OVH_CLOUD_PROJECT_SSH_NAME_TEST"))
}

// testAccPreCheckCloudExtNetIP is the PreCheck for the Ext-Net IP data
// source tests, which spawn an instance attached to a public network.
func testAccPreCheckCloudExtNetIP(t *testing.T) {
	testAccPreCheckCloudPublicIP(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_SSH_NAME_TEST")
}

func TestAccDataSourceCloudExtNetIP_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavor, image, err := getFlavorAndImage(serviceName, region)
	if err != nil {
		t.Skipf("failed to retrieve a flavor and an image: %s", err)
	}

	config := testAccCloudExtNetIPInstanceConfig(serviceName, region, image, flavor) + fmt.Sprintf(`
data "ovh_cloud_ext_net_ip" "test" {
  service_name = "%s"
  id           = [for address in ovh_cloud_project_instance.instance.addresses : address.ip if address.version == 4][0]
}
`, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudExtNetIP(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_ext_net_ip.test", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_ext_net_ip.test", "id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_ext_net_ip.test", "resource_status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_ext_net_ip.test", "current_state.ip"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_ext_net_ip.test", "current_state.id"),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_ext_net_ip.test", "id",
						"data.ovh_cloud_ext_net_ip.test", "current_state.ip",
					),
				),
			},
		},
	})
}

// TestAccDataSourceCloudExtNetIP_serviceNameFromEnv mirrors the _basic test but
// omits service_name from the data block: the data source must resolve the
// project id from the OVH_CLOUD_PROJECT_SERVICE environment variable.
func TestAccDataSourceCloudExtNetIP_serviceNameFromEnv(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavor, image, err := getFlavorAndImage(serviceName, region)
	if err != nil {
		t.Skipf("failed to retrieve a flavor and an image: %s", err)
	}

	t.Setenv("OVH_CLOUD_PROJECT_SERVICE", serviceName)

	// The instance keeps an explicit service_name so setup stays deterministic;
	// only the data block relies on the environment fallback.
	config := testAccCloudExtNetIPInstanceConfig(serviceName, region, image, flavor) + `
data "ovh_cloud_ext_net_ip" "test" {
  id = [for address in ovh_cloud_project_instance.instance.addresses : address.ip if address.version == 4][0]
}
`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudExtNetIP(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_ext_net_ip.test", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_ext_net_ip.test", "id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_ext_net_ip.test", "resource_status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_ext_net_ip.test", "current_state.ip"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_ext_net_ip.test", "current_state.id"),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_ext_net_ip.test", "id",
						"data.ovh_cloud_ext_net_ip.test", "current_state.ip",
					),
				),
			},
		},
	})
}
