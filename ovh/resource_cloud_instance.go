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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudInstanceResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudInstanceResource)(nil)
	_ resource.ResourceWithImportState = (*cloudInstanceResource)(nil)
)

func NewCloudInstanceResource() resource.Resource {
	return &cloudInstanceResource{}
}

type cloudInstanceResource struct {
	config *Config
}

func (r *cloudInstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance"
}

func (r *cloudInstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var instanceMutableAttrs = MutableAttrs{
	Strings:           []string{"name", "flavor_id", "image_id", "ssh_key_name"},
	Lists:             []string{"networks"},
	CustomStringLists: []string{"volume_ids"},
}

func (r *cloudInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates an instance in a public cloud project using the new publicCloud API.",
		Attributes: map[string]schema.Attribute{
			// Required attributes
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Instance name",
				MarkdownDescription: "Instance name",
			},
			"flavor_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Flavor ID for the instance",
				MarkdownDescription: "Flavor ID for the instance",
			},
			"image_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Image ID for the instance",
				MarkdownDescription: "Image ID for the instance",
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region where the instance will be created",
				MarkdownDescription: "Region where the instance will be created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Optional attributes
			"availability_zone": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "The availability zone where the instance will be created",
				MarkdownDescription: "The availability zone where the instance will be created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"networks": schema.ListNestedAttribute{
				Optional:    true,
				Description: "List of networks to attach to the instance (targetSpec.networks).",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Description:         "Network ID to attach.",
							MarkdownDescription: "Network ID to attach.",
						},
						"public": schema.BoolAttribute{
							Optional:            true,
							Description:         "Whether to attach the public network.",
							MarkdownDescription: "Whether to attach the public network.",
						},
						"subnet_id": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Description:         "Subnet ID for the private network, if any.",
							MarkdownDescription: "Subnet ID for the private network, if any.",
						},
						"floating_ip_id": schema.StringAttribute{
							CustomType:          ovhtypes.TfStringType{},
							Optional:            true,
							Description:         "Floating IP ID for the network.",
							MarkdownDescription: "Floating IP ID for the network.",
						},
					},
				},
			},

			"volume_ids": schema.ListAttribute{
				CustomType:          ovhtypes.NewTfListNestedType[ovhtypes.TfStringValue](ctx),
				Optional:            true,
				Description:         "List of volume IDs to attach to the instance",
				MarkdownDescription: "List of volume IDs to attach to the instance",
			},
			"ssh_key_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "SSH key name to associate with the instance",
				MarkdownDescription: "SSH key name to associate with the instance",
			},
			// Computed attributes
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Instance ID",
				MarkdownDescription: "Instance ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(instanceMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the instance",
				MarkdownDescription: "Creation date of the instance",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the instance",
				MarkdownDescription: "Last update date of the instance",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(instanceMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Instance readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, SUSPENDED, UPDATING)",
				MarkdownDescription: "Instance readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, SUSPENDED, UPDATING)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},

			// Current state (computed, read-only)
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the instance",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(instanceMutableAttrs),
				},
				Attributes: map[string]schema.Attribute{
					"flavor": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Flavor details",
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Flavor identifier",
							},
							"name": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Flavor name",
							},
							"vcpus": schema.Int64Attribute{
								Computed:    true,
								Description: "Number of vCPUs",
							},
							"ram": schema.Int64Attribute{
								Computed:    true,
								Description: "RAM in MB",
							},
							"disk": schema.Int64Attribute{
								Computed:    true,
								Description: "Local disk size in GB",
							},
						},
					},
					"image": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Image details",
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Image identifier",
							},
							"name": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Image name",
							},
							"status": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Image status",
							},
						},
					},
					"name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Instance name",
					},
					"host_id": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Host identifier",
					},
					"ssh_key_name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Associated SSH key name",
					},
					"project_id": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Project identifier",
					},
					"user_id": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "User identifier",
					},
					"networks": schema.ListNestedAttribute{
						Computed:    true,
						Description: "List of instance networks",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Network identifier",
								},
								"public": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether the network is public",
								},
								"subnet_id": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Associated subnet identifier, if any",
								},
								"gateway_id": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Gateway identifier",
								},
								"floating_ip_id": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Floating IP identifier",
								},
								"addresses": schema.ListNestedAttribute{
									Computed:    true,
									Description: "IP addresses",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"ip": schema.StringAttribute{
												CustomType:  ovhtypes.TfStringType{},
												Computed:    true,
												Description: "IP address",
											},
											"mac": schema.StringAttribute{
												CustomType:  ovhtypes.TfStringType{},
												Computed:    true,
												Description: "MAC address",
											},
											"type": schema.StringAttribute{
												CustomType:  ovhtypes.TfStringType{},
												Computed:    true,
												Description: "Address type (fixed, floating, ...)",
											},
											"version": schema.Int64Attribute{
												Computed:    true,
												Description: "IP version",
											},
										},
									},
								},
							},
						},
					},
					"volumes": schema.ListNestedAttribute{
						Computed:    true,
						Description: "Attached block volumes",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Volume identifier",
								},
								"name": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Volume name",
								},
								"size": schema.Int64Attribute{
									Computed:    true,
									Description: "Volume size in GB",
								},
							},
						},
					},
					"security_groups": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Description: "Security groups attached to the instance",
					},
				},
			},
		},
	}
}

func (r *cloudInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<instance_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudInstanceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the request payload
	createPayload := data.ToCreate()

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance"

	var responseData CloudInstanceAPIResponse
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

	// Wait for instance to be READY
	_, err := r.waitForInstanceReady(ctx, data.ServiceName.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for instance to be ready",
			err.Error(),
		)
		return
	}

	// Read the final state
	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(responseData.Id)
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Update model with response
	data.MergeWith(ctx, &responseData)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudInstanceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudInstanceAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudInstanceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the update payload with checksum for optimistic concurrency
	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudInstanceAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for instance to be READY
	_, err := r.waitForInstanceReady(ctx, data.ServiceName.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for instance to be ready after update",
			err.Error(),
		)
		return
	}

	// Read the final state
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	planData.MergeWith(ctx, &responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *cloudInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudInstanceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(data.Id.ValueString())

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
			res := &CloudInstanceAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/instance/" + url.PathEscape(data.Id.ValueString())
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
			"Error waiting for instance to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudInstanceResource) waitForInstanceReady(ctx context.Context, serviceName, instanceId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudInstanceAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/compute/instance/" + url.PathEscape(instanceId)
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
