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
	_ resource.Resource                = (*cloudGatewayResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudGatewayResource)(nil)
	_ resource.ResourceWithImportState = (*cloudGatewayResource)(nil)
)

func NewCloudGatewayResource() resource.Resource {
	return &cloudGatewayResource{}
}

type cloudGatewayResource struct {
	config *Config
}

func (r *cloudGatewayResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_gateway"
}

func (r *cloudGatewayResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var gatewayMutableAttrs = MutableAttrs{
	Strings:           []string{"name", "description"},
	Objects:           []string{"external_gateway"},
	CustomStringLists: []string{"subnet_ids"},
}

func (r *cloudGatewayResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a gateway in a public cloud project.",
		Attributes: map[string]schema.Attribute{
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
				Description:         "Region where the gateway will be created",
				MarkdownDescription: "Region where the gateway will be created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"availability_zone": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Availability zone for the gateway",
				MarkdownDescription: "Availability zone for the gateway",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Gateway name",
				MarkdownDescription: "Gateway name",
			},
			"description": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Gateway description",
				MarkdownDescription: "Gateway description",
			},
			"external_gateway": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "External gateway configuration",
				MarkdownDescription: "External gateway configuration",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Required:            true,
						Description:         "Whether the external gateway is enabled",
						MarkdownDescription: "Whether the external gateway is enabled",
					},
					"model": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "External gateway sizing model (S, M, L, XL, 2XL, 3XL). Required when enabled is true.",
						MarkdownDescription: "External gateway sizing model (`S`, `M`, `L`, `XL`, `2XL`, `3XL`). Required when enabled is true.",
					},
				},
			},
			"subnet_ids": schema.ListAttribute{
				CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
				Optional:            true,
				Description:         "List of subnet IDs to attach as router interfaces",
				MarkdownDescription: "List of subnet IDs to attach as router interfaces",
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Gateway ID",
				MarkdownDescription: "Gateway ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(gatewayMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the gateway",
				MarkdownDescription: "Creation date of the gateway",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the gateway",
				MarkdownDescription: "Last update date of the gateway",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(gatewayMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Gateway readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				MarkdownDescription: "Gateway readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the gateway",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(gatewayMutableAttrs),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Gateway name",
					},
					"description": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Gateway description",
					},
					"status": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "OpenStack router status (ACTIVE, BUILD, DOWN, ERROR)",
					},
					"external_gateway": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "External gateway configuration",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Computed:    true,
								Description: "Whether the external gateway is enabled",
							},
							"model": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "External gateway sizing model",
							},
						},
					},
					"external_ip": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "External IP address assigned to the gateway",
					},
					"subnets": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Currently attached subnets",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Subnet ID",
								},
							},
						},
					},
					"region": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Region",
					},
					"availability_zone": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Availability zone",
					},
				},
			},
		},
	}
}

func (r *cloudGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<gateway_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudGatewayModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/gateway"

	var responseData CloudGatewayAPIResponse
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

	// Wait for gateway to be READY
	_, err := r.waitForGatewayReady(ctx, data.ServiceName.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for gateway to be ready",
			err.Error(),
		)
		return
	}

	// Read final state
	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/gateway/" + url.PathEscape(responseData.Id)
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

func (r *cloudGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudGatewayModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/gateway/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudGatewayAPIResponse
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

func (r *cloudGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudGatewayModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/gateway/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudGatewayAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for gateway to be READY
	_, err := r.waitForGatewayReady(ctx, data.ServiceName.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for gateway to be ready after update",
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

func (r *cloudGatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudGatewayModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/gateway/" + url.PathEscape(data.Id.ValueString())

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
			res := &CloudGatewayAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/gateway/" + url.PathEscape(data.Id.ValueString())
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
			"Error waiting for gateway to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudGatewayResource) waitForGatewayReady(ctx context.Context, serviceName, gatewayId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudGatewayAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/gateway/" + url.PathEscape(gatewayId)
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
