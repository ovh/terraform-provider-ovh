package ovh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

func init() {
	resource.AddTestSweepers("ovh_dedicated_ceph_user", &resource.Sweeper{
		Name: "ovh_dedicated_ceph_user",
		F:    testSweepDedicatedCephUser,
	})
}

func testSweepDedicatedCephUser(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}
	serviceName := os.Getenv("OVH_DEDICATED_CEPH")
	if serviceName == "" {
		log.Print("[DEBUG] No OVH_DEDICATED_CEPH envvar specified. nothing to sweep")
		return nil
	}

	var users []DedicatedCephUser
	endpoint := "/dedicated/ceph/" + url.PathEscape(serviceName) + "/user"
	if err := client.Get(endpoint, &users); err != nil {
		return fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}

	for _, user := range users {
		if !strings.HasPrefix(user.Name, test_prefix) {
			continue
		}

		log.Printf("[INFO] Deleting user %s on dedicated ceph %s", user.Name, serviceName)
		if err := testAccDeleteDedicatedCephUser(client, serviceName, user.Name); err != nil {
			return err
		}
	}

	return nil
}

// testAccDeleteDedicatedCephUser deletes a user through the API.
func testAccDeleteDedicatedCephUser(client *ovhwrap.Client, serviceName, userName string) error {
	endpoint := dedicatedCephUserEndpoint(serviceName, userName)

	ctx := context.Background()

	if err := waitDedicatedCephIdle(ctx, client, serviceName, dedicatedCephUserTimeout); err != nil {
		return err
	}

	var taskId string
	err := retryDedicatedCephChange(ctx, dedicatedCephUserTimeout, func() error {
		return client.Delete(endpoint, &taskId)
	})
	if err != nil {
		return fmt.Errorf("error calling Delete %s: %w", endpoint, err)
	}
	return waitDedicatedCephTask(ctx, client, serviceName, taskId, dedicatedCephUserTimeout)
}

const testAccDedicatedCephUserConfigPools = `
resource "ovh_dedicated_ceph_pool" "pool1" {
  service_name = "%[1]s"
  name         = "%[2]s-1"
}

resource "ovh_dedicated_ceph_pool" "pool2" {
  service_name = "%[1]s"
  name         = "%[2]s-2"
}
`

const testAccDedicatedCephUserConfigNoPermission = testAccDedicatedCephUserConfigPools + `
resource "ovh_dedicated_ceph_user" "user" {
  service_name = "%[1]s"
  name         = "%[2]s"
}
`

const testAccDedicatedCephUserConfigPermissions = testAccDedicatedCephUserConfigPools + `
resource "ovh_dedicated_ceph_user" "user" {
  service_name = "%[1]s"
  name         = "%[2]s"

  pool_permissions = [
    {
      pool_name = ovh_dedicated_ceph_pool.pool1.name
      read      = true
      write     = true
    },
    {
      pool_name = ovh_dedicated_ceph_pool.pool2.name
      read      = true
    },
  ]
}
`

const testAccDedicatedCephUserConfigPermissionsUpdated = testAccDedicatedCephUserConfigPools + `
resource "ovh_dedicated_ceph_user" "user" {
  service_name = "%[1]s"
  name         = "%[2]s"

  pool_permissions = [
    {
      pool_name  = ovh_dedicated_ceph_pool.pool2.name
      read       = true
      execute    = true
      class_read = true
    },
  ]
}
`

func TestAccDedicatedCephUser_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_DEDICATED_CEPH")
	name := acctest.RandomWithPrefix(test_prefix)
	resourceName := "ovh_dedicated_ceph_user.user"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDedicatedCeph(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDedicatedCephUserConfigPermissions, serviceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "service_name", serviceName),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "id", serviceName+"/"+name),
					resource.TestCheckResourceAttrSet(resourceName, "key"),
					resource.TestCheckResourceAttrSet(resourceName, "mon_caps"),
					resource.TestCheckResourceAttr(resourceName, "pool_permissions.#", "2"),
					// The configured order is kept, whatever the order of the API.
					resource.TestCheckResourceAttr(resourceName, "pool_permissions.0.pool_name", name+"-1"),
					resource.TestCheckResourceAttr(resourceName, "pool_permissions.1.pool_name", name+"-2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "pool_permissions.*", map[string]string{
						"pool_name":   name + "-1",
						"read":        "true",
						"write":       "true",
						"execute":     "false",
						"class_read":  "false",
						"class_write": "false",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "pool_permissions.*", map[string]string{
						"pool_name": name + "-2",
						"read":      "true",
						"write":     "false",
					}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     serviceName + "/" + name,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(testAccDedicatedCephUserConfigPermissionsUpdated, serviceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "pool_permissions.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "pool_permissions.*", map[string]string{
						"pool_name":  name + "-2",
						"read":       "true",
						"write":      "false",
						"execute":    "true",
						"class_read": "true",
					}),
				),
			},
			{
				Config: fmt.Sprintf(testAccDedicatedCephUserConfigNoPermission, serviceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "pool_permissions"),
				),
			},
		},
	})
}

const testAccDedicatedCephUserConfigEmptyPermissions = testAccDedicatedCephUserConfigPools + `
resource "ovh_dedicated_ceph_user" "user" {
  service_name     = "%[1]s"
  name             = "%[2]s"
  pool_permissions = []
}
`

const testAccDedicatedCephUserConfigPermissionWithoutFlag = testAccDedicatedCephUserConfigPools + `
resource "ovh_dedicated_ceph_user" "user" {
  service_name = "%[1]s"
  name         = "%[2]s"

  pool_permissions = [
    {
      pool_name = ovh_dedicated_ceph_pool.pool1.name
      read      = true
    },
    {
      pool_name = ovh_dedicated_ceph_pool.pool2.name
    },
  ]
}
`

// An empty list must be kept as configured, otherwise the apply ends with an inconsistent result,
// while a permission granting nothing, which the API refuses, must be refused when planning.
func TestAccDedicatedCephUser_emptyPermissions(t *testing.T) {
	serviceName := os.Getenv("OVH_DEDICATED_CEPH")
	name := acctest.RandomWithPrefix(test_prefix)
	resourceName := "ovh_dedicated_ceph_user.user"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDedicatedCeph(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Refused when planning, so first: the resources are destroyed with the
				// configuration of the last step, which must then be valid.
				Config:      fmt.Sprintf(testAccDedicatedCephUserConfigPermissionWithoutFlag, serviceName, name),
				ExpectError: regexp.MustCompile("Pool permission granting nothing"),
			},
			{
				Config: fmt.Sprintf(testAccDedicatedCephUserConfigEmptyPermissions, serviceName, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "pool_permissions.#", "0"),
				),
			},
		},
	})
}

func TestAccDedicatedCephUser_disappears(t *testing.T) {
	serviceName := os.Getenv("OVH_DEDICATED_CEPH")
	name := acctest.RandomWithPrefix(test_prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckDedicatedCeph(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// The user deleted outside of Terraform must be planned for creation again.
				Config: fmt.Sprintf(testAccDedicatedCephUserConfigNoPermission, serviceName, name),
				Check: func(*terraform.State) error {
					return testAccDeleteDedicatedCephUser(testAccOVHClient, serviceName, name)
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
