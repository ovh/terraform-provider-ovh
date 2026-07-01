package ovh

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudNetworkPrivateSubnetResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudNetworkPrivateSubnetResource)(nil)
	_ resource.ResourceWithImportState = (*cloudNetworkPrivateSubnetResource)(nil)
)

func NewCloudNetworkPrivateVrackSubnetResource() resource.Resource {
	return &cloudNetworkPrivateSubnetResource{}
}

type cloudNetworkPrivateSubnetResource struct {
	config *Config
}

func (r *cloudNetworkPrivateSubnetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_network_private_vrack_subnet"
}

func (r *cloudNetworkPrivateSubnetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudNetworkPrivateSubnetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = CloudNetworkPrivateSubnetResourceSchema(ctx)
}

func (r *cloudNetworkPrivateSubnetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<network_id>/<id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)
}

func (r *cloudNetworkPrivateSubnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData CloudNetworkPrivateSubnetModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	base := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/network/" + url.PathEscape(data.NetworkId.ValueString()) + "/subnet"

	if err := r.config.OVHClient.Post(base, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", base),
			err.Error(),
		)
		return
	}

	// The response body does not carry the identity fields, re-attach them from
	// the request-side model so the ID is tracked even if the wait below fails.
	responseData.ServiceName = data.ServiceName
	responseData.NetworkId = data.NetworkId
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Wait for the subnet to be READY. The backend recover step rewrites the
	// targetSpec and checksum during this wait.
	itemEndpoint := base + "/" + url.PathEscape(responseData.Id.ValueString())
	if err := helpers.WaitForAPIv2ResourceStatusReady(ctx, r.config.OVHClient, itemEndpoint); err != nil {
		resp.Diagnostics.AddError("Error waiting for subnet to be ready", err.Error())
		return
	}

	// Read the final, post-recover state.
	var finalData CloudNetworkPrivateSubnetModel
	if err := r.config.OVHClient.Get(itemEndpoint, &finalData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", itemEndpoint),
			err.Error(),
		)
		return
	}

	finalData.ServiceName = data.ServiceName
	finalData.NetworkId = data.NetworkId

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalData)...)
}

func (r *cloudNetworkPrivateSubnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData CloudNetworkPrivateSubnetModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/network/" + url.PathEscape(data.NetworkId.ValueString()) + "/subnet/" + url.PathEscape(data.Id.ValueString())

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

	// Overwrite state with the API truth (re-attaching identity). We deliberately
	// do not call the fill-only MergeWith: overwriting target_spec from the API is
	// what surfaces drift.
	responseData.ServiceName = data.ServiceName
	responseData.NetworkId = data.NetworkId

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *cloudNetworkPrivateSubnetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, state CloudNetworkPrivateSubnetModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use the latest post-READY checksum from state for the concurrency control.
	planData.Checksum = state.Checksum

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(state.ServiceName.ValueString()) + "/network/" + url.PathEscape(state.NetworkId.ValueString()) + "/subnet/" + url.PathEscape(state.Id.ValueString())

	if err := r.config.OVHClient.Put(endpoint, planData.ToUpdate(), nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for the subnet to be READY after the update.
	if err := helpers.WaitForAPIv2ResourceStatusReady(ctx, r.config.OVHClient, endpoint); err != nil {
		resp.Diagnostics.AddError("Error waiting for subnet to be ready after update", err.Error())
		return
	}

	// Read the final, post-recover state.
	var finalData CloudNetworkPrivateSubnetModel
	if err := r.config.OVHClient.Get(endpoint, &finalData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	finalData.ServiceName = state.ServiceName
	finalData.NetworkId = state.NetworkId

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &finalData)...)
}

func (r *cloudNetworkPrivateSubnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudNetworkPrivateSubnetModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/network/" + url.PathEscape(data.NetworkId.ValueString()) + "/subnet/" + url.PathEscape(data.Id.ValueString())

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
			res := &CloudNetworkPrivateSubnetModel{}
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}
			return res, res.ResourceStatus.ValueString(), nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for subnet to be deleted",
			err.Error(),
		)
	}
}

func CloudNetworkPrivateSubnetResourceSchema(ctx context.Context) schema.Schema {
	attrs := map[string]schema.Attribute{
		"checksum": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Computed hash representing the current target specification value",
			MarkdownDescription: "Computed hash representing the current target specification value",
		},
		"created_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Creation date of the subnet",
			MarkdownDescription: "Creation date of the subnet",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"current_state": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{
				"allocation_pools": schema.ListNestedAttribute{
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"end": schema.StringAttribute{
								CustomType:          ovhtypes.TfStringType{},
								Computed:            true,
								Description:         "End IP address of the pool",
								MarkdownDescription: "End IP address of the pool",
							},
							"start": schema.StringAttribute{
								CustomType:          ovhtypes.TfStringType{},
								Computed:            true,
								Description:         "Start IP address of the pool",
								MarkdownDescription: "Start IP address of the pool",
							},
						},
						CustomType: CurrentStateAllocationPoolsType{
							ObjectType: types.ObjectType{
								AttrTypes: CurrentStateAllocationPoolsValue{}.AttributeTypes(ctx),
							},
						},
					},
					CustomType:          ovhtypes.NewTfListNestedType[CurrentStateAllocationPoolsValue](ctx),
					Computed:            true,
					Description:         "Allocation pools",
					MarkdownDescription: "Allocation pools",
				},
				"cidr": schema.StringAttribute{
					CustomType:          ovhtypes.TfStringType{},
					Computed:            true,
					Description:         "CIDR address range",
					MarkdownDescription: "CIDR address range",
				},
				"description": schema.StringAttribute{
					CustomType:          ovhtypes.TfStringType{},
					Computed:            true,
					Description:         "Subnet description",
					MarkdownDescription: "Subnet description",
				},
				"dhcp_enabled": schema.BoolAttribute{
					CustomType:          ovhtypes.TfBoolType{},
					Computed:            true,
					Description:         "Whether DHCP is enabled",
					MarkdownDescription: "Whether DHCP is enabled",
				},
				"dns_nameservers": schema.ListAttribute{
					CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
					Computed:            true,
					Description:         "DNS nameservers",
					MarkdownDescription: "DNS nameservers",
				},
				"gateway_ip": schema.StringAttribute{
					CustomType:          ovhtypes.TfStringType{},
					Computed:            true,
					Description:         "Default gateway IP",
					MarkdownDescription: "Default gateway IP",
				},
				"host_routes": schema.ListNestedAttribute{
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"destination": schema.StringAttribute{
								CustomType:          ovhtypes.TfStringType{},
								Computed:            true,
								Description:         "Destination CIDR",
								MarkdownDescription: "Destination CIDR",
							},
							"next_hop": schema.StringAttribute{
								CustomType:          ovhtypes.TfStringType{},
								Computed:            true,
								Description:         "Next hop IP address",
								MarkdownDescription: "Next hop IP address",
							},
						},
						CustomType: CurrentStateHostRoutesType{
							ObjectType: types.ObjectType{
								AttrTypes: CurrentStateHostRoutesValue{}.AttributeTypes(ctx),
							},
						},
					},
					CustomType:          ovhtypes.NewTfListNestedType[CurrentStateHostRoutesValue](ctx),
					Computed:            true,
					Description:         "Static host routes",
					MarkdownDescription: "Static host routes",
				},
				"location": schema.SingleNestedAttribute{
					Attributes: map[string]schema.Attribute{
						"availability_zone": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Computed:            true,
							Description:         "Availability zone within the region",
							MarkdownDescription: "Availability zone within the region",
						},
						"region": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Computed:            true,
							Description:         "Region code",
							MarkdownDescription: "Region code",
						},
					},
					CustomType: CurrentStateLocationType{
						ObjectType: types.ObjectType{
							AttrTypes: CurrentStateLocationValue{}.AttributeTypes(ctx),
						},
					},
					Computed:            true,
					Description:         "Location details",
					MarkdownDescription: "Location details",
				},
				"name": schema.StringAttribute{
					CustomType:          ovhtypes.TfStringType{},
					Computed:            true,
					Description:         "Subnet name",
					MarkdownDescription: "Subnet name",
				},
			},
			CustomType: CloudNetworkPrivateSubnetCurrentStateType{
				ObjectType: types.ObjectType{
					AttrTypes: CloudNetworkPrivateSubnetCurrentStateValue{}.AttributeTypes(ctx),
				},
			},
			Computed:            true,
			Description:         "Current state of the subnet",
			MarkdownDescription: "Current state of the subnet",
		},
		"current_tasks": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"errors": schema.ListNestedAttribute{
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"message": schema.StringAttribute{
									CustomType:          ovhtypes.TfStringType{},
									Computed:            true,
									Description:         "Error description",
									MarkdownDescription: "Error description",
								},
							},
							CustomType: CurrentTasksErrorsType{
								ObjectType: types.ObjectType{
									AttrTypes: CurrentTasksErrorsValue{}.AttributeTypes(ctx),
								},
							},
						},
						CustomType:          ovhtypes.NewTfListNestedType[CurrentTasksErrorsValue](ctx),
						Computed:            true,
						Description:         "Errors that occured on the task",
						MarkdownDescription: "Errors that occured on the task",
					},
					"id": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Identifier of the current task",
						MarkdownDescription: "Identifier of the current task",
					},
					"link": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Link to the task details",
						MarkdownDescription: "Link to the task details",
					},
					"status": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Current global status of the current task",
						MarkdownDescription: "Current global status of the current task",
					},
					"type": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Type of the current task",
						MarkdownDescription: "Type of the current task",
					},
				},
				CustomType: CloudNetworkPrivateSubnetCurrentTasksType{
					ObjectType: types.ObjectType{
						AttrTypes: CloudNetworkPrivateSubnetCurrentTasksValue{}.AttributeTypes(ctx),
					},
				},
			},
			CustomType:          ovhtypes.NewTfListNestedType[CloudNetworkPrivateSubnetCurrentTasksValue](ctx),
			Computed:            true,
			Description:         "Ongoing asynchronous tasks related to the subnet",
			MarkdownDescription: "Ongoing asynchronous tasks related to the subnet",
		},
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Unique identifier of the subnet",
			MarkdownDescription: "Unique identifier of the subnet",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"network_id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Network ID",
			MarkdownDescription: "Network ID",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name",
			MarkdownDescription: "Service name",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"resource_status": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Subnet readiness in the system",
			MarkdownDescription: "Subnet readiness in the system",
			PlanModifiers: []planmodifier.String{
				OutOfSyncPlanModifier(),
			},
		},
		"allocation_pools": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"end": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Computed:            true,
						Description:         "End IP address of the pool",
						MarkdownDescription: "End IP address of the pool",
					},
					"start": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Computed:            true,
						Description:         "Start IP address of the pool",
						MarkdownDescription: "Start IP address of the pool",
					},
				},
				CustomType: TargetSpecAllocationPoolsType{
					ObjectType: types.ObjectType{
						AttrTypes: TargetSpecAllocationPoolsValue{}.AttributeTypes(ctx),
					},
				},
			},
			CustomType:          ovhtypes.NewTfListNestedType[TargetSpecAllocationPoolsValue](ctx),
			Optional:            true,
			Computed:            true,
			Description:         "IP address allocation pools",
			MarkdownDescription: "IP address allocation pools",
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
		},
		"cidr": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "CIDR address range (immutable after creation)",
			MarkdownDescription: "CIDR address range (immutable after creation)",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"description": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Description:         "Subnet description",
			MarkdownDescription: "Subnet description",
		},
		"dhcp_enabled": schema.BoolAttribute{
			CustomType:          ovhtypes.TfBoolType{},
			Optional:            true,
			Computed:            true,
			Description:         "Whether DHCP is enabled",
			MarkdownDescription: "Whether DHCP is enabled",
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		},
		"dns_nameservers": schema.ListAttribute{
			CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
			Optional:            true,
			Description:         "DNS nameservers for the subnet",
			MarkdownDescription: "DNS nameservers for the subnet",
		},
		"gateway_ip": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Computed:            true,
			Description:         "Default gateway IP address",
			MarkdownDescription: "Default gateway IP address",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"region": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Region code",
			MarkdownDescription: "Region code",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
		"availability_zone": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Optional:            true,
			Description:         "Availability zone within the region",
			MarkdownDescription: "Availability zone within the region",
		},
		"name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Subnet name",
			MarkdownDescription: "Subnet name",
		},
		"updated_at": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Computed:            true,
			Description:         "Last update date of the subnet",
			MarkdownDescription: "Last update date of the subnet",
		},
	}

	return schema.Schema{
		Description: "",
		Attributes:  attrs,
	}
}

type CloudNetworkPrivateSubnetModel struct {
	Checksum       ovhtypes.TfStringValue                                                 `tfsdk:"checksum" json:"checksum"`
	CreatedAt      ovhtypes.TfStringValue                                                 `tfsdk:"created_at" json:"createdAt"`
	CurrentState   CloudNetworkPrivateSubnetCurrentStateValue                             `tfsdk:"current_state" json:"currentState"`
	CurrentTasks   ovhtypes.TfListNestedValue[CloudNetworkPrivateSubnetCurrentTasksValue] `tfsdk:"current_tasks" json:"currentTasks"`
	Id             ovhtypes.TfStringValue                                                 `tfsdk:"id" json:"id"`
	NetworkId      ovhtypes.TfStringValue                                                 `tfsdk:"network_id" json:"networkId"`
	ServiceName    ovhtypes.TfStringValue                                                 `tfsdk:"service_name" json:"projectId"`
	ResourceStatus ovhtypes.TfStringValue                                                 `tfsdk:"resource_status" json:"resourceStatus"`
	UpdatedAt      ovhtypes.TfStringValue                                                 `tfsdk:"updated_at" json:"updatedAt"`

	// Target-spec fields, lifted from the former nested target_spec object to
	// the resource root. They carry json:"-" because the targetSpec object is
	// (un)marshaled explicitly in MarshalJSON/UnmarshalJSON via the
	// CloudNetworkPrivateSubnetTargetSpecValue DTO, never at the JSON root.
	AllocationPools  ovhtypes.TfListNestedValue[TargetSpecAllocationPoolsValue] `tfsdk:"allocation_pools" json:"-"`
	Cidr             ovhtypes.TfStringValue                                     `tfsdk:"cidr" json:"-"`
	Description      ovhtypes.TfStringValue                                     `tfsdk:"description" json:"-"`
	DhcpEnabled      ovhtypes.TfBoolValue                                       `tfsdk:"dhcp_enabled" json:"-"`
	DnsNameservers   ovhtypes.TfListNestedValue[ovhtypes.TfStringValue]         `tfsdk:"dns_nameservers" json:"-"`
	GatewayIp        ovhtypes.TfStringValue                                     `tfsdk:"gateway_ip" json:"-"`
	Region           ovhtypes.TfStringValue                                     `tfsdk:"region" json:"-"`
	AvailabilityZone ovhtypes.TfStringValue                                     `tfsdk:"availability_zone" json:"-"`
	Name             ovhtypes.TfStringValue                                     `tfsdk:"name" json:"-"`
}

func (v *CloudNetworkPrivateSubnetModel) MergeWith(other *CloudNetworkPrivateSubnetModel) {

	if (v.Checksum.IsUnknown() || v.Checksum.IsNull()) && !other.Checksum.IsUnknown() {
		v.Checksum = other.Checksum
	}

	if (v.CreatedAt.IsUnknown() || v.CreatedAt.IsNull()) && !other.CreatedAt.IsUnknown() {
		v.CreatedAt = other.CreatedAt
	}

	if (v.CurrentState.IsUnknown() || v.CurrentState.IsNull()) && !other.CurrentState.IsUnknown() {
		v.CurrentState = other.CurrentState
	}

	if (v.CurrentTasks.IsUnknown() || v.CurrentTasks.IsNull()) && !other.CurrentTasks.IsUnknown() {
		v.CurrentTasks = other.CurrentTasks
	}

	if (v.Id.IsUnknown() || v.Id.IsNull()) && !other.Id.IsUnknown() {
		v.Id = other.Id
	}

	if (v.NetworkId.IsUnknown() || v.NetworkId.IsNull()) && !other.NetworkId.IsUnknown() {
		v.NetworkId = other.NetworkId
	}

	if (v.ServiceName.IsUnknown() || v.ServiceName.IsNull()) && !other.ServiceName.IsUnknown() {
		v.ServiceName = other.ServiceName
	}

	if (v.ResourceStatus.IsUnknown() || v.ResourceStatus.IsNull()) && !other.ResourceStatus.IsUnknown() {
		v.ResourceStatus = other.ResourceStatus
	}

	if (v.UpdatedAt.IsUnknown() || v.UpdatedAt.IsNull()) && !other.UpdatedAt.IsUnknown() {
		v.UpdatedAt = other.UpdatedAt
	}

	// Lifted target-spec fields.
	if (v.AllocationPools.IsUnknown() || v.AllocationPools.IsNull()) && !other.AllocationPools.IsUnknown() {
		v.AllocationPools = other.AllocationPools
	}

	if (v.Cidr.IsUnknown() || v.Cidr.IsNull()) && !other.Cidr.IsUnknown() {
		v.Cidr = other.Cidr
	}

	if (v.Description.IsUnknown() || v.Description.IsNull()) && !other.Description.IsUnknown() {
		v.Description = other.Description
	}

	if (v.DhcpEnabled.IsUnknown() || v.DhcpEnabled.IsNull()) && !other.DhcpEnabled.IsUnknown() {
		v.DhcpEnabled = other.DhcpEnabled
	}

	if (v.DnsNameservers.IsUnknown() || v.DnsNameservers.IsNull()) && !other.DnsNameservers.IsUnknown() {
		v.DnsNameservers = other.DnsNameservers
	}

	if (v.GatewayIp.IsUnknown() || v.GatewayIp.IsNull()) && !other.GatewayIp.IsUnknown() {
		v.GatewayIp = other.GatewayIp
	}

	if (v.Region.IsUnknown() || v.Region.IsNull()) && !other.Region.IsUnknown() {
		v.Region = other.Region
	}

	if (v.AvailabilityZone.IsUnknown() || v.AvailabilityZone.IsNull()) && !other.AvailabilityZone.IsUnknown() {
		v.AvailabilityZone = other.AvailabilityZone
	}

	if (v.Name.IsUnknown() || v.Name.IsNull()) && !other.Name.IsUnknown() {
		v.Name = other.Name
	}

}

// targetSpecValue assembles the lifted root fields back into the
// CloudNetworkPrivateSubnetTargetSpecValue DTO used for JSON (un)marshaling.
//
// location is create-only (immutable): the API rejects targetSpec.location on
// PUT. We therefore build a known location only when region is set. On update,
// the DTO's ToUpdate() drops the location, which clears the lifted region/AZ
// (see ToUpdate), so region is null here and the location stays null — and the
// DTO MarshalJSON omits a null location.
func (v CloudNetworkPrivateSubnetModel) targetSpecValue() CloudNetworkPrivateSubnetTargetSpecValue {
	location := NewTargetSpecLocationValueNull()
	if !v.Region.IsNull() && !v.Region.IsUnknown() {
		location = TargetSpecLocationValue{
			AvailabilityZone: v.AvailabilityZone,
			Region:           v.Region,
			state:            attr.ValueStateKnown,
		}
	}

	return CloudNetworkPrivateSubnetTargetSpecValue{
		AllocationPools: v.AllocationPools,
		Cidr:            v.Cidr,
		Description:     v.Description,
		DhcpEnabled:     v.DhcpEnabled,
		DnsNameservers:  v.DnsNameservers,
		GatewayIp:       v.GatewayIp,
		Location:        location,
		Name:            v.Name,
		state:           attr.ValueStateKnown,
	}
}

func (v CloudNetworkPrivateSubnetModel) ToCreate() *CloudNetworkPrivateSubnetModel {
	res := &CloudNetworkPrivateSubnetModel{}

	created := v.targetSpecValue().ToCreate()
	res.AllocationPools = created.AllocationPools
	res.Cidr = created.Cidr
	res.Description = created.Description
	res.DhcpEnabled = created.DhcpEnabled
	res.DnsNameservers = created.DnsNameservers
	res.GatewayIp = created.GatewayIp
	res.Region = created.Location.Region
	res.AvailabilityZone = created.Location.AvailabilityZone
	res.Name = created.Name

	return res
}

func (v CloudNetworkPrivateSubnetModel) ToUpdate() *CloudNetworkPrivateSubnetModel {
	res := &CloudNetworkPrivateSubnetModel{}

	if !v.Checksum.IsUnknown() {
		res.Checksum = v.Checksum
	}

	updated := v.targetSpecValue().ToUpdate()
	res.AllocationPools = updated.AllocationPools
	res.Cidr = updated.Cidr
	res.Description = updated.Description
	res.DhcpEnabled = updated.DhcpEnabled
	res.DnsNameservers = updated.DnsNameservers
	res.GatewayIp = updated.GatewayIp
	res.Region = updated.Location.Region
	res.AvailabilityZone = updated.Location.AvailabilityZone
	res.Name = updated.Name

	return res
}

func (v *CloudNetworkPrivateSubnetModel) MarshalJSON() ([]byte, error) {
	toMarshal := map[string]any{}
	if !v.Checksum.IsNull() && !v.Checksum.IsUnknown() {
		toMarshal["checksum"] = v.Checksum
	}

	// The lifted root fields are re-nested under targetSpec for the wire format.
	// The DTO's own MarshalJSON only emits the non-null/non-unknown sub-fields.
	targetSpec := v.targetSpecValue()
	toMarshal["targetSpec"] = &targetSpec

	return json.Marshal(toMarshal)
}

// UnmarshalJSON decodes the API envelope into the flat model. The nested
// targetSpec object is decoded into the lifted root fields via the
// CloudNetworkPrivateSubnetTargetSpecValue DTO; currentState / currentTasks
// keep their own value-type decoders.
func (v *CloudNetworkPrivateSubnetModel) UnmarshalJSON(data []byte) error {
	var shadow struct {
		Checksum       ovhtypes.TfStringValue                                                 `json:"checksum"`
		CreatedAt      ovhtypes.TfStringValue                                                 `json:"createdAt"`
		CurrentState   CloudNetworkPrivateSubnetCurrentStateValue                             `json:"currentState"`
		CurrentTasks   ovhtypes.TfListNestedValue[CloudNetworkPrivateSubnetCurrentTasksValue] `json:"currentTasks"`
		Id             ovhtypes.TfStringValue                                                 `json:"id"`
		NetworkId      ovhtypes.TfStringValue                                                 `json:"networkId"`
		ServiceName    ovhtypes.TfStringValue                                                 `json:"projectId"`
		ResourceStatus ovhtypes.TfStringValue                                                 `json:"resourceStatus"`
		UpdatedAt      ovhtypes.TfStringValue                                                 `json:"updatedAt"`
		TargetSpec     *CloudNetworkPrivateSubnetTargetSpecValue                              `json:"targetSpec"`
	}

	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}

	v.Checksum = shadow.Checksum
	v.CreatedAt = shadow.CreatedAt
	v.CurrentState = shadow.CurrentState
	v.CurrentTasks = shadow.CurrentTasks
	v.Id = shadow.Id
	v.NetworkId = shadow.NetworkId
	v.ServiceName = shadow.ServiceName
	v.ResourceStatus = shadow.ResourceStatus
	v.UpdatedAt = shadow.UpdatedAt

	if shadow.TargetSpec != nil {
		v.AllocationPools = shadow.TargetSpec.AllocationPools
		v.Cidr = shadow.TargetSpec.Cidr
		v.Description = shadow.TargetSpec.Description
		v.DhcpEnabled = shadow.TargetSpec.DhcpEnabled
		v.DnsNameservers = shadow.TargetSpec.DnsNameservers
		v.GatewayIp = shadow.TargetSpec.GatewayIp
		v.Region = shadow.TargetSpec.Location.Region
		v.AvailabilityZone = shadow.TargetSpec.Location.AvailabilityZone
		v.Name = shadow.TargetSpec.Name
	}

	return nil
}
