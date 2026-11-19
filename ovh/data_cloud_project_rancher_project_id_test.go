package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	"go.uber.org/ratelimit"
)

// TestCloudProjectRancherDataSourcesProjectIdFromEnv reads the
// ovh_cloud_project_rancher* data sources without project_id against a fake API
// answering 404: with OVH_CLOUD_PROJECT_SERVICE set, every request targets that
// project; without it, the read fails with "Missing project_id" before any
// request.
func TestCloudProjectRancherDataSourcesProjectIdFromEnv(t *testing.T) {
	dataSources := map[string]func() datasource.DataSource{
		"ovh_cloud_project_rancher":                      NewCloudProjectRancherDataSource,
		"ovh_cloud_project_rancher_plan":                 NewCloudProjectRancherPlanDataSource,
		"ovh_cloud_project_rancher_version":              NewCloudProjectRancherVersionDataSource,
		"ovh_cloud_project_rancher_capabilities_plan":    NewCloudProjectRancherCapabilitiesPlanDataSource,
		"ovh_cloud_project_rancher_capabilities_version": NewCloudProjectRancherCapabilitiesVersionDataSource,
	}

	var requests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/1.0/auth/time" {
			fmt.Fprint(w, time.Now().Unix())
			return
		}
		requests = append(requests, r.Method+" "+r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"class": "Client::NotFound", "message": "not found"}`)
	}))
	t.Cleanup(server.Close)

	client, err := ovh.NewClient(server.URL+"/1.0", "key", "secret", "consumer")
	if err != nil {
		t.Fatal(err)
	}
	config := &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}

	for _, env := range []string{"env-project", ""} {
		for typeName, newDataSource := range dataSources {
			t.Run(fmt.Sprintf("%s/env=%q", typeName, env), func(t *testing.T) {
				t.Setenv("OVH_CLOUD_PROJECT_SERVICE", env)
				// The client also reads these, and refuses two ways to authenticate.
				for _, name := range []string{"OVH_CLIENT_ID", "OVH_CLIENT_SECRET", "OVH_ACCESS_TOKEN"} {
					t.Setenv(name, "")
				}

				ctx := context.Background()
				d := newDataSource()
				d.(datasource.DataSourceWithConfigure).Configure(ctx, datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})

				var metadata datasource.MetadataResponse
				d.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "ovh"}, &metadata)
				if metadata.TypeName != typeName {
					t.Fatalf("got type name %s", metadata.TypeName)
				}

				var schemaResp datasource.SchemaResponse
				d.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
				if projectId := schemaResp.Schema.Attributes["project_id"]; projectId.IsRequired() || !projectId.IsComputed() {
					t.Errorf("project_id must be Optional and Computed")
				}

				// Every required argument is a string set to "x"; project_id is omitted.
				objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
				values := map[string]tftypes.Value{}
				for name, attributeType := range objectType.AttributeTypes {
					values[name] = tftypes.NewValue(attributeType, nil)
					if schemaResp.Schema.Attributes[name].IsRequired() {
						values[name] = tftypes.NewValue(attributeType, "x")
					}
				}
				raw := tftypes.NewValue(objectType, values)

				requests = nil
				resp := datasource.ReadResponse{State: tfsdk.State{Schema: schemaResp.Schema, Raw: raw}}
				d.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: raw}}, &resp)

				if env == "" {
					if !resp.Diagnostics.HasError() || resp.Diagnostics.Errors()[0].Summary() != "Missing project_id" {
						t.Errorf("expected a Missing project_id error, got %v", resp.Diagnostics)
					}
					if len(requests) != 0 {
						t.Errorf("expected no request, got %v", requests)
					}
					return
				}

				if len(requests) == 0 {
					t.Errorf("no API request (diagnostics: %v)", resp.Diagnostics)
				}
				for _, request := range requests {
					if !strings.Contains(request, "/v2/publicCloud/project/"+env+"/") {
						t.Errorf("%q does not target the project of OVH_CLOUD_PROJECT_SERVICE", request)
					}
				}
			})
		}
	}
}
