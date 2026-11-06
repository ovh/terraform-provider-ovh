package ovh

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudStorageObjectBucketLifecycleModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketId    ovhtypes.TfStringValue `tfsdk:"bucket_id"`
	Id          ovhtypes.TfStringValue `tfsdk:"id"`
	Rules       types.List             `tfsdk:"rules"`
}

type storageObjectBucketLifecycleRuleModel struct {
	Id                             ovhtypes.TfStringValue                                         `tfsdk:"id"`
	Status                         ovhtypes.TfStringValue                                         `tfsdk:"status"`
	Filter                         *storageObjectBucketLifecycleFilterModel                       `tfsdk:"filter"`
	Expiration                     *storageObjectBucketLifecycleExpirationModel                   `tfsdk:"expiration"`
	Transitions                    []storageObjectBucketLifecycleTransitionModel                  `tfsdk:"transitions"`
	NoncurrentVersionExpiration    *storageObjectBucketLifecycleNoncurrentVersionExpirationModel  `tfsdk:"noncurrent_version_expiration"`
	NoncurrentVersionTransitions   []storageObjectBucketLifecycleNoncurrentVersionTransitionModel `tfsdk:"noncurrent_version_transitions"`
	AbortIncompleteMultipartUpload *storageObjectBucketLifecycleAbortMultipartModel               `tfsdk:"abort_incomplete_multipart_upload"`
}

type storageObjectBucketLifecycleFilterModel struct {
	Prefix                ovhtypes.TfStringValue `tfsdk:"prefix"`
	Tags                  types.Map              `tfsdk:"tags"`
	ObjectSizeGreaterThan types.Int64            `tfsdk:"object_size_greater_than"`
	ObjectSizeLessThan    types.Int64            `tfsdk:"object_size_less_than"`
}

type storageObjectBucketLifecycleExpirationModel struct {
	Date                      ovhtypes.TfStringValue `tfsdk:"date"`
	Days                      types.Int64            `tfsdk:"days"`
	ExpiredObjectDeleteMarker types.Bool             `tfsdk:"expired_object_delete_marker"`
}

type storageObjectBucketLifecycleTransitionModel struct {
	Date         ovhtypes.TfStringValue `tfsdk:"date"`
	Days         types.Int64            `tfsdk:"days"`
	StorageClass ovhtypes.TfStringValue `tfsdk:"storage_class"`
}

type storageObjectBucketLifecycleNoncurrentVersionExpirationModel struct {
	NoncurrentDays          types.Int64 `tfsdk:"noncurrent_days"`
	NewerNoncurrentVersions types.Int64 `tfsdk:"newer_noncurrent_versions"`
}

type storageObjectBucketLifecycleNoncurrentVersionTransitionModel struct {
	NoncurrentDays          types.Int64            `tfsdk:"noncurrent_days"`
	NewerNoncurrentVersions types.Int64            `tfsdk:"newer_noncurrent_versions"`
	StorageClass            ovhtypes.TfStringValue `tfsdk:"storage_class"`
}

type storageObjectBucketLifecycleAbortMultipartModel struct {
	DaysAfterInitiation types.Int64 `tfsdk:"days_after_initiation"`
}

type CloudStorageObjectBucketLifecycleAPIFilter struct {
	Prefix                string            `json:"prefix,omitempty"`
	Tags                  map[string]string `json:"tags,omitempty"`
	ObjectSizeGreaterThan int64             `json:"objectSizeGreaterThan,omitempty"`
	ObjectSizeLessThan    int64             `json:"objectSizeLessThan,omitempty"`
}

type CloudStorageObjectBucketLifecycleAPIExpiration struct {
	Date                      string `json:"date,omitempty"`
	Days                      int64  `json:"days,omitempty"`
	ExpiredObjectDeleteMarker bool   `json:"expiredObjectDeleteMarker,omitempty"`
}

type CloudStorageObjectBucketLifecycleAPITransition struct {
	Date         string `json:"date,omitempty"`
	Days         int64  `json:"days,omitempty"`
	StorageClass string `json:"storageClass"`
}

type CloudStorageObjectBucketLifecycleAPINoncurrentVersionExpiration struct {
	NoncurrentDays          int64 `json:"noncurrentDays,omitempty"`
	NewerNoncurrentVersions int64 `json:"newerNoncurrentVersions,omitempty"`
}

type CloudStorageObjectBucketLifecycleAPINoncurrentVersionTransition struct {
	NoncurrentDays          int64  `json:"noncurrentDays,omitempty"`
	NewerNoncurrentVersions int64  `json:"newerNoncurrentVersions,omitempty"`
	StorageClass            string `json:"storageClass"`
}

type CloudStorageObjectBucketLifecycleAPIAbortMultipart struct {
	DaysAfterInitiation int64 `json:"daysAfterInitiation,omitempty"`
}

type CloudStorageObjectBucketLifecycleAPIRule struct {
	Id                             string                                                            `json:"id"`
	Status                         string                                                            `json:"status"`
	Filter                         *CloudStorageObjectBucketLifecycleAPIFilter                       `json:"filter,omitempty"`
	Expiration                     *CloudStorageObjectBucketLifecycleAPIExpiration                   `json:"expiration,omitempty"`
	Transitions                    []CloudStorageObjectBucketLifecycleAPITransition                  `json:"transitions,omitempty"`
	NoncurrentVersionExpiration    *CloudStorageObjectBucketLifecycleAPINoncurrentVersionExpiration  `json:"noncurrentVersionExpiration,omitempty"`
	NoncurrentVersionTransitions   []CloudStorageObjectBucketLifecycleAPINoncurrentVersionTransition `json:"noncurrentVersionTransitions,omitempty"`
	AbortIncompleteMultipartUpload *CloudStorageObjectBucketLifecycleAPIAbortMultipart               `json:"abortIncompleteMultipartUpload,omitempty"`
}

// CloudStorageObjectBucketLifecycleAPIRules is both the PUT payload and the GET/PUT response.
type CloudStorageObjectBucketLifecycleAPIRules struct {
	Rules []CloudStorageObjectBucketLifecycleAPIRule `json:"rules"`
}

func storageObjectBucketLifecycleFilterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"prefix":                   ovhtypes.TfStringType{},
		"tags":                     types.MapType{ElemType: ovhtypes.TfStringType{}},
		"object_size_greater_than": types.Int64Type,
		"object_size_less_than":    types.Int64Type,
	}
}

func storageObjectBucketLifecycleRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":     ovhtypes.TfStringType{},
		"status": ovhtypes.TfStringType{},
		"filter": types.ObjectType{AttrTypes: storageObjectBucketLifecycleFilterAttrTypes()},
		"expiration": types.ObjectType{AttrTypes: map[string]attr.Type{
			"date":                         ovhtypes.TfStringType{},
			"days":                         types.Int64Type,
			"expired_object_delete_marker": types.BoolType,
		}},
		"transitions": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"date":          ovhtypes.TfStringType{},
			"days":          types.Int64Type,
			"storage_class": ovhtypes.TfStringType{},
		}}},
		"noncurrent_version_expiration": types.ObjectType{AttrTypes: map[string]attr.Type{
			"noncurrent_days":           types.Int64Type,
			"newer_noncurrent_versions": types.Int64Type,
		}},
		"noncurrent_version_transitions": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"noncurrent_days":           types.Int64Type,
			"newer_noncurrent_versions": types.Int64Type,
			"storage_class":             ovhtypes.TfStringType{},
		}}},
		"abort_incomplete_multipart_upload": types.ObjectType{AttrTypes: map[string]attr.Type{
			"days_after_initiation": types.Int64Type,
		}},
	}
}

func (m *CloudStorageObjectBucketLifecycleModel) rules(ctx context.Context) ([]storageObjectBucketLifecycleRuleModel, diag.Diagnostics) {
	var rules []storageObjectBucketLifecycleRuleModel
	if m.Rules.IsNull() || m.Rules.IsUnknown() {
		return rules, nil
	}
	diags := m.Rules.ElementsAs(ctx, &rules, false)
	return rules, diags
}

func (m *CloudStorageObjectBucketLifecycleModel) ToPut(ctx context.Context) (*CloudStorageObjectBucketLifecycleAPIRules, diag.Diagnostics) {
	rules, diags := m.rules(ctx)
	if diags.HasError() {
		return nil, diags
	}

	payload := &CloudStorageObjectBucketLifecycleAPIRules{Rules: make([]CloudStorageObjectBucketLifecycleAPIRule, 0, len(rules))}
	for _, rule := range rules {
		apiRule := CloudStorageObjectBucketLifecycleAPIRule{
			Id:     rule.Id.ValueString(),
			Status: rule.Status.ValueString(),
		}

		if f := rule.Filter; f != nil {
			apiRule.Filter = &CloudStorageObjectBucketLifecycleAPIFilter{
				Prefix:                f.Prefix.ValueString(),
				ObjectSizeGreaterThan: f.ObjectSizeGreaterThan.ValueInt64(),
				ObjectSizeLessThan:    f.ObjectSizeLessThan.ValueInt64(),
			}
			if !f.Tags.IsNull() && !f.Tags.IsUnknown() {
				tags := make(map[string]ovhtypes.TfStringValue, len(f.Tags.Elements()))
				diags.Append(f.Tags.ElementsAs(ctx, &tags, false)...)
				apiRule.Filter.Tags = make(map[string]string, len(tags))
				for k, v := range tags {
					apiRule.Filter.Tags[k] = v.ValueString()
				}
			}
		}

		if e := rule.Expiration; e != nil {
			apiRule.Expiration = &CloudStorageObjectBucketLifecycleAPIExpiration{
				Date:                      e.Date.ValueString(),
				Days:                      e.Days.ValueInt64(),
				ExpiredObjectDeleteMarker: e.ExpiredObjectDeleteMarker.ValueBool(),
			}
		}

		for _, t := range rule.Transitions {
			apiRule.Transitions = append(apiRule.Transitions, CloudStorageObjectBucketLifecycleAPITransition{
				Date:         t.Date.ValueString(),
				Days:         t.Days.ValueInt64(),
				StorageClass: t.StorageClass.ValueString(),
			})
		}

		if e := rule.NoncurrentVersionExpiration; e != nil {
			apiRule.NoncurrentVersionExpiration = &CloudStorageObjectBucketLifecycleAPINoncurrentVersionExpiration{
				NoncurrentDays:          e.NoncurrentDays.ValueInt64(),
				NewerNoncurrentVersions: e.NewerNoncurrentVersions.ValueInt64(),
			}
		}

		for _, t := range rule.NoncurrentVersionTransitions {
			apiRule.NoncurrentVersionTransitions = append(apiRule.NoncurrentVersionTransitions, CloudStorageObjectBucketLifecycleAPINoncurrentVersionTransition{
				NoncurrentDays:          t.NoncurrentDays.ValueInt64(),
				NewerNoncurrentVersions: t.NewerNoncurrentVersions.ValueInt64(),
				StorageClass:            t.StorageClass.ValueString(),
			})
		}

		if a := rule.AbortIncompleteMultipartUpload; a != nil {
			apiRule.AbortIncompleteMultipartUpload = &CloudStorageObjectBucketLifecycleAPIAbortMultipart{
				DaysAfterInitiation: a.DaysAfterInitiation.ValueInt64(),
			}
		}

		payload.Rules = append(payload.Rules, apiRule)
	}

	return payload, diags
}

// MergeWith rebuilds rules from the API, ordered and shaped like prior (plan or state): S3 may
// return rules and transitions in another order, and drops zero values the config may spell out.
func (m *CloudStorageObjectBucketLifecycleModel) MergeWith(ctx context.Context, response *CloudStorageObjectBucketLifecycleAPIRules) diag.Diagnostics {
	m.Id = m.BucketId

	prior, diags := m.rules(ctx)
	if diags.HasError() {
		return diags
	}

	priorByID := make(map[string]*storageObjectBucketLifecycleRuleModel, len(prior))
	for i := range prior {
		priorByID[prior[i].Id.ValueString()] = &prior[i]
	}

	ordered := make([]CloudStorageObjectBucketLifecycleAPIRule, 0, len(response.Rules))
	seen := make(map[string]bool, len(response.Rules))
	apiByID := make(map[string]CloudStorageObjectBucketLifecycleAPIRule, len(response.Rules))
	for _, r := range response.Rules {
		apiByID[r.Id] = r
	}
	for _, p := range prior {
		if r, ok := apiByID[p.Id.ValueString()]; ok && !seen[r.Id] {
			ordered = append(ordered, r)
			seen[r.Id] = true
		}
	}
	for _, r := range response.Rules {
		if !seen[r.Id] {
			ordered = append(ordered, r)
			seen[r.Id] = true
		}
	}

	rules := make([]storageObjectBucketLifecycleRuleModel, 0, len(ordered))
	for _, r := range ordered {
		rules = append(rules, buildStorageObjectBucketLifecycleRule(r, priorByID[r.Id]))
	}

	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: storageObjectBucketLifecycleRuleAttrTypes()}, rules)
	diags.Append(d...)
	m.Rules = list
	return diags
}

func buildStorageObjectBucketLifecycleRule(r CloudStorageObjectBucketLifecycleAPIRule, prior *storageObjectBucketLifecycleRuleModel) storageObjectBucketLifecycleRuleModel {
	if prior == nil {
		prior = &storageObjectBucketLifecycleRuleModel{}
	}

	rule := storageObjectBucketLifecycleRuleModel{
		Id:     ovhtypes.TfStringValue{StringValue: types.StringValue(r.Id)},
		Status: ovhtypes.TfStringValue{StringValue: types.StringValue(r.Status)},
		Filter: buildStorageObjectBucketLifecycleFilter(r.Filter, prior.Filter),
	}

	if e := r.Expiration; e != nil {
		priorMarker := types.BoolNull()
		if prior.Expiration != nil {
			priorMarker = prior.Expiration.ExpiredObjectDeleteMarker
		}
		rule.Expiration = &storageObjectBucketLifecycleExpirationModel{
			Date:                      nullableTfString(e.Date),
			Days:                      lifecycleInt64(e.Days),
			ExpiredObjectDeleteMarker: lifecycleBool(e.ExpiredObjectDeleteMarker, priorMarker),
		}
	}

	transitions := orderByStorageClass(r.Transitions, prior.Transitions,
		func(t CloudStorageObjectBucketLifecycleAPITransition) string { return t.StorageClass },
		func(t storageObjectBucketLifecycleTransitionModel) string { return t.StorageClass.ValueString() })
	for _, t := range transitions {
		rule.Transitions = append(rule.Transitions, storageObjectBucketLifecycleTransitionModel{
			Date:         nullableTfString(t.Date),
			Days:         lifecycleInt64(t.Days),
			StorageClass: ovhtypes.TfStringValue{StringValue: types.StringValue(t.StorageClass)},
		})
	}

	if e := r.NoncurrentVersionExpiration; e != nil {
		rule.NoncurrentVersionExpiration = &storageObjectBucketLifecycleNoncurrentVersionExpirationModel{
			NoncurrentDays:          lifecycleInt64(e.NoncurrentDays),
			NewerNoncurrentVersions: lifecycleInt64(e.NewerNoncurrentVersions),
		}
	}

	noncurrentTransitions := orderByStorageClass(r.NoncurrentVersionTransitions, prior.NoncurrentVersionTransitions,
		func(t CloudStorageObjectBucketLifecycleAPINoncurrentVersionTransition) string { return t.StorageClass },
		func(t storageObjectBucketLifecycleNoncurrentVersionTransitionModel) string {
			return t.StorageClass.ValueString()
		})
	for _, t := range noncurrentTransitions {
		rule.NoncurrentVersionTransitions = append(rule.NoncurrentVersionTransitions, storageObjectBucketLifecycleNoncurrentVersionTransitionModel{
			NoncurrentDays:          lifecycleInt64(t.NoncurrentDays),
			NewerNoncurrentVersions: lifecycleInt64(t.NewerNoncurrentVersions),
			StorageClass:            ovhtypes.TfStringValue{StringValue: types.StringValue(t.StorageClass)},
		})
	}

	if a := r.AbortIncompleteMultipartUpload; a != nil {
		rule.AbortIncompleteMultipartUpload = &storageObjectBucketLifecycleAbortMultipartModel{
			DaysAfterInitiation: lifecycleInt64(a.DaysAfterInitiation),
		}
	}

	return rule
}

// The API sends a rule without filter to S3 as an empty prefix, and reads it back as an empty filter.
func buildStorageObjectBucketLifecycleFilter(f *CloudStorageObjectBucketLifecycleAPIFilter, prior *storageObjectBucketLifecycleFilterModel) *storageObjectBucketLifecycleFilterModel {
	empty := f == nil || (f.Prefix == "" && len(f.Tags) == 0 && f.ObjectSizeGreaterThan == 0 && f.ObjectSizeLessThan == 0)
	if empty && (prior == nil || !storageObjectBucketLifecycleFilterIsEmpty(prior)) {
		return nil
	}
	if f == nil {
		f = &CloudStorageObjectBucketLifecycleAPIFilter{}
	}
	if prior == nil {
		prior = &storageObjectBucketLifecycleFilterModel{}
	}

	filter := &storageObjectBucketLifecycleFilterModel{
		Prefix:                lifecycleString(f.Prefix, prior.Prefix),
		ObjectSizeGreaterThan: lifecycleInt64(f.ObjectSizeGreaterThan),
		ObjectSizeLessThan:    lifecycleInt64(f.ObjectSizeLessThan),
		Tags:                  types.MapNull(ovhtypes.TfStringType{}),
	}
	if len(f.Tags) > 0 || (!prior.Tags.IsNull() && !prior.Tags.IsUnknown() && len(prior.Tags.Elements()) == 0) {
		elems := make(map[string]attr.Value, len(f.Tags))
		for k, v := range f.Tags {
			elems[k] = ovhtypes.TfStringValue{StringValue: types.StringValue(v)}
		}
		filter.Tags, _ = types.MapValue(ovhtypes.TfStringType{}, elems)
	}
	return filter
}

func storageObjectBucketLifecycleFilterIsEmpty(f *storageObjectBucketLifecycleFilterModel) bool {
	return f.Prefix.ValueString() == "" &&
		len(f.Tags.Elements()) == 0 &&
		f.ObjectSizeGreaterThan.IsNull() &&
		f.ObjectSizeLessThan.IsNull()
}

func orderByStorageClass[A, P any](api []A, prior []P, apiClass func(A) string, priorClass func(P) string) []A {
	ordered := make([]A, 0, len(api))
	used := make([]bool, len(api))
	for _, p := range prior {
		for i, a := range api {
			if !used[i] && apiClass(a) == priorClass(p) {
				ordered = append(ordered, a)
				used[i] = true
				break
			}
		}
	}
	for i, a := range api {
		if !used[i] {
			ordered = append(ordered, a)
		}
	}
	return ordered
}

// Every integer of a lifecycle rule is validated >= 1, so the API's omitted zero means unset.
func lifecycleInt64(v int64) types.Int64 {
	if v == 0 {
		return types.Int64Null()
	}
	return types.Int64Value(v)
}

// The API omits "" and false; keep them when the prior value spelled them out.
func lifecycleString(v string, prior ovhtypes.TfStringValue) ovhtypes.TfStringValue {
	if v == "" && (prior.IsNull() || prior.IsUnknown() || prior.ValueString() != "") {
		return ovhtypes.TfStringValue{StringValue: types.StringNull()}
	}
	return ovhtypes.TfStringValue{StringValue: types.StringValue(v)}
}

func lifecycleBool(v bool, prior types.Bool) types.Bool {
	if !v && (prior.IsNull() || prior.IsUnknown() || prior.ValueBool()) {
		return types.BoolNull()
	}
	return types.BoolValue(v)
}
