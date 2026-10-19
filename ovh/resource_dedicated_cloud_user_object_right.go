package ovh

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.ResourceWithConfigure   = (*dedicatedCloudUserObjectRightResource)(nil)
	_ resource.ResourceWithImportState = (*dedicatedCloudUserObjectRightResource)(nil)
)

func NewDedicatedCloudUserObjectRightResource() resource.Resource {
	return &dedicatedCloudUserObjectRightResource{}
}

type dedicatedCloudUserObjectRightResource struct {
	config *Config
}

func (r *dedicatedCloudUserObjectRightResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_cloud_user_object_right"
}

func (r *dedicatedCloudUserObjectRightResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dedicatedCloudUserObjectRightResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DedicatedCloudUserObjectRightResourceSchema(ctx)
}

func dedicatedCloudUserObjectRightsEndpoint(serviceName string, userId int64) string {
	return fmt.Sprintf("/dedicatedCloud/%s/user/%d/objectRight", url.PathEscape(serviceName), userId)
}

func dedicatedCloudUserObjectRightEndpoint(serviceName string, userId, objectRightId int64) string {
	return fmt.Sprintf("%s/%d", dedicatedCloudUserObjectRightsEndpoint(serviceName, userId), objectRightId)
}

func listDedicatedCloudUserObjectRightIds(ctx context.Context, r *dedicatedCloudUserObjectRightResource, serviceName string, userId int64) (map[int64]bool, error) {
	var ids []int64
	endpoint := dedicatedCloudUserObjectRightsEndpoint(serviceName, userId)
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &ids); err != nil {
		return nil, fmt.Errorf("error calling Get %s: %w", endpoint, err)
	}

	set := make(map[int64]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	return set, nil
}

// dedicatedCloudUserObjectRightLocks serializes object right creation per service/user.
// The create Task response carries no objectRightId, so the new id is resolved by diffing
// the object right list before and after creation; running that sequence concurrently for
// the same user would make the diffs overlap.
var dedicatedCloudUserObjectRightLocks sync.Map

func lockDedicatedCloudUserObjectRights(serviceName string, userId int64) func() {
	m, _ := dedicatedCloudUserObjectRightLocks.LoadOrStore(fmt.Sprintf("%s/%d", serviceName, userId), &sync.Mutex{})
	mu := m.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

// resolveCreatedObjectRightId finds the objectRightId created by the preceding POST by
// diffing the object right id list before and after creation. Creation is serialized
// per user within this provider, but rights may still be added out of band (console,
// other tooling), so every candidate - even a sole one - is validated against the
// requested right/type/vmware_object_id.
func (r *dedicatedCloudUserObjectRightResource) resolveCreatedObjectRightId(ctx context.Context, serviceName string, userId int64, before map[int64]bool, after map[int64]bool, want *DedicatedCloudUserObjectRightModel) (int64, error) {
	var candidates []int64
	for id := range after {
		if !before[id] {
			candidates = append(candidates, id)
		}
	}

	if len(candidates) == 0 {
		return 0, fmt.Errorf("object right was created but no new object_right_id was found on %s/user/%d - the object right list may not be consistent yet, try importing it manually", serviceName, userId)
	}

	var matches []int64
	for _, candidateId := range candidates {
		var candidate DedicatedCloudUserObjectRightModel
		endpoint := dedicatedCloudUserObjectRightEndpoint(serviceName, userId, candidateId)
		if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &candidate); err != nil {
			return 0, fmt.Errorf("error calling Get %s: %w", endpoint, err)
		}

		if candidate.Right.ValueString() == want.Right.ValueString() &&
			candidate.Type.ValueString() == want.Type.ValueString() &&
			candidate.VmwareObjectId.ValueString() == want.VmwareObjectId.ValueString() {
			matches = append(matches, candidateId)
		}
	}

	if len(matches) != 1 {
		return 0, fmt.Errorf(
			"found %d new object right(s) on %s/user/%d of which %d match the requested right/type/vmware_object_id - "+
				"the created object right could not be unambiguously identified (was an object right added outside Terraform at the same time?), import it manually",
			len(candidates), serviceName, userId, len(matches),
		)
	}

	return matches[0], nil
}

func (r *dedicatedCloudUserObjectRightResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData DedicatedCloudUserObjectRightModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()
	userId := data.UserId.ValueInt64()

	unlock := lockDedicatedCloudUserObjectRights(serviceName, userId)
	defer unlock()

	idsBefore, err := listDedicatedCloudUserObjectRightIds(ctx, r, serviceName, userId)
	if err != nil {
		resp.Diagnostics.AddError("Error listing existing object rights", err.Error())
		return
	}

	var task DedicatedCloudTask
	endpoint := dedicatedCloudUserObjectRightsEndpoint(serviceName, userId)
	if err := r.config.OVHClient.PostWithContext(ctx, endpoint, data.ToCreate(), &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	if _, err := waitForDedicatedCloudTask(ctx, r.config.OVHClient, serviceName, task.TaskId); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error waiting for object right creation on %s/user/%d", serviceName, userId),
			err.Error(),
		)
		return
	}

	idsAfter, err := listDedicatedCloudUserObjectRightIds(ctx, r, serviceName, userId)
	if err != nil {
		resp.Diagnostics.AddError("Error listing object rights after creation", err.Error())
		return
	}

	objectRightId, err := r.resolveCreatedObjectRightId(ctx, serviceName, userId, idsBefore, idsAfter, &data)
	if err != nil {
		resp.Diagnostics.AddError("Error resolving created object_right_id", err.Error())
		return
	}

	endpoint = dedicatedCloudUserObjectRightEndpoint(serviceName, userId, objectRightId)
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)
	responseData.ID = ovhtypes.NewTfStringValue(fmt.Sprintf("%s/%d/%d", serviceName, userId, objectRightId))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *dedicatedCloudUserObjectRightResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData DedicatedCloudUserObjectRightModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := dedicatedCloudUserObjectRightEndpoint(data.ServiceName.ValueString(), data.UserId.ValueInt64(), data.ObjectRightId.ValueInt64())
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
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

	// API values take precedence so out-of-band changes show up as drift
	responseData.MergeWith(&data)
	responseData.ID = data.ID

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *dedicatedCloudUserObjectRightResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not implemented", "the dedicatedCloud objectRight API has no update endpoint: all writable attributes require replacement, this func should never be called")
}

func (r *dedicatedCloudUserObjectRightResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DedicatedCloudUserObjectRightModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()
	userId := data.UserId.ValueInt64()

	var task DedicatedCloudTask
	endpoint := dedicatedCloudUserObjectRightEndpoint(serviceName, userId, data.ObjectRightId.ValueInt64())
	if err := r.config.OVHClient.DeleteWithContext(ctx, endpoint, &task); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	if _, err := waitForDedicatedCloudTask(ctx, r.config.OVHClient, serviceName, task.TaskId); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error waiting for object right deletion on %s/user/%d", serviceName, userId),
			err.Error(),
		)
	}
}

func (r *dedicatedCloudUserObjectRightResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	serviceName, userId, objectRightId, err := splitDedicatedCloudUserObjectRightId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Given ID is malformed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), serviceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), userId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("object_right_id"), objectRightId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
