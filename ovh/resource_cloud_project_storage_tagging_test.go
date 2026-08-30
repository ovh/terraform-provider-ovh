package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// Unit tests

func TestCloudProjectStorageTaggingPayload_JSON_WithTags(t *testing.T) {
	ctx := context.Background()
	tags := tfStringMap(ctx, map[string]string{"env": "prod", "team": "platform"})
	payload := cloudProjectStorageTaggingPayload{Tags: &tags}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	raw, found := parsed["tags"]
	if !found {
		t.Fatalf("expected 'tags' in JSON, got: %s", data)
	}

	var tagMap map[string]string
	if err := json.Unmarshal(raw, &tagMap); err != nil {
		t.Fatalf("unexpected tags unmarshal error: %v", err)
	}
	if tagMap["env"] != "prod" {
		t.Errorf("expected tags[env]=prod, got: %v", tagMap)
	}
}

func TestCloudProjectStorageTaggingPayload_JSON_WithEmptyTags(t *testing.T) {
	ctx := context.Background()
	emptyTags, _ := ovhtypes.NewTfMapNestedValue[ovhtypes.TfStringValue](ctx, map[string]attr.Value{})
	payload := cloudProjectStorageTaggingPayload{Tags: &emptyTags}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}

	var parsed map[string]json.RawMessage
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	raw, found := parsed["tags"]
	if !found {
		t.Fatalf("expected 'tags' in JSON, got: %s", data)
	}
	if string(raw) != "{}" {
		t.Errorf("expected tags to be {}, got: %s", raw)
	}
}

func TestCompositeID(t *testing.T) {
	m := cloudProjectStorageTaggingModel{
		ServiceName: ovhtypes.NewTfStringValue("my-project"),
		RegionName:  ovhtypes.NewTfStringValue("GRA"),
		Name:        ovhtypes.NewTfStringValue("my-bucket"),
	}
	got := compositeID(m)
	want := "my-project/GRA/my-bucket"
	if got != want {
		t.Errorf("compositeID = %q, want %q", got, want)
	}
}

// Acceptance tests

func TestAccCloudProjectStorageTagging_basic(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create bucket then apply tags via the standalone resource
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "bucket" {
						service_name = "%s"
						region_name  = "GRA"
						name         = "%s"
					}

					resource "ovh_cloud_project_storage_tagging" "tags" {
						service_name = ovh_cloud_project_storage.bucket.service_name
						region_name  = ovh_cloud_project_storage.bucket.region_name
						name         = ovh_cloud_project_storage.bucket.name
						tags = {
							environment = "test"
							team        = "platform"
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_tagging.tags", "name", bucketName),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_tagging.tags", "tags.environment", "test"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_tagging.tags", "tags.team", "platform"),
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_tagging.tags", "id",
						fmt.Sprintf("%s/GRA/%s", serviceName, bucketName)),
				),
			},
			// Step 2: Update tags
			{
				Config: fmt.Sprintf(`
					resource "ovh_cloud_project_storage" "bucket" {
						service_name = "%s"
						region_name  = "GRA"
						name         = "%s"
					}

					resource "ovh_cloud_project_storage_tagging" "tags" {
						service_name = ovh_cloud_project_storage.bucket.service_name
						region_name  = ovh_cloud_project_storage.bucket.region_name
						name         = ovh_cloud_project_storage.bucket.name
						tags = {
							environment = "production"
						}
					}`, serviceName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_project_storage_tagging.tags", "tags.environment", "production"),
					resource.TestCheckNoResourceAttr("ovh_cloud_project_storage_tagging.tags", "tags.team"),
				),
			},
		},
	})
}
