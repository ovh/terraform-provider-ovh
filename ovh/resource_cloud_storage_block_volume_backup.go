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
	_ resource.Resource                = (*cloudStorageBlockVolumeBackupResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudStorageBlockVolumeBackupResource)(nil)
	_ resource.ResourceWithImportState = (*cloudStorageBlockVolumeBackupResource)(nil)
)

func NewCloudStorageBlockVolumeBackupResource() resource.Resource {
	return &cloudStorageBlockVolumeBackupResource{}
}

type cloudStorageBlockVolumeBackupResource struct {
	config *Config
}

func (r *cloudStorageBlockVolumeBackupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_block_volume_backup"
}

func (r *cloudStorageBlockVolumeBackupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var backupMutableAttrs = MutableAttrs{
	Strings: []string{"name", "description"},
}

func (r *cloudStorageBlockVolumeBackupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a block storage volume backup in a public cloud project using the publicCloud API.",
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
				Description:         "Backup name",
				MarkdownDescription: "Backup name",
			},
			"description": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Backup description",
				MarkdownDescription: "Backup description",
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region where the backup will be created",
				MarkdownDescription: "Region where the backup will be created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"volume_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the volume to back up",
				MarkdownDescription: "ID of the volume to back up",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Backup ID",
				MarkdownDescription: "Backup ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(backupMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the backup",
				MarkdownDescription: "Creation date of the backup",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the backup",
				MarkdownDescription: "Last update date of the backup",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(backupMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Backup readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				MarkdownDescription: "Backup readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the block storage volume backup",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(backupMutableAttrs),
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
						Description:         "Backup name",
						MarkdownDescription: "Backup name",
					},
					"description": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Backup description",
						MarkdownDescription: "Backup description",
					},
					"volume_id": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "ID of the backed-up volume",
						MarkdownDescription: "ID of the backed-up volume",
					},
					"size": schema.Int64Attribute{
						Computed:            true,
						Description:         "Size of the backup in GB",
						MarkdownDescription: "Size of the backup in GB",
					},
				},
			},
		},
	}
}

func (r *cloudStorageBlockVolumeBackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<backup_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudStorageBlockVolumeBackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudStorageBlockVolumeBackupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/backup"

	var responseData CloudStorageBlockVolumeBackupAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for backup to be READY
	_, err := r.waitForBackupReady(ctx, data.ServiceName.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for block storage backup to be ready",
			err.Error(),
		)
		return
	}

	// Read final state
	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/backup/" + url.PathEscape(responseData.Id)
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

func (r *cloudStorageBlockVolumeBackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudStorageBlockVolumeBackupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/backup/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudStorageBlockVolumeBackupAPIResponse
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

func (r *cloudStorageBlockVolumeBackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudStorageBlockVolumeBackupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/backup/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudStorageBlockVolumeBackupAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for backup to be READY
	_, err := r.waitForBackupReady(ctx, data.ServiceName.ValueString(), data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for block storage backup to be ready after update",
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

func (r *cloudStorageBlockVolumeBackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudStorageBlockVolumeBackupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/backup/" + url.PathEscape(data.Id.ValueString())

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
			res := &CloudStorageBlockVolumeBackupAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/block/backup/" + url.PathEscape(data.Id.ValueString())
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
			"Error waiting for block storage backup to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudStorageBlockVolumeBackupResource) waitForBackupReady(ctx context.Context, serviceName, backupId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudStorageBlockVolumeBackupAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/storage/block/backup/" + url.PathEscape(backupId)
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
