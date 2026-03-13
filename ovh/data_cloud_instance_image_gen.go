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

func CloudInstanceImageDataSourceSchema(ctx context.Context) schema.Schema {
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
		"image_id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Image ID",
			MarkdownDescription: "Image ID",
		},
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Image ID",
			MarkdownDescription: "Image ID",
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Image name",
			MarkdownDescription: "Image name",
		},
		"status": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Image status",
			MarkdownDescription: "Image status",
		},
		"visibility": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Image visibility",
			MarkdownDescription: "Image visibility",
		},
		"min_disk": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Minimum disk size in GB",
			MarkdownDescription: "Minimum disk size in GB",
		},
		"min_ram": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Minimum RAM in MB",
			MarkdownDescription: "Minimum RAM in MB",
		},
		"size": schema.Int64Attribute{
			CustomType:          ovhtypes.TfInt64Type{},
			Computed:            true,
			Description:         "Image size in bytes",
			MarkdownDescription: "Image size in bytes",
		},
		"created_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Image creation date",
			MarkdownDescription: "Image creation date",
		},
		"updated_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Image last update date",
			MarkdownDescription: "Image last update date",
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type CloudInstanceImageModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name" json:"-"`
	RegionName  ovhtypes.TfStringValue `tfsdk:"region_name" json:"-"`
	ImageId     ovhtypes.TfStringValue `tfsdk:"image_id" json:"-"`
	Id          ovhtypes.TfStringValue `tfsdk:"id" json:"id"`
	Name        ovhtypes.TfStringValue `tfsdk:"name" json:"name"`
	Status      ovhtypes.TfStringValue `tfsdk:"status" json:"status"`
	Visibility  ovhtypes.TfStringValue `tfsdk:"visibility" json:"visibility"`
	MinDisk     ovhtypes.TfInt64Value  `tfsdk:"min_disk" json:"minDisk"`
	MinRam      ovhtypes.TfInt64Value  `tfsdk:"min_ram" json:"minRam"`
	Size        ovhtypes.TfInt64Value  `tfsdk:"size" json:"size"`
	CreatedAt   ovhtypes.TfStringValue `tfsdk:"created_at" json:"createdAt"`
	UpdatedAt   ovhtypes.TfStringValue `tfsdk:"updated_at" json:"updatedAt"`
}

func CloudInstanceImagesDataSourceSchema(ctx context.Context) schema.Schema {
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
			Description:         "Filter images by name (regexp match, e.g. 'Debian 13')",
			MarkdownDescription: "Filter images by name (regexp match, e.g. `Debian 13`)",
		},
		"images": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Image ID",
						MarkdownDescription: "Image ID",
					},
					"name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Image name",
						MarkdownDescription: "Image name",
					},
					"status": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Image status",
						MarkdownDescription: "Image status",
					},
					"visibility": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Image visibility",
						MarkdownDescription: "Image visibility",
					},
					"min_disk": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "Minimum disk size in GB",
						MarkdownDescription: "Minimum disk size in GB",
					},
					"min_ram": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "Minimum RAM in MB",
						MarkdownDescription: "Minimum RAM in MB",
					},
					"size": schema.Int64Attribute{
						CustomType:          ovhtypes.TfInt64Type{},
						Computed:            true,
						Description:         "Image size in bytes",
						MarkdownDescription: "Image size in bytes",
					},
					"created_at": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Image creation date",
						MarkdownDescription: "Image creation date",
					},
					"updated_at": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Image last update date",
						MarkdownDescription: "Image last update date",
					},
				},
				CustomType: CloudInstanceImagesValueType{
					ObjectType: types.ObjectType{
						AttrTypes: CloudInstanceImagesValue{}.AttributeTypes(ctx),
					},
				},
			},
			CustomType: ovhtypes.NewTfListNestedType[CloudInstanceImagesValue](ctx),
			Computed:   true,
		},
	}

	return schema.Schema{
		Attributes: attrs,
	}
}

type CloudInstanceImagesModel struct {
	ServiceName ovhtypes.TfStringValue                               `tfsdk:"service_name" json:"-"`
	RegionName  ovhtypes.TfStringValue                               `tfsdk:"region_name" json:"-"`
	Name        ovhtypes.TfStringValue                               `tfsdk:"name" json:"-"`
	Images      ovhtypes.TfListNestedValue[CloudInstanceImagesValue] `tfsdk:"images" json:"images"`
}

// --- CloudInstanceImagesValue (nested object for list items) ---

var _ basetypes.ObjectTypable = CloudInstanceImagesValueType{}

type CloudInstanceImagesValueType struct {
	basetypes.ObjectType
}

func (t CloudInstanceImagesValueType) Equal(o attr.Type) bool {
	other, ok := o.(CloudInstanceImagesValueType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

func (t CloudInstanceImagesValueType) String() string {
	return "CloudInstanceImagesValueType"
}

func (t CloudInstanceImagesValueType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
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
	statusVal, ok := attributes["status"].(ovhtypes.TfStringValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`status expected to be ovhtypes.TfStringValue, was: %T`, attributes["status"]))
		return nil, diags
	}
	visibilityVal, ok := attributes["visibility"].(ovhtypes.TfStringValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`visibility expected to be ovhtypes.TfStringValue, was: %T`, attributes["visibility"]))
		return nil, diags
	}
	minDiskVal, ok := attributes["min_disk"].(ovhtypes.TfInt64Value)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`min_disk expected to be ovhtypes.TfInt64Value, was: %T`, attributes["min_disk"]))
		return nil, diags
	}
	minRamVal, ok := attributes["min_ram"].(ovhtypes.TfInt64Value)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`min_ram expected to be ovhtypes.TfInt64Value, was: %T`, attributes["min_ram"]))
		return nil, diags
	}
	sizeVal, ok := attributes["size"].(ovhtypes.TfInt64Value)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`size expected to be ovhtypes.TfInt64Value, was: %T`, attributes["size"]))
		return nil, diags
	}
	createdAtVal, ok := attributes["created_at"].(ovhtypes.TfStringValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`created_at expected to be ovhtypes.TfStringValue, was: %T`, attributes["created_at"]))
		return nil, diags
	}
	updatedAtVal, ok := attributes["updated_at"].(ovhtypes.TfStringValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf(`updated_at expected to be ovhtypes.TfStringValue, was: %T`, attributes["updated_at"]))
		return nil, diags
	}

	if diags.HasError() {
		return nil, diags
	}

	return CloudInstanceImagesValue{
		Id:         idVal,
		Name:       nameVal,
		Status:     statusVal,
		Visibility: visibilityVal,
		MinDisk:    minDiskVal,
		MinRam:     minRamVal,
		Size:       sizeVal,
		CreatedAt:  createdAtVal,
		UpdatedAt:  updatedAtVal,
		state:      attr.ValueStateKnown,
	}, diags
}

func (t CloudInstanceImagesValueType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return CloudInstanceImagesValue{state: attr.ValueStateNull}, nil
	}
	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}
	if !in.IsKnown() {
		return CloudInstanceImagesValue{state: attr.ValueStateUnknown}, nil
	}
	if in.IsNull() {
		return CloudInstanceImagesValue{state: attr.ValueStateNull}, nil
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

	result, diags := NewCloudInstanceImagesValue(CloudInstanceImagesValue{}.AttributeTypes(ctx), attributes)
	if diags.HasError() {
		diagsStrings := make([]string, 0, len(diags))
		for _, d := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf("%s | %s | %s", d.Severity(), d.Summary(), d.Detail()))
		}
		return nil, fmt.Errorf("error creating CloudInstanceImagesValue: %s", strings.Join(diagsStrings, "\n"))
	}
	return result, nil
}

func (t CloudInstanceImagesValueType) ValueType(ctx context.Context) attr.Value {
	return CloudInstanceImagesValue{}
}

var _ basetypes.ObjectValuable = CloudInstanceImagesValue{}

type CloudInstanceImagesValue struct {
	Id         ovhtypes.TfStringValue `tfsdk:"id" json:"id"`
	Name       ovhtypes.TfStringValue `tfsdk:"name" json:"name"`
	Status     ovhtypes.TfStringValue `tfsdk:"status" json:"status"`
	Visibility ovhtypes.TfStringValue `tfsdk:"visibility" json:"visibility"`
	MinDisk    ovhtypes.TfInt64Value  `tfsdk:"min_disk" json:"minDisk"`
	MinRam     ovhtypes.TfInt64Value  `tfsdk:"min_ram" json:"minRam"`
	Size       ovhtypes.TfInt64Value  `tfsdk:"size" json:"size"`
	CreatedAt  ovhtypes.TfStringValue `tfsdk:"created_at" json:"createdAt"`
	UpdatedAt  ovhtypes.TfStringValue `tfsdk:"updated_at" json:"updatedAt"`
	state      attr.ValueState
}

func (v *CloudInstanceImagesValue) UnmarshalJSON(data []byte) error {
	type JsonCloudInstanceImagesValue CloudInstanceImagesValue
	var tmp JsonCloudInstanceImagesValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Id = tmp.Id
	v.Name = tmp.Name
	v.Status = tmp.Status
	v.Visibility = tmp.Visibility
	v.MinDisk = tmp.MinDisk
	v.MinRam = tmp.MinRam
	v.Size = tmp.Size
	v.CreatedAt = tmp.CreatedAt
	v.UpdatedAt = tmp.UpdatedAt
	v.state = attr.ValueStateKnown
	return nil
}

func NewCloudInstanceImagesValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (CloudInstanceImagesValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]
		if !ok {
			diags.AddError("Missing Attribute", fmt.Sprintf("CloudInstanceImagesValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()))
			continue
		}
		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError("Invalid Attribute Type", fmt.Sprintf("Expected: %s, Given: %s", attributeType.String(), attribute.Type(ctx)))
		}
	}
	if diags.HasError() {
		return CloudInstanceImagesValue{state: attr.ValueStateUnknown}, diags
	}

	return CloudInstanceImagesValue{
		Id:         attributes["id"].(ovhtypes.TfStringValue),
		Name:       attributes["name"].(ovhtypes.TfStringValue),
		Status:     attributes["status"].(ovhtypes.TfStringValue),
		Visibility: attributes["visibility"].(ovhtypes.TfStringValue),
		MinDisk:    attributes["min_disk"].(ovhtypes.TfInt64Value),
		MinRam:     attributes["min_ram"].(ovhtypes.TfInt64Value),
		Size:       attributes["size"].(ovhtypes.TfInt64Value),
		CreatedAt:  attributes["created_at"].(ovhtypes.TfStringValue),
		UpdatedAt:  attributes["updated_at"].(ovhtypes.TfStringValue),
		state:      attr.ValueStateKnown,
	}, diags
}

func (v CloudInstanceImagesValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":         ovhtypes.TfStringType{},
		"name":       ovhtypes.TfStringType{},
		"status":     ovhtypes.TfStringType{},
		"visibility": ovhtypes.TfStringType{},
		"min_disk":   ovhtypes.TfInt64Type{},
		"min_ram":    ovhtypes.TfInt64Type{},
		"size":       ovhtypes.TfInt64Type{},
		"created_at": ovhtypes.TfStringType{},
		"updated_at": ovhtypes.TfStringType{},
	}
}

func (v CloudInstanceImagesValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"id":         v.Id,
		"name":       v.Name,
		"status":     v.Status,
		"visibility": v.Visibility,
		"min_disk":   v.MinDisk,
		"min_ram":    v.MinRam,
		"size":       v.Size,
		"created_at": v.CreatedAt,
		"updated_at": v.UpdatedAt,
	}
}

func (v CloudInstanceImagesValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 9)
	attrTypes["id"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["status"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["visibility"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["min_disk"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["min_ram"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["size"] = basetypes.Int64Type{}.TerraformType(ctx)
	attrTypes["created_at"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["updated_at"] = basetypes.StringType{}.TerraformType(ctx)

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

		val, err = v.Status.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["status"] = val

		val, err = v.Visibility.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["visibility"] = val

		val, err = v.MinDisk.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["min_disk"] = val

		val, err = v.MinRam.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["min_ram"] = val

		val, err = v.Size.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["size"] = val

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

func (v CloudInstanceImagesValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v CloudInstanceImagesValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v CloudInstanceImagesValue) String() string {
	return "CloudInstanceImagesValue"
}

func (v CloudInstanceImagesValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(v.AttributeTypes(ctx), v.Attributes())
}

func (v CloudInstanceImagesValue) Equal(o attr.Value) bool {
	other, ok := o.(CloudInstanceImagesValue)
	if !ok {
		return false
	}
	if v.state != other.state {
		return false
	}
	if v.state != attr.ValueStateKnown {
		return true
	}
	return v.Id.Equal(other.Id) && v.Name.Equal(other.Name) && v.Status.Equal(other.Status) &&
		v.Visibility.Equal(other.Visibility) && v.MinDisk.Equal(other.MinDisk) && v.MinRam.Equal(other.MinRam) &&
		v.Size.Equal(other.Size) && v.CreatedAt.Equal(other.CreatedAt) && v.UpdatedAt.Equal(other.UpdatedAt)
}

func (v CloudInstanceImagesValue) Type(ctx context.Context) attr.Type {
	return CloudInstanceImagesValueType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}
