package ovh

import (
	"context"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	dsschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// vpsTestAttributeInfo is what the tests compare between two schemas.
type vpsTestAttributeInfo struct {
	typ                     attr.Type
	description, markdown   string
	computed, optional, req bool
}

// vpsTestResourceAttributes flattens resource attributes by dotted path.
func vpsTestResourceAttributes(prefix string, attrs map[string]schema.Attribute, out map[string]vpsTestAttributeInfo) map[string]vpsTestAttributeInfo {
	for name, a := range attrs {
		p := prefix + name
		out[p] = vpsTestAttributeInfo{a.GetType(), a.GetDescription(), a.GetMarkdownDescription(), a.IsComputed(), a.IsOptional(), a.IsRequired()}
		switch nested := a.(type) {
		case schema.SingleNestedAttribute:
			vpsTestResourceAttributes(p+".", nested.Attributes, out)
		case schema.ListNestedAttribute:
			vpsTestResourceAttributes(p+"[].", nested.NestedObject.Attributes, out)
		}
	}
	return out
}

// vpsTestDataSourceAttributes flattens data source attributes by dotted path.
func vpsTestDataSourceAttributes(prefix string, attrs map[string]dsschema.Attribute, out map[string]vpsTestAttributeInfo) map[string]vpsTestAttributeInfo {
	for name, a := range attrs {
		p := prefix + name
		out[p] = vpsTestAttributeInfo{a.GetType(), a.GetDescription(), a.GetMarkdownDescription(), a.IsComputed(), a.IsOptional(), a.IsRequired()}
		switch nested := a.(type) {
		case dsschema.SingleNestedAttribute:
			vpsTestDataSourceAttributes(p+".", nested.Attributes, out)
		case dsschema.ListNestedAttribute:
			vpsTestDataSourceAttributes(p+"[].", nested.NestedObject.Attributes, out)
		}
	}
	return out
}

// vpsTestResourceAttribute returns the resource attribute at a dotted path.
func vpsTestResourceAttribute(t *testing.T, dotted string) schema.Attribute {
	t.Helper()
	attrs := vpsResourceAttributes()
	var a schema.Attribute
	for _, name := range strings.Split(dotted, ".") {
		if attrs == nil {
			t.Fatalf("%s: %T has no nested attribute", dotted, a)
		}
		var ok bool
		if a, ok = attrs[name]; !ok {
			t.Fatalf("%s: no attribute %q", dotted, name)
		}
		attrs = nil
		if nested, ok := a.(schema.SingleNestedAttribute); ok {
			attrs = nested.Attributes
		}
	}
	return a
}

func vpsTestSchema() schema.Schema {
	return schema.Schema{Attributes: vpsResourceAttributes()}
}

// vpsTestPlan returns a plan of the resource schema with the given values,
// null elsewhere.
func vpsTestPlan(t *testing.T, values map[string]attr.Value) tfsdk.Plan {
	t.Helper()
	ctx := context.Background()
	s := vpsTestSchema()
	plan := tfsdk.Plan{Schema: s, Raw: tftypes.NewValue(s.Type().TerraformType(ctx), nil)}
	for dotted, value := range values {
		p := path.Empty()
		for _, name := range strings.Split(dotted, ".") {
			p = p.AtName(name)
		}
		if diags := plan.SetAttribute(ctx, p, value); diags.HasError() {
			t.Fatalf("setting %s: %v", dotted, diags)
		}
	}
	return plan
}

func TestVPSResourceSchemaMatchesModel(t *testing.T) {
	ctx := context.Background()
	attrs := vpsResourceAttributes()
	attrTypes := vpsResourceAttrTypes()

	for name, a := range attrs {
		want, ok := attrTypes[name]
		if !ok {
			t.Errorf("attribute %s: missing from vpsResourceAttrTypes", name)
			continue
		}
		if got := a.GetType(); !got.Equal(want) {
			t.Errorf("attribute %s: schema type %s, vpsResourceAttrTypes %s", name, got, want)
		}
	}
	for name := range attrTypes {
		if _, ok := attrs[name]; !ok {
			t.Errorf("vpsResourceAttrTypes %s: no such attribute in the schema", name)
		}
	}

	var tags []string
	model := reflect.TypeOf(VPSResourceModel{})
	for i := range model.NumField() {
		tags = append(tags, model.Field(i).Tag.Get("tfsdk"))
	}
	for name := range attrs {
		if !slices.Contains(tags, name) {
			t.Errorf("attribute %s: no field of VPSResourceModel", name)
		}
	}
	if len(tags) != len(attrs) {
		t.Errorf("VPSResourceModel has %d fields, the schema %d attributes", len(tags), len(attrs))
	}

	// The model reads a plan of this schema.
	var m VPSResourceModel
	plan := vpsTestPlan(t, map[string]attr.Value{"plan_code": ovhtypes.NewTfStringValue("vps-2025-model1")})
	if diags := plan.Get(ctx, &m); diags.HasError() {
		t.Fatalf("reading a plan into VPSResourceModel: %v", diags)
	}
}

func TestVPSResourceCurrentStateMatchesDataSource(t *testing.T) {
	got := vpsTestResourceAttributes("", vpsResourceCurrentStateAttributes(), map[string]vpsTestAttributeInfo{})
	want := vpsTestDataSourceAttributes("", vpsCurrentStateAttributes(), map[string]vpsTestAttributeInfo{})

	if !reflect.DeepEqual(got, want) {
		for p, w := range want {
			if g, ok := got[p]; !ok || !reflect.DeepEqual(g, w) {
				t.Errorf("current_state.%s: got %+v, want %+v", p, g, w)
			}
		}
		for p := range got {
			if _, ok := want[p]; !ok {
				t.Errorf("current_state.%s: not in the data source", p)
			}
		}
	}

	attrTypes := vpsResourceAttrTypes()["current_state"].(types.ObjectType).AttrTypes
	if !reflect.DeepEqual(attrTypes, vpsCurrentStateAttrTypes()) {
		t.Errorf("current_state types: got %v, want %v", attrTypes, vpsCurrentStateAttrTypes())
	}
}

func TestVPSResourceSchemaHasNoServer(t *testing.T) {
	for p := range vpsTestResourceAttributes("", vpsResourceAttributes(), map[string]vpsTestAttributeInfo{}) {
		if strings.Contains(strings.ToLower(p), "server") {
			t.Errorf("attribute %s contains \"server\"", p)
		}
	}
}

func TestVPSResourceSchemaModes(t *testing.T) {
	attrs := vpsTestResourceAttributes("", vpsResourceAttributes(), map[string]vpsTestAttributeInfo{})
	// required, optional, computed
	want := map[string][3]bool{
		"id":                                    {false, false, true},
		"plan_code":                             {true, false, false},
		"location":                              {true, false, false},
		"location.region":                       {true, false, false},
		"display_name":                          {false, true, true},
		"power_state":                           {false, true, true},
		"installation":                          {true, false, false},
		"installation.from":                     {true, false, false},
		"installation.from.image_id":            {false, true, true},
		"installation.from.image_name":          {false, true, true},
		"installation.ssh_key":                  {false, true, false},
		"installation.post_installation_script": {false, true, false},
		"installation.send_password_by_email":   {false, true, true},
		"options":                               {false, true, true},
		"options.auto_backup":                   {false, true, true},
		"options.auto_backup.cron":              {false, true, true},
		"order":                                 {false, true, true},
		"order.duration":                        {false, true, true},
		"order.pricing_mode":                    {false, true, true},
		"order.order_id":                        {false, false, true},
		"timeouts":                              {false, true, false},
		"timeouts.create":                       {false, true, false},
		"timeouts.update":                       {false, true, false},
		"timeouts.delete":                       {false, true, false},
		"current_state":                         {false, false, true},
	}
	for p, w := range want {
		a, ok := attrs[p]
		if !ok {
			t.Errorf("%s: missing", p)
			continue
		}
		if got := [3]bool{a.req, a.optional, a.computed}; got != w {
			t.Errorf("%s: (required, optional, computed) = %v, want %v", p, got, w)
		}
	}
	for p := range attrs {
		if _, ok := want[p]; !ok && !strings.HasPrefix(p, "current_state.") {
			t.Errorf("%s: not in the plan of the resource", p)
		}
	}

	if !vpsTestResourceAttribute(t, "installation.post_installation_script").IsSensitive() {
		t.Error("installation.post_installation_script is not sensitive")
	}
}

func TestVPSRequiresReplaceIfPriorKnown(t *testing.T) {
	ctx := context.Background()
	someState := tfsdk.State{Raw: tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{}}, map[string]tftypes.Value{})}
	somePlan := tfsdk.Plan{Raw: someState.Raw}

	cases := []struct {
		name        string
		state, plan types.String
		replace     bool
	}{
		{"null state (import, migration)", types.StringNull(), types.StringValue("eu-west-gra"), false},
		{"changed", types.StringValue("eu-west-gra"), types.StringValue("eu-west-sbg"), true},
		{"unchanged", types.StringValue("eu-west-gra"), types.StringValue("eu-west-gra"), false},
	}
	for _, c := range cases {
		req := planmodifier.StringRequest{State: someState, Plan: somePlan, StateValue: c.state, PlanValue: c.plan, ConfigValue: c.plan}
		resp := &planmodifier.StringResponse{PlanValue: c.plan}
		vpsRequiresReplaceIfPriorKnown().PlanModifyString(ctx, req, resp)
		if resp.RequiresReplace != c.replace {
			t.Errorf("%s: RequiresReplace = %t, want %t", c.name, resp.RequiresReplace, c.replace)
		}
	}

	for _, p := range []string{"plan_code", "location.region"} {
		a := vpsTestResourceAttribute(t, p).(schema.StringAttribute)
		if len(a.PlanModifiers) != 1 || a.PlanModifiers[0].Description(ctx) != vpsRequiresReplaceIfPriorKnown().Description(ctx) {
			t.Errorf("%s: plan modifiers %v, want vpsRequiresReplaceIfPriorKnown", p, a.PlanModifiers)
		}
	}
}

// vpsTestValidateString runs the validators of a string attribute.
func vpsTestValidateString(t *testing.T, config tfsdk.Config, dotted string) diag.Diagnostics {
	t.Helper()
	ctx := context.Background()
	p := path.Empty()
	for _, name := range strings.Split(dotted, ".") {
		p = p.AtName(name)
	}
	var value ovhtypes.TfStringValue
	if diags := config.GetAttribute(ctx, p, &value); diags.HasError() {
		t.Fatalf("reading %s: %v", dotted, diags)
	}
	var diags diag.Diagnostics
	for _, v := range vpsTestResourceAttribute(t, dotted).(schema.StringAttribute).Validators {
		req := validator.StringRequest{Path: p, PathExpression: p.Expression(), Config: config, ConfigValue: value.StringValue}
		resp := &validator.StringResponse{}
		v.ValidateString(ctx, req, resp)
		diags.Append(resp.Diagnostics...)
	}
	return diags
}

func TestVPSInstallationFromExactlyOne(t *testing.T) {
	cases := []struct {
		name   string
		values map[string]attr.Value
		ok     bool
	}{
		{"image_name", map[string]attr.Value{"installation.from.image_name": ovhtypes.NewTfStringValue("Debian 12")}, true},
		{"image_id", map[string]attr.Value{"installation.from.image_id": ovhtypes.NewTfStringValue("3c1e5a4e-8d55-4f39-9b5c-9d9c3a1b2c3d")}, true},
		{"none", map[string]attr.Value{"installation.from": types.ObjectValueMust(vpsInstallationFromAttrTypes(), map[string]attr.Value{
			"image_id":   ovhtypes.TfStringValue{StringValue: types.StringNull()},
			"image_name": ovhtypes.TfStringValue{StringValue: types.StringNull()},
		})}, false},
		{"both", map[string]attr.Value{
			"installation.from.image_name": ovhtypes.NewTfStringValue("Debian 12"),
			"installation.from.image_id":   ovhtypes.NewTfStringValue("3c1e5a4e-8d55-4f39-9b5c-9d9c3a1b2c3d"),
		}, false},
	}
	for _, c := range cases {
		plan := vpsTestPlan(t, c.values)
		config := tfsdk.Config{Schema: plan.Schema, Raw: plan.Raw}
		var diags diag.Diagnostics
		diags.Append(vpsTestValidateString(t, config, "installation.from.image_name")...)
		diags.Append(vpsTestValidateString(t, config, "installation.from.image_id")...)
		if diags.HasError() == c.ok {
			t.Errorf("%s: errors %v, want ok = %t", c.name, diags, c.ok)
		}
	}
}

func TestVPSResourceStringValidators(t *testing.T) {
	cases := []struct {
		attribute, value string
		ok               bool
	}{
		{"display_name", strings.Repeat("a", 50), true},
		{"display_name", strings.Repeat("a", 51), false},
		{"power_state", "ACTIVE", true},
		{"power_state", "RESCUE", true},
		{"power_state", "SHUTOFF", true},
		{"power_state", "UNKNOWN", false},
		{"timeouts.create", "1h30m", true},
		{"timeouts.update", "45m", true},
		{"timeouts.delete", "30 minutes", false},
	}
	for _, c := range cases {
		plan := vpsTestPlan(t, map[string]attr.Value{c.attribute: ovhtypes.NewTfStringValue(c.value)})
		diags := vpsTestValidateString(t, tfsdk.Config{Schema: plan.Schema, Raw: plan.Raw}, c.attribute)
		if diags.HasError() == c.ok {
			t.Errorf("%s = %q: errors %v, want ok = %t", c.attribute, c.value, diags, c.ok)
		}
	}
}

func TestVPSResourceTimeout(t *testing.T) {
	ctx := context.Background()
	var m VPSResourceModel
	plan := vpsTestPlan(t, map[string]attr.Value{"timeouts.create": ovhtypes.NewTfStringValue("90m")})
	if diags := plan.Get(ctx, &m); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if got := m.timeout("create", time.Hour); got != 90*time.Minute {
		t.Errorf("create: got %s, want 1h30m", got)
	}
	if got := m.timeout("delete", 30*time.Minute); got != 30*time.Minute {
		t.Errorf("delete: got %s, want the default 30m", got)
	}

	m = VPSResourceModel{Timeouts: types.ObjectNull(vpsTimeoutsAttrTypes())}
	if got := m.timeout("update", time.Hour); got != time.Hour {
		t.Errorf("no timeouts: got %s, want the default 1h", got)
	}
}

// TestVPSCurrentStateUnknownWhenUpdatePlanned checks that current_state keeps
// its state value unless the plan changes the VPS.
func TestVPSCurrentStateUnknownWhenUpdatePlanned(t *testing.T) {
	ctx := context.Background()
	csType := types.ObjectType{AttrTypes: vpsCurrentStateAttrTypes()}
	csState := types.ObjectNull(vpsCurrentStateAttrTypes())
	csUnknown := types.ObjectUnknown(vpsCurrentStateAttrTypes())
	stateValues := map[string]attr.Value{
		"display_name":                        ovhtypes.NewTfStringValue("web-01"),
		"power_state":                         ovhtypes.NewTfStringValue("ACTIVE"),
		"installation.from.image_name":        ovhtypes.NewTfStringValue("Debian 12"),
		"installation.from.image_id":          ovhtypes.NewTfStringValue("3c1e5a4e-8d55-4f39-9b5c-9d9c3a1b2c3d"),
		"installation.send_password_by_email": types.BoolValue(true),
	}
	with := func(changes map[string]attr.Value) map[string]attr.Value {
		values := map[string]attr.Value{}
		for k, v := range stateValues {
			values[k] = v
		}
		for k, v := range changes {
			values[k] = v
		}
		return values
	}
	null := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	unknown := ovhtypes.TfStringValue{StringValue: types.StringUnknown()}

	cases := []struct {
		name         string
		config, plan map[string]attr.Value
		wantUnknown  bool
	}{
		{"no change", stateValues, stateValues, false},
		{"display_name changed", with(map[string]attr.Value{"display_name": ovhtypes.NewTfStringValue("web-02")}), with(map[string]attr.Value{"display_name": ovhtypes.NewTfStringValue("web-02")}), true},
		{"ssh_key set", with(map[string]attr.Value{"installation.ssh_key": ovhtypes.NewTfStringValue("ssh-ed25519 AAAA")}), with(map[string]attr.Value{"installation.ssh_key": ovhtypes.NewTfStringValue("ssh-ed25519 AAAA")}), true},
		{"only timeouts changed", with(map[string]attr.Value{"timeouts.create": ovhtypes.NewTfStringValue("2h")}), with(map[string]attr.Value{"timeouts.create": ovhtypes.NewTfStringValue("2h")}), false},
		{
			"computed left out, marked unknown by the framework",
			with(map[string]attr.Value{"installation.from.image_id": null, "timeouts.create": ovhtypes.NewTfStringValue("2h")}),
			with(map[string]attr.Value{"installation.from.image_id": unknown, "timeouts.create": ovhtypes.NewTfStringValue("2h")}),
			false,
		},
	}
	for _, c := range cases {
		state := vpsTestPlan(t, stateValues)
		config := vpsTestPlan(t, c.config)
		plan := vpsTestPlan(t, c.plan)
		req := planmodifier.ObjectRequest{
			State:       tfsdk.State{Schema: state.Schema, Raw: state.Raw},
			Plan:        plan,
			Config:      tfsdk.Config{Schema: config.Schema, Raw: config.Raw},
			StateValue:  csState,
			PlanValue:   csUnknown,
			ConfigValue: types.ObjectNull(csType.AttrTypes),
		}
		resp := &planmodifier.ObjectResponse{PlanValue: csUnknown}
		vpsCurrentStatePlanModifier().PlanModifyObject(ctx, req, resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("%s: unexpected error: %v", c.name, resp.Diagnostics)
		}
		if got := resp.PlanValue.IsUnknown(); got != c.wantUnknown {
			t.Errorf("%s: current_state unknown = %t, want %t", c.name, got, c.wantUnknown)
		}
	}

	// Creation: no state, current_state stays unknown.
	plan := vpsTestPlan(t, stateValues)
	req := planmodifier.ObjectRequest{
		State:      tfsdk.State{Schema: plan.Schema, Raw: tftypes.NewValue(plan.Raw.Type(), nil)},
		Plan:       plan,
		Config:     tfsdk.Config{Schema: plan.Schema, Raw: plan.Raw},
		StateValue: types.ObjectNull(csType.AttrTypes),
		PlanValue:  csUnknown,
	}
	resp := &planmodifier.ObjectResponse{PlanValue: csUnknown}
	vpsCurrentStatePlanModifier().PlanModifyObject(ctx, req, resp)
	if !resp.PlanValue.IsUnknown() {
		t.Errorf("creation: current_state = %v, want unknown", resp.PlanValue)
	}
}
