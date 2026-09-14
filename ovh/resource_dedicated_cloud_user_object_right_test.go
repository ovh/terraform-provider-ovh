package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccCheckOvhDedicatedCloudUserObjectRightConfig_basic = `
resource "ovh_dedicated_cloud_user" "testuser" {
	service_name = "%s"
	name         = "%s"
}

resource "ovh_dedicated_cloud_user_object_right" "testright" {
	service_name     = ovh_dedicated_cloud_user.testuser.service_name
	user_id          = ovh_dedicated_cloud_user.testuser.user_id
	type             = "cluster"
	vmware_object_id = "%s"
	right            = "readonly"
}
`

func TestAccDedicatedCloudUserObjectRight_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DEDICATED_CLOUD_TEST")
	objectId := os.Getenv("OVH_DEDICATED_CLOUD_OBJECT_ID_TEST")
	// The dedicatedCloud user "name" is validated server-side as a shortname:
	// hyphens (and thus acctest.RandomWithPrefix) are rejected with a 400.
	userName := "tfacctest" + acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDedicatedCloudUserObjectRight(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhDedicatedCloudUserObjectRightConfig_basic, serviceName, userName, objectId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dedicated_cloud_user_object_right.testright", "type", "cluster"),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_cloud_user_object_right.testright", "vmware_object_id", objectId),
					resource.TestCheckResourceAttr(
						"ovh_dedicated_cloud_user_object_right.testright", "right", "readonly"),
					resource.TestCheckResourceAttrSet(
						"ovh_dedicated_cloud_user_object_right.testright", "object_right_id"),
					resource.TestCheckResourceAttrSet(
						"ovh_dedicated_cloud_user_object_right.testright", "name"),
				),
			},
		},
	})
}

func TestAccDedicatedCloudUserObjectRight_importBasic(t *testing.T) {
	serviceName := os.Getenv("OVH_DEDICATED_CLOUD_TEST")
	objectId := os.Getenv("OVH_DEDICATED_CLOUD_OBJECT_ID_TEST")
	// The dedicatedCloud user "name" is validated server-side as a shortname:
	// hyphens (and thus acctest.RandomWithPrefix) are rejected with a 400.
	userName := "tfacctest" + acctest.RandString(8)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDedicatedCloudUserObjectRight(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckOvhDedicatedCloudUserObjectRightConfig_basic, serviceName, userName, objectId),
			},
			{
				ResourceName:                         "ovh_dedicated_cloud_user_object_right.testright",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "object_right_id",
				ImportStateIdFunc:                    testAccDedicatedCloudUserObjectRight_import("ovh_dedicated_cloud_user_object_right.testright"),
			},
		},
	})
}

func testAccDedicatedCloudUserObjectRight_import(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("ovh_dedicated_cloud_user_object_right not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s/%s",
			rs.Primary.Attributes["service_name"],
			rs.Primary.Attributes["user_id"],
			rs.Primary.Attributes["object_right_id"],
		), nil
	}
}
