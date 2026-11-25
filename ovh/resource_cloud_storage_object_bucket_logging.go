package ovh

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/ovh/go-ovh/ovh"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var (
	_ resource.Resource                = (*cloudStorageObjectBucketLoggingResource)(nil)
	_ resource.ResourceWithConfigure   = (*cloudStorageObjectBucketLoggingResource)(nil)
	_ resource.ResourceWithImportState = (*cloudStorageObjectBucketLoggingResource)(nil)
)

func NewCloudStorageObjectBucketLoggingResource() resource.Resource {
	return &cloudStorageObjectBucketLoggingResource{}
}

type cloudStorageObjectBucketLoggingResource struct {
	config *Config
}

func (r *cloudStorageObjectBucketLoggingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_logging"
}

func (r *cloudStorageObjectBucketLoggingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *cloudStorageObjectBucketLoggingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the server access logging of an S3 compatible object storage bucket in a public cloud project, through the S3 API " +
			"with the S3 credentials of the provider configuration. Destroying it disables the logging of the bucket.",
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
				Description:         "Bucket identifier, in the <REGION>_<name> form",
				MarkdownDescription: "Bucket identifier, in the `<REGION>_<name>` form",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Logging configuration identifier, the bucket identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"target_bucket": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Name of the bucket receiving the logs, in the same region, with the log-delivery-write ACL",
			},
			"target_prefix": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Prefix of the keys of the log objects",
			},
		},
	}
}

func (r *cloudStorageObjectBucketLoggingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	splits := strings.Split(req.ID, "/")
	if len(splits) != 2 {
		resp.Diagnostics.AddError("Given ID is malformed", "ID must be formatted like the following: <service_name>/<bucket_id>")
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_name"), splits[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("bucket_id"), splits[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), splits[1])...)
}

// resolve returns the S3 client and where the bucket is served. The error of a missing bucket is the 404 of the publicCloud API.
func (r *cloudStorageObjectBucketLoggingResource) resolve(ctx context.Context, data *CloudStorageObjectBucketLoggingModel) (*s3Client, *storageObjectBucketS3Target, error) {
	client, err := r.config.S3Client()
	if err != nil {
		return nil, nil, err
	}

	target, err := storageObjectBucketS3TargetRead(ctx, r.config.OVHClient, data.ServiceName.ValueString(), data.BucketId.ValueString())
	if err != nil {
		return nil, nil, err
	}

	return client, target, nil
}

// apply enables the logging then reads it back, as S3 may normalize the prefix.
func (r *cloudStorageObjectBucketLoggingResource) apply(ctx context.Context, data *CloudStorageObjectBucketLoggingModel) diag.Diagnostics {
	var diags diag.Diagnostics

	client, target, err := r.resolve(ctx, data)
	if err != nil {
		diags.AddError(fmt.Sprintf("Cannot reach bucket %s on the S3 API", data.BucketId.ValueString()), err.Error())
		return diags
	}

	if err := client.PutBucketLogging(ctx, target.Host, target.Region, data.ToS3()); err != nil {
		diags.AddError(fmt.Sprintf("Error enabling the logging of bucket %s", target.Host), err.Error())
		return diags
	}

	logging, err := client.GetBucketLogging(ctx, target.Host, target.Region)
	if err != nil {
		diags.AddError(fmt.Sprintf("Error reading the logging of bucket %s", target.Host), err.Error())
		return diags
	}
	if logging == nil {
		diags.AddError(fmt.Sprintf("Error enabling the logging of bucket %s", target.Host), "S3 reports the logging disabled right after enabling it")
		return diags
	}

	data.MergeWith(logging)
	return diags
}

func (r *cloudStorageObjectBucketLoggingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudStorageObjectBucketLoggingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.apply(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketLoggingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudStorageObjectBucketLoggingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, target, err := r.resolve(ctx, &data)
	if err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			resp.State.RemoveResource(ctx)
			return
		}
		endpoint := storageObjectBucketEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString())
		if storageObjectBucketReadUnavailable(err, endpoint, &resp.Diagnostics) {
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot reach bucket %s on the S3 API", data.BucketId.ValueString()), err.Error())
		return
	}

	logging, err := client.GetBucketLogging(ctx, target.Host, target.Region)
	if err != nil {
		if isS3NoSuchBucket(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Error reading the logging of bucket %s", target.Host), err.Error())
		return
	}

	if logging == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	data.MergeWith(logging)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketLoggingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudStorageObjectBucketLoggingModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(r.apply(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *cloudStorageObjectBucketLoggingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudStorageObjectBucketLoggingModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, target, err := r.resolve(ctx, &data)
	if err != nil {
		if errOvh, ok := err.(*ovh.APIError); ok && errOvh.Code == 404 {
			return
		}
		resp.Diagnostics.AddError(fmt.Sprintf("Cannot reach bucket %s on the S3 API", data.BucketId.ValueString()), err.Error())
		return
	}

	if err := client.PutBucketLogging(ctx, target.Host, target.Region, nil); err != nil && !isS3NoSuchBucket(err) {
		resp.Diagnostics.AddError(fmt.Sprintf("Error disabling the logging of bucket %s", target.Host), err.Error())
	}
}
