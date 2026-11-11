package ovh

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageObjectBucketReplicationJobsDataSource)(nil)

func NewCloudStorageObjectBucketReplicationJobsDataSource() datasource.DataSource {
	return &cloudStorageObjectBucketReplicationJobsDataSource{}
}

type cloudStorageObjectBucketReplicationJobsDataSource struct {
	config *Config
}

type cloudStorageObjectBucketReplicationJobsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue                       `tfsdk:"service_name"`
	BucketId    ovhtypes.TfStringValue                       `tfsdk:"bucket_id"`
	Jobs        []storageObjectBucketReplicationJobItemModel `tfsdk:"jobs"`
}

func (d *cloudStorageObjectBucketReplicationJobsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_replication_jobs"
}

func (d *cloudStorageObjectBucketReplicationJobsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageObjectBucketReplicationJobsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	jobAttributes := storageObjectBucketReplicationJobDataSourceAttributes()
	jobAttributes["id"] = schema.StringAttribute{
		CustomType:  ovhtypes.TfStringType{},
		Computed:    true,
		Description: "Unique identifier of the job",
	}

	resp.Schema = schema.Schema{
		Description: "List the batch replication jobs of an S3 compatible object storage bucket in a public cloud project.",
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
			"jobs": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "Batch replication jobs of the bucket",
				NestedObject: schema.NestedAttributeObject{Attributes: jobAttributes},
			},
		},
	}
}

func (d *cloudStorageObjectBucketReplicationJobsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageObjectBucketReplicationJobsDataSourceModel

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

	endpoint := storageObjectBucketReplicationJobsEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString())

	jobs, err := helpers.GetAllPagesV2[CloudStorageObjectBucketReplicationJobAPIResponse](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.Jobs = make([]storageObjectBucketReplicationJobItemModel, 0, len(jobs))
	for i := range jobs {
		data.Jobs = append(data.Jobs, newStorageObjectBucketReplicationJobItem(&jobs[i]))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
