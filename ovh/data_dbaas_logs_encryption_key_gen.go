package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func DbaasLogsEncryptionKeyDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"created_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Encryption key creation date",
			MarkdownDescription: "Encryption key creation date",
		},
		"encryption_key_id": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{},
			Optional:   true,
			Computed:   true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("title")),
			},
			Description:         "Encryption key ID",
			MarkdownDescription: "Encryption key ID",
		},
		"fingerprint": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "PGP key fingerprint",
			MarkdownDescription: "PGP key fingerprint",
		},
		"is_editable": schema.BoolAttribute{
			CustomType:          ovhtypes.TfBoolType{},
			Computed:            true,
			Description:         "Indicates if the key is editable",
			MarkdownDescription: "Indicates if the key is editable",
		},
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name",
			MarkdownDescription: "Service name",
		},
		"title": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{},
			Optional:   true,
			Computed:   true,
			Validators: []validator.String{
				stringvalidator.ConflictsWith(path.MatchRoot("encryption_key_id")),
			},
			Description:         "Encryption key title",
			MarkdownDescription: "Encryption key title",
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type DbaasLogsEncryptionKeyDataSourceModel struct {
	CreatedAt       ovhtypes.TfStringValue `tfsdk:"created_at" json:"createdAt"`
	EncryptionKeyId ovhtypes.TfStringValue `tfsdk:"encryption_key_id" json:"encryptionKeyId"`
	Fingerprint     ovhtypes.TfStringValue `tfsdk:"fingerprint" json:"fingerprint"`
	IsEditable      ovhtypes.TfBoolValue   `tfsdk:"is_editable" json:"isEditable"`
	ServiceName     ovhtypes.TfStringValue `tfsdk:"service_name" json:"serviceName"`
	Title           ovhtypes.TfStringValue `tfsdk:"title" json:"title"`
}

func (v *DbaasLogsEncryptionKeyDataSourceModel) MergeWith(other *DbaasLogsEncryptionKeyDataSourceModel) {
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
