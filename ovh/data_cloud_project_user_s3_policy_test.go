package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccDataCloudProjectUserS3PolicyConfig_basic = `
resource "ovh_cloud_project_user" "user" {
  service_name = "%s"
  description  = "my user for acceptance tests"
  role_names   = [
	"objectstore_operator"
  ]
}

resource "ovh_cloud_project_user_s3_credential" "s3_cred" {
  service_name  = ovh_cloud_project_user.user.service_name
  user_id       = ovh_cloud_project_user.user.id
}

resource "ovh_cloud_project_user_s3_policy" "policy" {
  service_name  = ovh_cloud_project_user.user.service_name
  user_id       = ovh_cloud_project_user.user.id
  policy        = <<EOF
%s
EOF
  depends_on    = [ ovh_cloud_project_user_s3_credential.s3_cred ]
}

data "ovh_cloud_project_user_s3_policy" "policy" {
 service_name = ovh_cloud_project_user.user.service_name
 user_id      = ovh_cloud_project_user.user.id
 depends_on   = [ovh_cloud_project_user_s3_policy.policy]
}
`

const normalizedPolicyRWBucket = "{\"Statement\":[{\"Sid\":\"RWContainer\",\"Effect\":\"Allow\",\"Action\":[\"s3:GetObject\",\"s3:PutObject\",\"s3:DeleteObject\",\"s3:ListBucket\",\"s3:ListMultipartUploadParts\",\"s3:ListBucketMultipartUploads\",\"s3:AbortMultipartUpload\",\"s3:GetBucketLocation\"],\"Resource\":[\"arn:aws:s3:::hp-bucket\",\"arn:aws:s3:::hp-bucket/*\"]}]}"

func TestAccDataCloudProjectUserS3Policy_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	config := fmt.Sprintf(testAccDataCloudProjectUserS3PolicyConfig_basic, serviceName, policyRWBucket)
	resourceName := "data.ovh_cloud_project_user_s3_policy.policy"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCredentials(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy"),
					resource.TestCheckResourceAttr(resourceName, "policy", normalizedPolicyRWBucket),
				),
			},
		},
	})
}
