package ovh

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// The logs of a KMS are forwarded to a Logs Data Platform stream with the payloads of the account
// log subscriptions (dbaas.logs.LogSubscriptionCreation, LogSubscriptionResponse and
// LogSubscription), under the path of the KMS.

var (
	_ resource.Resource                = (*okmsLogSubscriptionResource)(nil)
	_ resource.ResourceWithConfigure   = (*okmsLogSubscriptionResource)(nil)
	_ resource.ResourceWithImportState = (*okmsLogSubscriptionResource)(nil)
)

func NewOkmsLogSubscriptionResource() resource.Resource {
	return &okmsLogSubscriptionResource{}
}

type okmsLogSubscriptionResource struct {
	config *Config
}

type OkmsLogSubscriptionModel struct {
	ID             ovhtypes.TfStringValue `tfsdk:"id"`
	OkmsId         ovhtypes.TfStringValue `tfsdk:"okms_id"`
	Kind           ovhtypes.TfStringValue `tfsdk:"kind"`
	StreamId       ovhtypes.TfStringValue `tfsdk:"stream_id"`
	LdpServiceName ovhtypes.TfStringValue `tfsdk:"ldp_service_name"`
	ResourceName   ovhtypes.TfStringValue `tfsdk:"resource_name"`
	ResourceType   ovhtypes.TfStringValue `tfsdk:"resource_type"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
}

func (m *OkmsLogSubscriptionModel) MergeWith(sub *AccountLogSubscription) {
	m.ID = ovhtypes.NewTfStringValue(sub.SubscriptionId)
	m.Kind = ovhtypes.NewTfStringValue(sub.Kind)
	m.StreamId = ovhtypes.NewTfStringValue(sub.StreamId)
	m.LdpServiceName = ovhtypes.NewTfStringValue(sub.ServiceName)
	m.ResourceName = ovhtypes.NewTfStringValue(sub.Resource.Name)
	m.ResourceType = ovhtypes.NewTfStringValue(sub.Resource.Type)
	m.CreatedAt = ovhtypes.NewTfStringValue(sub.CreatedAt)
	m.UpdatedAt = ovhtypes.NewTfStringValue(sub.UpdatedAt)
}

// CreatedState returns the state of a subscription whose creation succeeded, before its attributes
// can be read: they are null rather than unknown, which a state cannot hold.
func (m OkmsLogSubscriptionModel) CreatedState(subscriptionId, ldpServiceName string) OkmsLogSubscriptionModel {
	m.ID = ovhtypes.NewTfStringValue(subscriptionId)
	m.LdpServiceName = ovhtypes.NewTfStringValue(ldpServiceName)
	m.ResourceName = ovhtypes.NewTfStringNull()
	m.ResourceType = ovhtypes.NewTfStringNull()
	m.CreatedAt = ovhtypes.NewTfStringNull()
	m.UpdatedAt = ovhtypes.NewTfStringNull()
	return m
}

func (r *okmsLogSubscriptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_okms_log_subscription"
}

func (r *okmsLogSubscriptionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *okmsLogSubscriptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	description := "Forwards the logs of an OVHcloud KMS to a Logs Data Platform stream."

	resp.Schema = schema.Schema{
		Description:         description,
		MarkdownDescription: description,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Identifier of the subscription",
				MarkdownDescription: "Identifier of the subscription",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"okms_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "ID of the KMS",
				MarkdownDescription: "ID of the KMS",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"kind": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Log kind name to subscribe to. The kinds available are listed by GET /v2/okms/resource/{okmsId}/log/kind",
				MarkdownDescription: "Log kind name to subscribe to. The kinds available are listed by `GET /v2/okms/resource/{okmsId}/log/kind`",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"stream_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Identifier of the Logs Data Platform stream the logs are forwarded to",
				MarkdownDescription: "Identifier of the Logs Data Platform stream the logs are forwarded to",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Computed
			"ldp_service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Name of the Logs Data Platform service holding the stream",
				MarkdownDescription: "Name of the Logs Data Platform service holding the stream",
			},
			"resource_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Name of the subscribed resource, where the logs come from",
				MarkdownDescription: "Name of the subscribed resource, where the logs come from",
			},
			"resource_type": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Type of the subscribed resource, where the logs come from",
				MarkdownDescription: "Type of the subscribed resource, where the logs come from",
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the subscription",
				MarkdownDescription: "Creation date of the subscription",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the subscription",
				MarkdownDescription: "Last update date of the subscription",
			},
		},
	}
}

// okmsLogSubscriptionEndpoint is the API path of the subscriptions of a KMS, without a trailing slash.
func okmsLogSubscriptionEndpoint(okmsId string) string {
	return "/v2/okms/resource/" + url.PathEscape(okmsId) + "/log/subscription"
}

func (r *okmsLogSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data OkmsLogSubscriptionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := &AccountLogSubscriptionCreateOpts{
		Kind:     data.Kind.ValueString(),
		StreamId: data.StreamId.ValueString(),
	}
	endpoint := okmsLogSubscriptionEndpoint(data.OkmsId.ValueString())
	var op AccountLogSubscriptionOperation
	if err := r.config.OVHClient.PostWithContext(ctx, endpoint, opts, &op); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// The subscription is created asynchronously, and only the operation tells its identifier.
	res, err := waitForDbaasLogsOperation(ctx, r.config.OVHClient, op.ServiceName, op.OperationId)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the log subscription to be created", err.Error())
		return
	}
	if res.SubscriptionID == nil || *res.SubscriptionID == "" {
		resp.Diagnostics.AddError(
			"Error reading the created log subscription",
			fmt.Sprintf("Operation %s of service %s succeeded without a subscription id", op.OperationId, op.ServiceName),
		)
		return
	}

	// Save the subscription as soon as it exists, so that Terraform tracks it, and replaces it on
	// the next apply, when reading it fails.
	created := data.CreatedState(*res.SubscriptionID, op.ServiceName)
	resp.Diagnostics.Append(resp.State.Set(ctx, &created)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sub AccountLogSubscription
	endpoint = endpoint + "/" + url.PathEscape(*res.SubscriptionID)
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &sub); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(&sub)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *okmsLogSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data OkmsLogSubscriptionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sub AccountLogSubscription
	endpoint := okmsLogSubscriptionEndpoint(data.OkmsId.ValueString()) + "/" + url.PathEscape(data.ID.ValueString())
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &sub); err != nil {
		var errOvh *ovh.APIError
		if errors.As(err, &errOvh) && errOvh.Code == http.StatusNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(&sub)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *okmsLogSubscriptionResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Every configurable attribute requires a replacement, so there is nothing to update.
	resp.Diagnostics.AddError("Update not supported", "A log subscription cannot be updated, it must be replaced.")
}

func (r *okmsLogSubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data OkmsLogSubscriptionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var op AccountLogSubscriptionOperation
	endpoint := okmsLogSubscriptionEndpoint(data.OkmsId.ValueString()) + "/" + url.PathEscape(data.ID.ValueString())
	if err := r.config.OVHClient.DeleteWithContext(ctx, endpoint, &op); err != nil {
		var errOvh *ovh.APIError
		if errors.As(err, &errOvh) && errOvh.Code == http.StatusNotFound {
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
		return
	}

	if _, err := waitForDbaasLogsOperation(ctx, r.config.OVHClient, op.ServiceName, op.OperationId); err != nil {
		resp.Diagnostics.AddError("Error waiting for the log subscription to be deleted", err.Error())
	}
}

func (r *okmsLogSubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids := strings.Split(req.ID, "/")
	if len(ids) != 2 || ids[0] == "" || ids[1] == "" {
		resp.Diagnostics.AddError(
			"Given ID is malformed",
			"ID must be formatted like the following: <okms_id>/<subscription_id>",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("okms_id"), ids[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), ids[1])...)
}
