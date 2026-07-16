package ovh

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// testAccPreCheckCloudInstanceNet is the shared PreCheck for cloud instance
// tests that also exercise private networking (they require a vRack). Tier 2
// reuses this for the private-network / storage compositions.
func testAccPreCheckCloudInstanceNet(t *testing.T) {
	testAccPreCheckCredentials(t)
	checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_INSTANCE_REGION_TEST")
	checkEnvOrSkip(t, "OVH_VRACK_SERVICE_TEST")
	checkEnvOrSkip(t, "OVH_INSTANCE_FLAVOR_ID_TEST")
	checkEnvOrSkip(t, "OVH_INSTANCE_IMAGE_ID_TEST")
}

// testAccPreCheckCloudInstanceE2E is the shared PreCheck for full end-to-end
// instance compositions. It requires the same environment as
// testAccPreCheckCloudInstanceNet; the optional resize/rebuild image and
// flavor ids are read opportunistically by individual tests (with fallbacks).
func testAccPreCheckCloudInstanceE2E(t *testing.T) {
	testAccPreCheckCloudInstanceNet(t)
}

func init() {
	resource.AddTestSweepers("ovh_cloud_instance", &resource.Sweeper{
		Name: "ovh_cloud_instance",
		F:    testSweepCloudInstance,
	})
	resource.AddTestSweepers("ovh_cloud_instance_group", &resource.Sweeper{
		Name: "ovh_cloud_instance_group",
		// Instances may reference a placement group, so sweep instances first.
		Dependencies: []string{"ovh_cloud_instance"},
		F:            testSweepCloudInstanceGroup,
	})
}

// cloudInstanceSweepItem is a minimal projection of the /compute/instance API
// response, used only by the sweeper to identify test-owned instances.
type cloudInstanceSweepItem struct {
	Id         string `json:"id"`
	TargetSpec struct {
		Name string `json:"name"`
	} `json:"targetSpec"`
}

// hasTestPrefix reports whether name looks like a resource created by the
// acceptance tests (either the shared test_prefix or the legacy "test-" prefix).
func hasTestPrefix(name string) bool {
	return strings.HasPrefix(name, test_prefix) || strings.HasPrefix(name, "test-")
}

// testSweepCloudInstance deletes leftover test instances. It is best-effort and
// never fails the sweep run: listing/deletion errors are logged and swallowed.
func testSweepCloudInstance(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_CLOUD_PROJECT_SERVICE_TEST is not set. No instance to sweep")
		return nil
	}

	endpoint := fmt.Sprintf("/v2/publicCloud/project/%s/compute/instance", serviceName)

	instances := []cloudInstanceSweepItem{}
	if err := client.Get(endpoint, &instances); err != nil {
		log.Printf("[DEBUG] error listing instances to sweep (GET %s): %s", endpoint, err)
		return nil
	}

	for _, inst := range instances {
		if !hasTestPrefix(inst.TargetSpec.Name) {
			continue
		}

		log.Printf("[DEBUG] sweeping instance %s (%q) from project %s", inst.Id, inst.TargetSpec.Name, serviceName)
		deleteEndpoint := fmt.Sprintf("/v2/publicCloud/project/%s/compute/instance/%s", serviceName, inst.Id)
		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			if err := client.Delete(deleteEndpoint, nil); err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[DEBUG] error deleting instance %s: %s", inst.Id, err)
		}
	}

	return nil
}

// testSweepCloudInstanceGroup deletes leftover test instance groups. Best-effort.
func testSweepCloudInstanceGroup(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_CLOUD_PROJECT_SERVICE_TEST is not set. No instance group to sweep")
		return nil
	}

	endpoint := fmt.Sprintf("/v2/publicCloud/project/%s/compute/instanceGroup", serviceName)

	groups := []cloudInstanceSweepItem{}
	if err := client.Get(endpoint, &groups); err != nil {
		log.Printf("[DEBUG] error listing instance groups to sweep (GET %s): %s", endpoint, err)
		return nil
	}

	for _, grp := range groups {
		if !hasTestPrefix(grp.TargetSpec.Name) {
			continue
		}

		log.Printf("[DEBUG] sweeping instance group %s (%q) from project %s", grp.Id, grp.TargetSpec.Name, serviceName)
		deleteEndpoint := fmt.Sprintf("/v2/publicCloud/project/%s/compute/instanceGroup/%s", serviceName, grp.Id)
		err = resource.Retry(10*time.Minute, func() *resource.RetryError {
			if err := client.Delete(deleteEndpoint, nil); err != nil {
				return resource.RetryableError(err)
			}
			return nil
		})
		if err != nil {
			log.Printf("[DEBUG] error deleting instance group %s: %s", grp.Id, err)
		}
	}

	return nil
}
