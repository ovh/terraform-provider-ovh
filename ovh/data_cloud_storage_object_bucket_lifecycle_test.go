package ovh

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// The data source reuses the resource model, so its schema must keep the same attribute types.
func TestCloudStorageObjectBucketLifecycleDataSourceSchemaMatchesModel(t *testing.T) {
	ctx := context.Background()

	var schemaResp datasource.SchemaResponse
	NewCloudStorageObjectBucketLifecycleDataSource().Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema: %v", schemaResp.Diagnostics)
	}

	model := CloudStorageObjectBucketLifecycleModel{
		ServiceName: ovhtypes.NewTfStringValue("project"),
		BucketId:    ovhtypes.NewTfStringValue("GRA_bucket"),
		Rules:       types.ListNull(types.ObjectType{AttrTypes: storageObjectBucketLifecycleRuleAttrTypes()}),
	}
	response := &CloudStorageObjectBucketLifecycleAPIRules{Rules: []CloudStorageObjectBucketLifecycleAPIRule{
		{
			Id:     "archive",
			Status: "ENABLED",
			Filter: &CloudStorageObjectBucketLifecycleAPIFilter{Prefix: "reports/", Tags: map[string]string{"retention": "long"}},
			Transitions: []CloudStorageObjectBucketLifecycleAPITransition{
				{Days: 90, StorageClass: "STANDARD_IA"},
			},
			NoncurrentVersionExpiration:    &CloudStorageObjectBucketLifecycleAPINoncurrentVersionExpiration{NoncurrentDays: 60},
			AbortIncompleteMultipartUpload: &CloudStorageObjectBucketLifecycleAPIAbortMultipart{DaysAfterInitiation: 7},
		},
	}}
	if diags := model.MergeWith(ctx, response); diags.HasError() {
		t.Fatalf("merge: %v", diags)
	}

	state := tfsdk.State{
		Schema: schemaResp.Schema,
		Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
	}
	if diags := state.Set(ctx, &model); diags.HasError() {
		t.Fatalf("set state: %v", diags)
	}
}

func TestAccDataSourceCloudStorageObjectBucketLifecycle_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	bucketName := acctest.RandomWithPrefix(test_prefix)
	dataSourceName := "data.ovh_cloud_storage_object_bucket_lifecycle.lifecycle"

	bucketConfig := fmt.Sprintf(`
resource "ovh_cloud_storage_object_bucket" "bucket" {
  service_name = %q
  name         = %q
  region       = %q

  versioning = {
    status = "ENABLED"
  }
}
`, serviceName, bucketName, region)

	withoutLifecycle := bucketConfig + `
data "ovh_cloud_storage_object_bucket_lifecycle" "lifecycle" {
  service_name = ovh_cloud_storage_object_bucket.bucket.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.bucket.id
}
`

	withLifecycle := bucketConfig + `
resource "ovh_cloud_storage_object_bucket_lifecycle" "lifecycle" {
  service_name = ovh_cloud_storage_object_bucket.bucket.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.bucket.id

  rules = [
    {
      id     = "expire-logs"
      status = "ENABLED"
      filter = {
        prefix = "logs/"
      }
      expiration = {
        days = 30
      }
    }
  ]
}

data "ovh_cloud_storage_object_bucket_lifecycle" "lifecycle" {
  service_name = ovh_cloud_storage_object_bucket_lifecycle.lifecycle.service_name
  bucket_id    = ovh_cloud_storage_object_bucket_lifecycle.lifecycle.bucket_id
}
`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucket(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketLifecycleDestroy,
		Steps: []resource.TestStep{
			{
				Config: withoutLifecycle,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", region+"_"+bucketName),
					resource.TestCheckResourceAttr(dataSourceName, "rules.#", "0"),
				),
			},
			{
				Config: withLifecycle,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "service_name", serviceName),
					resource.TestCheckResourceAttr(dataSourceName, "bucket_id", region+"_"+bucketName),
					resource.TestCheckResourceAttr(dataSourceName, "id", region+"_"+bucketName),
					resource.TestCheckResourceAttr(dataSourceName, "rules.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "rules.0.id", "expire-logs"),
					resource.TestCheckResourceAttr(dataSourceName, "rules.0.status", "ENABLED"),
					resource.TestCheckResourceAttr(dataSourceName, "rules.0.filter.prefix", "logs/"),
					resource.TestCheckResourceAttr(dataSourceName, "rules.0.expiration.days", "30"),
				),
			},
		},
	})
}
