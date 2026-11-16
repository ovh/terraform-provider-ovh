package ovh

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageObjectBucketMultipartUploadsDataSource)(nil)

func NewCloudStorageObjectBucketMultipartUploadsDataSource() datasource.DataSource {
	return &cloudStorageObjectBucketMultipartUploadsDataSource{}
}

type cloudStorageObjectBucketMultipartUploadsDataSource struct {
	config *Config
}

type CloudStorageObjectBucketMultipartUploadAPIResponse struct {
	Id           string `json:"id"`
	Key          string `json:"key"`
	StorageClass string `json:"storageClass"`
	CreatedAt    string `json:"createdAt,omitempty"`
}

type storageObjectBucketMultipartUploadItemModel struct {
	Id           ovhtypes.TfStringValue `tfsdk:"id"`
	Key          ovhtypes.TfStringValue `tfsdk:"key"`
	StorageClass ovhtypes.TfStringValue `tfsdk:"storage_class"`
	CreatedAt    ovhtypes.TfStringValue `tfsdk:"created_at"`
}

type cloudStorageObjectBucketMultipartUploadsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue                        `tfsdk:"service_name"`
	BucketId    ovhtypes.TfStringValue                        `tfsdk:"bucket_id"`
	Key         ovhtypes.TfStringValue                        `tfsdk:"key"`
	Uploads     []storageObjectBucketMultipartUploadItemModel `tfsdk:"uploads"`
}

func (d *cloudStorageObjectBucketMultipartUploadsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_multipart_uploads"
}

func (d *cloudStorageObjectBucketMultipartUploadsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageObjectBucketMultipartUploadsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List the incomplete multipart uploads of an S3 compatible object storage bucket in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Computed:            true,
				Description:         "Service name of the resource representing the id of the cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used.",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
			},
			"bucket_id": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Bucket identifier, in the <REGION>_<name> form",
				MarkdownDescription: "Bucket identifier, in the `<REGION>_<name>` form",
			},
			"key": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Optional:    true,
				Description: "If set, only the multipart uploads of the object with this key are returned",
			},
			"uploads": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Incomplete multipart uploads of the bucket",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Unique identifier of the multipart upload",
						},
						"key": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Object key (path within the bucket)",
						},
						"storage_class": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Storage class of the object",
						},
						"created_at": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Multipart upload initiation timestamp",
						},
					},
				},
			},
		},
	}
}

func (d *cloudStorageObjectBucketMultipartUploadsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageObjectBucketMultipartUploadsDataSourceModel

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
		"/storage/object/bucket/" + url.PathEscape(data.BucketId.ValueString())
	if !data.Key.IsNull() && !data.Key.IsUnknown() {
		endpoint += "/object/" + url.PathEscape(data.Key.ValueString())
	}
	endpoint += "/multipartUpload"

	uploads, err := helpers.GetAllPagesV2[CloudStorageObjectBucketMultipartUploadAPIResponse](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.Uploads = make([]storageObjectBucketMultipartUploadItemModel, 0, len(uploads))
	for _, upload := range uploads {
		data.Uploads = append(data.Uploads, storageObjectBucketMultipartUploadItemModel{
			Id:           ovhtypes.NewTfStringValue(upload.Id),
			Key:          ovhtypes.NewTfStringValue(upload.Key),
			StorageClass: ovhtypes.NewTfStringValue(upload.StorageClass),
			CreatedAt:    nullableTfString(upload.CreatedAt),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
