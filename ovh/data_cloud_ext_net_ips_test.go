package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudExtNetIPs_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")
	flavor, image, err := getFlavorAndImage(serviceName, region)
	if err != nil {
		t.Skipf("failed to retrieve a flavor and an image: %s", err)
	}

	config := testAccCloudExtNetIPInstanceConfig(serviceName, region, image, flavor) + fmt.Sprintf(`
data "ovh_cloud_ext_net_ips" "test" {
  service_name = "%s"

  depends_on = [ovh_cloud_project_instance.instance]
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
					resource.TestCheckResourceAttr("data.ovh_cloud_ext_net_ips.test", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_ext_net_ips.test", "ext_net_ips.#"),
					resource.TestCheckResourceAttrWith("data.ovh_cloud_ext_net_ips.test", "ext_net_ips.#", testAccCheckCloudPublicIPListNotEmpty),
				),
			},
		},
	})
}
