package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccdDedicatedCloudDatacenter_importBasic(t *testing.T) {
	resourceName := "ovh_vrack_dedicated_cloud_datacenter.vrack-dedicatedCloudDatacenter"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDedicatedCloudDatacenter(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				Config:            testAccDedicatedCloudDatacenterConfig,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s/%s", os.Getenv("OVH_VRACK_SERVICE_TEST"), os.Getenv("OVH_VRACK_DEDICATED_CLOUD_DATACENTER"), os.Getenv("OVH_VRACK_TARGET_SERVICE_TEST")), nil
				},
			},
		},
	})
}
