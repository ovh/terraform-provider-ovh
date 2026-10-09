package ovh

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	ovhsdk "github.com/ovh/go-ovh/ovh"
	"go.uber.org/ratelimit"

	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
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

	zoneName := os.Getenv("OVH_ZONE_TEST")
	if zoneName == "" {
		log.Print("[DEBUG] OVH_ZONE_TEST is not set. No zone to sweep")
		return nil
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

func TestAccDomainZoneRecord_Basic(t *testing.T) {
	var record OvhDomainZoneRecord
	zone := os.Getenv("OVH_ZONE_TEST")
	subdomain := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDomain(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOvhDomainZoneRecordDestroy,
		Steps: []resource.TestStep{
			// provider shall send an error if the TTL is less than 60
			{
				Config:      testAccCheckOvhDomainZoneRecordConfig_CNAME(zone, subdomain, "google.com.", 10),
				ExpectError: regexp.MustCompile(`must be either equal to 0 or, greater than or equal to 60`),
			},
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
				Config: testAccCheckOvhDomainZoneRecordConfig_A(zone, subdomain, "192.168.0.11", 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRecordExists("ovh_domain_zone_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", subdomain),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "target", "192.168.0.11"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "ttl", "0"),
				),
			},
			{
				Config: testAccCheckOvhDomainZoneRecordConfig_A_noTTL(zone, subdomain, "192.168.0.12"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainZoneRecordExists("ovh_domain_zone_record.foobar", &record),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "subdomain", subdomain),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "zone", zone),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "target", "192.168.0.12"),
					resource.TestCheckResourceAttr(
						"ovh_domain_zone_record.foobar", "ttl", "0"),
				),
			},
		},
	})
}

func TestAccDomainZoneRecord_Updated(t *testing.T) {
	record := OvhDomainZoneRecord{}
	zone := os.Getenv("OVH_ZONE_TEST")
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

func TestAccDomainZoneRecord_updateType(t *testing.T) {
	record := OvhDomainZoneRecord{}
	zone := os.Getenv("OVH_ZONE_TEST")
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

func TestAccDomainZoneRecord_EmptyPlanForTXT(t *testing.T) {
	zone := os.Getenv("OVH_ZONE_TEST")
	subdomain := acctest.RandomWithPrefix(test_prefix)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDomain(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOvhDomainZoneRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckOvhDomainZoneRecordConfig_TXT(zone, subdomain, "test", 3600),
			},
			{
				Config: testAccCheckOvhDomainZoneRecordConfig_TXT(zone, subdomain, "test", 3600),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func testAccCheckOvhDomainZoneRecordDestroy(s *terraform.State) error {
	provider := testAccProvider.Meta().(*Config)
	zone := os.Getenv("OVH_ZONE_TEST")

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
	zone := os.Getenv("OVH_ZONE_TEST")
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

		if strconv.FormatInt(record.Id, 10) != rs.Primary.ID {
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

func testAccCheckOvhDomainZoneRecordConfig_A_noTTL(zone, subdomain, target string) string {
	return fmt.Sprintf(`
resource "ovh_domain_zone_record" "foobar" {
	zone = "%s"
	subdomain = "%s"
	target = "%s"
	fieldtype = "A"
}`, zone, subdomain, target)
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

func testAccCheckOvhDomainZoneRecordConfig_TXT(zone, subdomain, target string, ttl int) string {
	return fmt.Sprintf(`
resource "ovh_domain_zone_record" "foobar" {
	zone = "%s"
	subdomain = "%s"
	target = "%s"
	fieldtype = "TXT"
	ttl = %d
}`, zone, subdomain, target, ttl)
}

// newRefreshTestConfig returns a Config whose OVH client talks to the given test
// server. Dummy credentials are supplied so go-ovh signs requests (which also
// triggers the /auth/time call served by the test server) without reading any
// ambient OVH configuration.
func newRefreshTestConfig(t *testing.T, serverURL string) *Config {
	t.Helper()
	client, err := ovhsdk.NewClient(serverURL, "appKey", "appSecret", "consumerKey")
	if err != nil {
		t.Fatalf("failed to create go-ovh client: %s", err)
	}
	return &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}
}

// refreshDeniedServer serves the record endpoints and rejects the zone refresh
// with a 403, mimicking credentials that lack the dnsZone:apiovh:refresh
// permission.
func refreshDeniedServer(t *testing.T, zone string, refreshCalls *int32) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/auth/time":
			fmt.Fprint(w, "1700000000")
		case r.URL.Path == "/domain/zone/"+zone+"/refresh" && r.Method == http.MethodPost:
			atomic.AddInt32(refreshCalls, 1)
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, `{"message":"This call has not been granted"}`)
		case r.URL.Path == "/domain/zone/"+zone+"/record" && r.Method == http.MethodPost:
			fmt.Fprint(w, `{"id":42,"zone":"`+zone+`","subDomain":"test","fieldType":"A","target":"192.0.2.1","ttl":3600}`)
		case r.URL.Path == "/domain/zone/"+zone+"/record/42" && r.Method == http.MethodDelete:
			fmt.Fprint(w, `null`)
		case r.URL.Path == "/domain/zone/"+zone+"/record/42" && r.Method == http.MethodGet:
			// Only reached if the refresh error is swallowed and the read runs anyway.
			fmt.Fprint(w, `{"id":42,"zone":"`+zone+`","subDomain":"test","fieldType":"A","target":"192.0.2.1","ttl":3600}`)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func testRecordResourceData(t *testing.T, zone string) *schema.ResourceData {
	t.Helper()
	return schema.TestResourceDataRaw(t, resourceOvhDomainZoneRecord().Schema, map[string]interface{}{
		"zone":      zone,
		"subdomain": "test",
		"fieldtype": "A",
		"target":    "192.0.2.1",
		"ttl":       3600,
	})
}

// A record is only served by the DNS servers once the zone is refreshed, so a
// denied refresh must fail the apply rather than be swallowed into a log line.
func TestDomainZoneRecordCreateFailsWhenRefreshDenied(t *testing.T) {
	const zone = "example.com"
	var refreshCalls int32

	server := refreshDeniedServer(t, zone, &refreshCalls)
	defer server.Close()

	d := testRecordResourceData(t, zone)

	err := resourceOvhDomainZoneRecordCreate(d, newRefreshTestConfig(t, server.URL))
	if err == nil {
		t.Fatal("expected create to fail when the zone refresh is denied, got nil")
	}
	if !strings.Contains(err.Error(), "refresh after record creation failed") {
		t.Errorf("error should name the failing operation, got: %s", err)
	}
	if got := atomic.LoadInt32(&refreshCalls); got != 1 {
		t.Errorf("expected exactly 1 refresh call, got %d", got)
	}
}

// On delete the refresh runs before the ID is cleared, so a denied refresh
// leaves the resource in state for the next apply to retry.
func TestDomainZoneRecordDeleteFailsWhenRefreshDeniedAndKeepsID(t *testing.T) {
	const zone = "example.com"
	var refreshCalls int32

	server := refreshDeniedServer(t, zone, &refreshCalls)
	defer server.Close()

	d := testRecordResourceData(t, zone)
	d.SetId("42")

	err := resourceOvhDomainZoneRecordDelete(d, newRefreshTestConfig(t, server.URL))
	if err == nil {
		t.Fatal("expected delete to fail when the zone refresh is denied, got nil")
	}
	if !strings.Contains(err.Error(), "refresh after record deletion failed") {
		t.Errorf("error should name the failing operation, got: %s", err)
	}
	if d.Id() != "42" {
		t.Errorf("id must be preserved so the next apply retries the delete, got %q", d.Id())
	}
	if got := atomic.LoadInt32(&refreshCalls); got != 1 {
		t.Errorf("expected exactly 1 refresh call, got %d", got)
	}
}
