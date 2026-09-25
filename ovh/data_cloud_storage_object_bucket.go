package ovh

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageObjectBucketDataSource)(nil)

func NewCloudStorageObjectBucketDataSource() datasource.DataSource {
	return &cloudStorageObjectBucketDataSource{}
}

type cloudStorageObjectBucketDataSource struct {
	config *Config
}

func (d *cloudStorageObjectBucketDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_bucket"
}

func (d *cloudStorageObjectBucketDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func s3BucketDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Bucket identifier",
		},
		"name": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Bucket name",
		},
		"location": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Target geographic region of the bucket",
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Region identifier",
				},
			},
		},
		"owner_user_id": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Owner user identifier",
		},
		"encryption": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Server-side encryption configuration",
			Attributes: map[string]schema.Attribute{
				"algorithm": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Encryption algorithm",
				},
			},
		},
		"versioning": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Versioning configuration",
			Attributes: map[string]schema.Attribute{
				"status": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Versioning status",
				},
			},
		},
		"object_lock": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Object lock (WORM) configuration",
			Attributes: map[string]schema.Attribute{
				"mode": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Object lock retention mode",
				},
				"retention_days": schema.Int64Attribute{
					Computed:    true,
					Description: "Number of days to retain objects",
				},
				"retention_years": schema.Int64Attribute{
					Computed:    true,
					Description: "Number of years to retain objects",
				},
			},
		},
		"tags": schema.MapAttribute{
			ElementType: ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Metadata tags for the bucket",
		},
		"checksum": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Computed hash representing the current target specification value",
		},
		"created_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Creation date of the bucket",
		},
		"updated_at": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Last update date of the bucket",
		},
		"resource_status": schema.StringAttribute{
			CustomType:  ovhtypes.TfStringType{},
			Computed:    true,
			Description: "Bucket readiness in the system (CREATING, DELETING, ERROR, OUT_OF_SYNC, READY, SUSPENDED, UNKNOWN, UPDATING)",
		},
		"current_state": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Current observed state of the bucket",
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Bucket name",
				},
				"location": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "Geographic region where the bucket is located",
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Region identifier",
						},
					},
				},
				"encryption": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "Current encryption configuration",
					Attributes: map[string]schema.Attribute{
						"algorithm": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Encryption algorithm",
						},
					},
				},
				"versioning": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "Current versioning configuration",
					Attributes: map[string]schema.Attribute{
						"status": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Versioning status",
						},
					},
				},
				"object_lock": schema.SingleNestedAttribute{
					Computed:    true,
					Description: "Current object lock configuration",
					Attributes: map[string]schema.Attribute{
						"mode": schema.StringAttribute{
							CustomType:  ovhtypes.TfStringType{},
							Computed:    true,
							Description: "Object lock retention mode",
						},
						"retention_days": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of days to retain objects",
						},
						"retention_years": schema.Int64Attribute{
							Computed:    true,
							Description: "Number of years to retain objects",
						},
					},
				},
				"tags": schema.MapAttribute{
					ElementType: ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Current metadata tags",
				},
				"virtual_host": schema.StringAttribute{
					CustomType:  ovhtypes.TfStringType{},
					Computed:    true,
					Description: "Bucket virtual host, as a hostname without scheme. Only returned on a single bucket read.",
				},
				"objects_count": schema.Int64Attribute{
					Computed:    true,
					Description: "Bucket total objects count. Only returned on a single bucket read.",
				},
				"objects_size": schema.Int64Attribute{
					Computed:    true,
					Description: "Bucket total objects size in bytes. Only returned on a single bucket read.",
				},
			},
		},
	}
}

func (d *cloudStorageObjectBucketDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"service_name": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Service name of the resource representing the id of the cloud project",
			MarkdownDescription: "Service name of the resource representing the id of the cloud project",
		},
		"id": schema.StringAttribute{
			CustomType:          ovhtypes.TfStringType{},
			Required:            true,
			Description:         "Bucket identifier, in the <REGION>_<name> form",
			MarkdownDescription: "Bucket identifier, in the `<REGION>_<name>` form",
		},
	}

	for name, attribute := range s3BucketDataSourceAttributes() {
		if name == "id" {
			continue
		}
		attrs[name] = attribute
	}

	resp.Schema = schema.Schema{
		Description:         "Get an S3 compatible object storage bucket in a public cloud project.",
		MarkdownDescription: "Get an S3 compatible object storage bucket in a public cloud project.",
		Attributes:          attrs,
	}
}

type cloudS3BucketDataSourceModel struct {
	ServiceName    ovhtypes.TfStringValue `tfsdk:"service_name"`
	Id             ovhtypes.TfStringValue `tfsdk:"id"`
	Name           ovhtypes.TfStringValue `tfsdk:"name"`
	Location       types.Object           `tfsdk:"location"`
	OwnerUserId    ovhtypes.TfStringValue `tfsdk:"owner_user_id"`
	Encryption     types.Object           `tfsdk:"encryption"`
	Versioning     types.Object           `tfsdk:"versioning"`
	ObjectLock     types.Object           `tfsdk:"object_lock"`
	Tags           types.Map              `tfsdk:"tags"`
	Checksum       ovhtypes.TfStringValue `tfsdk:"checksum"`
	CreatedAt      ovhtypes.TfStringValue `tfsdk:"created_at"`
	UpdatedAt      ovhtypes.TfStringValue `tfsdk:"updated_at"`
	ResourceStatus ovhtypes.TfStringValue `tfsdk:"resource_status"`
	CurrentState   types.Object           `tfsdk:"current_state"`
}

func (d *cloudStorageObjectBucketDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudS3BucketDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) +
		"/storage/object/bucket/" + url.PathEscape(data.Id.ValueString())

	var v CloudS3BucketAPIResponse
	if err := d.config.OVHClient.GetWithContext(ctx, endpoint, &v); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	mapS3BucketToDataSourceModel(ctx, &v, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func mapS3BucketToDataSourceModel(ctx context.Context, v *CloudS3BucketAPIResponse, data *cloudS3BucketDataSourceModel) {
	data.Id = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Id)}
	data.Checksum = ovhtypes.TfStringValue{StringValue: types.StringValue(v.Checksum)}
	data.CreatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(v.CreatedAt)}
	data.UpdatedAt = ovhtypes.TfStringValue{StringValue: types.StringValue(v.UpdatedAt)}
	data.ResourceStatus = ovhtypes.TfStringValue{StringValue: types.StringValue(v.ResourceStatus)}

	if v.TargetSpec != nil {
		data.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(v.TargetSpec.Name)}
		data.Location = buildS3BucketLocationObject(v.TargetSpec.Location)
		data.OwnerUserId = nullableTfString(v.TargetSpec.OwnerUserId)
		data.Encryption = buildS3BucketEncryptionObject(v.TargetSpec.Encryption)
		data.Versioning = buildS3BucketVersioningObject(v.TargetSpec.Versioning)
		data.ObjectLock = buildS3BucketObjectLockObject(v.TargetSpec.ObjectLock)
		data.Tags = buildS3BucketTagsMap(v.TargetSpec.Tags)
	} else {
		data.Location = types.ObjectNull(S3BucketLocationAttrTypes())
		data.OwnerUserId = nullableTfString("")
		data.Encryption = types.ObjectNull(S3BucketEncryptionAttrTypes())
		data.Versioning = types.ObjectNull(S3BucketVersioningAttrTypes())
		data.ObjectLock = types.ObjectNull(S3BucketObjectLockAttrTypes())
		data.Tags = types.MapNull(ovhtypes.TfStringType{})
	}

	if v.CurrentState != nil {
		if v.CurrentState.Name != "" {
			data.Name = ovhtypes.TfStringValue{StringValue: types.StringValue(v.CurrentState.Name)}
		}
		data.CurrentState = buildS3BucketCurrentStateObject(ctx, v.CurrentState)
	} else {
		data.CurrentState = types.ObjectNull(S3BucketCurrentStateAttrTypes())
	}
}
