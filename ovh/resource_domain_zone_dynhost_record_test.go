package ovh

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_domain_zone_dynhost_record", &resource.Sweeper{
		Name: "ovh_domain_zone_dynhost_record",
		F:    testSweepDomainZoneDynhostRecord,
	})
}

func testSweepDomainZoneDynhostRecord(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	zoneName := os.Getenv("OVH_ZONE_TEST")
	if zoneName == "" {
		log.Print("[DEBUG] OVH_ZONE_TEST is not set. No zone dynhost records to sweep")
		return nil
	}

	records := make([]string, 0)
	if err := client.Get(fmt.Sprintf("/domain/zone/%s/dynhost/record", url.PathEscape(zoneName)), &records); err != nil {
		return fmt.Errorf("Error calling /domain/zone/%s/dynhost/record:\n\t %q", zoneName, err)
	}

	if len(records) == 0 {
		log.Print("[DEBUG] No record to sweep")
		return nil
	}

	for _, rec := range records {
		record := &DomainZoneDynhostRecordModel{}

		if err := client.Get(fmt.Sprintf("/domain/zone/%s/dynhost/record/%s", url.PathEscape(zoneName), url.PathEscape(rec)), &record); err != nil {
			return fmt.Errorf("Error calling /domain/zone/%s/dynhost/record/%s:\n\t %q", zoneName, rec, err)
		}

		log.Printf("[DEBUG] record found %v", record)

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting record %v", record)
			if err := client.Delete(fmt.Sprintf("/domain/zone/%s/dynhost/record/%v", zoneName, rec), nil); err != nil {
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

func TestAccDomainZoneDynHostRecord_Basic(t *testing.T) {
	zone := os.Getenv("OVH_ZONE_TEST")
	subdomain := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDomain(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// provider should send an error if no zone is given
			{
				Config:      testAccOvhDomainZoneDynhostRecordNoZoneNameConfig(subdomain, "1.2.3.4"),
				ExpectError: regexp.MustCompile(`The argument "zone_name" is required, but no definition was found.`),
			},
			// provider should send an error if zone does not exist
			{
				Config:      testAccOvhDomainZoneDynhostRecordConfig("non-existing-domain.com", subdomain, "1.2.3.4"),
				ExpectError: regexp.MustCompile("Client::NotFound: \"This service does\nnot exist\""),
			},
			// provider should send an error if password IP is not valid
			{
				Config:      testAccOvhDomainZoneDynhostRecordConfig(zone, subdomain, "1.2.3.4.5"),
				ExpectError: regexp.MustCompile("'Invalid\nip'"),
			},
			// resource creation
			{
				Config: testAccOvhDomainZoneDynhostRecordConfig(zone, subdomain, "1.2.3.4"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_record.this", "ip", "1.2.3.4"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_record.this", "sub_domain", subdomain),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_record.this", "zone", zone),
				),
			},
			// resource ip update
			{
				Config: testAccOvhDomainZoneDynhostRecordConfig(zone, subdomain, "5.6.7.8"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_record.this", "ip", "5.6.7.8"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_record.this", "sub_domain", subdomain),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_record.this", "zone", zone),
				),
			},
			// resource subdomain update
			{
				Config: testAccOvhDomainZoneDynhostRecordConfig(zone, "other-subdomain", "5.6.7.8"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_record.this", "ip", "5.6.7.8"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_record.this", "sub_domain", "other-subdomain"),
					resource.TestCheckResourceAttr("ovh_domain_zone_dynhost_record.this", "zone", zone),
				),
			},
		},
	})
}

func testAccOvhDomainZoneDynhostRecordNoZoneNameConfig(subDomain, ip string) string {
	return fmt.Sprintf(`	
resource "ovh_domain_zone_dynhost_record" "this" {
  sub_domain = %q
  ip         = %q
}`, subDomain, ip)
}

func testAccOvhDomainZoneDynhostRecordConfig(zoneName, subDomain, ip string) string {
	return fmt.Sprintf(`	
resource "ovh_domain_zone_dynhost_record" "this" {
  zone_name  = %q
  sub_domain = %q
  ip         = %q
}`, zoneName, subDomain, ip)
}
