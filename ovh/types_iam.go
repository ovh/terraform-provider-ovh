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

type IamReferenceAction struct {
	Action       string   `json:"action"`
	Categories   []string `json:"categories"`
	Description  string   `json:"description"`
	ResourceType string   `json:"resourceType"`
}

func (a *IamReferenceAction) ToMap() map[string]any {
	out := make(map[string]any, 4)

	out["action"] = a.Action
	out["categories"] = a.Categories
	out["description"] = a.Description
	out["resource_type"] = a.ResourceType

	return out
}

type IamPolicy struct {
	Id                string            `json:"id,omitempty"`
	Name              string            `json:"name"`
	Description       string            `json:"description,omitempty"`
	Identities        []string          `json:"identities"`
	Resources         []IamResource     `json:"resources"`
	Permissions       IamPermissions    `json:"permissions"`
	PermissionsGroups []PermissionGroup `json:"permissionsGroups"`
	CreatedAt         string            `json:"createdAt,omitempty"`
	UpdatedAt         string            `json:"updatedAt,omitempty"`
	ReadOnly          bool              `json:"readOnly,omitempty"`
	Owner             string            `json:"owner,omitempty"`
}

type PermissionGroup struct {
	Urn string `json:"urn"`
}

func (p IamPolicy) ToMap() map[string]any {
	out := make(map[string]any, 0)
	out["name"] = p.Name

	out["owner"] = p.Owner
	out["created_at"] = p.CreatedAt
	out["identities"] = p.Identities
	var resources []string
	for _, r := range p.Resources {
		resources = append(resources, r.URN)
	}
	out["resources"] = resources

	// inline allow, except and deny
	allow, except, deny := p.Permissions.ToLists()
	if len(allow) != 0 {
		out["allow"] = allow
	}
	if len(except) != 0 {
		out["except"] = except
	}
	if len(deny) != 0 {
		out["deny"] = deny
	}

	if len(p.PermissionsGroups) != 0 {
		var permGrps []string
		for _, grp := range p.PermissionsGroups {
			permGrps = append(permGrps, grp.Urn)
		}

		out["permissions_groups"] = permGrps
	}

	if p.Description != "" {
		out["description"] = p.Description
	}
	if p.ReadOnly {
		out["read_only"] = p.ReadOnly
	}
	if p.UpdatedAt != "" {
		out["updated_at"] = p.UpdatedAt
	}

	return out
}

// IamResource represent a possible information returned when viewing a policy
type IamResource struct {
	// URN is always returned and is the urn of the resource or resource group
	// can also be a globing urn, for example "urn:v1:eu:resource:*"
	URN string `json:"urn,omitempty"`
	// If the urn is a resourceGroup. the details are also present
	Group *IamResourceGroup `json:"group,omitempty"`
	// If the urn is an IAM resource, the details are also present
	Resource *IamResourceDetails `json:"resource,omitempty"`
}

type IamResourceDetails struct {
	Id          string            `json:"id,omitempty"`
	URN         string            `json:"urn,omitempty"`
	Name        string            `json:"name,omitempty"`
	DisplayName string            `json:"displayName,omitempty"`
	Owner       string            `json:"owner,omitempty"`
	Type        string            `json:"type,omitempty"`
	Tags        map[string]string `json:"tags,omitempty"`
}

type IamPermissions struct {
	Allow  []IamAction `json:"allow"`
	Except []IamAction `json:"except"`
	Deny   []IamAction `json:"deny"`
}

func (p IamPermissions) ToLists() ([]string, []string, []string) {
	var allow []string
	var except []string
	var deny []string

	for _, r := range p.Allow {
		allow = append(allow, r.Action)
	}

	for _, r := range p.Except {
		except = append(except, r.Action)
	}

	for _, r := range p.Deny {
		deny = append(deny, r.Action)
	}
	return allow, except, deny
}

type IamAction struct {
	Action string `json:"action"`
}

type IamResourceGroup struct {
	ID        string               `json:"id,omitempty"`
	Name      string               `json:"name"`
	Resources []IamResourceDetails `json:"resources"`
	URN       string               `json:"urn,omitempty"`
	Owner     string               `json:"owner,omitempty"`
	CreatedAt string               `json:"createdAt,omitempty"`
	UpdatedAt string               `json:"updatedAt,omitempty"`
	ReadOnly  bool                 `json:"readOnly,omitempty"`
}

type IamResourceGroupCreate struct {
	Name      string               `json:"name"`
	Resources []IamResourceDetails `json:"resources"`
}

type IamPermissionsGroup struct {
	Id          string         `json:"id,omitempty"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Permissions IamPermissions `json:"permissions"`
	CreatedAt   string         `json:"createdAt,omitempty"`
	UpdatedAt   string         `json:"updatedAt,omitempty"`
	Urn         string         `json:"urn,omitempty"`
	Owner       string         `json:"owner,omitempty"`
}

func (p IamPermissionsGroup) ToMap() map[string]any {
	out := make(map[string]any, 0)
	out["name"] = p.Name

	out["owner"] = p.Owner
	out["created_at"] = p.CreatedAt

	// inline allow, except and deny
	allow, except, deny := p.Permissions.ToLists()
	if len(allow) != 0 {
		out["allow"] = allow
	}
	if len(except) != 0 {
		out["except"] = except
	}
	if len(deny) != 0 {
		out["deny"] = deny
	}

	if p.Description != "" {
		out["description"] = p.Description
	}
	if p.UpdatedAt != "" {
		out["updated_at"] = p.UpdatedAt
	}

	out["urn"] = p.Urn

	return out
}

func AppendIamDatasourceSchema(attrs map[string]schema.Attribute, ctx context.Context) {
	attrs["iam"] = schema.SingleNestedAttribute{
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
	}
}

var _ basetypes.ObjectTypable = IamType{}

type IamType struct {
	basetypes.ObjectType
}

func (t IamType) Equal(o attr.Type) bool {
	other, ok := o.(IamType)

	if !ok {
		return false
	}

	return t.ObjectType.Equal(other.ObjectType)
}

func (t IamType) String() string {
	return "IamType"
}

func (t IamType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics

	attributes := in.Attributes()

	displayNameAttribute, ok := attributes["display_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`display_name is missing from object`)

		return nil, diags
	}

	displayNameVal, ok := displayNameAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`display_name expected to be ovhtypes.TfStringValue, was: %T`, displayNameAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return nil, diags
	}

	idVal, ok := idAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be ovhtypes.TfStringValue, was: %T`, idAttribute))
	}

	tagsAttribute, ok := attributes["tags"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`tags is missing from object`)

		return nil, diags
	}

	tagsVal, ok := tagsAttribute.(ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`tags expected to be ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue], was: %T`, tagsAttribute))
	}

	urnAttribute, ok := attributes["urn"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`urn is missing from object`)

		return nil, diags
	}

	urnVal, ok := urnAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`urn expected to be ovhtypes.TfStringValue, was: %T`, urnAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return IamValue{
		DisplayName: displayNameVal,
		Id:          idVal,
		Tags:        tagsVal,
		Urn:         urnVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func NewIamValueNull() IamValue {
	return IamValue{
		state: attr.ValueStateNull,
	}
}

func NewIamValueUnknown() IamValue {
	return IamValue{
		state: attr.ValueStateUnknown,
	}
}

func NewIamValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (IamValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/521
	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]

		if !ok {
			diags.AddError(
				"Missing IamValue Attribute Value",
				"While creating a IamValue value, a missing attribute value was detected. "+
					"A IamValue must contain values for all attributes, even if null or unknown. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("IamValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)

			continue
		}

		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid IamValue Attribute Type",
				"While creating a IamValue value, an invalid attribute value was detected. "+
					"A IamValue must use a matching attribute type for the value. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("IamValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("IamValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		_, ok := attributeTypes[name]

		if !ok {
			diags.AddError(
				"Extra IamValue Attribute Value",
				"While creating a IamValue value, an extra attribute value was detected. "+
					"A IamValue must not contain values beyond the expected attribute types. "+
					"This is always an issue with the provider and should be reported to the provider developers.\n\n"+
					fmt.Sprintf("Extra IamValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewIamValueUnknown(), diags
	}

	displayNameAttribute, ok := attributes["display_name"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`display_name is missing from object`)

		return NewIamValueUnknown(), diags
	}

	displayNameVal, ok := displayNameAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`display_name expected to be ovhtypes.TfStringValue, was: %T`, displayNameAttribute))
	}

	idAttribute, ok := attributes["id"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`id is missing from object`)

		return NewIamValueUnknown(), diags
	}

	idVal, ok := idAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`id expected to be ovhtypes.TfStringValue, was: %T`, idAttribute))
	}

	tagsAttribute, ok := attributes["tags"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`tags is missing from object`)

		return NewIamValueUnknown(), diags
	}

	tagsVal, ok := tagsAttribute.(ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue])

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`tags expected to be ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue], was: %T`, tagsAttribute))
	}

	urnAttribute, ok := attributes["urn"]

	if !ok {
		diags.AddError(
			"Attribute Missing",
			`urn is missing from object`)

		return NewIamValueUnknown(), diags
	}

	urnVal, ok := urnAttribute.(ovhtypes.TfStringValue)

	if !ok {
		diags.AddError(
			"Attribute Wrong Type",
			fmt.Sprintf(`urn expected to be ovhtypes.TfStringValue, was: %T`, urnAttribute))
	}

	if diags.HasError() {
		return NewIamValueUnknown(), diags
	}

	return IamValue{
		DisplayName: displayNameVal,
		Id:          idVal,
		Tags:        tagsVal,
		Urn:         urnVal,
		state:       attr.ValueStateKnown,
	}, diags
}

func NewIamValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) IamValue {
	object, diags := NewIamValue(attributeTypes, attributes)

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

		panic("NewIamValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}

	return object
}

func (t IamType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewIamValueNull(), nil
	}

	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}

	if !in.IsKnown() {
		return NewIamValueUnknown(), nil
	}

	if in.IsNull() {
		return NewIamValueNull(), nil
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

	return NewIamValueMust(IamValue{}.AttributeTypes(ctx), attributes), nil
}

func (t IamType) ValueType(ctx context.Context) attr.Value {
	return IamValue{}
}

var _ basetypes.ObjectValuable = IamValue{}

type IamValue struct {
	DisplayName ovhtypes.TfStringValue                            `tfsdk:"display_name" json:"displayName"`
	Id          ovhtypes.TfStringValue                            `tfsdk:"id" json:"id"`
	Tags        ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue] `tfsdk:"tags" json:"tags"`
	Urn         ovhtypes.TfStringValue                            `tfsdk:"urn" json:"urn"`
	state       attr.ValueState
}

type IamWritableValue struct {
	*IamValue   `json:"-"`
	DisplayName *ovhtypes.TfStringValue                            `json:"displayName,omitempty"`
	Id          *ovhtypes.TfStringValue                            `json:"id,omitempty"`
	Tags        *ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue] `json:"tags,omitempty"`
	Urn         *ovhtypes.TfStringValue                            `json:"urn,omitempty"`
}

func (v IamValue) ToCreate() *IamWritableValue {
	res := &IamWritableValue{}

	if !v.Id.IsNull() {
		res.Id = &v.Id
	}

	if !v.Tags.IsNull() {
		res.Tags = &v.Tags
	}

	if !v.Urn.IsNull() {
		res.Urn = &v.Urn
	}

	if !v.DisplayName.IsNull() {
		res.DisplayName = &v.DisplayName
	}

	return res
}

func (v *IamValue) UnmarshalJSON(data []byte) error {
	type JsonIamValue IamValue

	var tmp JsonIamValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.DisplayName = tmp.DisplayName
	v.Id = tmp.Id
	v.Tags = tmp.Tags
	v.Urn = tmp.Urn

	v.state = attr.ValueStateKnown

	return nil
}

func (v *IamValue) MergeWith(other *IamValue) {

	if (v.DisplayName.IsUnknown() || v.DisplayName.IsNull()) && !other.DisplayName.IsUnknown() {
		v.DisplayName = other.DisplayName
	}

	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}

	if (v.Tags.IsUnknown() || v.Tags.IsNull()) && !other.Tags.IsUnknown() {
		v.Tags = other.Tags
	}

	if (v.Urn.IsUnknown() || v.Urn.IsNull()) && !other.Urn.IsUnknown() {
		v.Urn = other.Urn
	}

	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v IamValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"displayName": v.DisplayName,
		"id":          v.Id,
		"tags":        v.Tags,
		"urn":         v.Urn,
	}
}
func (v IamValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 4)

	var val tftypes.Value
	var err error

	attrTypes["display_name"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["id"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["tags"] = basetypes.MapType{
		ElemType: types.StringType,
	}.TerraformType(ctx)
	attrTypes["urn"] = basetypes.StringType{}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 4)

		val, err = v.DisplayName.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["display_name"] = val

		val, err = v.Id.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["id"] = val

		val, err = v.Tags.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["tags"] = val

		val, err = v.Urn.ToTerraformValue(ctx)

		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}

		vals["urn"] = val

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

func (v IamValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}

func (v IamValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}

func (v IamValue) String() string {
	return "IamValue"
}

func (v IamValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	objVal, diags := types.ObjectValue(
		map[string]attr.Type{
			"display_name": ovhtypes.TfStringType{},
			"id":           ovhtypes.TfStringType{},
			"tags":         ovhtypes.NewTfMapNestedType[ovhtypes.TfStringValue](ctx),
			"urn":          ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"display_name": v.DisplayName,
			"id":           v.Id,
			"tags":         v.Tags,
			"urn":          v.Urn,
		})

	return objVal, diags
}

func (v IamValue) Equal(o attr.Value) bool {
	other, ok := o.(IamValue)

	if !ok {
		return false
	}

	if v.state != other.state {
		return false
	}

	if v.state != attr.ValueStateKnown {
		return true
	}

	if !v.DisplayName.Equal(other.DisplayName) {
		return false
	}

	if !v.Id.Equal(other.Id) {
		return false
	}

	if !v.Tags.Equal(other.Tags) {
		return false
	}

	if !v.Urn.Equal(other.Urn) {
		return false
	}

	return true
}

func (v IamValue) Type(ctx context.Context) attr.Type {
	return IamType{
		basetypes.ObjectType{
			AttrTypes: v.AttributeTypes(ctx),
		},
	}
}

func (v IamValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"display_name": ovhtypes.TfStringType{},
		"id":           ovhtypes.TfStringType{},
		"tags":         ovhtypes.NewTfMapNestedType[ovhtypes.TfStringValue](ctx),
		"urn":          ovhtypes.TfStringType{},
	}
}
