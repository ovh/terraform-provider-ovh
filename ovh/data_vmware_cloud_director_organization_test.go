package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVMwareCloudDirectorOrganizationData(t *testing.T) {
	serviceName := os.Getenv("OVH_VCD_ORGANIZATION")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			checkEnvOrSkip(t, "OVH_VCD_ORGANIZATION")
			testAccPreCheckCredentials(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_vmware_cloud_director_organization" "org" {
						organization_id = "%s"
					}
				`, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_vmware_cloud_director_organization.org", "id", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_organization.org", "resource_status"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_organization.org", "current_state.api_url"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_organization.org", "current_state.name"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_organization.org", "current_state.region"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_organization.org", "target_spec.full_name"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_organization.org", "iam.urn"),
				),
			},
		},
	})
}
