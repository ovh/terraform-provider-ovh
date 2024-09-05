package ovh

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"
)

func OkmsServiceKeyAttributes(ctx context.Context) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"created_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Creation time of the key",
			MarkdownDescription: "Creation time of the key",
		},
		"curve": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Curve type for Elliptic Curve (EC) keys",
			MarkdownDescription: "Curve type for Elliptic Curve (EC) keys",
		},
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Key ID",
			MarkdownDescription: "Key ID",
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Key name",
			MarkdownDescription: "Key name",
		},
		"operations": schema.ListAttribute{
			CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
			Computed:            true,
			Description:         "The operations for which the key is intended to be used",
			MarkdownDescription: "The operations for which the key is intended to be used",
		},
		"size": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Size of the key",
			MarkdownDescription: "Size of the key",
		},
		"state": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "State of the key",
			MarkdownDescription: "State of the key",
		},
		"type": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Key type",
			MarkdownDescription: "Key type",
		},
	}
}

func OkmsServiceKeyDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := OkmsServiceKeyAttributes(ctx)
	attrs["okms_id"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Required:            true,
		Description:         "Okms ID",
		MarkdownDescription: "Okms ID",
	}
	return schema.Schema{
		Attributes:  attrs,
		Description: "Use this data source to retrieve information about a KMS service key.",
	}
}

type OkmsServiceKeyModel struct {
	CreatedAt  ovhtypes.TfStringValue                             `tfsdk:"created_at" json:"createdAt"`
	Curve      ovhtypes.TfStringValue                             `tfsdk:"curve" json:"curve"`
	Id         ovhtypes.TfStringValue                             `tfsdk:"id" json:"id"`
	Name       ovhtypes.TfStringValue                             `tfsdk:"name" json:"name"`
	OkmsId     ovhtypes.TfStringValue                             `tfsdk:"okms_id" json:"okmsId"`
	Operations ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"operations" json:"operations"`
	Size       ovhtypes.TfInt64Value                              `tfsdk:"size" json:"size"`
	State      ovhtypes.TfStringValue                             `tfsdk:"state" json:"state"`
	Type       ovhtypes.TfStringValue                             `tfsdk:"type" json:"type"`
}

func (v *OkmsServiceKeyModel) MergeWith(other *OkmsServiceKeyModel) {

	if (v.CreatedAt.IsUnknown() || v.CreatedAt.IsNull()) && !other.CreatedAt.IsUnknown() {
		v.CreatedAt = other.CreatedAt
	}

	if (v.Curve.IsUnknown() || v.Curve.IsNull()) && !other.Curve.IsUnknown() {
		v.Curve = other.Curve
	}

	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}

	if (v.Name.IsUnknown() || v.Name.IsNull()) && !other.Name.IsUnknown() {
		v.Name = other.Name
	}

	if (v.OkmsId.IsUnknown() || v.OkmsId.IsNull()) && !other.OkmsId.IsUnknown() {
		v.OkmsId = other.OkmsId
	}

	if (v.Operations.IsUnknown() || v.Operations.IsNull()) && !other.Operations.IsUnknown() {
		v.Operations = other.Operations
	}

	if (v.Size.IsUnknown() || v.Size.IsNull()) && !other.Size.IsUnknown() {
		v.Size = other.Size
	}

	if (v.State.IsUnknown() || v.State.IsNull()) && !other.State.IsUnknown() {
		v.State = other.State
	}

	if (v.Type.IsUnknown() || v.Type.IsNull()) && !other.Type.IsUnknown() {
		v.Type = other.Type
	}
}
