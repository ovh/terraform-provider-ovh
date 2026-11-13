package ovh

import (
	"context"
	"fmt"
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

var (
	_ resource.Resource                = (*cloudStorageObjectBucketReplicationJobResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudStorageObjectBucketReplicationJobResource)(nil)
	_ resource.ResourceWithImportState = (*cloudStorageObjectBucketReplicationJobResource)(nil)
)

func NewCloudStorageObjectBucketReplicationJobResource() resource.Resource {
	return &cloudStorageObjectBucketReplicationJobResource{}
}

type cloudStorageObjectBucketReplicationJobResource struct {
	config *Config
}

func (r *cloudStorageObjectBucketReplicationJobResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_replication_job"
}

func (r *cloudStorageObjectBucketReplicationJobResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudStorageObjectBucketReplicationJobResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Starts a batch replication job on an S3 compatible object storage bucket of a public cloud project using the publicCloud API: " +
			"the current replication rules of the bucket are applied to the objects already in it. " +
			"A job cannot be cancelled: destroying the resource only removes it from the state.",
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
			"bucket_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Source bucket identifier, in the <REGION>_<name> form",
				MarkdownDescription: "Source bucket identifier, in the `<REGION>_<name>` form",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Unique identifier of the job",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Job status (ACTIVE, COMPLETED, CREATED, FAILED, PREPARING)"},
			"total":      schema.Int64Attribute{Computed: true, Description: "Total number of objects to replicate"},
			"completed":  schema.Int64Attribute{Computed: true, Description: "Number of successfully replicated objects"},
			"failed":     schema.Int64Attribute{Computed: true, Description: "Number of objects that failed to replicate"},
			"stuck":      schema.Int64Attribute{Computed: true, Description: "Number of stuck objects"},
			"created_at": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Job creation date"},
		},
	}
}

func (r *cloudStorageObjectBucketReplicationJobResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 3 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<bucket_id>/<job_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[2])...)
}

func storageObjectBucketReplicationJobsEndpoint(serviceName, bucketId string) string {
	return storageObjectBucketReplicationEndpoint(serviceName, bucketId) + "/trigger"
}

func storageObjectBucketReplicationJobEndpoint(serviceName, bucketId, jobId string) string {
	return storageObjectBucketReplicationJobsEndpoint(serviceName, bucketId) + "/" + url.PathEscape(jobId)
}

func (r *cloudStorageObjectBucketReplicationJobResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudStorageObjectBucketReplicationJobModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := storageObjectBucketReplicationJobsEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString())

	var responseData CloudStorageObjectBucketReplicationJobAPIResponse
	if err := r.config.OVHClient.PostWithContext(ctx, endpoint, nil, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Post %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketReplicationJobResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudStorageObjectBucketReplicationJobModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := storageObjectBucketReplicationJobEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString(), data.Id.ValueString())

	var responseData CloudStorageObjectBucketReplicationJobAPIResponse
	if err := r.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			// A job that was already read once is history: dropping it from the state would make
			// the next apply start a new job over every object of the bucket.
			if data.Status.IsNull() {
				resp.State.RemoveResource(ctx)
			}
			return
		}
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update is never reached with a plan change: every argument forces a replacement.
func (r *cloudStorageObjectBucketReplicationJobResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudStorageObjectBucketReplicationJobModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete only forgets the job: the API has no route to cancel or delete one.
func (r *cloudStorageObjectBucketReplicationJobResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
