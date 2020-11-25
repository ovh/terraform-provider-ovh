package ovh

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	testAccIpLoadbalancingTcpFarmConfig = `
data "ovh_iploadbalancing" "iplb" {
  service_name = "%s"
}
resource "ovh_iploadbalancing_tcp_farm" "testfarm" {
  service_name     = data.ovh_iploadbalancing.iplb.id
  display_name     = "%s"
  port             = "%d"
  zone             = "%s"
  balance 		   = "roundrobin"
  probe {
        interval = 30
        type = "oco"
  }
}
`
)

func init() {
	resource.AddTestSweepers("ovh_iploadbalancing_tcp_farm", &resource.Sweeper{
		Name: "ovh_iploadbalancing_tcp_farm",
		Dependencies: []string{
			"ovh_iploadbalancing_tcp_farm_server",
		},
		F: testSweepIploadbalancingTcpFarm,
	})
}

func testSweepIploadbalancingTcpFarm(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	iplb := os.Getenv("OVH_IPLB_SERVICE")
	if iplb == "" {
		log.Print("[DEBUG] OVH_IPLB_SERVICE is not set. No iploadbalancing_vrack_network to sweep")
		return nil
	}

	farms := make([]int64, 0)
	if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm", iplb), &farms); err != nil {
		return fmt.Errorf("Error calling /ipLoadbalancing/%s/tcp/farm:\n\t %q", iplb, err)
	}

	if len(farms) == 0 {
		log.Print("[DEBUG] No tcp farm to sweep")
		return nil
	}

	for _, f := range farms {
		farm := &IpLoadbalancingFarm{}

		if err := client.Get(fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%d", iplb, f), &farm); err != nil {
			return fmt.Errorf("Error calling /ipLoadbalancing/%s/tcp/farm/%d:\n\t %q", iplb, f, err)
		}

		if !strings.HasPrefix(*farm.DisplayName, test_prefix) {
			continue
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			if err := client.Delete(fmt.Sprintf("/ipLoadbalancing/%s/tcp/farm/%d", iplb, f), nil); err != nil {
				return resource.RetryableError(err)
			}
			// Successful delete
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func TestAccIpLoadbalancingTcpFarmBasicCreate(t *testing.T) {
	displayName1 := acctest.RandomWithPrefix(test_prefix)
	displayName2 := acctest.RandomWithPrefix(test_prefix)
	config1 := fmt.Sprintf(
		testAccIpLoadbalancingTcpFarmConfig,
		os.Getenv("OVH_IPLB_SERVICE"),
		displayName1,
		12345,
		"all",
	)
	config2 := fmt.Sprintf(
		testAccIpLoadbalancingTcpFarmConfig,
		os.Getenv("OVH_IPLB_SERVICE"),
		displayName2,
		12346,
		"all",
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckIpLoadbalancing(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "display_name", displayName1),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "zone", "all"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "port", "12345"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "probe.0.interval", "30"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "display_name", displayName2),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "zone", "all"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "port", "12346"),
					resource.TestCheckResourceAttr("ovh_iploadbalancing_tcp_farm.testfarm", "probe.0.interval", "30"),
				),
			},
		},
	})
}
