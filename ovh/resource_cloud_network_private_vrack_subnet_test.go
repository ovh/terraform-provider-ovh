package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

const testAccResourceCloudNetworkPrivateSubnetNamePrefix = "tf-test-subnet-"

func testAccPreCheckCloudNetworkPrivateSubnet(t *testing.T) {
	testAccPreCheckCredentials(t)
	if os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_SERVICE_TEST must be set for acceptance tests")
	}
	if os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST") == "" {
		t.Skip("OVH_CLOUD_PROJECT_REGION_TEST must be set for acceptance tests")
	}
	// Creating a regional private network requires a vRack.
	testAccPreCheckVRack(t)
}

func testAccCloudNetworkPrivateSubnetImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s/%s",
			rs.Primary.Attributes["service_name"],
			rs.Primary.Attributes["network_id"],
			rs.Primary.Attributes["id"],
		), nil
	}
}

func TestAccCloudNetworkPrivateVrackSubnet_basic(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id = ovh_cloud_network_private_vrack.network.id

  name = "%s"
  cidr = "10.0.0.0/24"
  region       = "%s"
}
`, serviceName, region, subnetName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateSubnet(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "service_name", serviceName),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "network_id"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "region", region),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "checksum"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "created_at"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "resource_status", "READY"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "current_state.name"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "current_state.cidr"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "current_state.location.region"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_network_private_vrack_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudNetworkPrivateSubnetImportStateIdFunc("ovh_cloud_network_private_vrack_subnet.test"),
			},
		},
	})
}

// TestAccCloudNetworkPrivateVrackSubnet_withAllOptions creates a subnet exercising all
// optional attributes: dhcp_enabled, gateway_ip, dns_nameservers and
// allocation_pools.
func TestAccCloudNetworkPrivateVrackSubnet_withAllOptions(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id = ovh_cloud_network_private_vrack.network.id

  name            = "%s"
  cidr            = "10.0.0.0/24"
  dhcp_enabled    = true
  gateway_ip      = "10.0.0.1"
  dns_nameservers = ["1.1.1.1", "8.8.8.8"]
  region       = "%s"
  allocation_pools = [
    {
      start = "10.0.0.10"
      end   = "10.0.0.100"
    },
  ]
}
`, serviceName, region, subnetName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateSubnet(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "cidr", "10.0.0.0/24"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dhcp_enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "gateway_ip", "10.0.0.1"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.#", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.1", "8.8.8.8"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "allocation_pools.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "allocation_pools.0.start", "10.0.0.10"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "allocation_pools.0.end", "10.0.0.100"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "resource_status", "READY"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "current_state.gateway_ip", "10.0.0.1"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "current_state.dhcp_enabled", "true"),
				),
			},
			// Test import
			{
				ResourceName:      "ovh_cloud_network_private_vrack_subnet.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCloudNetworkPrivateSubnetImportStateIdFunc("ovh_cloud_network_private_vrack_subnet.test"),
			},
		},
	})
}

// TestAccCloudNetworkPrivateVrackSubnet_updateMutableFields toggles dhcp_enabled and
// updates dns_nameservers / description / name in place (all mutable, no replacement).
func TestAccCloudNetworkPrivateVrackSubnet_updateMutableFields(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id = ovh_cloud_network_private_vrack.network.id

  name            = "%s"
  cidr            = "10.0.0.0/24"
  description     = "initial description"
  dhcp_enabled    = false
  dns_nameservers = ["1.1.1.1"]
  region       = "%s"
}
`, serviceName, region, subnetName, region)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id = ovh_cloud_network_private_vrack.network.id

  name            = "%s"
  cidr            = "10.0.0.0/24"
  description     = "updated description"
  dhcp_enabled    = true
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
  region       = "%s"
}
`, serviceName, region, updatedName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateSubnet(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "description", "initial description"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dhcp_enabled", "false"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.#", "1"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.0", "1.1.1.1"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "id"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "description", "updated description"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dhcp_enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.#", "2"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.0", "8.8.8.8"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "dns_nameservers.1", "8.8.4.4"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "current_state.dhcp_enabled", "true"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "resource_status", "READY"),
				),
			},
		},
	})
}

func TestAccCloudNetworkPrivateVrackSubnet_update(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)
	updatedName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)

	config := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id = ovh_cloud_network_private_vrack.network.id

  name        = "%s"
  cidr        = "10.0.0.0/24"
  description = "initial description"
  region       = "%s"
}
`, serviceName, region, subnetName, region)

	updatedConfig := fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id = ovh_cloud_network_private_vrack.network.id

  name        = "%s"
  cidr        = "10.0.0.0/24"
  description = "updated description"
  region       = "%s"
}
`, serviceName, region, updatedName, region)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateSubnet(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", subnetName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "description", "initial description"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "checksum"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "name", updatedName),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "description", "updated description"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "id"),
					resource.TestCheckResourceAttrSet("ovh_cloud_network_private_vrack_subnet.test", "checksum"),
					resource.TestCheckResourceAttr("ovh_cloud_network_private_vrack_subnet.test", "resource_status", "READY"),
				),
			},
		},
	})
}

func testSubnetHostRoutes(ctx context.Context, routes ...[2]string) ovhtypes.TfListNestedValue[TargetSpecHostRoutesValue] {
	elems := make([]attr.Value, 0, len(routes))
	for _, route := range routes {
		elems = append(elems, NewTargetSpecHostRoutesValueMust(
			TargetSpecHostRoutesValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"destination": ovhtypes.NewTfStringValue(route[0]),
				"next_hop":    ovhtypes.NewTfStringValue(route[1]),
			},
		))
	}
	return ovhtypes.TfListNestedValue[TargetSpecHostRoutesValue]{
		ListValue: basetypes.NewListValueMust(TargetSpecHostRoutesValue{}.Type(ctx), elems),
	}
}

func testSubnetTargetSpecJSON(t *testing.T, payload *CloudNetworkPrivateSubnetModel) map[string]json.RawMessage {
	t.Helper()
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %s", err)
	}
	var body struct {
		TargetSpec map[string]json.RawMessage `json:"targetSpec"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("unmarshal %s: %s", raw, err)
	}
	return body.TargetSpec
}

func TestCloudNetworkPrivateSubnetModel_hostRoutesPayload(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		name       string
		hostRoutes ovhtypes.TfListNestedValue[TargetSpecHostRoutesValue]
		want       string
	}{
		{"unknown is omitted", ovhtypes.NewListNestedObjectValueOfUnknown[TargetSpecHostRoutesValue](ctx), ""},
		{"null is omitted", ovhtypes.NewListNestedObjectValueOfNull[TargetSpecHostRoutesValue](ctx), ""},
		{"empty list is sent as []", testSubnetHostRoutes(ctx), "[]"},
		{"list is sent", testSubnetHostRoutes(ctx, [2]string{"192.168.10.0/24", "10.0.0.254"}), `[{"destination":"192.168.10.0/24","nextHop":"10.0.0.254"}]`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := CloudNetworkPrivateSubnetModel{
				Name:       ovhtypes.NewTfStringValue("subnet"),
				Cidr:       ovhtypes.NewTfStringValue("10.0.0.0/24"),
				Region:     ovhtypes.NewTfStringValue("GRA11"),
				HostRoutes: tc.hostRoutes,
			}
			for kind, payload := range map[string]*CloudNetworkPrivateSubnetModel{"create": model.ToCreate(), "update": model.ToUpdate()} {
				got, ok := testSubnetTargetSpecJSON(t, payload)["hostRoutes"]
				if tc.want == "" {
					if ok {
						t.Errorf("%s: hostRoutes = %s, want omitted", kind, got)
					}
					continue
				}
				if string(got) != tc.want {
					t.Errorf("%s: hostRoutes = %s, want %s", kind, got, tc.want)
				}
			}
		})
	}
}

func TestCloudNetworkPrivateSubnetModel_hostRoutesFromAPI(t *testing.T) {
	cases := []struct {
		name string
		body string
		want int
	}{
		{"absent", `{"targetSpec":{"name":"subnet"}}`, -1},
		{"null", `{"targetSpec":{"name":"subnet","hostRoutes":null}}`, -1},
		{"empty", `{"targetSpec":{"name":"subnet","hostRoutes":[]}}`, 0},
		{"one route", `{"targetSpec":{"name":"subnet","hostRoutes":[{"destination":"192.168.10.0/24","nextHop":"10.0.0.254"}]}}`, 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var model CloudNetworkPrivateSubnetModel
			if err := json.Unmarshal([]byte(tc.body), &model); err != nil {
				t.Fatalf("unmarshal: %s", err)
			}
			if tc.want < 0 {
				if !model.HostRoutes.IsNull() {
					t.Fatalf("host_routes = %s, want null", model.HostRoutes)
				}
				return
			}
			if got := len(model.HostRoutes.Elements()); got != tc.want {
				t.Fatalf("len(host_routes) = %d, want %d", got, tc.want)
			}
			if tc.want == 1 {
				route := model.HostRoutes.Elements()[0].(TargetSpecHostRoutesValue)
				if route.Destination.ValueString() != "192.168.10.0/24" || route.NextHop.ValueString() != "10.0.0.254" {
					t.Fatalf("route = %s/%s", route.Destination.ValueString(), route.NextHop.ValueString())
				}
			}
		})
	}
}

func TestCloudNetworkPrivateSubnetModel_alignEmptyHostRoutes(t *testing.T) {
	ctx := context.Background()
	null := ovhtypes.NewListNestedObjectValueOfNull[TargetSpecHostRoutesValue](ctx)
	unknown := ovhtypes.NewListNestedObjectValueOfUnknown[TargetSpecHostRoutesValue](ctx)
	empty := testSubnetHostRoutes(ctx)
	one := testSubnetHostRoutes(ctx, [2]string{"192.168.10.0/24", "10.0.0.254"})

	cases := []struct {
		name     string
		got      ovhtypes.TfListNestedValue[TargetSpecHostRoutesValue]
		expected ovhtypes.TfListNestedValue[TargetSpecHostRoutesValue]
		want     ovhtypes.TfListNestedValue[TargetSpecHostRoutesValue]
	}{
		{"null from API keeps planned []", null, empty, empty},
		{"absent from API keeps planned []", ovhtypes.TfListNestedValue[TargetSpecHostRoutesValue]{}, empty, empty},
		{"[] from API keeps prior null", empty, null, null},
		{"null from API with attribute left out stays null", null, null, null},
		{"absent from API with unknown plan is a typed null", ovhtypes.TfListNestedValue[TargetSpecHostRoutesValue]{}, unknown, null},
		{"[] from API with unknown plan stays []", empty, unknown, empty},
		{"routes from API are kept", one, empty, one},
		{"empty from API with routes planned is not hidden", empty, one, empty},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := CloudNetworkPrivateSubnetModel{HostRoutes: tc.got}
			model.alignEmptyHostRoutes(ctx, tc.expected)
			if !model.HostRoutes.Equal(tc.want) {
				t.Fatalf("host_routes = %s, want %s", model.HostRoutes, tc.want)
			}
			if _, err := model.HostRoutes.ToTerraformValue(ctx); err != nil {
				t.Fatalf("ToTerraformValue: %s", err)
			}
		})
	}
}

func TestCloudNetworkPrivateSubnetResourceSchema_hostRoutesOptionalComputed(t *testing.T) {
	ctx := context.Background()
	attribute, ok := CloudNetworkPrivateSubnetResourceSchema(ctx).Attributes["host_routes"].(schema.ListNestedAttribute)
	if !ok {
		t.Fatal("host_routes is not a ListNestedAttribute")
	}
	if !attribute.Optional || !attribute.Computed || len(attribute.PlanModifiers) != 1 {
		t.Fatalf("host_routes: Optional=%t Computed=%t PlanModifiers=%d, want Optional+Computed with one plan modifier", attribute.Optional, attribute.Computed, len(attribute.PlanModifiers))
	}
	if got, want := attribute.PlanModifiers[0].Description(ctx), listplanmodifier.UseStateForUnknown().Description(ctx); got != want {
		t.Fatalf("host_routes plan modifier = %q, want UseStateForUnknown", got)
	}
	for name, nested := range attribute.NestedObject.Attributes {
		if !nested.IsComputed() {
			t.Fatalf("host_routes.%s is not computed: a left-out host_routes would be planned as null", name)
		}
	}
}

func TestSubnetHostRouteFieldsRequired(t *testing.T) {
	ctx := context.Background()
	attrTypes := map[string]attr.Type{"destination": types.StringType, "next_hop": types.StringType}
	for _, tc := range []struct {
		name        string
		destination types.String
		nextHop     types.String
		wantErrors  int
	}{
		{"both set", types.StringValue("192.168.10.0/24"), types.StringValue("10.0.0.254"), 0},
		{"next_hop missing", types.StringValue("192.168.10.0/24"), types.StringNull(), 1},
		{"both missing", types.StringNull(), types.StringNull(), 2},
		{"next_hop unknown", types.StringValue("192.168.10.0/24"), types.StringUnknown(), 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := validator.ObjectRequest{
				Path:        path.Root("host_routes").AtListIndex(0),
				ConfigValue: types.ObjectValueMust(attrTypes, map[string]attr.Value{"destination": tc.destination, "next_hop": tc.nextHop}),
			}
			resp := &validator.ObjectResponse{}
			subnetHostRouteFieldsRequired{}.ValidateObject(ctx, req, resp)
			if got := resp.Diagnostics.ErrorsCount(); got != tc.wantErrors {
				t.Fatalf("errors = %d, want %d: %v", got, tc.wantErrors, resp.Diagnostics)
			}
		})
	}
}

func testSubnetAllocationPools(ctx context.Context, start, end string) ovhtypes.TfListNestedValue[TargetSpecAllocationPoolsValue] {
	pool := NewTargetSpecAllocationPoolsValueMust(
		TargetSpecAllocationPoolsValue{}.AttributeTypes(ctx),
		map[string]attr.Value{
			"start": ovhtypes.NewTfStringValue(start),
			"end":   ovhtypes.NewTfStringValue(end),
		},
	)
	return ovhtypes.TfListNestedValue[TargetSpecAllocationPoolsValue]{
		ListValue: basetypes.NewListValueMust(TargetSpecAllocationPoolsValue{}.Type(ctx), []attr.Value{pool}),
	}
}

func TestCloudNetworkPrivateSubnetModel_updateRenameOnly(t *testing.T) {
	ctx := context.Background()
	pools := testSubnetAllocationPools(ctx, "10.0.0.10", "10.0.0.100")

	plan := CloudNetworkPrivateSubnetModel{
		Name:            ovhtypes.NewTfStringValue("renamed"),
		Cidr:            ovhtypes.NewTfStringValue("10.0.0.0/24"),
		GatewayIp:       ovhtypes.NewTfStringValue("10.0.0.1"),
		AllocationPools: pools,
		DhcpEnabled:     ovhtypes.NewTfBoolValue(true),
		HostRoutes:      testSubnetHostRoutes(ctx, [2]string{"192.168.10.0/24", "10.0.0.254"}),
	}
	plan.fillDhcpEnabledForUpdate(CloudNetworkPrivateSubnetModel{DhcpEnabled: ovhtypes.NewTfBoolValue(false)})

	body := testSubnetTargetSpecJSON(t, plan.ToUpdate())
	if got := string(body["dhcpEnabled"]); got != "true" {
		t.Errorf("dhcpEnabled = %q, want true", got)
	}
	if got := string(body["allocationPools"]); got != `[{"end":"10.0.0.100","start":"10.0.0.10"}]` && got != `[{"start":"10.0.0.10","end":"10.0.0.100"}]` {
		t.Errorf("allocationPools = %s, want the pool from the plan", got)
	}
	if got := string(body["gatewayIp"]); got != `"10.0.0.1"` {
		t.Errorf("gatewayIp = %s", got)
	}
	if got := string(body["hostRoutes"]); got != `[{"destination":"192.168.10.0/24","nextHop":"10.0.0.254"}]` {
		t.Errorf("hostRoutes = %s, want the routes from the plan", got)
	}
	for _, key := range []string{"cidr", "location"} {
		if got, ok := body[key]; ok {
			t.Errorf("%s = %s, want omitted on PUT", key, got)
		}
	}
}

func TestCloudNetworkPrivateSubnetModel_fillDhcpEnabledForUpdate(t *testing.T) {
	enabled := ovhtypes.NewTfBoolValue(true)
	disabled := ovhtypes.NewTfBoolValue(false)
	null := ovhtypes.TfBoolValue{BoolValue: basetypes.NewBoolNull()}
	unknown := ovhtypes.TfBoolValue{BoolValue: basetypes.NewBoolUnknown()}

	cases := []struct {
		name    string
		plan    ovhtypes.TfBoolValue
		state   ovhtypes.TfBoolValue
		current ovhtypes.TfBoolValue
		want    ovhtypes.TfBoolValue
	}{
		{"planned value is kept", disabled, enabled, enabled, disabled},
		{"unknown plan takes state", unknown, disabled, enabled, disabled},
		{"null plan takes state", null, disabled, enabled, disabled},
		{"unknown plan and null state take current state", unknown, null, enabled, enabled},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			model := CloudNetworkPrivateSubnetModel{DhcpEnabled: tc.plan}
			state := CloudNetworkPrivateSubnetModel{
				DhcpEnabled:  tc.state,
				CurrentState: CloudNetworkPrivateSubnetCurrentStateValue{DhcpEnabled: tc.current, state: attr.ValueStateKnown},
			}
			model.fillDhcpEnabledForUpdate(state)
			if !model.DhcpEnabled.Equal(tc.want) {
				t.Fatalf("dhcp_enabled = %s, want %s", model.DhcpEnabled, tc.want)
			}
		})
	}
}

func TestAccCloudNetworkPrivateVrackSubnet_hostRoutes(t *testing.T) {
	serviceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE_TEST")
	region := os.Getenv("OVH_CLOUD_PROJECT_REGION_TEST")

	subnetName := acctest.RandomWithPrefix(testAccResourceCloudNetworkPrivateSubnetNamePrefix)

	config := func(description, hostRoutes string) string {
		return fmt.Sprintf(`
resource "ovh_cloud_network_private_vrack" "network" {
  service_name = "%s"
  name         = "terraform_testacc_private_net"
  region       = "%s"
}

resource "ovh_cloud_network_private_vrack_subnet" "test" {
  service_name = ovh_cloud_network_private_vrack.network.service_name
  network_id   = ovh_cloud_network_private_vrack.network.id

  name        = "%s"
  cidr        = "10.0.0.0/24"
  region      = "%s"
  description = "%s"
%s
}
`, serviceName, region, subnetName, region, description, hostRoutes)
	}

	route := func(destination, nextHop string) string {
		return fmt.Sprintf(`    {
      destination = "%s"
      next_hop    = "%s"
    },
`, destination, nextHop)
	}

	hostRoutes := func(routes ...string) string {
		return "  host_routes = [\n" + strings.Join(routes, "") + "  ]"
	}

	const resourceName = "ovh_cloud_network_private_vrack_subnet.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckCloudNetworkPrivateSubnet(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("initial", hostRoutes(route("192.168.10.0/24", "10.0.0.254"))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "host_routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "host_routes.0.destination", "192.168.10.0/24"),
					resource.TestCheckResourceAttr(resourceName, "host_routes.0.next_hop", "10.0.0.254"),
					resource.TestCheckResourceAttr(resourceName, "current_state.host_routes.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "resource_status", "READY"),
				),
			},
			{
				Config: config("initial", hostRoutes(route("192.168.20.0/24", "10.0.0.253"), route("192.168.30.0/24", "10.0.0.253"))),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "host_routes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "host_routes.0.destination", "192.168.20.0/24"),
					resource.TestCheckResourceAttr(resourceName, "host_routes.1.destination", "192.168.30.0/24"),
					resource.TestCheckResourceAttr(resourceName, "current_state.host_routes.#", "2"),
				),
			},
			{
				Config: config("updated", ""),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "updated"),
					resource.TestCheckResourceAttr(resourceName, "host_routes.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "host_routes.0.destination", "192.168.20.0/24"),
					resource.TestCheckResourceAttr(resourceName, "host_routes.1.destination", "192.168.30.0/24"),
					resource.TestCheckResourceAttr(resourceName, "current_state.host_routes.#", "2"),
				),
			},
			{
				Config: config("updated", "  host_routes = []"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "host_routes.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "current_state.host_routes.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCloudNetworkPrivateSubnetImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{"host_routes"},
			},
		},
	})
}
