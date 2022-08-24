package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccDataCloudProjectUserS3CredentialsConfig = fmt.Sprintf(`
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
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

func TestAccDataCloudProjectUserS3Credentials_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataCloudProjectUserS3CredentialsConfig,
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
