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

// The account logs (audit, API activity and IAM access policy logs) are forwarded to a Logs Data
// Platform stream through three endpoints that share the same payloads, so a single implementation
// serves the three resources.

var (
	_ resource.Resource                = (*accountLogSubscriptionResource)(nil)
	_ resource.ResourceWithConfigure   = (*accountLogSubscriptionResource)(nil)
	_ resource.ResourceWithImportState = (*accountLogSubscriptionResource)(nil)
)

func NewMeLogsAuditLogSubscriptionResource() resource.Resource {
	return &accountLogSubscriptionResource{
		typeName: "me_logs_audit_log_subscription",
		endpoint: "/me/logs/audit/log/subscription",
		logs:     "the audit logs",
	}
}

func NewMeApiLogSubscriptionResource() resource.Resource {
	return &accountLogSubscriptionResource{
		typeName: "me_api_log_subscription",
		endpoint: "/me/api/log/subscription",
		logs:     "the activity logs",
	}
}

func NewIamLogSubscriptionResource() resource.Resource {
	return &accountLogSubscriptionResource{
		typeName: "iam_log_subscription",
		endpoint: "/v2/iam/log/subscription",
		logs:     "the access policy logs",
	}
}

type accountLogSubscriptionResource struct {
	config *Config

	// typeName is the resource type name without the provider prefix.
	typeName string
	// endpoint is the API path of the subscriptions, without a trailing slash.
	endpoint string
	// logs names the forwarded logs in the description of the resource.
	logs string
}

// AccountLogSubscriptionCreateOpts mirrors dbaas.logs.LogSubscriptionCreation.
type AccountLogSubscriptionCreateOpts struct {
	Kind     string `json:"kind"`
	StreamId string `json:"streamId"`
}

// AccountLogSubscriptionOperation mirrors dbaas.logs.LogSubscriptionResponse, returned when a
// subscription is created or deleted.
type AccountLogSubscriptionOperation struct {
	OperationId string `json:"operationId"`
	ServiceName string `json:"serviceName"`
}

type AccountLogSubscriptionResource struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// AccountLogSubscription mirrors dbaas.logs.LogSubscription.
type AccountLogSubscription struct {
	CreatedAt      string                         `json:"createdAt"`
	Kind           string                         `json:"kind"`
	Resource       AccountLogSubscriptionResource `json:"resource"`
	ServiceName    string                         `json:"serviceName"`
	StreamId       string                         `json:"streamId"`
	SubscriptionId string                         `json:"subscriptionId"`
	UpdatedAt      string                         `json:"updatedAt"`
}

type AccountLogSubscriptionModel struct {
	ID             ovhtypes.TfStringValue `tfsdk:"id"`
	Kind           ovhtypes.TfStringValue `tfsdk:"kind"`
	StreamId       ovhtypes.TfStringValue `tfsdk:"stream_id"`
	LdpServiceName ovhtypes.TfStringValue `tfsdk:"ldp_service_name"`
	ResourceName   ovhtypes.TfStringValue `tfsdk:"resource_name"`
	ResourceType   ovhtypes.TfStringValue `tfsdk:"resource_type"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
}

func (m *AccountLogSubscriptionModel) MergeWith(sub *AccountLogSubscription) {
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
func (m AccountLogSubscriptionModel) CreatedState(subscriptionId, ldpServiceName string) AccountLogSubscriptionModel {
	m.ID = ovhtypes.NewTfStringValue(subscriptionId)
	m.LdpServiceName = ovhtypes.NewTfStringValue(ldpServiceName)
	m.ResourceName = ovhtypes.NewTfStringNull()
	m.ResourceType = ovhtypes.NewTfStringNull()
	m.CreatedAt = ovhtypes.NewTfStringNull()
	m.UpdatedAt = ovhtypes.NewTfStringNull()
	return m
}

func (r *accountLogSubscriptionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName
}

func (r *accountLogSubscriptionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *accountLogSubscriptionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	description := fmt.Sprintf("Forwards %s of the OVHcloud account to a Logs Data Platform stream.", r.logs)

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
			"kind": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         fmt.Sprintf("Log kind name to subscribe to. The kinds available are listed by GET %s", r.kindEndpoint()),
				MarkdownDescription: fmt.Sprintf("Log kind name to subscribe to. The kinds available are listed by `GET %s`", r.kindEndpoint()),
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

// kindEndpoint is the API path listing the log kinds that can be subscribed to.
func (r *accountLogSubscriptionResource) kindEndpoint() string {
	return strings.TrimSuffix(r.endpoint, "/subscription") + "/kind"
}

func (r *accountLogSubscriptionResource) subscriptionEndpoint(id string) string {
	return r.endpoint + "/" + url.PathEscape(id)
}

func (r *accountLogSubscriptionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AccountLogSubscriptionModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := &AccountLogSubscriptionCreateOpts{
		Kind:     data.Kind.ValueString(),
		StreamId: data.StreamId.ValueString(),
	}
	var op AccountLogSubscriptionOperation
	if err := r.config.OVHClient.PostWithContext(ctx, r.endpoint, opts, &op); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", r.endpoint), err.Error())
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
	endpoint := r.subscriptionEndpoint(*res.SubscriptionID)
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &sub); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(&sub)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *accountLogSubscriptionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AccountLogSubscriptionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sub AccountLogSubscription
	endpoint := r.subscriptionEndpoint(data.ID.ValueString())
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

func (r *accountLogSubscriptionResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Every configurable attribute requires a replacement, so there is nothing to update.
	resp.Diagnostics.AddError("Update not supported", "A log subscription cannot be updated, it must be replaced.")
}

func (r *accountLogSubscriptionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data AccountLogSubscriptionModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var op AccountLogSubscriptionOperation
	endpoint := r.subscriptionEndpoint(data.ID.ValueString())
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

func (r *accountLogSubscriptionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
