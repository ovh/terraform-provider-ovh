package ovh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCloudManagedDatabasePostgresql_basic(t *testing.T) {
	name := acctest.RandomWithPrefix(test_prefix)
	config := testAccCloudManagedDatabasePostgresqlConfig(name, fmt.Sprintf("  description = %q", name), "") + `
data "ovh_cloud_managed_database" "db" {
  service_name = ovh_cloud_managed_database.db.service_name
  engine       = ovh_cloud_managed_database.db.engine
  id           = ovh_cloud_managed_database.db.id
}
`

	var checks []resource.TestCheckFunc
	for _, attribute := range []string{"region", "version", "plan", "flavor", "nodes_count", "description", "resource_status", "current_state.endpoints.0.uri"} {
		checks = append(checks, resource.TestCheckResourceAttrPair("data.ovh_cloud_managed_database.db", attribute, testAccCloudManagedDatabasePostgresqlResource, attribute))
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudManagedDatabasePostgresql(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudManagedServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}
