package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataCloudProjectUserS3CredentialConfig_basic = `
resource "ovh_cloud_project_user" "user" {
 service_name = "%s"
 description  = "my user for acceptance tests"
}

resource "ovh_cloud_project_user_s3_credential" "s3_cred_1" {
 service_name = ovh_cloud_project_user.user.service_name
 user_id      = ovh_cloud_project_user.user.id
}

resource "ovh_cloud_project_user_s3_credential" "s3_cred_2" {
 service_name = ovh_cloud_project_user.user.service_name
 user_id      = ovh_cloud_project_user.user.id
}

data "ovh_cloud_project_user_s3_credential" "s3_cred_key_2" {
 service_name  = ovh_cloud_project_user.user.service_name
 user_id       = ovh_cloud_project_user.user.id
 access_key_id = ovh_cloud_project_user_s3_credential.s3_cred_2.access_key_id
 depends_on    = [ovh_cloud_project_user_s3_credential.s3_cred_1, ovh_cloud_project_user_s3_credential.s3_cred_2]
}

data "ovh_cloud_project_user_s3_credential" "s3_cred_key_1" {
 service_name  = ovh_cloud_project_user.user.service_name
 user_id       = ovh_cloud_project_user.user.id
 access_key_id = ovh_cloud_project_user_s3_credential.s3_cred_1.access_key_id
 depends_on    = [ovh_cloud_project_user_s3_credential.s3_cred_1, ovh_cloud_project_user_s3_credential.s3_cred_2]
}

output "same_secret_key_cred_1" {
    value = data.ovh_cloud_project_user_s3_credential.s3_cred_key_1.secret_access_key == ovh_cloud_project_user_s3_credential.s3_cred_1.secret_access_key
    sensitive=true
}

output "same_secret_key_cred_2" {
    value = data.ovh_cloud_project_user_s3_credential.s3_cred_key_2.secret_access_key == ovh_cloud_project_user_s3_credential.s3_cred_2.secret_access_key
    sensitive=true
}
`

func TestAccDataCloudProjectUserS3Credential_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	config := fmt.Sprintf(testAccDataCloudProjectUserS3CredentialConfig_basic, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_user_s3_credential.s3_cred_key_1",
						"secret_access_key",
					),
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_user_s3_credential.s3_cred_key_2",
						"secret_access_key",
					),
					resource.TestCheckOutput(
						"same_secret_key_cred_1", "true"),
					resource.TestCheckOutput(
						"same_secret_key_cred_2", "true"),
				),
			},
		},
	})
}
