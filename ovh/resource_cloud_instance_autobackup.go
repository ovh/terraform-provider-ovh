package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ resource.ResourceWithConfigure = (*cloudInstanceAutobackupResource)(nil)
var _ resource.ResourceWithImportState = (*cloudInstanceAutobackupResource)(nil)

func NewCloudInstanceAutobackupResource() resource.Resource {
	return &cloudInstanceAutobackupResource{}
}

type cloudInstanceAutobackupResource struct {
	config *Config
}

func (r *cloudInstanceAutobackupResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_instance_autobackup"
}

func (r *cloudInstanceAutobackupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudInstanceAutobackupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates an autobackup (Mistral crontrigger) for a Public Cloud instance.",
		Attributes: map[string]schema.Attribute{
			// Required — immutable
			"project_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Public Cloud project ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Crontrigger name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Backup image name prefix",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cron": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Cron schedule pattern (e.g. '0 2 * * *')",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rotation": schema.Int64Attribute{
				CustomType:  ovhtypes.TfInt64Type{},
				Required:    true,
				Description: "Number of backup versions to keep (minimum 1)",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "OpenStack region",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"instance_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "ID of the source instance",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			// Optional — immutable
			"distant": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Cross-region backup configuration",
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.RequiresReplace(),
				},
				Attributes: map[string]schema.Attribute{
					"region": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Required:    true,
						Description: "Target region for cross-region backup",
					},
					"image_name": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Required:    true,
						Description: "Image name for the cross-region backup copy",
					},
				},
			},

			// Computed
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Autobackup (crontrigger) ID",
			},
			"checksum": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Resource checksum for concurrency control",
			},
			"resource_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Resource status (READY, CREATING, DELETING, ERROR)",
			},
			"created_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Creation timestamp (RFC 3339)",
			},
			"updated_at": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Last update timestamp (RFC 3339)",
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current state from the infrastructure",
				Attributes:  cloudInstanceAutobackupCurrentStateSchemaAttributes(),
			},
		},
	}
}

func cloudInstanceAutobackupCurrentStateSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Crontrigger name",
		},
		"image_name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Backup image name prefix",
		},
		"cron": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Cron schedule pattern",
		},
		"rotation": schema.Int64Attribute{
			CustomType:  ovhtypes.TfInt64Type{},
			Computed:    true,
			Description: "Number of backup versions to keep",
		},
		"region": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "OpenStack region",
		},
		"instance_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Source instance ID",
		},
		"workflow_name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "OpenStack Mistral workflow name",
		},
		"next_execution_time": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Next scheduled execution time",
		},
		"distant": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Cross-region backup configuration",
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Target region for cross-region backup",
				},
				"image_name": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Image name for the cross-region backup copy",
				},
			},
		},
	}
}

func (r *cloudInstanceAutobackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <project_id>/<autobackup_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudInstanceAutobackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudInstanceAutobackupModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/compute/autobackup"

	var responseData CloudInstanceAutobackupAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Post %s", endpoint), err.Error())
		return
	}

	// Wait for resource to be ready
	endpoint = "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/compute/autobackup/" + url.PathEscape(responseData.Id)
	if err := helpers.WaitForAPIv2ResourceStatusReady(ctx, r.config.OVHClient, endpoint); err != nil {
		resp.Diagnostics.AddError("Error waiting for autobackup to be ready", err.Error())
		return
	}

	// Fetch up-to-date resource
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudInstanceAutobackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudInstanceAutobackupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/compute/autobackup/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudInstanceAutobackupAPIResponse
	if err := r.config.OVHClient.Get(endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Get %s", endpoint), err.Error())
		return
	}

	data.MergeWith(ctx, &responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update is not supported — autobackup is immutable (replace-only via RequiresReplace plan modifiers).
func (r *cloudInstanceAutobackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Autobackup resources are immutable. Changes require replacement.")
}

func (r *cloudInstanceAutobackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudInstanceAutobackupModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ProjectId.ValueString()) + "/compute/autobackup/" + url.PathEscape(data.Id.ValueString())
	if err := r.config.OVHClient.Delete(endpoint, nil); err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calling Delete %s", endpoint), err.Error())
	}
}
