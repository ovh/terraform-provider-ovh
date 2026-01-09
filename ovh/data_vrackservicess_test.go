package ovh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVrackservicessDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCredentials(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: "data ovh_vrackservicess vrackservicess {}",
				Check:  resource.TestCheckResourceAttrSet("data.ovh_vrackservicess.vrackservicess", "vrackservicess.#"),
			},
		},
	})
}
