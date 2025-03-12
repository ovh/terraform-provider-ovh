package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVMwareCloudDirectorBackupData(t *testing.T) {
	serviceName := os.Getenv("OVH_VCD_BACKUP")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			checkEnvOrSkip(t, "OVH_VCD_BACKUP")
			testAccPreCheckCredentials(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "ovh_vmware_cloud_director_backup" "backup" {
						backup_id = "%s"
					}
				`, serviceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_vmware_cloud_director_backup.backup", "id", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_backup.backup", "resource_status"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_backup.backup", "current_state.az_name"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_backup.backup", "current_state.offers.#"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_backup.backup", "current_state.offers.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_backup.backup", "current_state.offers.0.protection_primary_region"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_backup.backup", "target_spec.offers.#"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_backup.backup", "target_spec.offers.0.name"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_backup.backup", "target_spec.offers.0.quota_in_tb"),
					resource.TestCheckResourceAttrSet("data.ovh_vmware_cloud_director_backup.backup", "iam.urn"),
				),
			},
		},
	})
}
