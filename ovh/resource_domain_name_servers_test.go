package ovh

import (
	"fmt"
	"github.com/ovh/go-ovh/ovh"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_domain_name_servers", &resource.Sweeper{
		Name: "ovh_domain_name_servers",
		F:    testSweepOvhDomainNameServers,
	})
}

func testSweepOvhDomainNameServers(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	domain := os.Getenv("OVH_DOMAIN_TEST")
	if domain == "" {
		log.Print("[DEBUG] OVH_DOMAIN_TEST is not set. No domain to sweep")
		return nil
	}

	// Check we have access to the domain
	err = client.Get(fmt.Sprintf("/domain/%s", domain), nil)

	if err != nil {
		if err.(*ovh.APIError).Code == 404 {
			log.Printf("[DEBUG] OVH domain %s does not exist. No domain to sweep", domain)
			return nil
		}
		return fmt.Errorf("error getting domain: %s", err)
	}

	return nil
}

func TestAccDomainNameServers_Basic(t *testing.T) {
	domain := os.Getenv("OVH_DOMAIN_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDomain(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// provider shall send an error if the TTL is less than 60
			{
				Config: testAccCheckOvhDomainNameServersConfig(domain, []string{"dns104.ovh.net", "ns104.ovh.net"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_domain_name_servers.foobar", "domain", domain),
				),
			},
		},
	})
}

func testAccCheckOvhDomainNameServersConfig(domain string, nameServers []string) string {
	return `
resource "ovh_domain_name_servers" "foobar" {
	  domain = "` + domain + `"
	` + testAccCheckOvhDomainNameServersConfigNameServers(nameServers) + `
}
`
}

func testAccCheckOvhDomainNameServersConfigNameServers(nameServers []string) string {
	var config string
	for _, nameServer := range nameServers {
		config += `
	servers {
		host = "` + nameServer + `"
	}
	`
	}
	return config
}
