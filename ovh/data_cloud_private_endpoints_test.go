package ovh

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// TestPrivateEndpointListItemAttrTypesMatchSchema guards the drift that the
// plural data source cannot detect on its own: privateEndpointListItemAttrTypes
// and the nested object declared in the schema must agree, otherwise every
// element fails its type check at apply time with an opaque error.
func TestPrivateEndpointListItemAttrTypesMatchSchema(t *testing.T) {
	ctx := context.Background()

	var resp datasource.SchemaResponse
	(&cloudPrivateEndpointsDataSource{}).Schema(ctx, datasource.SchemaRequest{}, &resp)

	listAttr, ok := resp.Schema.Attributes["private_endpoints"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatalf("private_endpoints is not a ListNestedAttribute, got %T", resp.Schema.Attributes["private_endpoints"])
	}

	fromSchema, ok := listAttr.NestedObject.Type().(types.ObjectType)
	if !ok {
		t.Fatalf("nested object type is not an ObjectType, got %T", listAttr.NestedObject.Type())
	}

	declared := privateEndpointListItemAttrTypes()
	schemaTypes := fromSchema.AttributeTypes()

	for name, want := range declared {
		got, present := schemaTypes[name]
		if !present {
			t.Errorf("%q declared in privateEndpointListItemAttrTypes but absent from the schema", name)
			continue
		}
		if !got.Equal(want) {
			t.Errorf("%q type mismatch: schema has %s, privateEndpointListItemAttrTypes has %s", name, got, want)
		}
	}
	for name := range schemaTypes {
		if _, present := declared[name]; !present {
			t.Errorf("%q present in the schema but missing from privateEndpointListItemAttrTypes", name)
		}
	}
}

// TestBuildPrivateEndpointListItemObject verifies one API response maps onto a
// well-typed list element, flattening the placement the platform resolved.
func TestBuildPrivateEndpointListItemObject(t *testing.T) {
	response := &CloudPrivateEndpointAPIResponse{
		Id:             "pe-1",
		Checksum:       "checksum-abc",
		ResourceStatus: "READY",
		TargetSpec: &CloudPrivateEndpointTargetSpec{
			Name:             "my-endpoint",
			Location:         &CloudPrivateEndpointLocation{Region: "GRA11"},
			Service:          &CloudPrivateEndpointServiceRef{Name: "s3", Tier: "S"},
			Subnet:           &CloudPrivateEndpointRef{Id: "subnet-456"},
			PublicDnsEnabled: true,
		},
		CurrentState: &CloudPrivateEndpointCurrentState{
			EndpointId:       "pe-0123abcd",
			Service:          &CloudPrivateEndpointServiceState{Name: "s3", Tier: "S", Port: 443, Transport: "TCP"},
			Network:          &CloudPrivateEndpointRef{Id: "network-123"},
			Subnet:           &CloudPrivateEndpointRef{Id: "subnet-456"},
			VipAddress:       "10.0.0.7",
			Hostname:         "endpoint.example.net",
			PublicDnsEnabled: true,
		},
	}

	obj, diags := buildPrivateEndpointListItemObject(context.Background(), response)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	attrs := obj.Attributes()
	for name, want := range map[string]string{
		"id":              "pe-1",
		"name":            "my-endpoint",
		"region":          "GRA11",
		"network_id":      "network-123",
		"subnet_id":       "subnet-456",
		"vip_address":     "10.0.0.7",
		"checksum":        "checksum-abc",
		"resource_status": "READY",
	} {
		if got := attrs[name].(ovhtypes.TfStringValue).ValueString(); got != want {
			t.Errorf("unexpected %s: got %q, want %q", name, got, want)
		}
	}

	if !attrs["public_dns_enabled"].(types.Bool).ValueBool() {
		t.Error("expected public_dns_enabled to be true")
	}

	service := attrs["service"].(types.Object).Attributes()
	if got := service["tier"].(ovhtypes.TfStringValue).ValueString(); got != "S" {
		t.Errorf("unexpected service.tier: got %q", got)
	}

	currentState := attrs["current_state"].(types.Object).Attributes()
	if got := currentState["hostname"].(ovhtypes.TfStringValue).ValueString(); got != "endpoint.example.net" {
		t.Errorf("unexpected current_state.hostname: got %q", got)
	}
	if got := currentState["endpoint_id"].(ovhtypes.TfStringValue).ValueString(); got != "pe-0123abcd" {
		t.Errorf("unexpected current_state.endpoint_id: got %q", got)
	}
	if !currentState["public_dns_enabled"].(types.Bool).ValueBool() {
		t.Error("expected current_state.public_dns_enabled to be true")
	}
	csService := currentState["service"].(types.Object).Attributes()
	if got := csService["tier"].(ovhtypes.TfStringValue).ValueString(); got != "S" {
		t.Errorf("unexpected current_state.service.tier: got %q", got)
	}
	if !csService["bandwidth"].(types.Int64).IsNull() {
		t.Error("expected current_state.service.bandwidth to be null when the API omits it")
	}
}

// TestBuildPrivateEndpointListItemObject_EmptyEnvelope verifies an element with
// neither targetSpec nor currentState still types correctly. MergeWith leaves
// those attributes untouched, so the builder has to seed them with typed nulls.
func TestBuildPrivateEndpointListItemObject_EmptyEnvelope(t *testing.T) {
	obj, diags := buildPrivateEndpointListItemObject(context.Background(), &CloudPrivateEndpointAPIResponse{Id: "pe-1"})
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	attrs := obj.Attributes()
	if !attrs["service"].(types.Object).IsNull() {
		t.Error("expected service to be null")
	}
	if !attrs["current_state"].(types.Object).IsNull() {
		t.Error("expected current_state to be null")
	}
	if got := attrs["id"].(ovhtypes.TfStringValue).ValueString(); got != "pe-1" {
		t.Errorf("unexpected id: got %q", got)
	}
}

func TestAccDataSourceCloudPrivateEndpoints_basic(t *testing.T) {
	env := privateEndpointTestEnv()
	name := acctest.RandomWithPrefix(testAccResourceCloudPrivateEndpointNamePrefix)

	config := testAccCloudPrivateEndpointConfig(env, name, "listed", "", env.subnetID, "") + `
data "ovh_cloud_private_endpoints" "endpoints" {
  service_name = ovh_cloud_private_endpoint.endpoint.service_name
}

data "ovh_cloud_private_endpoint" "endpoint" {
  service_name = ovh_cloud_private_endpoint.endpoint.service_name
  id           = ovh_cloud_private_endpoint.endpoint.id
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheckCloudPrivateEndpoint(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCloudPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					// The singular data source must agree with the resource it reads.
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_private_endpoint.endpoint", "name",
						testAccPrivateEndpointResource, "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_private_endpoint.endpoint", "vip_address",
						testAccPrivateEndpointResource, "vip_address",
					),
					resource.TestCheckResourceAttrPair(
						"data.ovh_cloud_private_endpoint.endpoint", "subnet_id",
						testAccPrivateEndpointResource, "subnet_id",
					),
					resource.TestCheckResourceAttr("data.ovh_cloud_private_endpoint.endpoint", "description", "listed"),
					resource.TestCheckResourceAttr("data.ovh_cloud_private_endpoint.endpoint", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("data.ovh_cloud_private_endpoint.endpoint", "current_state.hostname"),
					// The endpoint just created must appear in the listing.
					testAccCheckPrivateEndpointListed("data.ovh_cloud_private_endpoints.endpoints", name),
				),
			},
		},
	})
}

// testAccCheckPrivateEndpointListed asserts the plural data source returned an
// element whose name matches, without depending on its index in the list.
func testAccCheckPrivateEndpointListed(dataSourceName, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("not found: %s", dataSourceName)
		}

		count := rs.Primary.Attributes["private_endpoints.#"]
		if count == "" || count == "0" {
			return fmt.Errorf("%s returned no private endpoint", dataSourceName)
		}

		for i := 0; ; i++ {
			got, present := rs.Primary.Attributes[fmt.Sprintf("private_endpoints.%d.name", i)]
			if !present {
				return fmt.Errorf("private endpoint %q not found in %s (%s listed)", name, dataSourceName, count)
			}
			if got == name {
				return nil
			}
		}
	}
}

// TestPrivateEndpointDataSourceSchemaMatchesModel guards the singular data
// source, which reads straight into the resource model: a field missing from
// the schema, or a schema attribute missing from the model, only surfaces as a
// Config.Get error at read time.
func TestPrivateEndpointDataSourceSchemaMatchesModel(t *testing.T) {
	ctx := context.Background()

	var resp datasource.SchemaResponse
	(&cloudPrivateEndpointDataSource{}).Schema(ctx, datasource.SchemaRequest{}, &resp)

	modelTypes := modelAttrTypes(ctx, CloudPrivateEndpointModel{})
	rt := reflect.TypeOf(CloudPrivateEndpointModel{})
	fields := map[string]bool{}
	for i := range rt.NumField() {
		name := rt.Field(i).Tag.Get("tfsdk")
		fields[name] = true
		a, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("%q declared in CloudPrivateEndpointModel but absent from the schema", name)
			continue
		}
		if want, known := modelTypes[name]; known && !a.GetType().Equal(want) {
			t.Errorf("%q type mismatch: schema has %s, model expects %s", name, a.GetType(), want)
		}
	}
	for name := range resp.Schema.Attributes {
		if !fields[name] {
			t.Errorf("%q present in the schema but missing from CloudPrivateEndpointModel", name)
		}
	}
}
