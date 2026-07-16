package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudInstanceImageAPIResponse struct {
	Id         string                    `json:"id"`
	Name       string                    `json:"name"`
	Status     string                    `json:"status"`
	Visibility string                    `json:"visibility"`
	MinDisk    int64                     `json:"minDisk"`
	MinRam     int64                     `json:"minRam"`
	Size       int64                     `json:"size"`
	CreatedAt  string                    `json:"createdAt"`
	UpdatedAt  string                    `json:"updatedAt"`
	Location   *cloudInstanceAPILocation `json:"location,omitempty"`
}

type CloudInstanceImageModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id          ovhtypes.TfStringValue `tfsdk:"id"`
	Name        ovhtypes.TfStringValue `tfsdk:"name"`
	Status      ovhtypes.TfStringValue `tfsdk:"status"`
	Visibility  ovhtypes.TfStringValue `tfsdk:"visibility"`
	MinDisk     types.Int64            `tfsdk:"min_disk"`
	MinRam      types.Int64            `tfsdk:"min_ram"`
	Size        types.Int64            `tfsdk:"size"`
	CreatedAt   ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt   ovhtypes.TfStringValue `tfsdk:"updated_at"`
	Location    types.Object           `tfsdk:"location"`
}

func instanceImageItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         ovhtypes.TfStringType{},
		"name":       ovhtypes.TfStringType{},
		"status":     ovhtypes.TfStringType{},
		"visibility": ovhtypes.TfStringType{},
		"min_disk":   types.Int64Type,
		"min_ram":    types.Int64Type,
		"size":       types.Int64Type,
		"created_at": ovhtypes.TfStringType{},
		"updated_at": ovhtypes.TfStringType{},
		"location":   types.ObjectType{AttrTypes: instanceLocationAttrTypes()},
	}
}

func (m *CloudInstanceImageModel) MergeWith(ctx context.Context, response *CloudInstanceImageAPIResponse) {
	m.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)}
	m.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Name)}
	m.Status = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Status)}
	m.Visibility = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Visibility)}
	m.MinDisk = types.Int64Value(response.MinDisk)
	m.MinRam = types.Int64Value(response.MinRam)
	m.Size = types.Int64Value(response.Size)
	m.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)}
	m.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)}
	m.Location = buildCloudInstanceLocationObject(response.Location)
}
