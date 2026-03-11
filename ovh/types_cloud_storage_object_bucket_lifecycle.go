package ovh

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudStorageObjectBucketLifecycleModel is the Terraform state model for the lifecycle resource.
type CloudStorageObjectBucketLifecycleModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketName  ovhtypes.TfStringValue `tfsdk:"bucket_name"`
	Id          ovhtypes.TfStringValue `tfsdk:"id"`
	Rules       types.List             `tfsdk:"rules"`
}

// API request/response types

type CloudBucketLifecycleRequest struct {
	Rules []CloudBucketLifecycleRule `json:"rules"`
}

type CloudBucketLifecycleResponse struct {
	Rules []CloudBucketLifecycleRule `json:"rules"`
}

type CloudBucketLifecycleRule struct {
	ID                             string                                `json:"id"`
	Status                         string                                `json:"status"`
	Filter                         *CloudBucketLifecycleFilter           `json:"filter,omitempty"`
	Expiration                     *CloudBucketLifecycleExpiration       `json:"expiration,omitempty"`
	Transitions                    []CloudBucketLifecycleTransition      `json:"transitions,omitempty"`
	NoncurrentVersionExpiration    *CloudBucketLifecycleNoncurrentExp    `json:"noncurrentVersionExpiration,omitempty"`
	NoncurrentVersionTransitions   []CloudBucketLifecycleNoncurrentTrans `json:"noncurrentVersionTransitions,omitempty"`
	AbortIncompleteMultipartUpload *CloudBucketLifecycleAbortIncomplete  `json:"abortIncompleteMultipartUpload,omitempty"`
}

type CloudBucketLifecycleFilter struct {
	Prefix                string            `json:"prefix,omitempty"`
	Tags                  map[string]string `json:"tags,omitempty"`
	ObjectSizeGreaterThan int64             `json:"objectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan    int64             `json:"objectSizeLessThan,omitempty"`
}

type CloudBucketLifecycleExpiration struct {
	Date                      string `json:"date,omitempty"`
	Days                      int64  `json:"days,omitempty"`
	ExpiredObjectDeleteMarker bool   `json:"expiredObjectDeleteMarker,omitempty"`
}

type CloudBucketLifecycleTransition struct {
	Date         string `json:"date,omitempty"`
	Days         int64  `json:"days,omitempty"`
	StorageClass string `json:"storageClass"`
}

type CloudBucketLifecycleNoncurrentExp struct {
	NoncurrentDays          int64 `json:"noncurrentDays,omitempty"`
	NewerNoncurrentVersions int64 `json:"newerNoncurrentVersions,omitempty"`
}

type CloudBucketLifecycleNoncurrentTrans struct {
	NoncurrentDays          int64  `json:"noncurrentDays,omitempty"`
	NewerNoncurrentVersions int64  `json:"newerNoncurrentVersions,omitempty"`
	StorageClass            string `json:"storageClass"`
}

type CloudBucketLifecycleAbortIncomplete struct {
	DaysAfterInitiation int64 `json:"daysAfterInitiation,omitempty"`
}

// Attribute type helpers

func lifecycleFilterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"prefix":                   ovhtypes.TfStringType{},
		"tags":                     types.MapType{ElemType: types.StringType},
		"object_size_greater_than": types.Int64Type,
		"object_size_less_than":    types.Int64Type,
	}
}

func lifecycleExpirationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"date":                         ovhtypes.TfStringType{},
		"days":                         types.Int64Type,
		"expired_object_delete_marker": types.BoolType,
	}
}

func lifecycleTransitionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"date":          ovhtypes.TfStringType{},
		"days":          types.Int64Type,
		"storage_class": ovhtypes.TfStringType{},
	}
}

func lifecycleNoncurrentExpAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"noncurrent_days":           types.Int64Type,
		"newer_noncurrent_versions": types.Int64Type,
	}
}

func lifecycleNoncurrentTransAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"noncurrent_days":           types.Int64Type,
		"newer_noncurrent_versions": types.Int64Type,
		"storage_class":             ovhtypes.TfStringType{},
	}
}

func lifecycleAbortIncompleteAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"days_after_initiation": types.Int64Type,
	}
}

func lifecycleRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                                ovhtypes.TfStringType{},
		"status":                            ovhtypes.TfStringType{},
		"filter":                            types.ObjectType{AttrTypes: lifecycleFilterAttrTypes()},
		"expiration":                        types.ObjectType{AttrTypes: lifecycleExpirationAttrTypes()},
		"transitions":                       types.ListType{ElemType: types.ObjectType{AttrTypes: lifecycleTransitionAttrTypes()}},
		"noncurrent_version_expiration":     types.ObjectType{AttrTypes: lifecycleNoncurrentExpAttrTypes()},
		"noncurrent_version_transitions":    types.ListType{ElemType: types.ObjectType{AttrTypes: lifecycleNoncurrentTransAttrTypes()}},
		"abort_incomplete_multipart_upload": types.ObjectType{AttrTypes: lifecycleAbortIncompleteAttrTypes()},
	}
}

// apiLifecycleRulesToTFList converts API rules to a Terraform list, reordered to match the
// state list order (by rule ID) to prevent spurious diffs.
func apiLifecycleRulesToTFList(apiRules []CloudBucketLifecycleRule, stateRules types.List) types.List {
	ruleElemType := types.ObjectType{AttrTypes: lifecycleRuleAttrTypes()}

	// Index API rules by ID for reordering
	ruleByID := make(map[string]CloudBucketLifecycleRule, len(apiRules))
	for _, r := range apiRules {
		ruleByID[r.ID] = r
	}

	// Follow state order first
	var ordered []CloudBucketLifecycleRule
	if !stateRules.IsNull() && !stateRules.IsUnknown() {
		for _, elem := range stateRules.Elements() {
			obj, ok := elem.(types.Object)
			if !ok {
				continue
			}
			idVal, ok := obj.Attributes()["id"].(ovhtypes.TfStringValue)
			if !ok {
				continue
			}
			if r, found := ruleByID[idVal.ValueString()]; found {
				ordered = append(ordered, r)
				delete(ruleByID, idVal.ValueString())
			}
		}
	}

	// Append any new rules from the API not present in state
	for _, r := range apiRules {
		if _, remaining := ruleByID[r.ID]; remaining {
			ordered = append(ordered, r)
		}
	}

	elems := make([]attr.Value, 0, len(ordered))
	for _, r := range ordered {
		elems = append(elems, apiLifecycleRuleToTFObject(r))
	}

	list, _ := types.ListValue(ruleElemType, elems)
	return list
}

func apiLifecycleRuleToTFObject(r CloudBucketLifecycleRule) basetypes.ObjectValue {
	// filter
	filterObj := types.ObjectNull(lifecycleFilterAttrTypes())
	if r.Filter != nil {
		var tagsVal basetypes.MapValue
		if len(r.Filter.Tags) > 0 {
			tagVals := make(map[string]attr.Value, len(r.Filter.Tags))
			for k, v := range r.Filter.Tags {
				tagVals[k] = types.StringValue(v)
			}
			tagsVal, _ = types.MapValue(types.StringType, tagVals)
		} else {
			tagsVal = types.MapNull(types.StringType)
		}
		filterObj, _ = types.ObjectValue(lifecycleFilterAttrTypes(), map[string]attr.Value{
			"prefix":                   ovhtypes.TfStringValue{StringValue: types.StringValue(r.Filter.Prefix)},
			"tags":                     tagsVal,
			"object_size_greater_than": types.Int64Value(r.Filter.ObjectSizeGreaterThan),
			"object_size_less_than":    types.Int64Value(r.Filter.ObjectSizeLessThan),
		})
	}

	// expiration
	expirationObj := types.ObjectNull(lifecycleExpirationAttrTypes())
	if r.Expiration != nil {
		expirationObj, _ = types.ObjectValue(lifecycleExpirationAttrTypes(), map[string]attr.Value{
			"date":                         ovhtypes.TfStringValue{StringValue: types.StringValue(r.Expiration.Date)},
			"days":                         types.Int64Value(r.Expiration.Days),
			"expired_object_delete_marker": types.BoolValue(r.Expiration.ExpiredObjectDeleteMarker),
		})
	}

	// transitions
	transElems := make([]attr.Value, 0, len(r.Transitions))
	for _, t := range r.Transitions {
		transObj, _ := types.ObjectValue(lifecycleTransitionAttrTypes(), map[string]attr.Value{
			"date":          ovhtypes.TfStringValue{StringValue: types.StringValue(t.Date)},
			"days":          types.Int64Value(t.Days),
			"storage_class": ovhtypes.TfStringValue{StringValue: types.StringValue(t.StorageClass)},
		})
		transElems = append(transElems, transObj)
	}
	transitionsVal, _ := types.ListValue(types.ObjectType{AttrTypes: lifecycleTransitionAttrTypes()}, transElems)

	// noncurrent_version_expiration
	ncveObj := types.ObjectNull(lifecycleNoncurrentExpAttrTypes())
	if r.NoncurrentVersionExpiration != nil {
		ncveObj, _ = types.ObjectValue(lifecycleNoncurrentExpAttrTypes(), map[string]attr.Value{
			"noncurrent_days":           types.Int64Value(r.NoncurrentVersionExpiration.NoncurrentDays),
			"newer_noncurrent_versions": types.Int64Value(r.NoncurrentVersionExpiration.NewerNoncurrentVersions),
		})
	}

	// noncurrent_version_transitions
	ncvtElems := make([]attr.Value, 0, len(r.NoncurrentVersionTransitions))
	for _, t := range r.NoncurrentVersionTransitions {
		ncvtObj, _ := types.ObjectValue(lifecycleNoncurrentTransAttrTypes(), map[string]attr.Value{
			"noncurrent_days":           types.Int64Value(t.NoncurrentDays),
			"newer_noncurrent_versions": types.Int64Value(t.NewerNoncurrentVersions),
			"storage_class":             ovhtypes.TfStringValue{StringValue: types.StringValue(t.StorageClass)},
		})
		ncvtElems = append(ncvtElems, ncvtObj)
	}
	ncvtVal, _ := types.ListValue(types.ObjectType{AttrTypes: lifecycleNoncurrentTransAttrTypes()}, ncvtElems)

	// abort_incomplete_multipart_upload
	abortObj := types.ObjectNull(lifecycleAbortIncompleteAttrTypes())
	if r.AbortIncompleteMultipartUpload != nil {
		abortObj, _ = types.ObjectValue(lifecycleAbortIncompleteAttrTypes(), map[string]attr.Value{
			"days_after_initiation": types.Int64Value(r.AbortIncompleteMultipartUpload.DaysAfterInitiation),
		})
	}

	obj, _ := types.ObjectValue(lifecycleRuleAttrTypes(), map[string]attr.Value{
		"id":                                ovhtypes.TfStringValue{StringValue: types.StringValue(r.ID)},
		"status":                            ovhtypes.TfStringValue{StringValue: types.StringValue(r.Status)},
		"filter":                            filterObj,
		"expiration":                        expirationObj,
		"transitions":                       transitionsVal,
		"noncurrent_version_expiration":     ncveObj,
		"noncurrent_version_transitions":    ncvtVal,
		"abort_incomplete_multipart_upload": abortObj,
	})
	return obj
}

// tfListToAPILifecycleRules converts the Terraform rules list to the API slice.
func tfListToAPILifecycleRules(list types.List) []CloudBucketLifecycleRule {
	if list.IsNull() || list.IsUnknown() {
		return []CloudBucketLifecycleRule{}
	}

	rules := make([]CloudBucketLifecycleRule, 0, len(list.Elements()))
	for _, elem := range list.Elements() {
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}
		rules = append(rules, tfObjectToAPILifecycleRule(obj))
	}
	return rules
}

func tfObjectToAPILifecycleRule(obj types.Object) CloudBucketLifecycleRule {
	attrs := obj.Attributes()
	rule := CloudBucketLifecycleRule{}

	if v, ok := attrs["id"].(ovhtypes.TfStringValue); ok {
		rule.ID = v.ValueString()
	}
	if v, ok := attrs["status"].(ovhtypes.TfStringValue); ok {
		rule.Status = v.ValueString()
	}

	if filterObj, ok := attrs["filter"].(types.Object); ok && !filterObj.IsNull() && !filterObj.IsUnknown() {
		fa := filterObj.Attributes()
		filter := &CloudBucketLifecycleFilter{}
		if v, ok := fa["prefix"].(ovhtypes.TfStringValue); ok && !v.IsNull() {
			filter.Prefix = v.ValueString()
		}
		if v, ok := fa["object_size_greater_than"].(basetypes.Int64Value); ok && !v.IsNull() {
			filter.ObjectSizeGreaterThan = v.ValueInt64()
		}
		if v, ok := fa["object_size_less_than"].(basetypes.Int64Value); ok && !v.IsNull() {
			filter.ObjectSizeLessThan = v.ValueInt64()
		}
		if tagsMap, ok := fa["tags"].(types.Map); ok && !tagsMap.IsNull() {
			tags := make(map[string]string, len(tagsMap.Elements()))
			for k, v := range tagsMap.Elements() {
				if sv, ok := v.(types.String); ok {
					tags[k] = sv.ValueString()
				}
			}
			if len(tags) > 0 {
				filter.Tags = tags
			}
		}
		rule.Filter = filter
	}

	if expObj, ok := attrs["expiration"].(types.Object); ok && !expObj.IsNull() && !expObj.IsUnknown() {
		ea := expObj.Attributes()
		exp := &CloudBucketLifecycleExpiration{}
		if v, ok := ea["date"].(ovhtypes.TfStringValue); ok && !v.IsNull() {
			exp.Date = v.ValueString()
		}
		if v, ok := ea["days"].(basetypes.Int64Value); ok && !v.IsNull() {
			exp.Days = v.ValueInt64()
		}
		if v, ok := ea["expired_object_delete_marker"].(basetypes.BoolValue); ok && !v.IsNull() {
			exp.ExpiredObjectDeleteMarker = v.ValueBool()
		}
		rule.Expiration = exp
	}

	if transList, ok := attrs["transitions"].(types.List); ok && !transList.IsNull() && !transList.IsUnknown() {
		for _, elem := range transList.Elements() {
			transObj, ok := elem.(types.Object)
			if !ok {
				continue
			}
			ta := transObj.Attributes()
			trans := CloudBucketLifecycleTransition{}
			if v, ok := ta["date"].(ovhtypes.TfStringValue); ok && !v.IsNull() {
				trans.Date = v.ValueString()
			}
			if v, ok := ta["days"].(basetypes.Int64Value); ok && !v.IsNull() {
				trans.Days = v.ValueInt64()
			}
			if v, ok := ta["storage_class"].(ovhtypes.TfStringValue); ok {
				trans.StorageClass = v.ValueString()
			}
			rule.Transitions = append(rule.Transitions, trans)
		}
	}

	if ncveObj, ok := attrs["noncurrent_version_expiration"].(types.Object); ok && !ncveObj.IsNull() && !ncveObj.IsUnknown() {
		na := ncveObj.Attributes()
		ncve := &CloudBucketLifecycleNoncurrentExp{}
		if v, ok := na["noncurrent_days"].(basetypes.Int64Value); ok && !v.IsNull() {
			ncve.NoncurrentDays = v.ValueInt64()
		}
		if v, ok := na["newer_noncurrent_versions"].(basetypes.Int64Value); ok && !v.IsNull() {
			ncve.NewerNoncurrentVersions = v.ValueInt64()
		}
		rule.NoncurrentVersionExpiration = ncve
	}

	if ncvtList, ok := attrs["noncurrent_version_transitions"].(types.List); ok && !ncvtList.IsNull() && !ncvtList.IsUnknown() {
		for _, elem := range ncvtList.Elements() {
			ncvtObj, ok := elem.(types.Object)
			if !ok {
				continue
			}
			na := ncvtObj.Attributes()
			ncvt := CloudBucketLifecycleNoncurrentTrans{}
			if v, ok := na["noncurrent_days"].(basetypes.Int64Value); ok && !v.IsNull() {
				ncvt.NoncurrentDays = v.ValueInt64()
			}
			if v, ok := na["newer_noncurrent_versions"].(basetypes.Int64Value); ok && !v.IsNull() {
				ncvt.NewerNoncurrentVersions = v.ValueInt64()
			}
			if v, ok := na["storage_class"].(ovhtypes.TfStringValue); ok {
				ncvt.StorageClass = v.ValueString()
			}
			rule.NoncurrentVersionTransitions = append(rule.NoncurrentVersionTransitions, ncvt)
		}
	}

	if abortObj, ok := attrs["abort_incomplete_multipart_upload"].(types.Object); ok && !abortObj.IsNull() && !abortObj.IsUnknown() {
		aa := abortObj.Attributes()
		abort := &CloudBucketLifecycleAbortIncomplete{}
		if v, ok := aa["days_after_initiation"].(basetypes.Int64Value); ok && !v.IsNull() {
			abort.DaysAfterInitiation = v.ValueInt64()
		}
		rule.AbortIncompleteMultipartUpload = abort
	}

	return rule
}
