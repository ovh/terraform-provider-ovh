package ovh

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/go-ovh/ovh"
)

func init() {
	resource.AddTestSweepers("ovh_domain_name_servers", &resource.Sweeper{
		Name: "ovh_domain_name_servers",
		F:    testSweepDomainNameServers,
	})
}

func testSweepDomainNameServers(region string) error {
	client, err := sharedClientForRegion(region)

	if err != nil {
		return fmt.Errorf("error getting client:\n\t%v", err)
	}

	domainName := os.Getenv("OVH_ZONE_TEST")

	if domainName == "" {
		log.Println("[DEBUG] OVH_ZONE_TEST is not set. No domain to sweep")
		return nil
	}

	domainNameServerTypeOpts := &DomainNameServerTypeOpts{
		NameServerType: "hosted",
	}

	log.Printf("[DEBUG] Will reset domain name servers to defaults: %s\n", domainName)

	endpoint := fmt.Sprintf("/domain/%s", url.PathEscape(domainName))

	if err := client.Put(endpoint, domainNameServerTypeOpts, nil); err != nil {
		// Ignore error "Your domaine nsType is already 'hosted'" as it means there is no sweep to perform
		if ovhErr, ok := err.(ovh.APIError); ok && ovhErr.Code == http.StatusConflict && ovhErr.Class == "Client::Conflict::DomMsuUnknownError" {
			return nil
		}

		return fmt.Errorf("error calling PUT on %s:\n\t%v", endpoint, err)
	}

	return nil
}

func TestAccDomainNameServers_Basic(t *testing.T) {
	domainName := os.Getenv("OVH_ZONE_TEST")

	nameServer1Host := os.Getenv("OVH_DOMAIN_NS1_HOST_TEST")
	nameServer1Ip := os.Getenv("OVH_DOMAIN_NS1_IP_TEST")
	nameServer2Host := os.Getenv("OVH_DOMAIN_NS2_HOST_TEST")
	nameServer3Host := os.Getenv("OVH_DOMAIN_NS3_HOST_TEST")

	fmt.Printf("[INFO] Will update test domain name servers: %s\n", domainName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDomain(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOvhDomainNameServersDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckOvhDomainNameServersConfig_Invalid(domainName, nameServer1Host, nameServer1Ip),
				ExpectError: regexp.MustCompile(`2 "servers" blocks are required`),
			},
			{
				Config: testAccCheckOvhDomainNameServersConfig(domainName, nameServer1Host, nameServer1Ip, nameServer2Host),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainNameServersCurrent("ovh_domain_name_servers.test", nameServer1Host, nameServer1Ip, nameServer2Host),
					resource.TestCheckResourceAttr("ovh_domain_name_servers.test", "domain", domainName),
					resource.TestCheckTypeSetElemNestedAttrs("ovh_domain_name_servers.test", "servers.*", map[string]string{
						"host": nameServer1Host,
						"ip":   nameServer1Ip,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("ovh_domain_name_servers.test", "servers.*", map[string]string{
						"host": nameServer2Host,
					}),
				),
			},
			{
				Config: testAccCheckOvhDomainNameServersConfig(domainName, nameServer2Host, "", nameServer3Host),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainNameServersCurrent("ovh_domain_name_servers.test", nameServer2Host, "", nameServer3Host),
					resource.TestCheckResourceAttr("ovh_domain_name_servers.test", "domain", domainName),
					resource.TestCheckTypeSetElemNestedAttrs("ovh_domain_name_servers.test", "servers.*", map[string]string{
						"host": nameServer2Host,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("ovh_domain_name_servers.test", "servers.*", map[string]string{
						"host": nameServer3Host,
					}),
				),
			},
		},
	})
}

func testAccCheckOvhDomainNameServersConfig_Invalid(domainName string, nameServer1Host string, nameServer1Ip string) string {
	return fmt.Sprintf(`
resource "ovh_domain_name_servers" "invalid" {
	domain = "%s"

	servers {
		host = "%s"
		ip = "%s"
	}
}
`, domainName, nameServer1Host, nameServer1Ip)
}

func testAccCheckOvhDomainNameServersConfig(domainName string, nameServer1Host string, nameServer1Ip string, nameServer2Host string) string {
	return fmt.Sprintf(`
resource "ovh_domain_name_servers" "test" {
	domain = "%s"

	servers {
		host = "%s"
		ip = "%s"
	}

	servers {
		host = "%s"
		ip = ""
    }
}
`, domainName, nameServer1Host, nameServer1Ip, nameServer2Host)
}

func testAccCheckOvhDomainNameServersCurrent(resourceName string, nameServer1Host string, nameServer1Ip string, nameServer2Host string) resource.TestCheckFunc {
	domainName := os.Getenv("OVH_ZONE_TEST")

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("no resource found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf(`no "id" is set: %s`, resourceName)
		}

		provider := testAccProvider.Meta().(*Config)

		var nameservers []int
		endpoint := fmt.Sprintf("/domain/%s/nameServer", url.PathEscape(domainName))

		if err := provider.OVHClient.Get(endpoint, &nameservers); err != nil {
			return fmt.Errorf("error while calling GET on %s:\n\t%v", endpoint, err)
		}

		found := 0

		for _, nameServerId := range nameservers {
			responseData := DomainNameServer{}
			endpoint := fmt.Sprintf("/domain/%s/nameServer/%d", url.PathEscape(domainName), nameServerId)

			if err := provider.OVHClient.Get(endpoint, &responseData); err != nil {
				return fmt.Errorf("error while calling GET on %s:\n\t%v", endpoint, err)
			}

			if (responseData.Host == nameServer1Host && responseData.Ip == nameServer1Ip) || (responseData.Host == nameServer2Host && responseData.Ip == "") {
				found += 1
			}
		}

		if found != 2 {
			return fmt.Errorf("domain name servers not configured properly")
		}

		return nil
	}
}

func testAccCheckOvhDomainNameServersDestroy(s *terraform.State) error {
	provider := testAccProvider.Meta().(*Config)
	domainName := os.Getenv("OVH_ZONE_TEST")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_domain_name_servers" {
			continue
		}

		resultRecord := DomainNameServerTypeOpts{}
		endpoint := fmt.Sprintf("/domain/%s", url.PathEscape(domainName))

		if err := provider.OVHClient.Get(endpoint, &resultRecord); err != nil {
			return fmt.Errorf("error while calling GET on %s:\n\t%v", endpoint, err)
		}

		if resultRecord.NameServerType == "external" {
			return fmt.Errorf(`domain name servers type is still "external"`)
		}
	}

	return nil
}
