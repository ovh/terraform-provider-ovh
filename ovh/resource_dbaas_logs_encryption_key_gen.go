package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func DbaasLogsEncryptionKeyResourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Unique identifier for the resource",
		},
		"content": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Sensitive:           true,
			Description:         "PGP public key content",
			MarkdownDescription: "PGP public key content",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"created_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Encryption key creation date",
			MarkdownDescription: "Encryption key creation date",
		},
		"encryption_key_id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Encryption key ID",
			MarkdownDescription: "Encryption key ID",
		},
		"fingerprint": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "PGP key fingerprint",
			MarkdownDescription: "PGP key fingerprint",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"is_editable": schema.BoolAttribute{
			CustomType:          ovhtypes.TfBoolType{},
			Computed:            true,
			Description:         "Indicates if the key is editable",
			MarkdownDescription: "Indicates if the key is editable",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name",
			MarkdownDescription: "Service name",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"title": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Encryption key title",
			MarkdownDescription: "Encryption key title",
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type DbaasLogsEncryptionKeyModel struct {
	ID              ovhtypes.TfStringValue `tfsdk:"id" json:"-"`
	Content         ovhtypes.TfStringValue `tfsdk:"content" json:"content"`
	CreatedAt       ovhtypes.TfStringValue `tfsdk:"created_at" json:"createdAt"`
	EncryptionKeyId ovhtypes.TfStringValue `tfsdk:"encryption_key_id" json:"encryptionKeyId"`
	Fingerprint     ovhtypes.TfStringValue `tfsdk:"fingerprint" json:"fingerprint"`
	IsEditable      ovhtypes.TfBoolValue   `tfsdk:"is_editable" json:"isEditable"`
	ServiceName     ovhtypes.TfStringValue `tfsdk:"service_name" json:"serviceName"`
	Title           ovhtypes.TfStringValue `tfsdk:"title" json:"title"`
}

func (v *DbaasLogsEncryptionKeyModel) MergeWith(other *DbaasLogsEncryptionKeyModel) {
	if (v.ID.IsUnknown() || v.ID.IsNull()) && !other.ID.IsUnknown() {
		v.ID = other.ID
	}

	if (v.Content.IsUnknown() || v.Content.IsNull()) && !other.Content.IsUnknown() {
		v.Content = other.Content
	}

	if (v.CreatedAt.IsUnknown() || v.CreatedAt.IsNull()) && !other.CreatedAt.IsUnknown() {
		v.CreatedAt = other.CreatedAt
	}

	if (v.EncryptionKeyId.IsUnknown() || v.EncryptionKeyId.IsNull()) && !other.EncryptionKeyId.IsUnknown() {
		v.EncryptionKeyId = other.EncryptionKeyId
	}

	if (v.Fingerprint.IsUnknown() || v.Fingerprint.IsNull()) && !other.Fingerprint.IsUnknown() {
		v.Fingerprint = other.Fingerprint
	}

	if (v.IsEditable.IsUnknown() || v.IsEditable.IsNull()) && !other.IsEditable.IsUnknown() {
		v.IsEditable = other.IsEditable
	}

	if (v.ServiceName.IsUnknown() || v.ServiceName.IsNull()) && !other.ServiceName.IsUnknown() {
		v.ServiceName = other.ServiceName
	}

	if (v.Title.IsUnknown() || v.Title.IsNull()) && !other.Title.IsUnknown() {
		v.Title = other.Title
	}
}

type DbaasLogsEncryptionKeyWritableModel struct {
	Content     *ovhtypes.TfStringValue `tfsdk:"content" json:"content,omitempty"`
	Fingerprint *ovhtypes.TfStringValue `tfsdk:"fingerprint" json:"fingerprint,omitempty"`
	Title       *ovhtypes.TfStringValue `tfsdk:"title" json:"title,omitempty"`
}

func (v DbaasLogsEncryptionKeyModel) ToCreate() *DbaasLogsEncryptionKeyWritableModel {
	res := &DbaasLogsEncryptionKeyWritableModel{}

	if !v.Content.IsUnknown() {
		res.Content = &v.Content
	}

	if !v.Fingerprint.IsUnknown() {
		res.Fingerprint = &v.Fingerprint
	}

	if !v.Title.IsUnknown() {
		res.Title = &v.Title
	}

	return res
}

func (v DbaasLogsEncryptionKeyModel) ToUpdate() *DbaasLogsEncryptionKeyWritableModel {
	res := &DbaasLogsEncryptionKeyWritableModel{}

	if !v.Title.IsUnknown() {
		res.Title = &v.Title
	}

	return res
}
