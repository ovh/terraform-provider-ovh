package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccDataCloudProjectUserS3CredentialsConfig_basic = `
resource "ovh_cloud_project_user" "user" {
 service_name = "%s"
 description  = "my user for acceptance tests"
}

resource "ovh_cloud_project_user_s3_credential" "s3_cred" {
 service_name = ovh_cloud_project_user.user.service_name
 user_id      = ovh_cloud_project_user.user.id
}

data "ovh_cloud_project_user_s3_credentials" "keys" {
 service_name = ovh_cloud_project_user.user.service_name
 user_id      = ovh_cloud_project_user.user.id
 depends_on   = [ovh_cloud_project_user_s3_credential.s3_cred]
}

output "access_key_ids_count" {
    value = length(data.ovh_cloud_project_user_s3_credentials.keys.access_key_ids)
}
`

func TestAccDataCloudProjectUserS3Credentials_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	config := fmt.Sprintf(testAccDataCloudProjectUserS3CredentialsConfig_basic, serviceName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.ovh_cloud_project_user_s3_credentials.keys",
						"access_key_ids.#",
					),
					resource.TestCheckOutput(
						"access_key_ids_count", "1"),
				),
			},
		},
	})
}
