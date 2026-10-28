package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudPrivateEndpointResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudPrivateEndpointResource)(nil)
	_ resource.ResourceWithImportState = (*cloudPrivateEndpointResource)(nil)
)

func NewCloudPrivateEndpointResource() resource.Resource {
	return &cloudPrivateEndpointResource{}
}

type cloudPrivateEndpointResource struct {
	config *Config
}

func (r *cloudPrivateEndpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_private_endpoint"
}

func (r *cloudPrivateEndpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Only name and description are mutable; every other targetSpec field forces a
// replacement.
var privateEndpointMutableAttrs = MutableAttrs{
	Strings: []string{"name", "description"},
}

func (r *cloudPrivateEndpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a private endpoint in a public cloud project using the publicCloud API.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Service name of the resource representing the id of the cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used.",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
				PlanModifiers: []planmodifier.String{
					EnvDefaultString("OVH_CLOUD_PROJECT_SERVICE", true),
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Private endpoint name, unique in the project for a given subnet (case-insensitive)",
				MarkdownDescription: "Private endpoint name, unique in the project for a given subnet (case-insensitive)",
			},
			"description": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Private endpoint description",
				MarkdownDescription: "Private endpoint description",
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region where the private endpoint is created. Changing this value recreates the resource.",
				MarkdownDescription: "Region where the private endpoint is created. **Changing this value recreates the resource.**",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service": schema.SingleNestedAttribute{
				Required:            true,
				Description:         "Service to expose through the endpoint. Changing this value recreates the resource.",
				MarkdownDescription: "Service to expose through the endpoint. **Changing this value recreates the resource.**",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Required:            true,
						Description:         "Name of the service",
						MarkdownDescription: "Name of the service",
					},
					"tier": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Required:            true,
						Description:         "Tier of the service, as listed by the ovh_cloud_private_endpoint_service_availabilities data source",
						MarkdownDescription: "Tier of the service, as listed by the `ovh_cloud_private_endpoint_service_availabilities` data source",
					},
				},
			},
			"network_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Network to allocate the endpoint port in. At least one of network_id or subnet_id must be set; both may be set, and the subnet must then belong to the network. When omitted, the network owning the given subnet is resolved by the platform. Changing this value recreates the resource.",
				MarkdownDescription: "Network to allocate the endpoint port in. At least one of `network_id` or `subnet_id` must be set; both may be set, and the subnet must then belong to the network. When omitted, the network owning the given subnet is resolved by the platform. **Changing this value recreates the resource.**",
				// AtLeastOneOf reports its own path too, so declaring it here covers
				// subnet_id as well.
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.MatchRoot("subnet_id")),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Subnet to allocate the endpoint port in. At least one of network_id or subnet_id must be set; both may be set, and the subnet must then belong to the network. When omitted, the network must hold exactly one IPv4 subnet, which is resolved by the platform. Changing this value recreates the resource.",
				MarkdownDescription: "Subnet to allocate the endpoint port in. At least one of `network_id` or `subnet_id` must be set; both may be set, and the subnet must then belong to the network. When omitted, the network must hold exactly one IPv4 subnet, which is resolved by the platform. **Changing this value recreates the resource.**",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vip_address": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Private IP address to pin on the endpoint port. Requires subnet_id. When omitted, an address is allocated from the subnet allocation pool. Changing this value recreates the resource.",
				MarkdownDescription: "Private IP address to pin on the endpoint port. Requires `subnet_id`. When omitted, an address is allocated from the subnet allocation pool. **Changing this value recreates the resource.**",
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.MatchRoot("subnet_id")),
				},
				// UseStateForUnknown is load-bearing: MarkComputedNilsAsUnknown marks an
				// unconfigured Optional+Computed attribute unknown before plan modifiers
				// run, so without it the planned value never equals state and replacement
				// fires on every plan.
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_dns_enabled": schema.BoolAttribute{
				Required:            true,
				Description:         "Register the endpoint hostname in the public DNS zone. Changing this value recreates the resource.",
				MarkdownDescription: "Register the endpoint hostname in the public DNS zone. **Changing this value recreates the resource.**",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Unique identifier of the private endpoint",
				MarkdownDescription: "Unique identifier of the private endpoint",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(privateEndpointMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the private endpoint",
				MarkdownDescription: "Creation date of the private endpoint",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the private endpoint",
				MarkdownDescription: "Last update date of the private endpoint",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(privateEndpointMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Private endpoint readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				MarkdownDescription: "Private endpoint readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the private endpoint",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(privateEndpointMutableAttrs),
				},
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Private endpoint name",
					},
					"description": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Private endpoint description",
					},
					"endpoint_id": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Identifier of the endpoint on the service side, as pe-<id>",
					},
					"location": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Location details",
						Attributes: map[string]schema.Attribute{
							"region": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Region",
							},
						},
					},
					"service": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Service exposed through the endpoint",
						Attributes: map[string]schema.Attribute{
							"name": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Name of the service",
							},
							"tier": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Tier of the service",
							},
							"port": schema.Int64Attribute{
								Computed:    true,
								Description: "Port the service listens on through the endpoint",
							},
							"transport": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Transport protocol of the service",
							},
							"bandwidth": schema.Int64Attribute{
								Computed:    true,
								Description: "Bandwidth of the tier, in Mbit/s. Null when the catalog does not declare it",
							},
						},
					},
					"network": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Network the endpoint port is attached to",
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Network identifier",
							},
						},
					},
					"subnet": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Subnet the endpoint port is attached to",
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								CustomType:  ovhtypes.TfStringType{},
								Computed:    true,
								Description: "Subnet identifier",
							},
						},
					},
					"vip_address": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Private IP address of the endpoint in the customer network",
					},
					"hostname": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Fully qualified domain name of the endpoint. Null until the service publishes it. See public_dns_enabled for public resolvability",
					},
					"public_dns_enabled": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether the endpoint hostname is registered in the public DNS zone",
					},
				},
			},
		},
	}
}

func (r *cloudPrivateEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<private_endpoint_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudPrivateEndpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudPrivateEndpointModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/privateEndpoint"

	var responseData CloudPrivateEndpointAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save state immediately so the resource ID is tracked even if the workflow
	// fails before reaching READY.
	data.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if err := r.waitForPrivateEndpointReady(ctx, data.ServiceName.ValueString(), responseData.Id); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for private endpoint to be ready",
			err.Error(),
		)
		return
	}

	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/privateEndpoint/" + url.PathEscape(responseData.Id)
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

func (r *cloudPrivateEndpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudPrivateEndpointModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/privateEndpoint/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudPrivateEndpointAPIResponse
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

func (r *cloudPrivateEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudPrivateEndpointModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/privateEndpoint/" + url.PathEscape(data.Id.ValueString())

	// The checksum comes from state, not from a fresh read: it is a
	// compare-and-swap token, so sending the value Terraform last observed is
	// what makes the server reject a concurrent modification.
	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	var responseData CloudPrivateEndpointAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	if err := r.waitForPrivateEndpointReady(ctx, data.ServiceName.ValueString(), data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for private endpoint to be ready after update",
			err.Error(),
		)
		return
	}

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

func (r *cloudPrivateEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudPrivateEndpointModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/privateEndpoint/" + url.PathEscape(data.Id.ValueString())

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

	stateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (any, string, error) {
			res := &CloudPrivateEndpointAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/privateEndpoint/" + url.PathEscape(data.Id.ValueString())
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}
			if res.ResourceStatus == "ERROR" {
				return res, res.ResourceStatus, cloudResourceErrorFromTasks("private endpoint", data.Id.ValueString(), res.CurrentTasks)
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for private endpoint to be deleted",
			err.Error(),
		)
	}
}

// waitForPrivateEndpointReady polls until the endpoint reaches READY. An ERROR
// status is terminal: the ongoing task error is surfaced to the caller.
func (r *cloudPrivateEndpointResource) waitForPrivateEndpointReady(ctx context.Context, serviceName, id string) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (any, string, error) {
			res := &CloudPrivateEndpointAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/privateEndpoint/" + url.PathEscape(id)
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			if res.ResourceStatus == "ERROR" {
				return res, res.ResourceStatus, cloudResourceErrorFromTasks("private endpoint", id, res.CurrentTasks)
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
