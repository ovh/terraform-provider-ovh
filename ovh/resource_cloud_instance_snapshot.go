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
	_ resource.Resource                = (*cloudInstanceSnapshotResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudInstanceSnapshotResource)(nil)
	_ resource.ResourceWithImportState = (*cloudInstanceSnapshotResource)(nil)
)

func NewCloudInstanceSnapshotResource() resource.Resource {
	return &cloudInstanceSnapshotResource{}
}

type cloudInstanceSnapshotResource struct {
	config *Config
}

func (r *cloudInstanceSnapshotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_snapshot"
}

func (r *cloudInstanceSnapshotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudInstanceSnapshotResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates an on-demand instance snapshot in a public cloud project using the publicCloud API.",
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
				Description:         "Snapshot name",
				MarkdownDescription: "Snapshot name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region where the snapshot will be created",
				MarkdownDescription: "Region where the snapshot will be created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the instance to snapshot",
				MarkdownDescription: "ID of the instance to snapshot",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Snapshot ID",
				MarkdownDescription: "Snapshot ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the snapshot",
				MarkdownDescription: "Creation date of the snapshot",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the snapshot",
				MarkdownDescription: "Last update date of the snapshot",
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Snapshot readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY)",
				MarkdownDescription: "Snapshot readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY)",
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the instance snapshot",
				Attributes: map[string]schema.Attribute{
					"instance": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Source instance reference",
						Attributes: map[string]schema.Attribute{
							"id": schema.StringAttribute{
								CustomType:          ovhtypes.TfStringType{},
								Computed:            true,
								Description:         "Instance unique identifier",
								MarkdownDescription: "Instance unique identifier",
							},
						},
					},
					"location": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Current location",
						Attributes: map[string]schema.Attribute{
							"region": schema.StringAttribute{
								CustomType:          ovhtypes.TfStringType{},
								Computed:            true,
								Description:         "Region",
								MarkdownDescription: "Region",
							},
						},
					},
					"min_disk": schema.Int64Attribute{
						Computed:            true,
						Description:         "Minimum disk size in GB required to boot",
						MarkdownDescription: "Minimum disk size in GB required to boot",
					},
					"min_ram": schema.Int64Attribute{
						Computed:            true,
						Description:         "Minimum RAM in MB required to boot",
						MarkdownDescription: "Minimum RAM in MB required to boot",
					},
					"name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Current snapshot name",
						MarkdownDescription: "Current snapshot name",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						Description:         "Image size in bytes",
						MarkdownDescription: "Image size in bytes",
					},
					"status": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Image status in the backend",
						MarkdownDescription: "Image status in the backend",
					},
					"visibility": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Image visibility",
						MarkdownDescription: "Image visibility",
					},
				},
			},
		},
	}
}

func (r *cloudInstanceSnapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<snapshot_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudInstanceSnapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudInstanceSnapshotModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/snapshot"

	var responseData CloudInstanceSnapshotAPIResponse
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

	_, err := r.waitForInstanceSnapshotReady(ctx, data.ServiceName.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for instance snapshot to be ready",
			err.Error(),
		)
		return
	}

	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/snapshot/" + url.PathEscape(responseData.Id)
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

func (r *cloudInstanceSnapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudInstanceSnapshotModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/snapshot/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudInstanceSnapshotAPIResponse
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

func (r *cloudInstanceSnapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes are immutable (RequiresReplace); the API exposes no update endpoint.
	resp.Diagnostics.AddError(
		"Update not supported",
		"ovh_cloud_instance_snapshot does not support in-place updates; all attributes require replacement.",
	)
}

func (r *cloudInstanceSnapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudInstanceSnapshotModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/snapshot/" + url.PathEscape(data.Id.ValueString())

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
		Refresh: func() (interface{}, string, error) {
			res := &CloudInstanceSnapshotAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/compute/snapshot/" + url.PathEscape(data.Id.ValueString())
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
			"Error waiting for instance snapshot to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudInstanceSnapshotResource) waitForInstanceSnapshotReady(ctx context.Context, serviceName, snapshotId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudInstanceSnapshotAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/compute/snapshot/" + url.PathEscape(snapshotId)
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
