package ovh

import (
	"context"

	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func handleCsrReplace(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
	var data, planData OkmsCredentialResourceModel

	resp.RequiresReplace = true
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.FromCsr.ValueBool() && (data.Csr.IsNull() || data.Csr.IsUnknown()) && planData.Csr != data.Csr {
		// This credential was created from a CSR but it's gone from the state.
		// This can happen if we remove and then reimport this ressource,
		// because the API doesn't return the original CSR.
		// In that case let's just update the state with the CSR present in the config,
		// there's no update to do on the server side.
		resp.RequiresReplace = false
	}
}

func OkmsCredentialResourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"certificate_pem": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Certificate PEM of the credential",
			MarkdownDescription: "Certificate PEM of the credential",
		},
		"created_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Creation time of the credential",
			MarkdownDescription: "Creation time of the credential",
		},
		"csr": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Computed:            true,
			Description:         "Valid Certificate Signing Request",
			MarkdownDescription: "Valid Certificate Signing Request",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIf(
					handleCsrReplace,
					"description",
					"markdown description",
				),
			},
		},
		"description": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Computed:            true,
			Description:         "Description of the credential (max 200)",
			MarkdownDescription: "Description of the credential (max 200)",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"expired_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Expiration time of the credential",
			MarkdownDescription: "Expiration time of the credential",
		},
		"from_csr": schema.BoolAttribute{
			CustomType:          ovhtypes.TfBoolType{},
			Computed:            true,
			Description:         "Is the credential generated from CSR",
			MarkdownDescription: "Is the credential generated from CSR",
		},
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "ID of the credential",
			MarkdownDescription: "ID of the credential",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"identity_urns": schema.ListAttribute{
			CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
			Required:            true,
			Description:         "List of identity URNs associated with the credential (max 25)",
			MarkdownDescription: "List of identity URNs associated with the credential (max 25)",
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Name of the credential (max 50)",
			MarkdownDescription: "Name of the credential (max 50)",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"okms_id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Okms ID",
			MarkdownDescription: "Okms ID",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"private_key_pem": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Sensitive:           true,
			Description:         "Private Key PEM of the credential if no CSR is provided (cannot be retrieve later)",
			MarkdownDescription: "Private Key PEM of the credential if no CSR is provided (cannot be retrieve later)",
		},
		"status": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Status of the credential",
			MarkdownDescription: "Status of the credential",
		},
		"validity": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Optional:            true,
			Computed:            true,
			Description:         "Validity in days (default 365, max 365)",
			MarkdownDescription: "Validity in days (default 365, max 365)",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.RequiresReplaceIfConfigured(),
			},
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type OkmsCredentialResourceModel struct {
	CertificatePem ovhtypes.TfStringValue                             `tfsdk:"certificate_pem" json:"certificatePem"`
	CreatedAt      ovhtypes.TfStringValue                             `tfsdk:"created_at" json:"createdAt"`
	Csr            ovhtypes.TfStringValue                             `tfsdk:"csr" json:"csr"`
	Description    ovhtypes.TfStringValue                             `tfsdk:"description" json:"description"`
	ExpiredAt      ovhtypes.TfStringValue                             `tfsdk:"expired_at" json:"expiredAt"`
	FromCsr        ovhtypes.TfBoolValue                               `tfsdk:"from_csr" json:"fromCsr"`
	Id             ovhtypes.TfStringValue                             `tfsdk:"id" json:"id"`
	IdentityUrns   ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"identity_urns" json:"identityURNs"`
	Name           ovhtypes.TfStringValue                             `tfsdk:"name" json:"name"`
	OkmsId         ovhtypes.TfStringValue                             `tfsdk:"okms_id" json:"okmsId"`
	PrivateKeyPem  ovhtypes.TfStringValue                             `tfsdk:"private_key_pem" json:"privateKeyPem"`
	Status         ovhtypes.TfStringValue                             `tfsdk:"status" json:"status"`
	Validity       ovhtypes.TfInt64Value                              `tfsdk:"validity" json:"validity"`
}

func (o *OkmsCredentialResourceModel) MergeWith(other *OkmsCredentialResourceModel) {

	if (o.CertificatePem.IsUnknown() || o.CertificatePem.IsNull()) && !other.CertificatePem.IsUnknown() {
		o.CertificatePem = other.CertificatePem
	}

	if (o.CreatedAt.IsUnknown() || o.CreatedAt.IsNull()) && !other.CreatedAt.IsUnknown() {
		o.CreatedAt = other.CreatedAt
	}

	if (o.Csr.IsUnknown() || o.Csr.IsNull()) && !other.Csr.IsUnknown() {
		o.Csr = other.Csr
	}

	if (o.Description.IsUnknown() || o.Description.IsNull()) && !other.Description.IsUnknown() {
		o.Description = other.Description
	}

	if (o.ExpiredAt.IsUnknown() || o.ExpiredAt.IsNull()) && !other.ExpiredAt.IsUnknown() {
		o.ExpiredAt = other.ExpiredAt
	}

	if (o.FromCsr.IsUnknown() || o.FromCsr.IsNull()) && !other.FromCsr.IsUnknown() {
		o.FromCsr = other.FromCsr
	}

	if (o.Id.IsUnknown() || o.Id.IsNull()) && !other.Id.IsUnknown() {
		o.Id = other.Id
	}

	if (o.IdentityUrns.IsUnknown() || o.IdentityUrns.IsNull()) && !other.IdentityUrns.IsUnknown() {
		o.IdentityUrns = other.IdentityUrns
	}

	if (o.Name.IsUnknown() || o.Name.IsNull()) && !other.Name.IsUnknown() {
		o.Name = other.Name
	}

	if (o.OkmsId.IsUnknown() || o.OkmsId.IsNull()) && !other.OkmsId.IsUnknown() {
		o.OkmsId = other.OkmsId
	}

	if (o.PrivateKeyPem.IsUnknown() || o.PrivateKeyPem.IsNull()) && !other.PrivateKeyPem.IsUnknown() {
		o.PrivateKeyPem = other.PrivateKeyPem
	}

	if (o.Status.IsUnknown() || o.Status.IsNull()) && !other.Status.IsUnknown() {
		o.Status = other.Status
	}

	if (o.Validity.IsUnknown() || o.Validity.IsNull()) && !other.Validity.IsUnknown() {
		o.Validity = other.Validity
	}

}

type OkmsCredentialWritableResourceModel struct {
	Csr          *ovhtypes.TfStringValue                             `tfsdk:"csr" json:"csr,omitempty"`
	Description  *ovhtypes.TfStringValue                             `tfsdk:"description" json:"description,omitempty"`
	IdentityUrns *ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"identity_urns" json:"identityURNs,omitempty"`
	Name         *ovhtypes.TfStringValue                             `tfsdk:"name" json:"name,omitempty"`
	Validity     *ovhtypes.TfInt64Value                              `tfsdk:"validity" json:"validity,omitempty"`
}

func (v OkmsCredentialResourceModel) ToCreate() *OkmsCredentialWritableResourceModel {
	res := &OkmsCredentialWritableResourceModel{}

	if !v.Csr.IsUnknown() {
		res.Csr = &v.Csr
	}

	if !v.Description.IsUnknown() {
		res.Description = &v.Description
	}

	if !v.IdentityUrns.IsUnknown() {
		res.IdentityUrns = &v.IdentityUrns
	}

	if !v.Name.IsUnknown() {
		res.Name = &v.Name
	}

	if !v.Validity.IsUnknown() {
		res.Validity = &v.Validity
	}

	return res
}

func (v OkmsCredentialResourceModel) ToUpdate() *OkmsCredentialWritableResourceModel {
	res := &OkmsCredentialWritableResourceModel{}

	if !v.Csr.IsUnknown() {
		res.Csr = &v.Csr
	}

	if !v.Description.IsUnknown() {
		res.Description = &v.Description
	}

	if !v.IdentityUrns.IsUnknown() {
		res.IdentityUrns = &v.IdentityUrns
	}

	if !v.Name.IsUnknown() {
		res.Name = &v.Name
	}

	if !v.Validity.IsUnknown() {
		res.Validity = &v.Validity
	}

	return res
}
