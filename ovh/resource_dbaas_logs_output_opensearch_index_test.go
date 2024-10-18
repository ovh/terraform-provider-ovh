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
	resource.AddTestSweepers("ovh_dbaas_logs_output_opensearch_index", &resource.Sweeper{
		Name: "ovh_dbaas_logs_output_opensearch_index",
		F:    testSweepDbaasOutputOpensearchIndex,
	})
}

func testSweepDbaasOutputOpensearchIndex(region string) error {
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
		"/dbaas/logs/%s/output/opensearch/index",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}

	if len(res) == 0 {
		log.Print("[DEBUG] No ovh_dbaas_logs_output_opensearch_index to sweep")
		return nil
	}

	for _, id := range res {
		log.Printf("[DEBUG] Will read dbaas logs output opensearch index id: %s/%s", serviceName, id)

		endpoint := fmt.Sprintf(
			"/dbaas/logs/%s/output/opensearch/index/%s",
			url.PathEscape(serviceName),
			url.PathEscape(id),
		)

		Index := &DbaasLogsOutputOpensearchIndex{}
		if err := config.OVHClient.Get(endpoint, &Index); err != nil {
			return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
		}

		if !strings.HasPrefix(Index.Description, test_prefix) {
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

func TestAccResourceDbaasLogsOutputOpensearchIndex_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	desc := acctest.RandomWithPrefix(test_prefix)
	suffix := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
		resource "ovh_dbaas_logs_output_opensearch_index" "index" {
			service_name = "%s"
			description  = "%s"
			suffix        = "%s"
			nb_shard = 1
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
						"ovh_dbaas_logs_output_opensearch_index.index",
						"description",
						desc,
					),
				),
			},
		},
	})
}
