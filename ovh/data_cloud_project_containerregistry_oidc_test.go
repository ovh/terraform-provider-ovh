package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCloudProjectContainerRegistryOIDCDataSource_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_CONTAINERREGISTRY_REGION_TEST")
	registryName := acctest.RandomWithPrefix(test_prefix)
	oidcEndpoint := os.Getenv("OVH_CLOUD_PROJECT_CONTAINERREGISTRY_OIDC_ENDPOINT_TEST")

	config := fmt.Sprintf(
		testAccCloudProjectContainerRegistryOIDCDataSourceConfig,
		serviceName,
		region,
		registryName,
		oidcEndpoint,
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
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_containerregistry_oidc.oidcData", "oidc_name", "name"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_containerregistry_oidc.oidcData", "oidc_endpoint", oidcEndpoint),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_containerregistry_oidc.oidcData", "oidc_client_id", "clientID"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_containerregistry_oidc.oidcData", "oidc_scope", "openid,profile,email,offline_access"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_containerregistry_oidc.oidcData", "oidc_groups_claim", "groupsClaim"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_containerregistry_oidc.oidcData", "oidc_admin_group", "adminGroup"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_containerregistry_oidc.oidcData", "oidc_verify_cert", "true"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_containerregistry_oidc.oidcData", "oidc_auto_onboard", "true"),
					resource.TestCheckResourceAttr(
						"data.ovh_cloud_project_containerregistry_oidc.oidcData", "oidc_user_claim", "userClaim"),
				),
			},
		},
	})
}

var testAccCloudProjectContainerRegistryOIDCDataSourceConfig = `
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

resource "ovh_cloud_project_containerregistry_oidc" "oidc" {
	service_name = ovh_cloud_project_containerregistry.registry.service_name
	registry_id  = ovh_cloud_project_containerregistry.registry.id
	
	oidc_name = "name"
	oidc_endpoint = "%s"
	oidc_client_id = "clientID"
	oidc_client_secret = "clientSecret"
	oidc_scope = "openid,profile,email,offline_access"
	oidc_groups_claim = "groupsClaim"
	oidc_admin_group = "adminGroup"
	oidc_verify_cert = "true"
	oidc_auto_onboard = "true"
	oidc_user_claim = "userClaim"

	depends_on = [
		ovh_cloud_project_containerregistry.registry
	]
}

data "ovh_cloud_project_containerregistry_oidc" "oidcData" {
    service_name = ovh_cloud_project_containerregistry.registry.service_name
    registry_id = ovh_cloud_project_containerregistry.registry.id

	depends_on = [
		ovh_cloud_project_containerregistry_oidc.oidc
	]
}
`
