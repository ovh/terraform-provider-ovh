package ovh

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudFlavorAPILocation struct {
	Region string `json:"region"`
}

type CloudFlavorAPIResponse struct {
	Id          string                  `json:"id"`
	Name        string                  `json:"name"`
	Vcpus       int64                   `json:"vcpus"`
	Ram         int64                   `json:"ram"`
	Disk        int64                   `json:"disk"`
	Swap        int64                   `json:"swap"`
	Ephemeral   int64                   `json:"ephemeral"`
	IsPublic    bool                    `json:"isPublic"`
	Description string                  `json:"description"`
	Location    *CloudFlavorAPILocation `json:"location,omitempty"`
}

type CloudFlavorModel struct {
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
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
}

func flavorAttrTypes() map[string]attr.Type {
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
		"region":      ovhtypes.TfStringType{},
	}
}

func flavorRegion(r *CloudFlavorAPILocation) ovhtypes.TfStringValue {
	if r == nil {
		return ovhtypes.TfStringValue{StringValue: types.StringNull()}
	}
	return ovhtypes.TfStringValue{StringValue: types.StringValue(r.Region)}
}

func buildFlavorModel(r *CloudFlavorAPIResponse) CloudFlavorModel {
	return CloudFlavorModel{
		Id:          ovhtypes.TfStringValue{StringValue: types.StringValue(r.Id)},
		Name:        ovhtypes.TfStringValue{StringValue: types.StringValue(r.Name)},
		Vcpus:       types.Int64Value(r.Vcpus),
		Ram:         types.Int64Value(r.Ram),
		Disk:        types.Int64Value(r.Disk),
		Swap:        types.Int64Value(r.Swap),
		Ephemeral:   types.Int64Value(r.Ephemeral),
		IsPublic:    types.BoolValue(r.IsPublic),
		Description: ovhtypes.TfStringValue{StringValue: types.StringValue(r.Description)},
		Region:      flavorRegion(r.Location),
	}
}
