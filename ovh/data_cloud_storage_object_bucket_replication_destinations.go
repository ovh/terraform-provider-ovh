package ovh

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageObjectBucketReplicationDestinationsDataSource)(nil)

func NewCloudStorageObjectBucketReplicationDestinationsDataSource() datasource.DataSource {
	return &cloudStorageObjectBucketReplicationDestinationsDataSource{}
}

type cloudStorageObjectBucketReplicationDestinationsDataSource struct {
	config *Config
}

type cloudStorageObjectBucketReplicationDestinationsDataSourceModel struct {
	ServiceName  ovhtypes.TfStringValue                               `tfsdk:"service_name"`
	BucketId     ovhtypes.TfStringValue                               `tfsdk:"bucket_id"`
	Destinations []storageObjectBucketReplicationDestinationItemModel `tfsdk:"destinations"`
}

type CloudStorageObjectBucketReplicationDestinationAPIResponse struct {
	Eligible bool   `json:"eligible"`
	Name     string `json:"name"`
	Reason   string `json:"reason"`
	Region   string `json:"region"`
}

type storageObjectBucketReplicationDestinationItemModel struct {
	Eligible types.Bool             `tfsdk:"eligible"`
	Name     ovhtypes.TfStringValue `tfsdk:"name"`
	Reason   ovhtypes.TfStringValue `tfsdk:"reason"`
	Region   ovhtypes.TfStringValue `tfsdk:"region"`
}

func newStorageObjectBucketReplicationDestinationItem(response *CloudStorageObjectBucketReplicationDestinationAPIResponse) storageObjectBucketReplicationDestinationItemModel {
	return storageObjectBucketReplicationDestinationItemModel{
		Eligible: types.BoolValue(response.Eligible),
		Name:     ovhtypes.NewTfStringValue(response.Name),
		Reason:   nullableTfString(response.Reason),
		Region:   ovhtypes.NewTfStringValue(response.Region),
	}
}

func (d *cloudStorageObjectBucketReplicationDestinationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_replication_destinations"
}

func (d *cloudStorageObjectBucketReplicationDestinationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageObjectBucketReplicationDestinationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List the other buckets of a public cloud project and whether each one can receive the objects replicated from an S3 compatible object storage bucket.",
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
				Description:         "Source bucket identifier, in the <REGION>_<name> form",
				MarkdownDescription: "Source bucket identifier, in the `<REGION>_<name>` form",
			},
			"destinations": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Buckets of the project evaluated as replication destinations of the source bucket",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Bucket name",
						},
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region the bucket is located in",
						},
						"eligible": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether this bucket can be used as the destination of a replication rule on the source bucket",
						},
						"reason": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Why the bucket cannot be used as a destination (OBJECT_LOCK_INCOMPATIBLE, REGION_NOT_REPLICATION_TARGET, VERSIONING_DISABLED). Null when eligible is true",
						},
					},
				},
			},
		},
	}
}

func (d *cloudStorageObjectBucketReplicationDestinationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageObjectBucketReplicationDestinationsDataSourceModel

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

	endpoint := storageObjectBucketReplicationEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString()) + "/destination"

	destinations, err := helpers.GetAllPagesV2[CloudStorageObjectBucketReplicationDestinationAPIResponse](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.Destinations = make([]storageObjectBucketReplicationDestinationItemModel, 0, len(destinations))
	for i := range destinations {
		data.Destinations = append(data.Destinations, newStorageObjectBucketReplicationDestinationItem(&destinations[i]))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
