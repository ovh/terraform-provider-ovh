package ovh

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// vpsTestServer is the VPS of newVPSTestAPI, with the automated backup option.
func vpsTestServer() *VPSAPIServer {
	active := "ACTIVE"
	imageId := "3c1e5a4e-8d55-4f39-9b5c-9d9c3a1b2c3d"
	imageName := "Debian 12"
	region := "eu-west-gra"
	cron := "0 3 * * *"
	rotation := int64(7)
	from := VPSAPIInstallationFrom{ImageId: &imageId, ImageName: &imageName}

	return &VPSAPIServer{
		Id: vpsTestId,
		TargetSpec: VPSAPIReadTargetSpec{
			DisplayName:  "web-01",
			PowerState:   &active,
			Installation: VPSAPIReadInstallation{From: from},
			Options:      VPSAPIReadOptions{AutoBackup: &VPSAPIAutoBackup{Cron: cron}},
		},
		CurrentState: &VPSAPICurrentState{
			DisplayName:  "web-01",
			PowerState:   "ACTIVE",
			Installation: VPSAPICurrentInstallation{From: from},
			Location:     VPSAPILocation{Region: &region},
			Model:        VPSAPIHardware{Vcpus: 4, Ram: 8192, Disk: 160},
			IpAddresses:  []VPSAPIIPAddress{{Ip: "51.75.1.2", Version: 4}, {Ip: "2001:41d0::1", Version: 6}},
			Options: VPSAPIOptions{
				AutoBackup: VPSAPIAutoBackupOption{Enabled: true, Cron: &cron, Rotation: &rotation},
				Snapshot:   VPSAPISnapshotOption{Enabled: true},
			},
		},
	}
}

// vpsTestAttr walks an object along names, failing the test on a missing or
// null step.
func vpsTestAttr(t *testing.T, obj basetypes.ObjectValue, names ...string) attr.Value {
	t.Helper()
	var value attr.Value = obj
	for i, name := range names {
		o, ok := value.(basetypes.ObjectValue)
		if !ok || o.IsNull() {
			t.Fatalf("%v: not an object or null before %q", names[:i], name)
		}
		value, ok = o.Attributes()[name]
		if !ok {
			t.Fatalf("%v: no attribute %q", names[:i], name)
		}
	}
	return value
}

func vpsTestString(t *testing.T, value attr.Value) string {
	t.Helper()
	s, ok := value.(basetypes.StringValuable)
	if !ok {
		t.Fatalf("%v is not a string", value)
	}
	v, _ := s.ToStringValue(context.Background())
	return v.ValueString()
}

func TestVPSDataModelMergeWith(t *testing.T) {
	var m VPSDataModel
	m.MergeWith(context.Background(), vpsTestServer())

	if got := m.Id.ValueString(); got != vpsTestId {
		t.Errorf("id = %q, want %q", got, vpsTestId)
	}
	if got := m.DisplayName.ValueString(); got != "web-01" {
		t.Errorf("display_name = %q, want web-01", got)
	}
	if got := m.PowerState.ValueString(); got != "ACTIVE" {
		t.Errorf("power_state = %q, want ACTIVE", got)
	}
	if got := vpsTestString(t, vpsTestAttr(t, m.Location, "region")); got != "eu-west-gra" {
		t.Errorf("location.region = %q, want eu-west-gra", got)
	}
	if got := vpsTestString(t, vpsTestAttr(t, m.Installation, "from", "image_name")); got != "Debian 12" {
		t.Errorf("installation.from.image_name = %q, want Debian 12", got)
	}
	if got := vpsTestString(t, vpsTestAttr(t, m.Options, "auto_backup", "cron")); got != "0 3 * * *" {
		t.Errorf("options.auto_backup.cron = %q, want 0 3 * * *", got)
	}

	cs := m.CurrentState
	if got := vpsTestString(t, vpsTestAttr(t, cs, "location", "region")); got != "eu-west-gra" {
		t.Errorf("current_state.location.region = %q, want eu-west-gra", got)
	}
	if got := vpsTestString(t, vpsTestAttr(t, cs, "power_state")); got != "ACTIVE" {
		t.Errorf("current_state.power_state = %q, want ACTIVE", got)
	}
	if got := vpsTestAttr(t, cs, "model", "ram").(types.Int64).ValueInt64(); got != 8192 {
		t.Errorf("current_state.model.ram = %d, want 8192", got)
	}
	ips := vpsTestAttr(t, cs, "ip_addresses").(types.List).Elements()
	if len(ips) != 2 {
		t.Fatalf("current_state.ip_addresses has %d elements, want 2", len(ips))
	}
	if got := vpsTestAttr(t, ips[1].(types.Object), "version").(types.Int64).ValueInt64(); got != 6 {
		t.Errorf("current_state.ip_addresses[1].version = %d, want 6", got)
	}
	if got := vpsTestAttr(t, cs, "options", "auto_backup", "rotation").(types.Int64).ValueInt64(); got != 7 {
		t.Errorf("current_state.options.auto_backup.rotation = %d, want 7", got)
	}
	if !vpsTestAttr(t, cs, "options", "snapshot", "enabled").(types.Bool).ValueBool() {
		t.Error("current_state.options.snapshot.enabled = false, want true")
	}
	if size := vpsTestAttr(t, cs, "options", "additional_disks", "size"); !size.IsNull() {
		t.Errorf("current_state.options.additional_disks.size = %v, want null", size)
	}
}

func TestVPSDataModelMergeWithUnknownValues(t *testing.T) {
	// API v1: transitional state, unknown image and region, no automated backup.
	srv := vpsTestServer()
	srv.TargetSpec.PowerState = nil
	srv.TargetSpec.Installation.From = VPSAPIInstallationFrom{}
	srv.TargetSpec.Options.AutoBackup = nil
	srv.CurrentState.PowerState = "UNKNOWN"
	srv.CurrentState.Installation.From = VPSAPIInstallationFrom{}
	srv.CurrentState.Location.Region = nil

	var m VPSDataModel
	m.MergeWith(context.Background(), srv)

	if !m.PowerState.IsNull() {
		t.Errorf("power_state = %q, want null", m.PowerState.ValueString())
	}
	if region := vpsTestAttr(t, m.Location, "region"); !region.IsNull() {
		t.Errorf("location.region = %v, want null", region)
	}
	if name := vpsTestAttr(t, m.Installation, "from", "image_name"); !name.IsNull() {
		t.Errorf("installation.from.image_name = %v, want null", name)
	}
	if autoBackup := vpsTestAttr(t, m.Options, "auto_backup"); !autoBackup.IsNull() {
		t.Errorf("options.auto_backup = %v, want null", autoBackup)
	}
	if id := vpsTestAttr(t, m.CurrentState, "installation", "from", "image_id"); !id.IsNull() {
		t.Errorf("current_state.installation.from.image_id = %v, want null", id)
	}
}

func TestVPSDataModelMergeWithNilCurrentState(t *testing.T) {
	srv := vpsTestServer()
	srv.CurrentState = nil

	var m VPSDataModel
	m.MergeWith(context.Background(), srv)

	if !m.CurrentState.IsNull() {
		t.Errorf("current_state = %v, want null", m.CurrentState)
	}
	if !m.Location.IsNull() {
		t.Errorf("location = %v, want null", m.Location)
	}
	if got := m.DisplayName.ValueString(); got != "web-01" {
		t.Errorf("display_name = %q, want web-01", got)
	}
}

// vpsTestSetState writes value into a state of the data source schema, which
// fails when the model and the schema disagree.
func vpsTestSetState(t *testing.T, d datasource.DataSource, value any) {
	t.Helper()
	ctx := context.Background()
	var schemaResp datasource.SchemaResponse
	d.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("schema: %v", schemaResp.Diagnostics)
	}
	if diags := schemaResp.Schema.ValidateImplementation(ctx); diags.HasError() {
		t.Fatalf("schema: %v", diags)
	}

	state := tfsdk.State{
		Schema: schemaResp.Schema,
		Raw:    tftypes.NewValue(schemaResp.Schema.Type().TerraformType(ctx), nil),
	}
	if diags := state.Set(ctx, value); diags.HasError() {
		t.Fatalf("setting the state: %v", diags)
	}
}

func TestVPSDataSourceSchemaMatchesModel(t *testing.T) {
	var m VPSDataModel
	m.MergeWith(context.Background(), vpsTestServer())
	vpsTestSetState(t, NewVPSDataSource(), &m)

	srv := vpsTestServer()
	srv.CurrentState = nil
	var empty VPSDataModel
	empty.MergeWith(context.Background(), srv)
	vpsTestSetState(t, NewVPSDataSource(), &empty)
}

func TestVPSsDataSourceSchemaMatchesModel(t *testing.T) {
	second := vpsTestServer()
	second.Id = "vps-5e6f7a8b.vps.ovh.net"
	second.CurrentState = nil

	m := newVPSsDataModel(context.Background(), []VPSAPIServer{*vpsTestServer(), *second})
	if n := len(m.Vps.Elements()); n != 2 {
		t.Fatalf("vps has %d elements, want 2", n)
	}
	vpsTestSetState(t, NewVPSsDataSource(), &m)

	empty := newVPSsDataModel(context.Background(), nil)
	if empty.Vps.IsNull() || len(empty.Vps.Elements()) != 0 {
		t.Errorf("vps = %v, want an empty list", empty.Vps)
	}
	vpsTestSetState(t, NewVPSsDataSource(), &empty)
}
