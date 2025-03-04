// Code generated by terraform-plugin-framework-generator DO NOT EDIT.

package ovh

import (
	"bytes"
	"context"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func IploadbalancingSslResourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"certificate": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{},
			Required:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIf(
					checkIfFingerprintChange,
					"Check if fingerprint change",
					"Check if fingerprint change",
				),
			},
			Description:         "Certificate",
			MarkdownDescription: "Certificate",
		},
		"chain": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{},
			Optional:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIf(
					checkIfFingerprintChange,
					"Check if fingerprint change",
					"Check if fingerprint change",
				),
			},
			Description:         "Certificate chain",
			MarkdownDescription: "Certificate chain",
		},
		"display_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Description:         "Human readable name for your ssl certificate, this field is for you",
			MarkdownDescription: "Human readable name for your ssl certificate, this field is for you",
		},
		"expire_date": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Expire date of your SSL certificate",
			MarkdownDescription: "Expire date of your SSL certificate",
		},
		"fingerprint": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Fingerprint of your SSL certificate",
			MarkdownDescription: "Fingerprint of your SSL certificate",
		},
		"id": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Id of your SSL certificate",
			MarkdownDescription: "Id of your SSL certificate",
		},
		"key": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{},
			Required:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIf(
					checkIfFingerprintChange,
					"Check if fingerprint change",
					"Check if fingerprint change",
				),
			},
			Description:         "Certificate key",
			MarkdownDescription: "Certificate key",
			Sensitive:           true,
		},
		"san": schema.ListAttribute{
			CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
			Computed:            true,
			Description:         "Subject Alternative Name of your SSL certificate",
			MarkdownDescription: "Subject Alternative Name of your SSL certificate",
		},
		"serial": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Serial of your SSL certificate (Deprecated, use fingerprint instead!)",
			MarkdownDescription: "Serial of your SSL certificate (Deprecated, use fingerprint instead!)",
		},
		"service_name": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{},
			Required:   true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
			Description:         "The internal name of your IP load balancing",
			MarkdownDescription: "The internal name of your IP load balancing",
		},
		"subject": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Subject of your SSL certificate",
			MarkdownDescription: "Subject of your SSL certificate",
		},
		"type": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Type of your SSL certificate.\n'built' for SSL certificates managed by the IP Load Balancing. 'custom' for user manager certificates.",
			MarkdownDescription: "Type of your SSL certificate.\n'built' for SSL certificates managed by the IP Load Balancing. 'custom' for user manager certificates.",
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type IploadbalancingSslModel struct {
	Certificate ovhtypes.TfStringValue                             `tfsdk:"certificate" json:"certificate"`
	Chain       ovhtypes.TfStringValue                             `tfsdk:"chain" json:"chain"`
	DisplayName ovhtypes.TfStringValue                             `tfsdk:"display_name" json:"displayName"`
	ExpireDate  ovhtypes.TfStringValue                             `tfsdk:"expire_date" json:"expireDate"`
	Fingerprint ovhtypes.TfStringValue                             `tfsdk:"fingerprint" json:"fingerprint"`
	Id          ovhtypes.TfInt64Value                              `tfsdk:"id" json:"id"`
	Key         ovhtypes.TfStringValue                             `tfsdk:"key" json:"key"`
	San         ovhtypes.TfListNestedValue[ovhtypes.TfStringValue] `tfsdk:"san" json:"san"`
	Serial      ovhtypes.TfStringValue                             `tfsdk:"serial" json:"serial"`
	ServiceName ovhtypes.TfStringValue                             `tfsdk:"service_name" json:"serviceName"`
	Subject     ovhtypes.TfStringValue                             `tfsdk:"subject" json:"subject"`
	Type        ovhtypes.TfStringValue                             `tfsdk:"type" json:"type"`
}

func (v *IploadbalancingSslModel) MergeWith(other *IploadbalancingSslModel) {

	if (v.Certificate.IsUnknown() || v.Certificate.IsNull()) && !other.Certificate.IsUnknown() {
		v.Certificate = other.Certificate
	}

	if (v.Chain.IsUnknown() || v.Chain.IsNull()) && !other.Chain.IsUnknown() {
		v.Chain = other.Chain
	}

	if (v.DisplayName.IsUnknown() || v.DisplayName.IsNull()) && !other.DisplayName.IsUnknown() {
		v.DisplayName = other.DisplayName
	}

	if (v.ExpireDate.IsUnknown() || v.ExpireDate.IsNull()) && !other.ExpireDate.IsUnknown() {
		v.ExpireDate = other.ExpireDate
	}

	if (v.Fingerprint.IsUnknown() || v.Fingerprint.IsNull()) && !other.Fingerprint.IsUnknown() {
		v.Fingerprint = other.Fingerprint
	}

	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}

	if (v.Key.IsUnknown() || v.Key.IsNull()) && !other.Key.IsUnknown() {
		v.Key = other.Key
	}

	if (v.San.IsUnknown() || v.San.IsNull()) && !other.San.IsUnknown() {
		v.San = other.San
	}

	if (v.Serial.IsUnknown() || v.Serial.IsNull()) && !other.Serial.IsUnknown() {
		v.Serial = other.Serial
	}

	if (v.ServiceName.IsUnknown() || v.ServiceName.IsNull()) && !other.ServiceName.IsUnknown() {
		v.ServiceName = other.ServiceName
	}

	if (v.Subject.IsUnknown() || v.Subject.IsNull()) && !other.Subject.IsUnknown() {
		v.Subject = other.Subject
	}

	if (v.Type.IsUnknown() || v.Type.IsNull()) && !other.Type.IsUnknown() {
		v.Type = other.Type
	}

}

type IploadbalancingSslWritableModel struct {
	Certificate *ovhtypes.TfStringValue `tfsdk:"certificate" json:"certificate,omitempty"`
	Chain       *ovhtypes.TfStringValue `tfsdk:"chain" json:"chain,omitempty"`
	DisplayName *ovhtypes.TfStringValue `tfsdk:"display_name" json:"displayName,omitempty"`
	Key         *ovhtypes.TfStringValue `tfsdk:"key" json:"key,omitempty"`
}

func (v IploadbalancingSslModel) ToCreate() *IploadbalancingSslWritableModel {
	res := &IploadbalancingSslWritableModel{}

	if !v.Certificate.IsUnknown() {
		res.Certificate = &v.Certificate
	}

	if !v.Chain.IsUnknown() {
		res.Chain = &v.Chain
	}

	if !v.DisplayName.IsUnknown() {
		res.DisplayName = &v.DisplayName
	}

	if !v.Key.IsUnknown() {
		res.Key = &v.Key
	}

	return res
}

func (v IploadbalancingSslModel) ToUpdate() *IploadbalancingSslWritableModel {
	res := &IploadbalancingSslWritableModel{}

	if !v.DisplayName.IsUnknown() {
		res.DisplayName = &v.DisplayName
	}

	return res
}

func checkIfFingerprintChange(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
	var data, planData IploadbalancingSslModel

	resp.RequiresReplace = true
	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// compute fingerprint from certificate
	block, _ := pem.Decode([]byte(planData.Certificate.ValueString()))
	if block == nil {
		tflog.Error(ctx, "Failed to parse pem file")
		return
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		tflog.Error(ctx, err.Error())
		return
	}

	fingerprint := sha1.Sum(cert.Raw)

	var buf bytes.Buffer
	for i, f := range fingerprint {
		if i > 0 {
			fmt.Fprintf(&buf, ":")
		}
		fmt.Fprintf(&buf, "%02X", f)
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if buf.String() == data.Fingerprint.ValueString() {
		// This ssl was created from a CERT/KEY/CHAIN but it's gone from the state.
		// This happens when we import this resource,
		// because the API doesn't return the original CERT/KEY/CHAIN.
		// In that case let's just update the state with the CERT/KEY/CHAIN present in the config,
		// there's no update to do on the server side.
		resp.RequiresReplace = false
	}
}
