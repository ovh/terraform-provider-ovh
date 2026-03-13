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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudStorageBlockVolumeResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudStorageBlockVolumeResource)(nil)
	_ resource.ResourceWithImportState = (*cloudStorageBlockVolumeResource)(nil)
)

func NewCloudStorageBlockVolumeResource() resource.Resource {
	return &cloudStorageBlockVolumeResource{}
}

type cloudStorageBlockVolumeResource struct {
	config *Config
}

func (r *cloudStorageBlockVolumeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_block_volume"
}

func (r *cloudStorageBlockVolumeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var volumeMutableAttrs = MutableAttrs{
	Strings: []string{"name", "volume_type"},
	Int64s:  []string{"size"},
	Bools:   []string{"bootable"},
}

func (r *cloudStorageBlockVolumeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a block storage volume in a public cloud project using the publicCloud API.",
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
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Volume name",
				MarkdownDescription: "Volume name",
			},
			"size": schema.Int64Attribute{
				Required:            true,
				Description:         "Size of the volume in GB",
				MarkdownDescription: "Size of the volume in GB",
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region where the volume will be created",
				MarkdownDescription: "Region where the volume will be created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"volume_type": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Volume type (CLASSIC, CLASSIC_LUKS, CLASSIC_MULTIATTACH, HIGH_SPEED, HIGH_SPEED_GEN2, HIGH_SPEED_GEN2_LUKS, HIGH_SPEED_LUKS). Can be changed after creation (triggers online retype).",
				MarkdownDescription: "Volume type (`CLASSIC`, `CLASSIC_LUKS`, `CLASSIC_MULTIATTACH`, `HIGH_SPEED`, `HIGH_SPEED_GEN2`, `HIGH_SPEED_GEN2_LUKS`, `HIGH_SPEED_LUKS`). Can be changed after creation (triggers online retype).",
			},
			"bootable": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Whether the volume is bootable",
				MarkdownDescription: "Whether the volume is bootable",
			},
			"create_from": schema.SingleNestedAttribute{
				Optional:            true,
				Description:         "Source to create the volume from (e.g. restore from a backup). Changing this value recreates the resource.",
				MarkdownDescription: "Source to create the volume from (e.g. restore from a backup). **Changing this value recreates the resource.**",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"backup_id": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Optional:            true,
						Description:         "Identifier of a backup to restore the volume from",
						MarkdownDescription: "Identifier of a backup to restore the volume from",
					},
				},
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Volume ID",
				MarkdownDescription: "Volume ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(volumeMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the volume",
				MarkdownDescription: "Creation date of the volume",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the volume",
				MarkdownDescription: "Last update date of the volume",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(volumeMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Volume readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				MarkdownDescription: "Volume readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			// targetSpec fields are exposed at root: `name`, `size`, `region`
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the block storage volume",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(volumeMutableAttrs),
				},
				Attributes: map[string]schema.Attribute{
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
					"name": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Volume name",
						MarkdownDescription: "Volume name",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						Description:         "Size of the volume in GB",
						MarkdownDescription: "Size of the volume in GB",
					},
					"volume_type": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Volume type (CLASSIC, CLASSIC_LUKS, CLASSIC_MULTIATTACH, HIGH_SPEED, HIGH_SPEED_GEN2, HIGH_SPEED_GEN2_LUKS, HIGH_SPEED_LUKS)",
						MarkdownDescription: "Volume type (`CLASSIC`, `CLASSIC_LUKS`, `CLASSIC_MULTIATTACH`, `HIGH_SPEED`, `HIGH_SPEED_GEN2`, `HIGH_SPEED_GEN2_LUKS`, `HIGH_SPEED_LUKS`)",
					},
					"bootable": schema.BoolAttribute{
						Computed:            true,
						Description:         "Whether the volume is bootable",
						MarkdownDescription: "Whether the volume is bootable",
					},
					"status": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Volume status (AVAILABLE, IN_USE, CREATING, DELETING, ATTACHING, DETACHING, EXTENDING, ERROR, ERROR_DELETING, ERROR_BACKING_UP, ERROR_RESTORING, ERROR_EXTENDING)",
						MarkdownDescription: "Volume status (`AVAILABLE`, `IN_USE`, `CREATING`, `DELETING`, `ATTACHING`, `DETACHING`, `EXTENDING`, `ERROR`, `ERROR_DELETING`, `ERROR_BACKING_UP`, `ERROR_RESTORING`, `ERROR_EXTENDING`)",
					},
				},
			},
		},
	}
}

func (r *cloudStorageBlockVolumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<volume_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudStorageBlockVolumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudStorageBlockVolumeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/volume"

	var responseData CloudStorageBlockVolumeAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for volume to be READY
	_, err := r.waitForVolumeReady(ctx, data.ServiceName.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for block storage to be ready",
			err.Error(),
		)
		return
	}

	// Read final state
	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/volume/" + url.PathEscape(responseData.Id)
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

func (r *cloudStorageBlockVolumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudStorageBlockVolumeModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/volume/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudStorageBlockVolumeAPIResponse
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

func (r *cloudStorageBlockVolumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudStorageBlockVolumeModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/volume/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudStorageBlockVolumeAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for volume to be READY
	_, err := r.waitForVolumeReady(ctx, data.ServiceName.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for block storage to be ready after update",
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

func (r *cloudStorageBlockVolumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudStorageBlockVolumeModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/volume/" + url.PathEscape(data.Id.ValueString())

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
			res := &CloudStorageBlockVolumeAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/volume/" + url.PathEscape(data.Id.ValueString())
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
			"Error waiting for block storage volume to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudStorageBlockVolumeResource) waitForVolumeReady(ctx context.Context, serviceName, volumeId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudStorageBlockVolumeAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/storage/block/volume/" + url.PathEscape(volumeId)
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
