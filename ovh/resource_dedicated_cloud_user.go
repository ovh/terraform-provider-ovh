package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.ResourceWithConfigure   = (*dedicatedCloudUserResource)(nil)
	_ resource.ResourceWithImportState = (*dedicatedCloudUserResource)(nil)
)

func NewDedicatedCloudUserResource() resource.Resource {
	return &dedicatedCloudUserResource{}
}

type dedicatedCloudUserResource struct {
	config *Config
}

func (r *dedicatedCloudUserResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dedicated_cloud_user"
}

func (r *dedicatedCloudUserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *dedicatedCloudUserResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = DedicatedCloudUserResourceSchema(ctx)
}

func dedicatedCloudUserEndpoint(serviceName string, userId int64) string {
	return fmt.Sprintf("/dedicatedCloud/%s/user/%d", url.PathEscape(serviceName), userId)
}

func (r *dedicatedCloudUserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data, responseData DedicatedCloudUserModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()

	var task DedicatedCloudTask
	endpoint := "/dedicatedCloud/" + url.PathEscape(serviceName) + "/user"
	if err := r.config.OVHClient.PostWithContext(ctx, endpoint, data.ToCreate(), &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	finishedTask, err := waitForDedicatedCloudTask(ctx, r.config.OVHClient, serviceName, task.TaskId)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error waiting for user creation on %s", serviceName),
			err.Error(),
		)
		return
	}

	endpoint = dedicatedCloudUserEndpoint(serviceName, finishedTask.UserId)
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&data)
	responseData.ID = ovhtypes.NewTfStringValue(fmt.Sprintf("%s/%d", serviceName, responseData.UserId.ValueInt64()))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *dedicatedCloudUserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data, responseData DedicatedCloudUserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := dedicatedCloudUserEndpoint(data.ServiceName.ValueString(), data.UserId.ValueInt64())
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
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

func (r *dedicatedCloudUserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var planData, stateData, responseData DedicatedCloudUserModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := stateData.ServiceName.ValueString()
	userId := stateData.UserId.ValueInt64()

	var propertiesTask DedicatedCloudTask
	endpoint := dedicatedCloudUserEndpoint(serviceName, userId) + "/changeProperties"
	if err := r.config.OVHClient.PostWithContext(ctx, endpoint, planData.ToChangeProperties(), &propertiesTask); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	if _, err := waitForDedicatedCloudTask(ctx, r.config.OVHClient, serviceName, propertiesTask.TaskId); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error waiting for user update on %s", serviceName),
			err.Error(),
		)
		return
	}

	if !planData.Password.IsUnknown() && !planData.Password.Equal(stateData.Password) {
		var passwordTask DedicatedCloudTask
		passwordPayload := struct {
			Password *string `json:"password,omitempty"`
		}{}

		if !planData.Password.IsNull() {
			val := planData.Password.ValueString()
			passwordPayload.Password = &val
		}

		passwordEndpoint := dedicatedCloudUserEndpoint(serviceName, userId) + "/changePassword"
		if err := r.config.OVHClient.PostWithContext(ctx, passwordEndpoint, passwordPayload, &passwordTask); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error calling Post %s", passwordEndpoint),
				err.Error(),
			)
			return
		}

		if _, err := waitForDedicatedCloudTask(ctx, r.config.OVHClient, serviceName, passwordTask.TaskId); err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error waiting for password update on %s", serviceName),
				err.Error(),
			)
			return
		}
	}

	endpoint = dedicatedCloudUserEndpoint(serviceName, userId)
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	responseData.MergeWith(&planData)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &responseData)...)
}

func (r *dedicatedCloudUserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DedicatedCloudUserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceName := data.ServiceName.ValueString()

	var task DedicatedCloudTask
	endpoint := dedicatedCloudUserEndpoint(serviceName, data.UserId.ValueInt64())
	if err := r.config.OVHClient.DeleteWithContext(ctx, endpoint, &task); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Delete %s", endpoint),
			err.Error(),
		)
		return
	}

	if _, err := waitForDedicatedCloudTask(ctx, r.config.OVHClient, serviceName, task.TaskId); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error waiting for user deletion on %s", serviceName),
			err.Error(),
		)
	}
}

func (r *dedicatedCloudUserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	serviceName, userId, err := splitDedicatedCloudUserId(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Given ID is malformed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), serviceName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("user_id"), userId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}
