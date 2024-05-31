package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccCloudProjectContainerRegistryIPRestrictionsRegistryConfig = `
data "ovh_cloud_project_capabilities_containerregistry_filter" "registryCap" {
  service_name = "%s"
  plan_name    = "SMALL"
  region       = "%s"
}

resource "ovh_cloud_project_containerregistry" "registry" {
  service_name = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.service_name
  plan_id      = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.id
  name         = "%s"
  region       = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.region
}

resource "ovh_cloud_project_containerregistry_ip_restrictions_registry" "my-registry-iprestrictions" {
  service_name = ovh_cloud_project_containerregistry.registry.service_name
  registry_id  = ovh_cloud_project_containerregistry.registry.id
	
  ip_restrictions = [
    {
      ip_block = "121.121.121.121/32"
      description = "my first ip restriction"
    },
    {
      ip_block = "122.122.122.122/32"
      description = "my second ip description"
    }
  ]
}
`

const testAccCloudProjectContainerRegistryIPRestrictionsRegistryConfigUpdated = `
data "ovh_cloud_project_capabilities_containerregistry_filter" "registryCap" {
  service_name = "%s"
  plan_name    = "SMALL"
  region       = "%s"
}

resource "ovh_cloud_project_containerregistry" "registry" {
  service_name = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.service_name
  plan_id      = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.id
  name         = "%s"
  region       = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.region
}

resource "ovh_cloud_project_containerregistry_ip_restrictions_registry" "my-registry-iprestrictions" {
  service_name = ovh_cloud_project_containerregistry.registry.service_name
  registry_id  = ovh_cloud_project_containerregistry.registry.id

  ip_restrictions = [
    {
      ip_block = "121.121.121.121/32"
      description = "my new first ip restriction"
    },
    {
      ip_block = "122.122.122.122/32"
      description = "my second ip description"
    }
  ]
}
`

const testAccCloudProjectContainerRegistryIPRestrictionsRegistryReorderIpRestrictions = `
data "ovh_cloud_project_capabilities_containerregistry_filter" "registryCap" {
  service_name = "%s"
  plan_name    = "SMALL"
  region       = "%s"
}

resource "ovh_cloud_project_containerregistry" "registry" {
  service_name = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.service_name
  plan_id      = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.id
  name         = "%s"
  region       = data.ovh_cloud_project_capabilities_containerregistry_filter.registryCap.region
}

resource "ovh_cloud_project_containerregistry_ip_restrictions_registry" "my-registry-iprestrictions" {
  service_name = ovh_cloud_project_containerregistry.registry.service_name
  registry_id  = ovh_cloud_project_containerregistry.registry.id
	
  ip_restrictions = [
    {
      ip_block = "122.122.122.122/32"
      description = "my second ip description"
    },
    {
      ip_block = "121.121.121.121/32"
      description = "my new first ip restriction"
    }
  ]
}
`

func TestAccCloudProjectContainerRegistryIPRestrictionsRegistry_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	registryName := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_CONTAINERREGISTRY_REGION_TEST")
	resourceName := "ovh_cloud_project_containerregistry_ip_restrictions_registry.my-registry-iprestrictions"

	config := fmt.Sprintf(
		testAccCloudProjectContainerRegistryIPRestrictionsRegistryConfig,
		serviceName,
		region,
		registryName,
	)

	configUpdated := fmt.Sprintf(
		testAccCloudProjectContainerRegistryIPRestrictionsRegistryConfigUpdated,
		serviceName,
		region,
		registryName,
	)

	configUpdatedIpRestrictions := fmt.Sprintf(
		testAccCloudProjectContainerRegistryIPRestrictionsRegistryReorderIpRestrictions,
		serviceName,
		region,
		registryName,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckContainerRegistry(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_ip_restrictions_registry.my-registry-iprestrictions", "ip_restrictions.0.ip_block", "121.121.121.121/32"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_ip_restrictions_registry.my-registry-iprestrictions", "ip_restrictions.0.description", "my first ip restriction"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_ip_restrictions_registry.my-registry-iprestrictions", "ip_restrictions.0.ip_block", "121.121.121.121/32"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_ip_restrictions_registry.my-registry-iprestrictions", "ip_restrictions.0.description", "my new first ip restriction"),
				),
			},
			{
				// We check that a change in the ip restrictions order does not affect the plan
				Config:             configUpdatedIpRestrictions,
				ExpectNonEmptyPlan: false,
			},
			{
				Config:       configUpdatedIpRestrictions,
				Destroy:      true,
				ResourceName: resourceName,
			},
		},
	})
}
