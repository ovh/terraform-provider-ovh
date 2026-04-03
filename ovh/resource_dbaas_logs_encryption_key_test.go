package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func init() {
	resource.AddTestSweepers("ovh_dbaas_logs_encryption_key", &resource.Sweeper{
		Name: "ovh_dbaas_logs_encryption_key",
		F:    testSweepDbaasLogsEncryptionKey,
	})
}

func testSweepDbaasLogsEncryptionKey(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_DBAAS_LOGS_SERVICE_TEST is not set. No ovh_dbaas_logs_encryption_key to sweep")
		return nil
	}

	res := []string{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/encryptionKey",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}

	if len(res) == 0 {
		log.Print("[DEBUG] No ovh_dbaas_logs_encryption_key to sweep")
		return nil
	}

	for _, id := range res {
		endpoint := fmt.Sprintf(
			"/dbaas/logs/%s/encryptionKey/%s",
			url.PathEscape(serviceName),
			url.PathEscape(id),
		)

		key := &DbaasLogsEncryptionKeyModel{}
		if err := config.OVHClient.Get(endpoint, &key); err != nil {
			return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
		}

		if !strings.HasPrefix(key.Title.ValueString(), test_prefix) {
			continue
		}

		opRes := &DbaasLogsOperation{}
		ctx := context.Background()
		err = retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
			log.Printf("[INFO] Will delete dbaas logs encryption key: %s/%s", serviceName, id)
			if err := config.OVHClient.Delete(endpoint, opRes); err != nil {
				return retry.RetryableError(err)
			}

			if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, opRes.OperationId); err != nil {
				return retry.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func testAccPreCheckDbaasLogsEncryptionKey(t *testing.T) {
	testAccPreCheckDbaasLogs(t)
	checkEnvOrSkip(t, "OVH_DBAAS_LOGS_ENCRYPTION_KEY_CONTENT_TEST")
	checkEnvOrSkip(t, "OVH_DBAAS_LOGS_ENCRYPTION_KEY_FINGERPRINT_TEST")
}

func TestAccResourceDbaasLogsEncryptionKey_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	pgpContent := os.Getenv("OVH_DBAAS_LOGS_ENCRYPTION_KEY_CONTENT_TEST")
	pgpFingerprint := os.Getenv("OVH_DBAAS_LOGS_ENCRYPTION_KEY_FINGERPRINT_TEST")
	title := acctest.RandomWithPrefix(test_prefix)
	titleUpdated := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
		resource "ovh_dbaas_logs_encryption_key" "key" {
			service_name = "%s"
			title        = "%s"
			content      = <<-EOT
%s
EOT
			fingerprint  = "%s"
		}
	`, serviceName, title, pgpContent, pgpFingerprint)

	configUpdated := fmt.Sprintf(`
		resource "ovh_dbaas_logs_encryption_key" "key" {
			service_name = "%s"
			title        = "%s"
			content      = <<-EOT
%s
EOT
			fingerprint  = "%s"
		}
	`, serviceName, titleUpdated, pgpContent, pgpFingerprint)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDbaasLogsEncryptionKey(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_dbaas_logs_encryption_key.key", "service_name", serviceName),
					resource.TestCheckResourceAttr("ovh_dbaas_logs_encryption_key.key", "title", title),
					resource.TestCheckResourceAttrSet("ovh_dbaas_logs_encryption_key.key", "encryption_key_id"),
					resource.TestCheckResourceAttr("ovh_dbaas_logs_encryption_key.key", "fingerprint", pgpFingerprint),
					resource.TestCheckResourceAttrSet("ovh_dbaas_logs_encryption_key.key", "created_at"),
					resource.TestCheckResourceAttrSet("ovh_dbaas_logs_encryption_key.key", "id"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_dbaas_logs_encryption_key.key", "title", titleUpdated),
				),
			},
			{
				Config:            config,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"content",
				},
				ResourceName:      "ovh_dbaas_logs_encryption_key.key",
				ImportStateIdFunc: testAccDbaasLogsEncryptionKeyImportId("ovh_dbaas_logs_encryption_key.key"),
			},
		},
	})
}

func testAccDbaasLogsEncryptionKeyImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		ds, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		return fmt.Sprintf(
			"%s/%s",
			ds.Primary.Attributes["service_name"],
			ds.Primary.Attributes["encryption_key_id"],
		), nil
	}
}

func TestAccDataSourceDbaasLogsEncryptionKey_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	pgpContent := os.Getenv("OVH_DBAAS_LOGS_ENCRYPTION_KEY_CONTENT_TEST")
	pgpFingerprint := os.Getenv("OVH_DBAAS_LOGS_ENCRYPTION_KEY_FINGERPRINT_TEST")
	title := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
		resource "ovh_dbaas_logs_encryption_key" "key" {
			service_name = "%s"
			title        = "%s"
			content      = <<-EOT
%s
EOT
			fingerprint  = "%s"
		}

		data "ovh_dbaas_logs_encryption_key" "key" {
			service_name      = ovh_dbaas_logs_encryption_key.key.service_name
			title             = ovh_dbaas_logs_encryption_key.key.title
		}
	`, serviceName, title, pgpContent, pgpFingerprint)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDbaasLogsEncryptionKey(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_dbaas_logs_encryption_key.key", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_dbaas_logs_encryption_key.key", "title", title),
					resource.TestCheckResourceAttrSet("data.ovh_dbaas_logs_encryption_key.key", "encryption_key_id"),
					resource.TestCheckResourceAttr("data.ovh_dbaas_logs_encryption_key.key", "fingerprint", pgpFingerprint),
					resource.TestCheckResourceAttrSet("data.ovh_dbaas_logs_encryption_key.key", "created_at"),
				),
			},
		},
	})
}
