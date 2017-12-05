package ovh

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOVHRecord_Basic(t *testing.T) {
	var record Record
	zone := os.Getenv("OVH_ZONE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOVHRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOVHRecordConfig_basic, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOVHRecordExists("ovh_domain_zone_record.foobar", &record),
					testAccCheckOVHRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", "terraform"),
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

func TestAccOVHRecord_Updated(t *testing.T) {
	record := Record{}
	zone := os.Getenv("OVH_ZONE")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOVHRecordDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOVHRecordConfig_basic, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOVHRecordExists("ovh_domain_zone_record.foobar", &record),
					testAccCheckOVHRecordAttributes(&record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", "terraform"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "target", "192.168.0.10"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "ttl", "3600"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOVHRecordConfig_new_value_1, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOVHRecordExists("ovh_domain_zone_record.foobar", &record),
					testAccCheckOVHRecordAttributesUpdated_1(&record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", "terraform"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "target", "192.168.0.11"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "ttl", "3600"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOVHRecordConfig_new_value_2, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOVHRecordExists("ovh_domain_zone_record.foobar", &record),
					testAccCheckOVHRecordAttributesUpdated_2(&record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", "terraform2"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "target", "192.168.0.11"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "ttl", "3600"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckOVHRecordConfig_new_value_3, zone),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOVHRecordExists("ovh_domain_zone_record.foobar", &record),
					testAccCheckOVHRecordAttributesUpdated_3(&record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", "terraform3"),
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

func testAccCheckOVHRecordDestroy(s *terraform.State) error {
	provider := testAccProvider.Meta().(*Config)
	zone := os.Getenv("OVH_ZONE")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_domain_zone_record" {
			continue
		}

		resultRecord := Record{}
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

func testAccCheckOVHRecordExists(n string, record *Record) resource.TestCheckFunc {
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

func testAccCheckOVHRecordAttributes(record *Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Target != "192.168.0.10" {
			return fmt.Errorf("Bad content: %#v", record)
		}

		return nil
	}
}

func testAccCheckOVHRecordAttributesUpdated_1(record *Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Target != "192.168.0.11" {
			return fmt.Errorf("Bad content: %#v", record)
		}

		return nil
	}
}

func testAccCheckOVHRecordAttributesUpdated_2(record *Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Target != "192.168.0.11" {
			return fmt.Errorf("Bad content: %#v", record)
		}

		if record.SubDomain != "terraform2" {
			return fmt.Errorf("Bad content: %#v", record)
		}

		return nil
	}
}

func testAccCheckOVHRecordAttributesUpdated_3(record *Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if record.Target != "192.168.0.13" {
			return fmt.Errorf("Bad content: %#v", record)
		}

		if record.SubDomain != "terraform3" {
			return fmt.Errorf("Bad content: %#v", record)
		}

		if record.Ttl != 3604 {
			return fmt.Errorf("Bad content: %#v", record)
		}
		return nil
	}
}

const testAccCheckOVHRecordConfig_basic = `
resource "ovh_domain_zone_record" "foobar" {
	zone = "%s"
	subdomain = "terraform"
	target = "192.168.0.10"
	fieldtype = "A"
	ttl = 3600
}`

const testAccCheckOVHRecordConfig_new_value_1 = `
resource "ovh_domain_zone_record" "foobar" {
	zone = "%s"
	subdomain = "terraform"
	target = "192.168.0.11"
	fieldtype = "A"
	ttl = 3600
}
`
const testAccCheckOVHRecordConfig_new_value_2 = `
resource "ovh_domain_zone_record" "foobar" {
	zone = "%s"
	subdomain = "terraform2"
	target = "192.168.0.11"
	fieldtype = "A"
	ttl = 3600
}
`
const testAccCheckOVHRecordConfig_new_value_3 = `
resource "ovh_domain_zone_record" "foobar" {
	zone = "%s"
	subdomain = "terraform3"
	target = "192.168.0.13"
	fieldtype = "A"
	ttl = 3604
}`
