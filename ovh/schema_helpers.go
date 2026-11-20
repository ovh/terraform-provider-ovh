package ovh

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	dschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/numberplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// Optional attributes are also Computed with UseStateForUnknown: the API fills what the
// config omits (defaults, imported values) and the state keeps it with no plan diff.

func optionalString(desc string, mods ...planmodifier.String) schema.StringAttribute {
	return schema.StringAttribute{
		CustomType: ovhtypes.TfStringType{}, Optional: true, Computed: true,
		Description:   desc,
		PlanModifiers: append([]planmodifier.String{stringplanmodifier.UseStateForUnknown()}, mods...),
	}
}

func requiredString(desc string, mods ...planmodifier.String) schema.StringAttribute {
	return schema.StringAttribute{
		CustomType: ovhtypes.TfStringType{}, Required: true,
		Description: desc, PlanModifiers: mods,
	}
}

func withOneOf(a schema.StringAttribute, values ...string) schema.StringAttribute {
	a.Validators = []validator.String{stringvalidator.OneOf(values...)}
	return a
}

// withVariants rejects a variant block that does not match the discriminator value;
// each block is named after its lowercased value (CILIUM → cilium).
func withVariants(a schema.SingleNestedAttribute, discriminator string, variants ...string) schema.SingleNestedAttribute {
	a.Validators = []validator.Object{variantValidator{discriminator, variants}}
	return a
}

type variantValidator struct {
	discriminator string
	variants      []string
}

func (v variantValidator) Description(context.Context) string {
	return fmt.Sprintf("only the block named after %s may be set among %s", v.discriminator, strings.Join(v.variants, ", "))
}

func (v variantValidator) MarkdownDescription(ctx context.Context) string { return v.Description(ctx) }

func (v variantValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	attrs := req.ConfigValue.Attributes()
	d, ok := attrs[v.discriminator].(basetypes.StringValuable)
	if !ok {
		return
	}
	value, diags := d.ToStringValue(ctx)
	if diags.HasError() || value.IsNull() || value.IsUnknown() {
		return
	}
	for _, name := range v.variants {
		if b := attrs[name]; name != strings.ToLower(value.ValueString()) && b != nil && !b.IsNull() && !b.IsUnknown() {
			resp.Diagnostics.AddAttributeError(req.Path.AtName(name), "Variant block does not match "+v.discriminator,
				fmt.Sprintf("%s is set while %s is %s", name, v.discriminator, value.ValueString()))
		}
	}
}

func optionalBool(desc string, mods ...planmodifier.Bool) schema.BoolAttribute {
	return schema.BoolAttribute{
		CustomType: ovhtypes.TfBoolType{}, Optional: true, Computed: true,
		Description:   desc,
		PlanModifiers: append([]planmodifier.Bool{boolplanmodifier.UseStateForUnknown()}, mods...),
	}
}

func optionalNumber(desc string) schema.NumberAttribute {
	return schema.NumberAttribute{
		CustomType: ovhtypes.TfNumberType{}, Optional: true, Computed: true,
		Description:   desc,
		PlanModifiers: []planmodifier.Number{numberplanmodifier.UseStateForUnknown()},
	}
}

func optionalInt64(desc string) schema.Int64Attribute {
	return schema.Int64Attribute{
		CustomType: ovhtypes.TfInt64Type{}, Optional: true, Computed: true,
		Description:   desc,
		PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
	}
}

func optionalStringList(desc string) schema.ListAttribute {
	return schema.ListAttribute{
		ElementType: ovhtypes.TfStringType{}, Optional: true, Computed: true,
		Description:   desc,
		PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
	}
}

func optionalStringMap(desc string) schema.MapAttribute {
	return schema.MapAttribute{
		ElementType: ovhtypes.TfStringType{}, Optional: true, Computed: true,
		Description:   desc,
		PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
	}
}

func optionalObject(desc string, attrs map[string]schema.Attribute) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional: true, Computed: true, Description: desc, Attributes: attrs,
		PlanModifiers: []planmodifier.Object{objectplanmodifier.UseStateForUnknown()},
	}
}

func optionalNestedList(desc string, attrs map[string]schema.Attribute) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		Optional: true, Computed: true, Description: desc,
		NestedObject:  schema.NestedAttributeObject{Attributes: attrs},
		PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
	}
}

func optionalNestedMap(desc string, attrs map[string]schema.Attribute) schema.MapNestedAttribute {
	return schema.MapNestedAttribute{
		Optional: true, Computed: true, Description: desc,
		NestedObject:  schema.NestedAttributeObject{Attributes: attrs},
		PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
	}
}

// computedAttribute deep-copies a resource attribute as Computed-only, without
// plan modifiers or validators.
func computedAttribute(a schema.Attribute) schema.Attribute {
	d, md := a.GetDescription(), a.GetMarkdownDescription()
	nested := func(in map[string]schema.Attribute) map[string]schema.Attribute {
		out := make(map[string]schema.Attribute, len(in))
		for k, v := range in {
			out[k] = computedAttribute(v)
		}
		return out
	}
	switch v := a.(type) {
	case schema.StringAttribute:
		return schema.StringAttribute{CustomType: v.CustomType, Computed: true, Sensitive: v.Sensitive, Description: d, MarkdownDescription: md}
	case schema.BoolAttribute:
		return schema.BoolAttribute{CustomType: v.CustomType, Computed: true, Description: d, MarkdownDescription: md}
	case schema.Int64Attribute:
		return schema.Int64Attribute{CustomType: v.CustomType, Computed: true, Description: d, MarkdownDescription: md}
	case schema.ListAttribute:
		return schema.ListAttribute{ElementType: v.ElementType, Computed: true, Description: d, MarkdownDescription: md}
	case schema.NumberAttribute:
		return schema.NumberAttribute{CustomType: v.CustomType, Computed: true, Description: d, MarkdownDescription: md}
	case schema.MapAttribute:
		return schema.MapAttribute{ElementType: v.ElementType, Computed: true, Description: d, MarkdownDescription: md}
	case schema.SingleNestedAttribute:
		return schema.SingleNestedAttribute{Attributes: nested(v.Attributes), Computed: true, Description: d, MarkdownDescription: md}
	case schema.ListNestedAttribute:
		return schema.ListNestedAttribute{NestedObject: schema.NestedAttributeObject{Attributes: nested(v.NestedObject.Attributes)}, Computed: true, Description: d, MarkdownDescription: md}
	case schema.MapNestedAttribute:
		return schema.MapNestedAttribute{NestedObject: schema.NestedAttributeObject{Attributes: nested(v.NestedObject.Attributes)}, Computed: true, Description: d, MarkdownDescription: md}
	}
	panic(fmt.Sprintf("computedAttribute: unsupported attribute type %T", a))
}

func cloudServiceNameAttribute() dschema.StringAttribute {
	return dataSourceRequiredString("Service name of the resource representing the id of the cloud project")
}

func dataSourceRequiredString(desc string) dschema.StringAttribute {
	return dschema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Required: true, Description: desc}
}

func dataSourceComputedString(desc string) dschema.StringAttribute {
	return dschema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: desc}
}

func dataSourceComputedInt64(desc string) dschema.Int64Attribute {
	return dschema.Int64Attribute{CustomType: ovhtypes.TfInt64Type{}, Computed: true, Description: desc}
}

func dataSourceComputedStringList(desc string) dschema.ListAttribute {
	return dschema.ListAttribute{ElementType: ovhtypes.TfStringType{}, Computed: true, Description: desc}
}

func dataSourceComputedBool(desc string) dschema.BoolAttribute {
	return dschema.BoolAttribute{CustomType: ovhtypes.TfBoolType{}, Computed: true, Description: desc}
}

func dataSourceComputedObject(desc string, attrs map[string]dschema.Attribute) dschema.SingleNestedAttribute {
	return dschema.SingleNestedAttribute{Computed: true, Description: desc, Attributes: attrs}
}

func dataSourceComputedList(desc string, attrs map[string]dschema.Attribute) dschema.ListNestedAttribute {
	return dschema.ListNestedAttribute{Computed: true, Description: desc, NestedObject: dschema.NestedAttributeObject{Attributes: attrs}}
}

// dataSourceAttribute converts a resource attribute into a Computed-only data
// source attribute of the same type.
func dataSourceAttribute(a schema.Attribute) dschema.Attribute {
	d, md := a.GetDescription(), a.GetMarkdownDescription()
	nested := func(in map[string]schema.Attribute) map[string]dschema.Attribute {
		out := make(map[string]dschema.Attribute, len(in))
		for k, v := range in {
			out[k] = dataSourceAttribute(v)
		}
		return out
	}
	switch v := a.(type) {
	case schema.StringAttribute:
		return dschema.StringAttribute{CustomType: v.CustomType, Computed: true, Sensitive: v.Sensitive, Description: d, MarkdownDescription: md}
	case schema.BoolAttribute:
		return dschema.BoolAttribute{CustomType: v.CustomType, Computed: true, Description: d, MarkdownDescription: md}
	case schema.Int64Attribute:
		return dschema.Int64Attribute{CustomType: v.CustomType, Computed: true, Description: d, MarkdownDescription: md}
	case schema.ListAttribute:
		return dschema.ListAttribute{ElementType: v.ElementType, Computed: true, Description: d, MarkdownDescription: md}
	case schema.NumberAttribute:
		return dschema.NumberAttribute{CustomType: v.CustomType, Computed: true, Description: d, MarkdownDescription: md}
	case schema.MapAttribute:
		return dschema.MapAttribute{ElementType: v.ElementType, Computed: true, Description: d, MarkdownDescription: md}
	case schema.SingleNestedAttribute:
		return dschema.SingleNestedAttribute{Attributes: nested(v.Attributes), Computed: true, Description: d, MarkdownDescription: md}
	case schema.ListNestedAttribute:
		return dschema.ListNestedAttribute{NestedObject: dschema.NestedAttributeObject{Attributes: nested(v.NestedObject.Attributes)}, Computed: true, Description: d, MarkdownDescription: md}
	case schema.MapNestedAttribute:
		return dschema.MapNestedAttribute{NestedObject: dschema.NestedAttributeObject{Attributes: nested(v.NestedObject.Attributes)}, Computed: true, Description: d, MarkdownDescription: md}
	}
	panic(fmt.Sprintf("dataSourceAttribute: unsupported attribute type %T", a))
}

// objectAttrTypes returns the attribute types of every object attribute.
func objectAttrTypes(attrs map[string]schema.Attribute) map[string]map[string]attr.Type {
	out := map[string]map[string]attr.Type{}
	for k, a := range attrs {
		if o, ok := a.GetType().(types.ObjectType); ok {
			out[k] = o.AttrTypes
		}
	}
	return out
}

// objectAs reads a Terraform object into a struct pointer; null or unknown yields nil.
func objectAs[T any](ctx context.Context, obj types.Object, diags *diag.Diagnostics) *T {
	var out *T
	diags.Append(obj.As(ctx, &out, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true})...)
	return out
}

// objectFrom builds the Terraform object of a root attribute; a nil pointer yields a typed null.
func objectFrom[T any](ctx context.Context, attrTypes map[string]attr.Type, v *T, diags *diag.Diagnostics) types.Object {
	obj, d := types.ObjectValueFrom(ctx, attrTypes, v)
	diags.Append(d...)
	return obj
}
