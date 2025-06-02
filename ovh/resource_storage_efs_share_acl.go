package ovh

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/ovhwrap"
)

var _ resource.ResourceWithConfigure = (*storageEfsShareAclResource)(nil)

func NewStorageEfsShareAclResource() resource.Resource {
	return &storageEfsShareAclResource{}
}

type storageEfsShareAclResource struct {
	config *Config
}

func (r *storageEfsShareAclResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_efs_share_acl"
}

func (d *storageEfsShareAclResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	d.config = config
}

func (d *storageEfsShareAclResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = StorageEfsShareAclResourceSchema(ctx)
}

func (r *storageEfsShareAclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData StorageEfsShareAclModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.ShareId.ValueString()) + "/acl"
	if err := r.config.OVHClient.PostWithContext(ctx, endpoint, data.ToCreate(), &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	resWait, err := r.WaitForAclCreation(ctx, r.config.OVHClient, data.ServiceName.ValueString(), data.ShareId.ValueString(), responseData.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Operation"),
			err.Error(),
		)
		return
	}

	resAcl := &StorageEfsShareAclModel{}
	endpoint = "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.ShareId.ValueString()) + "/acl/" + url.PathEscape(resWait.(string))
	if err := r.config.OVHClient.Get(endpoint, resAcl); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(resAcl)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *storageEfsShareAclResource) WaitForAclCreation(ctx context.Context, client *ovhwrap.Client, serviceName, shareID, aclID string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"queued_to_apply", "applying"},
		Target:  []string{"active"},
		Refresh: func() (interface{}, string, error) {
			res := &StorageEfsShareAclModel{}
			endpoint := "/storage/netapp/" + url.PathEscape(serviceName) + "/share/" + url.PathEscape(shareID) + "/acl/" + url.PathEscape(aclID)
			err := client.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			return res.Id.ValueString(), res.Status.ValueString(), nil
		},
		Timeout:    360 * time.Second,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	res, err := stateConf.WaitForStateContext(ctx)

	return res, err
}

func (r *storageEfsShareAclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData StorageEfsShareAclModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.ShareId.ValueString()) + "/acl/" + url.PathEscape(data.Id.ValueString()) + ""

	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *storageEfsShareAclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("not implemented", "this function should never be called")
}

func (r *storageEfsShareAclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StorageEfsShareAclModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	endpoint := "/storage/netapp/" + url.PathEscape(data.ServiceName.ValueString()) + "/share/" + url.PathEscape(data.ShareId.ValueString()) + "/acl/" + url.PathEscape(data.Id.ValueString()) + ""
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
	}
}
