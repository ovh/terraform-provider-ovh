package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	testAccIpLoadbalancingHttpFarmServerConfig_templ = `
data ovh_iploadbalancing iplb {
  service_name = "%s"
}

resource ovh_iploadbalancing_http_farm testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  display_name     = "%s"
  port = 8080
  zone = "all"
  probe {
    port     = 8080
    interval = 30
    type     = "http"
  }
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step0 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name = data.ovh_iploadbalancing.iplb.id
  farm_id      = ovh_iploadbalancing_http_farm.testacc.id
  address      = "10.0.0.11"
  status       = "active"
  display_name = "testBackendA"
  port         = 80
  weight       = 3
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step1 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendA"
  port = 8080
  weight = 3
  probe = false
  backup = false
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step2 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  weight = 2
  probe = true
  backup = true
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step3 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "inactive"
  display_name = "testBackendB"
  port = 80
  probe = false
  backup = false
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step4 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  proxy_protocol_version = "v2"
  ssl = true
  weight = 2
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step5 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  proxy_protocol_version = "v1"
  ssl    = true
  backup = false
}
`
	testAccIpLoadbalancingHttpFarmServerConfig_step6 = `
%s

resource ovh_iploadbalancing_http_farm_server testacc {
  service_name     = data.ovh_iploadbalancing.iplb.id
  farm_id = ovh_iploadbalancing_http_farm.testacc.id
  address = "10.0.0.11"
  status = "active"
  display_name = "testBackendB"
  port = 8080
  ssl = true
  backup = true
}
`
)

func TestAccIpLoadbalancingHttpFarmServerBasic(t *testing.T) {
	displayName := acctest.RandomWithPrefix(test_prefix)
	prefix := fmt.Sprintf(
		testAccIpLoadbalancingHttpFarmServerConfig_templ,
		os.Getenv("OVH_IPLB_SERVICE"),
		displayName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step0, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "display_name", "testBackendA"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "80"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "weight", "3"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "probe", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step1, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "display_name", "testBackendA"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "weight", "3"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "probe", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step2, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "display_name", "testBackendB"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "weight", "2"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "probe", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "backup", "true"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step3, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "inactive"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "80"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "probe", "false"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "backup", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step4, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "weight", "2"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "ssl", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "proxy_protocol_version", "v2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step5, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "weight", "1"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "ssl", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "backup", "false"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "proxy_protocol_version", "v1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccIpLoadbalancingHttpFarmServerConfig_step6, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "address", "10.0.0.11"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "status", "active"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "port", "8080"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "weight", "1"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "ssl", "true"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_http_farm_server.testacc", "backup", "true"),
				),
			},
		},
	})
}
