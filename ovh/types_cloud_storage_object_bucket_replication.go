package ovh

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// CloudStorageObjectBucketReplicationModel is the Terraform state model for the replication resource.
type CloudStorageObjectBucketReplicationModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketName  ovhtypes.TfStringValue `tfsdk:"bucket_name"`
	Id          ovhtypes.TfStringValue `tfsdk:"id"`
	Rules       types.List             `tfsdk:"rules"`
}

// API request/response types

type CloudBucketReplicationRequest struct {
	Rules []CloudBucketReplicationRule `json:"rules"`
}

type CloudBucketReplicationResponse struct {
	Rules []CloudBucketReplicationRule `json:"rules"`
}

type CloudBucketReplicationRule struct {
	ID                      string                            `json:"id"`
	Status                  string                            `json:"status"`
	Priority                int64                             `json:"priority"`
	Filter                  *CloudBucketReplicationFilter     `json:"filter,omitempty"`
	Destination             CloudBucketReplicationDestination `json:"destination"`
	DeleteMarkerReplication string                            `json:"deleteMarkerReplication"`
}

type CloudBucketReplicationFilter struct {
	Prefix string            `json:"prefix,omitempty"`
	Tags   map[string]string `json:"tags,omitempty"`
}

type CloudBucketReplicationDestination struct {
	Name         string `json:"name"`
	Region       string `json:"region"`
	StorageClass string `json:"storageClass,omitempty"`
}

// Attribute type helpers

func replicationFilterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"prefix": ovhtypes.TfStringType{},
		"tags":   types.MapType{ElemType: types.StringType},
	}
}

func replicationDestinationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          ovhtypes.TfStringType{},
		"region":        ovhtypes.TfStringType{},
		"storage_class": ovhtypes.TfStringType{},
	}
}

func replicationRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                        ovhtypes.TfStringType{},
		"status":                    ovhtypes.TfStringType{},
		"priority":                  types.Int64Type,
		"filter":                    types.ObjectType{AttrTypes: replicationFilterAttrTypes()},
		"destination":               types.ObjectType{AttrTypes: replicationDestinationAttrTypes()},
		"delete_marker_replication": ovhtypes.TfStringType{},
	}
}

// apiReplicationRulesToTFList converts API rules to a Terraform list, reordered to match the
// state list order (by rule ID) to prevent spurious diffs.
func apiReplicationRulesToTFList(apiRules []CloudBucketReplicationRule, stateRules types.List) types.List {
	ruleElemType := types.ObjectType{AttrTypes: replicationRuleAttrTypes()}

	ruleByID := make(map[string]CloudBucketReplicationRule, len(apiRules))
	for _, r := range apiRules {
		ruleByID[r.ID] = r
	}

	var ordered []CloudBucketReplicationRule
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

	for _, r := range apiRules {
		if _, remaining := ruleByID[r.ID]; remaining {
			ordered = append(ordered, r)
		}
	}

	// Build a map of state rule objects by ID for fallback values
	stateRuleByID := make(map[string]types.Object)
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
			stateRuleByID[idVal.ValueString()] = obj
		}
	}

	elems := make([]attr.Value, 0, len(ordered))
	for _, r := range ordered {
		stateObj, hasState := stateRuleByID[r.ID]
		if !hasState {
			stateObj = types.ObjectNull(replicationRuleAttrTypes())
		}
		elems = append(elems, apiReplicationRuleToTFObject(r, stateObj))
	}

	list, _ := types.ListValue(ruleElemType, elems)
	return list
}

func apiReplicationRuleToTFObject(r CloudBucketReplicationRule, stateRule types.Object) basetypes.ObjectValue {
	// filter — treat effectively empty filter as null (same pattern as lifecycle)
	filterObj := types.ObjectNull(replicationFilterAttrTypes())
	if r.Filter != nil && (r.Filter.Prefix != "" || len(r.Filter.Tags) > 0) {
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
		filterObj, _ = types.ObjectValue(replicationFilterAttrTypes(), map[string]attr.Value{
			"prefix": ovhtypes.TfStringValue{StringValue: types.StringValue(r.Filter.Prefix)},
			"tags":   tagsVal,
		})
	}

	// destination — preserve region from state when API returns empty
	regionVal := ovhtypes.TfStringValue{StringValue: types.StringValue(r.Destination.Region)}
	if r.Destination.Region == "" && !stateRule.IsNull() && !stateRule.IsUnknown() {
		if destObj, ok := stateRule.Attributes()["destination"].(types.Object); ok && !destObj.IsNull() {
			if stateRegion, ok := destObj.Attributes()["region"].(ovhtypes.TfStringValue); ok && !stateRegion.IsNull() {
				regionVal = stateRegion
			}
		}
	}

	// storage_class — null when empty
	storageClassVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	if r.Destination.StorageClass != "" {
		storageClassVal = ovhtypes.TfStringValue{StringValue: types.StringValue(r.Destination.StorageClass)}
	}

	destinationObj, _ := types.ObjectValue(replicationDestinationAttrTypes(), map[string]attr.Value{
		"name":          ovhtypes.TfStringValue{StringValue: types.StringValue(r.Destination.Name)},
		"region":        regionVal,
		"storage_class": storageClassVal,
	})

	obj, _ := types.ObjectValue(replicationRuleAttrTypes(), map[string]attr.Value{
		"id":                        ovhtypes.TfStringValue{StringValue: types.StringValue(r.ID)},
		"status":                    ovhtypes.TfStringValue{StringValue: types.StringValue(r.Status)},
		"priority":                  types.Int64Value(r.Priority),
		"filter":                    filterObj,
		"destination":               destinationObj,
		"delete_marker_replication": ovhtypes.TfStringValue{StringValue: types.StringValue(r.DeleteMarkerReplication)},
	})
	return obj
}

// tfListToAPIReplicationRules converts the Terraform rules list to the API slice.
func tfListToAPIReplicationRules(list types.List) []CloudBucketReplicationRule {
	if list.IsNull() || list.IsUnknown() {
		return []CloudBucketReplicationRule{}
	}

	rules := make([]CloudBucketReplicationRule, 0, len(list.Elements()))
	for _, elem := range list.Elements() {
		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}
		rules = append(rules, tfObjectToAPIReplicationRule(obj))
	}
	return rules
}

func tfObjectToAPIReplicationRule(obj types.Object) CloudBucketReplicationRule {
	attrs := obj.Attributes()
	rule := CloudBucketReplicationRule{}

	if v, ok := attrs["id"].(ovhtypes.TfStringValue); ok {
		rule.ID = v.ValueString()
	}
	if v, ok := attrs["status"].(ovhtypes.TfStringValue); ok {
		rule.Status = v.ValueString()
	}
	if v, ok := attrs["priority"].(basetypes.Int64Value); ok && !v.IsNull() {
		rule.Priority = v.ValueInt64()
	}
	if v, ok := attrs["delete_marker_replication"].(ovhtypes.TfStringValue); ok {
		rule.DeleteMarkerReplication = v.ValueString()
	}

	if filterObj, ok := attrs["filter"].(types.Object); ok && !filterObj.IsNull() && !filterObj.IsUnknown() {
		fa := filterObj.Attributes()
		filter := &CloudBucketReplicationFilter{}
		if v, ok := fa["prefix"].(ovhtypes.TfStringValue); ok && !v.IsNull() {
			filter.Prefix = v.ValueString()
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

	if destObj, ok := attrs["destination"].(types.Object); ok && !destObj.IsNull() && !destObj.IsUnknown() {
		da := destObj.Attributes()
		if v, ok := da["name"].(ovhtypes.TfStringValue); ok {
			rule.Destination.Name = v.ValueString()
		}
		if v, ok := da["region"].(ovhtypes.TfStringValue); ok {
			rule.Destination.Region = v.ValueString()
		}
		if v, ok := da["storage_class"].(ovhtypes.TfStringValue); ok && !v.IsNull() {
			rule.Destination.StorageClass = v.ValueString()
		}
	}

	return rule
}
