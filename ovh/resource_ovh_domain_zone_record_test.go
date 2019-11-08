package ovh

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"strings"
	"time"
)

func init() {
	resource.AddTestSweepers("ovh_domain_zone_record", &resource.Sweeper{
		Name: "ovh_domain_zone_record",
		F:    testSweepDomainZoneRecord,
	})
}

func testSweepDomainZoneRecord(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	zoneName := os.Getenv("OVH_ZONE")
	if zoneName == "" {
		return fmt.Errorf("OVH_ZONE must be set")
	}

	dz := &DomainZone{}

	if err := client.Get(fmt.Sprintf("/domain/zone/%s", zoneName), &dz); err != nil {
		return fmt.Errorf("Error calling /domain/zone/%s:\n\t %q", zoneName, err)
	}

	records := make([]int64, 0)
	if err := client.Get(fmt.Sprintf("/domain/zone/%s/record", zoneName), &records); err != nil {
		return fmt.Errorf("Error calling /domain/zone/%s:\n\t %q", zoneName, err)
	}

	if len(records) == 0 {
		log.Print("[DEBUG] No record to sweep")
		return nil
	}

	for _, rec := range records {
		record := &OvhDomainZoneRecord{}

		if err := client.Get(fmt.Sprintf("/domain/zone/%s/record/%v", zoneName, rec), &record); err != nil {
			return fmt.Errorf("Error calling /domain/zone/%s/record/%v:\n\t %q", zoneName, rec, err)
		}

		log.Printf("[DEBUG] record found %v", record)
		if !strings.HasPrefix(record.SubDomain, test_prefix) {
			continue
		}

		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting record %v", record)
			if err := client.Delete(fmt.Sprintf("/domain/zone/%s/record/%v", zoneName, rec), nil); err != nil {
				return resource.RetryableError(err)
			}
			// Successful delete
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Refreshing zone %s", zoneName)

		err := client.Post(
			fmt.Sprintf("/domain/zone/%s/refresh", zoneName),
			nil,
			nil,
		)

		if err != nil {
			return resource.RetryableError(fmt.Errorf("Error refresh OVH Zone: %s", err))
		}
		// Successful refresh
		return nil
	})

	return nil
}

func TestAccOvhDomainZoneRecord_Basic(t *testing.T) {
	var record OvhDomainZoneRecord
	zone := os.Getenv("OVH_ZONE")
	subdomain := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDomain(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOvhDomainZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOvhDomainZoneRecordConfig_A(zone, subdomain, "192.168.0.10", 3600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRecordExists("ovh_domain_zone_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", subdomain),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "target", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "ttl", "3600"),
				),
			},
		},
	})
}

func TestAccOvhDomainZoneRecord_Updated(t *testing.T) {
	record := OvhDomainZoneRecord{}
	zone := os.Getenv("OVH_ZONE")
	subdomain := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDomain(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOvhDomainZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOvhDomainZoneRecordConfig_A(zone, subdomain, "192.168.0.10", 3600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRecordExists("ovh_domain_zone_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", subdomain),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "target", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "ttl", "3600"),
				),
			},
			{
				Config: testAccCheckOvhDomainZoneRecordConfig_A(zone, subdomain, "192.168.0.11", 3600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRecordExists("ovh_domain_zone_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", subdomain),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "target", "192.168.0.11"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "ttl", "3600"),
				),
			},
			{
				Config: testAccCheckOvhDomainZoneRecordConfig_A(zone, fmt.Sprintf("%s2", subdomain),
					"192.168.0.11", 3600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRecordExists("ovh_domain_zone_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", fmt.Sprintf("%s2", subdomain)),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "target", "192.168.0.11"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "ttl", "3600"),
				),
			},
			{
				Config: testAccCheckOvhDomainZoneRecordConfig_A(zone, fmt.Sprintf("%s3", subdomain),
					"192.168.0.13", 3604),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRecordExists("ovh_domain_zone_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", fmt.Sprintf("%s3", subdomain)),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "target", "192.168.0.13"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "ttl", "3604"),
				),
			},
		},
	})
}

func TestAccOvhDomainZoneRecord_updateType(t *testing.T) {
	record := OvhDomainZoneRecord{}
	zone := os.Getenv("OVH_ZONE")
	subdomain := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDomain(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOvhDomainZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOvhDomainZoneRecordConfig_A(zone, subdomain, "192.168.0.1", 3600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRecordExists("ovh_domain_zone_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", subdomain),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "target", "192.168.0.1"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "fieldtype", "A"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "ttl", "3600"),
				),
			},
			{
				Config: testAccCheckOvhDomainZoneRecordConfig_CNAME(zone, subdomain, "google.com.", 3600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRecordExists("ovh_domain_zone_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", subdomain),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "target", "google.com."),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "fieldtype", "CNAME"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "ttl", "3600"),
				),
			},
		},
	})
}

func testAccCheckOvhDomainZoneRecordDestroy(s *terraform.State) error {
	provider := testAccProvider.Meta().(*Config)
	zone := os.Getenv("OVH_ZONE")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_domain_zone_record" {
			continue
		}

		resultRecord := OvhDomainZoneRecord{}
		err := provider.OVHClient.Get(
			fmt.Sprintf("/domain/zone/%s/record/%s", zone, rs.Primary.ID),
			&resultRecord,
		)

		if err == nil {
			return fmt.Errorf("Record still exists")
		}
	}

	return nil
}

func testAccCheckOvhDomainZoneRecordExists(n string, record *OvhDomainZoneRecord) resource.TestCheckFunc {
	zone := os.Getenv("OVH_ZONE")
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		provider := testAccProvider.Meta().(*Config)

		err := provider.OVHClient.Get(
			fmt.Sprintf("/domain/zone/%s/record/%s", zone, rs.Primary.ID),
			record,
		)

		if err != nil {
			return err
		}

		if strconv.Itoa(record.Id) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		return nil
	}
}

func testAccCheckOvhDomainZoneRecordConfig_A(zone, subdomain, target string, ttl int) string {
	return fmt.Sprintf(`
resource "ovh_domain_zone_record" "foobar" {
	zone = "%s"
	subdomain = "%s"
	target = "%s"
	fieldtype = "A"
	ttl = %d
}`, zone, subdomain, target, ttl)
}

func testAccCheckOvhDomainZoneRecordConfig_CNAME(zone, subdomain, target string, ttl int) string {
	return fmt.Sprintf(`
resource "ovh_domain_zone_record" "foobar" {
	zone = "%s"
	subdomain = "%s"
	target = "%s"
	fieldtype = "CNAME"
	ttl = %d
}`, zone, subdomain, target, ttl)
}
