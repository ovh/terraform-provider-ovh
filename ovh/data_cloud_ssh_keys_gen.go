package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func CloudSshKeysDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"ssh_keys": schema.SetNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "SSH key name",
						MarkdownDescription: "SSH key name",
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
				},
				CustomType: CloudSshKeysType{
					ObjectType: types.ObjectType{
						AttrTypes: CloudSshKeysValue{}.AttributeTypes(ctx),
					},
				},
			},
			CustomType:          ovhtypes.NewTfListNestedType[CloudSshKeysValue](ctx),
			Computed:            true,
			Description:         "List of SSH keys of the project",
			MarkdownDescription: "List of SSH keys of the project",
		},
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Description:         "Service name. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
			MarkdownDescription: "Service name. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
		},
	}

	return schema.Schema{
		Description: "List the SSH keys of a Public Cloud project (API v2).",
		Attributes:  attrs,
	}
}

type CloudSshKeysModel struct {
	SshKeys     ovhtypes.TfListNestedValue[CloudSshKeysValue] `tfsdk:"ssh_keys" json:"sshKeys"`
	ServiceName ovhtypes.TfStringValue                        `tfsdk:"service_name" json:"-"`
}

func (v *CloudSshKeysModel) MergeWith(other *CloudSshKeysModel) {
	if (v.SshKeys.IsUnknown() || v.SshKeys.IsNull()) && !other.SshKeys.IsUnknown() {
		v.SshKeys = other.SshKeys
	}

	if (v.ServiceName.IsUnknown() || v.ServiceName.IsNull()) && !other.ServiceName.IsUnknown() {
		v.ServiceName = other.ServiceName
	}
}

var _ basetypes.ObjectTypable = CloudSshKeysType{}

type CloudSshKeysType struct {
	basetypes.ObjectType
}

func (t CloudSshKeysType) Equal(o attr.Type) bool {
	other, ok := o.(CloudSshKeysType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t CloudSshKeysType) String() string {
	return "CloudSshKeysType"
}

func (t CloudSshKeysType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return nil, diags
	}

	nameVal, ok := nameAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be ovhtypes.TfStringValue, was: %T`, nameAttribute))
	}

	publicKeyAttribute, ok := attributes["public_key"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`public_key is missing from object`)

		return nil, diags
	}

	publicKeyVal, ok := publicKeyAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`public_key expected to be ovhtypes.TfStringValue, was: %T`, publicKeyAttribute))
	}

	createdAtAttribute, ok := attributes["created_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_at is missing from object`)

		return nil, diags
	}

	createdAtVal, ok := createdAtAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_at expected to be ovhtypes.TfStringValue, was: %T`, createdAtAttribute))
	}

	updatedAtAttribute, ok := attributes["updated_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`updated_at is missing from object`)

		return nil, diags
	}

	updatedAtVal, ok := updatedAtAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`updated_at expected to be ovhtypes.TfStringValue, was: %T`, updatedAtAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return CloudSshKeysValue{
		Name:      nameVal,
		PublicKey: publicKeyVal,
		CreatedAt: createdAtVal,
		UpdatedAt: updatedAtVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewCloudSshKeysValueNull() CloudSshKeysValue {
	return CloudSshKeysValue{
		state: attr.ValueStateNull,
	}
}

func NewCloudSshKeysValueUnknown() CloudSshKeysValue {
	return CloudSshKeysValue{
		state: attr.ValueStateUnknown,
	}
}

func NewCloudSshKeysValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CloudSshKeysValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing CloudSshKeysValue Attribute Value",
				"While creating a CloudSshKeysValue value, a missing attribute value was detected. "+
					"A CloudSshKeysValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CloudSshKeysValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid CloudSshKeysValue Attribute Type",
				"While creating a CloudSshKeysValue value, an invalid attribute value was detected. "+
					"A CloudSshKeysValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("CloudSshKeysValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("CloudSshKeysValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra CloudSshKeysValue Attribute Value",
				"While creating a CloudSshKeysValue value, an extra attribute value was detected. "+
					"A CloudSshKeysValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra CloudSshKeysValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewCloudSshKeysValueUnknown(), diags
	}

	nameAttribute, ok := attributes["name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`name is missing from object`)

		return NewCloudSshKeysValueUnknown(), diags
	}

	nameVal, ok := nameAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`name expected to be ovhtypes.TfStringValue, was: %T`, nameAttribute))
	}

	publicKeyAttribute, ok := attributes["public_key"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`public_key is missing from object`)

		return NewCloudSshKeysValueUnknown(), diags
	}

	publicKeyVal, ok := publicKeyAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`public_key expected to be ovhtypes.TfStringValue, was: %T`, publicKeyAttribute))
	}

	createdAtAttribute, ok := attributes["created_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`created_at is missing from object`)

		return NewCloudSshKeysValueUnknown(), diags
	}

	createdAtVal, ok := createdAtAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`created_at expected to be ovhtypes.TfStringValue, was: %T`, createdAtAttribute))
	}

	updatedAtAttribute, ok := attributes["updated_at"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`updated_at is missing from object`)

		return NewCloudSshKeysValueUnknown(), diags
	}

	updatedAtVal, ok := updatedAtAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`updated_at expected to be ovhtypes.TfStringValue, was: %T`, updatedAtAttribute))
	}

	if diags.HasError() {
		return NewCloudSshKeysValueUnknown(), diags
	}

	return CloudSshKeysValue{
		Name:      nameVal,
		PublicKey: publicKeyVal,
		CreatedAt: createdAtVal,
		UpdatedAt: updatedAtVal,
		state:     attr.ValueStateKnown,
	}, diags
}

func NewCloudSshKeysValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) CloudSshKeysValue {
	object, diags := NewCloudSshKeysValue(attributeTypes, attributes)

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

		panic("NewCloudSshKeysValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t CloudSshKeysType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewCloudSshKeysValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewCloudSshKeysValueUnknown(), nil
	}

	if in.IsNull() {
		return NewCloudSshKeysValueNull(), nil
	}

	attributes := map[string]attr.Value{}

	val := map[string]tftypes.Value{}

	err := in.As(&val)

	if err != nil {
		return nil, err
	}

	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)

		if err != nil {
			return nil, err
		}

		attributes[k] = a
	}

	return NewCloudSshKeysValueMust(CloudSshKeysValue{}.AttributeTypes(ctx), attributes), nil
}

func (t CloudSshKeysType) ValueType(ctx context.Context) attr.Value {
	return CloudSshKeysValue{}
}

var _ basetypes.ObjectValuable = CloudSshKeysValue{}

type CloudSshKeysValue struct {
	Name      ovhtypes.TfStringValue `tfsdk:"name" json:"name"`
	PublicKey ovhtypes.TfStringValue `tfsdk:"public_key" json:"publicKey"`
	CreatedAt ovhtypes.TfStringValue `tfsdk:"created_at" json:"createdAt"`
	UpdatedAt ovhtypes.TfStringValue `tfsdk:"updated_at" json:"updatedAt"`
	state     attr.ValueState
}

func (v *CloudSshKeysValue) UnmarshalJSON(data []byte) error {
	type JsonCloudSshKeysValue CloudSshKeysValue

	var tmp JsonCloudSshKeysValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Name = tmp.Name
	v.PublicKey = tmp.PublicKey
	v.CreatedAt = tmp.CreatedAt
	v.UpdatedAt = tmp.UpdatedAt

	v.state = attr.ValueStateKnown

	return nil
}

func (v *CloudSshKeysValue) MergeWith(other *CloudSshKeysValue) {

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

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v CloudSshKeysValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"name":       v.Name,
		"public_key": v.PublicKey,
		"created_at": v.CreatedAt,
		"updated_at": v.UpdatedAt,
	}
}

func (v CloudSshKeysValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 4)

	var val tftypes.Value
	var err error

	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["public_key"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["created_at"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["updated_at"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 4)

		val, err = v.Name.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["name"] = val

		val, err = v.PublicKey.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["public_key"] = val

		val, err = v.CreatedAt.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["created_at"] = val

		val, err = v.UpdatedAt.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["updated_at"] = val

		if err := tftypes.ValidateValue(objectType, vals); err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		return tftypes.NewValue(objectType, vals), nil
	case attr.ValueStateNull:
		return tftypes.NewValue(objectType, nil), nil
	case attr.ValueStateUnknown:
		return tftypes.NewValue(objectType, tftypes.UnknownValue), nil
	default:
		panic(fmt.Sprintf("unhandled Object state in ToTerraformValue: %s", v.state))
	}
}

func (v CloudSshKeysValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CloudSshKeysValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CloudSshKeysValue) String() string {
	return "CloudSshKeysValue"
}

func (v CloudSshKeysValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"name":       ovhtypes.TfStringType{},
			"public_key": ovhtypes.TfStringType{},
			"created_at": ovhtypes.TfStringType{},
			"updated_at": ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"name":       v.Name,
			"public_key": v.PublicKey,
			"created_at": v.CreatedAt,
			"updated_at": v.UpdatedAt,
		})

	return objVal, diags
}

func (v CloudSshKeysValue) Equal(o attr.Value) bool {
	other, ok := o.(CloudSshKeysValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.Name.Equal(other.Name) {
		return false
	}

	if !v.PublicKey.Equal(other.PublicKey) {
		return false
	}

	if !v.CreatedAt.Equal(other.CreatedAt) {
		return false
	}

	if !v.UpdatedAt.Equal(other.UpdatedAt) {
		return false
	}

	return true
}

func (v CloudSshKeysValue) Type(ctx context.Context) attr.Type {
	return CloudSshKeysType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v CloudSshKeysValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"name":       ovhtypes.TfStringType{},
		"public_key": ovhtypes.TfStringType{},
		"created_at": ovhtypes.TfStringType{},
		"updated_at": ovhtypes.TfStringType{},
	}
}
