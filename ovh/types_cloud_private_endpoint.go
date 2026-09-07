package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudPrivateEndpointModel is the Terraform model for the private
// endpoint resource. targetSpec fields are exposed at the root; the resolved
// state stays in the nested current_state object.
type CloudPrivateEndpointModel struct {
	ServiceName      ovhtypes.TfStringValue `tfsdk:"service_name"`
	Name             ovhtypes.TfStringValue `tfsdk:"name"`
	Description      ovhtypes.TfStringValue `tfsdk:"description"`
	Region           ovhtypes.TfStringValue `tfsdk:"region"`
	Service          types.Object           `tfsdk:"service"`
	NetworkId        ovhtypes.TfStringValue `tfsdk:"network_id"`
	SubnetId         ovhtypes.TfStringValue `tfsdk:"subnet_id"`
	VipAddress       ovhtypes.TfStringValue `tfsdk:"vip_address"`
	PublicDnsEnabled types.Bool             `tfsdk:"public_dns_enabled"`

	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

// API types mirror the publicCloud.privateEndpoint models.

type CloudPrivateEndpointRef struct {
	Id string `json:"id,omitempty"`
}

// id returns the referenced identifier, or an empty string when the API omitted
// the reference.
func (r *CloudPrivateEndpointRef) id() string {
	if r == nil {
		return ""
	}
	return r.Id
}

type CloudPrivateEndpointServiceRef struct {
	Name string `json:"name"`
	Tier string `json:"tier"`
}

type CloudPrivateEndpointServiceState struct {
	Name      string `json:"name"`
	Tier      string `json:"tier"`
	Port      int64  `json:"port,omitempty"`
	Transport string `json:"transport,omitempty"`
	// Null when the catalog does not declare it.
	Bandwidth *int64 `json:"bandwidth,omitempty"`
}

type CloudPrivateEndpointLocation struct {
	Region string `json:"region,omitempty"`
}

type CloudPrivateEndpointTargetSpec struct {
	Name             string                          `json:"name,omitempty"`
	Description      string                          `json:"description,omitempty"`
	Location         *CloudPrivateEndpointLocation   `json:"location,omitempty"`
	Service          *CloudPrivateEndpointServiceRef `json:"service,omitempty"`
	Network          *CloudPrivateEndpointRef        `json:"network,omitempty"`
	Subnet           *CloudPrivateEndpointRef        `json:"subnet,omitempty"`
	VipAddress       string                          `json:"vipAddress,omitempty"`
	PublicDnsEnabled bool                            `json:"publicDnsEnabled"`
}

type CloudPrivateEndpointCurrentState struct {
	Name             string                            `json:"name,omitempty"`
	Description      string                            `json:"description,omitempty"`
	EndpointId       string                            `json:"endpointId,omitempty"`
	Location         *CloudPrivateEndpointLocation     `json:"location,omitempty"`
	Service          *CloudPrivateEndpointServiceState `json:"service,omitempty"`
	Network          *CloudPrivateEndpointRef          `json:"network,omitempty"`
	Subnet           *CloudPrivateEndpointRef          `json:"subnet,omitempty"`
	VipAddress       string                            `json:"vipAddress,omitempty"`
	Hostname         string                            `json:"hostname,omitempty"`
	PublicDnsEnabled bool                              `json:"publicDnsEnabled"`
}

type CloudPrivateEndpointAPIResponse struct {
	Id             string                            `json:"id"`
	Checksum       string                            `json:"checksum"`
	CreatedAt      string                            `json:"createdAt"`
	UpdatedAt      string                            `json:"updatedAt"`
	ResourceStatus string                            `json:"resourceStatus"`
	CurrentState   *CloudPrivateEndpointCurrentState `json:"currentState,omitempty"`
	TargetSpec     *CloudPrivateEndpointTargetSpec   `json:"targetSpec,omitempty"`
	CurrentTasks   []CloudResourceTask               `json:"currentTasks"`
}

// Discovery types for the reference/privateEndpoint/serviceAvailability listing.

type CloudPrivateEndpointServiceAvailabilityTier struct {
	Name      string                        `json:"name"`
	Location  *CloudPrivateEndpointLocation `json:"location,omitempty"`
	Transport string                        `json:"transport,omitempty"`
	Port      int64                         `json:"port,omitempty"`
	Bandwidth *int64                        `json:"bandwidth,omitempty"`
}

type CloudPrivateEndpointServiceAvailability struct {
	Name        string                                        `json:"name"`
	Description string                                        `json:"description,omitempty"`
	Tiers       []CloudPrivateEndpointServiceAvailabilityTier `json:"tiers"`
}

// Payloads

type CloudPrivateEndpointCreatePayload struct {
	TargetSpec *CloudPrivateEndpointTargetSpec `json:"targetSpec"`
}

type CloudPrivateEndpointUpdateTargetSpec struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CloudPrivateEndpointUpdatePayload struct {
	Checksum   string                                `json:"checksum"`
	TargetSpec *CloudPrivateEndpointUpdateTargetSpec `json:"targetSpec"`
}

// Attribute type helpers

func privateEndpointRefAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id": ovhtypes.TfStringType{},
	}
}

func privateEndpointLocationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"region": ovhtypes.TfStringType{},
	}
}

func privateEndpointServiceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": ovhtypes.TfStringType{},
		"tier": ovhtypes.TfStringType{},
	}
}

func privateEndpointCurrentStateServiceAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":      ovhtypes.TfStringType{},
		"tier":      ovhtypes.TfStringType{},
		"port":      types.Int64Type,
		"transport": ovhtypes.TfStringType{},
		"bandwidth": types.Int64Type,
	}
}

func privateEndpointCurrentStateAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":               ovhtypes.TfStringType{},
		"description":        ovhtypes.TfStringType{},
		"endpoint_id":        ovhtypes.TfStringType{},
		"location":           types.ObjectType{AttrTypes: privateEndpointLocationAttrTypes()},
		"service":            types.ObjectType{AttrTypes: privateEndpointCurrentStateServiceAttrTypes()},
		"network":            types.ObjectType{AttrTypes: privateEndpointRefAttrTypes()},
		"subnet":             types.ObjectType{AttrTypes: privateEndpointRefAttrTypes()},
		"vip_address":        ovhtypes.TfStringType{},
		"hostname":           ovhtypes.TfStringType{},
		"public_dns_enabled": types.BoolType,
	}
}

// ToCreate builds the POST payload. It carries every field, including the
// immutable placement, location and service references.
func (m *CloudPrivateEndpointModel) ToCreate() *CloudPrivateEndpointCreatePayload {
	target := &CloudPrivateEndpointTargetSpec{
		Name:             m.Name.ValueString(),
		Location:         &CloudPrivateEndpointLocation{Region: m.Region.ValueString()},
		PublicDnsEnabled: m.PublicDnsEnabled.ValueBool(),
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		target.Description = m.Description.ValueString()
	}

	if !m.Service.IsNull() && !m.Service.IsUnknown() {
		attrs := m.Service.Attributes()
		svc := &CloudPrivateEndpointServiceRef{}
		if v, ok := attrs["name"].(ovhtypes.TfStringValue); ok {
			svc.Name = v.ValueString()
		}
		if v, ok := attrs["tier"].(ovhtypes.TfStringValue); ok {
			svc.Tier = v.ValueString()
		}
		target.Service = svc
	}

	if !m.NetworkId.IsNull() && !m.NetworkId.IsUnknown() && m.NetworkId.ValueString() != "" {
		target.Network = &CloudPrivateEndpointRef{Id: m.NetworkId.ValueString()}
	}
	if !m.SubnetId.IsNull() && !m.SubnetId.IsUnknown() && m.SubnetId.ValueString() != "" {
		target.Subnet = &CloudPrivateEndpointRef{Id: m.SubnetId.ValueString()}
	}
	if !m.VipAddress.IsNull() && !m.VipAddress.IsUnknown() && m.VipAddress.ValueString() != "" {
		target.VipAddress = m.VipAddress.ValueString()
	}

	return &CloudPrivateEndpointCreatePayload{TargetSpec: target}
}

// ToUpdate builds the PUT payload. Only name and description are mutable, so
// no placement, location or service reference is included.
func (m *CloudPrivateEndpointModel) ToUpdate(checksum string) *CloudPrivateEndpointUpdatePayload {
	target := &CloudPrivateEndpointUpdateTargetSpec{
		Name: m.Name.ValueString(),
	}
	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		target.Description = m.Description.ValueString()
	}

	return &CloudPrivateEndpointUpdatePayload{
		Checksum:   checksum,
		TargetSpec: target,
	}
}

// resolvePrivateEndpointComputed fills an Optional+Computed attribute from the
// target spec, falling back to what the platform resolved in the current state.
// Such an attribute must never stay unknown, which would fail the apply on a
// missing value, so an unresolved one becomes null.
func resolvePrivateEndpointComputed(current ovhtypes.TfStringValue, spec, state string) ovhtypes.TfStringValue {
	if spec == "" {
		spec = state
	}
	if spec != "" {
		return ovhtypes.TfStringValue{StringValue: types.StringValue(spec)}
	}
	if current.IsUnknown() {
		return ovhtypes.TfStringValue{StringValue: types.StringNull()}
	}
	return current
}

func buildPrivateEndpointRefObject(ref *CloudPrivateEndpointRef) types.Object {
	if ref == nil {
		return types.ObjectNull(privateEndpointRefAttrTypes())
	}
	obj, _ := types.ObjectValue(privateEndpointRefAttrTypes(), map[string]attr.Value{
		"id": ovhtypes.TfStringValue{StringValue: types.StringValue(ref.Id)},
	})
	return obj
}

func buildPrivateEndpointCurrentState(cs *CloudPrivateEndpointCurrentState) types.Object {
	locObj := types.ObjectNull(privateEndpointLocationAttrTypes())
	if cs.Location != nil {
		locObj, _ = types.ObjectValue(privateEndpointLocationAttrTypes(), map[string]attr.Value{
			"region": ovhtypes.TfStringValue{StringValue: types.StringValue(cs.Location.Region)},
		})
	}

	svcObj := types.ObjectNull(privateEndpointCurrentStateServiceAttrTypes())
	if cs.Service != nil {
		svcObj, _ = types.ObjectValue(privateEndpointCurrentStateServiceAttrTypes(), map[string]attr.Value{
			"name":      ovhtypes.TfStringValue{StringValue: types.StringValue(cs.Service.Name)},
			"tier":      ovhtypes.TfStringValue{StringValue: types.StringValue(cs.Service.Tier)},
			"port":      types.Int64Value(cs.Service.Port),
			"transport": ovhtypes.TfStringValue{StringValue: types.StringValue(cs.Service.Transport)},
			"bandwidth": types.Int64PointerValue(cs.Service.Bandwidth),
		})
	}

	obj, _ := types.ObjectValue(privateEndpointCurrentStateAttrTypes(), map[string]attr.Value{
		"name":        ovhtypes.TfStringValue{StringValue: types.StringValue(cs.Name)},
		"description": ovhtypes.TfStringValue{StringValue: types.StringValue(cs.Description)},
		"endpoint_id": ovhtypes.TfStringValue{StringValue: types.StringValue(cs.EndpointId)},
		"location":    locObj,
		"service":     svcObj,
		"network":     buildPrivateEndpointRefObject(cs.Network),
		"subnet":      buildPrivateEndpointRefObject(cs.Subnet),
		"vip_address": ovhtypes.TfStringValue{StringValue: types.StringValue(cs.VipAddress)},
		// Null until the service publishes it.
		"hostname":           nullableTfString(cs.Hostname),
		"public_dns_enabled": types.BoolValue(cs.PublicDnsEnabled),
	})
	return obj
}

// MergeWith copies the API response into the Terraform model. Root fields come
// from targetSpec, except for the Optional+Computed ones that the platform may
// resolve on the caller's behalf, which fall back to currentState.
func (m *CloudPrivateEndpointModel) MergeWith(ctx context.Context, response *CloudPrivateEndpointAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)}
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)}

	if response.CurrentState != nil {
		m.CurrentState = buildPrivateEndpointCurrentState(response.CurrentState)
	} else {
		m.CurrentState = types.ObjectNull(privateEndpointCurrentStateAttrTypes())
	}

	ts, cs := response.TargetSpec, response.CurrentState
	if ts == nil {
		ts = &CloudPrivateEndpointTargetSpec{}
	}
	if cs == nil {
		cs = &CloudPrivateEndpointCurrentState{}
	}

	// vip_address, network_id and subnet_id are Optional+Computed: targetSpec
	// echoes back only what the caller supplied, so fall back to the address
	// allocated and the counterpart reference resolved by the platform, both
	// reported in currentState.
	m.VipAddress = resolvePrivateEndpointComputed(m.VipAddress, ts.VipAddress, cs.VipAddress)
	m.NetworkId = resolvePrivateEndpointComputed(m.NetworkId, ts.Network.id(), cs.Network.id())
	m.SubnetId = resolvePrivateEndpointComputed(m.SubnetId, ts.Subnet.id(), cs.Subnet.id())

	if response.TargetSpec == nil {
		return
	}

	m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(ts.Name)}

	// description is Optional and never defaulted by the server: preserve a null
	// config instead of forcing an empty string (avoids a phantom diff), while
	// still clearing a value the server dropped, so drift is detected.
	if ts.Description != "" || (!m.Description.IsNull() && !m.Description.IsUnknown()) {
		m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(ts.Description)}
	} else if m.Description.IsUnknown() {
		m.Description = ovhtypes.TfStringValue{StringValue: types.StringNull()}
	}

	if ts.Location != nil {
		m.Region = ovhtypes.TfStringValue{StringValue: types.StringValue(ts.Location.Region)}
	}

	if ts.Service != nil {
		svcObj, _ := types.ObjectValue(privateEndpointServiceAttrTypes(), map[string]attr.Value{
			"name": ovhtypes.TfStringValue{StringValue: types.StringValue(ts.Service.Name)},
			"tier": ovhtypes.TfStringValue{StringValue: types.StringValue(ts.Service.Tier)},
		})
		m.Service = svcObj
	}

	m.PublicDnsEnabled = types.BoolValue(ts.PublicDnsEnabled)
}
