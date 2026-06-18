package ovh

import (
	"context"

	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func CloudSshKeyResourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"name": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{},
			Required:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Description:         "SSH key name (unique per project, used as identifier)",
			MarkdownDescription: "SSH key name (unique per project, used as identifier)",
		},
		"public_key": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{},
			Required:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
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
		"service_name": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{},
			Optional:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Description:         "Service name. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
			MarkdownDescription: "Service name. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
		},
	}

	return schema.Schema{
		Description: "Creates an SSH key in a Public Cloud project (API v2). Keys are stored in the database and synced to OpenStack lazily on instance creation.",
		Attributes:  attrs,
	}
}

// CloudSshKeyModel represents the Terraform state for an SSH key (API v2).
type CloudSshKeyModel struct {
	Name        ovhtypes.TfStringValue `tfsdk:"name" json:"-"`
	PublicKey   ovhtypes.TfStringValue `tfsdk:"public_key" json:"-"`
	CreatedAt   ovhtypes.TfStringValue `tfsdk:"created_at" json:"-"`
	UpdatedAt   ovhtypes.TfStringValue `tfsdk:"updated_at" json:"-"`
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name" json:"-"`
}

func (v *CloudSshKeyModel) MergeWith(other *CloudSshKeyModel) {
	if (v.Name.IsUnknown() || v.Name.IsNull()) && !other.Name.IsUnknown() {
		v.Name = other.Name
	}
	if (v.PublicKey.IsUnknown() || v.PublicKey.IsNull()) && !other.PublicKey.IsUnknown() {
		v.PublicKey = other.PublicKey
	}
	if (v.CreatedAt.IsUnknown() || v.CreatedAt.IsNull()) && !other.CreatedAt.IsUnknown() {
		v.CreatedAt = other.CreatedAt
	}
	if (v.UpdatedAt.IsUnknown() || v.UpdatedAt.IsNull()) && !other.UpdatedAt.IsUnknown() {
		v.UpdatedAt = other.UpdatedAt
	}
	if (v.ServiceName.IsUnknown() || v.ServiceName.IsNull()) && !other.ServiceName.IsUnknown() {
		v.ServiceName = other.ServiceName
	}
}

// cloudSshKeyCreateRequest is the POST body for creating an SSH key via API v2.
type cloudSshKeyCreateRequest struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
}

// cloudSshKeyResponse is the flat response returned by API v2.
type cloudSshKeyResponse struct {
	Name      string `json:"name"`
	PublicKey string `json:"publicKey"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// toModel converts an API v2 response to a Terraform model.
func (r *cloudSshKeyResponse) toModel() *CloudSshKeyModel {
	return &CloudSshKeyModel{
		Name:      ovhtypes.NewTfStringValue(r.Name),
		PublicKey: ovhtypes.NewTfStringValue(r.PublicKey),
		CreatedAt: ovhtypes.NewTfStringValue(r.CreatedAt),
		UpdatedAt: ovhtypes.NewTfStringValue(r.UpdatedAt),
	}
}
