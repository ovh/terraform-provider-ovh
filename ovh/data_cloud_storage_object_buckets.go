package ovh

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ovh/terraform-provider-ovh/v2/ovh/helpers"
	ovhtypes "github.com/ovh/terraform-provider-ovh/v2/ovh/types"
)

var _ datasource.DataSourceWithConfigure = (*cloudStorageObjectBucketsDataSource)(nil)

func NewCloudStorageObjectBucketsDataSource() datasource.DataSource {
	return &cloudStorageObjectBucketsDataSource{}
}

type cloudStorageObjectBucketsDataSource struct {
	config *Config
}

func (d *cloudStorageObjectBucketsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_storage_object_buckets"
}

func (d *cloudStorageObjectBucketsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *cloudStorageObjectBucketsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "List the S3 compatible object storage buckets in a public cloud project.",
		MarkdownDescription: "List the S3 compatible object storage buckets in a public cloud project.",
		Attributes: map[string]schema.Attribute{
			"service_name": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Required:            true,
				Description:         "Service name of the resource representing the id of the cloud project",
				MarkdownDescription: "Service name of the resource representing the id of the cloud project",
			},
			"region": schema.StringAttribute{
				CustomType:          ovhtypes.TfStringType{},
				Optional:            true,
				Description:         "If set, only buckets located in this region are returned",
				MarkdownDescription: "If set, only buckets located in this region are returned",
			},
			"buckets": schema.ListNestedAttribute{
				Computed:            true,
				Description:         "List of buckets",
				MarkdownDescription: "List of buckets",
				NestedObject: schema.NestedAttributeObject{
					Attributes: storageObjectBucketDataSourceAttributes(),
				},
			},
		},
	}
}

// cloudStorageObjectBucketsDataSourceModel is the Terraform state model for this data source.
type cloudStorageObjectBucketsDataSourceModel struct {
	ServiceName ovhtypes.TfStringValue `tfsdk:"service_name"`
	Region      ovhtypes.TfStringValue `tfsdk:"region"`
	Buckets     types.List             `tfsdk:"buckets"`
}

// storageObjectBucketListItemAttrTypes returns the attribute types for a single bucket item in the list.
func storageObjectBucketListItemAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":              ovhtypes.TfStringType{},
		"name":            ovhtypes.TfStringType{},
		"location":        types.ObjectType{AttrTypes: StorageObjectBucketLocationAttrTypes()},
		"owner_user_id":   ovhtypes.TfStringType{},
		"encryption":      types.ObjectType{AttrTypes: StorageObjectBucketEncryptionAttrTypes()},
		"versioning":      types.ObjectType{AttrTypes: StorageObjectBucketVersioningAttrTypes()},
		"object_lock":     types.ObjectType{AttrTypes: StorageObjectBucketObjectLockAttrTypes()},
		"tags":            types.MapType{ElemType: ovhtypes.TfStringType{}},
		"checksum":        ovhtypes.TfStringType{},
		"created_at":      ovhtypes.TfStringType{},
		"updated_at":      ovhtypes.TfStringType{},
		"resource_status": ovhtypes.TfStringType{},
		"current_state":   types.ObjectType{AttrTypes: StorageObjectBucketCurrentStateAttrTypes()},
	}
}

func (d *cloudStorageObjectBucketsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data cloudStorageObjectBucketsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := "/v2/publicCloud/project/" + url.PathEscape(data.ServiceName.ValueString()) + "/storage/object/bucket"

	apiBuckets, err := helpers.GetAllPagesV2[CloudStorageObjectBucketAPIResponse](ctx, d.config.OVHClient, endpoint)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error calling Get %s", endpoint),
			err.Error(),
		)
		return
	}

	regionFilter := ""
	if !data.Region.IsNull() && !data.Region.IsUnknown() {
		regionFilter = data.Region.ValueString()
	}

	bucketObjs := make([]attr.Value, 0, len(apiBuckets))
	for i := range apiBuckets {
		v := apiBuckets[i]

		region := ""
		if v.TargetSpec != nil && v.TargetSpec.Location != nil {
			region = v.TargetSpec.Location.Region
		}
		if v.CurrentState != nil && v.CurrentState.Location != nil && v.CurrentState.Location.Region != "" {
			region = v.CurrentState.Location.Region
		}
		if regionFilter != "" && !strings.EqualFold(region, regionFilter) {
			continue
		}

		// Reuse the singular mapping to keep the shape consistent.
		var item cloudStorageObjectBucketDataSourceModel
		mapStorageObjectBucketToDataSourceModel(ctx, &v, &item)

		obj, diags := types.ObjectValue(
			storageObjectBucketListItemAttrTypes(),
			map[string]attr.Value{
				"id":              item.Id,
				"name":            item.Name,
				"location":        item.Location,
				"owner_user_id":   item.OwnerUserId,
				"encryption":      item.Encryption,
				"versioning":      item.Versioning,
				"object_lock":     item.ObjectLock,
				"tags":            item.Tags,
				"checksum":        item.Checksum,
				"created_at":      item.CreatedAt,
				"updated_at":      item.UpdatedAt,
				"resource_status": item.ResourceStatus,
				"current_state":   item.CurrentState,
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		bucketObjs = append(bucketObjs, obj)
	}

	bucketsList, diags := types.ListValue(
		types.ObjectType{AttrTypes: storageObjectBucketListItemAttrTypes()},
		bucketObjs,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Buckets = bucketsList

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
