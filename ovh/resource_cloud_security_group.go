package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudSecurityGroupResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudSecurityGroupResource)(nil)
	_ resource.ResourceWithImportState = (*cloudSecurityGroupResource)(nil)
)

func NewCloudSecurityGroupResource() resource.Resource {
	return &cloudSecurityGroupResource{}
}

type cloudSecurityGroupResource struct {
	config *Config
}

func (r *cloudSecurityGroupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_security_group"
}

func (r *cloudSecurityGroupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.config = config
}

var securityGroupMutableAttrs = MutableAttrs{
	Strings: []string{"name", "description"},
	Lists:   []string{"rule"},
}

// securityGroupStateRuleSchemaAttributes returns the schema attributes for a computed rule in current_state.
func securityGroupStateRuleSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{}, Computed: true,
			Description: "Rule ID", MarkdownDescription: "Rule ID",
		},
		"direction": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{}, Computed: true,
			Description: "Direction of the rule", MarkdownDescription: "Direction of the rule",
		},
		"ethernet_type": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{}, Computed: true,
			Description: "Ethernet type", MarkdownDescription: "Ethernet type",
		},
		"protocol": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{}, Computed: true,
			Description: "Protocol", MarkdownDescription: "Protocol",
		},
		"port_range_min": schema.Int64Attribute{
			Computed: true, Description: "Minimum port number", MarkdownDescription: "Minimum port number",
		},
		"port_range_max": schema.Int64Attribute{
			Computed: true, Description: "Maximum port number", MarkdownDescription: "Maximum port number",
		},
		"remote_group_id": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{}, Computed: true,
			Description: "Remote security group ID", MarkdownDescription: "Remote security group ID",
		},
		"remote_ip_prefix": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{}, Computed: true,
			Description: "Remote IP prefix", MarkdownDescription: "Remote IP prefix",
		},
		"description": schema.StringAttribute{
			CustomType: ovhtypes.TfStringType{}, Computed: true,
			Description: "Description of the rule", MarkdownDescription: "Description of the rule",
		},
	}
}

func (r *cloudSecurityGroupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a security group in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			// Required — immutable
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region where the security group will be created",
				MarkdownDescription: "Region where the security group will be created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Required — mutable
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Name of the security group",
				MarkdownDescription: "Name of the security group",
			},

			// Optional — mutable
			"description": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Description of the security group",
				MarkdownDescription: "Description of the security group",
			},
			"rule": schema.ListNestedAttribute{
				Optional:    true,
				Description: "List of security group rules",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"direction": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Required:            true,
							Description:         "Direction of the rule (INGRESS or EGRESS)",
							MarkdownDescription: "Direction of the rule (`INGRESS` or `EGRESS`)",
						},
						"ethernet_type": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Required:            true,
							Description:         "Ethernet type (IPV4 or IPV6)",
							MarkdownDescription: "Ethernet type (`IPV4` or `IPV6`)",
						},
						"protocol": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Description:         "Protocol (TCP, UDP, ICMP, etc.)",
							MarkdownDescription: "Protocol (`TCP`, `UDP`, `ICMP`, etc.)",
						},
						"port_range_min": schema.Int64Attribute{
							Optional:            true,
							Description:         "Minimum port number",
							MarkdownDescription: "Minimum port number",
						},
						"port_range_max": schema.Int64Attribute{
							Optional:            true,
							Description:         "Maximum port number",
							MarkdownDescription: "Maximum port number",
						},
						"remote_group_id": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Description:         "Remote security group ID",
							MarkdownDescription: "Remote security group ID",
						},
						"remote_ip_prefix": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Description:         "Remote IP prefix (CIDR notation)",
							MarkdownDescription: "Remote IP prefix (CIDR notation)",
						},
						"description": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Description:         "Description of the rule",
							MarkdownDescription: "Description of the rule",
						},
					},
				},
			},

			// Computed
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Security group ID",
				MarkdownDescription: "Security group ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(securityGroupMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the security group",
				MarkdownDescription: "Creation date of the security group",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the security group",
				MarkdownDescription: "Last update date of the security group",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(securityGroupMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Security group readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				MarkdownDescription: "Security group readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the security group",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(securityGroupMutableAttrs),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Name of the security group",
						MarkdownDescription: "Name of the security group",
					},
					"description": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Description of the security group",
						MarkdownDescription: "Description of the security group",
					},
					"region": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Region of the security group",
						MarkdownDescription: "Region of the security group",
					},
					"rules": schema.ListNestedAttribute{
						Computed:    true,
						Description: "User-specified security group rules with their IDs",
						NestedObject: schema.NestedAttributeObject{
							Attributes: securityGroupStateRuleSchemaAttributes(),
						},
					},
					"default_rules": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Default egress rules auto-created by OpenStack",
						NestedObject: schema.NestedAttributeObject{
							Attributes: securityGroupStateRuleSchemaAttributes(),
						},
					},
				},
			},
		},
	}
}

func (r *cloudSecurityGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<security_group_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudSecurityGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudSecurityGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/securityGroup"

	var responseData CloudSecurityGroupAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save state immediately so the resource ID is tracked even if the workflow fails
	data.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Wait for security group to be READY
	_, err := r.waitForSecurityGroupReady(ctx, data.ServiceName.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for security group to be ready",
			err.Error(),
		)
		return
	}

	// Read final state
	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/securityGroup/" + url.PathEscape(responseData.Id)
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudSecurityGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudSecurityGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/securityGroup/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudSecurityGroupAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudSecurityGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudSecurityGroupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/securityGroup/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudSecurityGroupAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for security group to be READY
	_, err := r.waitForSecurityGroupReady(ctx, data.ServiceName.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for security group to be ready after update",
			err.Error(),
		)
		return
	}

	// Read final state
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	planData.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *cloudSecurityGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudSecurityGroupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/securityGroup/" + url.PathEscape(data.Id.ValueString())

	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for deletion to complete
	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudSecurityGroupAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/securityGroup/" + url.PathEscape(data.Id.ValueString())
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for security group to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudSecurityGroupResource) waitForSecurityGroupReady(ctx context.Context, serviceName, securityGroupId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudSecurityGroupAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/securityGroup/" + url.PathEscape(securityGroupId)
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}
