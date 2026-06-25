package ovh

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

// computeConsumerID builds the API consumer ID (URL-safe base64 of service:resourceType:resourceId).
func computeConsumerID(service, resourceType, resourceId string) string {
	raw := service + ":" + resourceType + ":" + resourceId
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

var (
	_ resource.Resource                = (*cloudKeyManagerSecretConsumerResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudKeyManagerSecretConsumerResource)(nil)
	_ resource.ResourceWithImportState = (*cloudKeyManagerSecretConsumerResource)(nil)
)

func NewCloudKeyManagerSecretConsumerResource() resource.Resource {
	return &cloudKeyManagerSecretConsumerResource{}
}

type cloudKeyManagerSecretConsumerResource struct {
	config *Config
}

func (r *cloudKeyManagerSecretConsumerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_key_manager_secret_consumer"
}

func (r *cloudKeyManagerSecretConsumerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudKeyManagerSecretConsumerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Registers a consumer on a Barbican Key Manager secret for a public cloud project.",
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
			"secret_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "UUID of the secret to register the consumer on",
				MarkdownDescription: "UUID of the secret to register the consumer on",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"service": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "OpenStack service type of the consumer (COMPUTE, IMAGE, LOADBALANCER, NETWORK)",
				MarkdownDescription: "OpenStack service type of the consumer (`COMPUTE`, `IMAGE`, `LOADBALANCER`, `NETWORK`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_type": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Type of the resource consuming the secret (IMAGE, INSTANCE, LOADBALANCER)",
				MarkdownDescription: "Type of the resource consuming the secret (`IMAGE`, `INSTANCE`, `LOADBALANCER`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"resource_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "UUID of the resource consuming the secret",
				MarkdownDescription: "UUID of the resource consuming the secret",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Consumer ID as returned by the API (URL-safe base64 of service:resource_type:resource_id)",
				MarkdownDescription: "Consumer ID as returned by the API (URL-safe base64 of `service:resource_type:resource_id`)",
			},
		},
	}
}

func (r *cloudKeyManagerSecretConsumerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import format: service_name/secret_id/service/resource_type/resource_id
	splits := strings.SplitN(req.ID, "/", 5)
	if len(splits) != 5 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like: <service_name>/<secret_id>/<service>/<resource_type>/<resource_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("secret_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service"), splits[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_type"), splits[3])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_id"), splits[4])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), computeConsumerID(splits[2], splits[3], splits[4]))...)
}

func (r *cloudKeyManagerSecretConsumerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudKeyManagerSecretConsumerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/keyManager/secret/" + url.PathEscape(data.SecretId.ValueString()) + "/consumer"

	payload := CloudKeyManagerSecretConsumerPayload{
		Service:      data.Service.ValueString(),
		ResourceType: data.ResourceType.ValueString(),
		ResourceId:   data.ResourceId.ValueString(),
	}

	var responseData CloudKeyManagerSecretConsumerAPIResponse
	if err := r.config.OVHClient.Post(endpoint, payload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	data.Service = ovhtypes.TfStringValue{StringValue: types.StringValue(responseData.Service)}
	data.ResourceType = ovhtypes.TfStringValue{StringValue: types.StringValue(responseData.ResourceType)}
	data.ResourceId = ovhtypes.TfStringValue{StringValue: types.StringValue(responseData.ResourceId)}

	// Use the raw API consumer ID so it can be fed directly to the singular
	// consumer data source. Fall back to deriving it when the API omits it.
	consumerId := responseData.Id
	if consumerId == "" {
		consumerId = computeConsumerID(responseData.Service, responseData.ResourceType, responseData.ResourceId)
	}
	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(consumerId)}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudKeyManagerSecretConsumerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudKeyManagerSecretConsumerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	consumerId := computeConsumerID(data.Service.ValueString(), data.ResourceType.ValueString(), data.ResourceId.ValueString())
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/keyManager/secret/" + url.PathEscape(data.SecretId.ValueString()) +
		"/consumer/" + url.PathEscape(consumerId)

	var consumer CloudKeyManagerSecretConsumerAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &consumer); err != nil {
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

	// Reconcile drift by copying the fetched consumer into the state model
	data.Service = ovhtypes.TfStringValue{StringValue: types.StringValue(consumer.Service)}
	data.ResourceType = ovhtypes.TfStringValue{StringValue: types.StringValue(consumer.ResourceType)}
	data.ResourceId = ovhtypes.TfStringValue{StringValue: types.StringValue(consumer.ResourceId)}
	if consumer.Id != "" {
		data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(consumer.Id)}
	} else {
		data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(consumerId)}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudKeyManagerSecretConsumerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Secret consumers cannot be updated. All fields require replacement.")
}

func (r *cloudKeyManagerSecretConsumerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudKeyManagerSecretConsumerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	consumerId := computeConsumerID(data.Service.ValueString(), data.ResourceType.ValueString(), data.ResourceId.ValueString())
	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/keyManager/secret/" + url.PathEscape(data.SecretId.ValueString()) +
		"/consumer/" + url.PathEscape(consumerId)

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
}
