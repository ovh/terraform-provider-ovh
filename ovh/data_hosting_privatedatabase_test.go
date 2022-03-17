package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHostingPrivateDatabase_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckHostingPrivateDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHostingPrivateDatabaseConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase.database",
						"cpu",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase.database",
						"datacenter",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase.database",
						"offer",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase.database",
						"type",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_hosting_privatedatabase.database",
						"version_number",
					),
				),
			},
		},
	})
}

var testAccHostingPrivateDatabaseConfig = fmt.Sprintf(`
data "ovh_hosting_privatedatabase" "database" {
  service_name = "%s"
}
`, os.Getenv("OVH_HOSTING_PRIVATEDABASE_SERVICE_TEST"))
