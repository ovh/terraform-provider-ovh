package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVpsMigrationDataSourceConfig = `
data "ovh_vps_migration" "mig" {
  service_name = "%s"
}
`

func TestAccVpsMigrationDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_VPS")
	if serviceName == "" {
		t.Skip("OVH_VPS must be set for this acceptance test")
	}

	config := fmt.Sprintf(testAccVpsMigrationDataSourceConfig, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_vps_migration.mig", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("data.ovh_vps_migration.mig", "status"),
					resource.TestCheckResourceAttrSet("data.ovh_vps_migration.mig", "current_plan"),
				),
			},
		},
	})
}
