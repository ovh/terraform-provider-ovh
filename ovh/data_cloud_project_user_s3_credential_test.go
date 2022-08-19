package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccDataCloudProjectUserS3CredentialConfig = fmt.Sprintf(`
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

data "ovh_cloud_project_user_s3_credentials" "keys" {
 service_name = ovh_cloud_project_user.user.service_name
 user_id      = ovh_cloud_project_user.user.id
 depends_on   = [ovh_cloud_project_user_s3_credential.s3_cred_1, ovh_cloud_project_user_s3_credential.s3_cred_2]
}

data "ovh_cloud_project_user_s3_credential" "s3_cred_key_2" {
 service_name  = ovh_cloud_project_user.user.service_name
 user_id       = ovh_cloud_project_user.user.id
 access_key_id = data.ovh_cloud_project_user_s3_credentials.keys.access_key_ids[1]
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
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

func TestAccDataCloudProjectUserS3Credential_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataCloudProjectUserS3CredentialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput(
						"same_secret_key_cred_1", "true"),
					resource.TestCheckOutput(
						"same_secret_key_cred_2", "true"),
				),
			},
		},
	})
}
