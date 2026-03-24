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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudKeymanagerSecretResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudKeymanagerSecretResource)(nil)
	_ resource.ResourceWithImportState = (*cloudKeymanagerSecretResource)(nil)
)

func NewCloudKeymanagerSecretResource() resource.Resource {
	return &cloudKeymanagerSecretResource{}
}

type cloudKeymanagerSecretResource struct {
	config *Config
}

func (r *cloudKeymanagerSecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_keymanager_secret"
}

func (r *cloudKeymanagerSecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudKeymanagerSecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a secret in the Barbican Key Manager service for a public cloud project.",
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
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Region where the secret will be created",
				MarkdownDescription: "Region where the secret will be created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Name of the secret",
				MarkdownDescription: "Name of the secret",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"secret_type": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Type of the secret (SYMMETRIC, PUBLIC, PRIVATE, PASSPHRASE, CERTIFICATE, OPAQUE)",
				MarkdownDescription: "Type of the secret (`SYMMETRIC`, `PUBLIC`, `PRIVATE`, `PASSPHRASE`, `CERTIFICATE`, `OPAQUE`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"algorithm": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Algorithm associated with the secret (e.g., AES, RSA)",
				MarkdownDescription: "Algorithm associated with the secret (e.g., `AES`, `RSA`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bit_length": schema.Int64Attribute{
				Optional:            true,
				Description:         "Bit length of the secret (e.g., 256)",
				MarkdownDescription: "Bit length of the secret (e.g., `256`)",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"mode": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Mode of the secret algorithm (e.g., CBC)",
				MarkdownDescription: "Mode of the secret algorithm (e.g., `CBC`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"payload": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Sensitive:           true,
				Description:         "Secret payload data (base64-encoded). Write-only, never returned in responses. Requires payload_content_type.",
				MarkdownDescription: "Secret payload data (base64-encoded). Write-only, never returned in responses. Requires `payload_content_type`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"payload_content_type": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Content type of the secret payload (TEXT_PLAIN, APPLICATION_OCTET_STREAM, APPLICATION_PKIX_CERT, APPLICATION_PKCS8)",
				MarkdownDescription: "Content type of the secret payload (`TEXT_PLAIN`, `APPLICATION_OCTET_STREAM`, `APPLICATION_PKIX_CERT`, `APPLICATION_PKCS8`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"expiration": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "Expiration date of the secret in RFC3339 format",
				MarkdownDescription: "Expiration date of the secret in RFC3339 format",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metadata": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				Description:         "Key-value metadata for the secret. This is the only mutable field on a secret.",
				MarkdownDescription: "Key-value metadata for the secret. This is the only mutable field on a secret.",
			},

			// Computed
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Secret ID",
				MarkdownDescription: "Secret ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current resource state",
				MarkdownDescription: "Computed hash representing the current resource state",
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the secret",
				MarkdownDescription: "Creation date of the secret",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the secret",
				MarkdownDescription: "Last update date of the secret",
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Secret readiness status (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, UPDATING)",
				MarkdownDescription: "Secret readiness status (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `UPDATING`)",
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state of the secret as reported by OpenStack Barbican",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"secret_type": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"algorithm": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"bit_length": schema.Int64Attribute{
						Computed: true,
					},
					"mode": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"payload_content_type": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"expiration": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"secret_ref": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"status": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"region": schema.StringAttribute{
						CustomType: ovhtypes.TfStringType{},
						Computed:   true,
					},
					"metadata": schema.MapAttribute{
						ElementType: types.StringType,
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *cloudKeymanagerSecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<secret_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudKeymanagerSecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudKeymanagerSecretModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate(ctx)

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/keyManager/secret"

	var responseData CloudKeymanagerSecretAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save state immediately so the resource ID is tracked
	data.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Wait for secret to be READY
	_, err := r.waitForReady(ctx, data.ServiceName.ValueString(), responseData.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for KMS secret to be ready",
			err.Error(),
		)
		return
	}

	// Read final state
	getEndpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/keyManager/secret/" + url.PathEscape(responseData.Id)
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

func (r *cloudKeymanagerSecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudKeymanagerSecretModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/keyManager/secret/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudKeymanagerSecretAPIResponse
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

func (r *cloudKeymanagerSecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudKeymanagerSecretModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state CloudKeymanagerSecretModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use checksum from current state for concurrency control
	plan.Checksum = state.Checksum
	plan.Id = state.Id

	updatePayload := plan.ToUpdate(ctx)
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(plan.ServiceName.ValueString()) + "/keyManager/secret/" + url.PathEscape(plan.Id.ValueString())

	var responseData CloudKeymanagerSecretAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	// Wait for secret to be READY after update
	_, err := r.waitForReady(ctx, plan.ServiceName.ValueString(), plan.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for KMS secret to be ready after update",
			err.Error(),
		)
		return
	}

	// Read final state
	getEndpoint := "/v2/publicCloud/project/" + url.PathEscape(plan.ServiceName.ValueString()) + "/keyManager/secret/" + url.PathEscape(plan.Id.ValueString())
	if err := r.config.OVHClient.Get(getEndpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", getEndpoint),
			err.Error(),
		)
		return
	}

	plan.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *cloudKeymanagerSecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudKeymanagerSecretModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/keyManager/secret/" + url.PathEscape(data.Id.ValueString())

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
			res := &CloudKeymanagerSecretAPIResponse{}
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
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for KMS secret to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudKeymanagerSecretResource) waitForReady(ctx context.Context, serviceName, secretId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"CREATING", "UPDATING", "OUT_OF_SYNC"},
		Target:  []string{"READY"},
		Refresh: func() (interface{}, string, error) {
			res := &CloudKeymanagerSecretAPIResponse{}
			endpoint := "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/keyManager/secret/" + url.PathEscape(secretId)
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}
