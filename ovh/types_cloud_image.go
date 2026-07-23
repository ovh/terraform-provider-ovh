package ovh

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudImageAPIResponse struct {
	Id         string                  `json:"id"`
	Name       string                  `json:"name"`
	Status     string                  `json:"status"`
	Visibility string                  `json:"visibility"`
	MinDisk    int64                   `json:"minDisk"`
	MinRam     int64                   `json:"minRam"`
	Size       int64                   `json:"size"`
	CreatedAt  string                  `json:"createdAt"`
	UpdatedAt  string                  `json:"updatedAt"`
	Location   *CloudFlavorAPILocation `json:"location,omitempty"`
}

type CloudImageModel struct {
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
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
}

func imageAttrTypes() map[string]attr.Type {
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
		"region":     ovhtypes.TfStringType{},
	}
}

func buildImageModel(r *CloudImageAPIResponse) CloudImageModel {
	return CloudImageModel{
		Id:         ovhtypes.TfStringValue{StringValue: types.StringValue(r.Id)},
		Name:       ovhtypes.TfStringValue{StringValue: types.StringValue(r.Name)},
		Status:     ovhtypes.TfStringValue{StringValue: types.StringValue(r.Status)},
		Visibility: ovhtypes.TfStringValue{StringValue: types.StringValue(r.Visibility)},
		MinDisk:    types.Int64Value(r.MinDisk),
		MinRam:     types.Int64Value(r.MinRam),
		Size:       types.Int64Value(r.Size),
		CreatedAt:  ovhtypes.TfStringValue{StringValue: types.StringValue(r.CreatedAt)},
		UpdatedAt:  ovhtypes.TfStringValue{StringValue: types.StringValue(r.UpdatedAt)},
		Region:     flavorRegion(r.Location),
	}
}
