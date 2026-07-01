package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudSshKeysDataSource_basic(t *testing.T) {
	keyName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(`
	resource "ovh_cloud_ssh_key" "key" {
		service_name = "%s"
		name         = "%s"
		public_key   = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQC9xPpdqP3sx2H+gcBm65tJEaUbuifQ1uGkgrWtNY0PRKNNPdy+3yoVOtxk6Vjo4YZ0EU/JhmQfnrK7X7Q5vhqYxmozi0LiTRt0BxgqHJ+4hWTWMIOgr+C2jLx7ZsCReRk+fy5AHr6h0PHQEuXVLXeUy/TDyuY2JPtUZ5jcqvLYgQ== my-nuclear-power-plant"
	}

	data "ovh_cloud_ssh_keys" "keys" {
		service_name = ovh_cloud_ssh_key.key.service_name
	}
	`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), keyName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_cloud_ssh_keys.keys", "ssh_keys.#"),
					resource.TestCheckTypeSetElemNestedAttrs("data.ovh_cloud_ssh_keys.keys", "ssh_keys.*", map[string]string{
						"name": keyName,
					}),
				),
			},
		},
	})
}
