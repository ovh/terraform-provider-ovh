package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

var testAccCloudProjectContainerRegistryOIDCConfig = `
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

resource "ovh_cloud_project_containerregistry_oidc" "my-oidc" {
	service_name = ovh_cloud_project_containerregistry.registry.service_name
	registry_id  = ovh_cloud_project_containerregistry.registry.id

	delete_users = true
	oidc_name = "name"
	oidc_endpoint = "%s"
	oidc_client_id = "%s"
	oidc_client_secret = "clientSecret"
	oidc_scope = "openid,profile,email,offline_access"
	oidc_group_filter = "groupFilter"
	oidc_groups_claim = "groupsClaim"
	oidc_admin_group = "adminGroup"
	oidc_verify_cert = "true"
	oidc_auto_onboard = "true"
	oidc_user_claim = "userClaim"
}
`

func TestAccCloudProjectContainerRegistryOIDC_full(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	registryName := acctest.RandomWithPrefix(test_prefix)
	region := os.Getenv("OVH_CLOUD_PROJECT_CONTAINERREGISTRY_REGION_TEST")
	oidcEndpoint := os.Getenv("OVH_CLOUD_PROJECT_CONTAINERREGISTRY_OIDC_ENDPOINT_TEST")
	resourceName := "ovh_cloud_project_containerregistry_oidc.my-oidc"

	config := fmt.Sprintf(
		testAccCloudProjectContainerRegistryOIDCConfig,
		serviceName,
		region,
		registryName,
		oidcEndpoint,
		"clientID",
	)

	configUpdated := fmt.Sprintf(
		testAccCloudProjectContainerRegistryOIDCConfig,
		serviceName,
		region,
		registryName,
		oidcEndpoint,
		"clientIDModified",
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckContainerRegistryOIDC(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "delete_users", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_name", "name"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_endpoint", oidcEndpoint),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_client_id", "clientID"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_client_secret", "clientSecret"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_scope", "openid,profile,email,offline_access"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_group_filter", "groupFilter"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_groups_claim", "groupsClaim"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_admin_group", "adminGroup"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_verify_cert", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_auto_onboard", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_user_claim", "userClaim"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_users", "oidc_client_secret"},
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					return fmt.Sprintf("%s/%s", serviceName, state.RootModule().Resources[resourceName].Primary.Attributes["registry_id"]), nil
				},
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "delete_users", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_name", "name"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_endpoint", oidcEndpoint),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_client_id", "clientIDModified"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_client_secret", "clientSecret"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_scope", "openid,profile,email,offline_access"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_group_filter", "groupFilter"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_groups_claim", "groupsClaim"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_admin_group", "adminGroup"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_verify_cert", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_auto_onboard", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_project_containerregistry_oidc.my-oidc", "oidc_user_claim", "userClaim"),
				),
			},
			{
				Config:       configUpdated,
				Destroy:      true,
				ResourceName: resourceName,
			},
		},
	})
}
