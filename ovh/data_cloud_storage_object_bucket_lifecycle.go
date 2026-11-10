package ovh

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageObjectBucketLifecycleDataSource)(nil)

func NewCloudStorageObjectBucketLifecycleDataSource() datasource.DataSource {
	return &cloudStorageObjectBucketLifecycleDataSource{}
}

type cloudStorageObjectBucketLifecycleDataSource struct {
	config *Config
}

func (d *cloudStorageObjectBucketLifecycleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket_lifecycle"
}

func (d *cloudStorageObjectBucketLifecycleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageObjectBucketLifecycleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get the lifecycle configuration of an S3 compatible object storage bucket in a public cloud project. " +
			"A bucket without lifecycle configuration has an empty rules list.",
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
			"id": schema.StringAttribute{
				CustomType:  ovhtypes.TfStringType{},
				Computed:    true,
				Description: "Lifecycle configuration identifier, the bucket identifier",
			},
			"rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Lifecycle rules of the bucket",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Lifecycle rule unique identifier",
						},
						"status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Whether the rule is enabled or disabled (ENABLED, DISABLED)",
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
								"object_size_greater_than": schema.Int64Attribute{
									Computed:    true,
									Description: "Minimum object size in bytes for the rule to apply",
								},
								"object_size_less_than": schema.Int64Attribute{
									Computed:    true,
									Description: "Maximum object size in bytes for the rule to apply",
								},
							},
						},
						"expiration": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Expiration settings for current version objects",
							Attributes: map[string]schema.Attribute{
								"date": schema.StringAttribute{
									CustomType:  ovhtypes.TfStringType{},
									Computed:    true,
									Description: "Date (YYYY-MM-DD) when objects expire",
								},
								"days": schema.Int64Attribute{
									Computed:    true,
									Description: "Number of days after creation when objects expire",
								},
								"expired_object_delete_marker": schema.BoolAttribute{
									Computed:    true,
									Description: "Whether to remove expired object delete markers",
								},
							},
						},
						"transitions": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Transition rules for current version objects",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"date": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Date (YYYY-MM-DD) when objects transition to another storage class",
									},
									"days": schema.Int64Attribute{
										Computed:    true,
										Description: "Number of days after creation when objects transition",
									},
									"storage_class": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Target storage class for the transition (DEEP_ARCHIVE, GLACIER_IR, STANDARD, STANDARD_IA)",
									},
								},
							},
						},
						"noncurrent_version_expiration": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Expiration settings for noncurrent version objects",
							Attributes: map[string]schema.Attribute{
								"noncurrent_days": schema.Int64Attribute{
									Computed:    true,
									Description: "Number of days after an object version becomes noncurrent before it expires",
								},
								"newer_noncurrent_versions": schema.Int64Attribute{
									Computed:    true,
									Description: "Number of newer noncurrent versions to retain",
								},
							},
						},
						"noncurrent_version_transitions": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Transition rules for noncurrent version objects",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"noncurrent_days": schema.Int64Attribute{
										Computed:    true,
										Description: "Number of days after an object version becomes noncurrent before it transitions",
									},
									"newer_noncurrent_versions": schema.Int64Attribute{
										Computed:    true,
										Description: "Number of newer noncurrent versions to retain before transitioning",
									},
									"storage_class": schema.StringAttribute{
										CustomType:  ovhtypes.TfStringType{},
										Computed:    true,
										Description: "Target storage class for the transition (DEEP_ARCHIVE, GLACIER_IR, STANDARD, STANDARD_IA)",
									},
								},
							},
						},
						"abort_incomplete_multipart_upload": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Settings for aborting incomplete multipart uploads",
							Attributes: map[string]schema.Attribute{
								"days_after_initiation": schema.Int64Attribute{
									Computed:    true,
									Description: "Number of days after initiation to abort incomplete multipart uploads",
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *cloudStorageObjectBucketLifecycleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CloudStorageObjectBucketLifecycleModel

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

	endpoint := storageObjectBucketLifecycleEndpoint(data.ServiceName.ValueString(), data.BucketId.ValueString())

	// A bucket without lifecycle configuration answers 200 with an empty rule list, kept as is.
	var responseData CloudStorageObjectBucketLifecycleAPIRules
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
