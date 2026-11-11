package ovh

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageObjectBucketReplicationDataSource)(nil)

func NewCloudStorageObjectBucketReplicationDataSource() datasource.DataSource {
	return &cloudStorageObjectBucketReplicationDataSource{}
}

type cloudStorageObjectBucketReplicationDataSource struct {
	config *Config
}

func (d *cloudStorageObjectBucketReplicationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_replication"
}

func (d *cloudStorageObjectBucketReplicationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageObjectBucketReplicationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get the replication configuration of an S3 compatible object storage bucket in a public cloud project. " +
			"A bucket without replication configuration has an empty rules list.",
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
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Replication configuration identifier, the bucket identifier",
			},
			"rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Replication rules of the bucket",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Replication rule unique identifier",
						},
						"status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Whether the rule is enabled or disabled (DISABLED, ENABLED)",
						},
						"priority": schema.Int64Attribute{
							Computed:    true,
							Description: "Priority of the replication rule",
						},
						"delete_marker_replication": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Whether delete markers are replicated (DISABLED, ENABLED)",
						},
						"destination": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Destination bucket and settings for replicated objects",
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Name of the destination bucket",
								},
								"region": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Region of the destination bucket. Null when the destination name matches no bucket of the project, or several",
								},
								"storage_class": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Storage class for replicated objects in the destination bucket (" + strings.Join(storageObjectBucketReplicationStorageClassValues, ", ") + ")",
								},
							},
						},
						"filter": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Filter criteria for the rule. Null when the rule applies to every object of the bucket",
							Attributes: map[string]schema.Attribute{
								"prefix": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Key prefix to filter objects",
								},
								"tags": schema.MapAttribute{
									ElementType: ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Key-value tag pairs to filter objects",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *cloudStorageObjectBucketReplicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudStorageObjectBucketReplicationModel

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

	endpoint := storageObjectBucketReplicationEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString())

	// A bucket without replication configuration answers 200 with an empty rule list, kept as is.
	var responseData CloudStorageObjectBucketReplicationAPIRules
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &responseData); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(data.MergeWith(ctx, &responseData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
