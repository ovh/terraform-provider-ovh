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
)

func init() {
	resource.AddTestSweepers("ovh_dbaas_logs_output_opensearch_alias", &resource.Sweeper{
		Name: "ovh_dbaas_logs_output_opensearch_alias",
		F:    testSweepDbaasOutputOpensearchAlias,
	})
}

func testSweepDbaasOutputOpensearchAlias(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_DBAAS_LOGS_SERVICE_TEST is not set. No ovh_dbaas_output_graylog_stream to sweep")
		return nil
	}

	res := []string{}
	endpoint := fmt.Sprintf(
		"/dbaas/logs/%s/output/opensearch/alias",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}

	if len(res) == 0 {
		log.Print("[DEBUG] No ovh_dbaas_logs_output_opensearch_alias to sweep")
		return nil
	}

	for _, id := range res {
		endpoint := fmt.Sprintf(
			"/dbaas/logs/%s/output/opensearch/alias/%s/index",
			url.PathEscape(serviceName),
			url.PathEscape(id),
		)

		streams := []string{}
		if err := config.OVHClient.Get(endpoint, &streams); err != nil {
			return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
		}

		for _, s := range streams {
			res := &DbaasLogsOperation{}
			ctx := context.Background()
			err = retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
				log.Printf("[INFO] Will detach stream from opensearch alias: %s/%s", serviceName, id)
				if err := config.OVHClient.Delete(fmt.Sprintf("%s/%s", endpoint, url.PathEscape(s)), res); err != nil {
					return retry.RetryableError(err)
				}

				if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId); err != nil {
					return retry.RetryableError(err)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}

		log.Printf("[DEBUG] Will read dbaas logs output opensearch alias id: %s/%s", serviceName, id)

		endpoint = fmt.Sprintf(
			"/dbaas/logs/%s/output/opensearch/alias/%s",
			url.PathEscape(serviceName),
			url.PathEscape(id),
		)

		alias := &DbaasLogsOutputOpensearchAlias{}
		if err := config.OVHClient.Get(endpoint, &alias); err != nil {
			return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
		}

		if !strings.HasPrefix(alias.Description, test_prefix) {
			continue
		}

		res := &DbaasLogsOperation{}
		ctx := context.Background()
		err = retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
			log.Printf("[INFO] Will delete dbaas logs output graylog stream: %s/%s", serviceName, id)
			if err := config.OVHClient.Delete(endpoint, res); err != nil {
				return retry.RetryableError(err)
			}

			// Wait for operation status
			if _, err := waitForDbaasLogsOperation(ctx, config.OVHClient, serviceName, res.OperationId); err != nil {
				return retry.RetryableError(err)
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

func TestAccResourceDbaasLogsOutputOpensearchAlias_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	desc := acctest.RandomWithPrefix(test_prefix)
	suffix := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
		resource "ovh_dbaas_logs_output_opensearch_alias" "alias" {
			service_name = "%s"
			description  = "%s"
			suffix        = "%s"
		}
		`,
		serviceName,
		desc,
		suffix,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogs(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dbaas_logs_output_opensearch_alias.alias",
						"description",
						desc,
					),
				),
			},
		},
	})
}

func TestAccResourceDbaasLogsOutputOpensearchAlias_withIndex(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	desc := acctest.RandomWithPrefix(test_prefix)
	suffix := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
		resource "ovh_dbaas_logs_output_opensearch_index" "idx" {
			service_name = "%s"
			description  = "%s"
			suffix        = "%s"
			nb_shard = 1
		}

		resource "ovh_dbaas_logs_output_opensearch_alias" "aliasWithIdx" {
			service_name = "%s"
			description  = "%s"
			suffix        = "%s"
			indexes = [ovh_dbaas_logs_output_opensearch_index.idx.index_id]
		}
		`,
		serviceName,
		desc,
		suffix,
		serviceName,
		desc,
		suffix,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogs(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dbaas_logs_output_opensearch_alias.aliasWithIdx",
						"description",
						desc,
					),
					resource.TestCheckResourceAttrSet(
						"ovh_dbaas_logs_output_opensearch_alias.aliasWithIdx",
						"indexes.#",
					),
				),
			},
		},
	})
}

func TestAccResourceDbaasLogsOutputOpensearchAlias_withStream(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	desc := acctest.RandomWithPrefix(test_prefix)
	title := acctest.RandomWithPrefix(test_prefix)
	suffix := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
		resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
			service_name = "%s"
			description  = "%s"
			title        = "%s"
		}
		resource "ovh_dbaas_logs_output_opensearch_alias" "aliasWithStream" {
			service_name = "%s"
			description  = "%s"
			suffix        = "%s"
			streams = [ovh_dbaas_logs_output_graylog_stream.stream.stream_id]
		}
		`,
		serviceName,
		desc,
		title,
		serviceName,
		desc,
		suffix,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogs(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dbaas_logs_output_opensearch_alias.aliasWithStream",
						"description",
						desc,
					),
					resource.TestCheckResourceAttrSet(
						"ovh_dbaas_logs_output_opensearch_alias.aliasWithStream",
						"streams.#",
					),
				),
			},
		},
	})
}
