package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"

	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func OkmsServiceKeyJwkResourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"context": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Description:         "Context of the key",
			MarkdownDescription: "Context of the key",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplaceIfConfigured(),
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"created_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Creation time of the key",
			MarkdownDescription: "Creation time of the key",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"deactivation_reason": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Key deactivation reason",
			MarkdownDescription: "Key deactivation reason",
			Validators: []validator.String{
				stringvalidator.OneOf(
					"AFFILIATION_CHANGED",
					"CA_COMPROMISE",
					"CESSATION_OF_OPERATION",
					"KEY_COMPROMISE",
					"PRIVILEGE_WITHDRAWN",
					"SUPERSEDED",
					"UNSPECIFIED",
				),
			},
		},
		"iam": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{
				"display_name": schema.StringAttribute{
					CustomType:          ovhtypes.TfStringType{},
					Computed:            true,
					Description:         "Resource display name",
					MarkdownDescription: "Resource display name",
				},
				"id": schema.StringAttribute{
					CustomType:          ovhtypes.TfStringType{},
					Computed:            true,
					Description:         "Unique identifier of the resource",
					MarkdownDescription: "Unique identifier of the resource",
				},
				"tags": schema.MapAttribute{
					CustomType:          ovhtypes.NewTfMapNestedType[ovhtypes.TfStringValue](ctx),
					Computed:            true,
					Description:         "Resource tags. Tags that were internally computed are prefixed with ovh:",
					MarkdownDescription: "Resource tags. Tags that were internally computed are prefixed with ovh:",
				},
				"urn": schema.StringAttribute{
					CustomType:          ovhtypes.TfStringType{},
					Computed:            true,
					Description:         "Unique resource name used in policies",
					MarkdownDescription: "Unique resource name used in policies",
				},
			},
			CustomType: IamType{
				ObjectType: types.ObjectType{
					AttrTypes: IamValue{}.AttributeTypes(ctx),
				},
			},
			Computed:            true,
			Description:         "IAM resource metadata",
			MarkdownDescription: "IAM resource metadata",
		},
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Key ID",
			MarkdownDescription: "Key ID",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"keys": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"alg": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "The algorithm intended to be used with the key",
						MarkdownDescription: "The algorithm intended to be used with the key",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"ES256",
								"ES384",
								"ES512",
								"PS256",
								"PS384",
								"PS512",
								"RS256",
								"RS384",
								"RS512",
							),
						},
					},
					"crv": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "The cryptographic curve used with the key",
						MarkdownDescription: "The cryptographic curve used with the key",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"P-256",
								"P-384",
								"P-521",
							),
						},
					},
					"d": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "The RSA or EC private exponent",
						MarkdownDescription: "The RSA or EC private exponent",
						Sensitive:           true,
					},
					"dp": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Sensitive:           true,
						Description:         "The RSA private key's first factor CRT exponent",
						MarkdownDescription: "The RSA private key's first factor CRT exponent",
					},
					"dq": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Sensitive:           true,
						Description:         "The RSA private key's second factor CRT exponent",
						MarkdownDescription: "The RSA private key's second factor CRT exponent",
					},
					"e": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "The exponent value for the RSA public key",
						MarkdownDescription: "The exponent value for the RSA public key",
					},
					"k": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Sensitive:           true,
						Description:         "The value of the symmetric (or other single-valued) key",
						MarkdownDescription: "The value of the symmetric (or other single-valued) key",
					},
					"key_ops": schema.ListAttribute{
						CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
						Required:            true,
						Description:         "The operation for which the key is intended to be used",
						MarkdownDescription: "The operation for which the key is intended to be used",
					},
					"kid": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "key ID parameter used to match a specific key",
						MarkdownDescription: "key ID parameter used to match a specific key",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"kty": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Required:            true,
						Description:         "Key type parameter identifies the cryptographic algorithm family used with the key, such as RSA or EC",
						MarkdownDescription: "Key type parameter identifies the cryptographic algorithm family used with the key, such as RSA or EC",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"EC",
								"RSA",
								"oct",
							),
						},
					},
					"n": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "The modulus value for the RSA public key",
						MarkdownDescription: "The modulus value for the RSA public key",
					},
					"p": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Sensitive:           true,
						Description:         "The first prime factor of the RSA private key",
						MarkdownDescription: "The first prime factor of the RSA private key",
					},
					"q": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Sensitive:           true,
						Description:         "The second prime factor of the RSA private key",
						MarkdownDescription: "The second prime factor of the RSA private key",
					},
					"qi": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Sensitive:           true,
						Description:         "The CRT coefficient of the second factor of the RSA private key",
						MarkdownDescription: "The CRT coefficient of the second factor of the RSA private key",
					},
					"use": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "The intended use of the public key",
						MarkdownDescription: "The intended use of the public key",
						Validators: []validator.String{
							stringvalidator.OneOf(
								"enc",
								"sig",
							),
						},
					},
					"x": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "The x coordinate for the Elliptic Curve point",
						MarkdownDescription: "The x coordinate for the Elliptic Curve point",
					},
					"y": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "The y coordinate for the Elliptic Curve point",
						MarkdownDescription: "The y coordinate for the Elliptic Curve point",
					},
				},
				CustomType: JwkFullType{
					ObjectType: types.ObjectType{
						AttrTypes: JwkFullValue{}.AttributeTypes(ctx),
					},
				},
			},
			CustomType: ovhtypes.NewTfListNestedType[JwkFullValue](ctx),
			PlanModifiers: []planmodifier.List{
				listplanmodifier.RequiresReplace(),
			},
			Required:            true,
			Description:         "Set of JSON Web Keys to import",
			MarkdownDescription: "Set of JSON Web Keys to import",
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
				listvalidator.SizeAtLeast(1),
			},
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Key name",
			MarkdownDescription: "Key name",
		},
		"okms_id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Okms ID",
			MarkdownDescription: "Okms ID",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"size": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Size of the key to be created",
			MarkdownDescription: "Size of the key to be created",
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"state": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "State of the key",
			MarkdownDescription: "State of the key",
			Validators: []validator.String{
				stringvalidator.OneOf(
					"ACTIVE",
					"COMPROMISED",
					"DEACTIVATED",
				),
			},
		},
		"type": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Type of the key to be created",
			MarkdownDescription: "Type of the key to be created",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Validators: []validator.String{
				stringvalidator.OneOf(
					"EC",
					"RSA",
					"oct",
				),
			},
		},
	}

	return schema.Schema{
		Attributes:  attrs,
		Description: "Import an existing JWK in an OVHcloud KMS.",
	}
}

type OkmsServiceKeyJwkResourceModel struct {
	Context            ovhtypes.TfStringValue                   `tfsdk:"context" json:"context"`
	CreatedAt          ovhtypes.TfStringValue                   `tfsdk:"created_at" json:"createdAt"`
	DeactivationReason ovhtypes.TfStringValue                   `tfsdk:"deactivation_reason" json:"deactivationReason"`
	Iam                IamValue                                 `tfsdk:"iam" json:"iam"`
	Id                 ovhtypes.TfStringValue                   `tfsdk:"id" json:"id"`
	Keys               ovhtypes.TfListNestedValue[JwkFullValue] `tfsdk:"keys" json:"keys"`
	Name               ovhtypes.TfStringValue                   `tfsdk:"name" json:"name"`
	OkmsId             ovhtypes.TfStringValue                   `tfsdk:"okms_id" json:"okmsId"`
	Size               ovhtypes.TfInt64Value                    `tfsdk:"size" json:"size"`
	State              ovhtypes.TfStringValue                   `tfsdk:"state" json:"state"`
	Type               ovhtypes.TfStringValue                   `tfsdk:"type" json:"type"`
}

func (v *OkmsServiceKeyJwkResourceModel) MergeWith(other *OkmsServiceKeyJwkResourceModel) {

	if (v.Context.IsUnknown() || v.Context.IsNull()) && !other.Context.IsUnknown() {
		v.Context = other.Context
	}

	if (v.CreatedAt.IsUnknown() || v.CreatedAt.IsNull()) && !other.CreatedAt.IsUnknown() {
		v.CreatedAt = other.CreatedAt
	}

	if (v.DeactivationReason.IsUnknown() || v.DeactivationReason.IsNull()) && !other.DeactivationReason.IsUnknown() {
		v.DeactivationReason = other.DeactivationReason
	}

	if v.Iam.IsUnknown() && !other.Iam.IsUnknown() {
		v.Iam = other.Iam
	} else if !other.Iam.IsUnknown() {
		v.Iam.MergeWith(&other.Iam)
	}

	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}

	if (v.Keys.IsUnknown() || v.Keys.IsNull()) && !other.Keys.IsUnknown() {
		v.Keys = other.Keys
	} else if !other.Keys.IsUnknown() && !other.Keys.IsNull() {
		elems := v.Keys.Elements()
		newElems := other.Keys.Elements()

		if len(newElems) > 0 && len(elems) != len(newElems) {
			v.Keys = other.Keys
		} else {
			newSlice := make([]attr.Value, len(elems))
			for idx, e := range elems {
				tmp := e.(JwkFullValue)
				tmp2 := newElems[idx].(JwkFullValue)
				tmp.MergeWith(&tmp2)
				newSlice[idx] = tmp
			}

			v.Keys = ovhtypes.TfListNestedValue[JwkFullValue]{
				ListValue: basetypes.NewListValueMust(JwkFullValue{}.Type(context.Background()), newSlice),
			}
		}
	}

	if (v.Name.IsUnknown() || v.Name.IsNull()) && !other.Name.IsUnknown() {
		v.Name = other.Name
	}

	if (v.OkmsId.IsUnknown() || v.OkmsId.IsNull()) && !other.OkmsId.IsUnknown() {
		v.OkmsId = other.OkmsId
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

type OkmsServiceKeyJwkWritableModel struct {
	Context            *ovhtypes.TfStringValue                           `tfsdk:"context" json:"context,omitempty"`
	DeactivationReason *ovhtypes.TfStringValue                           `tfsdk:"deactivation_reason" json:"deactivationReason,omitempty"`
	Keys               *ovhtypes.TfListNestedValue[JwkFullWritableValue] `tfsdk:"keys" json:"keys,omitempty"`
	Name               *ovhtypes.TfStringValue                           `tfsdk:"name" json:"name,omitempty"`
	State              *ovhtypes.TfStringValue                           `tfsdk:"state" json:"state,omitempty"`
}

func (v OkmsServiceKeyJwkResourceModel) ToCreate(diag *diag.Diagnostics, ctx context.Context) *OkmsServiceKeyJwkWritableModel {
	res := &OkmsServiceKeyJwkWritableModel{}

	if !v.Context.IsUnknown() {
		res.Context = &v.Context
	}

	if !v.Keys.IsUnknown() {
		var createKeys []JwkFullWritableValue
		for _, elem := range v.Keys.Elements() {
			createKeys = append(createKeys, *elem.(JwkFullValue).ToCreate())
		}

		newKeys, d := basetypes.NewListValueFrom(context.Background(), JwkFullWritableValue{
			JwkFullValue: &JwkFullValue{},
		}.Type(context.Background()), createKeys)
		diag.Append(d...)
		res.Keys = &ovhtypes.TfListNestedValue[JwkFullWritableValue]{
			ListValue: newKeys,
		}
	}

	if !v.Name.IsUnknown() {
		res.Name = &v.Name
	}

	return res
}

func (v OkmsServiceKeyJwkResourceModel) ToUpdate() *OkmsServiceKeyJwkWritableModel {
	res := &OkmsServiceKeyJwkWritableModel{}

	if !v.DeactivationReason.IsUnknown() {
		res.DeactivationReason = &v.DeactivationReason
	}

	if !v.Name.IsUnknown() {
		res.Name = &v.Name
	}

	if !v.State.IsUnknown() {
		res.State = &v.State
	}

	return res
}

func NewJwkFullValueNull() JwkFullValue {
	return JwkFullValue{
		state: attr.ValueStateNull,
	}
}

func NewJwkFullValueUnknown() JwkFullValue {
	return JwkFullValue{
		state: attr.ValueStateUnknown,
	}
}

func NewJwkFullValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (JwkFullValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing JwkFullValue Attribute Value",
				"While creating a JwkFullValue value, a missing attribute value was detected. "+
					"A JwkFullValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("JwkFullValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid JwkFullValue Attribute Type",
				"While creating a JwkFullValue value, an invalid attribute value was detected. "+
					"A JwkFullValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("JwkFullValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("JwkFullValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra JwkFullValue Attribute Value",
				"While creating a JwkFullValue value, an extra attribute value was detected. "+
					"A JwkFullValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra JwkFullValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewJwkFullValueUnknown(), diags
	}

	algAttribute, ok := attributes["alg"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`alg is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	algVal, ok := algAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`alg expected to be ovhtypes.TfStringValue, was: %T`, algAttribute))
	}

	crvAttribute, ok := attributes["crv"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`crv is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	crvVal, ok := crvAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`crv expected to be ovhtypes.TfStringValue, was: %T`, crvAttribute))
	}

	dAttribute, ok := attributes["d"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`d is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	dVal, ok := dAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`d expected to be ovhtypes.TfStringValue, was: %T`, dAttribute))
	}

	dpAttribute, ok := attributes["dp"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dp is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	dpVal, ok := dpAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dp expected to be ovhtypes.TfStringValue, was: %T`, dpAttribute))
	}

	dqAttribute, ok := attributes["dq"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dq is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	dqVal, ok := dqAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dq expected to be ovhtypes.TfStringValue, was: %T`, dqAttribute))
	}

	eAttribute, ok := attributes["e"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`e is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	eVal, ok := eAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`e expected to be ovhtypes.TfStringValue, was: %T`, eAttribute))
	}

	kAttribute, ok := attributes["k"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`k is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	kVal, ok := kAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`k expected to be ovhtypes.TfStringValue, was: %T`, kAttribute))
	}

	keyOpsAttribute, ok := attributes["key_ops"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`key_ops is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	keyOpsVal, ok := keyOpsAttribute.(ovhtypes.TfListNestedValue[ovhtypes.TfStringValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`key_ops expected to be ovhtypes.TfListNestedValue[ovhtypes.TfStringValue], was: %T`, keyOpsAttribute))
	}

	kidAttribute, ok := attributes["kid"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`kid is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	kidVal, ok := kidAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`kid expected to be ovhtypes.TfStringValue, was: %T`, kidAttribute))
	}

	ktyAttribute, ok := attributes["kty"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`kty is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	ktyVal, ok := ktyAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`kty expected to be ovhtypes.TfStringValue, was: %T`, ktyAttribute))
	}

	nAttribute, ok := attributes["n"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`n is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	nVal, ok := nAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`n expected to be ovhtypes.TfStringValue, was: %T`, nAttribute))
	}

	pAttribute, ok := attributes["p"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`p is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	pVal, ok := pAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`p expected to be ovhtypes.TfStringValue, was: %T`, pAttribute))
	}

	qAttribute, ok := attributes["q"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`q is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	qVal, ok := qAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`q expected to be ovhtypes.TfStringValue, was: %T`, qAttribute))
	}

	qiAttribute, ok := attributes["qi"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`qi is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	qiVal, ok := qiAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`qi expected to be ovhtypes.TfStringValue, was: %T`, qiAttribute))
	}

	useAttribute, ok := attributes["use"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`use is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	useVal, ok := useAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`use expected to be ovhtypes.TfStringValue, was: %T`, useAttribute))
	}

	xAttribute, ok := attributes["x"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`x is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	xVal, ok := xAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`x expected to be ovhtypes.TfStringValue, was: %T`, xAttribute))
	}

	yAttribute, ok := attributes["y"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`y is missing from object`)

		return NewJwkFullValueUnknown(), diags
	}

	yVal, ok := yAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`y expected to be ovhtypes.TfStringValue, was: %T`, yAttribute))
	}

	if diags.HasError() {
		return NewJwkFullValueUnknown(), diags
	}

	return JwkFullValue{
		Alg:    algVal,
		Crv:    crvVal,
		D:      dVal,
		Dp:     dpVal,
		Dq:     dqVal,
		E:      eVal,
		K:      kVal,
		KeyOps: keyOpsVal,
		Kid:    kidVal,
		Kty:    ktyVal,
		N:      nVal,
		P:      pVal,
		Q:      qVal,
		Qi:     qiVal,
		Use:    useVal,
		X:      xVal,
		Y:      yVal,
		state:  attr.ValueStateKnown,
	}, diags
}

func NewJwkFullValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) JwkFullValue {
	object, diags := NewJwkFullValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewJwkFullValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func NewJwkFullWritableValueNull() JwkFullWritableValue {
	return JwkFullWritableValue{
		state: attr.ValueStateNull,
	}
}

func NewJwkFullWritableValueUnknown() JwkFullWritableValue {
	return JwkFullWritableValue{
		state: attr.ValueStateUnknown,
	}
}

func NewJwkFullWritableValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (JwkFullWritableValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing JwkFullWritableValue Attribute Value",
				"While creating a JwkFullWritableValue value, a missing attribute value was detected. "+
					"A JwkFullWritableValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("JwkFullWritableValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid JwkFullWritableValue Attribute Type",
				"While creating a JwkFullWritableValue value, an invalid attribute value was detected. "+
					"A JwkFullWritableValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("JwkFullWritableValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("JwkFullWritableValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra JwkFullWritableValue Attribute Value",
				"While creating a JwkFullWritableValue value, an extra attribute value was detected. "+
					"A JwkFullWritableValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra JwkFullWritableValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewJwkFullWritableValueUnknown(), diags
	}

	algAttribute, ok := attributes["alg"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`alg is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	algVal, ok := algAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`alg expected to be ovhtypes.TfStringValue, was: %T`, algAttribute))
	}

	crvAttribute, ok := attributes["crv"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`crv is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	crvVal, ok := crvAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`crv expected to be ovhtypes.TfStringValue, was: %T`, crvAttribute))
	}

	dAttribute, ok := attributes["d"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`d is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	dVal, ok := dAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`d expected to be ovhtypes.TfStringValue, was: %T`, dAttribute))
	}

	dpAttribute, ok := attributes["dp"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dp is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	dpVal, ok := dpAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dp expected to be ovhtypes.TfStringValue, was: %T`, dpAttribute))
	}

	dqAttribute, ok := attributes["dq"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`dq is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	dqVal, ok := dqAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`dq expected to be ovhtypes.TfStringValue, was: %T`, dqAttribute))
	}

	eAttribute, ok := attributes["e"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`e is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	eVal, ok := eAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`e expected to be ovhtypes.TfStringValue, was: %T`, eAttribute))
	}

	kAttribute, ok := attributes["k"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`k is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	kVal, ok := kAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`k expected to be ovhtypes.TfStringValue, was: %T`, kAttribute))
	}

	keyOpsAttribute, ok := attributes["key_ops"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`key_ops is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	keyOpsVal, ok := keyOpsAttribute.(ovhtypes.TfListNestedValue[ovhtypes.TfStringValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`key_ops expected to be ovhtypes.TfListNestedValue[ovhtypes.TfStringValue], was: %T`, keyOpsAttribute))
	}

	kidAttribute, ok := attributes["kid"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`kid is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	kidVal, ok := kidAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`kid expected to be ovhtypes.TfStringValue, was: %T`, kidAttribute))
	}

	ktyAttribute, ok := attributes["kty"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`kty is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	ktyVal, ok := ktyAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`kty expected to be ovhtypes.TfStringValue, was: %T`, ktyAttribute))
	}

	nAttribute, ok := attributes["n"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`n is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	nVal, ok := nAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`n expected to be ovhtypes.TfStringValue, was: %T`, nAttribute))
	}

	pAttribute, ok := attributes["p"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`p is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	pVal, ok := pAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`p expected to be ovhtypes.TfStringValue, was: %T`, pAttribute))
	}

	qAttribute, ok := attributes["q"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`q is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	qVal, ok := qAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`q expected to be ovhtypes.TfStringValue, was: %T`, qAttribute))
	}

	qiAttribute, ok := attributes["qi"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`qi is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	qiVal, ok := qiAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`qi expected to be ovhtypes.TfStringValue, was: %T`, qiAttribute))
	}

	useAttribute, ok := attributes["use"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`use is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	useVal, ok := useAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`use expected to be ovhtypes.TfStringValue, was: %T`, useAttribute))
	}

	xAttribute, ok := attributes["x"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`x is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	xVal, ok := xAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`x expected to be ovhtypes.TfStringValue, was: %T`, xAttribute))
	}

	yAttribute, ok := attributes["y"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`y is missing from object`)

		return NewJwkFullWritableValueUnknown(), diags
	}

	yVal, ok := yAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`y expected to be ovhtypes.TfStringValue, was: %T`, yAttribute))
	}

	if diags.HasError() {
		return NewJwkFullWritableValueUnknown(), diags
	}

	return JwkFullWritableValue{
		Alg:    &algVal,
		Crv:    &crvVal,
		D:      &dVal,
		Dp:     &dpVal,
		Dq:     &dqVal,
		E:      &eVal,
		K:      &kVal,
		KeyOps: &keyOpsVal,
		Kid:    &kidVal,
		Kty:    &ktyVal,
		N:      &nVal,
		P:      &pVal,
		Q:      &qVal,
		Qi:     &qiVal,
		Use:    &useVal,
		X:      &xVal,
		Y:      &yVal,
		state:  attr.ValueStateKnown,
	}, diags
}

func NewJwkFullWritableValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) JwkFullWritableValue {
	object, diags := NewJwkFullWritableValue(attributeTypes, attributes)

	if diags.HasError() {
		// This could potentially be added to the diag package.
		diagsStrings := make([]string, 0, len(diags))

		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}

		panic("NewJwkFullValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}
