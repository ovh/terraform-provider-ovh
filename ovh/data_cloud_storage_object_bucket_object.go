package ovh

import (
	"context"
	"fmt"
	"net/url"
	"os"

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

func (d *cloudStorageObjectBucketObjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
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

func (d *cloudStorageObjectBucketObjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := storageObjectBucketObjectDataSourceAttributes()
	delete(attrs, "is_common_prefix")

	attrs["service_name"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Optional:            true,
		Computed:            true,
		Description:         "Service name of the resource representing the id of the cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used.",
		MarkdownDescription: "Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
	}
	attrs["bucket_id"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Required:            true,
		Description:         "Bucket identifier, in the <REGION>_<name> form",
		MarkdownDescription: "Bucket identifier, in the `<REGION>_<name>` form",
	}
	attrs["key"] = schema.StringAttribute{
		CustomType:  ovhtypes.TfStringType{},
		Required:    true,
		Description: "Object key (path within the bucket)",
	}

	resp.Schema = schema.Schema{
		Description:         "Get an object of an S3 compatible object storage bucket in a public cloud project.",
		MarkdownDescription: "Get an object of an S3 compatible object storage bucket in a public cloud project.",
		Attributes:          attrs,
	}
}

type cloudStorageObjectBucketObjectDataSourceModel struct {
	ServiceName       ovhtypes.TfStringValue `tfsdk:"service_name"`
	BucketId          ovhtypes.TfStringValue `tfsdk:"bucket_id"`
	Key               ovhtypes.TfStringValue `tfsdk:"key"`
	Size              types.Int64            `tfsdk:"size"`
	Etag              ovhtypes.TfStringValue `tfsdk:"etag"`
	LastModifiedAt    ovhtypes.TfStringValue `tfsdk:"last_modified_at"`
	StorageClass      ovhtypes.TfStringValue `tfsdk:"storage_class"`
	VersionId         ovhtypes.TfStringValue `tfsdk:"version_id"`
	IsLatest          types.Bool             `tfsdk:"is_latest"`
	IsDeleteMarker    types.Bool             `tfsdk:"is_delete_marker"`
	LegalHold         ovhtypes.TfStringValue `tfsdk:"legal_hold"`
	Lock              types.Object           `tfsdk:"lock"`
	RestoreStatus     types.Object           `tfsdk:"restore_status"`
	ReplicationStatus ovhtypes.TfStringValue `tfsdk:"replication_status"`
}

func (d *cloudStorageObjectBucketObjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageObjectBucketObjectDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.ServiceName.IsNull() || data.ServiceName.ValueString() == "" {
		envServiceName := os.Getenv("OVH_CLOUD_PROJECT_SERVICE")
		if envServiceName == "" {
			resp.Diagnostics.AddError(
				"Missing service_name",
				"The service_name attribute is required. Provide it in the data source "+
					"configuration or set the OVH_CLOUD_PROJECT_SERVICE environment variable.",
			)
			return
		}
		data.ServiceName = ovhtypes.NewTfStringValue(envServiceName)
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/object/bucket/" + url.PathEscape(data.BucketId.ValueString()) +
		"/object/" + url.PathEscape(data.Key.ValueString())

	var v CloudStorageObjectBucketObjectAPIResponse
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &v); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.Size = types.Int64Value(v.Size)
	data.Etag = nullableTfString(v.Etag)
	data.LastModifiedAt = nullableTfString(v.LastModifiedAt)
	data.StorageClass = nullableTfString(v.StorageClass)
	data.VersionId = nullableTfString(v.VersionId)
	data.IsLatest = types.BoolValue(v.IsLatest)
	data.IsDeleteMarker = types.BoolValue(v.IsDeleteMarker)
	data.LegalHold = nullableTfString(v.LegalHold)
	data.Lock = buildStorageObjectBucketObjectLock(v.Lock)
	data.RestoreStatus = buildStorageObjectBucketObjectRestoreStatus(v.RestoreStatus)
	data.ReplicationStatus = nullableTfString(v.ReplicationStatus)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
