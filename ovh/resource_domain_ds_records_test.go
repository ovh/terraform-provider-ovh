package ovh

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/go-ovh/ovh"
)

func init() {
	resource.AddTestSweepers("ovh_domain_ds_records", &resource.Sweeper{
		Name: "ovh_domain_ds_records",
		F:    testSweepDomainDsRecords,
	})
}

func testSweepDomainDsRecords(region string) error {
	client, err := sharedClientForRegion(region)

	if err != nil {
		return fmt.Errorf("error getting client:\n\t%v", err)
	}

	domainName := os.Getenv("OVH_ZONE_TEST")

	if domainName == "" {
		log.Println("[DEBUG] OVH_ZONE_TEST is not set. No domain to sweep")
		return nil
	}

	domainDsRecordsUpdateOpts := DomainDsRecordsUpdateOpts{
		DsRecords: []DomainDsRecord{},
	}

	log.Printf("[DEBUG] Will delete all domain DS records: %s\n", domainName)

	endpoint := fmt.Sprintf("/domain/%s/dsRecord", url.PathEscape(domainName))

	if err := client.Post(endpoint, domainDsRecordsUpdateOpts, nil); err != nil {
		if ovhErr, ok := err.(*ovh.APIError); ok &&
			ovhErr.Code == http.StatusForbidden &&
			ovhErr.Class == "Client::Forbidden::DomDoaForbidden" {
			log.Printf("[DEBUG] Cannot update DS records on the given domain, skipping sweeper")
			return nil
		}
		return fmt.Errorf("error calling POST on %s:\n\t%v", endpoint, err)
	}

	return nil
}

func TestAccDomainDsRecords_Basic(t *testing.T) {
	domainName := os.Getenv("OVH_ZONE_TEST")
	resourceName := "ovh_domain_ds_records.test"

	recordAlgorithm := os.Getenv("OVH_DOMAIN_DS_RECORD_ALGORITHM_TEST")
	recordFlags := "KEY_SIGNING_KEY"
	recordPublicKey := os.Getenv("OVH_DOMAIN_DS_RECORD_PUBLIC_KEY_TEST")
	recordTag, _ := strconv.ParseInt(os.Getenv("OVH_DOMAIN_DS_RECORD_TAG_TEST"), 10, 0)

	fmt.Printf("[INFO] Will update test domain DS records: %s\n", domainName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDomain(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOvhDomainDsRecordsDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckOvhDomainDsRecordsConfig_Invalid(domainName),
				ExpectError: regexp.MustCompile(`1 "ds_records" blocks are required`),
			},
			{
				Config: testAccCheckOvhDomainDsRecordsConfig(resourceName, domainName, recordAlgorithm, recordFlags, recordPublicKey, int(recordTag)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOvhDomainDsRecordsCurrent(resourceName, recordAlgorithm, recordFlags, recordPublicKey, int(recordTag)),
					resource.TestCheckResourceAttr(resourceName, "domain", domainName),
					resource.TestCheckResourceAttr(resourceName, "ds_records.0.algorithm", recordAlgorithm),
					resource.TestCheckResourceAttr(resourceName, "ds_records.0.flags", recordFlags),
					resource.TestCheckResourceAttr(resourceName, "ds_records.0.public_key", recordPublicKey),
					resource.TestCheckResourceAttr(resourceName, "ds_records.0.tag", fmt.Sprintf("%d", recordTag)),
				),
			},
		},
	})
}

func testAccCheckOvhDomainDsRecordsConfig_Invalid(domainName string) string {
	return fmt.Sprintf(`
resource "ovh_domain_ds_records" "invalid" {
	domain = "%s"
}
`, domainName)
}

func testAccCheckOvhDomainDsRecordsConfig(resourceFullName string, domainName string, recordAlgorithm string, recordFlags string, recordPublicKey string, recordTag int) string {
	resourceType, resourceName, _ := strings.Cut(resourceFullName, ".")

	return fmt.Sprintf(`
resource "%s" "%s" {
	domain = "%s"

	ds_records {
		algorithm = "%s"
		flags = "%s"
		public_key = "%s"
		tag = %d
	}
}
`, resourceType, resourceName, domainName, recordAlgorithm, recordFlags, recordPublicKey, recordTag)
}

func testAccReadOvhDomainDsRecords(domainName string) ([]DomainDsRecord, error) {
	provider := testAccProvider.Meta().(*Config)

	responseData := &[]int{}
	endpoint := fmt.Sprintf("/domain/%s/dsRecord", url.PathEscape(domainName))

	if err := provider.OVHClient.Get(endpoint, &responseData); err != nil {
		return nil, fmt.Errorf("error calling GET on %s:\n\t%v", endpoint, err)
	}

	var domainDsRecordList []DomainDsRecord

	for _, dsRecordId := range *responseData {
		responseData := &DomainDsRecord{}
		endpoint := fmt.Sprintf("/domain/%s/dsRecord/%d", url.PathEscape(domainName), dsRecordId)

		if err := provider.OVHClient.Get(endpoint, &responseData); err != nil {
			return nil, fmt.Errorf("error calling GET on %s:\n\t%v", endpoint, err)
		}

		domainDsRecordList = append(domainDsRecordList, *responseData)
	}

	return domainDsRecordList, nil
}

func testAccCheckOvhDomainDsRecordsCurrent(resourceName string, recordAlgorithm string, recordFlags string, recordPublicKey string, recordTag int) resource.TestCheckFunc {
	domainName := os.Getenv("OVH_ZONE_TEST")

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("no resource found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set: %s", resourceName)
		}

		domainDsRecordList, err := testAccReadOvhDomainDsRecords(domainName)

		if err != nil {
			return err
		}

		if len(domainDsRecordList) != 1 ||
			domainDsRecordList[0].Algorithm != DsRecordAlgorithmValuesMap[recordAlgorithm] ||
			domainDsRecordList[0].Flags != DsRecordFlagValuesMap[recordFlags] ||
			domainDsRecordList[0].PublicKey != recordPublicKey ||
			domainDsRecordList[0].Tag != recordTag {
			return fmt.Errorf("domain DS records not configured properly")
		}

		return nil
	}
}

func testAccCheckOvhDomainDsRecordsDestroy(s *terraform.State) error {
	domainName := os.Getenv("OVH_ZONE_TEST")

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ovh_domain_ds_records" {
			continue
		}

		domainDsRecordList, err := testAccReadOvhDomainDsRecords(domainName)

		if err != nil {
			return err
		}

		if len(domainDsRecordList) != 0 {
			return fmt.Errorf("domain DS records still exists")
		}
	}

	return nil
}
