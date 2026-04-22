package ovh

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageObjectBucketObjectDataSource)(nil)

func NewCloudStorageObjectBucketObjectDataSource() datasource.DataSource {
	return &cloudStorageObjectBucketObjectDataSource{}
}

type cloudStorageObjectBucketObjectDataSource struct {
	config *Config
}

// cloudStorageObjectBucketObjectDataSourceModel is the Terraform state model
// for the single-object data source.
type cloudStorageObjectBucketObjectDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketName  ovhtypes.TfStringValue `tfsdk:"bucket_name"`
	ObjectKey   ovhtypes.TfStringValue `tfsdk:"object_key"`
	VersionId   ovhtypes.TfStringValue `tfsdk:"version_id"`

	// Computed
	ETag              ovhtypes.TfStringValue `tfsdk:"e_tag"`
	IsCommonPrefix    types.Bool             `tfsdk:"is_common_prefix"`
	IsDeleteMarker    types.Bool             `tfsdk:"is_delete_marker"`
	IsLatest          types.Bool             `tfsdk:"is_latest"`
	Key               ovhtypes.TfStringValue `tfsdk:"key"`
	LastModified      ovhtypes.TfStringValue `tfsdk:"last_modified"`
	LegalHold         ovhtypes.TfStringValue `tfsdk:"legal_hold"`
	ReplicationStatus ovhtypes.TfStringValue `tfsdk:"replication_status"`
	RestoreStatus     types.Object           `tfsdk:"restore_status"`
	Retention         types.Object           `tfsdk:"retention"`
	Size              types.Int64            `tfsdk:"size"`
	StorageClass      ovhtypes.TfStringValue `tfsdk:"storage_class"`
	ResponseVersionId ovhtypes.TfStringValue `tfsdk:"response_version_id"`
}

func (d *cloudStorageObjectBucketObjectDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_object"
}

func (d *cloudStorageObjectBucketObjectDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	config, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.config = config
}

func (d *cloudStorageObjectBucketObjectDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads a single S3 object's metadata (optionally a specific version) using the v2 API.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Service name of the resource representing the id of the cloud project",
			},
			"bucket_name": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Name of the bucket",
			},
			"object_key": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Required:    true,
				Description: "Object key (path within the bucket)",
			},
			"version_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "Version identifier. When set, the data source reads metadata for that specific version.",
			},

			// Computed — exposing every field of the S3Object API model.
			"e_tag": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Entity tag (ETag) of the object",
			},
			"is_common_prefix": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this entry is a common-prefix grouping rather than an object",
			},
			"is_delete_marker": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this version is a delete marker",
			},
			"is_latest": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this version is the latest version of the object",
			},
			"key": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Object key (path within the bucket) as returned by the API",
			},
			"last_modified": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Date and time at which the object was last modified",
			},
			"legal_hold": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Current legal hold status of the object (ON, OFF)",
			},
			"replication_status": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Replication status reported for this object (COMPLETED, FAILED, NONE, PENDING, REPLICA)",
			},
			"restore_status": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Restore status for objects stored in archive tiers",
				Attributes: map[string]schema.Attribute{
					"expire_date": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Date and time at which the restored copy expires",
					},
					"in_progress": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether a restore operation is currently in progress",
					},
				},
			},
			"retention": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Current retention (object lock) configuration",
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Retention mode applied to the object (COMPLIANCE, GOVERNANCE, NONE)",
					},
					"retain_until_date": schema.StringAttribute{
						CustomType:  ovhtypes.TfStringType{},
						Computed:    true,
						Description: "Date and time until which the object is retained",
					},
				},
			},
			"size": schema.Int64Attribute{
				Computed:    true,
				Description: "Size of the object in bytes",
			},
			"storage_class": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Storage class of the object (DEEP_ARCHIVE, GLACIER, GLACIER_IR, HIGH_PERF, INTELLIGENT_TIERING, ONEZONE_IA, STANDARD, STANDARD_IA)",
			},
			"response_version_id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Version identifier returned by the API for the object",
			},
		},
	}
}

func (d *cloudStorageObjectBucketObjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageObjectBucketObjectDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := objectEndpoint(
		data.ServiceName.ValueString(),
		data.BucketName.ValueString(),
		data.ObjectKey.ValueString(),
		data.VersionId.ValueString(),
	)

	var response CloudBucketObjectAPI
	if err := d.config.OVHClient.Get(endpoint, &response); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	// Populate computed fields.
	data.ETag = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ETag)}
	data.IsCommonPrefix = types.BoolValue(response.IsCommonPrefix)
	data.IsDeleteMarker = types.BoolValue(response.IsDeleteMarker)
	data.IsLatest = types.BoolValue(response.IsLatest)
	data.Key = ovhtypes.TfStringValue{StringValue: types.StringValue(response.Key)}
	data.LastModified = ovhtypes.TfStringValue{StringValue: types.StringValue(response.LastModified)}
	data.LegalHold = ovhtypes.TfStringValue{StringValue: types.StringValue(response.LegalHold)}
	data.ReplicationStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(response.ReplicationStatus)}
	data.RestoreStatus = apiRestoreStatusToTFObject(response.RestoreStatus)
	data.Retention = apiObjectRetentionToTFObject(response.Retention)
	data.Size = types.Int64Value(response.Size)
	data.StorageClass = ovhtypes.TfStringValue{StringValue: types.StringValue(response.StorageClass)}
	data.ResponseVersionId = ovhtypes.TfStringValue{StringValue: types.StringValue(response.VersionID)}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
