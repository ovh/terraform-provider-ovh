package ovh

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/ovh/terraform-provider-ovh/ovh/helpers"
)

const testAccCloudProjectUserS3PolicyConfig_basic = `
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

resource "ovh_cloud_project_user_s3_policy" "s3_policy" {
 service_name  = ovh_cloud_project_user.user.service_name
 user_id       = ovh_cloud_project_user.user.id
 policy        = <<EOF
%s
EOF
 depends_on    = [ ovh_cloud_project_user_s3_credential.s3_cred ]
}
`

const (
	policyRWBucket = `{
		"Statement":[{
			"Sid": "RWContainer",
			"Effect": "Allow",
			"Action":["s3:GetObject", "s3:PutObject", "s3:DeleteObject", "s3:ListBucket", "s3:ListMultipartUploadParts", "s3:ListBucketMultipartUploads", "s3:AbortMultipartUpload", "s3:GetBucketLocation"],
			"Resource":["arn:aws:s3:::hp-bucket", "arn:aws:s3:::hp-bucket/*"]
		}]
	}`

	policyRWBucketNoUpdate = `{
		"Statement":[{
		  "Sid":   "RWContainer",
		  "Effect":   "Allow",
		  "Action":  ["s3:GetObject", "s3:PutObject", "s3:DeleteObject", "s3:ListBucket", "s3:ListMultipartUploadParts", "s3:ListBucketMultipartUploads", "s3:AbortMultipartUpload", "s3:GetBucketLocation"],
		  "Resource":  ["arn:aws:s3:::hp-bucket", "arn:aws:s3:::hp-bucket/*"]
		}]
	}`

	policyRWBucketUpdated = `{
		"Statement":[{
		  "Sid": "RWContainer",
		  "Effect": "Deny",
		  "Action": ["s3:GetObject"],
		  "Resource": ["arn:aws:s3:::hp-bucket/*"]
		}]
	}`
)

func TestAccCloudProjectUserS3Policy_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	resourceName := "ovh_cloud_project_user_s3_policy.s3_policy"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCloudProjectUserS3PolicyConfig_basic, serviceName, policyRWBucket),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy"),
					resource.TestCheckResourceAttrWith(resourceName, "policy", func(value string) error {
						if !helpers.JSONStringsEqual(value, policyRWBucket) {
							return fmt.Errorf("policies are not equal, expected %q, got %q", value, policyRWBucket)
						}
						return nil
					}),
				),
			},
			{
				Config: fmt.Sprintf(testAccCloudProjectUserS3PolicyConfig_basic, serviceName, policyRWBucketNoUpdate),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						ExpectEmptyPlan(),
					},
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(testAccCloudProjectUserS3PolicyConfig_basic, serviceName, policyRWBucketUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "policy"),
					resource.TestCheckResourceAttrWith(resourceName, "policy", func(value string) error {
						if !helpers.JSONStringsEqual(value, policyRWBucketUpdated) {
							return fmt.Errorf("policies are not equal, expected %q, got %q", value, policyRWBucketUpdated)
						}
						return nil
					}),
				),
			},
		},
	})
}
