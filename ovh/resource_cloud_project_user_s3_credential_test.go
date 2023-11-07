package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccCloudProjectUserS3CredentialConfig_basic = `
resource "ovh_cloud_project_user" "user" {
 service_name = "%s"
 description  = "my user for acceptance tests"
}

resource "ovh_cloud_project_user_s3_credential" "s3_cred" {
 service_name  = ovh_cloud_project_user.user.service_name
 user_id       = ovh_cloud_project_user.user.id
}
`

func TestAccCloudProjectUserS3Credential_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	config := fmt.Sprintf(testAccCloudProjectUserS3CredentialConfig_basic, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_user_s3_credential.s3_cred", "service_name"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_user_s3_credential.s3_cred", "user_id"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_user_s3_credential.s3_cred", "internal_user_id"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_user_s3_credential.s3_cred", "access_key_id"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_user_s3_credential.s3_cred", "secret_access_key"),
				),
			},
		},
	})
}
