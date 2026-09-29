package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

func init() {
	resource.AddTestSweepers("ovh_dedicated_ceph_pool", &resource.Sweeper{
		Name:         "ovh_dedicated_ceph_pool",
		Dependencies: []string{"ovh_dedicated_ceph_user"},
		F:            testSweepDedicatedCephPool,
	})
}

func testSweepDedicatedCephPool(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	serviceName := os.Getenv("OVH_DEDICATED_CEPH")
	if serviceName == "" {
		log.Print("[DEBUG] No OVH_DEDICATED_CEPH envvar specified. nothing to sweep")
		return nil
	}

	var pools []DedicatedCephPool
	endpoint := "/dedicated/ceph/" + url.PathEscape(serviceName) + "/pool"
	if err := client.Get(endpoint, &pools); err != nil {
		return fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}

	for _, pool := range pools {
		if !strings.HasPrefix(pool.Name, test_prefix) {
			continue
		}

		log.Printf("[INFO] Deleting pool %s on dedicated ceph %s", pool.Name, serviceName)
		if err := testAccDeleteDedicatedCephPool(client, serviceName, pool.Name); err != nil {
			return err
		}
	}

	return nil
}

// testAccDeleteDedicatedCephPool deletes a pool through the API, opening its deletion window first.
func testAccDeleteDedicatedCephPool(client *ovhwrap.Client, serviceName, poolName string) error {
	ctx := context.Background()
	endpoint := dedicatedCephPoolEndpoint(serviceName, poolName)

	var taskId string
	err := retryDedicatedCephChange(ctx, dedicatedCephPoolTimeout, func() error {
		if err := client.Put(endpoint+"/allowDeletion", nil, nil); err != nil {
			return fmt.Errorf("error calling Put %s/allowDeletion: %w", endpoint, err)
		}
		if err := client.Delete(endpoint, &taskId); err != nil {
			return fmt.Errorf("error calling Delete %s: %w", endpoint, err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return waitDedicatedCephTask(ctx, client, serviceName, taskId, dedicatedCephPoolTimeout)
}

const testAccDedicatedCephPoolConfig = `
resource "ovh_dedicated_ceph_pool" "pool" {
  service_name = "%s"
  name         = "%s"
}
`

func TestAccDedicatedCephPool_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DEDICATED_CEPH")
	poolName := acctest.RandomWithPrefix(test_prefix)
	resourceName := "ovh_dedicated_ceph_pool.pool"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDedicatedCeph(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDedicatedCephPoolConfig, serviceName, poolName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "service_name", serviceName),
					resource.TestCheckResourceAttr(resourceName, "name", poolName),
					resource.TestCheckResourceAttr(resourceName, "id", serviceName+"/"+poolName),
					resource.TestCheckResourceAttrSet(resourceName, "pool_type"),
					resource.TestCheckResourceAttrSet(resourceName, "replica_count"),
					resource.TestCheckResourceAttrSet(resourceName, "min_active_replicas"),
					resource.TestCheckResourceAttrSet(resourceName, "backup"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     serviceName + "/" + poolName,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccDedicatedCephPool_disappears(t *testing.T) {
	serviceName := os.Getenv("OVH_DEDICATED_CEPH")
	poolName := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDedicatedCeph(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// The pool deleted outside of Terraform must be planned for creation again.
				Config: fmt.Sprintf(testAccDedicatedCephPoolConfig, serviceName, poolName),
				Check: func(*terraform.State) error {
					return testAccDeleteDedicatedCephPool(testAccOVHClient, serviceName, poolName)
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
