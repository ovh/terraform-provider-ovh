package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// cloudInstanceAPILocation is publicCloud.common.Location, shared by every
// /reference/instance catalog entry (flavor, image).
type cloudInstanceAPILocation struct {
	Region           string `json:"region"`
	AvailabilityZone string `json:"availabilityZone,omitempty"`
}

type CloudInstanceFlavorAPIResponse struct {
	Id          string                    `json:"id"`
	Name        string                    `json:"name"`
	Vcpus       int64                     `json:"vcpus"`
	Ram         int64                     `json:"ram"`
	Disk        int64                     `json:"disk"`
	Swap        int64                     `json:"swap"`
	Ephemeral   int64                     `json:"ephemeral"`
	IsPublic    bool                      `json:"isPublic"`
	Description string                    `json:"description"`
	Location    *cloudInstanceAPILocation `json:"location,omitempty"`
}

type CloudInstanceFlavorModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id          ovhtypes.TfStringValue `tfsdk:"id"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	Vcpus       types.Int64            `tfsdk:"vcpus"`
	Ram         types.Int64            `tfsdk:"ram"`
	Disk        types.Int64            `tfsdk:"disk"`
	Swap        types.Int64            `tfsdk:"swap"`
	Ephemeral   types.Int64            `tfsdk:"ephemeral"`
	IsPublic    types.Bool             `tfsdk:"is_public"`
	Description ovhtypes.TfStringValue `tfsdk:"description"`
	Location    types.Object           `tfsdk:"location"`
}

func instanceFlavorItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":          ovhtypes.TfStringType{},
		"name":        ovhtypes.TfStringType{},
		"vcpus":       types.Int64Type,
		"ram":         types.Int64Type,
		"disk":        types.Int64Type,
		"swap":        types.Int64Type,
		"ephemeral":   types.Int64Type,
		"is_public":   types.BoolType,
		"description": ovhtypes.TfStringType{},
		"location":    types.ObjectType{AttrTypes: instanceLocationAttrTypes()},
	}
}

func buildCloudInstanceLocationObject(l *cloudInstanceAPILocation) types.Object {
	if l == nil {
		return types.ObjectNull(instanceLocationAttrTypes())
	}
	obj, _ := types.ObjectValue(instanceLocationAttrTypes(), map[string]attr.Value{
		"region":            ovhtypes.TfStringValue{StringValue: types.StringValue(l.Region)},
		"availability_zone": ovhtypes.TfStringValue{StringValue: types.StringValue(l.AvailabilityZone)},
	})
	return obj
}

func (m *CloudInstanceFlavorModel) MergeWith(ctx context.Context, response *CloudInstanceFlavorAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Name)}
	m.Vcpus = types.Int64Value(response.Vcpus)
	m.Ram = types.Int64Value(response.Ram)
	m.Disk = types.Int64Value(response.Disk)
	m.Swap = types.Int64Value(response.Swap)
	m.Ephemeral = types.Int64Value(response.Ephemeral)
	m.IsPublic = types.BoolValue(response.IsPublic)
	m.Description = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Description)}
	m.Location = buildCloudInstanceLocationObject(response.Location)
}
