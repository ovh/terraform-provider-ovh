package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAccCloudProjectUserS3CredentialConfig = fmt.Sprintf(`
resource "ovh_cloud_project_user" "user" {
 service_name = "%s"
 description  = "my user for acceptance tests"
}

resource "ovh_cloud_project_user_s3_credential" "s3_cred" {
 service_name  = ovh_cloud_project_user.user.service_name
 user_id       = ovh_cloud_project_user.user.id
}
`, os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST"))

func TestAccCloudProjectUserS3Credential_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudProjectUserS3CredentialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_user_s3_credential.s3_cred", "access_key_id"),
					resource.TestCheckResourceAttrSet(
						"ovh_cloud_project_user_s3_credential.s3_cred", "secret_access_key"),
				),
			},
		},
	})
}
