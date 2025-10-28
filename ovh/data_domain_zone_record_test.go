package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainZoneRecordDataSource_basic(t *testing.T) {
	zoneName := os.Getenv("OVH_ZONE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data ovh_domain_zone_records recs {
						zone_name = %q
					}

					data ovh_domain_zone_record rec {
						zone_name = data.ovh_domain_zone_records.recs.zone_name
						id        = data.ovh_domain_zone_records.recs.ids[0]
					}`, zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_domain_zone_record.rec", "field_type"),
					resource.TestCheckResourceAttrSet("data.ovh_domain_zone_record.rec", "target"),
					resource.TestCheckResourceAttrSet("data.ovh_domain_zone_record.rec", "ttl"),
					resource.TestCheckResourceAttrSet("data.ovh_domain_zone_record.rec", "zone"),
				),
			},
		},
	})
}
