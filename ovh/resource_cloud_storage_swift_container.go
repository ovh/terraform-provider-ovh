package ovh

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudStorageSwiftContainerResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudStorageSwiftContainerResource)(nil)
	_ resource.ResourceWithImportState = (*cloudStorageSwiftContainerResource)(nil)
)

func NewCloudStorageSwiftContainerResource() resource.Resource {
	return &cloudStorageSwiftContainerResource{}
}

type cloudStorageSwiftContainerResource struct {
	config *Config
}

func (r *cloudStorageSwiftContainerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_swift_container"
}

func (r *cloudStorageSwiftContainerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

var storageSwiftContainerMutableAttrs = MutableAttrs{
	Strings: []string{"type"},
	Sets:    []string{"cors"},
}

func storageSwiftContainerEndpoint(serviceName string) string {
	return "/v2/publicCloud/project/" + url.PathEscape(serviceName) + "/storage/swift/container"
}

func (r *cloudStorageSwiftContainerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Creates a Swift object storage container in a public cloud project using the publicCloud API.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Service name of the resource representing the id of the cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used.",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
				PlanModifiers: []planmodifier.String{
					EnvDefaultString("OVH_CLOUD_PROJECT_SERVICE", true),
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Container name, unique within the region: 1 to 256 characters, must not contain a slash (/)",
				MarkdownDescription: "Container name, unique within the region: 1 to 256 characters, must not contain a slash (`/`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 256),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[^/]+$`),
						"must not contain a slash",
					),
				},
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Upper-case region identifier hosting the container (e.g. GRA, SBG, BHS)",
				MarkdownDescription: "Upper-case region identifier hosting the container (e.g. `GRA`, `SBG`, `BHS`)",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					// The API upper-cases the region; a lower-case value would never match state.
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Z0-9-]+$`),
						"must be upper-case (e.g. \"GRA\")",
					),
				},
			},
			"archive": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				Description:         "Create the container with the Public Cloud Archive storage policy instead of the standard policy. Objects in an archive container are sealed after upload and must be unsealed before download. Defaults to false. Changing this value recreates the container.",
				MarkdownDescription: "Create the container with the Public Cloud Archive storage policy instead of the standard policy. Objects in an archive container are sealed after upload and must be unsealed before download. Defaults to `false`. Changing this value recreates the container.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplaceIfConfigured(),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Access type of the container (PRIVATE, PUBLIC, STATIC). Defaults to PRIVATE.",
				MarkdownDescription: "Access type of the container (`PRIVATE`, `PUBLIC`, `STATIC`). Defaults to `PRIVATE`.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("PRIVATE", "PUBLIC", "STATIC"),
				},
			},
			"cors": schema.SetAttribute{
				ElementType:         ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Origins allowed to access the container from a browser (Access-Control-Allow-Origin container metadata). Defaults to no origin.",
				MarkdownDescription: "Origins allowed to access the container from a browser (`Access-Control-Allow-Origin` container metadata). Defaults to no origin.",
				PlanModifiers: []planmodifier.Set{
					UseStateForUnknownSet(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Container identifier, in the <REGION>_<name> form",
				MarkdownDescription: "Container identifier, in the `<REGION>_<name>` form",
			},
			"checksum": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Computed hash representing the current target specification value",
				MarkdownDescription: "Computed hash representing the current target specification value",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(storageSwiftContainerMutableAttrs),
				},
			},
			"created_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Creation date of the container",
				MarkdownDescription: "Creation date of the container",
			},
			"updated_at": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Last update date of the container",
				MarkdownDescription: "Last update date of the container",
				PlanModifiers: []planmodifier.String{
					UnknownDuringUpdateStringModifier(storageSwiftContainerMutableAttrs),
				},
			},
			"resource_status": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Computed:            true,
				Description:         "Container readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, SUSPENDED, UNKNOWN, UPDATING)",
				MarkdownDescription: "Container readiness in the system (`CREATING`, `DELETING`, `ERROR`, `OUT_OF_SYNC`, `READY`, `SUSPENDED`, `UNKNOWN`, `UPDATING`)",
				PlanModifiers: []planmodifier.String{
					OutOfSyncPlanModifier(),
				},
			},
			"current_state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current observed state of the container, read from the Swift backend",
				PlanModifiers: []planmodifier.Object{
					UnknownDuringUpdateObjectModifier(storageSwiftContainerMutableAttrs),
				},
				Attributes: storageSwiftContainerCurrentStateResourceAttributes(),
			},
		},
	}
}

// storageSwiftContainerCurrentStateResourceAttributes returns the current_state attributes of
// the resource schema. The data sources declare the same shape with their own
// schema package.
func storageSwiftContainerCurrentStateResourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Container name",
		},
		"location": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Region hosting the container",
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Region identifier",
				},
			},
		},
		"archive": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether the container uses the Public Cloud Archive storage policy",
		},
		"type": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Access type of the container (PRIVATE, PUBLIC, STATIC)",
		},
		"cors": schema.SetAttribute{
			ElementType: ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Origins currently allowed to access the container from a browser",
		},
		"objects_count": schema.Int64Attribute{
			Computed:    true,
			Description: "Number of objects stored in the container",
		},
		"objects_size": schema.Int64Attribute{
			Computed:    true,
			Description: "Total size of the objects stored in the container, in bytes",
		},
		"static_url": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Public URL of the container, only set when type is PUBLIC or STATIC",
		},
	}
}

func (r *cloudStorageSwiftContainerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<container_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

func (r *cloudStorageSwiftContainerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudStorageSwiftContainerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createPayload := data.ToCreate()

	endpoint := storageSwiftContainerEndpoint(data.ServiceName.ValueString())

	var responseData CloudStorageSwiftContainerAPIResponse
	if err := r.config.OVHClient.Post(endpoint, createPayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	// Save state immediately so the container id is tracked even if the workflow fails
	data.MergeWith(ctx, &responseData)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if _, err := r.waitForStorageSwiftContainerReady(ctx, data.ServiceName.ValueString(), responseData.Id); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for container to be ready",
			err.Error(),
		)
		return
	}

	endpoint = storageSwiftContainerEndpoint(data.ServiceName.ValueString()) + "/" + url.PathEscape(responseData.Id)
	// json.Unmarshal merges into a non-zero struct: stale slices/pointers would leak.
	responseData = CloudStorageSwiftContainerAPIResponse{}
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

func (r *cloudStorageSwiftContainerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudStorageSwiftContainerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := storageSwiftContainerEndpoint(data.ServiceName.ValueString()) + "/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudStorageSwiftContainerAPIResponse
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

func (r *cloudStorageSwiftContainerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, planData CloudStorageSwiftContainerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	updatePayload := planData.ToUpdate(data.Checksum.ValueString())

	endpoint := storageSwiftContainerEndpoint(data.ServiceName.ValueString()) + "/" + url.PathEscape(data.Id.ValueString())

	var responseData CloudStorageSwiftContainerAPIResponse
	if err := r.config.OVHClient.Put(endpoint, updatePayload, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Put %s", endpoint),
			err.Error(),
		)
		return
	}

	if _, err := r.waitForStorageSwiftContainerReady(ctx, data.ServiceName.ValueString(), data.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for container to be ready after update",
			err.Error(),
		)
		return
	}

	responseData = CloudStorageSwiftContainerAPIResponse{}
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

func (r *cloudStorageSwiftContainerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudStorageSwiftContainerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := storageSwiftContainerEndpoint(data.ServiceName.ValueString()) + "/" + url.PathEscape(data.Id.ValueString())

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
			res := &CloudStorageSwiftContainerAPIResponse{}
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
					return res, "DELETED", nil
				}
				return res, "", err
			}
			if res.ResourceStatus == "ERROR" {
				return res, res.ResourceStatus, cloudResourceErrorFromTasks("container", data.Id.ValueString(), res.CurrentTasks)
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for container to be deleted",
			err.Error(),
		)
	}
}

func (r *cloudStorageSwiftContainerResource) waitForStorageSwiftContainerReady(ctx context.Context, serviceName, containerId string) (any, error) {
	stateConf := &retry.StateChangeConf{
		// UNKNOWN: transient Swift 5xx on the container HEAD.
		Pending: []string{"CREATING", "UPDATING", "PENDING", "OUT_OF_SYNC", "UNKNOWN"},
		Target:  []string{"READY"},
		Refresh: func() (any, string, error) {
			res := &CloudStorageSwiftContainerAPIResponse{}
			endpoint := storageSwiftContainerEndpoint(serviceName) + "/" + url.PathEscape(containerId)
			err := r.config.OVHClient.GetWithContext(ctx, endpoint, res)
			if err != nil {
				return res, "", err
			}
			// ERROR is terminal: surface it with the task reason instead of a generic unexpected-state.
			if res.ResourceStatus == "ERROR" {
				return res, res.ResourceStatus, cloudResourceErrorFromTasks("container", containerId, res.CurrentTasks)
			}
			if res.ResourceStatus == "SUSPENDED" {
				return res, res.ResourceStatus, fmt.Errorf("container %s is SUSPENDED: its region is in maintenance, retry once the maintenance is over", containerId)
			}
			return res, res.ResourceStatus, nil
		},
		Timeout:    20 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}
