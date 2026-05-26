package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccVpsMigrationConfig = `
resource "ovh_vps_migration" "mig" {
  service_name = "%s"
  target_plan  = "%s"
}
`

const testAccVpsMigrationScheduledConfig = `
resource "ovh_vps_migration" "mig" {
  service_name   = "%s"
  target_plan    = "%s"
  scheduled_date = "%s"
}
`

func TestAccVpsMigration_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_VPS")
	targetPlan := os.Getenv("OVH_VPS_MIGRATION_TARGET_PLAN")

	if serviceName == "" || targetPlan == "" {
		t.Skip("OVH_VPS and OVH_VPS_MIGRATION_TARGET_PLAN must be set for this acceptance test")
	}

	config := fmt.Sprintf(testAccVpsMigrationConfig, serviceName, targetPlan)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vps_migration.mig", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_vps_migration.mig", "target_plan", targetPlan),
					resource.TestCheckResourceAttrSet("ovh_vps_migration.mig", "status"),
					resource.TestCheckResourceAttrSet("ovh_vps_migration.mig", "current_plan"),
				),
			},
		},
	})
}

func TestAccVpsMigration_scheduled(t *testing.T) {
	serviceName := os.Getenv("OVH_VPS")
	targetPlan := os.Getenv("OVH_VPS_MIGRATION_TARGET_PLAN")
	scheduledDate := os.Getenv("OVH_VPS_MIGRATION_DATE")

	if serviceName == "" || targetPlan == "" || scheduledDate == "" {
		t.Skip("OVH_VPS, OVH_VPS_MIGRATION_TARGET_PLAN and OVH_VPS_MIGRATION_DATE must be set for this acceptance test")
	}

	config := fmt.Sprintf(testAccVpsMigrationScheduledConfig, serviceName, targetPlan, scheduledDate)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_vps_migration.mig", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_vps_migration.mig", "target_plan", targetPlan),
					resource.TestCheckResourceAttr("ovh_vps_migration.mig", "scheduled_date", scheduledDate),
				),
			},
		},
	})
}
