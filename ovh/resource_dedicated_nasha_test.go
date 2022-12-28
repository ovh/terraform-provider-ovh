package ovh

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDedicatedNASHA(t *testing.T) {
	serviceName := os.Getenv("OVH_NASHA_SERVICE_TEST")
	partitionName := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck: func() {
			testAccPreCheckCredentials(t)
			checkEnvOrFail(t, "OVH_NASHA_SERVICE_TEST")
		},
		CheckDestroy: testAccNashaPartitionDestroy("ovh_dedicated_nasha_partition.testacc"),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "ovh_dedicated_nasha_partition" "testacc" {
						service_name = "%s"
						name = "%s"
						description = "test description"
						protocol = "NFS"
						size = 10
					}

					resource "ovh_dedicated_nasha_partition_snapshot" "testacc" {
						service_name = "${ovh_dedicated_nasha_partition.testacc.service_name}"
						partition_name = "${ovh_dedicated_nasha_partition.testacc.name}"
						type = "day-3"
					}

					resource "ovh_dedicated_nasha_partition_access" "testacc" {
						service_name = "${ovh_dedicated_nasha_partition.testacc.service_name}"
						partition_name = "${ovh_dedicated_nasha_partition.testacc.name}"
						ip = "127.0.0.1/32"
						type = "readonly"
					}
				`, serviceName, partitionName),
				Check: resource.ComposeTestCheckFunc(
					testAccNashaPartitionCheck("ovh_dedicated_nasha_partition.testacc"),
					resource.TestCheckResourceAttr("ovh_dedicated_nasha_partition.testacc", "name", partitionName),
					resource.TestCheckResourceAttr("ovh_dedicated_nasha_partition.testacc", "description", "test description"),
					resource.TestCheckResourceAttr("ovh_dedicated_nasha_partition.testacc", "protocol", "NFS"),
					resource.TestCheckResourceAttr("ovh_dedicated_nasha_partition.testacc", "size", "10"),

					testAccNashaPartitionSnapshotCheck("ovh_dedicated_nasha_partition_snapshot.testacc"),
					resource.TestCheckResourceAttr("ovh_dedicated_nasha_partition_snapshot.testacc", "partition_name", partitionName),
					resource.TestCheckResourceAttr("ovh_dedicated_nasha_partition_snapshot.testacc", "type", "day-3"),

					testAccNashaPartitionAccessCheck("ovh_dedicated_nasha_partition_access.testacc"),
					resource.TestCheckResourceAttr("ovh_dedicated_nasha_partition_access.testacc", "partition_name", partitionName),
					resource.TestCheckResourceAttr("ovh_dedicated_nasha_partition_access.testacc", "ip", "127.0.0.1/32"),
					resource.TestCheckResourceAttr("ovh_dedicated_nasha_partition_access.testacc", "type", "readonly"),
				),
			},
		},
	})
}

func testAccNashaPartitionCheck(n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		partitionResource, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		serviceName := partitionResource.Primary.Attributes["service_name"]
		partitionName := partitionResource.Primary.Attributes["name"]

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s", serviceName, partitionName)
		partitionResponse := &DedicatedNASHAPartition{}
		err := config.OVHClient.Get(endpoint, partitionResponse)
		if err != nil {
			return fmt.Errorf("calling GET %s:\n\t %s", endpoint, err.Error())
		}

		fmt.Printf("HA-NAS partition: %+v\n", partitionResponse)

		return nil
	}
}

func testAccNashaPartitionSnapshotCheck(n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		partitionResource, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		serviceName := partitionResource.Primary.Attributes["service_name"]
		partitionName := partitionResource.Primary.Attributes["partition_name"]
		snapshotType := partitionResource.Primary.Attributes["type"]

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s/snapshot/%s", serviceName, partitionName, snapshotType)
		partitionSnapshotResponse := &DedicatedNASHAPartitionSnapshot{}
		err := config.OVHClient.Get(endpoint, partitionSnapshotResponse)
		if err != nil {
			return fmt.Errorf("calling GET %s:\n\t %s", endpoint, err.Error())
		}

		fmt.Printf("HA-NAS partition snapshot: %+v\n", partitionSnapshotResponse)

		return nil
	}
}

func testAccNashaPartitionAccessCheck(n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		partitionResource, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		serviceName := partitionResource.Primary.Attributes["service_name"]
		partitionName := partitionResource.Primary.Attributes["partition_name"]
		accessIp := partitionResource.Primary.Attributes["ip"]

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s/access/%s", serviceName, partitionName, url.PathEscape(accessIp))
		partitionAccessResponse := &DedicatedNASHAPartitionAccess{}
		err := config.OVHClient.Get(endpoint, partitionAccessResponse)
		if err != nil {
			return fmt.Errorf("calling GET %s:\n\t %s", endpoint, err.Error())
		}

		fmt.Printf("HA-NAS partition access:%+v\n", partitionAccessResponse)

		return nil
	}
}

func testAccNashaPartitionDestroy(n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		nashaPartitionAccessResource, ok := state.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		serviceName := nashaPartitionAccessResource.Primary.Attributes["service_name"]
		partitionName := nashaPartitionAccessResource.Primary.Attributes["name"]

		config := testAccProvider.Meta().(*Config)
		endpoint := fmt.Sprintf("/dedicated/nasha/%s/partition/%s", serviceName, partitionName)

		err := config.OVHClient.Get(endpoint, nil)
		if err == nil {
			return fmt.Errorf("HA-NAS Partition (%s) still exists", partitionName)
		}

		return nil
	}
}
