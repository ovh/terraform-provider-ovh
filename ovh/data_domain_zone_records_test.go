package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainZoneRecordsDataSource_basic(t *testing.T) {
	zoneName := os.Getenv("OVH_ZONE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data ovh_domain_zone_records recs {
						zone_name = %q
					}`, zoneName),
				Check: resource.TestCheckResourceAttrSet("data.ovh_domain_zone_records.recs", "ids.#"),
			},
		},
	})
}
