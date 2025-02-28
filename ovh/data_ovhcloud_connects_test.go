package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceOvhcloudConnects_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "ovh_ovhcloud_connects" "occs" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connects.occs", "occs.0.uuid"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connects.occs", "occs.0.bandwidth"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connects.occs", "occs.0.description"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connects.occs", "occs.0.interface_list.0"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connects.occs", "occs.0.pop"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connects.occs", "occs.0.port_quantity"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connects.occs", "occs.0.product"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connects.occs", "occs.0.provider_name"),
					resource.TestCheckResourceAttrSet("data.ovh_ovhcloud_connects.occs", "occs.0.status"),
				),
			},
		},
	})
}
