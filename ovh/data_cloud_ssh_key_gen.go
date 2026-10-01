package ovh

import (
	"context"

	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func CloudSshKeyDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Description:         "Service name. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
			MarkdownDescription: "Service name. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "SSH key name (unique per project, used as identifier)",
			MarkdownDescription: "SSH key name (unique per project, used as identifier)",
		},
		"public_key": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "SSH public key content",
			MarkdownDescription: "SSH public key content",
		},
		"created_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Creation date of the SSH key",
			MarkdownDescription: "Creation date of the SSH key",
		},
		"updated_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Last update date of the SSH key",
			MarkdownDescription: "Last update date of the SSH key",
		},
	}

	return schema.Schema{
		Description: "Get an SSH key of a Public Cloud project (API v2) by its name.",
		Attributes:  attrs,
	}
}

// CloudSshKeyDataSourceModel represents the Terraform state for a single SSH key
// datasource (API v2).
type CloudSshKeyDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name" json:"-"`
	Name        ovhtypes.TfStringValue `tfsdk:"name" json:"name"`
	PublicKey   ovhtypes.TfStringValue `tfsdk:"public_key" json:"publicKey"`
	CreatedAt   ovhtypes.TfStringValue `tfsdk:"created_at" json:"createdAt"`
	UpdatedAt   ovhtypes.TfStringValue `tfsdk:"updated_at" json:"updatedAt"`
}
