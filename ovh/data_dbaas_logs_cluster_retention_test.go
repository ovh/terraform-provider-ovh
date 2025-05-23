package ovh

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceDbaasLogsClusterRetention_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	clusterId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_ID")
	retentionId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_RETENTION_ID")

	config := fmt.Sprintf(`
		data "ovh_dbaas_logs_cluster_retention" "ret" {
			service_name = "%s"
			cluster_id   = "%s"
			retention_id = "%s"
		}`,
		serviceName,
		clusterId,
		retentionId,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDbaasLogsClusterRetention(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster_retention.ret",
						"service_name",
						serviceName,
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster_retention.ret",
						"duration",
						"P1Y",
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster_retention.ret",
						"is_supported",
						"true",
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster_retention.ret",
						"retention_type",
						"LOGS_INDEXING",
					),
				),
			},
		},
	})
}

func TestAccDataSourceDbaasLogsClusterRetention_by_duration(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	clusterId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_ID")
	retentionId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_RETENTION_ID")

	config := fmt.Sprintf(`
		data "ovh_dbaas_logs_cluster_retention" "ret" {
			service_name = "%s"
			cluster_id   = "%s"
			duration     = "P1Y"
		}`,
		serviceName,
		clusterId,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDbaasLogsClusterRetention(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster_retention.ret",
						"service_name",
						serviceName,
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster_retention.ret",
						"retention_id",
						retentionId,
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster_retention.ret",
						"is_supported",
						"true",
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster_retention.ret",
						"retention_type",
						"LOGS_INDEXING",
					),
				),
			},
		},
	})
}

func TestAccDataSourceDbaasLogsClusterRetention_by_duration_and_type(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	clusterId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_ID")
	retentionId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_RETENTION_ID")

	config := fmt.Sprintf(`
		data "ovh_dbaas_logs_cluster_retention" "ret" {
			service_name   = "%s"
			cluster_id     = "%s"
			duration       = "P1Y"
			retention_type = "LOGS_INDEXING"
		}`,
		serviceName,
		clusterId,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDbaasLogsClusterRetention(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster_retention.ret",
						"service_name",
						serviceName,
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster_retention.ret",
						"retention_id",
						retentionId,
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster_retention.ret",
						"is_supported",
						"true",
					),
					resource.TestCheckResourceAttr(
						"data.ovh_dbaas_logs_cluster_retention.ret",
						"retention_type",
						"LOGS_INDEXING",
					),
				),
			},
		},
	})
}

func TestAccDataSourceDbaasLogsClusterRetention_by_duration_not_found(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	clusterId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_ID")

	config := fmt.Sprintf(`
		data "ovh_dbaas_logs_cluster_retention" "ret" {
			service_name = "%s"
			cluster_id   = "%s"
			duration     = "P1000Y"
		}`,
		serviceName,
		clusterId,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDbaasLogsCluster(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("no retention was found with duration P1000Y and type LOGS_INDEXING"),
			},
		},
	})
}

func TestAccDataSourceDbaasLogsClusterRetention_by_duration_and_type_not_found(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	clusterId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_ID")

	config := fmt.Sprintf(`
		data "ovh_dbaas_logs_cluster_retention" "ret" {
			service_name   = "%s"
			cluster_id     = "%s"
			duration       = "P1Y"
			retention_type = "METRICS_TENANT"
		}`,
		serviceName,
		clusterId,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDbaasLogsClusterRetention(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("no retention was found with duration P1Y and type METRICS_TENANT"),
			},
		},
	})
}

func TestAccDataSourceDbaasLogsClusterRetention_missing_params(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	clusterId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_ID")

	config := fmt.Sprintf(`
		data "ovh_dbaas_logs_cluster_retention" "ret" {
			service_name = "%s"
			cluster_id   = "%s"
		}`,
		serviceName,
		clusterId,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDbaasLogsCluster(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("missing retention_id or duration"),
			},
		},
	})
}

func TestAccDataSourceDbaasLogsClusterRetention_invalid_params(t *testing.T) {
	serviceName := os.Getenv("OVH_DBAAS_LOGS_SERVICE_TEST")
	clusterId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_ID")
	retentionId := os.Getenv("OVH_DBAAS_LOGS_CLUSTER_RETENTION_ID")

	config := fmt.Sprintf(`
		data "ovh_dbaas_logs_cluster_retention" "ret" {
			service_name   = "%s"
			cluster_id     = "%s"
			retention_id   = "%s"
			retention_type = "LOGS_INDEXING"
		}`,
		serviceName,
		clusterId,
		retentionId,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDbaasLogsCluster(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`Attribute \"retention_type\" cannot be specified when \"retention_id\"`),
			},
		},
	})
}
