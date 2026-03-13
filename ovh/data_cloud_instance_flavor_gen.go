package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

func CloudInstanceFlavorDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name",
			MarkdownDescription: "Service name",
		},
		"region_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Region name",
			MarkdownDescription: "Region name",
		},
		"flavor_id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Flavor ID",
			MarkdownDescription: "Flavor ID",
		},
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Flavor ID",
			MarkdownDescription: "Flavor ID",
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Flavor name",
			MarkdownDescription: "Flavor name",
		},
		"vcpus": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Number of virtual CPUs",
			MarkdownDescription: "Number of virtual CPUs",
		},
		"ram": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "RAM in MB",
			MarkdownDescription: "RAM in MB",
		},
		"disk": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Root disk size in GB",
			MarkdownDescription: "Root disk size in GB",
		},
		"swap": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Swap space in MB",
			MarkdownDescription: "Swap space in MB",
		},
		"ephemeral": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Ephemeral disk size in GB",
			MarkdownDescription: "Ephemeral disk size in GB",
		},
		"is_public": schema.BoolAttribute{
			CustomType:          ovhtypes.TfBoolType{},
			Computed:            true,
			Description:         "Whether the flavor is publicly available",
			MarkdownDescription: "Whether the flavor is publicly available",
		},
		"description": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Flavor description",
			MarkdownDescription: "Flavor description",
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type CloudInstanceFlavorModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name" json:"-"`
	RegionName  ovhtypes.TfStringValue `tfsdk:"region_name" json:"-"`
	FlavorId    ovhtypes.TfStringValue `tfsdk:"flavor_id" json:"-"`
	Id          ovhtypes.TfStringValue `tfsdk:"id" json:"id"`
	Name        ovhtypes.TfStringValue `tfsdk:"name" json:"name"`
	VCPUs       ovhtypes.TfInt64Value  `tfsdk:"vcpus" json:"vcpus"`
	RAM         ovhtypes.TfInt64Value  `tfsdk:"ram" json:"ram"`
	Disk        ovhtypes.TfInt64Value  `tfsdk:"disk" json:"disk"`
	Swap        ovhtypes.TfInt64Value  `tfsdk:"swap" json:"swap"`
	Ephemeral   ovhtypes.TfInt64Value  `tfsdk:"ephemeral" json:"ephemeral"`
	IsPublic    ovhtypes.TfBoolValue   `tfsdk:"is_public" json:"isPublic"`
	Description ovhtypes.TfStringValue `tfsdk:"description" json:"description"`
}

func CloudInstanceFlavorsDataSourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name",
			MarkdownDescription: "Service name",
		},
		"region_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Region name",
			MarkdownDescription: "Region name",
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Description:         "Filter flavors by name (regexp match, e.g. 'b2-.*')",
			MarkdownDescription: "Filter flavors by name (regexp match, e.g. `b2-.*`)",
		},
		"flavors": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Flavor ID",
						MarkdownDescription: "Flavor ID",
					},
					"name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Flavor name",
						MarkdownDescription: "Flavor name",
					},
					"vcpus": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "Number of virtual CPUs",
						MarkdownDescription: "Number of virtual CPUs",
					},
					"ram": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "RAM in MB",
						MarkdownDescription: "RAM in MB",
					},
					"disk": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "Root disk size in GB",
						MarkdownDescription: "Root disk size in GB",
					},
					"swap": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "Swap space in MB",
						MarkdownDescription: "Swap space in MB",
					},
					"ephemeral": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "Ephemeral disk size in GB",
						MarkdownDescription: "Ephemeral disk size in GB",
					},
					"is_public": schema.BoolAttribute{
						CustomType:          ovhtypes.TfBoolType{},
						Computed:            true,
						Description:         "Whether the flavor is publicly available",
						MarkdownDescription: "Whether the flavor is publicly available",
					},
					"description": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Flavor description",
						MarkdownDescription: "Flavor description",
					},
				},
				CustomType: CloudInstanceFlavorsValueType{
					ObjectType: types.ObjectType{
						AttrTypes: CloudInstanceFlavorsValue{}.AttributeTypes(ctx),
					},
				},
			},
			CustomType: ovhtypes.NewTfListNestedType[CloudInstanceFlavorsValue](ctx),
			Computed:   true,
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type CloudInstanceFlavorsModel struct {
	ServiceName ovhtypes.TfStringValue                                `tfsdk:"service_name" json:"-"`
	RegionName  ovhtypes.TfStringValue                                `tfsdk:"region_name" json:"-"`
	Name        ovhtypes.TfStringValue                                `tfsdk:"name" json:"-"`
	Flavors     ovhtypes.TfListNestedValue[CloudInstanceFlavorsValue] `tfsdk:"flavors" json:"flavors"`
}

// --- CloudInstanceFlavorsValue (nested object for list items) ---

var _ basetypes.ObjectTypable = CloudInstanceFlavorsValueType{}

type CloudInstanceFlavorsValueType struct {
	basetypes.ObjectType
}

func (t CloudInstanceFlavorsValueType) Equal(o attr.Type) bool {
	other, ok := o.(CloudInstanceFlavorsValueType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

func (t CloudInstanceFlavorsValueType) String() string {
	return "CloudInstanceFlavorsValueType"
}

func (t CloudInstanceFlavorsValueType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics
	attributes := in.Attributes()

	idVal, ok := attributes["id"].(ovhtypes.TfStringValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`id expected to be ovhtypes.TfStringValue, was: %T`, attributes["id"]))
		return nil, diags
	}
	nameVal, ok := attributes["name"].(ovhtypes.TfStringValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`name expected to be ovhtypes.TfStringValue, was: %T`, attributes["name"]))
		return nil, diags
	}
	vcpusVal, ok := attributes["vcpus"].(ovhtypes.TfInt64Value)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`vcpus expected to be ovhtypes.TfInt64Value, was: %T`, attributes["vcpus"]))
		return nil, diags
	}
	ramVal, ok := attributes["ram"].(ovhtypes.TfInt64Value)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`ram expected to be ovhtypes.TfInt64Value, was: %T`, attributes["ram"]))
		return nil, diags
	}
	diskVal, ok := attributes["disk"].(ovhtypes.TfInt64Value)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`disk expected to be ovhtypes.TfInt64Value, was: %T`, attributes["disk"]))
		return nil, diags
	}
	swapVal, ok := attributes["swap"].(ovhtypes.TfInt64Value)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`swap expected to be ovhtypes.TfInt64Value, was: %T`, attributes["swap"]))
		return nil, diags
	}
	ephemeralVal, ok := attributes["ephemeral"].(ovhtypes.TfInt64Value)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`ephemeral expected to be ovhtypes.TfInt64Value, was: %T`, attributes["ephemeral"]))
		return nil, diags
	}
	isPublicVal, ok := attributes["is_public"].(ovhtypes.TfBoolValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`is_public expected to be ovhtypes.TfBoolValue, was: %T`, attributes["is_public"]))
		return nil, diags
	}
	descriptionVal, ok := attributes["description"].(ovhtypes.TfStringValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`description expected to be ovhtypes.TfStringValue, was: %T`, attributes["description"]))
		return nil, diags
	}

	if diags.HasError() {
		return nil, diags
	}

	return CloudInstanceFlavorsValue{
		Id:          idVal,
		Name:        nameVal,
		VCPUs:       vcpusVal,
		RAM:         ramVal,
		Disk:        diskVal,
		Swap:        swapVal,
		Ephemeral:   ephemeralVal,
		IsPublic:    isPublicVal,
		Description: descriptionVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func (t CloudInstanceFlavorsValueType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return CloudInstanceFlavorsValue{state: attr.ValueStateNull}, nil
	}
	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}
	if !in.IsKnown() {
		return CloudInstanceFlavorsValue{state: attr.ValueStateUnknown}, nil
	}
	if in.IsNull() {
		return CloudInstanceFlavorsValue{state: attr.ValueStateNull}, nil
	}

	attributes := map[string]attr.Value{}
	val := map[string]tftypes.Value{}
	if err := in.As(&val); err != nil {
		return nil, err
	}
	for k, v := range val {
		a, err := t.AttrTypes[k].ValueFromTerraform(ctx, v)
		if err != nil {
			return nil, err
		}
		attributes[k] = a
	}

	result, diags := NewCloudInstanceFlavorsValue(CloudInstanceFlavorsValue{}.AttributeTypes(ctx), attributes)
	if diags.HasError() {
		diagsStrings := make([]string, 0, len(diags))
		for _, d := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf("%s | %s | %s", d.Severity(), d.Summary(), d.Detail()))
		}
		return nil, fmt.Errorf("error creating CloudInstanceFlavorsValue: %s", strings.Join(diagsStrings, "\n"))
	}
	return result, nil
}

func (t CloudInstanceFlavorsValueType) ValueType(ctx context.Context) attr.Value {
	return CloudInstanceFlavorsValue{}
}

var _ basetypes.ObjectValuable = CloudInstanceFlavorsValue{}

type CloudInstanceFlavorsValue struct {
	Id          ovhtypes.TfStringValue `tfsdk:"id" json:"id"`
	Name        ovhtypes.TfStringValue `tfsdk:"name" json:"name"`
	VCPUs       ovhtypes.TfInt64Value  `tfsdk:"vcpus" json:"vcpus"`
	RAM         ovhtypes.TfInt64Value  `tfsdk:"ram" json:"ram"`
	Disk        ovhtypes.TfInt64Value  `tfsdk:"disk" json:"disk"`
	Swap        ovhtypes.TfInt64Value  `tfsdk:"swap" json:"swap"`
	Ephemeral   ovhtypes.TfInt64Value  `tfsdk:"ephemeral" json:"ephemeral"`
	IsPublic    ovhtypes.TfBoolValue   `tfsdk:"is_public" json:"isPublic"`
	Description ovhtypes.TfStringValue `tfsdk:"description" json:"description"`
	state       attr.ValueState
}

func (v *CloudInstanceFlavorsValue) UnmarshalJSON(data []byte) error {
	type JsonCloudInstanceFlavorsValue CloudInstanceFlavorsValue
	var tmp JsonCloudInstanceFlavorsValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Id = tmp.Id
	v.Name = tmp.Name
	v.VCPUs = tmp.VCPUs
	v.RAM = tmp.RAM
	v.Disk = tmp.Disk
	v.Swap = tmp.Swap
	v.Ephemeral = tmp.Ephemeral
	v.IsPublic = tmp.IsPublic
	v.Description = tmp.Description
	v.state = attr.ValueStateKnown
	return nil
}

func NewCloudInstanceFlavorsValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CloudInstanceFlavorsValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]
		if !ok {
			diags.AddError("Missing Attribute", fmt.Sprintf("CloudInstanceFlavorsValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()))
			continue
		}
		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError("Invalid Attribute Type", fmt.Sprintf("Expected: %s, Given: %s", attributeType.String(), attribute.Type(ctx)))
		}
	}
	if diags.HasError() {
		return CloudInstanceFlavorsValue{state: attr.ValueStateUnknown}, diags
	}

	return CloudInstanceFlavorsValue{
		Id:          attributes["id"].(ovhtypes.TfStringValue),
		Name:        attributes["name"].(ovhtypes.TfStringValue),
		VCPUs:       attributes["vcpus"].(ovhtypes.TfInt64Value),
		RAM:         attributes["ram"].(ovhtypes.TfInt64Value),
		Disk:        attributes["disk"].(ovhtypes.TfInt64Value),
		Swap:        attributes["swap"].(ovhtypes.TfInt64Value),
		Ephemeral:   attributes["ephemeral"].(ovhtypes.TfInt64Value),
		IsPublic:    attributes["is_public"].(ovhtypes.TfBoolValue),
		Description: attributes["description"].(ovhtypes.TfStringValue),
		state:       attr.ValueStateKnown,
	}, diags
}

func (v CloudInstanceFlavorsValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":          ovhtypes.TfStringType{},
		"name":        ovhtypes.TfStringType{},
		"vcpus":       ovhtypes.TfInt64Type{},
		"ram":         ovhtypes.TfInt64Type{},
		"disk":        ovhtypes.TfInt64Type{},
		"swap":        ovhtypes.TfInt64Type{},
		"ephemeral":   ovhtypes.TfInt64Type{},
		"is_public":   ovhtypes.TfBoolType{},
		"description": ovhtypes.TfStringType{},
	}
}

func (v CloudInstanceFlavorsValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"id":          v.Id,
		"name":        v.Name,
		"vcpus":       v.VCPUs,
		"ram":         v.RAM,
		"disk":        v.Disk,
		"swap":        v.Swap,
		"ephemeral":   v.Ephemeral,
		"is_public":   v.IsPublic,
		"description": v.Description,
	}
}

func (v CloudInstanceFlavorsValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 9)
	attrTypes["id"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["vcpus"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["ram"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["disk"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["swap"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["ephemeral"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["is_public"] = basetypes.BoolType{}.TerraformType(ctx)
	attrTypes["description"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 9)
		var val tftypes.Value
		var err error

		val, err = v.Id.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["id"] = val

		val, err = v.Name.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["name"] = val

		val, err = v.VCPUs.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["vcpus"] = val

		val, err = v.RAM.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["ram"] = val

		val, err = v.Disk.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["disk"] = val

		val, err = v.Swap.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["swap"] = val

		val, err = v.Ephemeral.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["ephemeral"] = val

		val, err = v.IsPublic.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["is_public"] = val

		val, err = v.Description.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["description"] = val

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

func (v CloudInstanceFlavorsValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CloudInstanceFlavorsValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CloudInstanceFlavorsValue) String() string {
	return "CloudInstanceFlavorsValue"
}

func (v CloudInstanceFlavorsValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(v.AttributeTypes(ctx), v.Attributes())
}

func (v CloudInstanceFlavorsValue) Equal(o attr.Value) bool {
	other, ok := o.(CloudInstanceFlavorsValue)
	if !ok {
		return false
	}
	if v.state != other.state {
		return false
	}
	if v.state != attr.ValueStateKnown {
		return true
	}
	return v.Id.Equal(other.Id) && v.Name.Equal(other.Name) && v.VCPUs.Equal(other.VCPUs) &&
		v.RAM.Equal(other.RAM) && v.Disk.Equal(other.Disk) && v.Swap.Equal(other.Swap) &&
		v.Ephemeral.Equal(other.Ephemeral) && v.IsPublic.Equal(other.IsPublic) && v.Description.Equal(other.Description)
}

func (v CloudInstanceFlavorsValue) Type(ctx context.Context) attr.Type {
	return CloudInstanceFlavorsValueType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}
