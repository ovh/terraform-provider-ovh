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
	resource.AddTestSweepers("ovh_dbaas_logs_output_graylog_stream", &resource.Sweeper{
		Name: "ovh_dbaas_logs_output_graylog_stream",
		F:    testSweepDbaasOutputGraylogStream,
	})
}

func testSweepDbaasOutputGraylogStream(region string) error {
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
		"/dbaas/logs/%s/output/graylog/stream",
		url.PathEscape(serviceName),
	)
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
	}

	if len(res) == 0 {
		log.Print("[DEBUG] No ovh_dbaas_output_graylog_stream to sweep")
		return nil
	}

	for _, id := range res {
		log.Printf("[DEBUG] Will read dbaas logs output graylog stream id : %s/%s", serviceName, id)

		endpoint := fmt.Sprintf(
			"/dbaas/logs/%s/output/graylog/stream/%s",
			url.PathEscape(serviceName),
			url.PathEscape(id),
		)

		stream := &DbaasLogsOutputGraylogStream{}
		if err := config.OVHClient.Get(endpoint, &stream); err != nil {
			return fmt.Errorf("Error calling Get %s:\n\t %q", endpoint, err)
		}

		if !strings.HasPrefix(stream.Title, test_prefix) {
			continue
		}

		res := &DbaasLogsOperation{}
		ctx := context.Background()
		err = retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
			log.Printf("[INFO] Will delete dbaas logs output graylog stream : %s/%s", serviceName, id)
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

func TestAccResourceDbaasLogsOutputGraylogStream_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	title := acctest.RandomWithPrefix(test_prefix)
	desc := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
		resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
			service_name = "%s"
			title        = "%s"
			description  = "%s"
		}
		`,
		serviceName,
		title,
		desc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogs(t) },

		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dbaas_logs_output_graylog_stream.stream",
						"description",
						desc,
					),
					resource.TestCheckResourceAttr(
						"ovh_dbaas_logs_output_graylog_stream.stream",
						"title",
						title,
					),
					resource.TestCheckResourceAttrSet(
						"ovh_dbaas_logs_output_graylog_stream.stream",
						"write_token",
					),
				),
			},
		},
	})
}

func TestAccResourceDbaasLogsOutputGraylogStream_with_retention(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	clusterId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_ID")
	retentionId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_RETENTION_ID")
	title := acctest.RandomWithPrefix(test_prefix)
	desc := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
		data "ovh_dbaas_logs_cluster_retention" "retention" {
			service_name = "%s"
			cluster_id   = "%s"
			duration     = "P1Y"
		}

		resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
			service_name = "%s"
			title        = "%s"
			description  = "%s"
			retention_id = data.ovh_dbaas_logs_cluster_retention.retention.retention_id
		}
		`,
		serviceName,
		clusterId,
		serviceName,
		title,
		desc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheckDbaasLogsClusterRetention(t) },

		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ovh_dbaas_logs_output_graylog_stream.stream",
						"description",
						desc,
					),
					resource.TestCheckResourceAttr(
						"ovh_dbaas_logs_output_graylog_stream.stream",
						"title",
						title,
					),
					resource.TestCheckResourceAttrSet(
						"ovh_dbaas_logs_output_graylog_stream.stream",
						"write_token",
					),
					resource.TestCheckResourceAttr(
						"ovh_dbaas_logs_output_graylog_stream.stream",
						"retention_id",
						retentionId,
					),
				),
			},
		},
	})
}
