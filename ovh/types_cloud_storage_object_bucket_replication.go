package ovh

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

type CloudStorageObjectBucketReplicationModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketId    ovhtypes.TfStringValue `tfsdk:"bucket_id"`
	Id          ovhtypes.TfStringValue `tfsdk:"id"`
	Rules       types.List             `tfsdk:"rules"`
}

type storageObjectBucketReplicationRuleModel struct {
	Id                      ovhtypes.TfStringValue                         `tfsdk:"id"`
	Status                  ovhtypes.TfStringValue                         `tfsdk:"status"`
	Priority                types.Int64                                    `tfsdk:"priority"`
	DeleteMarkerReplication ovhtypes.TfStringValue                         `tfsdk:"delete_marker_replication"`
	Destination             storageObjectBucketReplicationDestinationModel `tfsdk:"destination"`
	Filter                  *storageObjectBucketReplicationFilterModel     `tfsdk:"filter"`
}

type storageObjectBucketReplicationDestinationModel struct {
	Name         ovhtypes.TfStringValue `tfsdk:"name"`
	Region       ovhtypes.TfStringValue `tfsdk:"region"`
	StorageClass ovhtypes.TfStringValue `tfsdk:"storage_class"`
}

type storageObjectBucketReplicationFilterModel struct {
	Prefix ovhtypes.TfStringValue `tfsdk:"prefix"`
	Tags   types.Map              `tfsdk:"tags"`
}

type CloudStorageObjectBucketReplicationAPIFilter struct {
	Prefix string            `json:"prefix,omitempty"`
	Tags   map[string]string `json:"tags,omitempty"`
}

type CloudStorageObjectBucketReplicationAPIDestination struct {
	Name         string `json:"name"`
	Region       string `json:"region"`
	StorageClass string `json:"storageClass,omitempty"`
}

type CloudStorageObjectBucketReplicationAPIRule struct {
	Id                      string                                            `json:"id"`
	Status                  string                                            `json:"status"`
	Priority                int64                                             `json:"priority"`
	DeleteMarkerReplication string                                            `json:"deleteMarkerReplication"`
	Destination             CloudStorageObjectBucketReplicationAPIDestination `json:"destination"`
	Filter                  *CloudStorageObjectBucketReplicationAPIFilter     `json:"filter,omitempty"`
}

// CloudStorageObjectBucketReplicationAPIRules is both the PUT payload and the GET/PUT response.
type CloudStorageObjectBucketReplicationAPIRules struct {
	Rules []CloudStorageObjectBucketReplicationAPIRule `json:"rules"`
}

func storageObjectBucketReplicationRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                        ovhtypes.TfStringType{},
		"status":                    ovhtypes.TfStringType{},
		"priority":                  types.Int64Type,
		"delete_marker_replication": ovhtypes.TfStringType{},
		"destination": types.ObjectType{AttrTypes: map[string]attr.Type{
			"name":          ovhtypes.TfStringType{},
			"region":        ovhtypes.TfStringType{},
			"storage_class": ovhtypes.TfStringType{},
		}},
		"filter": types.ObjectType{AttrTypes: map[string]attr.Type{
			"prefix": ovhtypes.TfStringType{},
			"tags":   types.MapType{ElemType: ovhtypes.TfStringType{}},
		}},
	}
}

func (m *CloudStorageObjectBucketReplicationModel) rules(ctx context.Context) ([]storageObjectBucketReplicationRuleModel, diag.Diagnostics) {
	var rules []storageObjectBucketReplicationRuleModel
	if m.Rules.IsNull() || m.Rules.IsUnknown() {
		return rules, nil
	}
	diags := m.Rules.ElementsAs(ctx, &rules, false)
	return rules, diags
}

func (m *CloudStorageObjectBucketReplicationModel) ToPut(ctx context.Context) (*CloudStorageObjectBucketReplicationAPIRules, diag.Diagnostics) {
	rules, diags := m.rules(ctx)
	if diags.HasError() {
		return nil, diags
	}

	payload := &CloudStorageObjectBucketReplicationAPIRules{Rules: make([]CloudStorageObjectBucketReplicationAPIRule, 0, len(rules))}
	for _, rule := range rules {
		apiRule := CloudStorageObjectBucketReplicationAPIRule{
			Id:                      rule.Id.ValueString(),
			Status:                  rule.Status.ValueString(),
			Priority:                rule.Priority.ValueInt64(),
			DeleteMarkerReplication: rule.DeleteMarkerReplication.ValueString(),
			Destination: CloudStorageObjectBucketReplicationAPIDestination{
				Name:         rule.Destination.Name.ValueString(),
				Region:       rule.Destination.Region.ValueString(),
				StorageClass: rule.Destination.StorageClass.ValueString(),
			},
		}

		if f := rule.Filter; f != nil {
			apiRule.Filter = &CloudStorageObjectBucketReplicationAPIFilter{Prefix: f.Prefix.ValueString()}
			if !f.Tags.IsNull() && !f.Tags.IsUnknown() {
				tags := make(map[string]ovhtypes.TfStringValue, len(f.Tags.Elements()))
				diags.Append(f.Tags.ElementsAs(ctx, &tags, false)...)
				apiRule.Filter.Tags = make(map[string]string, len(tags))
				for k, v := range tags {
					apiRule.Filter.Tags[k] = v.ValueString()
				}
			}
		}

		payload.Rules = append(payload.Rules, apiRule)
	}

	return payload, diags
}

// MergeWith rebuilds rules from the API, ordered and shaped like prior (plan or state): S3 may
// return rules in another order, and drops empty values the config may spell out.
func (m *CloudStorageObjectBucketReplicationModel) MergeWith(ctx context.Context, response *CloudStorageObjectBucketReplicationAPIRules) diag.Diagnostics {
	m.Id = m.BucketId

	prior, diags := m.rules(ctx)
	if diags.HasError() {
		return diags
	}

	priorByID := make(map[string]*storageObjectBucketReplicationRuleModel, len(prior))
	for i := range prior {
		priorByID[prior[i].Id.ValueString()] = &prior[i]
	}

	apiByID := make(map[string]CloudStorageObjectBucketReplicationAPIRule, len(response.Rules))
	for _, r := range response.Rules {
		apiByID[r.Id] = r
	}
	ordered := make([]CloudStorageObjectBucketReplicationAPIRule, 0, len(response.Rules))
	seen := make(map[string]bool, len(response.Rules))
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

	rules := make([]storageObjectBucketReplicationRuleModel, 0, len(ordered))
	for _, r := range ordered {
		rules = append(rules, buildStorageObjectBucketReplicationRule(r, priorByID[r.Id]))
	}

	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: storageObjectBucketReplicationRuleAttrTypes()}, rules)
	diags.Append(d...)
	m.Rules = list
	return diags
}

func buildStorageObjectBucketReplicationRule(r CloudStorageObjectBucketReplicationAPIRule, prior *storageObjectBucketReplicationRuleModel) storageObjectBucketReplicationRuleModel {
	if prior == nil {
		prior = &storageObjectBucketReplicationRuleModel{}
	}

	// The API rebuilds the region from the bucket listing (upper case, or empty when ambiguous), so a matching prior is kept.
	region := nullableTfString(r.Destination.Region)
	priorRegion := prior.Destination.Region
	if !priorRegion.IsUnknown() && (r.Destination.Region == "" || strings.EqualFold(priorRegion.ValueString(), r.Destination.Region)) {
		region = priorRegion
	}

	return storageObjectBucketReplicationRuleModel{
		Id:                      ovhtypes.TfStringValue{StringValue: types.StringValue(r.Id)},
		Status:                  ovhtypes.TfStringValue{StringValue: types.StringValue(r.Status)},
		Priority:                types.Int64Value(r.Priority),
		DeleteMarkerReplication: ovhtypes.TfStringValue{StringValue: types.StringValue(r.DeleteMarkerReplication)},
		Destination: storageObjectBucketReplicationDestinationModel{
			Name:         ovhtypes.TfStringValue{StringValue: types.StringValue(r.Destination.Name)},
			Region:       region,
			StorageClass: nullableTfString(r.Destination.StorageClass),
		},
		Filter: buildStorageObjectBucketReplicationFilter(r.Filter, prior.Filter),
	}
}

// The API sends a rule without filter to S3 as an empty prefix, and reads it back as an empty filter.
func buildStorageObjectBucketReplicationFilter(f *CloudStorageObjectBucketReplicationAPIFilter, prior *storageObjectBucketReplicationFilterModel) *storageObjectBucketReplicationFilterModel {
	empty := f == nil || (f.Prefix == "" && len(f.Tags) == 0)
	priorEmpty := prior != nil && prior.Prefix.ValueString() == "" && len(prior.Tags.Elements()) == 0
	if empty && !priorEmpty {
		return nil
	}
	if f == nil {
		f = &CloudStorageObjectBucketReplicationAPIFilter{}
	}
	if prior == nil {
		prior = &storageObjectBucketReplicationFilterModel{}
	}

	filter := &storageObjectBucketReplicationFilterModel{
		Prefix: lifecycleString(f.Prefix, prior.Prefix),
		Tags:   types.MapNull(ovhtypes.TfStringType{}),
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
