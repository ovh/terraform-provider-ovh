package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

func TestAccCloudProjectStorageLifecycleConfigurationDataSource_basic(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := "GRA"

	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "bucket" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
	}

	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage.bucket.service_name
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "expiration-rule"
				status = "enabled"
				expiration = {
					days = 30
				}
			},
			{
				id     = "abort-multipart-rule"
				status = "enabled"
				abort_incomplete_multipart_upload = {
					days_after_initiation = 7
				}
			}
		]
	}

	data "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle.service_name
		region_name    = ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle.region_name
		container_name = ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle.container_name
	}
	`, serviceName, region, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloud(t); testAccCheckCloudProjectExists(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "service_name", serviceName),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "region_name", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "container_name", bucketName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "id"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "2"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.id", "expiration-rule"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.status", "enabled"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.expiration.days", "30"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.id", "abort-multipart-rule"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.1.abort_incomplete_multipart_upload.days_after_initiation", "7"),
				),
			},
		},
	})
}

// TestAccCloudProjectStorageLifecycleConfigurationDataSource_noServiceName validates that when
// service_name is not set in the data source configuration, the provider falls back to the
// OVH_CLOUD_PROJECT_SERVICE environment variable.
func TestAccCloudProjectStorageLifecycleConfigurationDataSource_noServiceName(t *testing.T) {
	bucketName := acctest.RandomWithPrefix(test_prefix)
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := "GRA"

	config := fmt.Sprintf(`
	resource "ovh_cloud_project_storage" "bucket" {
		service_name = "%s"
		region_name  = "%s"
		name         = "%s"
	}

	resource "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		service_name   = ovh_cloud_project_storage.bucket.service_name
		region_name    = ovh_cloud_project_storage.bucket.region_name
		container_name = ovh_cloud_project_storage.bucket.name

		rules = [
			{
				id     = "expiration-rule"
				status = "enabled"
				expiration = {
					days = 30
				}
			}
		]
	}

	data "ovh_cloud_project_storage_object_bucket_lifecycle_configuration" "lifecycle" {
		region_name    = ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle.region_name
		container_name = ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle.container_name
	}
	`, serviceName, region, bucketName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloud(t)
			testAccCheckCloudProjectExists(t)
			checkEnvOrSkip(t, "OVH_CLOUD_PROJECT_SERVICE")
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "service_name", os.Getenv("OVH_CLOUD_PROJECT_SERVICE")),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "region_name", region),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "container_name", bucketName),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "id"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.#", "1"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.id", "expiration-rule"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.status", "enabled"),
					resource.TestCheckResourceAttr("data.ovh_cloud_project_storage_object_bucket_lifecycle_configuration.lifecycle", "rules.0.expiration.days", "30"),
				),
			},
		},
	})
}

// TestCloudProjectStorageLifecycleConfigurationDataSourceRead_serviceName reads
// the data source without service_name against a fake API: the project comes
// from OVH_CLOUD_PROJECT_SERVICE, and without it the read fails before any
// request instead of calling /cloud/project//region/...
func TestCloudProjectStorageLifecycleConfigurationDataSourceRead_serviceName(t *testing.T) {
	for _, tt := range []struct {
		name    string
		env     string
		request string
	}{
		{"environment variable set", "env-project", "GET /1.0/cloud/project/env-project/region/GRA/storage/my-bucket/lifecycle"},
		{"environment variable unset", "", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("OVH_CLOUD_PROJECT_SERVICE", tt.env)

			var requests []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/1.0/auth/time" {
					fmt.Fprint(w, time.Now().Unix())
					return
				}
				requests = append(requests, r.Method+" "+r.URL.Path)
				fmt.Fprint(w, `{"rules": []}`)
			}))
			defer server.Close()

			client, err := ovh.NewClient(server.URL+"/1.0", "key", "secret", "consumer")
			if err != nil {
				t.Fatal(err)
			}

			ctx := context.Background()
			d := &cloudProjectStorageLifecycleConfigurationDataSource{config: &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}}
			var schemaResp datasource.SchemaResponse
			d.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)

			config := tfsdk.State{Schema: schemaResp.Schema, Raw: tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil)}
			config.SetAttribute(ctx, path.Root("region_name"), "GRA")
			config.SetAttribute(ctx, path.Root("container_name"), "my-bucket")

			resp := datasource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema, Raw: config.Raw.Copy()}}
			d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: config.Raw}}, &resp)

			if tt.request == "" {
				if !resp.Diagnostics.HasError() || resp.Diagnostics.Errors()[0].Summary() != "Missing service_name" {
					t.Errorf("expected a Missing service_name error, got %v", resp.Diagnostics)
				}
				if len(requests) != 0 {
					t.Errorf("expected no request, got %v", requests)
				}
				return
			}

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if len(requests) != 1 || requests[0] != tt.request {
				t.Errorf("got requests %v, want [%s]", requests, tt.request)
			}
			var serviceName string
			resp.State.GetAttribute(ctx, path.Root("service_name"), &serviceName)
			if serviceName != tt.env {
				t.Errorf("got service_name %q in state, want %q", serviceName, tt.env)
			}
		})
	}
}
