package ovh

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
	"go.uber.org/ratelimit"
)

func TestCloudProjectServiceNameOrEnv(t *testing.T) {
	t.Run("configured value wins over the environment", func(t *testing.T) {
		t.Setenv("OVH_CLOUD_PROJECT_SERVICE", "from-env")
		var diags diag.Diagnostics
		got := cloudProjectServiceNameOrEnv(ovhtypes.NewTfStringValue("from-config"), &diags)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if got.ValueString() != "from-config" {
			t.Fatalf("got %q, want %q", got.ValueString(), "from-config")
		}
	})

	t.Run("null value falls back to the environment", func(t *testing.T) {
		t.Setenv("OVH_CLOUD_PROJECT_SERVICE", "from-env")
		var diags diag.Diagnostics
		got := cloudProjectServiceNameOrEnv(ovhtypes.NewTfStringNull(), &diags)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if got.ValueString() != "from-env" {
			t.Fatalf("got %q, want %q", got.ValueString(), "from-env")
		}
	})

	t.Run("empty value falls back to the environment", func(t *testing.T) {
		t.Setenv("OVH_CLOUD_PROJECT_SERVICE", "from-env")
		var diags diag.Diagnostics
		got := cloudProjectServiceNameOrEnv(ovhtypes.NewTfStringValue(""), &diags)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if got.ValueString() != "from-env" {
			t.Fatalf("got %q, want %q", got.ValueString(), "from-env")
		}
	})

	t.Run("null value without environment is an error", func(t *testing.T) {
		t.Setenv("OVH_CLOUD_PROJECT_SERVICE", "")
		var diags diag.Diagnostics
		cloudProjectServiceNameOrEnv(ovhtypes.NewTfStringNull(), &diags)
		if !diags.HasError() {
			t.Fatal("expected an error diagnostic")
		}
	})

	t.Run("plain string value follows the same rules", func(t *testing.T) {
		t.Setenv("OVH_CLOUD_PROJECT_SERVICE", "from-env")
		var diags diag.Diagnostics
		configured := cloudProjectServiceNameStringOrEnv(types.StringValue("from-config"), &diags)
		omitted := cloudProjectServiceNameStringOrEnv(types.StringNull(), &diags)
		if diags.HasError() {
			t.Fatalf("unexpected diagnostics: %v", diags)
		}
		if configured.ValueString() != "from-config" || omitted.ValueString() != "from-env" {
			t.Fatalf("got %q and %q, want %q and %q", configured.ValueString(), omitted.ValueString(), "from-config", "from-env")
		}
	})
}

// cloudDataSourceWithServiceName is a framework ovh_cloud_* data source that
// takes a service_name.
type cloudDataSourceWithServiceName struct {
	typeName   string
	dataSource datasource.DataSource
	schema     datasource.SchemaResponse
}

// cloudDataSourcesWithServiceName asks the provider itself for the framework
// ovh_cloud_* data sources it registers, and keeps those taking a service_name.
func cloudDataSourcesWithServiceName(ctx context.Context) []cloudDataSourceWithServiceName {
	var found []cloudDataSourceWithServiceName

	for _, newDataSource := range (&OvhProvider{}).DataSources(ctx) {
		dataSource := newDataSource()

		var metadata datasource.MetadataResponse
		dataSource.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "ovh"}, &metadata)
		if !strings.HasPrefix(metadata.TypeName, "ovh_cloud_") {
			continue
		}

		var schemaResp datasource.SchemaResponse
		dataSource.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
		if _, ok := schemaResp.Schema.Attributes["service_name"]; !ok {
			continue
		}

		found = append(found, cloudDataSourceWithServiceName{metadata.TypeName, dataSource, schemaResp})
	}

	return found
}

// TestCloudDataSourcesServiceNameIsOptional checks every ovh_cloud_* data source
// the provider registers, framework and SDKv2 alike: one that takes a
// service_name must not require it, so that it can fall back to the
// OVH_CLOUD_PROJECT_SERVICE environment variable.
func TestCloudDataSourcesServiceNameIsOptional(t *testing.T) {
	t.Setenv("OVH_CLOUD_PROJECT_SERVICE", "from-env")

	checked := 0

	for _, d := range cloudDataSourcesWithServiceName(context.Background()) {
		checked++
		if d.schema.Schema.Attributes["service_name"].IsRequired() {
			t.Errorf("%s: service_name is Required", d.typeName)
		}
	}

	for typeName, dataSource := range Provider().DataSourcesMap {
		if !strings.HasPrefix(typeName, "ovh_cloud_") {
			continue
		}

		serviceName, ok := dataSource.Schema["service_name"]
		if !ok {
			continue
		}

		// With the SDK, a Required attribute is satisfied by its DefaultFunc.
		checked++
		if serviceName.DefaultFunc == nil {
			t.Errorf("%s: service_name has no DefaultFunc", typeName)
			continue
		}
		if value, err := serviceName.DefaultFunc(); err != nil || value != "from-env" {
			t.Errorf("%s: service_name does not default to OVH_CLOUD_PROJECT_SERVICE (got %v, err %v)", typeName, value, err)
		}
	}

	if checked == 0 {
		t.Fatal("no ovh_cloud_* data source with a service_name was found")
	}
	t.Logf("checked %d ovh_cloud_* data sources", checked)
}

// TestCloudDataSourcesReadServiceNameFromEnv reads every framework ovh_cloud_*
// data source that takes a service_name, without service_name in its
// configuration, against a fake API answering 404 to everything: each of them
// must call the API, and only for the project of OVH_CLOUD_PROJECT_SERVICE.
//
// It is the runtime counterpart of TestCloudDataSourcesServiceNameIsOptional: a
// service_name that is optional in the schema but that Read does not resolve
// would request /v2/publicCloud/project//... .
//
// A data source that cannot be read this way (it validates its arguments
// before any request, or does not go through Config.OVHClient) fails here
// with "no API request": give it readable dummy arguments below rather than
// dropping it from the check.
func TestCloudDataSourcesReadServiceNameFromEnv(t *testing.T) {
	const project = "project-from-env"
	t.Setenv("OVH_CLOUD_PROJECT_SERVICE", project)
	// The client below also reads these, and refuses two ways to authenticate.
	for _, name := range []string{"OVH_CLIENT_ID", "OVH_CLIENT_SECRET", "OVH_ACCESS_TOKEN"} {
		t.Setenv(name, "")
	}

	var mu sync.Mutex
	var requests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/1.0/auth/time" {
			fmt.Fprint(w, time.Now().Unix())
			return
		}
		mu.Lock()
		requests = append(requests, r.Method+" "+r.URL.Path)
		mu.Unlock()
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{"class": "Client::NotFound", "message": "not found"}`)
	}))
	t.Cleanup(server.Close)

	client, err := ovh.NewClient(server.URL+"/1.0", "key", "secret", "consumer")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	config := &Config{OVHClient: ovhwrap.NewClient(client, ratelimit.NewUnlimited())}

	ctx := context.Background()
	dataSources := cloudDataSourcesWithServiceName(ctx)

	for _, d := range dataSources {
		if configurable, ok := d.dataSource.(datasource.DataSourceWithConfigure); ok {
			configurable.Configure(ctx, datasource.ConfigureRequest{ProviderData: config}, &datasource.ConfigureResponse{})
		}

		// The configuration sets the required arguments, all of them strings,
		// and leaves everything else null: service_name is omitted.
		objectType := d.schema.Schema.Type().TerraformType(ctx).(tftypes.Object)
		values := make(map[string]tftypes.Value, len(objectType.AttributeTypes))
		for name, attributeType := range objectType.AttributeTypes {
			values[name] = tftypes.NewValue(attributeType, nil)
			if d.schema.Schema.Attributes[name].IsRequired() && attributeType.Is(tftypes.String) {
				values[name] = tftypes.NewValue(attributeType, "x")
			}
		}
		raw := tftypes.NewValue(objectType, values)

		mu.Lock()
		requests = nil
		mu.Unlock()

		resp := datasource.ReadResponse{State: tfsdk.State{Schema: d.schema.Schema, Raw: raw}}
		d.dataSource.Read(ctx, datasource.ReadRequest{Config: tfsdk.Config{Schema: d.schema.Schema, Raw: raw}}, &resp)

		mu.Lock()
		if len(requests) == 0 {
			t.Errorf("%s: no API request (diagnostics: %v)", d.typeName, resp.Diagnostics)
		}
		for _, request := range requests {
			if !strings.Contains(request+"/", "/"+project+"/") {
				t.Errorf("%s: %q does not target the project of OVH_CLOUD_PROJECT_SERVICE", d.typeName, request)
			}
		}
		mu.Unlock()
	}

	if len(dataSources) == 0 {
		t.Fatal("no ovh_cloud_* data source with a service_name was found")
	}
	t.Logf("read %d ovh_cloud_* data sources", len(dataSources))
}
