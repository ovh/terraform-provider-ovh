package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectSSHKey_basic(t *testing.T) {
	keyName := acctest.RandomWithPrefix(test_prefix)
	config := fmt.Sprintf(`
	resource "ovh_cloud_project_ssh_key" "key" {
		service_name = "%s"
		region       = "GRA9"
		name         = "%s"
		public_key   = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQC9xPpdqP3sx2H+gcBm65tJEaUbuifQ1uGkgrWtNY0PRKNNPdy+3yoVOtxk6Vjo4YZ0EU/JhmQfnrK7X7Q5vhqYxmozi0LiTRt0BxgqHJ+4hWTWMIOgr+C2jLx7ZsCReRk+fy5AHr6h0PHQEuXVLXeUy/TDyuY2JPtUZ5jcqvLYgQ== my-nuclear-power-plant"
	}
	`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"), keyName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_ssh_key.key", "name", keyName),
					resource.TestCheckResourceAttr("ovh_cloud_project_ssh_key.key", "regions.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_project_ssh_key.key", "regions.0", "GRA9"),
				),
			},
		},
	})
}
