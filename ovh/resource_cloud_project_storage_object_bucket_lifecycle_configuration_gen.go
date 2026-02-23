package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func CloudProjectStorageLifecycleConfigurationResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Unique identifier for the resource (service_name/region_name/container_name)",
			},
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "The ID of the public cloud project.",
				MarkdownDescription: "The ID of the public cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"region_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region name of the storage container.",
				MarkdownDescription: "Region name of the storage container.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"container_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Name of the storage container.",
				MarkdownDescription: "Name of the storage container.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"rules": schema.ListNestedAttribute{
				CustomType: ovhtypes.NewTfListNestedType[LifecycleRuleValue](ctx),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Computed:            true,
							Description:         "Rule ID.",
							MarkdownDescription: "Rule ID.",
						},
						"status": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Required:            true,
							Description:         "Rule status.",
							MarkdownDescription: "Rule status.",
							Validators: []validator.String{
								stringvalidator.OneOf("enabled", "disabled"),
							},
						},
						"filter": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"prefix": schema.StringAttribute{
									CustomType:          ovhtypes.TfStringType{},
									Optional:            true,
									Description:         "Prefix filter.",
									MarkdownDescription: "Prefix filter.",
								},
								"object_size_greater_than": schema.Int64Attribute{
									CustomType:          ovhtypes.TfInt64Type{},
									Optional:            true,
									Description:         "Minimum object size in bytes to which the rule applies.",
									MarkdownDescription: "Minimum object size in bytes to which the rule applies.",
								},
								"object_size_less_than": schema.Int64Attribute{
									CustomType:          ovhtypes.TfInt64Type{},
									Optional:            true,
									Description:         "Maximum object size in bytes to which the rule applies.",
									MarkdownDescription: "Maximum object size in bytes to which the rule applies.",
								},
								"tags": schema.MapAttribute{
									CustomType:          ovhtypes.NewTfMapNestedType[ovhtypes.TfStringValue](ctx),
									Optional:            true,
									Description:         "Tags filter.",
									MarkdownDescription: "Tags filter.",
								},
							},
							CustomType: LifecycleRuleFilterType{
								ObjectType: types.ObjectType{
									AttrTypes: LifecycleRuleFilterValue{}.AttributeTypes(ctx),
								},
							},
							Optional:            true,
							Description:         "Rule filters.",
							MarkdownDescription: "Rule filters.",
						},
						"expiration": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"days": schema.Int64Attribute{
									CustomType:          ovhtypes.TfInt64Type{},
									Optional:            true,
									Description:         "Objects will be deleted past this lifetime (in days).",
									MarkdownDescription: "Objects will be deleted past this lifetime (in days).",
								},
								"date": schema.StringAttribute{
									CustomType:          ovhtypes.TfStringType{},
									Optional:            true,
									Description:         "Indicates at what date the objects will be deleted (YYYY-MM-DD).",
									MarkdownDescription: "Indicates at what date the objects will be deleted (YYYY-MM-DD).",
									Validators: []validator.String{
										stringvalidator.RegexMatches(
											regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
											"must be a date in YYYY-MM-DD format",
										),
									},
								},
								"expired_object_delete_marker": schema.BoolAttribute{
									CustomType:          ovhtypes.TfBoolType{},
									Optional:            true,
									Description:         "Indicates whether a delete marker with no noncurrent versions will be removed.",
									MarkdownDescription: "Indicates whether a delete marker with no noncurrent versions will be removed. This cannot be specified with Days or Date.",
								},
							},
							CustomType: LifecycleRuleExpirationFilterType{
								ObjectType: types.ObjectType{
									AttrTypes: LifecycleRuleExpirationValue{}.AttributeTypes(ctx),
								},
							},
							Optional:            true,
							Description:         "Lifecycle rule expiration configuration.",
							MarkdownDescription: "Lifecycle rule expiration configuration.",
						},
						"abort_incomplete_multipart_upload": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"days_after_initiation": schema.Int64Attribute{
									CustomType:          ovhtypes.TfInt64Type{},
									Optional:            true,
									Description:         "Specifies the number of days after which an incomplete multipart upload is aborted.",
									MarkdownDescription: "Specifies the number of days after which an incomplete multipart upload is aborted.",
								},
							},
							CustomType: LifecycleRuleAbortIncompleteMultipartUploadType{
								ObjectType: types.ObjectType{
									AttrTypes: LifecycleRuleAbortIncompleteMultipartUploadValue{}.AttributeTypes(ctx),
								},
							},
							Optional:            true,
							Description:         "Abort incomplete multipart upload configuration.",
							MarkdownDescription: "Abort incomplete multipart upload configuration.",
						},
						"noncurrent_version_expiration": schema.SingleNestedAttribute{
							Attributes: map[string]schema.Attribute{
								"noncurrent_days": schema.Int64Attribute{
									CustomType:          ovhtypes.TfInt64Type{},
									Optional:            true,
									Description:         "Specifies the number of days an object is noncurrent before it can be expired.",
									MarkdownDescription: "Specifies the number of days an object is noncurrent before it can be expired.",
								},
								"newer_noncurrent_versions": schema.Int64Attribute{
									CustomType:          ovhtypes.TfInt64Type{},
									Optional:            true,
									Description:         "Specifies how many noncurrent versions to retain.",
									MarkdownDescription: "Specifies how many noncurrent versions to retain.",
								},
							},
							CustomType: LifecycleRuleNoncurrentVersionExpirationType{
								ObjectType: types.ObjectType{
									AttrTypes: LifecycleRuleNoncurrentVersionExpirationValue{}.AttributeTypes(ctx),
								},
							},
							Optional:            true,
							Description:         "Specifies when noncurrent object versions expire.",
							MarkdownDescription: "Specifies when noncurrent object versions expire.",
						},
						"transitions": schema.ListNestedAttribute{
							CustomType: ovhtypes.NewTfListNestedType[LifecycleRuleTransitionValue](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"days": schema.Int64Attribute{
										CustomType:          ovhtypes.TfInt64Type{},
										Optional:            true,
										Description:         "Indicates the number of days after creation when objects are transitioned.",
										MarkdownDescription: "Indicates the number of days after creation when objects are transitioned.",
									},
									"date": schema.StringAttribute{
										CustomType:          ovhtypes.TfStringType{},
										Optional:            true,
										Description:         "Indicates when objects are transitioned to the specified storage class (YYYY-MM-DD).",
										MarkdownDescription: "Indicates when objects are transitioned to the specified storage class (YYYY-MM-DD).",
										Validators: []validator.String{
											stringvalidator.RegexMatches(
												regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`),
												"must be a date in YYYY-MM-DD format",
											),
										},
									},
									"storage_class": schema.StringAttribute{
										CustomType:          ovhtypes.TfStringType{},
										Required:            true,
										Description:         "The storage class to which you want the object to transition.",
										MarkdownDescription: "The storage class to which you want the object to transition.",
										Validators: []validator.String{
											stringvalidator.OneOf(
												"STANDARD",
												"STANDARD_IA",
											),
										},
									},
								},
								CustomType: LifecycleRuleTransitionType{
									ObjectType: types.ObjectType{
										AttrTypes: LifecycleRuleTransitionValue{}.AttributeTypes(ctx),
									},
								},
							},
							Optional:            true,
							Description:         "Specifies when an object transitions to a specified storage class.",
							MarkdownDescription: "Specifies when an object transitions to a specified storage class.",
						},
						"noncurrent_version_transitions": schema.ListNestedAttribute{
							CustomType: ovhtypes.NewTfListNestedType[LifecycleRuleNoncurrentVersionTransitionValue](ctx),
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"noncurrent_days": schema.Int64Attribute{
										CustomType:          ovhtypes.TfInt64Type{},
										Optional:            true,
										Description:         "Specifies the number of days an object is noncurrent before the transition.",
										MarkdownDescription: "Specifies the number of days an object is noncurrent before the transition.",
									},
									"newer_noncurrent_versions": schema.Int64Attribute{
										CustomType:          ovhtypes.TfInt64Type{},
										Optional:            true,
										Description:         "Specifies how many noncurrent versions to retain in the same storage class.",
										MarkdownDescription: "Specifies how many noncurrent versions to retain in the same storage class.",
									},
									"storage_class": schema.StringAttribute{
										CustomType:          ovhtypes.TfStringType{},
										Optional:            true,
										Description:         "The storage class to which you want the noncurrent object to transition.",
										MarkdownDescription: "The storage class to which you want the noncurrent object to transition.",
									},
								},
								CustomType: LifecycleRuleNoncurrentVersionTransitionType{
									ObjectType: types.ObjectType{
										AttrTypes: LifecycleRuleNoncurrentVersionTransitionValue{}.AttributeTypes(ctx),
									},
								},
							},
							Optional:            true,
							Description:         "Specifies the transition rule for noncurrent object versions.",
							MarkdownDescription: "Specifies the transition rule for noncurrent object versions.",
						},
					},
					CustomType: LifecycleRuleType{
						ObjectType: types.ObjectType{
							AttrTypes: LifecycleRuleValue{}.AttributeTypes(ctx),
						},
					},
				},
				Required:            true,
				Description:         "List of lifecycle rules.",
				MarkdownDescription: "List of lifecycle rules.",
			},
		},
		Description:         "Manage S3-compatible storage container lifecycle configuration.",
		MarkdownDescription: "Manage S3-compatible storage container lifecycle configuration for an OVHcloud Object Storage bucket.",
	}
}

// ---------------------------------------------------------------------------
// Top-level model
// ---------------------------------------------------------------------------

type CloudProjectStorageLifecycleConfigurationModel struct {
	ID            ovhtypes.TfStringValue                         `tfsdk:"id" json:"-"`
	ServiceName   ovhtypes.TfStringValue                         `tfsdk:"service_name" json:"-"`
	RegionName    ovhtypes.TfStringValue                         `tfsdk:"region_name" json:"-"`
	ContainerName ovhtypes.TfStringValue                         `tfsdk:"container_name" json:"-"`
	Rules         ovhtypes.TfListNestedValue[LifecycleRuleValue] `tfsdk:"rules" json:"rules"`
}

func (v *CloudProjectStorageLifecycleConfigurationModel) MergeWith(other *CloudProjectStorageLifecycleConfigurationModel) {
	if (v.ID.IsUnknown() || v.ID.IsNull()) && !other.ID.IsUnknown() {
		v.ID = other.ID
	}
	if (v.ServiceName.IsUnknown() || v.ServiceName.IsNull()) && !other.ServiceName.IsUnknown() {
		v.ServiceName = other.ServiceName
	}
	if (v.RegionName.IsUnknown() || v.RegionName.IsNull()) && !other.RegionName.IsUnknown() {
		v.RegionName = other.RegionName
	}
	if (v.ContainerName.IsUnknown() || v.ContainerName.IsNull()) && !other.ContainerName.IsUnknown() {
		v.ContainerName = other.ContainerName
	}
	if (v.Rules.IsUnknown() || v.Rules.IsNull()) && !other.Rules.IsUnknown() {
		v.Rules = other.Rules
	} else if !other.Rules.IsUnknown() && !other.Rules.IsNull() {
		newSlice := make([]attr.Value, 0)
		elems := v.Rules.Elements()
		newElems := other.Rules.Elements()
		if len(elems) != len(newElems) {
			v.Rules = other.Rules
		} else {
			for idx, e := range elems {
				tmp := e.(LifecycleRuleValue)
				tmp2 := newElems[idx].(LifecycleRuleValue)
				tmp.MergeWith(&tmp2)
				newSlice = append(newSlice, tmp)
			}
			v.Rules = ovhtypes.TfListNestedValue[LifecycleRuleValue]{
				ListValue: basetypes.NewListValueMust(LifecycleRuleValue{}.Type(context.Background()), newSlice),
			}
		}
	}
}

type CloudProjectStorageLifecycleConfigurationWritableModel struct {
	Rules *ovhtypes.TfListNestedValue[LifecycleRuleValue] `json:"rules"`
}

func (v CloudProjectStorageLifecycleConfigurationModel) ToCreate() *CloudProjectStorageLifecycleConfigurationWritableModel {
	res := &CloudProjectStorageLifecycleConfigurationWritableModel{}
	if !v.Rules.IsUnknown() {
		res.Rules = &v.Rules
	}
	return res
}

func (v CloudProjectStorageLifecycleConfigurationModel) ToUpdate() *CloudProjectStorageLifecycleConfigurationWritableModel {
	return v.ToCreate()
}

// ---------------------------------------------------------------------------
// LifecycleRuleType / LifecycleRuleValue
// ---------------------------------------------------------------------------

var _ basetypes.ObjectTypable = LifecycleRuleType{}

type LifecycleRuleType struct {
	basetypes.ObjectType
}

func (t LifecycleRuleType) Equal(o attr.Type) bool {
	other, ok := o.(LifecycleRuleType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

func (t LifecycleRuleType) String() string {
	return "LifecycleRuleType"
}

func (t LifecycleRuleType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	var diags diag.Diagnostics
	attributes := in.Attributes()

	idVal, _ := attributes["id"].(ovhtypes.TfStringValue)
	statusVal, _ := attributes["status"].(ovhtypes.TfStringValue)

	filterAttribute, ok := attributes["filter"]
	if !ok {
		diags.AddError("Attribute Missing", "filter is missing from object")
		return nil, diags
	}
	filterVal, ok := filterAttribute.(LifecycleRuleFilterValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf("filter expected to be LifecycleRuleFilterValue, was: %T", filterAttribute))
	}

	expirationAttribute, ok := attributes["expiration"]
	if !ok {
		diags.AddError("Attribute Missing", "expiration is missing from object")
		return nil, diags
	}
	expirationVal, ok := expirationAttribute.(LifecycleRuleExpirationValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf("expiration expected to be LifecycleRuleExpirationValue, was: %T", expirationAttribute))
	}

	abortAttribute, ok := attributes["abort_incomplete_multipart_upload"]
	if !ok {
		diags.AddError("Attribute Missing", "abort_incomplete_multipart_upload is missing from object")
		return nil, diags
	}
	abortVal, ok := abortAttribute.(LifecycleRuleAbortIncompleteMultipartUploadValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf("abort_incomplete_multipart_upload expected to be LifecycleRuleAbortIncompleteMultipartUploadValue, was: %T", abortAttribute))
	}

	nveAttribute, ok := attributes["noncurrent_version_expiration"]
	if !ok {
		diags.AddError("Attribute Missing", "noncurrent_version_expiration is missing from object")
		return nil, diags
	}
	nveVal, ok := nveAttribute.(LifecycleRuleNoncurrentVersionExpirationValue)
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf("noncurrent_version_expiration expected to be LifecycleRuleNoncurrentVersionExpirationValue, was: %T", nveAttribute))
	}

	transitionsAttribute, ok := attributes["transitions"]
	if !ok {
		diags.AddError("Attribute Missing", "transitions is missing from object")
		return nil, diags
	}
	transitionsVal, ok := transitionsAttribute.(ovhtypes.TfListNestedValue[LifecycleRuleTransitionValue])
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf("transitions expected to be TfListNestedValue[LifecycleRuleTransitionValue], was: %T", transitionsAttribute))
	}

	nvtAttribute, ok := attributes["noncurrent_version_transitions"]
	if !ok {
		diags.AddError("Attribute Missing", "noncurrent_version_transitions is missing from object")
		return nil, diags
	}
	nvtVal, ok := nvtAttribute.(ovhtypes.TfListNestedValue[LifecycleRuleNoncurrentVersionTransitionValue])
	if !ok {
		diags.AddError("Attribute Wrong Type", fmt.Sprintf("noncurrent_version_transitions expected to be TfListNestedValue[LifecycleRuleNoncurrentVersionTransitionValue], was: %T", nvtAttribute))
	}

	if diags.HasError() {
		return nil, diags
	}

	return LifecycleRuleValue{
		Id:                             idVal,
		Status:                         statusVal,
		Filter:                         filterVal,
		Expiration:                     expirationVal,
		AbortIncompleteMultipartUpload: abortVal,
		NoncurrentVersionExpiration:    nveVal,
		Transitions:                    transitionsVal,
		NoncurrentVersionTransitions:   nvtVal,
		state:                          attr.ValueStateKnown,
	}, diags
}

func NewLifecycleRuleValueNull() LifecycleRuleValue {
	return LifecycleRuleValue{state: attr.ValueStateNull}
}

func NewLifecycleRuleValueUnknown() LifecycleRuleValue {
	return LifecycleRuleValue{state: attr.ValueStateUnknown}
}

func NewLifecycleRuleValue(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) (LifecycleRuleValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	ctx := context.Background()

	for name, attributeType := range attributeTypes {
		attribute, ok := attributes[name]
		if !ok {
			diags.AddError(
				"Missing LifecycleRuleValue Attribute Value",
				"While creating a LifecycleRuleValue value, a missing attribute value was detected. "+
					fmt.Sprintf("LifecycleRuleValue Attribute Name (%s) Expected Type: %s", name, attributeType.String()),
			)
			continue
		}
		if !attributeType.Equal(attribute.Type(ctx)) {
			diags.AddError(
				"Invalid LifecycleRuleValue Attribute Type",
				fmt.Sprintf("LifecycleRuleValue Attribute Name (%s) Expected Type: %s\n", name, attributeType.String())+
					fmt.Sprintf("LifecycleRuleValue Attribute Name (%s) Given Type: %s", name, attribute.Type(ctx)),
			)
		}
	}

	for name := range attributes {
		if _, ok := attributeTypes[name]; !ok {
			diags.AddError(
				"Extra LifecycleRuleValue Attribute Value",
				fmt.Sprintf("Extra LifecycleRuleValue Attribute Name: %s", name),
			)
		}
	}

	if diags.HasError() {
		return NewLifecycleRuleValueUnknown(), diags
	}

	obj, d := LifecycleRuleType{}.ValueFromObject(ctx, basetypes.NewObjectValueMust(attributeTypes, attributes))
	diags.Append(d...)
	if diags.HasError() {
		return NewLifecycleRuleValueUnknown(), diags
	}

	return obj.(LifecycleRuleValue), diags
}

func NewLifecycleRuleValueMust(attributeTypes map[string]attr.Type, attributes map[string]attr.Value) LifecycleRuleValue {
	object, diags := NewLifecycleRuleValue(attributeTypes, attributes)
	if diags.HasError() {
		diagsStrings := make([]string, 0, len(diags))
		for _, diagnostic := range diags {
			diagsStrings = append(diagsStrings, fmt.Sprintf(
				"%s | %s | %s",
				diagnostic.Severity(),
				diagnostic.Summary(),
				diagnostic.Detail()))
		}
		panic("NewLifecycleRuleValueMust received error(s): " + strings.Join(diagsStrings, "\n"))
	}
	return object
}

func (t LifecycleRuleType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewLifecycleRuleValueNull(), nil
	}
	if !in.Type().Equal(t.TerraformType(ctx)) {
		return nil, fmt.Errorf("expected %s, got %s", t.TerraformType(ctx), in.Type())
	}
	if !in.IsKnown() {
		return NewLifecycleRuleValueUnknown(), nil
	}
	if in.IsNull() {
		return NewLifecycleRuleValueNull(), nil
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
	return NewLifecycleRuleValueMust(LifecycleRuleValue{}.AttributeTypes(ctx), attributes), nil
}

func (t LifecycleRuleType) ValueType(ctx context.Context) attr.Value {
	return LifecycleRuleValue{}
}

var _ basetypes.ObjectValuable = LifecycleRuleValue{}

type LifecycleRuleValue struct {
	Id                             ovhtypes.TfStringValue                                                    `tfsdk:"id" json:"id"`
	Status                         ovhtypes.TfStringValue                                                    `tfsdk:"status" json:"status"`
	Filter                         LifecycleRuleFilterValue                                                  `tfsdk:"filter" json:"filter"`
	Expiration                     LifecycleRuleExpirationValue                                              `tfsdk:"expiration" json:"expiration"`
	AbortIncompleteMultipartUpload LifecycleRuleAbortIncompleteMultipartUploadValue                          `tfsdk:"abort_incomplete_multipart_upload" json:"abortIncompleteMultipartUpload"`
	NoncurrentVersionExpiration    LifecycleRuleNoncurrentVersionExpirationValue                             `tfsdk:"noncurrent_version_expiration" json:"noncurrentVersionExpiration"`
	Transitions                    ovhtypes.TfListNestedValue[LifecycleRuleTransitionValue]                  `tfsdk:"transitions" json:"transitions"`
	NoncurrentVersionTransitions   ovhtypes.TfListNestedValue[LifecycleRuleNoncurrentVersionTransitionValue] `tfsdk:"noncurrent_version_transitions" json:"noncurrentVersionTransitions"`
	state                          attr.ValueState
}

type LifecycleRuleWritableValue struct {
	Id                             *ovhtypes.TfStringValue                                                    `json:"id,omitempty"`
	Status                         *ovhtypes.TfStringValue                                                    `json:"status,omitempty"`
	Filter                         *LifecycleRuleFilterValue                                                  `json:"filter,omitempty"`
	Expiration                     *LifecycleRuleExpirationValue                                              `json:"expiration,omitempty"`
	AbortIncompleteMultipartUpload *LifecycleRuleAbortIncompleteMultipartUploadValue                          `json:"abortIncompleteMultipartUpload,omitempty"`
	NoncurrentVersionExpiration    *LifecycleRuleNoncurrentVersionExpirationValue                             `json:"noncurrentVersionExpiration,omitempty"`
	Transitions                    *ovhtypes.TfListNestedValue[LifecycleRuleTransitionValue]                  `json:"transitions,omitempty"`
	NoncurrentVersionTransitions   *ovhtypes.TfListNestedValue[LifecycleRuleNoncurrentVersionTransitionValue] `json:"noncurrentVersionTransitions,omitempty"`
}

func (v LifecycleRuleValue) ToCreate() *LifecycleRuleWritableValue {
	res := &LifecycleRuleWritableValue{}
	if !v.Id.IsNull() {
		res.Id = &v.Id
	}
	if !v.Status.IsNull() {
		res.Status = &v.Status
	}
	if !v.Filter.IsNull() {
		res.Filter = &v.Filter
	}
	if !v.Expiration.IsNull() {
		res.Expiration = &v.Expiration
	}
	if !v.AbortIncompleteMultipartUpload.IsNull() {
		res.AbortIncompleteMultipartUpload = &v.AbortIncompleteMultipartUpload
	}
	if !v.NoncurrentVersionExpiration.IsNull() {
		res.NoncurrentVersionExpiration = &v.NoncurrentVersionExpiration
	}
	if !v.Transitions.IsNull() {
		res.Transitions = &v.Transitions
	}
	if !v.NoncurrentVersionTransitions.IsNull() {
		res.NoncurrentVersionTransitions = &v.NoncurrentVersionTransitions
	}
	return res
}

func (v LifecycleRuleValue) ToUpdate() *LifecycleRuleWritableValue {
	return v.ToCreate()
}

func (v *LifecycleRuleValue) UnmarshalJSON(data []byte) error {
	type JsonLifecycleRuleValue LifecycleRuleValue
	var tmp JsonLifecycleRuleValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Id = tmp.Id
	v.Status = tmp.Status
	v.Filter = tmp.Filter
	v.Expiration = tmp.Expiration
	v.AbortIncompleteMultipartUpload = tmp.AbortIncompleteMultipartUpload
	v.NoncurrentVersionExpiration = tmp.NoncurrentVersionExpiration
	v.Transitions = tmp.Transitions
	v.NoncurrentVersionTransitions = tmp.NoncurrentVersionTransitions
	v.state = attr.ValueStateKnown
	return nil
}

func (v LifecycleRuleValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.ToCreate())
}

func (v *LifecycleRuleValue) MergeWith(other *LifecycleRuleValue) {
	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}
	if (v.Status.IsUnknown() || v.Status.IsNull()) && !other.Status.IsUnknown() {
		v.Status = other.Status
	}
	if v.Filter.IsUnknown() && !other.Filter.IsUnknown() {
		v.Filter = other.Filter
	} else if !other.Filter.IsUnknown() {
		v.Filter.MergeWith(&other.Filter)
	}
	if v.Expiration.IsUnknown() && !other.Expiration.IsUnknown() {
		v.Expiration = other.Expiration
	} else if !other.Expiration.IsUnknown() {
		v.Expiration.MergeWith(&other.Expiration)
	}
	if v.AbortIncompleteMultipartUpload.IsUnknown() && !other.AbortIncompleteMultipartUpload.IsUnknown() {
		v.AbortIncompleteMultipartUpload = other.AbortIncompleteMultipartUpload
	} else if !other.AbortIncompleteMultipartUpload.IsUnknown() {
		v.AbortIncompleteMultipartUpload.MergeWith(&other.AbortIncompleteMultipartUpload)
	}
	if v.NoncurrentVersionExpiration.IsUnknown() && !other.NoncurrentVersionExpiration.IsUnknown() {
		v.NoncurrentVersionExpiration = other.NoncurrentVersionExpiration
	} else if !other.NoncurrentVersionExpiration.IsUnknown() {
		v.NoncurrentVersionExpiration.MergeWith(&other.NoncurrentVersionExpiration)
	}
	if (v.Transitions.IsUnknown() || v.Transitions.IsNull()) && !other.Transitions.IsUnknown() {
		v.Transitions = other.Transitions
	}
	if (v.NoncurrentVersionTransitions.IsUnknown() || v.NoncurrentVersionTransitions.IsNull()) && !other.NoncurrentVersionTransitions.IsUnknown() {
		v.NoncurrentVersionTransitions = other.NoncurrentVersionTransitions
	}
	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v LifecycleRuleValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"id":                                v.Id,
		"status":                            v.Status,
		"filter":                            v.Filter,
		"expiration":                        v.Expiration,
		"abort_incomplete_multipart_upload": v.AbortIncompleteMultipartUpload,
		"noncurrent_version_expiration":     v.NoncurrentVersionExpiration,
		"transitions":                       v.Transitions,
		"noncurrent_version_transitions":    v.NoncurrentVersionTransitions,
	}
}

func (v LifecycleRuleValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := make(map[string]tftypes.Type, 8)
	attrTypes["id"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["status"] = basetypes.StringType{}.TerraformType(ctx)
	attrTypes["filter"] = basetypes.ObjectType{AttrTypes: LifecycleRuleFilterValue{}.AttributeTypes(ctx)}.TerraformType(ctx)
	attrTypes["expiration"] = basetypes.ObjectType{AttrTypes: LifecycleRuleExpirationValue{}.AttributeTypes(ctx)}.TerraformType(ctx)
	attrTypes["abort_incomplete_multipart_upload"] = basetypes.ObjectType{AttrTypes: LifecycleRuleAbortIncompleteMultipartUploadValue{}.AttributeTypes(ctx)}.TerraformType(ctx)
	attrTypes["noncurrent_version_expiration"] = basetypes.ObjectType{AttrTypes: LifecycleRuleNoncurrentVersionExpirationValue{}.AttributeTypes(ctx)}.TerraformType(ctx)
	attrTypes["transitions"] = basetypes.ListType{ElemType: LifecycleRuleTransitionValue{}.Type(ctx)}.TerraformType(ctx)
	attrTypes["noncurrent_version_transitions"] = basetypes.ListType{ElemType: LifecycleRuleNoncurrentVersionTransitionValue{}.Type(ctx)}.TerraformType(ctx)

	objectType := tftypes.Object{AttributeTypes: attrTypes}

	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 8)
		var val tftypes.Value
		var err error

		val, err = v.Id.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["id"] = val

		val, err = v.Status.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["status"] = val

		val, err = v.Filter.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["filter"] = val

		val, err = v.Expiration.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["expiration"] = val

		val, err = v.AbortIncompleteMultipartUpload.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["abort_incomplete_multipart_upload"] = val

		val, err = v.NoncurrentVersionExpiration.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["noncurrent_version_expiration"] = val

		val, err = v.Transitions.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["transitions"] = val

		val, err = v.NoncurrentVersionTransitions.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["noncurrent_version_transitions"] = val

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

func (v LifecycleRuleValue) IsNull() bool    { return v.state == attr.ValueStateNull }
func (v LifecycleRuleValue) IsUnknown() bool { return v.state == attr.ValueStateUnknown }
func (v LifecycleRuleValue) String() string  { return "LifecycleRuleValue" }

func (v LifecycleRuleValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	filterType := LifecycleRuleFilterType{ObjectType: types.ObjectType{AttrTypes: LifecycleRuleFilterValue{}.AttributeTypes(ctx)}}
	expirationFilterType := LifecycleRuleExpirationFilterType{ObjectType: types.ObjectType{AttrTypes: LifecycleRuleExpirationValue{}.AttributeTypes(ctx)}}
	abortType := LifecycleRuleAbortIncompleteMultipartUploadType{ObjectType: types.ObjectType{AttrTypes: LifecycleRuleAbortIncompleteMultipartUploadValue{}.AttributeTypes(ctx)}}
	nveType := LifecycleRuleNoncurrentVersionExpirationType{ObjectType: types.ObjectType{AttrTypes: LifecycleRuleNoncurrentVersionExpirationValue{}.AttributeTypes(ctx)}}

	return types.ObjectValue(
		map[string]attr.Type{
			"id":                                ovhtypes.TfStringType{},
			"status":                            ovhtypes.TfStringType{},
			"filter":                            filterType,
			"expiration":                        expirationFilterType,
			"abort_incomplete_multipart_upload": abortType,
			"noncurrent_version_expiration":     nveType,
			"transitions":                       ovhtypes.NewTfListNestedType[LifecycleRuleTransitionValue](ctx),
			"noncurrent_version_transitions":    ovhtypes.NewTfListNestedType[LifecycleRuleNoncurrentVersionTransitionValue](ctx),
		},
		map[string]attr.Value{
			"id":                                v.Id,
			"status":                            v.Status,
			"filter":                            v.Filter,
			"expiration":                        v.Expiration,
			"abort_incomplete_multipart_upload": v.AbortIncompleteMultipartUpload,
			"noncurrent_version_expiration":     v.NoncurrentVersionExpiration,
			"transitions":                       v.Transitions,
			"noncurrent_version_transitions":    v.NoncurrentVersionTransitions,
		},
	)
}

func (v LifecycleRuleValue) Equal(o attr.Value) bool {
	other, ok := o.(LifecycleRuleValue)
	if !ok {
		return false
	}
	if v.state != other.state {
		return false
	}
	if v.state != attr.ValueStateKnown {
		return true
	}
	return v.Id.Equal(other.Id) &&
		v.Status.Equal(other.Status) &&
		v.Filter.Equal(other.Filter) &&
		v.Expiration.Equal(other.Expiration) &&
		v.AbortIncompleteMultipartUpload.Equal(other.AbortIncompleteMultipartUpload) &&
		v.NoncurrentVersionExpiration.Equal(other.NoncurrentVersionExpiration) &&
		v.Transitions.Equal(other.Transitions) &&
		v.NoncurrentVersionTransitions.Equal(other.NoncurrentVersionTransitions)
}

func (v LifecycleRuleValue) Type(ctx context.Context) attr.Type {
	return LifecycleRuleType{basetypes.ObjectType{AttrTypes: v.AttributeTypes(ctx)}}
}

func (v LifecycleRuleValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"id":                                ovhtypes.TfStringType{},
		"status":                            ovhtypes.TfStringType{},
		"filter":                            LifecycleRuleFilterType{ObjectType: types.ObjectType{AttrTypes: LifecycleRuleFilterValue{}.AttributeTypes(ctx)}},
		"expiration":                        LifecycleRuleExpirationFilterType{ObjectType: types.ObjectType{AttrTypes: LifecycleRuleExpirationValue{}.AttributeTypes(ctx)}},
		"abort_incomplete_multipart_upload": LifecycleRuleAbortIncompleteMultipartUploadType{ObjectType: types.ObjectType{AttrTypes: LifecycleRuleAbortIncompleteMultipartUploadValue{}.AttributeTypes(ctx)}},
		"noncurrent_version_expiration":     LifecycleRuleNoncurrentVersionExpirationType{ObjectType: types.ObjectType{AttrTypes: LifecycleRuleNoncurrentVersionExpirationValue{}.AttributeTypes(ctx)}},
		"transitions":                       ovhtypes.NewTfListNestedType[LifecycleRuleTransitionValue](ctx),
		"noncurrent_version_transitions":    ovhtypes.NewTfListNestedType[LifecycleRuleNoncurrentVersionTransitionValue](ctx),
	}
}

// ---------------------------------------------------------------------------
// LifecycleRuleFilterType / LifecycleRuleFilterValue
// ---------------------------------------------------------------------------

var _ basetypes.ObjectTypable = LifecycleRuleFilterType{}

type LifecycleRuleFilterType struct {
	basetypes.ObjectType
}

func (t LifecycleRuleFilterType) Equal(o attr.Type) bool {
	other, ok := o.(LifecycleRuleFilterType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

func (t LifecycleRuleFilterType) String() string { return "LifecycleRuleFilterType" }

func (t LifecycleRuleFilterType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	attributes := in.Attributes()

	prefixVal, _ := attributes["prefix"].(ovhtypes.TfStringValue)
	sizeGTVal, _ := attributes["object_size_greater_than"].(ovhtypes.TfInt64Value)
	sizeLTVal, _ := attributes["object_size_less_than"].(ovhtypes.TfInt64Value)
	tagsVal, _ := attributes["tags"].(ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue])

	return LifecycleRuleFilterValue{
		Prefix:                prefixVal,
		ObjectSizeGreaterThan: sizeGTVal,
		ObjectSizeLessThan:    sizeLTVal,
		Tags:                  tagsVal,
		state:                 attr.ValueStateKnown,
	}, nil
}

func NewLifecycleRuleFilterValueNull() LifecycleRuleFilterValue {
	return LifecycleRuleFilterValue{state: attr.ValueStateNull}
}

func NewLifecycleRuleFilterValueUnknown() LifecycleRuleFilterValue {
	return LifecycleRuleFilterValue{state: attr.ValueStateUnknown}
}

func (t LifecycleRuleFilterType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewLifecycleRuleFilterValueNull(), nil
	}
	if !in.IsKnown() {
		return NewLifecycleRuleFilterValueUnknown(), nil
	}
	if in.IsNull() {
		return NewLifecycleRuleFilterValueNull(), nil
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
	obj, diags := t.ValueFromObject(ctx, basetypes.NewObjectValueMust(LifecycleRuleFilterValue{}.AttributeTypes(ctx), attributes))
	if diags.HasError() {
		return nil, fmt.Errorf("error converting object: %v", diags)
	}
	return obj, nil
}

func (t LifecycleRuleFilterType) ValueType(ctx context.Context) attr.Value {
	return LifecycleRuleFilterValue{}
}

var _ basetypes.ObjectValuable = LifecycleRuleFilterValue{}

type LifecycleRuleFilterValue struct {
	Prefix                ovhtypes.TfStringValue                            `tfsdk:"prefix" json:"prefix"`
	ObjectSizeGreaterThan ovhtypes.TfInt64Value                             `tfsdk:"object_size_greater_than" json:"objectSizeGreaterThan"`
	ObjectSizeLessThan    ovhtypes.TfInt64Value                             `tfsdk:"object_size_less_than" json:"objectSizeLessThan"`
	Tags                  ovhtypes.TfMapNestedValue[ovhtypes.TfStringValue] `tfsdk:"tags" json:"tags"`
	state                 attr.ValueState
}

func (v *LifecycleRuleFilterValue) UnmarshalJSON(data []byte) error {
	type JsonLifecycleRuleFilterValue LifecycleRuleFilterValue
	var tmp JsonLifecycleRuleFilterValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Prefix = tmp.Prefix
	v.ObjectSizeGreaterThan = tmp.ObjectSizeGreaterThan
	v.ObjectSizeLessThan = tmp.ObjectSizeLessThan
	v.Tags = tmp.Tags
	v.state = attr.ValueStateKnown
	return nil
}

func (v *LifecycleRuleFilterValue) MergeWith(other *LifecycleRuleFilterValue) {
	if (v.Prefix.IsUnknown() || v.Prefix.IsNull()) && !other.Prefix.IsUnknown() {
		v.Prefix = other.Prefix
	}
	if (v.ObjectSizeGreaterThan.IsUnknown() || v.ObjectSizeGreaterThan.IsNull()) && !other.ObjectSizeGreaterThan.IsUnknown() {
		v.ObjectSizeGreaterThan = other.ObjectSizeGreaterThan
	}
	if (v.ObjectSizeLessThan.IsUnknown() || v.ObjectSizeLessThan.IsNull()) && !other.ObjectSizeLessThan.IsUnknown() {
		v.ObjectSizeLessThan = other.ObjectSizeLessThan
	}
	if (v.Tags.IsUnknown() || v.Tags.IsNull()) && !other.Tags.IsUnknown() {
		v.Tags = other.Tags
	}
	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v LifecycleRuleFilterValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"prefix":                   v.Prefix,
		"object_size_greater_than": v.ObjectSizeGreaterThan,
		"object_size_less_than":    v.ObjectSizeLessThan,
		"tags":                     v.Tags,
	}
}

func (v LifecycleRuleFilterValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := map[string]tftypes.Type{
		"prefix":                   basetypes.StringType{}.TerraformType(ctx),
		"object_size_greater_than": basetypes.NumberType{}.TerraformType(ctx),
		"object_size_less_than":    basetypes.NumberType{}.TerraformType(ctx),
		"tags":                     basetypes.MapType{ElemType: basetypes.StringType{}}.TerraformType(ctx),
	}
	objectType := tftypes.Object{AttributeTypes: attrTypes}
	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 4)
		var val tftypes.Value
		var err error

		val, err = v.Prefix.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["prefix"] = val

		val, err = v.ObjectSizeGreaterThan.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["object_size_greater_than"] = val

		val, err = v.ObjectSizeLessThan.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["object_size_less_than"] = val

		val, err = v.Tags.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["tags"] = val

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

func (v LifecycleRuleFilterValue) IsNull() bool    { return v.state == attr.ValueStateNull }
func (v LifecycleRuleFilterValue) IsUnknown() bool { return v.state == attr.ValueStateUnknown }
func (v LifecycleRuleFilterValue) String() string  { return "LifecycleRuleFilterValue" }

func (v LifecycleRuleFilterValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(
		map[string]attr.Type{
			"prefix":                   ovhtypes.TfStringType{},
			"object_size_greater_than": ovhtypes.TfInt64Type{},
			"object_size_less_than":    ovhtypes.TfInt64Type{},
			"tags":                     ovhtypes.NewTfMapNestedType[ovhtypes.TfStringValue](ctx),
		},
		map[string]attr.Value{
			"prefix":                   v.Prefix,
			"object_size_greater_than": v.ObjectSizeGreaterThan,
			"object_size_less_than":    v.ObjectSizeLessThan,
			"tags":                     v.Tags,
		},
	)
}

func (v LifecycleRuleFilterValue) Equal(o attr.Value) bool {
	other, ok := o.(LifecycleRuleFilterValue)
	if !ok {
		return false
	}
	if v.state != other.state {
		return false
	}
	if v.state != attr.ValueStateKnown {
		return true
	}
	return v.Prefix.Equal(other.Prefix) &&
		v.ObjectSizeGreaterThan.Equal(other.ObjectSizeGreaterThan) &&
		v.ObjectSizeLessThan.Equal(other.ObjectSizeLessThan) &&
		v.Tags.Equal(other.Tags)
}

func (v LifecycleRuleFilterValue) Type(ctx context.Context) attr.Type {
	return LifecycleRuleFilterType{basetypes.ObjectType{AttrTypes: v.AttributeTypes(ctx)}}
}

func (v LifecycleRuleFilterValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"prefix":                   ovhtypes.TfStringType{},
		"object_size_greater_than": ovhtypes.TfInt64Type{},
		"object_size_less_than":    ovhtypes.TfInt64Type{},
		"tags":                     ovhtypes.NewTfMapNestedType[ovhtypes.TfStringValue](ctx),
	}
}

// ---------------------------------------------------------------------------
// LifecycleRuleExpirationFilterType / LifecycleRuleExpirationValue
// ---------------------------------------------------------------------------

var _ basetypes.ObjectTypable = LifecycleRuleExpirationFilterType{}

type LifecycleRuleExpirationFilterType struct {
	basetypes.ObjectType
}

func (t LifecycleRuleExpirationFilterType) Equal(o attr.Type) bool {
	other, ok := o.(LifecycleRuleExpirationFilterType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

func (t LifecycleRuleExpirationFilterType) String() string {
	return "LifecycleRuleExpirationFilterType"
}

func (t LifecycleRuleExpirationFilterType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	attributes := in.Attributes()
	daysVal, _ := attributes["days"].(ovhtypes.TfInt64Value)
	dateVal, _ := attributes["date"].(ovhtypes.TfStringValue)
	markerVal, _ := attributes["expired_object_delete_marker"].(ovhtypes.TfBoolValue)

	return LifecycleRuleExpirationValue{
		Days:                      daysVal,
		Date:                      dateVal,
		ExpiredObjectDeleteMarker: markerVal,
		state:                     attr.ValueStateKnown,
	}, nil
}

func NewLifecycleRuleExpirationValueNull() LifecycleRuleExpirationValue {
	return LifecycleRuleExpirationValue{state: attr.ValueStateNull}
}

func NewLifecycleRuleExpirationValueUnknown() LifecycleRuleExpirationValue {
	return LifecycleRuleExpirationValue{state: attr.ValueStateUnknown}
}

func (t LifecycleRuleExpirationFilterType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewLifecycleRuleExpirationValueNull(), nil
	}
	if !in.IsKnown() {
		return NewLifecycleRuleExpirationValueUnknown(), nil
	}
	if in.IsNull() {
		return NewLifecycleRuleExpirationValueNull(), nil
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
	obj, diags := t.ValueFromObject(ctx, basetypes.NewObjectValueMust(LifecycleRuleExpirationValue{}.AttributeTypes(ctx), attributes))
	if diags.HasError() {
		return nil, fmt.Errorf("error converting object: %v", diags)
	}
	return obj, nil
}

func (t LifecycleRuleExpirationFilterType) ValueType(ctx context.Context) attr.Value {
	return LifecycleRuleExpirationValue{}
}

var _ basetypes.ObjectValuable = LifecycleRuleExpirationValue{}

type LifecycleRuleExpirationValue struct {
	Days                      ovhtypes.TfInt64Value  `tfsdk:"days" json:"days"`
	Date                      ovhtypes.TfStringValue `tfsdk:"date" json:"date"`
	ExpiredObjectDeleteMarker ovhtypes.TfBoolValue   `tfsdk:"expired_object_delete_marker" json:"expiredObjectDeleteMarker"`
	state                     attr.ValueState
}

func (v *LifecycleRuleExpirationValue) UnmarshalJSON(data []byte) error {
	type JsonLifecycleRuleExpirationValue LifecycleRuleExpirationValue
	var tmp JsonLifecycleRuleExpirationValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Days = tmp.Days
	v.Date = tmp.Date
	v.ExpiredObjectDeleteMarker = tmp.ExpiredObjectDeleteMarker
	v.state = attr.ValueStateKnown
	return nil
}

func (v *LifecycleRuleExpirationValue) MergeWith(other *LifecycleRuleExpirationValue) {
	if (v.Days.IsUnknown() || v.Days.IsNull()) && !other.Days.IsUnknown() {
		v.Days = other.Days
	}
	if (v.Date.IsUnknown() || v.Date.IsNull()) && !other.Date.IsUnknown() {
		v.Date = other.Date
	}
	if (v.ExpiredObjectDeleteMarker.IsUnknown() || v.ExpiredObjectDeleteMarker.IsNull()) && !other.ExpiredObjectDeleteMarker.IsUnknown() {
		v.ExpiredObjectDeleteMarker = other.ExpiredObjectDeleteMarker
	}
	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v LifecycleRuleExpirationValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"days":                         v.Days,
		"date":                         v.Date,
		"expired_object_delete_marker": v.ExpiredObjectDeleteMarker,
	}
}

func (v LifecycleRuleExpirationValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := map[string]tftypes.Type{
		"days":                         basetypes.NumberType{}.TerraformType(ctx),
		"date":                         basetypes.StringType{}.TerraformType(ctx),
		"expired_object_delete_marker": basetypes.BoolType{}.TerraformType(ctx),
	}
	objectType := tftypes.Object{AttributeTypes: attrTypes}
	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)
		var val tftypes.Value
		var err error

		val, err = v.Days.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["days"] = val

		val, err = v.Date.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["date"] = val

		val, err = v.ExpiredObjectDeleteMarker.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["expired_object_delete_marker"] = val

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

func (v LifecycleRuleExpirationValue) IsNull() bool    { return v.state == attr.ValueStateNull }
func (v LifecycleRuleExpirationValue) IsUnknown() bool { return v.state == attr.ValueStateUnknown }
func (v LifecycleRuleExpirationValue) String() string  { return "LifecycleRuleExpirationValue" }

func (v LifecycleRuleExpirationValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(
		map[string]attr.Type{
			"days":                         ovhtypes.TfInt64Type{},
			"date":                         ovhtypes.TfStringType{},
			"expired_object_delete_marker": ovhtypes.TfBoolType{},
		},
		map[string]attr.Value{
			"days":                         v.Days,
			"date":                         v.Date,
			"expired_object_delete_marker": v.ExpiredObjectDeleteMarker,
		},
	)
}

func (v LifecycleRuleExpirationValue) Equal(o attr.Value) bool {
	other, ok := o.(LifecycleRuleExpirationValue)
	if !ok {
		return false
	}
	if v.state != other.state {
		return false
	}
	if v.state != attr.ValueStateKnown {
		return true
	}
	return v.Days.Equal(other.Days) && v.Date.Equal(other.Date) && v.ExpiredObjectDeleteMarker.Equal(other.ExpiredObjectDeleteMarker)
}

func (v LifecycleRuleExpirationValue) Type(ctx context.Context) attr.Type {
	return LifecycleRuleExpirationFilterType{basetypes.ObjectType{AttrTypes: v.AttributeTypes(ctx)}}
}

func (v LifecycleRuleExpirationValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"days":                         ovhtypes.TfInt64Type{},
		"date":                         ovhtypes.TfStringType{},
		"expired_object_delete_marker": ovhtypes.TfBoolType{},
	}
}

// ---------------------------------------------------------------------------
// LifecycleRuleAbortIncompleteMultipartUploadType / Value
// ---------------------------------------------------------------------------

var _ basetypes.ObjectTypable = LifecycleRuleAbortIncompleteMultipartUploadType{}

type LifecycleRuleAbortIncompleteMultipartUploadType struct {
	basetypes.ObjectType
}

func (t LifecycleRuleAbortIncompleteMultipartUploadType) Equal(o attr.Type) bool {
	other, ok := o.(LifecycleRuleAbortIncompleteMultipartUploadType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

func (t LifecycleRuleAbortIncompleteMultipartUploadType) String() string {
	return "LifecycleRuleAbortIncompleteMultipartUploadType"
}

func (t LifecycleRuleAbortIncompleteMultipartUploadType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	attributes := in.Attributes()
	daysVal, _ := attributes["days_after_initiation"].(ovhtypes.TfInt64Value)
	return LifecycleRuleAbortIncompleteMultipartUploadValue{
		DaysAfterInitiation: daysVal,
		state:               attr.ValueStateKnown,
	}, nil
}

func NewLifecycleRuleAbortIncompleteMultipartUploadValueNull() LifecycleRuleAbortIncompleteMultipartUploadValue {
	return LifecycleRuleAbortIncompleteMultipartUploadValue{state: attr.ValueStateNull}
}

func NewLifecycleRuleAbortIncompleteMultipartUploadValueUnknown() LifecycleRuleAbortIncompleteMultipartUploadValue {
	return LifecycleRuleAbortIncompleteMultipartUploadValue{state: attr.ValueStateUnknown}
}

func (t LifecycleRuleAbortIncompleteMultipartUploadType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewLifecycleRuleAbortIncompleteMultipartUploadValueNull(), nil
	}
	if !in.IsKnown() {
		return NewLifecycleRuleAbortIncompleteMultipartUploadValueUnknown(), nil
	}
	if in.IsNull() {
		return NewLifecycleRuleAbortIncompleteMultipartUploadValueNull(), nil
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
	obj, diags := t.ValueFromObject(ctx, basetypes.NewObjectValueMust(LifecycleRuleAbortIncompleteMultipartUploadValue{}.AttributeTypes(ctx), attributes))
	if diags.HasError() {
		return nil, fmt.Errorf("error converting object: %v", diags)
	}
	return obj, nil
}

func (t LifecycleRuleAbortIncompleteMultipartUploadType) ValueType(ctx context.Context) attr.Value {
	return LifecycleRuleAbortIncompleteMultipartUploadValue{}
}

var _ basetypes.ObjectValuable = LifecycleRuleAbortIncompleteMultipartUploadValue{}

type LifecycleRuleAbortIncompleteMultipartUploadValue struct {
	DaysAfterInitiation ovhtypes.TfInt64Value `tfsdk:"days_after_initiation" json:"daysAfterInitiation"`
	state               attr.ValueState
}

func (v *LifecycleRuleAbortIncompleteMultipartUploadValue) UnmarshalJSON(data []byte) error {
	type JsonValue LifecycleRuleAbortIncompleteMultipartUploadValue
	var tmp JsonValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.DaysAfterInitiation = tmp.DaysAfterInitiation
	v.state = attr.ValueStateKnown
	return nil
}

func (v *LifecycleRuleAbortIncompleteMultipartUploadValue) MergeWith(other *LifecycleRuleAbortIncompleteMultipartUploadValue) {
	if (v.DaysAfterInitiation.IsUnknown() || v.DaysAfterInitiation.IsNull()) && !other.DaysAfterInitiation.IsUnknown() {
		v.DaysAfterInitiation = other.DaysAfterInitiation
	}
	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v LifecycleRuleAbortIncompleteMultipartUploadValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{"days_after_initiation": v.DaysAfterInitiation}
}

func (v LifecycleRuleAbortIncompleteMultipartUploadValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := map[string]tftypes.Type{
		"days_after_initiation": basetypes.NumberType{}.TerraformType(ctx),
	}
	objectType := tftypes.Object{AttributeTypes: attrTypes}
	switch v.state {
	case attr.ValueStateKnown:
		val, err := v.DaysAfterInitiation.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals := map[string]tftypes.Value{"days_after_initiation": val}
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

func (v LifecycleRuleAbortIncompleteMultipartUploadValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}
func (v LifecycleRuleAbortIncompleteMultipartUploadValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}
func (v LifecycleRuleAbortIncompleteMultipartUploadValue) String() string {
	return "LifecycleRuleAbortIncompleteMultipartUploadValue"
}

func (v LifecycleRuleAbortIncompleteMultipartUploadValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(
		map[string]attr.Type{"days_after_initiation": ovhtypes.TfInt64Type{}},
		map[string]attr.Value{"days_after_initiation": v.DaysAfterInitiation},
	)
}

func (v LifecycleRuleAbortIncompleteMultipartUploadValue) Equal(o attr.Value) bool {
	other, ok := o.(LifecycleRuleAbortIncompleteMultipartUploadValue)
	if !ok {
		return false
	}
	if v.state != other.state {
		return false
	}
	if v.state != attr.ValueStateKnown {
		return true
	}
	return v.DaysAfterInitiation.Equal(other.DaysAfterInitiation)
}

func (v LifecycleRuleAbortIncompleteMultipartUploadValue) Type(ctx context.Context) attr.Type {
	return LifecycleRuleAbortIncompleteMultipartUploadType{basetypes.ObjectType{AttrTypes: v.AttributeTypes(ctx)}}
}

func (v LifecycleRuleAbortIncompleteMultipartUploadValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{"days_after_initiation": ovhtypes.TfInt64Type{}}
}

// ---------------------------------------------------------------------------
// LifecycleRuleNoncurrentVersionExpirationType / Value
// ---------------------------------------------------------------------------

var _ basetypes.ObjectTypable = LifecycleRuleNoncurrentVersionExpirationType{}

type LifecycleRuleNoncurrentVersionExpirationType struct {
	basetypes.ObjectType
}

func (t LifecycleRuleNoncurrentVersionExpirationType) Equal(o attr.Type) bool {
	other, ok := o.(LifecycleRuleNoncurrentVersionExpirationType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

func (t LifecycleRuleNoncurrentVersionExpirationType) String() string {
	return "LifecycleRuleNoncurrentVersionExpirationType"
}

func (t LifecycleRuleNoncurrentVersionExpirationType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	attributes := in.Attributes()
	noncurrentDaysVal, _ := attributes["noncurrent_days"].(ovhtypes.TfInt64Value)
	newerVersionsVal, _ := attributes["newer_noncurrent_versions"].(ovhtypes.TfInt64Value)
	return LifecycleRuleNoncurrentVersionExpirationValue{
		NoncurrentDays:          noncurrentDaysVal,
		NewerNoncurrentVersions: newerVersionsVal,
		state:                   attr.ValueStateKnown,
	}, nil
}

func NewLifecycleRuleNoncurrentVersionExpirationValueNull() LifecycleRuleNoncurrentVersionExpirationValue {
	return LifecycleRuleNoncurrentVersionExpirationValue{state: attr.ValueStateNull}
}

func NewLifecycleRuleNoncurrentVersionExpirationValueUnknown() LifecycleRuleNoncurrentVersionExpirationValue {
	return LifecycleRuleNoncurrentVersionExpirationValue{state: attr.ValueStateUnknown}
}

func (t LifecycleRuleNoncurrentVersionExpirationType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewLifecycleRuleNoncurrentVersionExpirationValueNull(), nil
	}
	if !in.IsKnown() {
		return NewLifecycleRuleNoncurrentVersionExpirationValueUnknown(), nil
	}
	if in.IsNull() {
		return NewLifecycleRuleNoncurrentVersionExpirationValueNull(), nil
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
	obj, diags := t.ValueFromObject(ctx, basetypes.NewObjectValueMust(LifecycleRuleNoncurrentVersionExpirationValue{}.AttributeTypes(ctx), attributes))
	if diags.HasError() {
		return nil, fmt.Errorf("error converting object: %v", diags)
	}
	return obj, nil
}

func (t LifecycleRuleNoncurrentVersionExpirationType) ValueType(ctx context.Context) attr.Value {
	return LifecycleRuleNoncurrentVersionExpirationValue{}
}

var _ basetypes.ObjectValuable = LifecycleRuleNoncurrentVersionExpirationValue{}

type LifecycleRuleNoncurrentVersionExpirationValue struct {
	NoncurrentDays          ovhtypes.TfInt64Value `tfsdk:"noncurrent_days" json:"noncurrentDays"`
	NewerNoncurrentVersions ovhtypes.TfInt64Value `tfsdk:"newer_noncurrent_versions" json:"newerNoncurrentVersions"`
	state                   attr.ValueState
}

func (v *LifecycleRuleNoncurrentVersionExpirationValue) UnmarshalJSON(data []byte) error {
	type JsonValue LifecycleRuleNoncurrentVersionExpirationValue
	var tmp JsonValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.NoncurrentDays = tmp.NoncurrentDays
	v.NewerNoncurrentVersions = tmp.NewerNoncurrentVersions
	v.state = attr.ValueStateKnown
	return nil
}

func (v *LifecycleRuleNoncurrentVersionExpirationValue) MergeWith(other *LifecycleRuleNoncurrentVersionExpirationValue) {
	if (v.NoncurrentDays.IsUnknown() || v.NoncurrentDays.IsNull()) && !other.NoncurrentDays.IsUnknown() {
		v.NoncurrentDays = other.NoncurrentDays
	}
	if (v.NewerNoncurrentVersions.IsUnknown() || v.NewerNoncurrentVersions.IsNull()) && !other.NewerNoncurrentVersions.IsUnknown() {
		v.NewerNoncurrentVersions = other.NewerNoncurrentVersions
	}
	if (v.state == attr.ValueStateUnknown || v.state == attr.ValueStateNull) && other.state != attr.ValueStateUnknown {
		v.state = other.state
	}
}

func (v LifecycleRuleNoncurrentVersionExpirationValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"noncurrent_days":           v.NoncurrentDays,
		"newer_noncurrent_versions": v.NewerNoncurrentVersions,
	}
}

func (v LifecycleRuleNoncurrentVersionExpirationValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := map[string]tftypes.Type{
		"noncurrent_days":           basetypes.NumberType{}.TerraformType(ctx),
		"newer_noncurrent_versions": basetypes.NumberType{}.TerraformType(ctx),
	}
	objectType := tftypes.Object{AttributeTypes: attrTypes}
	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 2)
		var val tftypes.Value
		var err error

		val, err = v.NoncurrentDays.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["noncurrent_days"] = val

		val, err = v.NewerNoncurrentVersions.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["newer_noncurrent_versions"] = val

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

func (v LifecycleRuleNoncurrentVersionExpirationValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}
func (v LifecycleRuleNoncurrentVersionExpirationValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}
func (v LifecycleRuleNoncurrentVersionExpirationValue) String() string {
	return "LifecycleRuleNoncurrentVersionExpirationValue"
}

func (v LifecycleRuleNoncurrentVersionExpirationValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(
		map[string]attr.Type{
			"noncurrent_days":           ovhtypes.TfInt64Type{},
			"newer_noncurrent_versions": ovhtypes.TfInt64Type{},
		},
		map[string]attr.Value{
			"noncurrent_days":           v.NoncurrentDays,
			"newer_noncurrent_versions": v.NewerNoncurrentVersions,
		},
	)
}

func (v LifecycleRuleNoncurrentVersionExpirationValue) Equal(o attr.Value) bool {
	other, ok := o.(LifecycleRuleNoncurrentVersionExpirationValue)
	if !ok {
		return false
	}
	if v.state != other.state {
		return false
	}
	if v.state != attr.ValueStateKnown {
		return true
	}
	return v.NoncurrentDays.Equal(other.NoncurrentDays) && v.NewerNoncurrentVersions.Equal(other.NewerNoncurrentVersions)
}

func (v LifecycleRuleNoncurrentVersionExpirationValue) Type(ctx context.Context) attr.Type {
	return LifecycleRuleNoncurrentVersionExpirationType{basetypes.ObjectType{AttrTypes: v.AttributeTypes(ctx)}}
}

func (v LifecycleRuleNoncurrentVersionExpirationValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"noncurrent_days":           ovhtypes.TfInt64Type{},
		"newer_noncurrent_versions": ovhtypes.TfInt64Type{},
	}
}

// ---------------------------------------------------------------------------
// LifecycleRuleTransitionType / LifecycleRuleTransitionValue
// ---------------------------------------------------------------------------

var _ basetypes.ObjectTypable = LifecycleRuleTransitionType{}

type LifecycleRuleTransitionType struct {
	basetypes.ObjectType
}

func (t LifecycleRuleTransitionType) Equal(o attr.Type) bool {
	other, ok := o.(LifecycleRuleTransitionType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

func (t LifecycleRuleTransitionType) String() string { return "LifecycleRuleTransitionType" }

func (t LifecycleRuleTransitionType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	attributes := in.Attributes()
	daysVal, _ := attributes["days"].(ovhtypes.TfInt64Value)
	dateVal, _ := attributes["date"].(ovhtypes.TfStringValue)
	storageClassVal, _ := attributes["storage_class"].(ovhtypes.TfStringValue)
	return LifecycleRuleTransitionValue{
		Days:         daysVal,
		Date:         dateVal,
		StorageClass: storageClassVal,
		state:        attr.ValueStateKnown,
	}, nil
}

func NewLifecycleRuleTransitionValueNull() LifecycleRuleTransitionValue {
	return LifecycleRuleTransitionValue{state: attr.ValueStateNull}
}

func NewLifecycleRuleTransitionValueUnknown() LifecycleRuleTransitionValue {
	return LifecycleRuleTransitionValue{state: attr.ValueStateUnknown}
}

func (t LifecycleRuleTransitionType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewLifecycleRuleTransitionValueNull(), nil
	}
	if !in.IsKnown() {
		return NewLifecycleRuleTransitionValueUnknown(), nil
	}
	if in.IsNull() {
		return NewLifecycleRuleTransitionValueNull(), nil
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
	obj, diags := t.ValueFromObject(ctx, basetypes.NewObjectValueMust(LifecycleRuleTransitionValue{}.AttributeTypes(ctx), attributes))
	if diags.HasError() {
		return nil, fmt.Errorf("error converting object: %v", diags)
	}
	return obj, nil
}

func (t LifecycleRuleTransitionType) ValueType(ctx context.Context) attr.Value {
	return LifecycleRuleTransitionValue{}
}

var _ basetypes.ObjectValuable = LifecycleRuleTransitionValue{}

type LifecycleRuleTransitionValue struct {
	Days         ovhtypes.TfInt64Value  `tfsdk:"days" json:"days"`
	Date         ovhtypes.TfStringValue `tfsdk:"date" json:"date"`
	StorageClass ovhtypes.TfStringValue `tfsdk:"storage_class" json:"storageClass"`
	state        attr.ValueState
}

func (v *LifecycleRuleTransitionValue) UnmarshalJSON(data []byte) error {
	type JsonValue LifecycleRuleTransitionValue
	var tmp JsonValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.Days = tmp.Days
	v.Date = tmp.Date
	v.StorageClass = tmp.StorageClass
	v.state = attr.ValueStateKnown
	return nil
}

func (v LifecycleRuleTransitionValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"days":          v.Days,
		"date":          v.Date,
		"storage_class": v.StorageClass,
	}
}

func (v LifecycleRuleTransitionValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := map[string]tftypes.Type{
		"days":          basetypes.NumberType{}.TerraformType(ctx),
		"date":          basetypes.StringType{}.TerraformType(ctx),
		"storage_class": basetypes.StringType{}.TerraformType(ctx),
	}
	objectType := tftypes.Object{AttributeTypes: attrTypes}
	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)
		var val tftypes.Value
		var err error

		val, err = v.Days.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["days"] = val

		val, err = v.Date.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["date"] = val

		val, err = v.StorageClass.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["storage_class"] = val

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

func (v LifecycleRuleTransitionValue) IsNull() bool    { return v.state == attr.ValueStateNull }
func (v LifecycleRuleTransitionValue) IsUnknown() bool { return v.state == attr.ValueStateUnknown }
func (v LifecycleRuleTransitionValue) String() string  { return "LifecycleRuleTransitionValue" }

func (v LifecycleRuleTransitionValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(
		map[string]attr.Type{
			"days":          ovhtypes.TfInt64Type{},
			"date":          ovhtypes.TfStringType{},
			"storage_class": ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"days":          v.Days,
			"date":          v.Date,
			"storage_class": v.StorageClass,
		},
	)
}

func (v LifecycleRuleTransitionValue) Equal(o attr.Value) bool {
	other, ok := o.(LifecycleRuleTransitionValue)
	if !ok {
		return false
	}
	if v.state != other.state {
		return false
	}
	if v.state != attr.ValueStateKnown {
		return true
	}
	return v.Days.Equal(other.Days) && v.Date.Equal(other.Date) && v.StorageClass.Equal(other.StorageClass)
}

func (v LifecycleRuleTransitionValue) Type(ctx context.Context) attr.Type {
	return LifecycleRuleTransitionType{basetypes.ObjectType{AttrTypes: v.AttributeTypes(ctx)}}
}

func (v LifecycleRuleTransitionValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"days":          ovhtypes.TfInt64Type{},
		"date":          ovhtypes.TfStringType{},
		"storage_class": ovhtypes.TfStringType{},
	}
}

// ---------------------------------------------------------------------------
// LifecycleRuleNoncurrentVersionTransitionType / Value
// ---------------------------------------------------------------------------

var _ basetypes.ObjectTypable = LifecycleRuleNoncurrentVersionTransitionType{}

type LifecycleRuleNoncurrentVersionTransitionType struct {
	basetypes.ObjectType
}

func (t LifecycleRuleNoncurrentVersionTransitionType) Equal(o attr.Type) bool {
	other, ok := o.(LifecycleRuleNoncurrentVersionTransitionType)
	if !ok {
		return false
	}
	return t.ObjectType.Equal(other.ObjectType)
}

func (t LifecycleRuleNoncurrentVersionTransitionType) String() string {
	return "LifecycleRuleNoncurrentVersionTransitionType"
}

func (t LifecycleRuleNoncurrentVersionTransitionType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	attributes := in.Attributes()
	noncurrentDaysVal, _ := attributes["noncurrent_days"].(ovhtypes.TfInt64Value)
	newerVersionsVal, _ := attributes["newer_noncurrent_versions"].(ovhtypes.TfInt64Value)
	storageClassVal, _ := attributes["storage_class"].(ovhtypes.TfStringValue)
	return LifecycleRuleNoncurrentVersionTransitionValue{
		NoncurrentDays:          noncurrentDaysVal,
		NewerNoncurrentVersions: newerVersionsVal,
		StorageClass:            storageClassVal,
		state:                   attr.ValueStateKnown,
	}, nil
}

func NewLifecycleRuleNoncurrentVersionTransitionValueNull() LifecycleRuleNoncurrentVersionTransitionValue {
	return LifecycleRuleNoncurrentVersionTransitionValue{state: attr.ValueStateNull}
}

func NewLifecycleRuleNoncurrentVersionTransitionValueUnknown() LifecycleRuleNoncurrentVersionTransitionValue {
	return LifecycleRuleNoncurrentVersionTransitionValue{state: attr.ValueStateUnknown}
}

func (t LifecycleRuleNoncurrentVersionTransitionType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewLifecycleRuleNoncurrentVersionTransitionValueNull(), nil
	}
	if !in.IsKnown() {
		return NewLifecycleRuleNoncurrentVersionTransitionValueUnknown(), nil
	}
	if in.IsNull() {
		return NewLifecycleRuleNoncurrentVersionTransitionValueNull(), nil
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
	obj, diags := t.ValueFromObject(ctx, basetypes.NewObjectValueMust(LifecycleRuleNoncurrentVersionTransitionValue{}.AttributeTypes(ctx), attributes))
	if diags.HasError() {
		return nil, fmt.Errorf("error converting object: %v", diags)
	}
	return obj, nil
}

func (t LifecycleRuleNoncurrentVersionTransitionType) ValueType(ctx context.Context) attr.Value {
	return LifecycleRuleNoncurrentVersionTransitionValue{}
}

var _ basetypes.ObjectValuable = LifecycleRuleNoncurrentVersionTransitionValue{}

type LifecycleRuleNoncurrentVersionTransitionValue struct {
	NoncurrentDays          ovhtypes.TfInt64Value  `tfsdk:"noncurrent_days" json:"noncurrentDays"`
	NewerNoncurrentVersions ovhtypes.TfInt64Value  `tfsdk:"newer_noncurrent_versions" json:"newerNoncurrentVersions"`
	StorageClass            ovhtypes.TfStringValue `tfsdk:"storage_class" json:"storageClass"`
	state                   attr.ValueState
}

func (v *LifecycleRuleNoncurrentVersionTransitionValue) UnmarshalJSON(data []byte) error {
	type JsonValue LifecycleRuleNoncurrentVersionTransitionValue
	var tmp JsonValue
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	v.NoncurrentDays = tmp.NoncurrentDays
	v.NewerNoncurrentVersions = tmp.NewerNoncurrentVersions
	v.StorageClass = tmp.StorageClass
	v.state = attr.ValueStateKnown
	return nil
}

func (v LifecycleRuleNoncurrentVersionTransitionValue) Attributes() map[string]attr.Value {
	return map[string]attr.Value{
		"noncurrent_days":           v.NoncurrentDays,
		"newer_noncurrent_versions": v.NewerNoncurrentVersions,
		"storage_class":             v.StorageClass,
	}
}

func (v LifecycleRuleNoncurrentVersionTransitionValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	attrTypes := map[string]tftypes.Type{
		"noncurrent_days":           basetypes.NumberType{}.TerraformType(ctx),
		"newer_noncurrent_versions": basetypes.NumberType{}.TerraformType(ctx),
		"storage_class":             basetypes.StringType{}.TerraformType(ctx),
	}
	objectType := tftypes.Object{AttributeTypes: attrTypes}
	switch v.state {
	case attr.ValueStateKnown:
		vals := make(map[string]tftypes.Value, 3)
		var val tftypes.Value
		var err error

		val, err = v.NoncurrentDays.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["noncurrent_days"] = val

		val, err = v.NewerNoncurrentVersions.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["newer_noncurrent_versions"] = val

		val, err = v.StorageClass.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.NewValue(objectType, tftypes.UnknownValue), err
		}
		vals["storage_class"] = val

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

func (v LifecycleRuleNoncurrentVersionTransitionValue) IsNull() bool {
	return v.state == attr.ValueStateNull
}
func (v LifecycleRuleNoncurrentVersionTransitionValue) IsUnknown() bool {
	return v.state == attr.ValueStateUnknown
}
func (v LifecycleRuleNoncurrentVersionTransitionValue) String() string {
	return "LifecycleRuleNoncurrentVersionTransitionValue"
}

func (v LifecycleRuleNoncurrentVersionTransitionValue) ToObjectValue(ctx context.Context) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(
		map[string]attr.Type{
			"noncurrent_days":           ovhtypes.TfInt64Type{},
			"newer_noncurrent_versions": ovhtypes.TfInt64Type{},
			"storage_class":             ovhtypes.TfStringType{},
		},
		map[string]attr.Value{
			"noncurrent_days":           v.NoncurrentDays,
			"newer_noncurrent_versions": v.NewerNoncurrentVersions,
			"storage_class":             v.StorageClass,
		},
	)
}

func (v LifecycleRuleNoncurrentVersionTransitionValue) Equal(o attr.Value) bool {
	other, ok := o.(LifecycleRuleNoncurrentVersionTransitionValue)
	if !ok {
		return false
	}
	if v.state != other.state {
		return false
	}
	if v.state != attr.ValueStateKnown {
		return true
	}
	return v.NoncurrentDays.Equal(other.NoncurrentDays) &&
		v.NewerNoncurrentVersions.Equal(other.NewerNoncurrentVersions) &&
		v.StorageClass.Equal(other.StorageClass)
}

func (v LifecycleRuleNoncurrentVersionTransitionValue) Type(ctx context.Context) attr.Type {
	return LifecycleRuleNoncurrentVersionTransitionType{basetypes.ObjectType{AttrTypes: v.AttributeTypes(ctx)}}
}

func (v LifecycleRuleNoncurrentVersionTransitionValue) AttributeTypes(ctx context.Context) map[string]attr.Type {
	return map[string]attr.Type{
		"noncurrent_days":           ovhtypes.TfInt64Type{},
		"newer_noncurrent_versions": ovhtypes.TfInt64Type{},
		"storage_class":             ovhtypes.TfStringType{},
	}
}
