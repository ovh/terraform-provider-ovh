package ovh

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("ovh_iam_resource_group", &resource.Sweeper{
		Name: "ovh_iam_resource_group",
		F:    testSweepMeIdentityResourceGroup,
	})
}

func testSweepMeIdentityResourceGroup(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	var groups []IamResourceGroup
	if err := client.Get("/v2/iam/resourceGroup", &groups); err != nil {
		return fmt.Errorf("Error calling /v2/iam/resourceGroup:\n\t %q", err)
	}

	if len(groups) == 0 {
		log.Print("[DEBUG] No identity users to sweep")
		return nil
	}
	for _, resGrp := range groups {
		if !strings.HasPrefix(resGrp.Name, test_prefix) {
			continue
		}

		if resGrp.ReadOnly {
			continue
		}
		log.Printf("[DEBUG] IAM resource group found %s: %s", resGrp.Name, resGrp.ID)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			log.Printf("[INFO] Deleting IAM resource group %s: %s", resGrp.Name, resGrp.ID)
			if err := client.Delete(fmt.Sprintf("/v2/iam/resourceGroup/%s", resGrp.ID), nil); err != nil {
				return resource.RetryableError(err)
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

func TestAccIamResourceGroupResource_basic(t *testing.T) {
	resourceGroupName1 := acctest.RandomWithPrefix(test_prefix)
	resourceGroupName2 := acctest.RandomWithPrefix(test_prefix)

	resource1 := "urn:v1:eu:resource:vrack:" + os.Getenv("OVH_VRACK_SERVICE_TEST")
	resource2 := "urn:v1:eu:resource:vps:" + os.Getenv("OVH_VPS")

	config := fmt.Sprintf(
		testAccIamResourceGroupResourceConfig_preSetup,
		resourceGroupName1,
		resource1,
		resourceGroupName2,
		resource1,
		resource2,
	)

	checks := append(
		checkIamResourceGroupResourceAttr("ovh_iam_resource_group.resource_group_1", resourceGroupName1, resource1),
		checkIamResourceGroupResourceAttr("ovh_iam_resource_group.resource_group_2", resourceGroupName2, resource1, resource2)...,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckIamResourceGroup(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  resource.ComposeTestCheckFunc(checks...),
			}, {
				ResourceName:      "ovh_iam_resource_group.resource_group_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIamResourceGroupResourceConfig_preSetup = `
resource "ovh_iam_resource_group" "resource_group_1" {
	name        = "%s"
	resources   = ["%s"]
}

resource "ovh_iam_resource_group" "resource_group_2" {
	name        = "%s"
	resources   = ["%s", "%s"]
}
`
