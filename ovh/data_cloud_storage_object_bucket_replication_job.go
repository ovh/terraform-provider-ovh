package ovh

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageObjectBucketReplicationJobDataSource)(nil)

func NewCloudStorageObjectBucketReplicationJobDataSource() datasource.DataSource {
	return &cloudStorageObjectBucketReplicationJobDataSource{}
}

type cloudStorageObjectBucketReplicationJobDataSource struct {
	config *Config
}

func (d *cloudStorageObjectBucketReplicationJobDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_replication_job"
}

func (d *cloudStorageObjectBucketReplicationJobDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func storageObjectBucketReplicationJobDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"status":     schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Job status (ACTIVE, COMPLETED, CREATED, FAILED, PREPARING)"},
		"total":      schema.Int64Attribute{Computed: true, Description: "Total number of objects to replicate"},
		"completed":  schema.Int64Attribute{Computed: true, Description: "Number of successfully replicated objects"},
		"failed":     schema.Int64Attribute{Computed: true, Description: "Number of objects that failed to replicate"},
		"stuck":      schema.Int64Attribute{Computed: true, Description: "Number of stuck objects"},
		"created_at": schema.StringAttribute{CustomType: ovhtypes.TfStringType{}, Computed: true, Description: "Job creation date"},
	}
}

func (d *cloudStorageObjectBucketReplicationJobDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attributes := storageObjectBucketReplicationJobDataSourceAttributes()
	attributes["service_name"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Optional:            true,
		Computed:            true,
		Description:         "Service name of the resource representing the id of the cloud project. If omitted, the OVH_CLOUD_PROJECT_SERVICE environment variable is used.",
		MarkdownDescription: "Service name of the resource representing the id of the cloud project. If omitted, the `OVH_CLOUD_PROJECT_SERVICE` environment variable is used.",
	}
	attributes["bucket_id"] = schema.StringAttribute{
		CustomType:          ovhtypes.TfStringType{},
		Required:            true,
		Description:         "Source bucket identifier, in the <REGION>_<name> form",
		MarkdownDescription: "Source bucket identifier, in the `<REGION>_<name>` form",
	}
	attributes["id"] = schema.StringAttribute{
		CustomType:  ovhtypes.TfStringType{},
		Required:    true,
		Description: "Unique identifier of the job",
	}

	resp.Schema = schema.Schema{
		Description: "Get a batch replication job of an S3 compatible object storage bucket in a public cloud project.",
		Attributes:  attributes,
	}
}

func (d *cloudStorageObjectBucketReplicationJobDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudStorageObjectBucketReplicationJobModel

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

	endpoint := storageObjectBucketReplicationJobEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString(), data.Id.ValueString())

	var responseData CloudStorageObjectBucketReplicationJobAPIResponse
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	data.MergeWith(&responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
