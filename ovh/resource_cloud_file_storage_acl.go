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
	_ resource.Resource                = (*cloudFileStorageAclResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudFileStorageAclResource)(nil)
	_ resource.ResourceWithImportState = (*cloudFileStorageAclResource)(nil)
)

func NewCloudFileStorageAclResource() resource.Resource {
	return &cloudFileStorageAclResource{}
}

type cloudFileStorageAclResource struct {
	config *Config
}

func (r *cloudFileStorageAclResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_file_storage_acl"
}

func (r *cloudFileStorageAclResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var fileStorageAclMutableAttrs = MutableAttrs{
	Strings: []string{"access_level"},
}

func (r *cloudFileStorageAclResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates an access rule (ACL) on a public cloud file storage share, controlling which IP addresses can access it.",
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
			"share_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the file storage share the access rule applies to",
				MarkdownDescription: "ID of the file storage share the access rule applies to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"access_to": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "IP address or CIDR allowed to access the file storage share",
				MarkdownDescription: "IP address or CIDR allowed to access the file storage share",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"access_level": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Access level granted (READ_WRITE, READ_ONLY)",
				MarkdownDescription: "Access level granted (`READ_WRITE`, `READ_ONLY`)",
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Access rule ID",
				MarkdownDescription: "Access rule ID",
				// On access_level update the backend denies the old Manila rule
				// and grants a new one with a fresh id, so id changes.
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(fileStorageAclMutableAttrs),
				},
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(fileStorageAclMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the access rule",
				MarkdownDescription: "Creation date of the access rule",
				// The recreated rule carries a fresh createdAt on access_level update.
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(fileStorageAclMutableAttrs),
				},
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the access rule",
				MarkdownDescription: "Last update date of the access rule",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(fileStorageAclMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Access rule readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				MarkdownDescription: "Access rule readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current observed state of the access rule from the infrastructure",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(fileStorageAclMutableAttrs),
				},
				Attributes: map[string]schema.Attribute{
					"access_to": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "IP address or CIDR allowed to access the file storage share",
						MarkdownDescription: "IP address or CIDR allowed to access the file storage share",
					},
					"access_level": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Access level granted",
						MarkdownDescription: "Access level granted",
					},
					"state": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Current state of the access rule (ACTIVE, APPLYING, DENYING, ERROR)",
						MarkdownDescription: "Current state of the access rule (`ACTIVE`, `APPLYING`, `DENYING`, `ERROR`)",
					},
					"created_at": schema.StringAttribute{
						CustomType:          ovhtypes.TfStringType{},
						Computed:            true,
						Description:         "Creation date of the access rule",
						MarkdownDescription: "Creation date of the access rule",
					},
				},
			},
		},
	}
}

func (r *cloudFileStorageAclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<share_id>/<acl_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("share_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)
}

func fileStorageAclBaseEndpoint(serviceName, shareId string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/storage/file/share/" + url.PathEscape(shareId) + "/acl"
}

func (r *cloudFileStorageAclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudFileStorageAclModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := fileStorageAclBaseEndpoint(data.ServiceName.ValueString(), data.ShareId.ValueString())

	var responseData CloudFileStorageAclAPIResponse
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

	// Wait for the access rule to be READY (currentState.state ACTIVE)
	_, err := r.waitForFileStorageAclReady(ctx, data.ServiceName.ValueString(), data.ShareId.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for file storage access rule to be ready",
			err.Error(),
		)
		return
	}

	// Read final state
	getEndpoint := endpoint + "/" + url.PathEscape(responseData.Id)
	if err := r.config.OVHClient.Get(getEndpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", getEndpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudFileStorageAclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudFileStorageAclModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := fileStorageAclBaseEndpoint(data.ServiceName.ValueString(), data.ShareId.ValueString()) + "/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudFileStorageAclAPIResponse
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

// Update sends the accessLevel change to the current (state) acl ID. The
// backend implements this by deleting and recreating the underlying Manila
// access rule, so the PUT response carries a NEW resource ID — the old ID
// stops existing. We must not assume the ID is stable: capture the new ID
// from the PUT response and use it for the ready-poll, the final GET, and
// the stored state.
func (r *cloudFileStorageAclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudFileStorageAclModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	baseEndpoint := fileStorageAclBaseEndpoint(data.ServiceName.ValueString(), data.ShareId.ValueString())
	oldEndpoint := baseEndpoint + "/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudFileStorageAclAPIResponse
	if err := r.config.OVHClient.Put(oldEndpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", oldEndpoint),
			err.Error(),
		)
		return
	}

	// responseData.Id may differ from data.Id.ValueString() — always use the
	// ID the API just returned from here on.
	newAclId := responseData.Id

	// Persist the NEW id immediately so a wait timeout/error below doesn't leave
	// Terraform tracking the stale old id (whose row no longer exists → 404 →
	// orphaned rule).
	planData.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)

	_, err := r.waitForFileStorageAclReady(ctx, data.ServiceName.ValueString(), data.ShareId.ValueString(), newAclId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for file storage access rule to be ready after update",
			err.Error(),
		)
		return
	}

	// Read final state from the (possibly new) resource ID
	newEndpoint := baseEndpoint + "/" + url.PathEscape(newAclId)
	if err := r.config.OVHClient.Get(newEndpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", newEndpoint),
			err.Error(),
		)
		return
	}

	planData.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)
}

func (r *cloudFileStorageAclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudFileStorageAclModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	baseEndpoint := fileStorageAclBaseEndpoint(data.ServiceName.ValueString(), data.ShareId.ValueString())
	endpoint := baseEndpoint + "/" + url.PathEscape(data.Id.ValueString())

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
			res := &CloudFileStorageAclAPIResponse{}
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
			"Error waiting for file storage access rule to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudFileStorageAclResource) waitForFileStorageAclReady(ctx context.Context, serviceName, shareId, aclId string) (any, error) {
	endpoint := fileStorageAclBaseEndpoint(serviceName, shareId) + "/" + url.PathEscape(aclId)

	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (any, string, error) {
			res := &CloudFileStorageAclAPIResponse{}
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
