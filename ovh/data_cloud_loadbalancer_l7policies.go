package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudLoadbalancerL7PoliciesDataSource)(nil)

func NewCloudLoadbalancerL7PoliciesDataSource() datasource.DataSource {
	return &cloudLoadbalancerL7PoliciesDataSource{}
}

type cloudLoadbalancerL7PoliciesDataSource struct {
	config *Config
}

// CloudLoadbalancerL7PoliciesModel is the model for the plural L7 policies data source.
type CloudLoadbalancerL7PoliciesModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	LoadbalancerId ovhtypes.TfStringValue `tfsdk:"loadbalancer_id"`
	ListenerId     ovhtypes.TfStringValue `tfsdk:"listener_id"`
	L7Policies     types.List             `tfsdk:"l7policies"`
}

func (d *cloudLoadbalancerL7PoliciesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_loadbalancer_l7policies"
}

func (d *cloudLoadbalancerL7PoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.config = config
}

// L7PolicyListItemAttrTypes returns the attribute types for a single L7 policy
// element of the plural data source list.
func L7PolicyListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                 ovhtypes.TfStringType{},
		"name":               ovhtypes.TfStringType{},
		"description":        ovhtypes.TfStringType{},
		"action":             ovhtypes.TfStringType{},
		"position":           types.Int64Type,
		"redirect_prefix":    ovhtypes.TfStringType{},
		"redirect_url":       ovhtypes.TfStringType{},
		"redirect_http_code": types.Int64Type,
		"redirect_pool_id":   ovhtypes.TfStringType{},
		"rules": types.ListType{
			ElemType: l7PolicyRuleElementType(),
		},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state": types.ObjectType{
			AttrTypes: L7PolicyCurrentStateAttrTypes(),
		},
	}
}

func (d *cloudLoadbalancerL7PoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to list the L7 policies of a load balancer listener in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"loadbalancer_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "ID of the load balancer",
			},
			"listener_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "ID of the listener",
			},
			"l7policies": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of L7 policies",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "L7 policy ID",
						},
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Name of the L7 policy",
						},
						"description": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Description of the L7 policy",
						},
						"action": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Action of the L7 policy (REDIRECT_PREFIX, REDIRECT_TO_POOL, REDIRECT_TO_URL, REJECT)",
						},
						"position": schema.Int64Attribute{
							Computed:    true,
							Description: "Position of the L7 policy",
						},
						"redirect_prefix": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Redirect prefix for REDIRECT_PREFIX action",
						},
						"redirect_url": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Redirect URL for REDIRECT_TO_URL action",
						},
						"redirect_http_code": schema.Int64Attribute{
							Computed:    true,
							Description: "HTTP redirect code (301, 302, 303, 307, 308)",
						},
						"redirect_pool_id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "ID of the pool for REDIRECT_TO_POOL action",
						},
						"rules": schema.ListNestedAttribute{
							Computed:    true,
							Description: "List of L7 rules for this policy",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Type of the L7 rule (COOKIE, FILE_TYPE, HEADER, HOST_NAME, PATH)",
									},
									"compare_type": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Comparison type (CONTAINS, ENDS_WITH, EQUAL_TO, REGEX, STARTS_WITH)",
									},
									"value": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Value to compare against",
									},
									"key": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Key for COOKIE and HEADER rule types",
									},
									"invert": schema.BoolAttribute{
										Computed:    true,
										Description: "Whether to invert the rule match",
									},
								},
							},
						},
						"checksum": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Computed hash representing the current target specification value",
						},
						"created_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Creation date of the L7 policy",
						},
						"updated_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Last update date of the L7 policy",
						},
						"resource_status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "L7 policy readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
						},
						"current_state": loadbalancerL7PolicyDataSourceCurrentStateSchema(),
					},
				},
			},
		},
	}
}

func (d *cloudLoadbalancerL7PoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudLoadbalancerL7PoliciesModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/loadbalancer/" + url.PathEscape(data.LoadbalancerId.ValueString()) +
		"/listener/" + url.PathEscape(data.ListenerId.ValueString()) +
		"/l7policy"

	var responseData []CloudLoadbalancerL7PolicyAPIResponse
	if err := d.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	itemObjType := types.ObjectType{AttrTypes: L7PolicyListItemAttrTypes()}
	items := make([]attr.Value, 0, len(responseData))
	for i := range responseData {
		items = append(items, buildL7PolicyListItemObject(ctx, &responseData[i]))
	}

	data.L7Policies = types.ListValueMust(itemObjType, items)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// buildL7PolicyListItemObject builds a single L7 policy list element from an API
// response, reusing the resource helpers for nested objects.
func buildL7PolicyListItemObject(ctx context.Context, response *CloudLoadbalancerL7PolicyAPIResponse) basetypes.ObjectValue {
	nameVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	descVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	actionVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	positionVal := types.Int64Null()
	redirectPrefixVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	redirectUrlVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	redirectHttpCodeVal := types.Int64Null()
	redirectPoolIdVal := ovhtypes.TfStringValue{StringValue: types.StringNull()}
	rulesVal := types.ListNull(l7PolicyRuleElementType())

	if spec := response.TargetSpec; spec != nil {
		nameVal = lbStringOrNull(spec.Name)
		descVal = lbStringOrNull(spec.Description)
		actionVal = ovhtypes.TfStringValue{StringValue: types.StringValue(spec.Action)}
		redirectPrefixVal = lbStringOrNull(spec.RedirectPrefix)
		redirectUrlVal = lbStringOrNull(spec.RedirectUrl)

		if spec.Position != nil {
			positionVal = types.Int64Value(int64(*spec.Position))
		}

		if spec.RedirectHttpCode != nil {
			redirectHttpCodeVal = types.Int64Value(int64(*spec.RedirectHttpCode))
		}

		if spec.RedirectPool != nil && spec.RedirectPool.ID != "" {
			redirectPoolIdVal = ovhtypes.TfStringValue{StringValue: types.StringValue(spec.RedirectPool.ID)}
		}

		rulesVal = buildL7PolicyRulesListFromTargetSpec(spec.Rules)
	}

	currentStateVal := types.ObjectNull(L7PolicyCurrentStateAttrTypes())
	if response.CurrentState != nil {
		currentStateVal = buildL7PolicyCurrentStateObject(ctx, response.CurrentState)
	}

	obj, _ := types.ObjectValue(
		L7PolicyListItemAttrTypes(),
		map[string]attr.Value{
			"id":                 ovhtypes.TfStringValue{StringValue: types.StringValue(response.Id)},
			"name":               nameVal,
			"description":        descVal,
			"action":             actionVal,
			"position":           positionVal,
			"redirect_prefix":    redirectPrefixVal,
			"redirect_url":       redirectUrlVal,
			"redirect_http_code": redirectHttpCodeVal,
			"redirect_pool_id":   redirectPoolIdVal,
			"rules":              rulesVal,
			"checksum":           ovhtypes.TfStringValue{StringValue: types.StringValue(response.Checksum)},
			"created_at":         ovhtypes.TfStringValue{StringValue: types.StringValue(response.CreatedAt)},
			"updated_at":         ovhtypes.TfStringValue{StringValue: types.StringValue(response.UpdatedAt)},
			"resource_status":    ovhtypes.TfStringValue{StringValue: types.StringValue(response.ResourceStatus)},
			"current_state":      currentStateVal,
		},
	)

	return obj
}
