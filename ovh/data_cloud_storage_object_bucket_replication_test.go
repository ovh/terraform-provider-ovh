package ovh

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func assertDataSourceSchemaAcceptsModel(t *testing.T, ds datasource.DataSource, model any) {
	t.Helper()
	ctx := context.Background()

	var schemaResp datasource.SchemaResponse
	ds.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema: %v", schemaResp.Diagnostics)
	}

	state := tfsdk.State{
		Schema: schemaResp.Schema,
		Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
	}
	if diags := state.Set(ctx, model); diags.HasError() {
		t.Fatalf("set state: %v", diags)
	}
}

// The data source reuses the resource model, so its schema must keep the same attribute types.
func TestCloudStorageObjectBucketReplicationDataSourceSchemaMatchesModel(t *testing.T) {
	model := CloudStorageObjectBucketReplicationModel{
		ServiceName: ovhtypes.NewTfStringValue("project"),
		BucketId:    ovhtypes.NewTfStringValue("GRA_source"),
		Rules:       types.ListNull(types.ObjectType{AttrTypes: storageObjectBucketReplicationRuleAttrTypes()}),
	}
	response := &CloudStorageObjectBucketReplicationAPIRules{Rules: []CloudStorageObjectBucketReplicationAPIRule{
		{
			Id:                      "logs",
			Status:                  "ENABLED",
			Priority:                1,
			DeleteMarkerReplication: "DISABLED",
			Destination:             CloudStorageObjectBucketReplicationAPIDestination{Name: "replica", Region: "GRA", StorageClass: "STANDARD_IA"},
			Filter:                  &CloudStorageObjectBucketReplicationAPIFilter{Prefix: "logs/", Tags: map[string]string{"env": "prod"}},
		},
	}}
	if diags := model.MergeWith(context.Background(), response); diags.HasError() {
		t.Fatalf("merge: %v", diags)
	}

	assertDataSourceSchemaAcceptsModel(t, NewCloudStorageObjectBucketReplicationDataSource(), &model)
}

func TestAccDataSourceCloudStorageObjectBucketReplication_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	sourceRegion := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	replicaRegion := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REPLICATION_REGION_TEST"))
	sourceName := acctest.RandomWithPrefix(test_prefix)
	replicaName := acctest.RandomWithPrefix(test_prefix)
	dataSourceName := "data.ovh_cloud_storage_object_bucket_replication.replication"
	buckets := testAccCloudStorageObjectBucketReplicationBuckets(serviceName, sourceName, sourceRegion, replicaName, replicaRegion)

	withoutReplication := buckets + `
data "ovh_cloud_storage_object_bucket_replication" "replication" {
  service_name = ovh_cloud_storage_object_bucket.source.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.source.id
}
`

	withReplication := testAccCloudStorageObjectBucketReplicationConfig(buckets, `
    {
      id                        = "replicate-logs"
      status                    = "ENABLED"
      priority                  = 1
      delete_marker_replication = "DISABLED"
      filter = {
        prefix = "logs/"
      }
      destination = {
        name   = ovh_cloud_storage_object_bucket.replica.name
        region = ovh_cloud_storage_object_bucket.replica.region
      }
    }
`) + `
data "ovh_cloud_storage_object_bucket_replication" "replication" {
  service_name = ovh_cloud_storage_object_bucket_replication.replication.service_name
  bucket_id    = ovh_cloud_storage_object_bucket_replication.replication.bucket_id
}
`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucketReplication(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: withoutReplication,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", sourceRegion+"_"+sourceName),
					resource.TestCheckResourceAttr(dataSourceName, "rules.#", "0"),
				),
			},
			{
				Config: withReplication,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "service_name", serviceName),
					resource.TestCheckResourceAttr(dataSourceName, "bucket_id", sourceRegion+"_"+sourceName),
					resource.TestCheckResourceAttr(dataSourceName, "rules.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "rules.0.id", "replicate-logs"),
					resource.TestCheckResourceAttr(dataSourceName, "rules.0.status", "ENABLED"),
					resource.TestCheckResourceAttr(dataSourceName, "rules.0.priority", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "rules.0.filter.prefix", "logs/"),
					resource.TestCheckResourceAttr(dataSourceName, "rules.0.destination.name", replicaName),
					resource.TestCheckResourceAttr(dataSourceName, "rules.0.destination.region", replicaRegion),
				),
			},
		},
	})
}

func TestCloudStorageObjectBucketReplicationJobDataSourcesSchemaMatchModels(t *testing.T) {
	job := CloudStorageObjectBucketReplicationJobAPIResponse{
		Id:        "3f0c9c1e-6d4b-4c55-9f0e-0d3a1d0b8f11",
		Status:    "COMPLETED",
		Total:     10,
		Completed: 10,
		CreatedAt: "2026-11-11T10:00:00Z",
	}

	single := CloudStorageObjectBucketReplicationJobModel{
		ServiceName: ovhtypes.NewTfStringValue("project"),
		BucketId:    ovhtypes.NewTfStringValue("GRA_source"),
	}
	single.MergeWith(&job)
	assertDataSourceSchemaAcceptsModel(t, NewCloudStorageObjectBucketReplicationJobDataSource(), &single)

	list := cloudStorageObjectBucketReplicationJobsDataSourceModel{
		ServiceName: ovhtypes.NewTfStringValue("project"),
		BucketId:    ovhtypes.NewTfStringValue("GRA_source"),
		Jobs:        []storageObjectBucketReplicationJobItemModel{newStorageObjectBucketReplicationJobItem(&job)},
	}
	assertDataSourceSchemaAcceptsModel(t, NewCloudStorageObjectBucketReplicationJobsDataSource(), &list)
}

func TestAccDataSourceCloudStorageObjectBucketReplicationJobs_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	sourceRegion := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	replicaRegion := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REPLICATION_REGION_TEST"))
	sourceName := acctest.RandomWithPrefix(test_prefix)
	replicaName := acctest.RandomWithPrefix(test_prefix)

	config := testAccCloudStorageObjectBucketReplicationJobConfig(serviceName, sourceName, sourceRegion, replicaName, replicaRegion) + `
data "ovh_cloud_storage_object_bucket_replication_job" "job" {
  service_name = ovh_cloud_storage_object_bucket_replication_job.job.service_name
  bucket_id    = ovh_cloud_storage_object_bucket_replication_job.job.bucket_id
  id           = ovh_cloud_storage_object_bucket_replication_job.job.id
}

data "ovh_cloud_storage_object_bucket_replication_jobs" "jobs" {
  service_name = ovh_cloud_storage_object_bucket_replication_job.job.service_name
  bucket_id    = ovh_cloud_storage_object_bucket_replication_job.job.bucket_id
}
`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucketReplication(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketReplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.ovh_cloud_storage_object_bucket_replication_job.job", "id", testAccCloudStorageObjectBucketReplicationJobName, "id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_replication_job.job", "status"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_replication_job.job", "created_at"),
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_replication_jobs.jobs", "jobs.#", "1"),
					resource.TestCheckResourceAttrPair("data.ovh_cloud_storage_object_bucket_replication_jobs.jobs", "jobs.0.id", testAccCloudStorageObjectBucketReplicationJobName, "id"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_storage_object_bucket_replication_jobs.jobs", "jobs.0.status"),
				),
			},
		},
	})
}

func TestCloudStorageObjectBucketReplicationDestinationsDataSourceSchemaMatchesModel(t *testing.T) {
	model := cloudStorageObjectBucketReplicationDestinationsDataSourceModel{
		ServiceName: ovhtypes.NewTfStringValue("project"),
		BucketId:    ovhtypes.NewTfStringValue("GRA_source"),
		Destinations: []storageObjectBucketReplicationDestinationItemModel{
			newStorageObjectBucketReplicationDestinationItem(&CloudStorageObjectBucketReplicationDestinationAPIResponse{Eligible: true, Name: "replica", Region: "GRA"}),
			newStorageObjectBucketReplicationDestinationItem(&CloudStorageObjectBucketReplicationDestinationAPIResponse{Name: "plain", Region: "GRA", Reason: "VERSIONING_DISABLED"}),
		},
	}
	if !model.Destinations[0].Reason.IsNull() {
		t.Fatalf("reason of an eligible bucket must be null, got %v", model.Destinations[0].Reason)
	}
	assertDataSourceSchemaAcceptsModel(t, NewCloudStorageObjectBucketReplicationDestinationsDataSource(), &model)
}

func TestAccDataSourceCloudStorageObjectBucketReplicationDestinations_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	sourceRegion := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REGION_TEST"))
	replicaRegion := strings.ToUpper(os.Getenv("OVH_CLOUD_PROJECT_STORAGE_REPLICATION_REGION_TEST"))
	sourceName := acctest.RandomWithPrefix(test_prefix)
	replicaName := acctest.RandomWithPrefix(test_prefix)

	config := testAccCloudStorageObjectBucketReplicationBuckets(serviceName, sourceName, sourceRegion, replicaName, replicaRegion) + `
data "ovh_cloud_storage_object_bucket_replication_destinations" "destinations" {
  service_name = ovh_cloud_storage_object_bucket.source.service_name
  bucket_id    = ovh_cloud_storage_object_bucket.source.id

  depends_on = [ovh_cloud_storage_object_bucket.replica]
}
`

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudStorageObjectBucketReplication(t)
			testAccCheckCloudProjectExists(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudStorageObjectBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_storage_object_bucket_replication_destinations.destinations", "bucket_id", sourceRegion+"_"+sourceName),
					testAccCheckCloudStorageObjectBucketReplicationDestination(replicaName, replicaRegion),
				),
			},
		},
	})
}

// The project may hold other buckets, so the replica is looked up by name rather than by index.
func testAccCheckCloudStorageObjectBucketReplicationDestination(name, region string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.ovh_cloud_storage_object_bucket_replication_destinations.destinations"]
		if !ok {
			return fmt.Errorf("destinations data source not found")
		}
		attrs := rs.Primary.Attributes
		count, err := strconv.Atoi(attrs["destinations.#"])
		if err != nil {
			return fmt.Errorf("destinations.#: %w", err)
		}
		for i := 0; i < count; i++ {
			prefix := fmt.Sprintf("destinations.%d.", i)
			if attrs[prefix+"name"] != name {
				continue
			}
			if attrs[prefix+"region"] != region || attrs[prefix+"eligible"] != "true" {
				return fmt.Errorf("destination %s: region = %q, eligible = %q, reason = %q", name, attrs[prefix+"region"], attrs[prefix+"eligible"], attrs[prefix+"reason"])
			}
			return nil
		}
		return fmt.Errorf("destination %s not listed", name)
	}
}
