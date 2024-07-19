package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceDbaasLogsOutputGraylogStream_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	title := acctest.RandomWithPrefix(test_prefix)
	desc := acctest.RandomWithPrefix(test_prefix)

	config := fmt.Sprintf(`
		resource "ovh_dbaas_logs_output_graylog_stream" "stream" {
			service_name = "%s"
			title        = "%s"
			description  = "%s"
		}

		data "ovh_dbaas_logs_output_graylog_stream" "stream" {
			service_name = ovh_dbaas_logs_output_graylog_stream.stream.service_name
			title        = ovh_dbaas_logs_output_graylog_stream.stream.title
		}`,
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
						"data.ovh_dbaas_logs_output_graylog_stream.stream",
						"description",
						desc,
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_output_graylog_stream.stream",
						"title",
						title,
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_dbaas_logs_output_graylog_stream.stream",
						"write_token",
					),
				),
			},
		},
	})
}

func TestAccDataSourceDbaasLogsOutputGraylogStream_with_retention(t *testing.T) {
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

		data "ovh_dbaas_logs_output_graylog_stream" "stream" {
			service_name = ovh_dbaas_logs_output_graylog_stream.stream.service_name
			title        = ovh_dbaas_logs_output_graylog_stream.stream.title
		}`,
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
						"data.ovh_dbaas_logs_output_graylog_stream.stream",
						"description",
						desc,
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_output_graylog_stream.stream",
						"title",
						title,
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_dbaas_logs_output_graylog_stream.stream",
						"write_token",
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_output_graylog_stream.stream",
						"retention_id",
						retentionId,
					),
				),
			},
		},
	})
}
